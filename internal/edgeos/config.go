// Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.
package edgeos

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"

	"github.com/britannic/blacklist/internal/regx"
)

// tree is a map of top node Objects
type tree map[string]*source

// ConfLoader interface handles multiple configuration load methods
type ConfLoader interface {
	read() io.Reader
}

// Config is a struct of configuration fields
type Config struct {
	*Env
	tree
}

// list is a struct map of entry with a RW Mutex
type ctr struct {
	*sync.RWMutex
	stat
}
type stat map[string]*stats

type stats struct {
	dropped   int32
	extracted int32
	kept      int32
}

const (
	agent     = `curl/7.26.0`
	all       = "all"
	blackhole = "dns-redirect-ip"
	disabled  = "disabled"
	domains   = "domains"
	files     = "file"
	hosts     = "hosts"
	notknown  = "unknown"
	preNoun   = "pre-configured"
	roots     = "roots"
	rootNode  = "blacklist"
	src       = "source"
	urls      = "url"

	// ExcDomns is a string labels for domain exclusions
	ExcDomns = "whitelisted-subdomains"
	// ExcHosts is a string labels for host exclusions
	ExcHosts = "whitelisted-servers"
	// ExcRoots is a string labels for preconfigured global domain exclusions
	ExcRoots = "global-whitelisted-domains"
	// PreDomns is a string label for preconfigured whitelisted domains
	PreDomns = "blacklisted-subdomains"
	// PreHosts is a string label for preconfigured blacklisted hosts
	PreHosts = "blacklisted-servers"
	// PreRoots is a string label for preconfigured global blacklisted hosts
	PreRoots = "global-blacklisted-domains"
	// False is a string constant
	False = "false"
	// True is a string constant
	True = "true"
)

func (c *Config) nodeExists(n string) bool {
	_, ok := c.tree[n]
	return ok
}

func (c *Config) addExc(n string) *Objects {
	var (
		exc   = []string{}
		iface IFace
		ltype string
	)

	switch n {
	case domains:
		ltype = ExcDomns
		iface = ExDmObj
	case hosts:
		ltype = ExcHosts
		iface = ExHtObj
	case rootNode:
		ltype = ExcRoots
		iface = ExRtObj
	}

	if c.nodeExists(n) {
		exc = c.tree[n].exc
	}

	return &Objects{
		Env:   c.Env,
		iface: iface,
		src: []*source{
			{
				Env:   c.Env,
				desc:  getLtypeDesc(iface.String()),
				exc:   exc,
				ip:    c.tree.getIP(n),
				ltype: ltype,
				nType: getType(ltype).(ntype),
				name:  ltype,
			},
		},
	}
}

func (c *Config) addInc(n string) *source {
	var (
		iface IFace
		inc   = []string{}
		lt    string
		nt    ntype
	)

	if c.nodeExists(n) && len(c.tree[n].inc) > 0 {
		inc = c.tree[n].inc
	}

	switch n {
	case domains:
		lt = getType(preDomn).(string)
		nt = getType(lt).(ntype)
		iface = PreDObj
	case hosts:
		lt = getType(preHost).(string)
		nt = getType(lt).(ntype)
		iface = PreHObj
	case rootNode:
		lt = getType(preRoot).(string)
		nt = getType(lt).(ntype)
		iface = PreRObj
	}

	return &source{
		Env:   c.Env,
		desc:  getLtypeDesc(lt),
		inc:   inc,
		iface: iface,
		ip:    c.tree.getIP(n),
		ltype: lt,
		nType: nt,
		name:  lt,
	}
}

// GetTotalStats displays aggregate statistics for processed sources
func (c *Config) GetTotalStats() (dropped, extracted, kept int32) {
	ctr := c.ctr.stat
	for k := range ctr {
		if ctr[k].kept+ctr[k].dropped != 0 {
			dropped += ctr[k].dropped
			extracted += ctr[k].extracted
			kept += ctr[k].kept
		}
	}

	if kept+dropped != 0 {
		c.Log.Noticef("Total entries found: %d", extracted)
		c.Log.Noticef("Total entries extracted %d", kept)
		c.Log.Noticef("Total entries dropped %d", dropped)
	}
	return dropped, extracted, kept
}

// NewContent returns a Contenter interface of the requested IFace type
func (c *Config) NewContent(iface IFace) (Contenter, error) {
	switch iface {
	case ExDmObj:
		return &ExcDomnObjects{Objects: c.addExc(domains)}, nil
	case ExHtObj:
		return &ExcHostObjects{Objects: c.addExc(hosts)}, nil
	case ExRtObj:
		return &ExcRootObjects{Objects: c.addExc(rootNode)}, nil
	case FileObj:
		return &FIODataObjects{Objects: c.GetAll(iface.String())}, nil
	case PreDObj:
		return &PreDomnObjects{Objects: c.GetAll(iface.String())}, nil
	case PreRObj:
		return &PreRootObjects{Objects: c.GetAll(iface.String())}, nil
	case PreHObj:
		return &PreHostObjects{Objects: c.GetAll(iface.String())}, nil
	case URLdObj:
		return &URLDomnObjects{Objects: c.Get(domains).Filter(urls)}, nil
	case URLhObj:
		return &URLHostObjects{Objects: c.Get(hosts).Filter(urls)}, nil
	}
	return nil, errors.New("invalid interface requested")
}

// excludes returns a string array of excludes
func (c *Config) excludes(nx ...string) list {
	var exc [][]byte
	switch nx {
	case nil:
		for _, k := range c.sortKeys() {
			for _, v := range c.tree[k].exc {
				exc = append(exc, []byte(v))
			}

		}
	default:
		for _, n := range nx {
			for _, v := range c.tree[n].exc {
				exc = append(exc, []byte(v))
			}
		}
	}
	return updateEntry(exc)
}

// Get returns an *Object for a given node
func (c *Config) Get(nx string) *Objects {
	o := &Objects{Env: c.Env, src: []*source{}}
	switch nx {
	case all:
		for _, n := range c.sortKeys() {
			o.addObj(c, n)
		}
	default:
		o.addObj(c, nx)
	}
	return o
}

// GetAll returns a pointer to an Objects struct
func (c *Config) GetAll(ltypes ...string) *Objects {
	o := &Objects{Env: c.Env}
	for _, n := range c.sortKeys() {
		o.objects(c, n, ltypes...)
	}
	return o
}

// InSession returns true if VyOS/EdgeOS configure is in session
func (c *Config) InSession() bool {
	return os.ExpandEnv("$_OFR_CONFIGURE") == "ok"
}

func (c tree) keyExists(k string) bool {
	_, ok := c[k]
	return ok
}

// load reads the config using the EdgeOS/VyOS cli-shell-api
func (c *Config) load(act, lvl string) ([]byte, error) {
	cmd := exec.Command(c.Bash)
	s := fmt.Sprintf(
		"%v %v %v --show-working-only", c.API, apiCMD(act, c.InSession()), lvl,
	)
	cmd.Stdin = strings.NewReader(s)
	c.Debug(fmt.Sprintf("Running shell command: %v", s))
	return cmd.Output()
}

// Nodes returns an array of configured nodes
func (c *Config) Nodes() (n []string) {
	for k := range c.tree {
		n = append(n, k)
	}
	sort.Strings(n)
	return n
}

// isTnode returns true if node is a root or top node in the blacklist configuration
func isTnode(n string) bool {
	switch n {
	case rootNode, domains, hosts:
		return true
	}
	return false
}

func (c *Config) excinc(t [][]byte, n string) {
	switch string(t[1]) {
	case "exclude":
		if isTnode(n) {
			c.Debug("Whitelisting %s on node %s", string(t[2]), n)
			c.tree[n].exc = append(c.tree[n].exc, string(t[2]))
		}
	case "include":
		if isTnode(n) {
			c.Debug("Blacklisting %s on node %s", string(t[2]), n)
			c.tree[n].inc = append(c.tree[n].inc, string(t[2]))
		}
	}
}

func (c *Config) label(name [][]byte, o *source, n string) {
	switch string(name[1]) {
	case "description":
		o.desc = string(name[2])
	case blackhole:
		o.ip = string(name[2])
	case files:
		o.file = string(name[2])
		o.ltype = string(name[1])
		c.tree[n].src = append(c.tree[n].src, o)
	case "prefix":
		o.prefix = string(name[2])
	case urls:
		o.ltype = string(name[1])
		o.url = string(name[2])
		c.tree[n].src = append(c.tree[n].src, o)
	}
}

func (c *Config) addTnodeSource(n string) {
	if isTnode(n) {
		c.tree[n] = newSource()
		c.tree[n].name = n
		c.tree[n].nType = getType(n).(ntype)
	}
}

func (c *Config) disable(line []byte, n string, find *regx.OBJ) {
	if isTnode(n) {
		c.tree[n].disabled = strToBool(string(find.SubMatch(regx.DSBL, line)[1]))
		c.Env.Disabled = c.tree[n].disabled
	}
}

func (c *Config) redirect(line []byte, n string, find *regx.OBJ) {
	if isTnode(n) {
		c.tree[n].ip = string(find.SubMatch(regx.IPBH, line)[1])
	}
}

func (c *Config) sourcename(o *source, line []byte, n string, find *regx.OBJ) {
	if isTnode(n) {
		name := find.SubMatch(regx.NAME, line)
		if o != nil {
			c.label(name, o, n)
		}
	}
}

// ProcessContent processes the Contents array
func (c *Config) ProcessContent(cts ...Contenter) error {
	var (
		errs []string
		wg   sync.WaitGroup
	)

	if len(cts) < 1 {
		return errors.New("empty Contenter interface{} passed to ProcessContent()")
	}

	for _, ct := range cts {
		for _, s := range ct.GetList().src {
			if s.err != nil {
				errs = append(errs, s.err.Error())
			}
			wg.Add(1)

			go func(s *source) {
				s.ctr.Lock()
				s.ctr.stat[typeInt(s.nType)] = &stats{}
				s.ctr.Unlock()

				if err := s.process().writeFile(); err != nil {
					errs = append(errs, err.Error())
				}
				wg.Done()
			}(s)
		}
	}

	wg.Wait()

	if errs != nil {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

// Blacklist extracts blacklist nodes from a EdgeOS/VyOS configuration structure
func (c *Config) Blacklist(r ConfLoader) error {
	var (
		b     = bufio.NewScanner(r.read())
		nodes []string
		o     *source
		find  = regx.NewRegex()
		tnode string
	)

	for b.Scan() {
		line := bytes.TrimSpace(b.Bytes())
		c.Debug(fmt.Sprintf("%s\n", string(line)))
		switch {
		case find.RX[regx.MLTI].Match(line): // add include/exclude
			c.Debug(fmt.Sprintf("Adding incExc to %s: %s\n", tnode, string(line)))
			c.excinc(find.SubMatch(regx.MLTI, line), tnode)
		case find.RX[regx.NODE].Match(line): // add tnode
			tnode = string(find.SubMatch(regx.NODE, line)[1])
			nodes = append(nodes, tnode)
			c.Debug(fmt.Sprintf("Adding %s node: %s\n", tnode, string(line)))
			c.addTnodeSource(tnode)
		case find.RX[regx.LEAF].Match(line): // add source to root/domains/hosts
			c.Debug(fmt.Sprintf("Adding leaf to %s: %s\n", tnode, string(line)))
			srcName := find.SubMatch(regx.LEAF, line)
			nodes = append(nodes, string(srcName[1]))
			o = newSource()
			o.addSource(srcName, tnode)
		case find.RX[regx.DSBL].Match(line): // add disable blacklist flag
			c.Debug(fmt.Sprintf("Adding disable flag to %s: %s\n", tnode, string(line)))
			c.disable(line, tnode, find)
		case find.RX[regx.IPBH].Match(line) && isntSource(nodes): // add blackhole IP
			c.Debug(fmt.Sprintf("Adding blackhole IP to %s: %s\n", tnode, string(line)))
			c.redirect(line, tnode, find)
		case find.RX[regx.NAME].Match(line): // add source name
			c.Debug(fmt.Sprintf("Adding source to %s: %s\n", tnode, string(line)))
			c.sourcename(o, line, tnode, find)
		case find.RX[regx.RBRC].Match(line): // found closing bracket
			if len(nodes) > 1 {
				c.Debug(fmt.Sprintf("Matching closing bracket: %s\n", string(line)))
				nodes = nodes[:len(nodes)-1] // pop last node
				tnode = nodes[len(nodes)-1]  // capture top node
			}
		}
	}

	if len(c.tree) < 1 {
		return errors.New("no blacklist configuration has been detected")
	}

	c.Debug(fmt.Sprintf("Using router configuration %v", c.String()))

	return nil
}

// ReloadDNS reloads the dnsmasq configuration
func (c *Config) ReloadDNS() ([]byte, error) {
	cmd := exec.Command(c.Bash)
	cmd.Stdin = strings.NewReader(c.DNSsvc)
	return cmd.CombinedOutput()
}

// sortKeys returns a slice of keys in lexicographical sorted order.
func (c *Config) sortKeys() (pkeys sort.StringSlice) {
	pkeys = make(sort.StringSlice, len(c.tree))
	i := 0
	for k := range c.tree {
		pkeys[i] = k
		i++
	}
	pkeys.Sort()
	return pkeys
}

// String returns pretty print for the Blacklist struct
func (c *Config) String() (s string) {
	indent := 1
	ø := comma
	cnt := len(c.sortKeys())
	s += fmt.Sprintf("{\n%v%q: [{\n", tabs(indent), "nodes")

	for i, pkey := range c.sortKeys() {
		if i == cnt-1 {
			ø = null
		}

		indent++
		s += fmt.Sprintf("%v%q: {\n", tabs(indent), pkey)
		indent++

		s += fmt.Sprintf("%v%q: %q,\n", tabs(indent), disabled, booltoStr(c.tree[pkey].disabled))
		s = is(indent, s, "ip", c.tree[pkey].ip)
		s += getJSONArray(&cfgJSON{array: c.tree[pkey].exc, pk: pkey, leaf: "excludes", indent: indent})
		s += getJSONArray(&cfgJSON{array: c.tree[pkey].inc, pk: pkey, leaf: "includes", indent: indent})
		s += getJSONsrcArray(&cfgJSON{Config: c, pk: pkey, indent: indent})

		indent--
		s += fmt.Sprintf("%v}%v\n", tabs(indent), ø)
		indent--
	}

	s += tabs(indent) + "}]\n}"
	return s
}

func (c tree) getIP(node string) string {
	if c.keyExists(node) && c[node].ip != "" {
		return c[node].ip
	}
	if c.keyExists(rootNode) {
		return c[rootNode].ip
	}
	return "0.0.0.0"
}

func (c tree) validate(node string) *Objects {
	if c.keyExists(node) {
		for _, o := range c[node].src {
			if o.ip == "" {
				o.ip = c.getIP(node)
			}
		}
		return &c[node].Objects
	}
	return &Objects{}
}

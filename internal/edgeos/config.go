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

	"github.com/britannic/blacklist/internal/regx"
)

// tree is a map of top nodes
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

type ctr map[string]*stats

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

func (c *Config) nodeExists(node string) bool {
	if _, ok := c.tree[node]; ok {
		return ok
	}
	return false
}

func (c *Config) addExc(node string) *Objects {
	var (
		exc   = []string{}
		ltype string
	)

	switch node {
	case domains:
		ltype = ExcDomns
	case hosts:
		ltype = ExcHosts
	case rootNode:
		ltype = ExcRoots
	}

	if c.nodeExists(node) {
		exc = c.tree[node].exc
	}

	return &Objects{
		Env: c.Env,
		src: []*source{
			{
				Env:   c.Env,
				desc:  getLtypeDesc(ltype),
				exc:   exc,
				ip:    c.tree.getIP(node),
				ltype: ltype,
				nType: getType(ltype).(ntype),
				name:  ltype,
			},
		},
	}
}

func (c *Config) addInc(node string) *source {
	var (
		inc   = []string{}
		ltype string
		n     ntype
	)

	if c.nodeExists(node) {
		if len(c.tree[node].inc) > 0 {
			inc = c.tree[node].inc
		}
	}

	switch node {
	case domains:
		ltype = getType(preDomn).(string)
		n = getType(ltype).(ntype)
	case hosts:
		ltype = getType(preHost).(string)
		n = getType(ltype).(ntype)
	case rootNode:
		ltype = getType(preRoot).(string)
		n = getType(ltype).(ntype)
	}

	return &source{
		Env:   c.Env,
		desc:  getLtypeDesc(ltype),
		inc:   inc,
		ip:    c.tree.getIP(node),
		ltype: ltype,
		nType: n,
		name:  ltype,
	}
}

// GetTotalStats displays aggregate statistics for processed sources
func (c *Config) GetTotalStats() (dropped, extracted, kept int32) {
	var keys []string
	for k := range c.ctr {
		keys = append(keys, k)
	}

	for _, k := range keys {
		if c.ctr[k].kept+c.ctr[k].dropped != 0 {
			dropped += c.ctr[k].dropped
			extracted += c.ctr[k].extracted
			kept += c.ctr[k].kept
		}
	}

	if kept+dropped != 0 {
		c.Log.Noticef("Total entries found: %d", extracted)
		c.Log.Noticef("Total entries extracted %d", kept)
		c.Log.Noticef("Total entries dropped %d", dropped)
	}
	return dropped, extracted, kept
}

// NewContent returns an interface of the requested IFace type
func (c *Config) NewContent(iface IFace) (Contenter, error) {
	var ltype = iface.String()
	switch iface {
	case ExDmObj:
		return &ExcDomnObjects{Objects: c.addExc(domains)}, nil
	case ExHtObj:
		return &ExcHostObjects{Objects: c.addExc(hosts)}, nil
	case ExRtObj:
		return &ExcRootObjects{Objects: c.addExc(rootNode)}, nil
	case FileObj:
		return &FIODataObjects{Objects: c.GetAll(ltype)}, nil
	case PreDObj:
		return &PreDomnObjects{Objects: c.GetAll(ltype)}, nil
	case PreRObj:
		return &PreRootObjects{Objects: c.GetAll(ltype)}, nil
	case PreHObj:
		return &PreHostObjects{Objects: c.GetAll(ltype)}, nil
	case URLdObj:
		return &URLDomnObjects{Objects: c.Get(domains).Filter(urls)}, nil
	case URLhObj:
		return &URLHostObjects{Objects: c.Get(hosts).Filter(urls)}, nil
	}
	return nil, errors.New("Invalid interface requested")
}

// excludes returns a string array of excludes
func (c *Config) excludes(nodes ...string) list {
	var exc []string
	switch nodes {
	case nil:
		for _, k := range c.sortKeys() {
			if len(c.tree[k].exc) != 0 {
				exc = append(exc, c.tree[k].exc...)
			}
		}
	default:
		for _, node := range nodes {
			exc = append(exc, c.tree[node].exc...)
		}
	}

	entries := make([][]byte, len(exc))
	for i, v := range exc {
		entries[i] = []byte(v)
	}

	return updateEntry(entries)
}

// Get returns an *Object for a given node
func (c *Config) Get(node string) *Objects {
	o := &Objects{Env: c.Env, src: []*source{}}
	switch node {
	case all:
		for _, node := range c.sortKeys() {
			o.addObj(c, node)
		}
	default:
		o.addObj(c, node)
	}
	return o
}

// GetAll returns an array of Objects
func (c *Config) GetAll(ltypes ...string) *Objects {
	o := &Objects{Env: c.Env}
	for _, node := range c.sortKeys() {
		o.objects(c, node, ltypes...)
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
func (c *Config) Nodes() (nodes []string) {
	for k := range c.tree {
		nodes = append(nodes, k)
	}
	sort.Strings(nodes)
	return nodes
}

// isTnode returns true if tnode is part of the blacklist configuration
func isTnode(tnode string) bool {
	switch tnode {
	case rootNode, domains, hosts:
		return true
	}
	return false
}

func (c *Config) excinc(t [][]byte, tnode string) {
	switch string(t[1]) {
	case "exclude":
		if isTnode(tnode) {
			c.Debug("Whitelisting %s on node %s", string(t[2]), tnode)
			c.tree[tnode].exc = append(c.tree[tnode].exc, string(t[2]))
		}
	case "include":
		if isTnode(tnode) {
			c.Debug("Blacklisting %s on node %s", string(t[2]), tnode)
			c.tree[tnode].inc = append(c.tree[tnode].inc, string(t[2]))
		}
	}
}

func (c *Config) label(name [][]byte, o *source, tnode string) {
	switch string(name[1]) {
	case "description":
		o.desc = string(name[2])
	case blackhole:
		o.ip = string(name[2])
	case files:
		o.file = string(name[2])
		o.ltype = string(name[1])
		c.tree[tnode].src = append(c.tree[tnode].src, o)

	case "prefix":
		o.prefix = string(name[2])
	case urls:
		o.ltype = string(name[1])
		o.url = string(name[2])
		c.tree[tnode].src = append(c.tree[tnode].src, o)
	}
}

func (c *Config) addTnodeSource(tnode string) {
	if isTnode(tnode) {
		c.tree[tnode] = newSource()
		c.tree[tnode].name = tnode
		c.tree[tnode].nType = getType(tnode).(ntype)
	}
}

func (c *Config) disable(line []byte, tnode string, find *regx.OBJ) {
	if isTnode(tnode) {
		c.tree[tnode].disabled = strToBool(string(find.SubMatch(regx.DSBL, line)[1]))
		c.Env.Disabled = c.tree[tnode].disabled
	}
}

func (c *Config) redirect(line []byte, tnode string, find *regx.OBJ) {
	if isTnode(tnode) {
		c.tree[tnode].ip = string(find.SubMatch(regx.IPBH, line)[1])
	}
}

func (c *Config) sourcename(o *source, line []byte, tnode string, find *regx.OBJ) {
	if isTnode(tnode) {
		name := find.SubMatch(regx.NAME, line)
		if o != nil {
			c.label(name, o, tnode)
		}
	}
}

// ProcessContent processes the Contents array
func (c *Config) ProcessContent(cts ...Contenter) error {
	var (
		errs      []string
		getErrors chan error
	)

	if len(cts) < 1 {
		return errors.New("Empty Contenter interface{} passed to ProcessContent()")
	}

	for _, ct := range cts {
		var (
			a, b  int32
			area  string
			tally = &stats{dropped: a, kept: b}
		)

		for _, s := range ct.GetList().src {
			getErrors = make(chan error)

			if s.err != nil {
				errs = append(errs, s.err.Error())
			}

			go func(s *source) {
				area = typeInt(s.nType)
				c.ctr[area] = tally
				getErrors <- s.process().writeFile()
			}(s)

			for range cts {
				if err := <-getErrors; err != nil {
					errs = append(errs, err.Error())
				}
				close(getErrors)
			}
		}

		if area != "" {
			if c.ctr[area].kept+c.ctr[area].dropped != 0 {
				c.Log.Noticef("Total %s found: %d", area, c.ctr[area].extracted)
				c.Log.Noticef("Total %s extracted %d", area, c.ctr[area].kept)
				c.Log.Noticef("Total %s dropped %d", area, c.ctr[area].dropped)
			}
		}
	}

	if errs != nil {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}

	return nil
}

// ReadCfg extracts nodes from a EdgeOS/VyOS configuration structure
func (c *Config) ReadCfg(r ConfLoader) error {
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
			incExc := find.SubMatch(regx.MLTI, line)
			c.excinc(incExc, tnode)
		case find.RX[regx.NODE].Match(line): // add node
			node := find.SubMatch(regx.NODE, line)
			tnode = string(node[1])
			nodes = append(nodes, tnode)
			c.Debug(fmt.Sprintf("Adding %s node: %s\n", tnode, string(line)))
			c.addTnodeSource(tnode)
		case find.RX[regx.LEAF].Match(line): // add leaf node to root/domains/hosts
			c.Debug(fmt.Sprintf("Adding leaf to %s: %s\n", tnode, string(line)))
			srcName := find.SubMatch(regx.LEAF, line)
			nodes = append(nodes, string(srcName[1]))
			o = newSource()
			o.addLeaf(srcName, tnode)
		case find.RX[regx.DSBL].Match(line): // add disable blacklist flag
			c.Debug(fmt.Sprintf("Adding disable flag to %s: %s\n", tnode, string(line)))
			c.disable(line, tnode, find)
		case find.RX[regx.IPBH].Match(line) && isntSource(nodes): // add blackhole IP
			c.Debug(fmt.Sprintf("Adding blackhole IP to %s: %s\n", tnode, string(line)))
			c.redirect(line, tnode, find)
		case find.RX[regx.NAME].Match(line): // add source name
			c.Debug(fmt.Sprintf("Adding source to %s: %s\n", tnode, string(line)))
			c.sourcename(o, line, tnode, find)
		case find.RX[regx.RBRC].Match(line):
			if len(nodes) > 1 {
				c.Debug(fmt.Sprintf("Matching closing bracket: %s\n", string(line)))
				nodes = nodes[:len(nodes)-1] // pop last node
				tnode = nodes[len(nodes)-1]
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
	if c.keyExists(node) {
		if c[node].ip != "" {
			return c[node].ip
		}
		if c.keyExists(rootNode) {
			return c[rootNode].ip
		}
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

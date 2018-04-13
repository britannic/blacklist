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
	*Parms
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
	rootNode  = "blacklist"
	src       = "source"
	urls      = "url"

	// ExcDomns labels domain exclusions
	ExcDomns = "whitelisted-subdomains"
	// ExcHosts labels host exclusions
	ExcHosts = "whitelisted-servers"
	// ExcRoots labels global domain exclusions
	ExcRoots = "whitelisted-global"
	// PreDomns designates string label for preconfigured whitelisted domains
	PreDomns = "blacklisted-subdomains"
	// PreHosts designates string label for preconfigured blacklisted hosts
	PreHosts = "blacklisted-servers"
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
		Parms: c.Parms,
		src: []*source{
			{
				Parms: c.Parms,
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
		inc = c.tree[node].inc
	}

	switch node {
	case domains:
		ltype = getType(preDomn).(string)
		n = getType(ltype).(ntype)
	case hosts:
		ltype = getType(preHost).(string)
		n = getType(ltype).(ntype)
	}

	return &source{
		Parms: c.Parms,
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
	var (
		err   error
		ltype = iface.String()
		o     *Objects
	)

	switch ltype {
	case ExcDomns:
		o = c.addExc(domains)
	case ExcHosts:
		o = c.addExc(hosts)
	case ExcRoots:
		o = c.addExc(rootNode)
	case urls:
		switch iface {
		case URLdObj:
			return &URLDomnObjects{Objects: c.Get(domains).Filter(urls)}, nil
		case URLhObj:
			return &URLHostObjects{Objects: c.Get(hosts).Filter(urls)}, nil
		}
	case notknown:
		err = errors.New("Invalid interface requested")
	default:
		o = c.GetAll(ltype)
	}

	switch iface {
	case ExDmObj:
		return &ExcDomnObjects{Objects: o}, nil
	case ExHtObj:
		return &ExcHostObjects{Objects: o}, nil
	case ExRtObj:
		return &ExcRootObjects{Objects: o}, nil
	case FileObj:
		return &FIODataObjects{Objects: o}, nil
	case PreDObj:
		return &PreDomnObjects{Objects: o}, nil
	case PreHObj:
		return &PreHostObjects{Objects: o}, nil
	}

	return nil, err
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
	o := &Objects{Parms: c.Parms, src: []*source{}}
	switch node {
	case all:
		for _, node := range c.sortKeys() {
			if node == rootNode {
				continue
			}
			o.addObj(c, node)
		}
	default:
		o.addObj(c, node)
	}
	return o
}

// GetAll returns an array of Objects
func (c *Config) GetAll(ltypes ...string) *Objects {
	o := &Objects{Parms: c.Parms}

	for _, node := range c.sortKeys() {
		if node == rootNode {
			continue
		}
		o.objects(c, node, ltypes...)
	}
	return o
}

// InSession returns true if VyOS/EdgeOS configure is in session
func (c *Config) InSession() bool {
	return os.ExpandEnv("$_OFR_CONFIGURE") == "ok"
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
			c.tree[tnode].exc = append(c.tree[tnode].exc, string(t[2]))
		}
	case "include":
		if isTnode(tnode) {
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
		c.tree[tnode].Objects.src = append(c.tree[tnode].Objects.src, o)
	case "prefix":
		o.prefix = string(name[2])
	case urls:
		o.ltype = string(name[1])
		o.url = string(name[2])
		c.tree[tnode].Objects.src = append(c.tree[tnode].Objects.src, o)
	}
}

func (c *Config) addSource(tnode string) {
	if isTnode(tnode) {
		c.tree[tnode] = newSource()
	}
}

func (c *Config) disable(line []byte, tnode string, find *regx.OBJ) {
	if isTnode(tnode) {
		c.tree[tnode].disabled = strToBool(string(find.SubMatch(regx.DSBL, line)[1]))
		c.Parms.Disabled = c.tree[tnode].disabled
	}
}

func (c *Config) redirect(line []byte, tnode string, find *regx.OBJ) {
	if isTnode(tnode) {
		c.tree[tnode].ip = string(find.SubMatch(regx.IPBH, line)[1])
	}
}

func (c *Config) leafname(o *source, line []byte, tnode string, find *regx.OBJ) {
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

		for _, o := range ct.GetList().src {
			getErrors = make(chan error)

			if o.err != nil {
				errs = append(errs, o.err.Error())
			}

			go func(o *source) {
				area = typeInt(o.nType)
				c.ctr[area] = tally
				getErrors <- o.process().writeFile()
			}(o)

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

		switch {
		case find.RX[regx.MLTI].Match(line):
			incExc := find.SubMatch(regx.MLTI, line)
			c.excinc(incExc, tnode)
		case find.RX[regx.NODE].Match(line):
			node := find.SubMatch(regx.NODE, line)
			tnode = string(node[1])
			nodes = append(nodes, tnode)
			c.addSource(tnode)
		case find.RX[regx.LEAF].Match(line):
			srcName := find.SubMatch(regx.LEAF, line)
			nodes = append(nodes, string(srcName[1]))
			o = newSource()
			o.addLeaf(srcName, tnode)
		case find.RX[regx.DSBL].Match(line):
			c.disable(line, tnode, find)
		case find.RX[regx.IPBH].Match(line) && isntSource(nodes):
			c.redirect(line, tnode, find)
		case find.RX[regx.NAME].Match(line):
			c.leafname(o, line, tnode, find)
		// case find.RX[regx.DESC].Match(line), find.RX[regx.CMNT].Match(line), find.RX[regx.MISC].Match(line):
		// 	continue
		case find.RX[regx.RBRC].Match(line):
			if len(nodes) > 1 {
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

		s += fmt.Sprintf("%v%q: %q,\n", tabs(indent), disabled,
			booltoStr(c.tree[pkey].disabled))
		s = is(indent, s, "ip", c.tree[pkey].ip)
		s += getJSONArray(&cfgJSON{array: c.tree[pkey].exc, pk: pkey, leaf: "excludes", indent: indent})

		if pkey != rootNode {
			s += getJSONArray(&cfgJSON{array: c.tree[pkey].inc, pk: pkey, leaf: "includes", indent: indent})
			s += getJSONsrcArray(&cfgJSON{Config: c, pk: pkey, indent: indent})
		}

		indent--
		s += fmt.Sprintf("%v}%v\n", tabs(indent), ø)
		indent--
	}

	s += tabs(indent) + "}]\n}"
	return s
}

func (b tree) getIP(node string) string {
	if _, ok := b[node]; ok {
		if b[node].ip != "" {
			return b[node].ip
		}
		if _, ok := b[rootNode]; ok {
			return b[rootNode].ip
		}
	}
	return "0.0.0.0"
}

func (b tree) validate(node string) *Objects {
	if _, ok := b[node]; ok {
		for _, o := range b[node].Objects.src {
			if o.ip == "" {
				o.ip = b.getIP(node)
			}
		}
		return &b[node].Objects
	}
	return &Objects{}
}

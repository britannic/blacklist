// Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.
package edgeos

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/britannic/blacklist/internal/regx"
)

type ntype int

// Types determine load order and processing behavior for blacklist sources
const (
	unknown ntype = iota // denotes a coding error
	pre                  // Pre-configured backlisted domains/hosts
	domain               // Format e.g. address=/.d.com/0.0.0.0
	host                 // Format e.g. address=/www.d.com/0.0.0.0
	root                 // Topmost root node
	zone                 // Unused - future application
)

const (
	agent     = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/601.4.4 (KHTML, like Gecko) Version/9.0.3 Safari/601.4.4`
	all       = "all"
	notknown  = "unknown"
	blackhole = "dns-redirect-ip"
	blacklist = "blacklist"
	dbg       = false
	disabled  = "disabled"
	domains   = "domains"
	src       = "source"
	hosts     = "hosts"
	preConf   = "pre-configured"
	rootNode  = blacklist
	zones     = "zones"

	// False is a string constant
	False = "false"
	// True is a string constant
	True = "true"
)

// deleteFile removes a file if it exists
func deleteFile(f string) bool {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return true
	}

	if err := os.Remove(f); err != nil {
		return false
	}

	return true
}

// DiffArray returns the delta of two arrays
func DiffArray(a, b []string) (diff sort.StringSlice) {
	var biggest, smallest []string

	switch {
	case len(a) > len(b), len(a) == len(b):
		biggest = a
		smallest = b

	case len(a) < len(b):
		biggest = b
		smallest = a
	}

	dmap := make(List)
	for _, k := range smallest {
		dmap[k] = 0
	}

	for _, k := range biggest {
		if !dmap.keyExists(k) {
			diff = append(diff, k)
		}
	}
	diff.Sort()
	return diff
}

// Excludes returns a string array of excludes
func (c *Config) Excludes(node string) []string {
	var exc []string
	switch {
	case node == all:
		for _, k := range c.Nodes() {
			if len(c.bNodes[k].exc) != 0 {
				exc = append(exc, c.bNodes[k].exc...)
			}
		}
	default:
		exc = c.bNodes[node].exc
	}
	return exc
}

// Files returns a list of dnsmasq conf files from all srcs
func (o Objects) Files() *CFile {
	b := false
	c := CFile{Parms: o.Parms}
	obj := o.S
	for k := range obj {
		for sk := range obj[k].data {
			if !b {
				c.nType = obj[k].nType
			}

			src := obj[k].data[sk]
			format := src.Parms.Dir + "/%v.%v." + src.Parms.Ext
			c.names = append(c.names, fmt.Sprintf(format, getType(src.nType), src.name))
		}
	}
	sort.Strings(c.names)
	return &c
}

// Get returns an *Object for a given node
func (c *Config) Get(node string) (o *Object) {
	getObj := func(o *Object, node string) {
		for k := range o.data {
			o.data[k].Parms = c.Parms // TODO: check this is needed
			if o.data[k].ip == "" {
				o.data[k].ip = c.bNodes[node].ip
			}
			o.data[k].nType = getType(node).(ntype)
			switch {
			case o.data[k].url != "":
				o.data[k].ltype = "urls"
			case o.data[k].file != "":
				o.data[k].ltype = "files"
			}
		}
	}

	getInc := func(o *Object, node string) {
		if len(o.inc) > 0 {
			o.data[preConf] = &Object{
				desc:  preConf,
				inc:   o.inc,
				ip:    o.ip,
				ltype: preConf,
				name:  preConf,
				nType: getType(node).(ntype),
				Parms: c.Parms,
			}
		}
	}

	mergeList := func(a, b *Object) *Object {
		for k, v := range a.data {
			b.data[k] = v
		}
		return b
	}

	switch node {
	case all:
		o = &Object{Parms: c.Parms, data: make(data)}
		d := make([]*Object, len(c.Parms.Nodes))
		for i, node := range c.Parms.Nodes {
			d[i] = c.bNodes[node]
			getObj(d[i], node)
			getInc(d[i], node)
			o = mergeList(d[i], o)
		}

	default:
		o = c.bNodes[node]
		getObj(o, node)
		getInc(o, node)
	}

	// o = c.bNodes[node]
	// getObj(o, node)
	// getInc(o, node)

	return o
}

// GetAll returns and array of Objects
func (c *Config) GetAll() *Objects {
	o := Objects{Parms: c.Parms}
	for _, node := range c.Parms.Nodes {
		o.S = append(o.S, c.Get(node).Source(all).S...)
	}
	sort.Sort(&o)
	return &o
}

// getSubdomains returns a map of subdomains
func getSubdomains(s string) (l List) {
	l = make(List)
	keys := strings.Split(s, ".")
	for i := 0; i < len(keys)-1; i++ {
		key := strings.Join(keys[i:], ".")
		l[key] = 0
	}
	return l
}

// getType returns the converted "in" type
func getType(in interface{}) (out interface{}) {
	switch in.(type) {
	case ntype:
		out = typeInt(in.(ntype))
	case string:
		out = typeStr(in.(string))
	}
	return out
}

// Includes returns an io.Reader of blacklist Includes
func (o *Object) Includes() io.Reader {
	sort.Strings(o.inc)
	return bytes.NewBuffer([]byte(strings.Join(o.inc, "\n")))
}

// Load returns an EdgeOS CLI loaded configuration
func (c *CFGstatic) Load() io.Reader {
	return bytes.NewBufferString(c.Cfg)
}

func newObject() *Object {
	return &Object{
		data: make(data),
		exc:  make([]string, 0),
		inc:  make([]string, 0),
	}
}

// Nodes returns an array of configured nodes
func (c *Config) Nodes() (nodes []string) {
	for k := range c.bNodes {
		nodes = append(nodes, k)
	}
	sort.Strings(nodes)
	return nodes
}

// ReadCfg extracts nodes from a EdgeOS/VyOS configuration structure
func ReadCfg(r ConfLoader) (*Config, error) {
	var (
		tnode  string
		b      = bufio.NewScanner(r.Load())
		branch string
		nodes  = make([]string, 2)
		rx     = regx.Objects
		s      *Object
		sCfg   = Config{Parms: &Parms{}, bNodes: make(bNodes)}
	)

LINE:
	for b.Scan() {
		line := strings.TrimSpace(b.Text())

		switch {
		case rx.MLTI.MatchString(line):
			incExc := regx.Get("mlti", line)
			switch incExc[1] {
			case "exclude":
				sCfg.bNodes[tnode].exc = append(sCfg.bNodes[tnode].exc, incExc[2])

			case "include":
				sCfg.bNodes[tnode].inc = append(sCfg.bNodes[tnode].inc, incExc[2])
			}

		case rx.NODE.MatchString(line):
			node := regx.Get("node", line)
			tnode = node[1]
			nodes = append(nodes, tnode)
			s = newObject()
			sCfg.bNodes[tnode] = s

		case rx.LEAF.MatchString(line):
			srcName := regx.Get("leaf", line)
			branch = srcName[2]
			nodes = append(nodes, srcName[1])

			if srcName[1] == src {
				s.name = branch
				s.nType = getType(tnode).(ntype)
			}

		case rx.DSBL.MatchString(line):
			sCfg.bNodes[tnode].disabled = StrToBool(regx.Get("dsbl", line)[1])

		case rx.IPBH.MatchString(line) && nodes[len(nodes)-1] != src:
			sCfg.bNodes[tnode].ip = regx.Get("ipbh", line)[1]

		case rx.NAME.MatchString(line):
			name := regx.Get("name", line)

			switch name[1] {
			case "description":
				s.desc = name[2]

			case blackhole:
				s.ip = name[2]

			case "file":
				s.file = name[2]
				s.ltype = name[1]
				sCfg.bNodes[tnode].data[branch] = s
				s = newObject() // reset s for the next loop

			case "prefix":
				s.prefix = name[2]

			case "url":
				s.ltype = name[1]
				s.url = name[2]
				sCfg.bNodes[tnode].data[branch] = s
				s = newObject() // reset s for the next loop

			}

		case rx.DESC.MatchString(line) || rx.CMNT.MatchString(line) || rx.MISC.MatchString(line):
			continue LINE

		case rx.RBRC.MatchString(line):
			nodes = nodes[:len(nodes)-1] // pop last node
			tnode = nodes[len(nodes)-1]
		}
	}

	if len(sCfg.bNodes) < 1 {
		return &sCfg, errors.New("Configuration data is empty, cannot continue")
	}
	return &sCfg, nil
}

// Remove deletes a CFile array of file names
func (c *CFile) Remove() error {
	var got = make([]string, 5)
	dlist, err := ioutil.ReadDir(c.Dir)
	if err != nil {
		return err
	}

	for _, f := range dlist {
		if strings.Contains(f.Name(), getType(c.nType).(string)) && strings.Contains(f.Name(), c.Ext) {
			got = append(got, c.Dir+"/"+f.Name())
		}
	}

	return purgeFiles(DiffArray(c.names, got))
}

// Source returns a map of sources
func (d data) Source(ltype string) *Objects {
	b := false
	var p *Parms
	objs := []*Object{}
	for _, k := range d.sortSKeys() {
		if !b {
			if p = d[k].Parms; p.Dir != "" {
				b = true
			}
		}
		switch {
		case ltype == d[k].ltype:
			objs = append(objs, d[k])
		case ltype == "all":
			objs = append(objs, d[k])
		}
	}
	return &Objects{Parms: p, S: objs}
}

// String implements string method
func (c *CFile) String() string {
	return strings.Join(c.names, "\n")
}

// Strings returns a sorted array of strings.
func (c *CFile) Strings() []string {
	sort.Strings(c.names)
	return c.names
}

// STypes returns an array of configured nodes
func (c *Config) STypes() []string {
	return c.Parms.Stypes
}

// BooltoStr converts a boolean ("true" or "false") to a string equivalent
func BooltoStr(b bool) string {
	if b {
		return True
	}
	return False
}

// StrToBool converts a string ("true" or "false") to it's boolean equivalent
func StrToBool(s string) bool {
	if strings.ToLower(s) == True {
		return true
	}
	return false
}

func typeInt(i ntype) (s string) {
	switch i {
	case domain:
		s = domains
	case host:
		s = hosts
	case pre:
		s = preConf
	case root:
		s = blacklist
	case unknown:
		s = notknown
	case zone:
		s = zones
	}
	return s
}

func typeStr(s string) (i ntype) {
	switch s {
	case domains:
		i = domain
	case hosts:
		i = host
	case preConf:
		i = pre
	case blacklist:
		i = root
	case notknown:
		i = unknown
	case zones:
		i = zone
	}
	return i
}

// UpdateList converts []string to map of List
func UpdateList(data []string) (l List) {
	l = make(List)
	for _, k := range data {
		l[k] = 0
	}
	return l
}

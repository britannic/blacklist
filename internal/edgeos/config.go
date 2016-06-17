// Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.
package edgeos

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
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
func (o *Objects) Files() *CFile {
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
			o.data[k].Parms = c.Parms
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

	return o
}

// GetAll returns an array of Objects
func (c *Config) GetAll() *Objects {
	o := Objects{Parms: c.Parms}
	for _, node := range c.Parms.Nodes {
		o.S = append(o.S, c.Get(node).Source(all).S...)
	}
	sort.Sort(&o)
	return &o
}

// Load returns an EdgeOS CLI loaded configuration
func (c *CFGstatic) Load() io.Reader {
	return bytes.NewBufferString(c.Cfg)
}

// Nodes returns an array of configured nodes
func (c *Config) Nodes() (nodes []string) {
	for k := range c.bNodes {
		nodes = append(nodes, k)
	}
	sort.Strings(nodes)
	return nodes
}

// ReadDir implements OSinformer
func (c *CFile) ReadDir(dir string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dir)
}

// Remove deletes a CFile array of file names
func (c *CFile) Remove() error {
	var got = make([]string, 5)

	dlist, err := c.ReadDir(c.Dir)
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

// String returns pretty print for the Blacklist struct
func (c *Config) String() (result string) {
	indent := 1
	cmma := comma
	cnt := len(c.sortKeys())
	result += fmt.Sprintf("{\n%v%q: [{\n", tabs(indent), "nodes")

	for i, pkey := range c.sortKeys() {

		if i == cnt-1 {
			cmma = null
		}

		indent++
		result += fmt.Sprintf("%v%q: {\n", tabs(indent), pkey)

		indent++
		result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), disabled, getJSONdisabled(&cfgJSON{Config: c, pk: pkey}))

		result += tabs(indent) + getJSONsrcIP(c, pkey)

		result += getJSONArray(&cfgJSON{array: c.bNodes[pkey].exc, pk: pkey, leaf: "excludes", indent: indent})

		if pkey != rootNode {
			result += getJSONArray(&cfgJSON{array: c.bNodes[pkey].inc, pk: pkey, leaf: "includes", indent: indent})
		}

		if pkey != rootNode {
			result += getJSONsrcArray(&cfgJSON{Config: c, pk: pkey, indent: indent})
		}

		indent--
		result += fmt.Sprintf("%v}%v\n", tabs(indent), cmma)
		indent--
	}

	result += tabs(indent) + "}]\n}"
	return result
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

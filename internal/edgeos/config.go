// Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.
package edgeos

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
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
	files     = "file"
	src       = "source"
	hosts     = "hosts"
	preConf   = "pre-configured"
	rootNode  = blacklist
	urls      = "url"
	zones     = "zones"

	// False is a string constant
	False = "false"
	// True is a string constant
	True = "true"
)

// Excludes returns a List map of blacklist exclusions
func (c *Config) Excludes(node string) List {
	return UpdateList(c.excludes(node))
}

// Excludes returns a string array of excludes
func (c *Config) excludes(node string) []string {
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
	c := CFile{Parms: o.Parms}
	for _, obj := range o.S {
		c.nType = obj.nType
		format := o.Parms.Dir + "/%v.%v." + o.Parms.Ext
		c.names = append(c.names, fmt.Sprintf(format, getType(obj.nType), obj.name))
	}
	sort.Strings(c.names)
	return &c
}

func getInc(obj *Object, node string) []*Object {
	return []*Object{
		&Object{
			desc:  preConf,
			inc:   obj.inc,
			ip:    obj.ip,
			ltype: preConf,
			name:  preConf,
			nType: getType(node).(ntype),
			Parms: obj.Parms,
		},
	}
}

func (b bNodes) validate(node string) *Objects {
	for _, obj := range b[node].Objects.S {
		if obj.ip == "" {
			obj.ip = b[node].ip
		}
	}
	return &b[node].Objects
}

// Get returns an *Object for a given node
func (c *Config) Get(node string) *Objects {
	o := &Objects{Parms: c.Parms}

	switch node {
	case all:
		for _, node := range c.Parms.Nodes {
			o.addInc(c, node)
			o.addObj(c, node)
		}
	default:
		o.addInc(c, node)
		o.addObj(c, node)
	}
	return o
}

func (o *Objects) addInc(c *Config, node string) {
	if c.bNodes[node].inc != nil {
		o.S = append(o.S, getInc(&Object{Parms: c.Parms, inc: c.bNodes[node].inc, ip: c.bNodes[node].ip}, node)...)
	}
}

func (o *Objects) addObj(c *Config, node string) {
	o.S = append(o.S, c.bNodes.validate(node).S...)
}

// GetAll returns an array of Objects
func (c *Config) GetAll(ltypes ...string) *Objects {
	var (
		o = &Objects{Parms: c.Parms}
	)

	for _, node := range c.Parms.Nodes {
		switch ltypes {
		case nil:
			o.addInc(c, node)
			o.addObj(c, node)
		default:
			for _, ltype := range ltypes {
				switch ltype {
				case preConf:
					o.addInc(c, node)
				default:
					obj := c.bNodes[node].Objects.S
					for i := range obj {
						b := obj[i].ltype == ltype
						if b {
							o.S = append(o.S, obj[i])
						}
					}
				}
			}
		}
	}
	return o
}

// insession returns true if VyOS/EdgeOS configuration is in session
func (c *Config) insession() bool {
	var (
		cmd = exec.Command(c.API, "inSession")
		out bytes.Buffer
	)

	cmd.Stdout = &out
	if ok := cmd.Run(); ok == nil {
		return out.String() == "0"
	}
	return false
}

// Load returns an EdgeOS CLI loaded configuration
func (c *CFGstatic) Load() io.Reader {
	return bytes.NewBufferString(c.Cfg)
}

// NewConfig returns a new *Config initialized with the parameter options passed to it
func NewConfig(opts ...Option) *Config {
	c := Config{
		bNodes: make(bNodes),
		Parms: &Parms{
			Dex: make(List),
			Exc: make(List),
		},
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
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
func (c *Config) ReadCfg(r ConfLoader) error {
	var (
		tnode  string
		b      = bufio.NewScanner(r.Load())
		branch string
		nodes  = make([]string, 2)
		rx     = regx.Objects
		s      *Object
	)

LINE:
	for b.Scan() {
		line := strings.TrimSpace(b.Text())

		switch {
		case rx.MLTI.MatchString(line):
			incExc := regx.Get("mlti", line)
			switch incExc[1] {
			case "exclude":
				c.bNodes[tnode].exc = append(c.bNodes[tnode].exc, incExc[2])

			case "include":
				c.bNodes[tnode].inc = append(c.bNodes[tnode].inc, incExc[2])
			}

		case rx.NODE.MatchString(line):
			node := regx.Get("node", line)
			tnode = node[1]
			nodes = append(nodes, tnode)
			s = newObject()
			c.bNodes[tnode] = s

		case rx.LEAF.MatchString(line):
			srcName := regx.Get("leaf", line)
			branch = srcName[2]
			nodes = append(nodes, srcName[1])

			if srcName[1] == src {
				s.name = branch
				s.nType = getType(tnode).(ntype)
			}

		case rx.DSBL.MatchString(line):
			c.bNodes[tnode].disabled = StrToBool(regx.Get("dsbl", line)[1])

		case rx.IPBH.MatchString(line) && nodes[len(nodes)-1] != src:
			c.bNodes[tnode].ip = regx.Get("ipbh", line)[1]

		case rx.NAME.MatchString(line):
			name := regx.Get("name", line)

			switch name[1] {
			case "description":
				s.desc = name[2]

			case blackhole:
				s.ip = name[2]

			case files:
				s.file = name[2]
				s.ltype = name[1]
				c.bNodes[tnode].Objects.S = append(c.bNodes[tnode].Objects.S, s)
				s = newObject() // reset s for the next loop

			case "prefix":
				s.prefix = name[2]

			case urls:
				s.ltype = name[1]
				s.url = name[2]
				c.bNodes[tnode].Objects.S = append(c.bNodes[tnode].Objects.S, s)
				s = newObject() // reset s for the next loop

			}

		case rx.DESC.MatchString(line) || rx.CMNT.MatchString(line) || rx.MISC.MatchString(line):
			continue LINE

		case rx.RBRC.MatchString(line):
			nodes = nodes[:len(nodes)-1] // pop last node
			tnode = nodes[len(nodes)-1]
		}
	}

	if len(c.bNodes) < 1 {
		return errors.New("Configuration data is empty, cannot continue")
	}

	return nil
}

// ReadDir returns a listing of dnsmasq formatted blacklist configuration files
func (c *CFile) ReadDir(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

// Remove deletes a CFile array of file names
func (c *CFile) Remove() error {
	if c.Wildcard == (Wildcard{}) {
		c.Wildcard = Wildcard{node: "*s", name: "*"}
	}

	pattern := fmt.Sprintf(c.FnFmt, c.Dir, c.Wildcard.node, c.Wildcard.name, c.Parms.Ext)
	dlist, err := c.ReadDir(pattern)
	if err != nil {
		return err
	}
	return purgeFiles(DiffArray(c.names, dlist))
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
		result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), disabled,
			BooltoStr(c.bNodes[pkey].disabled))

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

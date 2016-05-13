// Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.
package edgeos

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/britannic/blacklist/regx"
)

var bnodes = []string{Root, Domains, Hosts}

const (
	agent     = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/601.4.4 (KHTML, like Gecko) Version/9.0.3 Safari/601.4.4`
	blackhole = "dns-redirect-ip"
	blacklist = "blacklist"
	dbg       = false
	disabled  = "disabled"
	source    = "source"

	// Fext sets the dnsmasq configuration file extension
	Fext = "blacklist.conf"

	// Domains sets the domains string
	Domains = "domains"

	// Hosts sets the hosts string
	Hosts = "hosts"

	// PreCon sets the string for pre-configured
	PreCon = "pre-configured"

	// Root is the topmost node
	Root = blacklist

	// False is a string constant
	False = "false"

	// True is a string constant
	True = "true"
)

// Types determine load order and processing behavior for blacklist sources
const (
	// Unknown shouldn't ever be used, it denotes a coding error
	unknown int = iota

	// Pre type is for pre-configured backlisted domains/hosts
	pre

	// Domain type sets which format is used for dnsmasq conf files
	// e.g. address=/.domain.com/0.0.0.0
	domain

	// Host type sets which format is used for dnsmasq conf files
	// e.g. address=/www.domain.com/0.0.0.0
	host

	// Root type is the topmost root node
	root

	// Zone type is for future use
	zone
)

// Leaf is a struct for EdgeOS configuration data
type Leaf struct {
	Data     map[string]*Srcs `json:"data, omitempty"`
	Disabled bool             `json:"disable"`
	Excludes []string         `json:"excludes, omitempty"`
	Includes []string         `json:"includes, omitempty"`
	IP       string           `json:"ip, omitempty"`
}

// Srcs holds download source information
type Srcs struct {
	Desc     string `json:"desc, omitempty"`
	Disabled bool   `json:"disabled, omitempty"`
	File     string `json:"file, omitempty"`
	IP       string `json:"ip, omitempty"`
	List     List   `json:"-"`
	Name     string `json:"name"`
	No       int    `json:"-"`
	Prefix   string `json:"prefix"`
	Type     int    `json:"type, omitempty"`
	URL      string `json:"url, omitempty"`
}

// NewNodes implements a new Node map
func NewNodes() Nodes {
	return make(Nodes)
}

// GetSubdomains returns a map of subdomains
func GetSubdomains(s string) (l List) {
	l = make(List)
	keys := strings.Split(s, ".")
	for i := 0; i < len(keys)-1; i++ {
		key := strings.Join(keys[i:], ".")
		l[key] = 0
	}
	return l
}

func typeInt(i int) (s string) {
	switch i {
	case domain:
		s = Domains
	case host:
		s = Hosts
	case pre:
		s = PreCon
	case root:
		s = blacklist
	case unknown:
		s = "unknown"
	}
	return s
}

func typeStr(s string) (i int) {
	switch s {
	case Domains:
		i = domain
	case Hosts:
		i = host
	case PreCon:
		i = pre
	case blacklist:
		i = root
	case "unknown":
		i = unknown
	}
	return i
}

// GetType returns the converted "in" type
func GetType(in interface{}) (out interface{}) {
	switch in.(type) {
	case int:
		out = typeInt(in.(int))
	case string:
		out = typeStr(in.(string))
	}
	return out
}

// ReadCfg extracts nodes from a EdgeOS/VyOS configuration structure
func ReadCfg(reader io.Reader) (Nodes, error) {
	var (
		tnode string
		b     = bufio.NewScanner(reader)
		leaf  string
		nodes = make([]string, 2)
		rx    = regx.Objects
		sCfg  = NewNodes()
		src   = &Srcs{}
	)

LINE:
	for b.Scan() {
		line := strings.TrimSpace(b.Text())

		switch {
		case rx.MLTI.MatchString(line):
			IncExc := regx.Get("mlti", line)
			switch IncExc[1] {
			case "exclude":
				sCfg[tnode].Excludes = append(sCfg[tnode].Excludes, IncExc[2])

			case "include":
				sCfg[tnode].Includes = append(sCfg[tnode].Includes, IncExc[2])

			}

		case rx.NODE.MatchString(line):
			node := regx.Get("node", line)
			tnode = node[1]
			nodes = append(nodes, tnode)

			sCfg[tnode] = &Leaf{Includes: make([]string, 0), Excludes: make([]string, 0), Data: make(map[string]*Srcs)}

		case rx.LEAF.MatchString(line):
			srcName := regx.Get("leaf", line)
			leaf = srcName[2]
			nodes = append(nodes, srcName[1])

			if srcName[1] == source {
				src.Name = leaf
				src.Type = GetType(tnode).(int)
			}

		case rx.DSBL.MatchString(line):
			sCfg[tnode].Disabled = ToBool(regx.Get("dsbl", line)[1])

		case rx.IPBH.MatchString(line) && nodes[len(nodes)-1] != source:
			sCfg[tnode].IP = regx.Get("ipbh", line)[1]

		case rx.NAME.MatchString(line):
			name := regx.Get("name", line)

			switch name[1] {
			case "description":
				src.Desc = name[2]

			case blackhole:
				src.IP = name[2]

			case "file":
				src.File = name[2]
				sCfg[tnode].Data[leaf] = src
				src = &Srcs{} // reset src for the next loop

			case "prefix":
				src.Prefix = name[2]

			case "url":
				src.URL = name[2]
				sCfg[tnode].Data[leaf] = src
				src = &Srcs{} // reset src for the next loop

			}

		case rx.DESC.MatchString(line) || rx.CMNT.MatchString(line) || rx.MISC.MatchString(line):
			continue LINE

		case rx.RBRC.MatchString(line):
			nodes = nodes[:len(nodes)-1] // pop last node
			tnode = nodes[len(nodes)-1]
		}
	}

	if len(sCfg) < 1 {
		return sCfg, errors.New("Configuration data is empty, cannot continue")
	}

	return sCfg, nil
}

// Nodes is a map of Leaf nodes
type Nodes map[string]*Leaf

// ToBool converts a string ("true" or "false") to it's boolean equivalent
func ToBool(s string) bool {
	if strings.ToLower(s) == "true" {
		return true
	}
	return false
}

// func tabs(t int) (r string) {
// 	if t <= 0 {
// 		return r
// 	}
// 	for i := 0; i < t; i++ {
// 		r += tab
// 	}
// 	return r
// }
//
// type cfgJSON struct {
// 	array        []string
// 	blist        Blacklist
// 	indent       int
// 	leaf, pk, sk string
// }
//
// func getJSONdisabled(c *cfgJSON) (d string) {
// 	switch c.sk {
// 	case null:
// 		d = c.blist[c.pk].Disable
// 	default:
// 		d = c.blist[c.pk].Source[c.sk].Disable
// 	}
//
// 	switch d {
// 	case False, True:
// 		return d
// 	default:
// 		return False
// 	}
// }
// func getJSONsrcIP(b Blacklist, pkey string) (result string) {
// 	if len(b[pkey].IP) > 0 {
// 		result += fmt.Sprintf("%q: %q,\n", "ip", b[pkey].IP)
// 	}
// 	return result
// }
//
// func getJSONArray(c *cfgJSON) (result string) {
// 	indent := c.indent
// 	cmma := comma
// 	ret := enter
// 	result += fmt.Sprintf("%v%q: [", tabs(indent), c.leaf)
// 	cnt := len(c.array)
//
// 	switch {
// 	case c.pk != root && cnt == 0:
// 		result += "],\n"
// 		return result
//
// 	// case cnt == 0:
// 	// 	result += "]\n"
// 	// 	return result
//
// 	case cnt == 1:
// 		ret = null
// 		indent = 0
//
// 	case cnt > 1:
// 		result += enter
// 		indent++
// 	}
//
// 	if cnt > 0 {
// 		for i, s := range c.array {
// 			if i == cnt-1 {
// 				cmma = null
// 			}
// 			result += fmt.Sprintf("%v%q%v%v", tabs(indent), s, cmma, ret)
// 		}
//
// 		cmma = comma
//
// 		if c.pk == root {
// 			cmma = null
// 		}
//
// 		result += fmt.Sprintf("%v]%v\n", tabs(indent), cmma)
// 	}
//
// 	return result
// }

// func getJSONsrcArray(c *cfgJSON) (result string) {
// 	var i int
// 	indent := c.indent
// 	skeys := c.blist.SortSKeys(c.pk)
// 	cnt := len(skeys)
// 	if cnt == 0 {
// 		result += fmt.Sprintf("%v%q: [{}]\n", tabs(c.indent), "sources")
// 		return result
// 	}
// 	result += fmt.Sprintf("%v%q: [{%v", tabs(c.indent), "sources", enter)
//
// 	for i, c.sk = range skeys {
// 		if _, ok := c.blist[c.pk].Source[c.sk]; ok {
// 			cmma := comma
// 			indent = c.indent + 1
//
// 			if i == cnt-1 {
// 				cmma = null
// 			}
//
// 			d := getJSONdisabled(&cfgJSON{blist: c.blist, pk: c.pk, sk: c.sk})
// 			s := c.blist[c.pk].Source[c.sk]
// 			// result += fmt.Sprintf("%v%q: {%v", tabs(indent), "source", enter)
// 			// indent++
// 			result += fmt.Sprintf("%v%q: {\n", tabs(indent), c.sk)
// 			indent++
// 			result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), disabled, d)
// 			result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), "description", s.Desc)
// 			result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), "prefix", s.Prfx)
// 			result += fmt.Sprintf("%v%q: %q\n", tabs(indent), "url", s.URL)
// 			indent--
// 			result += fmt.Sprintf("%v}%v%v", tabs(indent), cmma, enter)
// 		}
// 	}
//
// 	indent -= 2
// 	result += fmt.Sprintf("%v}]%v", tabs(indent), enter)
// 	return result
// }

// String returns pretty print for the Blacklist struct
// func (b Blacklist) String() (result string) {
// 	indent := 1
// 	cmma := comma
// 	cnt := len(b.SortKeys())
// 	result += fmt.Sprintf("{\n%v%q: [{\n", tabs(indent), "nodes")
//
// 	for i, pkey := range b.SortKeys() {
//
// 		if i == cnt-1 {
// 			cmma = null
// 		}
//
// 		indent++
// 		result += fmt.Sprintf("%v%q: {\n", tabs(indent), pkey)
//
// 		indent++
// 		result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), disabled, getJSONdisabled(&cfgJSON{blist: b, pk: pkey}))
//
// 		result += tabs(indent) + getJSONsrcIP(b, pkey)
//
// 		result += getJSONArray(&cfgJSON{array: b[pkey].Exclude, pk: pkey, leaf: "excludes", indent: indent})
//
// 		if pkey != root {
// 			result += getJSONArray(&cfgJSON{array: b[pkey].Include, pk: pkey, leaf: "includes", indent: indent})
// 		}
//
// 		if pkey != root {
// 			result += getJSONsrcArray(&cfgJSON{blist: b, pk: pkey, indent: indent})
// 		}
//
// 		indent--
// 		result += fmt.Sprintf("%v}%v\n", tabs(indent), cmma)
// 		indent--
// 	}
//
// 	result += tabs(indent) + "}]\n}"
// 	return result
// }

package edgeos

import (
	"fmt"
	"sort"
	"strings"
)

// Objects is a struct of []*source
type Objects struct {
	*Env
	src []*source
}

func (o *Objects) addObj(c *Config, node string) {
	o.src = append(o.src, c.addInc(node))
	o.src = append(o.src, c.tree.validate(node).src...)
}

// Files returns a list of dnsmasq conf files from all srcs
func (o *Objects) Files() *CFile {
	var c = CFile{Env: o.Env}
	if !o.Disabled {
		for _, obj := range o.src {
			c.Names = append(c.Names, obj.setFilePrefix(o.Env.Dir+"/%v.%v."+o.Env.Ext))
			c.nType = obj.nType
		}
		sort.Strings(c.Names)
	}
	return &c
}

// Filter returns a subset of Objects filtered by ltype
func (o *Objects) Filter(ltype string) *Objects {
	objx := Objects{Env: o.Env}
	switch ltype {
	case files:
		for _, obj := range o.src {
			if obj.ltype == files && obj.file != "" {
				objx.src = append(objx.src, obj)
			}
		}
	case urls:
		for _, obj := range o.src {
			if obj.ltype == urls && obj.url != "" {
				objx.src = append(objx.src, obj)
			}
		}
	}
	return &objx
}

// Find returns the int position of an Objects' element
func (o *Objects) Find(elem string) int {
	for i, obj := range o.src {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

func getLtypeDesc(s string) string {
	if s == "" {
		s = "unknown ltype"
	}
	return fmt.Sprintf("%s %s", preNoun, strings.Join(strings.Split(s, "-"), " "))
}

func (o *Objects) procltypes(c *Config, node string, ltypes ...string) {
	var (
		newDomns = true
		newHosts = true
		newRoots = true
	)
	switch ltypes {
	case nil:
		o.addObj(c, node)
	default:
		for _, ltype := range ltypes {
			switch ltype {
			case PreDomns:
				if newDomns && node == domains {
					o.src = append(o.src, c.addInc(node))
					newDomns = false
				}
			case PreHosts:
				if newHosts && node == hosts {
					o.src = append(o.src, c.addInc(node))
					newHosts = false
				}
			case PreRoots:
				if newRoots && node == rootNode {
					o.src = append(o.src, c.addInc(node))
					newRoots = false
				}
			default:
				obj := c.validate(node).src
				for i := range obj {
					if obj[i].ltype == ltype {
						o.src = append(o.src, obj[i])
					}
				}
			}
		}
	}
}

func (o *Objects) objects(c *Config, node string, ltypes ...string) {
	switch node {
	case domains:
		o.procltypes(c, node, ltypes...)
	case hosts:
		o.procltypes(c, node, ltypes...)
	case rootNode:
		o.procltypes(c, node, ltypes...)
	}
}

// Names returns a sorted slice of Objects names
func (o *Objects) Names() (s sort.StringSlice) {
	for _, obj := range o.src {
		s = append(s, obj.name)
	}
	sort.Sort(s)
	return s
}

// Stringer for Objects
func (o *Objects) String() (s string) {
	for _, src := range o.src {
		s += src.String()
	}
	return s
}

// Implement Sort Interface for Objects
func (o *Objects) Len() int           { return len(o.src) }
func (o *Objects) Less(i, j int) bool { return o.src[i].name < o.src[j].name }
func (o *Objects) Swap(i, j int)      { o.src[i], o.src[j] = o.src[j], o.src[i] }

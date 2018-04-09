package edgeos

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// source struct for normalizing EdgeOS data.
type source struct {
	*Parms
	desc     string
	disabled bool
	err      error
	exc      []string
	file     string
	inc      []string
	ip       string
	ltype    string
	name     string
	nType    ntype
	Objects
	prefix string
	r      io.Reader
	url    string
}

// Objects is a struct of []*source
type Objects struct {
	*Parms
	xx []*source
}

func (o *Objects) addObj(c *Config, node string) {
	obj := c.addInc(node)
	o.xx = append(o.xx, obj)
	o.xx = append(o.xx, c.tree.validate(node).xx...)
}

func (o *source) area() string {
	switch getType(o.nType).(string) {
	case domains, PreDomns:
		return domains
	}
	return hosts
}

// excludes returns an io.Reader of blacklist includes
func (o *source) excludes() io.Reader {
	sort.Strings(o.exc)
	return strings.NewReader(strings.Join(o.exc, "\n"))
}

func (o *source) setFilePrefix(format string) (f string) {
	switch o.nType {
	case excDomn, excRoot, preDomn:
		f = fmt.Sprintf(format, domains, o.name)
	case excHost, preHost:
		f = fmt.Sprintf(format, hosts, o.name)
	default:
		f = fmt.Sprintf(format, getType(o.nType), o.name)
	}
	return f
}

// Files returns a list of dnsmasq conf files from all srcs
func (o *Objects) Files() *CFile {
	var c = CFile{Parms: o.Parms}

	if !o.Disabled {
		for _, obj := range o.xx {
			f := obj.setFilePrefix(o.Parms.Dir + "/%v.%v." + o.Parms.Ext)
			c.Names = append(c.Names, f)
			c.nType = obj.nType
		}
		sort.Strings(c.Names)
	}
	return &c
}

// Filter returns a subset of Objects filtered by ltype
func (o *Objects) Filter(ltype string) *Objects {
	sources := Objects{Parms: o.Parms}

	switch ltype {
	case files:
		for _, obj := range o.xx {
			if obj.ltype == files && obj.file != "" {
				sources.xx = append(sources.xx, obj)
			}
		}
	case urls:
		for _, obj := range o.xx {
			if obj.ltype == urls && obj.url != "" {
				sources.xx = append(sources.xx, obj)
			}
		}
	default:
		sources = Objects{Parms: o.Parms}
	}
	return &sources
}

// Find returns the int position of an Objects' element
func (o *Objects) Find(elem string) int {
	for i, obj := range o.xx {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

func getLtypeDesc(l string) string {
	switch l {
	case ExcDomns:
		return preNoun + " whitelisted domains"
	case ExcHosts:
		return preNoun + " whitelisted hosts"
	case ExcRoots:
		return preNoun + " global whitelisted domains"
	case PreDomns:
		return preNoun + " blacklisted domains"
	case PreHosts:
		return preNoun + " blacklisted hosts"
	default:
		return "Unknown ltype"
	}
}

// includes returns an io.Reader of blacklist includes
func (o *source) includes() io.Reader {
	sort.Strings(o.inc)
	return strings.NewReader(strings.Join(o.inc, "\n"))
}

func (o *Objects) objects(c *Config, node string, ltypes ...string) {
	var (
		newDomns = true
		newHosts = true
	)

	switch ltypes {
	case nil:
		o.addObj(c, node)
	default:
		for _, ltype := range ltypes {
			switch ltype {
			case PreDomns:
				if newDomns && node == domains {
					o.xx = append(o.xx, c.addInc(node))
					newDomns = false
				}
			case PreHosts:
				if newHosts && node == hosts {
					o.xx = append(o.xx, c.addInc(node))
					newHosts = false
				}
			default:
				obj := c.validate(node).xx
				for i := range obj {
					if obj[i].ltype == ltype {
						o.xx = append(o.xx, obj[i])
					}
				}
			}
		}
	}
}

// Names returns a sorted slice of Objects names
func (o *Objects) Names() (s sort.StringSlice) {
	for _, obj := range o.xx {
		s = append(s, obj.name)
	}
	sort.Sort(s)
	return s
}

func newSource() *source {
	return &source{
		Objects: Objects{},
		exc:     []string{},
		inc:     []string{},
	}
}

// Stringer for Object
func (o *source) String() string {
	s := fmt.Sprintf("\nDesc:\t %q\n", o.desc)
	s += fmt.Sprintf("Disabled: %v\n", o.disabled)
	s += fmt.Sprintf("File:\t %q\n", o.file)
	s += fmt.Sprintf("IP:\t %q\n", o.ip)
	s += fmt.Sprintf("Ltype:\t %q\n", o.ltype)
	s += fmt.Sprintf("Name:\t %q\n", o.name)
	s += fmt.Sprintf("nType:\t %q\n", o.nType)
	s += fmt.Sprintf("Prefix:\t %q\n", o.prefix)
	s += fmt.Sprintf("Type:\t %q\n", getType(o.nType))
	s += fmt.Sprintf("URL:\t %q\n", o.url)
	return s
}

// Stringer for Objects
func (o *Objects) String() string {
	return fmt.Sprint(o.xx)
}

// Implement Sort Interface for Objects
func (o *Objects) Len() int           { return len(o.xx) }
func (o *Objects) Less(i, j int) bool { return o.xx[i].name < o.xx[j].name }
func (o *Objects) Swap(i, j int)      { o.xx[i], o.xx[j] = o.xx[j], o.xx[i] }

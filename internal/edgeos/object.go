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
	switch obj := c.addInc(node); obj {
	case nil:
		o.xx = append(o.xx, c.tree.validate(node).xx...)
	default:
		o.xx = append(o.xx, obj)
		o.xx = append(o.xx, c.tree.validate(node).xx...)
	}
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

// Files returns a list of dnsmasq conf files from all srcs
func (o *Objects) Files() *CFile {
	var f string
	c := CFile{Parms: o.Parms}
	format := o.Parms.Dir + "/%v.%v." + o.Parms.Ext
	if !o.Disabled {
		for _, obj := range o.xx {
			switch obj.nType {
			case excDomn, excRoot, preDomn:
				f = fmt.Sprintf(format, domains, obj.name)
			case excHost, preHost:
				f = fmt.Sprintf(format, hosts, obj.name)
			default:
				f = fmt.Sprintf(format, getType(obj.nType), obj.name)
			}
			c.Names = append(c.Names, f)
			c.nType = obj.nType
		}
		sort.Strings(c.Names)
	}
	return &c
}

// Filter returns a subset of Objects filtered by ltype
func (o *Objects) Filter(ltype string) *Objects {
	var (
		sources = Objects{Parms: o.Parms}
		xFiles  = "-" + files
		xURLs   = "-" + urls
	)

	switch ltype {
	case files:
		for _, obj := range o.xx {
			if obj.ltype == files && obj.file != "" {
				sources.xx = append(sources.xx, obj)
			}
		}
	case xFiles:
		for _, obj := range o.xx {
			if obj.ltype != files {
				sources.xx = append(sources.xx, obj)
			}
		}
	case urls:
		for _, obj := range o.xx {
			if obj.ltype == urls && obj.url != "" {
				sources.xx = append(sources.xx, obj)
			}
		}
	case xURLs:
		for _, obj := range o.xx {
			if obj.ltype != urls {
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
	case "":
		return "Unknown ltype"
	default:
		panic(fmt.Sprintf("getLtypeDesc(l) passed an illegal lType: %v", l))
	}
}

// includes returns an io.Reader of blacklist includes
func (o *source) includes() io.Reader {
	sort.Strings(o.inc)
	return strings.NewReader(strings.Join(o.inc, "\n"))
}

// Names returns a sorted slice of Objects names
func (o *Objects) Names() (s sort.StringSlice) {
	for _, obj := range o.xx {
		s = append(s, obj.name)
	}
	sort.Sort(s)
	return s
}

func newObject() *source {
	return &source{
		Objects: Objects{},
		exc:     make([]string, 0),
		inc:     make([]string, 0),
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

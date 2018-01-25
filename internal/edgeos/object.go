package edgeos

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// object struct for normalizing EdgeOS data.
type object struct {
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

// Objects is a struct of []*Object
type Objects struct {
	*Parms
	x []*object
}

func (o *Objects) addObj(c *Config, node string) {
	switch obj := c.addInc(node); obj {
	case nil:
		o.x = append(o.x, c.tree.validate(node).x...)
	default:
		o.x = append(o.x, obj)
		o.x = append(o.x, c.tree.validate(node).x...)
	}
}

// func (o *object) area() string {
// 	switch getType(o.nType).(string) {
// 	case PreDomns:
// 		return domains
// 	case PreHosts:
// 		return hosts
// 	}
// 	return getType(o.nType).(string)
// }

// excludes returns an io.Reader of blacklist includes
func (o *object) excludes() io.Reader {
	sort.Strings(o.exc)
	return strings.NewReader(strings.Join(o.exc, "\n"))
}

// Files returns a list of dnsmasq conf files from all srcs
func (o *Objects) Files() *CFile {
	c := CFile{Parms: o.Parms}
	if !o.Disabled {
		for _, obj := range o.x {
			c.nType = obj.nType
			format := o.Parms.Dir + "/%v.%v." + o.Parms.Ext
			c.Names = append(c.Names, fmt.Sprintf(format, getType(obj.nType), obj.name))
		}
		sort.Strings(c.Names)
	}
	return &c
}

// Filter returns a subset of Objects; ltypes with "-" prepended remove ltype
func (o *Objects) Filter(ltype string) *Objects {
	var (
		objects = Objects{Parms: o.Parms}
		xFiles  = "-" + files
		xURLs   = "-" + urls
	)

	switch ltype {
	case files:
		for _, obj := range o.x {
			if obj.ltype == files && obj.file != "" {
				objects.x = append(objects.x, obj)
			}
		}
	case xFiles:
		for _, obj := range o.x {
			if obj.ltype != files {
				objects.x = append(objects.x, obj)
			}
		}
	case urls:
		for _, obj := range o.x {
			if obj.ltype == urls && obj.url != "" {
				objects.x = append(objects.x, obj)
			}
		}
	case xURLs:
		for _, obj := range o.x {
			if obj.ltype != urls {
				objects.x = append(objects.x, obj)
			}
		}
	default:
		objects = Objects{Parms: o.Parms}
	}
	return &objects
}

// Find returns the int position of an Objects' element
func (o *Objects) Find(elem string) int {
	for i, obj := range o.x {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

// includes returns an io.Reader of blacklist includes
func (o *object) includes() io.Reader {
	sort.Strings(o.inc)
	return strings.NewReader(strings.Join(o.inc, "\n"))
}

func (o *object) isExclude() bool {
	switch o.name {
	case ExcDomns, ExcHosts, ExcRoots:
		return true
	}
	return false
}

// Names returns a sorted slice of Objects names
func (o *Objects) Names() (s sort.StringSlice) {
	for _, obj := range o.x {
		s = append(s, obj.name)
	}
	sort.Sort(s)
	return s
}

func newObject() *object {
	return &object{
		Objects: Objects{},
		exc:     make([]string, 0),
		inc:     make([]string, 0),
	}
}

// Stringer for Object
func (o *object) String() (s string) {
	s += fmt.Sprintf("\nDesc:\t %q\n", o.desc)
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
	return fmt.Sprint(o.x)
}

// Implement Sort Interface for Objects
func (o *Objects) Len() int           { return len(o.x) }
func (o *Objects) Less(i, j int) bool { return o.x[i].name < o.x[j].name }
func (o *Objects) Swap(i, j int)      { o.x[i], o.x[j] = o.x[j], o.x[i] }

package edgeos

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Object struct for normalizing EdgeOS data.
type Object struct {
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
	S []*Object
}

func (o *Objects) addObj(c *Config, node string) {
	switch obj := c.addInc(node); obj {
	case nil:
		o.S = append(o.S, c.bNodes.validate(node).S...)
	default:
		o.S = append(o.S, obj)
		o.S = append(o.S, c.bNodes.validate(node).S...)
	}
}

// Excludes returns an io.Reader of blacklist Includes
func (o *Object) Excludes() io.Reader {
	sort.Strings(o.exc)
	return strings.NewReader(strings.Join(o.exc, "\n"))
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

// Find returns the int position of an Objects' element in the StringSlice
func (o *Objects) Find(elem string) int {
	for i, obj := range o.S {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

// Includes returns an io.Reader of blacklist Includes
func (o *Object) Includes() io.Reader {
	sort.Strings(o.inc)
	return strings.NewReader(strings.Join(o.inc, "\n"))
}

// Names returns a sorted slice of Objects names
func (o *Objects) Names() (s sort.StringSlice) {
	for _, obj := range o.S {
		s = append(s, obj.name)
	}
	sort.Sort(s)
	return s
}

func newObject() *Object {
	return &Object{
		Objects: Objects{},
		exc:     make([]string, 0),
		inc:     make([]string, 0),
	}
}

// Stringer for Object
func (o *Object) String() (r string) {
	r += fmt.Sprintf("\nDesc:\t %q\n", o.desc)
	r += fmt.Sprintf("Disabled: %v\n", o.disabled)
	r += fmt.Sprintf("File:\t %q\n", o.file)
	r += fmt.Sprintf("IP:\t %q\n", o.ip)
	r += fmt.Sprintf("Ltype:\t %q\n", o.ltype)
	r += fmt.Sprintf("Name:\t %q\n", o.name)
	r += fmt.Sprintf("nType:\t %q\n", o.nType)
	r += fmt.Sprintf("Prefix:\t %q\n", o.prefix)
	r += fmt.Sprintf("Type:\t %q\n", getType(o.nType))
	r += fmt.Sprintf("URL:\t %q\n", o.url)
	return r
}

// Stringer for Objects
func (o *Objects) String() string {
	return fmt.Sprint(o.S)
}

// Implement Sort Interface for Objects
func (o *Objects) Len() int           { return len(o.S) }
func (o *Objects) Less(i, j int) bool { return o.S[i].name < o.S[j].name }
func (o *Objects) Swap(i, j int)      { o.S[i], o.S[j] = o.S[j], o.S[i] }

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
	objects
	prefix string
	r      io.Reader
	url    string
}

// objects is a struct of []*Object
type objects struct {
	*Parms
	obs []*object
}

func (o *objects) addObj(c *Config, node string) {
	switch obj := c.addInc(node); obj {
	case nil:
		o.obs = append(o.obs, c.bNodes.validate(node).obs...)
	default:
		o.obs = append(o.obs, obj)
		o.obs = append(o.obs, c.bNodes.validate(node).obs...)
	}
}

// excludes returns an io.Reader of blacklist includes
func (o *object) excludes() io.Reader {
	sort.Strings(o.exc)
	return strings.NewReader(strings.Join(o.exc, "\n"))
}

// Files returns a list of dnsmasq conf files from all srcs
func (o *objects) Files() *CFile {
	c := CFile{Parms: o.Parms}
	for _, obj := range o.obs {
		c.nType = obj.nType
		format := o.Parms.Dir + "/%v.%v." + o.Parms.Ext
		c.names = append(c.names, fmt.Sprintf(format, getType(obj.nType), obj.name))
	}
	sort.Strings(c.names)
	return &c
}

// Find returns the int position of an Objects' element in the StringSlice
func (o *objects) Find(elem string) int {
	for i, obj := range o.obs {
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

// Names returns a sorted slice of Objects names
func (o *objects) Names() (s sort.StringSlice) {
	for _, obj := range o.obs {
		s = append(s, obj.name)
	}
	sort.Sort(s)
	return s
}

func newObject() *object {
	return &object{
		objects: objects{},
		exc:     make([]string, 0),
		inc:     make([]string, 0),
	}
}

// Stringer for Object
func (o *object) String() (r string) {
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
func (o *objects) String() string {
	return fmt.Sprint(o.obs)
}

// Implement Sort Interface for Objects
func (o *objects) Len() int           { return len(o.obs) }
func (o *objects) Less(i, j int) bool { return o.obs[i].name < o.obs[j].name }
func (o *objects) Swap(i, j int)      { o.obs[i], o.obs[j] = o.obs[j], o.obs[i] }

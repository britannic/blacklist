package edgeos

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

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
	return bytes.NewBuffer([]byte(strings.Join(o.inc, "\n")))
}

func newObject() *Object {
	return &Object{
		Objects: Objects{},
		exc:     make([]string, 0),
		inc:     make([]string, 0),
	}
}

// Names returns a sorted slice of Objects names
func (o *Objects) Names() (s sort.StringSlice) {
	for _, obj := range o.S {
		s = append(s, obj.name)
	}
	sort.Sort(s)
	return s
}

// Source returns a map of sources
func (o *Objects) Source(ltype string) *Objects {
	objs := Objects{Parms: o.Parms}
	for _, obj := range o.S {
		switch ltype {
		case obj.ltype:
			objs.S = append(objs.S, obj)
		case "all":
			// objs.S = append(objs.S, obj)
			return o
		}
	}
	return &objs
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

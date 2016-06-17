package edgeos

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Excludes returns a List map of blacklist exclusions
func (o *Object) Excludes() List {
	return UpdateList(o.exc)
}

// Includes returns an io.Reader of blacklist Includes
func (o *Object) Includes() io.Reader {
	sort.Strings(o.inc)
	return bytes.NewBuffer([]byte(strings.Join(o.inc, "\n")))
}

func newObject() *Object {
	return &Object{
		data: make(data),
		exc:  make([]string, 0),
		inc:  make([]string, 0),
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

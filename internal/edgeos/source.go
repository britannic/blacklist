// Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.
package edgeos

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// source struct for normalizing EdgeOS data.
type source struct {
	*Env
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

func (o *source) addLeaf(srcName [][]byte, tnode string) {
	if bytes.Equal(srcName[1], []byte(src)) {
		o.name = string(srcName[2])
		o.nType = getType(tnode).(ntype)
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

// includes returns an io.Reader of blacklist includes
func (o *source) includes() io.Reader {
	sort.Strings(o.inc)
	return strings.NewReader(strings.Join(o.inc, "\n"))
}

func isntSource(nodes []string) bool {
	if len(nodes) == 0 {
		return true
	}
	return nodes[len(nodes)-1] != src
}

func newSource() *source {
	return &source{
		Objects: Objects{},
		exc:     []string{},
		inc:     []string{},
	}
}

func (o *source) setFilePrefix(format string) string {
	switch o.nType {
	case excDomn, excRoot, preDomn:
		return fmt.Sprintf(format, domains, o.name)
	case excHost, preHost:
		return fmt.Sprintf(format, hosts, o.name)
	}
	return fmt.Sprintf(format, getType(o.nType), o.name)
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

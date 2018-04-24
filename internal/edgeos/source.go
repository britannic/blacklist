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

func (s *source) addLeaf(srcName [][]byte, node string) {
	if bytes.Equal(srcName[1], []byte(src)) {
		s.name = string(srcName[2])
		s.nType = getType(node).(ntype)
	}
}

func (s *source) area() string {
	switch getType(s.nType).(string) {
	case domains, PreDomns:
		return domains
	}
	return hosts
}

// excludes returns an io.Reader of blacklist includes
func (s *source) excludes() io.Reader {
	sort.Strings(s.exc)
	return strings.NewReader(strings.Join(s.exc, "\n"))
}

func (s *source) filename(area string) string {
	switch s.nType {
	case excRoot, preRoot:
		return fmt.Sprintf(s.FnFmt, s.Dir, roots, s.name, s.Ext)
	case excDomn, preDomn:
		return fmt.Sprintf(s.FnFmt, s.Dir, domains, s.name, s.Ext)
	case excHost, preHost:
		return fmt.Sprintf(s.FnFmt, s.Dir, hosts, s.name, s.Ext)
	}
	return fmt.Sprintf(s.FnFmt, s.Dir, area, s.name, s.Ext)
}

// includes returns an io.Reader of blacklist includes
func (s *source) includes() io.Reader {
	sort.Strings(s.inc)
	return strings.NewReader(strings.Join(s.inc, "\n"))
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

func (s *source) setFilePrefix(format string) string {
	switch s.nType {
	case excDomn, preDomn:
		return fmt.Sprintf(format, domains, s.name)
	case excHost, preHost:
		return fmt.Sprintf(format, hosts, s.name)
	case excRoot, preRoot:
		return fmt.Sprintf(format, roots, s.name)
	}
	return fmt.Sprintf(format, getType(s.nType), s.name)
}

func printArray(a []string) (s string) {
	if len(a) == 0 {
		s = fmt.Sprintf("              %q\n", "**No entries found**")
		return s
	}
	for _, e := range a {
		s += fmt.Sprintf("              %q\n", e)
	}
	return s
}

func pad(s string) string {
	return fmt.Sprintf("%s %-*s", s, 13-len(s), " ")
}

// Stringer for *source
func (s *source) String() string {
	a := func(s string) string {
		if s == "" {
			return "**Undefined**"
		}
		return s
	}

	str := fmt.Sprintf("\n%s%q\n", pad("Desc:"), a(s.desc))
	str += fmt.Sprintf("%s\"%v\"\n", pad("Disabled:"), s.disabled)
	str += fmt.Sprintf("%s%q\n", pad("File:"), a(s.file))
	str += fmt.Sprintf("%s%q\n", pad("IP:"), a(s.ip))
	str += fmt.Sprintf("%s%q\n", pad("Ltype:"), a(s.ltype))
	str += fmt.Sprintf("%s%q\n", pad("Name:"), a(s.name))
	str += fmt.Sprintf("%s%q\n", pad("nType:"), s.nType)
	str += fmt.Sprintf("%s%q\n", pad("Prefix:"), a(s.prefix))
	str += fmt.Sprintf("%s%q\n", pad("Type:"), getType(s.nType))
	str += fmt.Sprintf("%s%q\n", pad("URL:"), a(s.url))
	str += fmt.Sprintf("Whitelist:\n%s", printArray(s.exc))
	str += fmt.Sprintf("Blacklist:\n%s", printArray(s.inc))
	return str
}

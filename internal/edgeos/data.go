package edgeos

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
)

// ntype for labeling blacklist source types
type ntype int

//go:generate stringer -type=ntype
// ntype label blacklist source types
const (
	unknown ntype = iota // denotes a coding error
	domn                 // Format type e.g. address=/d.com/0.0.0.0
	excDomn              // Excluded from domains
	excHost              // Excluded from hosts
	excRoot              // Excluded globally
	host                 // Format type e.g. server=/www.d.com/0.0.0.0
	preDomn              // Pre-configured blacklisted domains
	preHost              // Pre-configured blacklisted hosts
	preRoot              // Pre-configured global blacklist domains
	root                 // Topmost root node
)

// booltoStr converts a boolean ("true" or "false") to a string equivalent
func booltoStr(b bool) string {
	return fmt.Sprintf("%t", b)
}

// diffArray returns the delta of two arrays
func diffArray(a, b []string) (diff sort.StringSlice) {
	var largest, smallest []string
	switch {
	case len(a) > len(b), len(a) == len(b):
		largest, smallest = a, b
	case len(a) < len(b):
		largest, smallest = b, a
	}

	d := list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
	for _, k := range smallest {
		d.set([]byte(k), 0)
	}

	for _, k := range largest {
		if !d.keyExists([]byte(k)) {
			diff = append(diff, k)
		}
	}

	diff.Sort()
	return diff
}

// formatData returns an io.Reader loaded with dnsmasq formatted data
func formatData(s string, l list) io.Reader {
	var ls sort.StringSlice
	l.RLock()

	for k := range l.entry {
		ls = append(ls, fmt.Sprintf(s+"\n", k))
	}

	ls.Sort()
	l.RUnlock()
	return strings.NewReader(strings.Join(ls, ""))
}

// getDnsmasqPrefix returns the dnsmasq conf file delimiter
func getDnsmasqPrefix(s *source) string {
	switch s.nType {
	case domn, preDomn, preRoot, root:
		return s.Pfx.domain + "/%v/" + s.ip
	case excDomn, excHost, excRoot:
		return s.Pfx.host + "/%v/#"
	}
	return s.Pfx.domain + "/%v/" + s.ip
}

// getType returns the converted "in" type
func getType(in interface{}) (out interface{}) {
	switch in.(type) {
	case ntype:
		out = typeInt(in.(ntype))
	case string:
		out = typeStr(in.(string))
	}
	return out
}

// Iter iterates over ints - use it in for loops
func Iter(i int) []struct{} {
	return make([]struct{}, i)
}

// NewWriter returns an io.Writer
func NewWriter() io.Writer {
	return bufio.NewWriter(&bytes.Buffer{})
}

// strToBool converts a string ("true" or "false") to boolean
func strToBool(s string) bool {
	return strings.ToLower(s) == True
}

func typeInt(n ntype) string {
	switch n {
	case domn:
		return domains
	case excDomn:
		return ExcDomns
	case excHost:
		return ExcHosts
	case excRoot:
		return ExcRoots
	case host:
		return hosts
	case preDomn:
		return PreDomns
	case preHost:
		return PreHosts
	case preRoot:
		return PreRoots
	case root:
		return rootNode
	}
	return notknown
}

func typeStr(s string) ntype {
	switch s {
	case domains:
		return domn
	case ExcDomns:
		return excDomn
	case ExcHosts:
		return excHost
	case ExcRoots:
		return excRoot
	case hosts:
		return host
	case PreDomns:
		return preDomn
	case PreHosts:
		return preHost
	case PreRoots:
		return preRoot
	case rootNode:
		return root
	}
	return unknown
}

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
	root                 // Topmost root node
	zone                 // Unused - future application
)

// booltoStr converts a boolean ("true" or "false") to a string equivalent
func booltoStr(b bool) string {
	if b {
		return True
	}
	return False
}

// diffArray returns the delta of two arrays
func diffArray(a, b []string) (diff sort.StringSlice) {
	var biggest, smallest []string
	switch {
	case len(a) > len(b), len(a) == len(b):
		biggest, smallest = a, b
	case len(a) < len(b):
		biggest, smallest = b, a
	}

	dmap := list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
	for _, k := range smallest {
		dmap.set([]byte(k), 0)
	}

	for _, k := range biggest {
		if !dmap.keyExists([]byte(k)) {
			diff = append(diff, k)
		}
	}

	diff.Sort()
	return diff
}

// formatData returns an io.Reader loaded with dnsmasq formatted data
func formatData(fmttr string, l list) io.Reader {
	var lines sort.StringSlice
	l.RLock()

	for k := range l.entry {
		lines = append(lines, fmt.Sprintf(fmttr+"\n", k))
	}

	lines.Sort()
	l.RUnlock()
	return strings.NewReader(strings.Join(lines, ""))
}

// getDnsmasqPrefix returns the dnsmasq conf file delimiter
func getDnsmasqPrefix(o *object) string {
	switch o.nType {
	case domn, preDomn, root:
		return o.Pfx.domain + "/%v/" + o.ip
	case excDomn, excHost, excRoot:
		return o.Pfx.host + "/%v/#"
	}
	return o.Pfx.host + "/%v/" + o.ip
}

// getSubdomains returns a map of subdomains
func getSubdomains(b []byte) (l list) {
	l.entry = make(entry)
	keys := bytes.Split(b, []byte("."))
	for i := range Iter(len(keys) - 1) {
		k := bytes.Join(keys[i:], []byte("."))
		l.entry[string(k)] = 0
	}
	return l
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
	var b bytes.Buffer
	return bufio.NewWriter(&b)
}

// logIt writes to io.Writer
func logIt(w io.Writer, s string) {
	io.Copy(w, strings.NewReader(s))
}

// strToBool converts a string ("true" or "false") to boolean
func strToBool(s string) bool {
	return strings.ToLower(s) == True
}

func typeInt(n ntype) (s string) {
	switch n {
	case domn:
		s = domains
	case excDomn:
		s = ExcDomns
	case excHost:
		s = ExcHosts
	case excRoot:
		s = ExcRoots
	case host:
		s = hosts
	case preDomn:
		s = PreDomns
	case preHost:
		s = PreHosts
	case root:
		s = rootNode
	case unknown:
		s = notknown
	case zone:
		s = zones
	}
	return s
}

func typeStr(s string) (n ntype) {
	switch s {
	case domains:
		n = domn
	case ExcDomns:
		n = excDomn
	case ExcHosts:
		n = excHost
	case ExcRoots:
		n = excRoot
	case hosts:
		n = host
	case notknown:
		n = unknown
	case PreDomns:
		n = preDomn
	case PreHosts:
		n = preHost
	case rootNode:
		n = root
	case zones:
		n = zone
	}
	return n
}

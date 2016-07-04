package edgeos

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ntype for labeling blacklist source types
type ntype int

//go:generate stringer -type=ntype
// ntypes label blacklist source types
const (
	unknown ntype = iota // denotes a coding error
	domn                 // Format type e.g. address=/.d.com/0.0.0.0
	excDomn              // Won't be written to disk
	excHost              // Won't be written to disk
	excRoot              // Won't be written to disk
	host                 // Format type e.g. address=/www.d.com/0.0.0.0
	preDomn              // Pre-configured backlisted domains
	preHost              // Pre-configured backlisted hosts
	root                 // Topmost root node
	zone                 // Unused - future application
)

// BooltoStr converts a boolean ("true" or "false") to a string equivalent
func BooltoStr(b bool) string {
	if b {
		return True
	}
	return False
}

// DiffArray returns the delta of two arrays
func DiffArray(a, b []string) (diff sort.StringSlice) {
	var biggest, smallest []string

	switch {
	case len(a) > len(b), len(a) == len(b):
		biggest = a
		smallest = b
	case len(a) < len(b):
		biggest = b
		smallest = a
	}
	dmap := make(list)
	for _, k := range smallest {
		dmap[k] = 0
	}
	for _, k := range biggest {
		if !dmap.keyExists(k) {
			diff = append(diff, k)
		}
	}
	diff.Sort()
	return diff
}

// formatData returns an io.Reader loaded with dnsmasq formatted data
func formatData(fmttr string, data list) io.Reader {
	var lines sort.StringSlice
	for k := range data {
		lines = append(lines, fmt.Sprintf(fmttr+"\n", k))
	}
	lines.Sort()

	return strings.NewReader(strings.Join(lines, ""))
}

// getSeparator returns the dnsmasq conf file delimiter
func getSeparator(node string) string {
	if node == domains {
		return "/."
	}
	return "/"
}

// getSubdomains returns a map of subdomains
func getSubdomains(s string) (l list) {
	l = make(list)
	keys := strings.Split(s, ".")
	for i := 0; i < len(keys)-1; i++ {
		key := strings.Join(keys[i:], ".")
		l[key] = 0
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

// StrToBool converts a string ("true" or "false") to boolean
func StrToBool(s string) bool {
	if strings.ToLower(s) == True {
		return true
	}
	return false
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

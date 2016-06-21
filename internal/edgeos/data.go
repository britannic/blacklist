package edgeos

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
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

	dmap := make(List)
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
func formatData(fmttr string, data List) io.Reader {
	var lines sort.StringSlice
	for k := range data {
		lines = append(lines, fmt.Sprintf(fmttr+"\n", k))
	}
	lines.Sort()
	return bytes.NewBuffer([]byte(strings.Join(lines, "")))
}

// getSeparator returns the dnsmasq conf file delimiter
func getSeparator(node string) string {
	if node == domains {
		return "/."
	}
	return "/"
}

// getSubdomains returns a map of subdomains
func getSubdomains(s string) (l List) {
	l = make(List)
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

// StrToBool converts a string ("true" or "false") to it's boolean equivalent
func StrToBool(s string) bool {
	if strings.ToLower(s) == True {
		return true
	}
	return false
}

func typeInt(i ntype) (s string) {
	switch i {
	case domain:
		s = domains
	case host:
		s = hosts
	case pre:
		s = preConf
	case root:
		s = blacklist
	case unknown:
		s = notknown
	case zone:
		s = zones
	}
	return s
}

func typeStr(s string) (i ntype) {
	switch s {
	case domains:
		i = domain
	case hosts:
		i = host
	case preConf:
		i = pre
	case blacklist:
		i = root
	case notknown:
		i = unknown
	case zones:
		i = zone
	}
	return i
}

package edgeos

import (
	"fmt"
	"sort"
	"strings"
)

// List is a map of int
type List map[string]int

// keyExists returns true if the key exists
func (l List) keyExists(s string) bool {
	_, ok := l[s]
	return ok
}

// mergeList combines two List maps
func mergeList(a, b List) List {
	for k, v := range a {
		b[k] = v
	}
	return b
}

// String implements fmt.Print interface
func (l List) String() string {
	var lines sort.StringSlice
	for k, v := range l {
		lines = append(lines, fmt.Sprintf("%q:%v,\n", k, v))
	}
	lines.Sort()
	return strings.Join(lines, "")
}

// SubKeyExists returns true if part of all of the key matches
func (l List) subKeyExists(s string) bool {
	keys := getSubdomains(s)
	for k := range keys {
		if l.keyExists(k) {
			return true
		}
	}
	return l.keyExists(s)
}

// UpdateList converts []string to map of List
func UpdateList(data []string) (l List) {
	l = make(List)
	for _, k := range data {
		l[k] = 0
	}
	return l
}

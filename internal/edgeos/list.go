package edgeos

import (
	"fmt"
	"sort"
	"strings"
)

// list is a map of int
type list map[string]int

// keyExists returns true if the list key exists
func (l list) keyExists(s string) bool {
	_, ok := l[s]
	return ok
}

// mergeList combines two list maps
func mergeList(a, b list) list {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func (l list) String() string {
	var lines sort.StringSlice
	for k, v := range l {
		lines = append(lines, fmt.Sprintf("%q:%v,\n", k, v))
	}
	lines.Sort()
	return strings.Join(lines, "")
}

// SubKeyExists returns true if part of all of the key matches
func (l list) subKeyExists(s string) bool {
	keys := getSubdomains(s)
	for k := range keys {
		if l.keyExists(k) {
			return true
		}
	}
	return l.keyExists(s)
}

// updateList converts []string to map of List
func updateList(data []string) (l list) {
	l = make(list)
	for _, k := range data {
		l[k] = 0
	}
	return l
}

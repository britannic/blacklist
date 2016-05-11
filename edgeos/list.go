package edgeos

import (
	"fmt"
	"sort"
	"strings"
)

// List is a map of int
type List map[string]int

// Lister implements List methods
type Lister interface {
	KeyExists(s string) bool
	String() string
	SubKeyExists(s string) bool
}

// Keys is used for sorting operations on map keys
type Keys []string

// KeyExists returns true if the key exists
func (l List) KeyExists(s string) bool {
	_, ok := l[s]
	return ok
}

// Len returns length of Keys
func (k Keys) Len() int { return len(k) }

// Less returns the smallest element
func (k Keys) Less(i, j int) bool { return k[i] < k[j] }

// Swap swaps elements of a key array
func (k Keys) Swap(i, j int) { k[i], k[j] = k[j], k[i] }

// SortKeys returns an array of sorted strings
func (n Nodes) SortKeys() (pkeys Keys) {
	for pkey := range n {
		pkeys = append(pkeys, pkey)
	}
	sort.Sort(Keys(pkeys))
	return pkeys
}

// SortSKeys returns an array of sorted strings
func (n Nodes) SortSKeys(node string) (skeys Keys) {
	for skey := range n[node].Data {
		skeys = append(skeys, skey)
	}
	sort.Sort(Keys(skeys))
	return skeys
}

// MergeList combines two List maps
func MergeList(a, b List) List {
	for k, v := range a {
		b[k] = v
	}
	return b
}

func (l List) String() string {
	var lines []string
	for k, v := range l {
		lines = append(lines, fmt.Sprintf("%q:%v,\n", k, v))
	}
	sort.Strings(lines)
	return strings.Join(lines, "")
}

// SubKeyExists returns true if part of all of the key matches
func (l List) SubKeyExists(s string) bool {
	keys := GetSubdomains(s)
	for k := range keys {
		if l.KeyExists(k) {
			return true
		}
	}
	return l.KeyExists(s)
}

package edgeos

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type entry map[string]int

// list is a struct map of int
type list struct {
	*sync.RWMutex
	entry
}

// inc increments the entry by +1
// func (l list) inc(k string) {
// 	l.Lock()
// 	l.entry[k]++
// 	l.Unlock()
// }

// set sets the int value of entry
func (l list) keyExists(s string) bool {
	l.RLock()
	defer l.RUnlock()
	_, ok := l.entry[s]
	return ok
}

// keyExists returns true if the list key exists
func mergeList(a, b list) list {
	a.Lock()
	b.Lock()
	defer a.Unlock()
	defer b.Unlock()
	for k, v := range b.entry {
		a.entry[k] = v
	}
	return a
}

// mergeList combines two list maps
func (l list) set(k string, v int) {
	l.Lock()
	l.entry[k] = v
	l.Unlock()
}

func (l list) String() string {
	var lines sort.StringSlice
	for k, v := range l.entry {
		lines = append(lines, fmt.Sprintf("%q:%v,\n", k, v))
	}
	lines.Sort()
	return strings.Join(lines, "")
}

// subKeyExists returns true if part of all of the key matches
func (l list) subKeyExists(s string) bool {
	keys := getSubdomains(s)
	for k := range keys.entry {
		if l.keyExists(k) {
			return true
		}
	}
	return l.keyExists(s)
}

// updateEntry converts []string to map of List
func updateEntry(data []string) (l list) {
	l.entry = make(entry)
	for _, k := range data {
		l.entry[k] = 0
	}
	return l
}

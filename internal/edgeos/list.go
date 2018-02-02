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

// set sets the int value of entry
func (l list) keyExists(k []byte) bool {
	l.RLock()
	_, ok := l.entry[string(k)]
	l.RUnlock()
	return ok
}

// keyExists returns true if the list key exists
func mergeList(a, b list) list {
	a.Lock()
	b.Lock()

	for k, v := range b.entry {
		a.entry[k] = v
	}

	a.Unlock()
	b.Unlock()
	return a
}

// mergeList combines two list maps
func (l list) set(k []byte, v int) {
	l.Lock()
	l.entry[string(k)] = v
	l.Unlock()
}

func (l list) String() string {
	var lines sort.StringSlice
	for k, v := range l.entry {
		lines = append(lines, fmt.Sprintf("%q:%v,\n", string(k), v))
	}
	lines.Sort()
	return strings.Join(lines, "")
}

// subKeyExists returns true if part of all of the key matches
func (l list) subKeyExists(k []byte) bool {
	keys := getSubdomains(k)
	for k := range keys.entry {
		if l.keyExists([]byte(k)) {
			return true
		}
	}
	return l.keyExists(k)
}

// updateEntry converts []string to map of List
func updateEntry(data [][]byte) (l list) {
	l.entry = make(entry)
	for _, k := range data {
		l.entry[string(k)] = 0
	}
	return l
}

package edgeos

import (
	"bytes"
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

// merge returns a merge of two lists
func (l list) merge(a list) {
	l.Lock()
	for k, v := range a.entry {
		l.entry[k] = v
	}
	l.Unlock()
}

// set adds a list entry map member
func (l list) set(k []byte, v int) {
	l.Lock()
	l.entry[string(k)] = v
	l.Unlock()
}

func (l list) String() string {
	var ls sort.StringSlice
	for k, v := range l.entry {
		ls = append(ls, fmt.Sprintf("%q:%v,\n", string(k), v))
	}
	ls.Sort()
	return strings.Join(ls, "")
}

// subKeyExists returns true if part or all of the key matches
func (l list) subKeyExists(b []byte) bool {
	d := bytes.Split(b, []byte("."))
	for i := range Iter(len(d) - 1) {
		if l.keyExists(bytes.Join(d[i:], []byte("."))) {
			return true
		}
	}
	return l.keyExists(b)
}

// updateEntry converts [][]byte to map of List
func updateEntry(d [][]byte) (l list) {
	l.entry = make(entry)
	for _, k := range d {
		l.entry[string(k)] = 0
	}
	return l
}

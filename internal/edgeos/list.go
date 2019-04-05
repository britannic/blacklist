package edgeos

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type entry map[string]struct{}

// list is a struct map of entry with a RW Mutex
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
func (l list) set(k []byte) {
	l.Lock()
	l.entry[string(k)] = struct{}{}
	l.Unlock()
}

func (l list) String() string {
	var (
		i  int64
		ls = make(sort.StringSlice, len(l.entry))
	)
	for k, v := range l.entry {
		ls[i] = fmt.Sprintf("%q:%v,\n", string(k), v)
		i++
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
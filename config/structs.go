package config

import "strings"

// Blacklist type is a map of Nodes with string keys
type Blacklist map[string]*Node

// Source is a map of Srcs with string keys
type Source map[string]*Src

// Dict is a common string key map of ints
type Dict map[string]int

// GetSubdomains returns a map of subdomains
func GetSubdomains(s string) (d Dict) {
	d = make(Dict)
	keys := strings.Split(s, ".")
	for i := 0; i < len(keys)-1; i++ {
		key := strings.Join(keys[i:], ".")
		d[key] = 0
	}
	return
}

// SubKeyExists returns true if part of the key matches
func (d Dict) SubKeyExists(s string) bool {
	keys := GetSubdomains(s)
	for k := range keys {
		if d.KeyExists(k) {
			return true
		}
	}
	if d.KeyExists(s) {
		return true
	}
	return false
}

// KeyExists returns true if the key exists
func (d Dict) KeyExists(s string) bool {
	if _, exist := d[s]; exist {
		return true
	}
	return false
}

// Node configuration record
type Node struct {
	Disable          bool
	IP               string
	Exclude, Include []string
	Source           Source
}

// Src record struct for Source map
type Src struct {
	Desc    string
	Disable bool
	IP      string
	List    Dict
	Name    string
	No      int
	Prfx    string
	Type    string
	URL     string
}

// Keys is used for sorting operations on map keys
type Keys []string

// Len returns length of Keys
func (k Keys) Len() int { return len(k) }

// Swap swaps elements of a key array
func (k Keys) Swap(i, j int) { k[i], k[j] = k[j], k[i] }

// Less returns the smallest element
func (k Keys) Less(i, j int) bool { return k[i] < k[j] }

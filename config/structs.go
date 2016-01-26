// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

package config

// Blacklist type is a map of Nodes with string keys
type Blacklist map[string]*Node

// Source is a map of Srcs with string keys
type Source map[string]*Src

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
	List    map[string]int
	Name    string
	Prfx    string
	Type    string
	URL     string
}

// Area struct holds data on downloaded hosts and domains
type Area struct {
	Cntr, Dupe, Rcrd, Uniq int
	Trgt                   map[string]string
}

// Keys is used for sorting operations on map keys
type Keys []string

// Len returns length of Keys
func (k Keys) Len() int { return len(k) }

// Swap swaps elements of a key array
func (k Keys) Swap(i, j int) { k[i], k[j] = k[j], k[i] }

// Less returns the smallest element
func (k Keys) Less(i, j int) bool { return k[i] < k[j] }

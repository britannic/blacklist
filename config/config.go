// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package config provides methods and data structures for loading
// an EdgeOS/VyOS configuration
package config

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	r "github.com/britannic/blacklist/regx"
)

// API sets the path and executable for the EdgeOS shell API
const (
	API = "/bin/cli-shell-api"
)

// SortKeys returns an array of sorted strings
func (b Blacklist) SortKeys() (pkeys Keys) {
	for pkey := range b {
		pkeys = append(pkeys, pkey)
	}
	sort.Sort(Keys(pkeys))
	return
}

// SortSKeys returns an array of sorted strings
func (b Blacklist) SortSKeys() (skeys Keys) {
	for _, pkey := range b.SortKeys() {
		for skey := range b[pkey].Source {
			skeys = append(skeys, skey)
		}
	}
	sort.Sort(Keys(skeys))
	return
}

// String returns pretty print for the Blacklist struct
func (b Blacklist) String() (result string) {
	// cols, _, _ := terminal.GetSize(int(os.Stdout.Fd()))

	for _, pkey := range b.SortKeys() {
		result += fmt.Sprintf("Node: %v\n\tDisabled: %v\n\t", pkey, b[pkey].Disable)

		if k := b[pkey].IP; len(k) > 0 {
			result += fmt.Sprintf("Redirect IP: %v\n\t", b[pkey].IP)
		}

		if k := b[pkey].Exclude; len(k) > 1 {
			result += "Exclude(s):\n"
			for _, exclude := range b[pkey].Exclude {
				result += fmt.Sprintf("\t\t%v\n", exclude)
			}
		}

		if k := b[pkey].Include; len(k) > 1 {
			result += "\tInclude(s):\n"
			for _, include := range b[pkey].Include {
				result += fmt.Sprintf("\t\t%v\n", include)
			}
		}

		for _, skey := range b.SortSKeys() {
			if _, ok := b[pkey].Source[skey]; ok {
				result += fmt.Sprintf("\tSource: %v\n\t\tDisabled: %v\n\t\tDescription: %v\n\t\tPrefix: %q\n\t\tURL: %v\n", skey, b[pkey].Source[skey].Disable, b[pkey].Source[skey].Desc, b[pkey].Source[skey].Prfx, b[pkey].Source[skey].URL)
			}
		}
		result += "\n"
	}
	return
}

// APICmd returns a map of CLI commands
func APICmd() (r map[string]string) {
	r = make(map[string]string)
	l := []string{
		"cfExists",
		"cfReturnValue",
		"cfReturnValues",
		"exists",
		"existsActive",
		"getNodeType",
		"inSession",
		"isLeaf",
		"isMulti",
		"isTag",
		"listActiveNodes",
		"listNodes",
		"returnActiveValue",
		"returnActiveValues",
		"returnValue",
		"returnValues",
		"showCfg",
		"showConfig",
	}
	for _, k := range l {
		r[k] = fmt.Sprintf("%v %v", API, k)
	}
	return
}

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

// Dict is a common string key map of ints
type Dict map[string]int

// KeyExists returns true if the key exists
func (d Dict) KeyExists(s string) bool {
	if _, exist := d[s]; exist {
		return true
	}
	return false
}

// SubKeyExists returns true if part of all of the key matches
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

// Keys is used for sorting operations on map keys
type Keys []string

// Len returns length of Keys
func (k Keys) Len() int { return len(k) }

// Swap swaps elements of a key array
func (k Keys) Swap(i, j int) { k[i], k[j] = k[j], k[i] }

// Less returns the smallest element
func (k Keys) Less(i, j int) bool { return k[i] < k[j] }

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
	List    Dict
	Name    string
	No      int
	Prfx    string
	Type    string
	URL     string
}

// ToBool converts a string ("true" or "false") to it's boolean equivalent
func ToBool(s string) bool {

	if strings.ToLower(s) == "true" {
		return true
	}
	return false
}

// SHcmd returns the appropriate command for non-tty or tty context
func SHcmd(a string) (action string) {
	if !Insession() {
		action = a
		switch a {
		case "listNodes":
			action = "listActiveNodes"
		case "returnValue":
			action = "returnActiveValue"
		case "returnValues":
			action = "returnActiveValues"
		case "exists":
			action = "existsActive"
		case "showConfig":
			action = "showCfg"
		}
	}
	return
}

// Load reads the config using the EdgeOS/VyOS cli-shell-api
func Load(action string, level string) (r string, err error) {
	action = SHcmd(action)
	x := APICmd()
	if _, ok := x[action]; !ok {
		return r, fmt.Errorf("API command %v is invalid", level)
	}

	cmd := exec.Command("/bin/bash")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%v %v", x[action], level))

	stdout, err := cmd.Output()
	if err != nil {
		return
	}

	r = string(stdout)
	return r, err
}

// Insession returns true if VyOS/EdgeOS configuration is in session
func Insession() bool {
	cmd := exec.Command(API, "inSession")
	var out bytes.Buffer
	cmd.Stdout = &out
	if ok := cmd.Run(); ok == nil {
		if out.String() == "0" {
			return true
		}
	}
	return false
}

// Get extracts nodes from a EdgeOS/VyOS configuration structure
func Get(cfg string, root string) (*Blacklist, error) {
	var (
		cfgtree     = make(Blacklist)
		err         error
		leaf, tnode string
	)

	if len(cfg) < 1 {
		err = errors.New("Configuration data is empty, cannot continue")
		return &cfgtree, err
	}

	nodes := make([]string, 5)
	rx := r.Regex

LINE:
	for _, line := range strings.Split(cfg, "\n") {
		line = strings.TrimSpace(line)

		switch {
		case rx.MLTI.MatchString(line):
			IncExc := r.Get("mlti", line)
			switch IncExc[1] {
			case "exclude":
				cfgtree[tnode].Exclude = append(cfgtree[tnode].Exclude, IncExc[2])
			case "include":
				cfgtree[tnode].Include = append(cfgtree[tnode].Include, IncExc[2])
			}

		case rx.NODE.MatchString(line):
			node := r.Get("node", line)
			tnode = node[1]
			cfgtree[tnode] = &Node{}
			cfgtree[tnode].Source = make(Source)
			nodes = append(nodes, tnode)

		case rx.LEAF.MatchString(line):
			src := r.Get("leaf", line)
			leaf = src[2]
			nodes = append(nodes, src[1])

			if src[1] == "source" {
				cfgtree[tnode].Source[leaf] = &Src{Type: tnode, Name: leaf}
			}

		case rx.DSBL.MatchString(line):
			cfgtree[tnode].Disable = ToBool(r.Get("dsbl", line)[1])

		case rx.NAME.MatchString(line):
			name := r.Get("name", line)

			switch name[1] {
			case "prefix":
				cfgtree[tnode].Source[leaf].Prfx = name[2]

			case "url":
				cfgtree[tnode].Source[leaf].URL = name[2]

			case "description":
				cfgtree[tnode].Source[leaf].Desc = name[2]

			case "dns-redirect-ip":
				cfgtree[tnode].IP = name[2]
			}

		case rx.DESC.MatchString(line) || rx.CMNT.MatchString(line) || rx.MISC.MatchString(line):
			continue LINE

		case rx.RBRC.MatchString(line):
			nodes = nodes[:len(nodes)-1] // pop last node
			tnode = nodes[len(nodes)-1]

		}
	}

	return &cfgtree, err
}

// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package config provides methods and data structures for loading
// an EdgeOS/VyOS configuration
package config

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"

	r "github.com/britannic/blacklist/regx"
)

const (
	api = "/bin/cli-shell-api"
)

// String returns pretty print for the Blacklist struct
func (b Blacklist) String() (result string) {
	cols, _, _ := terminal.GetSize(int(os.Stdout.Fd()))
	sortPkeys := func() (pkeys Keys) {
		for pkey := range b {
			pkeys = append(pkeys, pkey)
		}
		sort.Sort(Keys(pkeys))
		return
	}

	sortSkeys := func() (skeys Keys) {
		for _, pkey := range sortPkeys() {
			for skey := range b[pkey].Source {
				skeys = append(skeys, skey)
			}
		}
		sort.Sort(Keys(skeys))
		return
	}

	for _, pkey := range sortPkeys() {
		result += fmt.Sprintf("Node: %v\n\tDisabled: %v\n\tRedirect IP: %v\n\tExclude(s):\n", pkey, b[pkey].Disable, b[pkey].IP)
		for _, exclude := range b[pkey].Exclude {
			result += fmt.Sprintf("\t\t%v\n", exclude)
		}
		result += fmt.Sprintf("\tInclude(s):\n")
		for _, include := range b[pkey].Include {
			result += fmt.Sprintf("\t\t%v\n", include)
		}
		for _, skey := range sortSkeys() {
			if _, ok := b[pkey].Source[skey]; ok {
				result += fmt.Sprintf("\tSource: %v\n\t\tDisabled: %v\n\t\tDescription: %v\n\t\tPrefix: %v\n\t\tURL: %v\n", skey, b[pkey].Source[skey].Disable, b[pkey].Source[skey].Desc, b[pkey].Source[skey].Prfx, b[pkey].Source[skey].URL)
			}
		}
		result += fmt.Sprintln(strings.Repeat("-", cols/2))
	}
	return
}

// APICmd returns a map of CLI commands
func APICmd() (r map[string]string) {
	r = make(map[string]string)
	l := []string{"cfExists", "cfReturnValue", "cfReturnValues", "exists", "existsActive", "getNodeType", "inSession", "isLeaf", "isMulti", "isTag", "listActiveNodes", "listNodes", "returnActiveValue", "returnActiveValues", "returnValue", "returnValues", "showCfg", "showConfig"}
	for _, k := range l {
		r[k] = fmt.Sprintf("%v %v", api, k)
	}
	return
}

// ToBool converts a string ("true" or "false") to it's boolean equivalent
func ToBool(s string) (b bool) {
	if len(s) == 0 {
		log.Fatal("ERROR: variable empty, cannot convert to boolean")
	}
	switch s {
	case "false":
		b = false
	case "true":
		b = true
	}
	return
}

// shCmd returns the appropriate command for non-tty or tty context
func shCmd(a string) (action string) {
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
	action = shCmd(action)
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
	return
}

// Insession returns true if VyOS/EdgeOS configuration is in session
func Insession() bool {
	cmd := exec.Command(api, "inSession")
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
func Get(cfg string, root string) (b *Blacklist, err error) {

	if len(cfg) < 1 {
		err = fmt.Errorf("Configuration data is empty, cannot continue")
		return
	}

	cfgtree := make(Blacklist)
	err = nil
	nodes := make([]string, 5)
	rx := r.Regex()
	var leaf string
	var tnode string

	for _, line := range strings.Split(cfg, "\n") {
		line = strings.TrimSpace(line)

		switch {
		case rx.MLTI.MatchString(line):
			{
				IncExc := r.Get("mlti", line)
				switch IncExc[1] {
				case "exclude":
					cfgtree[tnode].Exclude = append(cfgtree[tnode].Exclude, IncExc[2])
				case "include":
					cfgtree[tnode].Include = append(cfgtree[tnode].Include, IncExc[2])
				}
			}
		case rx.NODE.MatchString(line):
			{
				node := r.Get("node", line)
				tnode = node[1]
				cfgtree[tnode] = &Node{}
				cfgtree[tnode].Source = make(Source)
				nodes = append(nodes, tnode)
			}
		case rx.LEAF.MatchString(line):
			src := r.Get("leaf", line)
			leaf = src[2]
			nodes = append(nodes, src[1])

			if src[1] == "source" {
				cfgtree[tnode].Source[leaf] = &Src{Type: tnode, Name: leaf}
			}
		case rx.DSBL.MatchString(line):
			{
				disabled := r.Get("dsbl", line)
				cfgtree[tnode].Disable = ToBool(disabled[1])
			}
		case rx.NAME.MatchString(line):
			{
				name := r.Get("name", line)
				switch name[1] {
				case "prefix":
					cfgtree[tnode].Source[leaf].Prfx = name[2]
				case "url":
					cfgtree[tnode].Source[leaf].URL = name[2]
				case "description":
					cfgtree[tnode].Source[leaf].Desc = name[2]
				case "dns-redirect-ip":
					if name[2] == "" {
						cfgtree[tnode].IP = cfgtree[root].IP
						break
					}
					cfgtree[tnode].IP = name[2]
				}
			}
		case rx.DESC.MatchString(line) || rx.CMNT.MatchString(line) || rx.MISC.MatchString(line):
			break
		case rx.RBRC.MatchString(line):
			{
				nodes = nodes[:len(nodes)-1] // pop last node
				tnode = nodes[len(nodes)-1]
			}
		}
	}

	return &cfgtree, err
}

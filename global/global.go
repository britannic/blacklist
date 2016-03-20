// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package global sets project scoped variables and constants
package global

import (
	"log"
	"os"
	"runtime"

	"github.com/britannic/blacklist/utils"
)

// Areas type string of configuration tree nodes
type Areas struct {
	Domains,
	Hosts,
	Root string
}

var (
	// Area defines top nodes for the configuration tree
	Area = &Areas{
		Domains: "domains",
		Hosts:   "hosts",
		Root:    "blacklist",
	}

	// Dbg sets the Debug flag
	Dbg = false

	// DmsqDir defines dnsmasq directory location
	DmsqDir string

	// Fext defines the blacklist filename extension
	Fext = ".blacklist.conf"

	// FStr provides a blacklist filename/path template
	FStr = "%v/%v.%v" + Fext

	// Logfile set the log path and filename
	Logfile string

	// TestOS sets a global value for the test environment
	TestOS = "darwin"

	// Program is the current binary's filename
	Program = utils.Basename(os.Args[0])

	// Root is the top level configuration Node
	Root = "blacklist"

	// WhatOS is the current operating system
	WhatOS string
)

func init() {
	WhatOS = runtime.GOOS
	switch WhatOS {
	case TestOS:
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Cannot determine current directory - exiting")
		}
		DmsqDir = cwd + "/testdata"
		Logfile = cwd + "/testdata/blacklist.log"
	default:
		DmsqDir = "/etc/dnsmasq.d"
		FStr = "%v/%v.%v" + Fext
		Logfile = "/var/log/blacklist.log"
	}
}

// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package global sets project scoped variables and constants
package global

import (
	"os"
	"runtime"

	"github.com/britannic/blacklist/utils"
)

var (
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

	// Program is the current binary's filename
	Program = utils.Basename(os.Args[0])

	// Root is the top level configuration Node
	Root = "blacklist"

	// WhatOS is the current operating system
	WhatOS = runtime.GOOS
)

func init() {
	WhatOS = runtime.GOOS
	switch WhatOS {
	case "darwin":
		DmsqDir = "/Users/Neil/go/src/github.com/britannic/blacklist/testdata"
		Logfile = "/Users/Neil/go/src/github.com/britannic/blacklist/testdata/blacklist.log"
	default:
		DmsqDir = "/etc/dnsmasq.d"
		FStr = "%v/%v.%v" + Fext
		Logfile = "/var/log/blacklist.log"
	}
}

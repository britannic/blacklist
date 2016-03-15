// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package global sets project scoped variables and constants
package global

import (
	"os"

	"github.com/britannic/blacklist/utils"
)

var (
	// Dbg sets the Debug flag
	Dbg = false

	// DmsqDir defines dnsmasq directory location
	DmsqDir = "/etc/dnsmasq.d"

	// Fext defines the blacklist filename extension
	Fext = ".blacklist.conf"

	// FStr provides a blacklist filename/path template
	FStr = "%v/%v.%v" + Fext

	// Logfile set the log path and filename
	Logfile = "/var/log/blacklist.log"

	// Program is the current binary's filename
	Program = utils.Basename(os.Args[0])

	// Root is the top level configuration Node
	Root = "blacklist"
)

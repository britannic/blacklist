// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package global sets project scoped variables and constants
package global

import (
	"log"
	"os"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/britannic/blacklist/utils"
)

// Areas type string of configuration tree nodes
type Areas struct {
	Domains,
	Hosts,
	Root string
}

const (

	// Fext defines the blacklist filename extension
	Fext = ".blacklist.conf"

	// TestOS sets a global value for the test environment
	TestOS = "darwin"
)

var (
	// Area defines top nodes for the configuration tree
	Area = &Areas{
		Domains: "domains",
		Hosts:   "hosts",
		Root:    "blacklist",
	}

	// Args provides a way to set flags within the unit test
	Args []string

	cwd string

	// Dbg sets the Debug flag
	Dbg = false

	// DNSRestart defines the dnsmasq restart command
	DNSRestart = "service dnsmasq restart"

	// DmsqDir defines dnsmasq directory location
	DmsqDir string

	// FStr provides a blacklist filename/path template
	FStr = "%v/%v.%v" + Fext

	// Log is a pointer to *log.Inst
	Log *logrus.Logger

	// LogFile sets the log path and filename
	LogFile string

	// LogMaxAge sets how old the back up file is before rotation
	LogMaxAge int

	// LogMaxBackups sets how many log file backups there are before purging
	LogMaxBackups int

	// LogMaxSize set the maximum log size
	LogMaxSize int

	// LogOutput sets where to log
	LogOutput = "file"

	// Program is the current binary's filename
	Program = utils.Basename(os.Args[0])

	// TargetArch is the host platforms CPU
	TargetArch = "mips64"

	// WhatArch is the current operating system
	WhatArch string

	// WhatOS is the current operating system
	WhatOS string
)

func init() {
	WhatOS = runtime.GOOS
	WhatArch = runtime.GOARCH
	SetVars(WhatArch)
}

// SetVars conditionally sets global variables based on the current OS
func SetVars(ARCH string) {

	switch ARCH {
	case TargetArch:
		DmsqDir = "/etc/dnsmasq.d"
		FStr = "%v/%v.%v" + Fext
		LogFile = "/var/log/blacklist.log"

	default:
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Cannot determine current directory - exiting")
		}
		LogFile = "/tmp/blacklist.log"
		DmsqDir = cwd + "/testdata"

		DNSRestart = "echo -n dnsmasq not implemented on " + WhatOS
	}
}

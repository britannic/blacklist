// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	c "github.com/britannic/blacklist/config"
)

// root constant defines the top level configuration node
const root = "blacklist"

var (
	build   = "UNKNOWN"
	cores   = runtime.NumCPU()
	dmsqDir = "/etc/dnsmasq.d"
	fSfx    = ".blacklist.conf"
	fStr    = "%v/%v.%v" + fSfx
	githash = "UNKNOWN"
	logfile = "/var/log/blacklist.log"
	program = basename(os.Args[0])
	version = "UNKNOWN"
)

func main() {
	runtime.GOMAXPROCS(cores)

	var (
		poll    = time.Second * 2 // poll
		timeout = time.Minute * 3
	)

	whatOS := runtime.GOOS
	if whatOS == "darwin" {
		dmsqDir = "/tmp"
		logfile = "/tmp/blacklist.log"
	}

	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE, 0755)
	if err == nil {
		log.SetFormatter(&log.TextFormatter{DisableColors: true})
		log.SetOutput(f)
	}

	o := getopts()

	switch {
	case *o.version:
		{
			log.Info("%s version: %s, build date: %s\n hash: %v\n", program, version, build, githash)
			os.Exit(0)
		}

	case *o.verb:
		log.SetFormatter(&log.TextFormatter{DisableColors: false})
		log.SetOutput(os.Stderr)

	case *o.poll != 5:
		poll = time.Duration(*o.poll) * time.Second
		log.Info("Poll duration", poll)
	}

	log.Info("CPU Cores: ", cores)

	blist, err := func() (b *c.Blacklist, err error) {
		switch whatOS {
		case "darwin":
			{
				b, err = c.Get(c.Testdata, root)
				if err != nil {
					return b, fmt.Errorf("unable to get configuration data, error code: %v\n", err)
				}
				return
			}
		default:
			{
				cfg, err := c.Load("showCfg", "service dns forwarding")
				if err != nil {
					return b, fmt.Errorf("unable to get configuration data, error code: %v\n", err)
				}
				b, err = c.Get(cfg, root)
				return b, err
			}
		}
	}()
	if err != nil {
		log.Fatal("Critical issue, exiting, error: ", err)
	}

	if !disabled(*blist, root) {
		urls := getURLs(*blist)

		if err := purgeFiles(urls); err != nil {
			log.Error("Error removing unused conf files", "error", err)
		}

		getBlacklists(timeout, getExcludes(*blist), urls)
	}
}

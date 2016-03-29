// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	c "github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/data"
	g "github.com/britannic/blacklist/global"
	"github.com/britannic/blacklist/utils"
)

var (
	cores   = runtime.NumCPU()
	build   = "UNKNOWN"
	githash = "UNKNOWN"
	version = "UNKNOWN"
)

func main() {
	// defer profile.Start(profile.CPUProfile, profile.MemProfile).Stop()
	runtime.GOMAXPROCS(cores)

	var (
		dnsmasq = g.DNSRestart
		poll    = time.Second * 2 // poll
		timeout = time.Minute * 30
	)

	f, err := os.OpenFile(g.Logfile, os.O_WRONLY|os.O_APPEND, 0755)
	if err == nil {
		log.SetFormatter(&log.TextFormatter{DisableColors: true})
		log.SetOutput(f)
	}

	o := getopts()
	a := os.Args[1:]
	if g.Args != nil {
		a = g.Args
	}
	flag.CommandLine.Parse(a)

	switch {
	case *o.Debug == true:
		g.Dbg = true

	case *o.Poll != 5:
		poll = time.Duration(*o.Poll) * time.Second
		log.Info("Poll duration", poll)

	case *o.Test:
		code := 0
		os.Exit(code)

	case *o.Version:

		fmt.Printf(" Version:\t\t%s\n Build date:\t\t%s\n Git short hash:\t%v\n", version, build, githash)
		os.Exit(0)

	case *o.Verb:
		log.SetFormatter(&log.TextFormatter{DisableColors: false})
		log.SetOutput(os.Stderr)

	}

	log.Info("CPU Cores: ", cores)

	blist, err := func() (b *c.Blacklist, err error) {
		switch g.WhatOS {
		case g.TestOS:
			b, err = c.Get(c.Testdata, g.Area.Root)
			if err != nil {
				return b, fmt.Errorf("unable to get configuration data, error code: %v\n", err)
			}
			return b, err

		default:

			cfg, err := c.Load("showCfg", "service dns forwarding")
			if err != nil {
				return b, fmt.Errorf("unable to get configuration data, error code: %v\n", err)
			}
			b, err = c.Get(cfg, g.Area.Root)
			return b, err
		}
	}()
	if err != nil {
		log.Fatal("Critical issue, exiting, error: ", err)
	}

	if !data.IsDisabled(*blist, g.Area.Root) {

		areas := data.GetURLs(*blist)

		if err = data.PurgeFiles(areas, g.DmsqDir); err != nil {
			log.Error("Error removing unused conf files: ", err)
		}

		ex := data.GetExcludes(*blist)
		dex := make(c.Dict)

		for _, k := range []string{g.Area.Domains, g.Area.Hosts} {
			getBlacklists(timeout, dex, ex, areas[k])
		}
	}

	log.Info("Reloading dnsmasq configuration...")
	s, err := utils.ReloadDNS(dnsmasq)
	if err != nil {
		log.Errorf("Error reloading dnsmasq configuration: %v", err)
	}
	log.Infof("dnsmasq command output: %v", s)
}

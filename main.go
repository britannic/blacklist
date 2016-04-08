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

	logger "github.com/Sirupsen/logrus"
	c "github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/data"
	g "github.com/britannic/blacklist/global"
	"github.com/britannic/blacklist/utils"
)

var (
	log     *logger.Logger
	cores   = runtime.NumCPU()
	build   = "UNKNOWN"
	githash = "UNKNOWN"
	version = "UNKNOWN"
)

func init() {
	log = g.Log
}

func main() {
	// defer profile.Start(profile.CPUProfile, profile.MemProfile).Stop()
	runtime.GOMAXPROCS(cores)

	var (
		dnsmasq = g.DNSRestart
		poll    = time.Second * 2 // poll
		timeout = time.Minute * 30
	)

	o := getopts()
	a := os.Args[1:]
	if g.Args != nil {
		a = g.Args
	}
	flag.CommandLine.Parse(a)

	switch {
	case *o.Debug:
		g.Dbg = true

	case *o.Poll != 5:
		poll = time.Duration(*o.Poll) * time.Second
		log.Infof("Poll duration %v", poll)

	case *o.Test:
		code := 0
		os.Exit(code)

	case *o.Version:

		fmt.Printf(" Version:\t\t%s\n Build date:\t\t%s\n Git short hash:\t%v\n", version, build, githash)
		os.Exit(0)

	case *o.Verb:
		g.LogOutput = "screen"
		s := &g.Set{
			Level:  logger.DebugLevel,
			Output: g.LogOutput,
		}

		log = g.LogInit(s)
	}

	log.Infof("CPU Cores: ", cores)

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
		log.Fatalf("Critical issue, exiting, error: %v", err)
	}

	if !data.IsDisabled(*blist, g.Area.Root) {

		areas := data.GetURLs(*blist)

		if err = data.PurgeFiles(areas, g.DmsqDir); err != nil {
			log.Errorf("Error removing unused conf files: %v", err)
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

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
	"github.com/britannic/blacklist/data"
	g "github.com/britannic/blacklist/global"
)

var (
	cores   = runtime.NumCPU()
	build   = "UNKNOWN"
	githash = "UNKNOWN"
	version = "UNKNOWN"
)

func main() {
	runtime.GOMAXPROCS(cores)

	var (
		poll    = time.Second * 2 // poll
		timeout = time.Minute * 30
	)

	whatOS := runtime.GOOS
	if whatOS == "darwin" {
		g.DmsqDir = "/tmp"
		g.Logfile = "/tmp/blacklist.log"
	}

	f, err := os.OpenFile(g.Logfile, os.O_WRONLY|os.O_CREATE, 0755)
	if err == nil {
		log.SetFormatter(&log.TextFormatter{DisableColors: true})
		log.SetOutput(f)
	}

	o := getopts()

	switch {
	case *o.debug == true:
		g.Dbg = true

	case *o.poll != 5:
		poll = time.Duration(*o.poll) * time.Second
		log.Info("Poll duration", poll)

	case *o.test:
		code := 0
		os.Exit(code)

	case *o.version:

		fmt.Printf(" Version:\t\t%s\n Build date:\t\t%s\n Git short hash:\t%v\n", version, build, githash)
		os.Exit(0)

	case *o.verb:
		log.SetFormatter(&log.TextFormatter{DisableColors: false})
		log.SetOutput(os.Stderr)

	}

	log.Info("CPU Cores: ", cores)

	blist, err := func() (b *c.Blacklist, err error) {
		switch whatOS {
		case "darwin":
			{
				b, err = c.Get(c.Testdata, g.Root)
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
				b, err = c.Get(cfg, g.Root)
				return b, err
			}
		}
	}()
	if err != nil {
		log.Fatal("Critical issue, exiting, error: ", err)
	}

	if !data.IsDisabled(*blist, g.Root) {
		areas := data.GetURLs(*blist)

		if err = data.PurgeFiles(areas); err != nil {
			log.Error("Error removing unused conf files", "error", err)
		}

		ex := data.GetExcludes(*blist)
		dex := make(c.Dict)
		getBlacklists(timeout, dex, ex, areas)
	}
}

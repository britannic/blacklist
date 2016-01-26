// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"log"
	"os"
	"runtime"

	c "github.com/britannic/blacklist/config"
)

var (
	build   = "UNKNOWN"
	cores   = runtime.NumCPU()
	dmsqDir = "/tmp"
	fSfx    = ".blacklist.conf"
	fStr    = "%v/%v.%v" + fSfx
	githash = "UNKNOWN"
	program = basename(os.Args[0])
	version = "UNKNOWN"
)

func main() {
	runtime.GOMAXPROCS(cores)
	var timeout int64 = 1e9 * 60 * 2 // 2 minutes

	o := getopts()
	if *o.Version {
		log.Printf("%s version: %s, build date: %s\n hash: %v\n", program, version, build, githash)
		os.Exit(0)
	}

	blist, err := func() (b *c.Blacklist, e error) {
		switch runtime.GOOS {
		case "darwin":
			{
				b, e = c.Get(c.Testdata, root)
				return
			}
		default:
			{
				cfg, err := c.Load("showCfg", "service dns forwarding")
				if err != nil {
					log.Printf("Unable to get configuration data, error code: %v\n", err)
				}
				b, e = c.Get(cfg, root)
				return
			}
		}
	}()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	if !disabled(*blist, root) {
		urls := getURLs(*blist)

		if err := purgeFiles(urls); err != nil {
			log.Println(err)
		}

		getBlacklists(timeout, getExcludes(*blist), urls)
	}
}

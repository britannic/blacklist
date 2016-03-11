// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

package main

import "flag"

// opts struct for command line options
type opts struct {
	debug   *bool
	file    *string
	poll    *int
	verb    *bool
	version *bool
}

// getopts returns legal command lines flags and values or displays help
func getopts() (options opts) {
	// options.file = flag.String("f", "/config/config.boot", "<file> # Load a configuration file")
	options.poll = flag.Int("i", 5, "Polling interval")
	options.debug = flag.Bool("d", false, "Enable debug mode")
	options.verb = flag.Bool("v", false, "Verbose display")
	options.version = flag.Bool("version", false, "# show program version number")
	flag.Parse()

	return
}

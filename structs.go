// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

package main

import c "github.com/britannic/blacklist/config"

// root constant defines the top level configuration node
const root = "blacklist"

// excludes holds fqdns that mustn't be blacklisted
type excludes map[string]int

// includes holds fqdns that should be blacklisted
type includes map[string]int

// purgeFileError contains the filename and err
type purgeFileError struct {
	file string
	err  error
}

// purgeErrors is a []*purgeFileError type
type purgeErrors []*purgeFileError

// Result holds returned data
type Result struct {
	Data  string
	Error error
	Src   *c.Src
}

// Job holds job information
type Job struct {
	Src     *c.Src
	results chan<- Result
}

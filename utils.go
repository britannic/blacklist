// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

// basename removes directory components and file extensions.
func basename(s string) string {
	// Discard last '/' and everything before.
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break
		}
	}
	// Preserve everything before last '.'.
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}
	return s
}

// getfile reads a file returns a []string array
func getfile(f string) (data []string, err error) {
	b, err := ioutil.ReadFile(f)
	if len(string(b)) > 0 {
		data = strings.Split(string(b), "\n")
	} else {
		data = []string{}
	}
	return
}

// isAdmin returns true if user has superuser privileges
func isAdmin() bool {
	u, _ := user.Current()
	switch u.Uid {
	case "0":
		return true
	default:
		return false
	}
}

// writeFile writes blacklist data to storage
func writeFile(fname string, data []byte) (err error) {
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("Unable to open file: %v for writing, error: %v", fname, err)
	}
	defer f.Close()

	b, err := f.Write(data)
	if err != nil {
		return fmt.Errorf("Unable to write to file, bytes written: %v, error: %v", b, err)
	}
	return
}

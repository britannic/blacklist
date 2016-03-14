// Package utils provides general utilities for blacklist
package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

// Basename removes directory components and file extensions.
func Basename(s string) string {
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

// Getfile reads a file returns a []string array
func Getfile(f string) (data []string, err error) {
	b, err := ioutil.ReadFile(f)
	if len(string(b)) > 0 {
		data = strings.Split(string(b), "\n")
	} else {
		data = []string{}
	}
	return
}

// IsAdmin returns true if user has superuser privileges
func IsAdmin() bool {
	u, _ := user.Current()
	switch u.Uid {
	case "0":
		return true
	default:
		return false
	}
}

// WriteFile writes blacklist data to storage
func WriteFile(fname string, data []byte) (err error) {
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

// Package utils provides general utilities for blacklist
package utils

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
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

// CmpHash compares the hashes of a to b and returns true if they're identical
func CmpHash(a, b []byte) bool {
	if md5.Sum(a) == md5.Sum(b) {
		return true
	}
	return false
}

// GetByteArray returns an array of []byte from a *bufio.Scanner
func GetByteArray(b *bufio.Scanner, s []byte) []byte {
	for b.Scan() {
		s = append(s, b.Bytes()...)
	}
	return s
}

// GetStringArray returns an array of []string from a *bufio.Scanner
func GetStringArray(b *bufio.Scanner, s []string) []string {
	for b.Scan() {
		s = append(s, b.Text())
	}
	return s
}

// GetFile reads a file and returns a *bufio.Scanner instance
func GetFile(fname string) (b *bufio.Scanner, err error) {
	f, err := os.Open(fname)
	b = bufio.NewScanner(f)
	return b, err
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

// ReloadDNS reloads the dnsmasq configuration
func ReloadDNS(d string) (string, error) {
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = strings.NewReader(d + "")
	// stdout, err := cmd.StdoutPipe()
	// stderr, err := cmd.StderrPipe()
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// cmd.Run()

	out, err := cmd.CombinedOutput()

	return string(out), err
}

// WriteFile writes blacklist data to storage
func WriteFile(fname string, data []byte) error {
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("Unable to open file: %v for writing, error: %v", fname, err)
	}

	defer func() error {
		if err = f.Close(); err != nil {
			return err
		}
		return err
	}()

	r := bytes.NewReader(data)
	w := bufio.NewWriter(f)
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		if _, err = w.Write(buf[:n]); err != nil {
			return err
		}

		if err = w.Flush(); err != nil {
			return err
		}
	}

	return err
}

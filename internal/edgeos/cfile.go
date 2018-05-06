// Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.
package edgeos

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// CFile holds an array of file names
type CFile struct {
	*Env
	Names []string
}

// readDir returns a listing of dnsmasq blacklist configuration files
func (c *CFile) readDir(pattern string) ([]string, error) {
	f, err := filepath.Glob(pattern)
	c.Debug(fmt.Sprintf("Files: %v\n: %v", pattern, f))
	return f, err
}

// Remove deletes a CFile array of file names
func (c *CFile) Remove() error {
	d, err := c.readDir(fmt.Sprintf(c.FnFmt, c.Dir, c.Wildcard.Node, c.Wildcard.Name, c.Ext))
	if err != nil {
		return err
	}
	f := diffArray(c.Names, d)
	c.Debug(fmt.Sprintf("Removing: %v", f))
	return purgeFiles(f)
}

// String implements string method
func (c *CFile) String() string {
	return strings.Join(c.Strings(), "\n")
}

// Strings returns a sorted array of strings.
func (c *CFile) Strings() []string {
	sort.Strings(c.Names)
	return c.Names
}

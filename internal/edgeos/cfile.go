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
	nType ntype
}

// readDir returns a listing of dnsmasq blacklist configuration files
func (c *CFile) readDir(pattern string) ([]string, error) {
	files, err := filepath.Glob(pattern)
	c.Debug(fmt.Sprintf("Files: %v\n: %v", pattern, files))
	return files, err
}

// Remove deletes a CFile array of file names
func (c *CFile) Remove() error {
	d, err := c.readDir(fmt.Sprintf(c.FnFmt, c.Dir, c.Wildcard.Node, c.Wildcard.Name, c.Ext))
	if err != nil {
		return err
	}
	files := diffArray(c.Names, d)
	c.Debug(fmt.Sprintf("Removing: %v", files))
	return purgeFiles(files)
}

// String implements string method
func (c *CFile) String() string {
	sort.Strings(c.Names)
	return strings.Join(c.Names, "\n")
}

// Strings returns a sorted array of strings.
func (c *CFile) Strings() []string {
	sort.Strings(c.Names)
	return c.Names
}

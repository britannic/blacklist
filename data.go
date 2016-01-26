// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	c "github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/regx"
)

// disabled returns true if blacklist is disabled
func disabled(d c.Blacklist, root string) bool {
	r := d[root].Disable
	return r
}

// getExcludes returns a map[string]int of excludes
func getExcludes(b c.Blacklist) (e excludes) {
	e = make(excludes)
	for pkey := range b {
		for _, skey := range b[pkey].Exclude {
			e[skey] = 0
		}
	}
	return
}

// getIncludes returns a map[string]int of includes
func getIncludes(n *c.Node) (i includes) {
	i = make(includes)
	for _, skey := range n.Include {
		i[skey] = 0
	}
	return
}

// getList returns a sorted []byte of blacklist entries
func getList(c *c.Src) (b []byte) {
	eq := "/"
	if c.Type == "domains" {
		eq = "/."
	}
	var lines []string
	for key := range c.List {
		line := fmt.Sprintf("address=%v%v/%v\n", eq, key, c.IP)
		lines = append(lines, line)
	}
	sort.Strings(lines)
	for _, line := range lines {
		b = append(b, line...)
	}
	return
}

// getURLs returns an array of config.Src structs with active urls
func getURLs(b c.Blacklist) (urls []*c.Src) {
	var inc includes
	for pkey := range b {
		if pkey != root {
			if len(getIncludes(b[pkey])) > 0 {
				inc = getIncludes(b[pkey])
			}
			b[pkey].Source["pre"] = &c.Src{List: inc, Name: "pre-configured", Type: pkey}
			if b[pkey].IP == "" {
				b[pkey].IP = b[root].IP
			}
			for skey := range b[pkey].Source {
				b[pkey].Source[skey].IP = b[pkey].IP
				urls = append(urls, b[pkey].Source[skey])
			}
		}
	}
	return
}

// process extracts hosts/domains from downloaded raw content
func process(s *c.Src, exc excludes, d string) *c.Src {
	rx := regx.Regex()
	s.List = make(map[string]int)

	for _, line := range strings.Split(d, "\n") {
		line, ok := stripPrefix(line, s.Prfx, rx)
		switch {
		case ok:
			{
				line = strings.TrimSpace(line)
				line = strings.TrimPrefix(line, s.Prfx)
				line = strings.ToLower(line)
				line = rx.SUFX.ReplaceAllString(line, "")
				fqdns := rx.FQDN.FindAllString(line, -1)
				for _, fqdn := range fqdns {
					fqdn = strings.TrimSpace(fqdn)
					i := strings.Count(fqdn, ".")
					switch {
					case i == 1:
						{
							if _, exists := exc[fqdn]; !exists {
								exc[fqdn] = 0
								s.List[fqdn] = 0
							} else {
								exc[fqdn]++
							}
						}
					case i > 1:
						{
							keys := strings.Split(fqdn, ".")
							for i := 0; i < len(keys)-1; i++ {
								key := strings.Join(keys[i:], ".")
								if _, exists := exc[key]; !exists {
									if len(key) > 5 {
										// fmt.Printf("fqdn: %v - keys: %v - key: %v\n", fqdn, keys, key)
										exc[key] = 0
										exc[fqdn] = 0
									}
									s.List[fqdn] = 0
								} else {
									exc[key]++
									exc[fqdn]++
									if _, exists := s.List[fqdn]; exists {
										s.List[fqdn]++
									}
								}
							}
						}
					default:
						break
					}
				}
			}
		default:
			break
		}
	}
	return s
}

func (p purgeErrors) String() (result string) {
	for _, e := range p {
		result += fmt.Sprintf("Error removing: %v: %v\n", e.file, e.err)
	}
	return
}

// listFiles returns a list of blacklist files
func listFiles(d string) (files []string) {
	dlist, err := ioutil.ReadDir(d)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range dlist {
		if strings.Contains(f.Name(), fSfx) {
			files = append(files, dmsqDir+"/"+f.Name())
		}
	}
	return
}

// diffArray returns the delta of two arrays
func diffArray(a, b []string) (diff []string) {
	dmap := make(map[string]int)
	for _, d := range b {
		dmap[d] = 0
	}

	for _, key := range a {
		if _, ok := dmap[key]; !ok {
			diff = append(diff, key)
		}
	}
	return
}

// purgeFiles removes any files that are no longer configured
func purgeFiles(c []*c.Src) error {
	var clist []string
	for _, s := range c {
		clist = append(clist, fmt.Sprintf(fStr, dmsqDir, s.Type, s.Name))
	}

	dlist := listFiles(dmsqDir)

	errors := make(purgeErrors, 0)
	for _, f := range diffArray(dlist, clist) {
		if err := os.Remove(f); err != nil {
			errors = append(errors, &purgeFileError{file: f, err: err})
		}
	}
	return fmt.Errorf("%v", errors)
}

// stripPrefix returns the modified line and true if it can strip the prefix
func stripPrefix(l string, p string, rx *regx.RGX) (string, bool) {
	switch {
	case p == "http":
		if !rx.HTTP.MatchString(l) {
			return l, false
		}
		return rx.HTTP.FindStringSubmatch(l)[1], true
	case p == "":
		break
	case strings.HasPrefix(l, p):
		return strings.TrimPrefix(l, p), true
	}
	return l, true
}

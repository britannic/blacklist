// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"

	c "github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/regx"
)

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

// disabled returns true if blacklist is disabled
func disabled(d c.Blacklist, root string) bool {
	r := d[root].Disable
	return r
}

// excludes stores fqdns that mustn't be blacklisted
type excludes map[string]int

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

// includes stores fqdns that should be blacklisted
type includes map[string]int

// getIncludes returns a map[string]int of includes
func getIncludes(n *c.Node) (i includes) {
	i = make(includes)
	for _, skey := range n.Include {
		i[skey] = 0
	}
	return
}

// getList returns a sorted []byte of blacklist entries
func getList(cf *c.Src) (b []byte) {
	eq := "/"
	if cf.Type == "domains" {
		eq = "/."
	}
	var lines []string

	sortKeys := func() (pkeys c.Keys) {
		for pkey := range cf.List {
			pkeys = append(pkeys, pkey)
		}
		sort.Sort(c.Keys(pkeys))
		return
	}

	for _, key := range sortKeys() {
		line := fmt.Sprintf("address=%v%v/%v\n", eq, key, cf.IP)
		lines = append(lines, line)
	}

	for _, line := range lines {
		b = append(b, line...)
	}
	return
}

// areaURLs is a map of c.Src
type areaURLs map[string][]*c.Src

// getURLs returns an array of config.Src structs with active urls
func getURLs(b c.Blacklist) (a areaURLs) {
	var inc includes
	a = make(areaURLs)

	for pkey := range b {
		var urls []*c.Src
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
			a[pkey] = urls
		}
	}
	return
}

// process extracts hosts/domains from downloaded raw content
func process(s *c.Src, exc excludes, d string) *c.Src {
	rx := regx.Regex()
	s.List = make(map[string]int)

NEXT:
	for _, line := range strings.Split(d, "\n") {
		switch {
		case strings.HasPrefix(line, "#"), strings.HasPrefix(line, "//"):
			continue NEXT
		case strings.HasPrefix(line, s.Prfx):
			var ok bool // We have to declare ok here, to fix var line shadow bug
			line, ok = stripPrefix(line, s.Prfx, rx)
			if ok {
				line = strings.ToLower(line)
				line = rx.SUFX.ReplaceAllString(line, "")
				line = strings.TrimSpace(line)
				fqdns := rx.FQDN.FindAllString(line, -1)
			FQDN:
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
								if _, exists := s.List[fqdn]; exists {
									s.List[fqdn]++
								}
							}
						}
					case i > 1:
						{
							keys := strings.Split(fqdn, ".")
							for i := 0; i < len(keys)-1; i++ {
								key := strings.Join(keys[i:], ".")
								if _, exists := exc[key]; !exists {
									// if len(key) > 5 && s.Type == "domains" {
									// fmt.Printf("fqdn: %v - keys: %v - key: %v\n", fqdn, keys, key)
									// 	exc[key] = 0
									// }
									exc[fqdn] = 0
									s.List[fqdn] = 0
								} else {
									// exc[key]++
									// exc[fqdn]++
									if _, exists := s.List[fqdn]; exists {
										s.List[fqdn]++
									} else {
										s.List[fqdn] = 0
									}
								}
							}
						}
					default:
						continue FQDN
					}
				}
			}
		default:
			continue NEXT
		}
	}

	if _, ok := s.List["localhost"]; ok {
		delete(s.List, "localhost")
	}

	fmt.Println(s)
	return s
}

// purgeFileError contains the filename and err
type purgeFileError struct {
	file string
	err  error
}

// purgeErrors is a []*purgeFileError type
type purgeErrors []*purgeFileError

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

// purgeFiles removes any files that are no longer configured
func purgeFiles(a areaURLs) error {
	var clist []string
	for k := range a {
		for _, s := range a[k] {
			clist = append(clist, fmt.Sprintf(fStr, dmsqDir, s.Type, s.Name))
		}
	}

	dlist := listFiles(dmsqDir)

	errors := make(purgeErrors, 0)
	for _, f := range diffArray(dlist, clist) {
		if err := os.Remove(f); err != nil {
			errors = append(errors, &purgeFileError{file: f, err: err})
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("%v", errors)
	}

	return nil
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
		return l, true
	case strings.HasPrefix(l, p):
		return strings.TrimPrefix(l, p), true
	}
	return l, false
}

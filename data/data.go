// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package data provides downloaded and configured data processing methods
package data

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"

	c "github.com/britannic/blacklist/config"
	g "github.com/britannic/blacklist/global"
	"github.com/britannic/blacklist/regx"
)

func debug(data []byte, err error) {
	if g.Dbg == false {
		return
	}
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

// DiffArray returns the delta of two arrays
func DiffArray(a, b []string) (diff []string) {
	dmap := make(c.Dict)
	for _, d := range b {
		dmap[d] = 0
	}

	for _, key := range a {
		if !dmap.KeyExists(key) {
			diff = append(diff, key)
		}
	}
	return
}

// IsDisabled returns true if blacklist is disabled
func IsDisabled(d c.Blacklist, root string) bool {
	r := d[root].Disable
	return r
}

// GetExcludes returns a map[string]int of excludes
func GetExcludes(b c.Blacklist) (e c.Dict) {
	e = make(c.Dict)
	for pkey := range b {
		for _, skey := range b[pkey].Exclude {
			e[skey] = 0
		}
	}
	return
}

// GetHTTP creates http requests to download data
func GetHTTP(URL string) (body []byte, err error) {
	const agent = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/601.4.4 (KHTML, like Gecko) Version/9.0.3 Safari/601.4.4`
	var (
		resp *http.Response
		req  *http.Request
	)

	req, err = http.NewRequest("GET", URL, nil)
	if err == nil {
		req.Header.Set("User-Agent", agent)
		// req.Header.Add("Content-Type", "application/json")
		debug(httputil.DumpRequestOut(req, true))
		resp, err = (&http.Client{}).Do(req)
	} else {
		log.Printf("Unable to form request for %s, error: %v", URL, err)
	}

	if err == nil {
		defer resp.Body.Close()
		debug(httputil.DumpResponse(resp, true))
		body, err = ioutil.ReadAll(resp.Body)
	}
	return
}

// GetIncludes returns a map[string]int of includes
func GetIncludes(n *c.Node) (r c.Dict) {
	r = make(c.Dict)
	for _, skey := range n.Include {
		r[skey] = 0
	}
	return
}

// GetList returns a sorted []byte of blacklist entries
func GetList(cf *c.Src) (b []byte) {
	eq := "/"
	if cf.Type == g.Area.Domains {
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

// AreaURLs is a map of c.Src
type AreaURLs map[string][]*c.Src

// GetURLs returns an array of config.Src structs with active urls
func GetURLs(b c.Blacklist) (a AreaURLs) {
	inc := make(c.Dict)
	a = make(AreaURLs)

	for pkey := range b {
		var urls []*c.Src
		if pkey != g.Root {
			inc = GetIncludes(b[pkey])

			b[pkey].Source["pre"] = &c.Src{List: inc, Name: "pre-configured", Type: pkey}
			if b[pkey].IP == "" {
				b[pkey].IP = b[g.Root].IP
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

// Process extracts hosts/domains from downloaded raw content
func Process(s *c.Src, dex c.Dict, ex c.Dict, d string) *c.Src {
	rx := regx.Regex()
	s.List = make(c.Dict)
	d = strings.ToLower(d)

NEXT:
	for _, line := range strings.Split(d, "\n") {
		switch {
		case strings.HasPrefix(line, "#"), strings.HasPrefix(line, "//"):
			continue NEXT
		case strings.HasPrefix(line, s.Prfx):
			var ok bool // We have to declare ok here, to fix var line shadow bug
			line, ok = StripPrefix(line, s.Prfx, rx)
			if ok {
				line = rx.SUFX.ReplaceAllString(line, "")
				line = strings.TrimSpace(line)
				fqdns := rx.FQDN.FindAllString(line, -1)
			FQDN:
				for _, fqdn := range fqdns {
					// fqdn = strings.TrimSpace(fqdn)
					i := strings.Count(fqdn, ".")
					isDEX := dex.SubKeyExists(fqdn)
					isEX := ex.KeyExists(fqdn)
					switch {
					case i == 1:
						{
							switch {
							case isDEX, isEX:
								continue NEXT
							case s.Type == g.Area.Domains:
								if !isDEX {
									dex[fqdn] = 0
									ex[fqdn] = 0
									s.List[fqdn] = 0
								}
							default:
								if !isDEX {
									if !isEX {
										ex[fqdn] = 0
										s.List[fqdn] = 0
									} else {
										ex[fqdn]++
										if s.List.KeyExists(fqdn) {
											s.List[fqdn]++
										}
									}
								}
							}
						}
					case i > 1:
						switch {
						case isDEX, isEX:
							continue NEXT
						case s.List.KeyExists(fqdn):
							s.List[fqdn]++
						default:
							ex[fqdn] = 0
							s.List[fqdn] = 0
							if s.Type == g.Area.Domains {
								dex[fqdn] = 0
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

	return s
}

// purgeFileError contains the filename and err
type purgeFileError struct {
	file string
	err  error
}

// purgeErrors is a []*purgeFileError type
type purgeErrors []*purgeFileError

// String returns a purgeErrors result string
func (p purgeErrors) String() (result string) {
	for _, e := range p {
		result += fmt.Sprintf("Error removing: %v: %v\n", e.file, e.err)
	}
	return
}

// ListFiles returns a list of blacklist files
func ListFiles(d string) (files []string) {
	dlist, err := ioutil.ReadDir(d)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range dlist {
		if strings.Contains(f.Name(), g.Fext) {
			files = append(files, g.DmsqDir+"/"+f.Name())
		}
	}
	return
}

// PurgeFiles removes any files that are no longer configured
func PurgeFiles(a AreaURLs) error {
	var clist []string
	for k := range a {
		for _, s := range a[k] {
			clist = append(clist, fmt.Sprintf(g.FStr, g.DmsqDir, s.Type, s.Name))
		}
	}

	dlist := ListFiles(g.DmsqDir)

	errors := make(purgeErrors, 0)
	for _, f := range DiffArray(dlist, clist) {
		if err := os.Remove(f); err != nil {
			errors = append(errors, &purgeFileError{file: f, err: err})
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("%v", errors)
	}

	return nil
}

// StripPrefix returns the modified line and true if it can strip the prefix
func StripPrefix(l string, p string, rx *regx.RGX) (string, bool) {
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

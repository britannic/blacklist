// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package data provides downloaded and configured data processing methods
package data

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
	"strings"

	c "github.com/britannic/blacklist/config"
	g "github.com/britannic/blacklist/global"
	"github.com/britannic/blacklist/regx"
)

var log = g.Log

func debug(data []byte, err error) {
	switch {
	case !g.Dbg:
		return
	case err == nil:
		fmt.Printf("%s\n\n", data)
	default:
		log.Errorf("%s\n\n", err)
	}
}

// DiffArray returns the delta of two arrays
func DiffArray(a, b []string) (diff []string) {
	var biggest, smallest []string

	switch {
	case len(a) > len(b), len(a) == len(b):
		biggest = a
		smallest = b

	case len(a) < len(b):
		biggest = b
		smallest = a
	}

	dmap := make(c.Dict)
	for _, k := range smallest {
		dmap[k] = 0
	}

	for _, k := range biggest {
		if !dmap.KeyExists(k) {
			diff = append(diff, k)
		}
	}

	sort.Strings(diff)
	return diff
}

// IsDisabled returns true if blacklist is disabled
func IsDisabled(c c.Blacklist, root string) bool {
	return c[root].Disable
}

// GetExcludes returns a map[string]int of excludes
func GetExcludes(b c.Blacklist) (ex c.Dict) {
	ex = make(c.Dict)
	for pkey := range b {
		for _, skey := range b[pkey].Exclude {
			ex[skey] = 0
		}
	}
	return ex
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
		log.Errorf("Unable to form request for %s, error: %v", URL, err)
	}

	if err == nil {
		defer resp.Body.Close()
		debug(httputil.DumpResponse(resp, true))
		body, err = ioutil.ReadAll(resp.Body)
	}
	return body, err
}

// GetIncludes returns a map[string]int of includes
func GetIncludes(n *c.Node) (r c.Dict) {
	r = make(c.Dict)
	for _, skey := range n.Include {
		r[skey] = 0
	}
	return r
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
		return pkeys
	}

	for _, key := range sortKeys() {
		line := fmt.Sprintf("address=%v%v/%v\n", eq, key, cf.IP)
		lines = append(lines, line)
	}

	for _, line := range lines {
		b = append(b, line...)
	}
	return b
}

// AreaURLs is a map of c.Src
type AreaURLs map[string][]*c.Src

// GetURLs returns an array of config.Src structs with active urls
func GetURLs(b c.Blacklist) (a AreaURLs) {
	a = make(AreaURLs)

	for pkey := range b {
		var urls []*c.Src
		if pkey != g.Area.Root {
			if inc := GetIncludes(b[pkey]); len(inc) > 0 {
				b[pkey].Source["pre"] = &c.Src{List: inc, Name: "pre-configured", Type: pkey}
			}

			if b[pkey].IP == "" {
				b[pkey].IP = b[g.Area.Root].IP
			}

			for skey := range b[pkey].Source {
				b[pkey].Source[skey].IP = b[pkey].IP
				urls = append(urls, b[pkey].Source[skey])
			}
			a[pkey] = urls
		}
	}
	return a
}

// ListFiles returns a list of blacklist files
func ListFiles(d string) (files []string, err error) {
	dlist, err := ioutil.ReadDir(d)
	if err != nil {
		return files, err
	}

	for _, f := range dlist {
		if strings.Contains(f.Name(), g.Fext) {
			files = append(files, d+"/"+f.Name())
		}
	}

	return files, err
}

// Process extracts hosts/domains from downloaded raw content
func Process(s *c.Src, dex c.Dict, ex c.Dict, b *bufio.Scanner) *c.Src {
	rx := regx.Regex
	s.List = make(c.Dict)

NEXT:
	for b.Scan() {
		line := strings.ToLower(b.Text())

		switch {
		case strings.HasPrefix(line, "#"), strings.HasPrefix(line, "//"):
			continue NEXT

		case strings.HasPrefix(line, s.Prfx):
			var ok bool // We have to declare ok here, to fix var line shadow bug
			line, ok = StripPrefixAndSuffix(line, s.Prfx, rx)
			if ok {
				fqdns := rx.FQDN.FindAllString(line, -1)

			FQDN:
				for _, fqdn := range fqdns {
					isDEX := dex.SubKeyExists(fqdn)
					isEX := ex.KeyExists(fqdn)
					isList := s.List.KeyExists(fqdn)

					switch {
					case isDEX:
						continue FQDN

					case isEX:
						ex[fqdn]++

					case isList:
						s.List[fqdn]++

					case s.Name == "pre-configured":
						switch {
						case s.Type == g.Area.Domains:
							dex[fqdn] = 0
							ex[fqdn] = 0
							s.List[fqdn] = 0

						default:
							ex[fqdn] = 0
							s.List[fqdn] = 0
						}

					case !isEX:
						ex[fqdn] = 0
						s.List[fqdn] = 0

					}
				}
			}
		default:
			continue NEXT
		}
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
func (p purgeErrors) String() (r string) {
	for _, e := range p {
		r += fmt.Sprintf("Error removing: %v: %v\n", e.file, e.err)
	}
	return r
}

// PurgeFiles removes any files that are no longer configured
func PurgeFiles(a AreaURLs, d string) error {
	var (
		clist, dlist []string
		err          error
	)

	if _, err = os.Stat(d); os.IsNotExist(err) {
		return fmt.Errorf("%v", err)
	}

	for k := range a {
		for _, s := range a[k] {
			clist = append(clist, fmt.Sprintf(g.FStr, d, s.Type, s.Name))
		}
	}

	if dlist, err = ListFiles(d); err != nil {
		return err
	}

	errors := make(purgeErrors, 0)
	for _, f := range DiffArray(dlist, clist) {
		if err := os.Remove(f); err != nil {
			errors = append(errors, &purgeFileError{file: f, err: err})
		}
	}

	if len(errors) > 0 {
		fmt.Println(errors)
		return fmt.Errorf("%v", errors)
	}

	return nil
}

// SortKeys returns a sorted list of c.Keys
func SortKeys(urls AreaURLs) (pkeys c.Keys) {
	for pkey := range urls {
		pkeys = append(pkeys, pkey)
	}
	sort.Sort(c.Keys(pkeys))
	return pkeys
}

// StripPrefixAndSuffix strips prefix and StripPrefixAndSuffix
func StripPrefixAndSuffix(s string, p string, rx *regx.RGX) (string, bool) {
	b := true

	switch {
	case p == "http":
		if !rx.HTTP.MatchString(s) {
			return s, false
		}
		s = rx.HTTP.FindStringSubmatch(s)[1]

	case strings.HasPrefix(s, p):
		s = strings.TrimPrefix(s, p)
	}

	s = rx.SUFX.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	s = strings.Replace(s, `"`, "", -1)
	return s, b
}

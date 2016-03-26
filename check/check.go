// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package check provides routines to sanity check blacklist is working correctly
package check

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/data"
	"github.com/britannic/blacklist/global"
	"github.com/britannic/blacklist/regx"
	"github.com/britannic/blacklist/utils"
)

// Args is a struct of check function parameters
type Args struct {
	Fname, Data, Dir, IP string
	Ex, Dex              config.Dict
}

// Cfg type of config.Blacklist
type Cfg struct {
	Blacklist *config.Blacklist
}

var rx = regx.Regex

// Blacklistings checks that only configured blacklisted includes are present in {domains,hosts}pre-configured.blacklist.conf
func (c *Cfg) Blacklistings(a *Args) error {
	var (
		b   *bufio.Scanner
		err error
		got []string
		l   = *c.Blacklist
	)

	for k := range l {
		if len(l[k].Include) > 0 {

			f := fmt.Sprintf(a.Fname, k)

			if b, err = utils.GetFile(f); err != nil {
				return err
			}

			for b.Scan() {
				got = append(got, b.Text())
			}

			got = ExtractHost(got)
			want := l[k].Include

			if len(data.DiffArray(want, got)) != 0 {
				err = fmt.Errorf("Includes not correct in %v\n\tGot: %v\n\tWant: %v", f, got, want)
				return err
			}
		}
	}
	return err
}

// Exclusions checks that configured exclusions are excluded from dnsmasq conf files
func (c *Cfg) Exclusions(a *Args) (pass bool) {
	var (
		b   *bufio.Scanner
		err error
		got []string
	)
	l := *c.Blacklist
	for k := range l {
		for sk := range l[k].Source {
			s := *l[k].Source[sk]
			f := fmt.Sprintf(global.FStr, a.Dir, s.Type, s.Name)

			if b, err = utils.GetFile(f); err != nil {
				log.Errorf("Error getting file: %v, error: %v\n", f, err)
				return pass
			}

			for b.Scan() {
				got = append(got, b.Text())
			}

			got = ExtractHost(got)
			pass = true

			for _, ex := range got {
				if _, ok := a.Ex[ex]; ok {
					pass = false
					log.Errorf("Found excluded entry %v, in %v\n", ex, f)
				}
			}
		}
	}

	return pass
}

// ExcludedDomains checks that domains are excluded from dnsmasq hosts conf files
func (c *Cfg) ExcludedDomains(a *Args) (pass bool) {
	var (
		b         *bufio.Scanner
		err       error
		got, want []string
		l         = *c.Blacklist
	)

	sortKeys := func() (pkeys config.Keys) {
		for pkey := range l {
			pkeys = append(pkeys, pkey)
			switch pkey {
			case global.Area.Domains, global.Area.Hosts:
				inc := data.GetIncludes(l[pkey])
				l[pkey].Source["pre"] = &config.Src{List: inc, Name: "pre-configured", Type: pkey}
			}
		}
		sort.Sort(config.Keys(pkeys))
		return pkeys
	}

	for _, k := range sortKeys() {
		for sk := range l[k].Source {
			s := *l[k].Source[sk]
			f := fmt.Sprintf(global.FStr, a.Dir, s.Type, s.Name)

			switch s.Type {
			case global.Area.Domains:
				if b, err = utils.GetFile(f); err != nil {
					log.Errorf("Error getting file: %v, error: %v\n", f, err)
					return pass
				}

				for b.Scan() {
					want = append(want, b.Text())
				}

				want = ExtractHost(want)
				for _, fqdn := range want {
					a.Dex[fqdn] = 0
				}

			default:
				if b, err = utils.GetFile(f); err != nil {
					log.Errorf("Error getting file: %v, error: %v\n", f, err)
					return pass
				}

				got = ExtractHost(utils.GetStringArray(b, got))
			}

			pass = true
			for _, ex := range got {
				if _, ok := a.Dex[ex]; ok {
					pass = false
					log.Errorf("Found excluded entry %v, in %v\n", ex, f)
				}
			}
		}
	}

	return pass
}

// ConfFiles checks that all blacklist sources have generated dnsmasq conf files and there aren't any orphans
func (c *Cfg) ConfFiles(a *Args) bool {
	var (
		err       error
		got, want []string
		l         = *c.Blacklist
	)

	if got, err = filepath.Glob(a.Fname); err != nil {
		log.Error(err)
		return false
	}

	for k := range l {
		for sk := range l[k].Source {
			s := *l[k].Source[sk]
			want = append(want, fmt.Sprintf(global.FStr, a.Dir, s.Type, s.Name))
		}
	}

	if diff := len(data.DiffArray(want, got)); diff != 0 {
		log.Errorf("Issues with files in %v/\n\tGot: %v\n\tWant: %v\nDiff: %v\n", global.DmsqDir, got, want, diff)
		return false
	}

	return true
}

// ConfIP checks configure IP matches redirected blackhole IP in dnsmasq conf files
func (c *Cfg) ConfIP(a *Args) bool {
	var (
		b        *bufio.Scanner
		err      error
		got, IPs []string
		l        = *c.Blacklist
		pass     = true
	)

	for k := range l {
		for sk := range l[k].Source {
			s := *l[k].Source[sk]
			f := fmt.Sprintf(global.FStr, a.Dir, s.Type, s.Name)
			if b, err = utils.GetFile(f); err != nil {
				log.Errorf("Error returned by Getfile(): %v", err)
				return false
			}

			IPs = ExtractIP(utils.GetStringArray(b, got))

			for _, ip := range IPs {
				if ip != l[global.Area.Root].IP {
					log.Errorf("Found incorrect redirection IP %v, in %v\n", ip, f)
					pass = false
				}
			}
		}
	}

	return pass
}

// ConfTemplates checks for existence/non-existence (governed by installation state) of the blacklist configuration templates
func ConfTemplates(a *Args) bool {
	var (
		err error
		got []byte
	)

	cmd := exec.Command("/bin/bash")
	find := "/usr/bin/find"
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%v %v", find, a.Dir))

	if got, err = cmd.Output(); err != nil {
		log.Error(err)
		return false
	}

	var want []byte
	want = append(want, a.Data...)

	if !utils.CmpHash(got, want) {
		fmt.Printf("Got: %v\nWant:%v\n", string(got), a.Data)
		return false
	}

	return true
}

// ExtractHost returns just the FQDN in a []string
func ExtractHost(s []string) (r []string) {
	for _, line := range s {
		if d := rx.HOST.FindStringSubmatch(line); len(d) > 0 {
			r = append(r, d[1])
		}
	}
	return r
}

// ExtractIP returns a map of unique IPs in []string of dnsmasq formatted entries
func ExtractIP(s []string) (r []string) {

	for _, line := range s {
		if k := rx.FLIP.FindStringSubmatch(line); len(k) > 0 {
			r = append(r, k[1])
		}
	}
	return r
}

// IPRedirection checks that each domain or host dnsmasq conf entry is redirected to the configured blackhole IP
func (c *Cfg) IPRedirection(a *Args) bool {
	var (
		b         *bufio.Scanner
		err       error
		l         = *c.Blacklist
		pass      = true
		rIP       = l[global.Area.Root].IP
		got, lIPs []string
	)

	for k := range l {
		for sk := range l[k].Source {
			s := *l[k].Source[sk]
			f := fmt.Sprintf(global.FStr, a.Dir, s.Type, s.Name)

			if b, err = utils.GetFile(f); err != nil {
				log.Error(err)
				return false
			}
			got = ExtractHost(utils.GetStringArray(b, got))

		HOST:
			for _, host := range got {
				if s.Type == global.Area.Domains {
					host = "www." + host
				}

				lIPs, err = net.LookupHost(host)

				switch {
				case err != nil:
					continue HOST

				default:
					for _, ip := range lIPs {
						if ip != rIP {
							log.Errorf("Host %v found in %v, resolves to %v - should be %v\n", host, f, ip, rIP)
							pass = false
						}
					}
				}
			}
		}
	}

	return pass
}

// IsDisabled checks that blacklist is actually disabled when the flag is true
// func (c *Cfg) IsDisabled(a *Args) error {
// 	var (
// 		b   bool
// 		e   string
// 		err error
// 		l   = *c.Blacklist
// 	)
//
// 	return fmt.Errorf(e)
// }

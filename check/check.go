// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package check provides routines to sanity check blacklist is working correctly
package check

import (
	"fmt"
	"net"
	"os/exec"
	"path/filepath"
	"strings"

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

var (
	e   string
	got []string
	rx  = regx.Regex()
)

// ConfBlacklistings checks that only configured blacklisted includes are present in {domains,hosts}pre-configured.blacklist.conf
func (c *Cfg) ConfBlacklistings(a *Args) (err error) {
	l := *c.Blacklist
	for k := range l {
		if len(l[k].Include) > 0 {

			f := fmt.Sprintf(a.Fname, k)

			got, err = utils.Getfile(f)
			if err != nil {
				return
			}

			got = ExtractHost(got)
			want := l[k].Include

			if len(data.DiffArray(want, got)) != 0 {
				err = fmt.Errorf("Includes not correct in %v\n\tGot: %v\n\tWant: %v", f, got, want)
				return
			}
		}
	}
	return
}

// ConfExclusions checks that configured exclusions are excluded from dnsmasq conf files
func (c *Cfg) ConfExclusions(a *Args) (err error) {
	l := *c.Blacklist
	for k := range l {
		for sk := range l[k].Source {
			s := *l[k].Source[sk]
			f := fmt.Sprintf(global.FStr, a.Dir, s.Type, s.Name)
			got, err = utils.Getfile(f)
			if err != nil {
				return err
			}
			got = ExtractHost(got)
			for _, ex := range got {
				if _, ok := a.Ex[ex]; ok {
					e += fmt.Sprintf("Found excluded entry %v, in %v\n", ex, f)
				}
			}
		}
	}
	if len(e) > 0 {
		err = fmt.Errorf(e)
	}
	return err
}

// ConfExcludedDomains checks that domains are excluded from dnsmasq hosts conf files
func (c *Cfg) ConfExcludedDomains(a *Args) (err error) {
	var (
		want []string
		l    = *c.Blacklist
	)

	for k := range l {
		for sk := range l[k].Source {
			s := *l[k].Source[sk]
			f := fmt.Sprintf(global.FStr, a.Dir, s.Type, s.Name)

			switch s.Type {
			case "domains":
				want, err = utils.Getfile(f)
				if err != nil {
					return err
				}
				want = ExtractHost(want)
				for _, fqdn := range want {
					a.Dex[fqdn] = 0
				}

			default:
				got, err = utils.Getfile(f)
				if err != nil {
					return err
				}
				got = ExtractHost(got)
			}
			for _, ex := range got {
				if _, ok := a.Dex[ex]; ok {
					e += fmt.Sprintf("Found excluded entry %v, in %v\n", ex, f)
				}
			}
		}
	}

	if len(e) > 0 {
		err = fmt.Errorf(e)
	}
	return err
}

// ConfFiles checks that all blacklist sources have generated dnsmasq conf files and there aren't any orphans
func (c *Cfg) ConfFiles(a *Args) (err error) {
	var (
		want, got []string
		l         = *c.Blacklist
	)

	got, err = filepath.Glob(a.Fname)
	if err != nil {
		return
	}

	for k := range l {
		for sk := range l[k].Source {
			s := *l[k].Source[sk]
			want = append(want, fmt.Sprintf(global.FStr, a.Dir, s.Type, s.Name))
		}
	}

	diff := data.DiffArray(want, got)
	if len(diff) != 0 {
		err = fmt.Errorf("Issues with files in %v/\n\tGot: %v\n\tWant: %v\nDiff: %v\n", global.DmsqDir, got, want, diff)
		return
	}

	return err
}

// ConfIP checks configure IP matches redirected blackhole IP in dnsmasq conf files
func (c *Cfg) ConfIP(a *Args) (err error) {
	var (
		e   string
		got []string
		IPs []string
		l   = *c.Blacklist
	)

	for k := range l {
		for sk := range l[k].Source {
			s := *l[k].Source[sk]
			f := fmt.Sprintf(global.FStr, a.Dir, s.Type, s.Name)
			got, err = utils.Getfile(f)
			if err != nil {
				return err
			}

			IPs = ExtractIP(got)

			for _, ip := range IPs {
				if ip != l[global.Area.Root].IP {
					e += fmt.Sprintf("Found incorrect redirection IP %v, in %v\n", ip, f)
				}
			}
		}
	}

	if len(e) > 0 {
		err = fmt.Errorf(e)
	}
	return err
}

// ConfTemplates checks for existence/non-existence (governed by installation state) of the blacklist configuration templates
func ConfTemplates(a *Args) (b bool, err error) {
	cmd := exec.Command("/bin/bash")
	find := "/usr/bin/find"
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%v %v", find, a.Dir))

	got, err := cmd.Output()
	if err != nil {
		return b, err
	}

	var want []byte
	want = append(want, a.Data...)

	if b = utils.CmpHash(got, want); !b {
		fmt.Printf("Got: %v\nWant:%v\n", string(got[:]), a.Data)
	}

	return b, err
}

// ExtractHost returns just the FQDN in a []string
func ExtractHost(s []string) (r []string) {
	for _, line := range s {
		d := rx.HOST.FindStringSubmatch(line)
		if len(d) > 0 {
			r = append(r, d[1])
		}
	}
	return
}

// ExtractIP returns a map of unique IPs in []string of dnsmasq formatted entries
func ExtractIP(s []string) (r []string) {

	for _, line := range s {
		k := rx.FLIP.FindStringSubmatch(line)
		if len(k) > 0 {
			r = append(r, k[1])
		}
	}
	return
}

// IPRedirection checks that each domain or host dnsmasq conf entry is redirected to the configured blackhole IP
func (c *Cfg) IPRedirection(a *Args) (err error) {
	var (
		l    = *c.Blacklist
		rIP  = l[global.Area.Root].IP
		lIPs []string
	)

	for k := range l {
		for sk := range l[k].Source {
			s := *l[k].Source[sk]
			f := fmt.Sprintf(global.FStr, a.Dir, s.Type, s.Name)

			got, err = utils.Getfile(f)
			if err != nil {
				return err
			}
			got = ExtractHost(got)

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
							e += fmt.Sprintf("Host %v found in %v, resolves to %v - should be %v\n", host, f, ip, rIP)
						}
					}
				}
			}
		}
	}

	if len(e) > 0 {
		err = fmt.Errorf(e)
	}
	return err
}

// IsDisabled checks that blacklist is actually disabled when the flag is true
func (c *Cfg) IsDisabled(a *Args) (err error) {

	return
}

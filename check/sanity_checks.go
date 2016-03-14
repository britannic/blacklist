// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package check provides routines to sanity check blacklist is working correctly
package check

import "github.com/britannic/blacklist/config"

// Args is a struct of check function parameters
type Args struct {
	Template string
	Ex, Dex  config.Dict
}

// LiveCfg type of config.Blacklist
type LiveCfg config.Blacklist

// ConfBlacklistings checks that only configured blacklisted includes are present in {domains,hosts}pre-configured.blacklist.conf
func (l LiveCfg) ConfBlacklistings(a Args) (b bool, err error) {
	for k := range l {
		if len(l[k].Include) > 0 {
			for inc := range l[k].Include {

			}
		}
	}
	return
}

// ConfExclusions checks that configured exclusions are excluded from dnsmasq conf files
func (l LiveCfg) ConfExclusions(a Args) (b bool, err error) {

	return
}

// ConfExcludedDomains checks that domains are excluded from dnsmasq hosts conf files
func (l LiveCfg) ConfExcludedDomains(a Args) (b bool, err error) {

	return
}

// ConfFiles checks that all blacklist sources have generated dnsmasq conf files and there aren't any orphahs
func (l LiveCfg) ConfFiles(a Args) (b bool, err error) {

	return
}

// ConfIP checks configure IP matches redirected blackhole IP in dnsmasq conf files
func (l LiveCfg) ConfIP(a Args) (b bool, err error) {

	return
}

// ConfTemplates checks that existence/non-existence (governed by installation state) of the blacklist configure templates
func (l LiveCfg) ConfTemplates(a Args) (b bool, err error) {

	return
}

// IPRedirection checks that each domain or host dnsmasq conf entry is redirected to the configured blackhole IP
func (l LiveCfg) IPRedirection(a Args) (b bool, err error) {

	return
}

// IsDisabled checks that blacklist is actually disabled when the flag is true
func (l LiveCfg) IsDisabled(a Args) (b bool, err error) {

	return
}

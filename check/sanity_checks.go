// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package check provides routines to sanity check blacklist is working correctly
package check

import "github.com/britannic/blacklist/config"

// LiveCfg type of config.Blacklist
type LiveCfg config.Blacklist

// ConfBlacklistings checks that only configured blacklisted includes are present in {domains,hostsreturn}pre-configured.blacklist.conf
func (l LiveCfg) ConfBlacklistings() (b bool, err error) {

	return
}

// ConfExclusions checks that configured exclusions are excluded from dnsmasq conf files
func (l LiveCfg) ConfExclusions() (b bool, err error) {

	return
}

// ConfExcludedDomains checks that domains are excluded from dnsmasq hosts conf files
func (l LiveCfg) ConfExcludedDomains() (b bool, err error) {

	return
}

// ConfFiles checks that all blacklist sources have generated dnsmasq conf files and there aren't any orphahs
func (l LiveCfg) ConfFiles() (b bool, err error) {

	return
}

// ConfIP checks configure IP matches redirected blackhole IP in dnsmasq conf files
func (l LiveCfg) ConfIP() (b bool, err error) {

	return
}

// ConfTemplates checks that existence/non-existence (governed by installation state) of the blacklist configure templates
func (l LiveCfg) ConfTemplates() (b bool, err error) {

	return
}

// IPRedirection checks that each domain or host dnsmasq conf entry is redirected to the configured blackhole IP
func (l LiveCfg) IPRedirection() (b bool, err error) {

	return
}

// IsDisabled checks that blacklist is actually disabled when the flag is true
func (l LiveCfg) IsDisabled() (b bool, err error) {

	return
}

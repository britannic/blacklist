package check_test

import (
	"log"
	"testing"

	"github.com/britannic/blacklist/check"
	"github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/data"
	"github.com/britannic/blacklist/global"
)

var (
	blacklist *config.Blacklist
	live      = &check.Cfg{Blacklist: blacklist}
	dmsqdir   string
)

func init() {
	switch global.WhatOS {
	case global.TestOS:
		dmsqdir = "../testdata"
	default:
		dmsqdir = global.DmsqDir
	}
	var err error
	live.Blacklist, err = config.Get(config.Testdata, global.Area.Root)
	if err != nil {
		log.Fatal("Couldn't load config.Testdata")
	}
}

func TestBlacklistings(t *testing.T) {
	a := &check.Args{
		Fname: dmsqdir + "/%v" + ".pre-configured" + global.Fext,
	}

	err := live.Blacklistings(a)
	if err != nil {
		t.Errorf("check.Blacklistings returned an error: %v ", err)
	}
}

func TestExclusions(t *testing.T) {
	a := &check.Args{
		Ex:  data.GetExcludes(*live.Blacklist),
		Dir: dmsqdir,
	}

	err := live.Exclusions(a)
	if err != nil {
		t.Errorf("Exclusions failure: %v", err)
	}
}

func TestExcludedDomains(t *testing.T) {
	a := &check.Args{
		Ex:  data.GetExcludes(*live.Blacklist),
		Dex: make(config.Dict),
		Dir: dmsqdir,
	}

	if err := live.ExcludedDomains(a); err != nil {
		t.Errorf("Excluded domains failure: %#v", err)
	}
}

func TestConfFiles(t *testing.T) {
	a := &check.Args{
		Dir:   dmsqdir,
		Fname: dmsqdir + `/*` + global.Fext,
	}

	if b, err := live.ConfFiles(a); !b {
		t.Errorf("Problems with dnsmasq configuration files: %v", err)
	}
}

func TestConfIP(t *testing.T) {
	a := &check.Args{
		Dir: dmsqdir,
	}

	if err := live.ConfIP(a); err != nil {
		t.Errorf("Problems with IP: %v", err)
	}
}

func TestConfTemplates(t *testing.T) {
	a := &check.Args{
		Data: fileManifest,
		Dir:  `../payload/blacklist`,
	}

	b, err := check.ConfTemplates(a)
	switch {
	case err != nil:
		t.Errorf("ConfTemplates returned an error: %v", err)
	case !b:
		t.Error("Configuration template nodes do not match")
	}
}

// func TestIsDisabled(t *testing.T) {
// 	a := &Args{
// 		Dir:   dmsqdir,
// 		Fname: dmsqdir + `/*` + global.Fext,
// 	}
//
// }

// func TestIPRedirection(t *testing.T) {
// 	a := &check.Args{
// 		Dir: dmsqdir,
// 	}
// 	if live.IPRedirection(a) != nil {
// 		t.Errorf("Problems with IP redirection: %v", err)
// 	}
// }

var fileManifest = `../payload/blacklist
../payload/blacklist/disabled
../payload/blacklist/disabled/node.def
../payload/blacklist/dns-redirect-ip
../payload/blacklist/dns-redirect-ip/node.def
../payload/blacklist/domains
../payload/blacklist/domains/dns-redirect-ip
../payload/blacklist/domains/dns-redirect-ip/node.def
../payload/blacklist/domains/exclude
../payload/blacklist/domains/exclude/node.def
../payload/blacklist/domains/include
../payload/blacklist/domains/include/node.def
../payload/blacklist/domains/node.def
../payload/blacklist/domains/source
../payload/blacklist/domains/source/node.def
../payload/blacklist/domains/source/node.tag
../payload/blacklist/domains/source/node.tag/description
../payload/blacklist/domains/source/node.tag/description/node.def
../payload/blacklist/domains/source/node.tag/prefix
../payload/blacklist/domains/source/node.tag/prefix/node.def
../payload/blacklist/domains/source/node.tag/url
../payload/blacklist/domains/source/node.tag/url/node.def
../payload/blacklist/exclude
../payload/blacklist/exclude/node.def
../payload/blacklist/hosts
../payload/blacklist/hosts/dns-redirect-ip
../payload/blacklist/hosts/dns-redirect-ip/node.def
../payload/blacklist/hosts/exclude
../payload/blacklist/hosts/exclude/node.def
../payload/blacklist/hosts/include
../payload/blacklist/hosts/include/node.def
../payload/blacklist/hosts/node.def
../payload/blacklist/hosts/source
../payload/blacklist/hosts/source/node.def
../payload/blacklist/hosts/source/node.tag
../payload/blacklist/hosts/source/node.tag/description
../payload/blacklist/hosts/source/node.tag/description/node.def
../payload/blacklist/hosts/source/node.tag/prefix
../payload/blacklist/hosts/source/node.tag/prefix/node.def
../payload/blacklist/hosts/source/node.tag/url
../payload/blacklist/hosts/source/node.tag/url/node.def
../payload/blacklist/node.def
`

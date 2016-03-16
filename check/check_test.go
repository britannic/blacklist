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
	err       error
	live      = &check.Cfg{Blacklist: blacklist}
	dmsqdir   string
)

func init() {
	switch global.WhatOS {
	case "darwin":
		dmsqdir = "../testdata"
	default:
		dmsqdir = global.DmsqDir
	}

	live.Blacklist, err = config.Get(config.Testdata, global.Root)
	if err != nil {
		log.Fatal("Couldn't load config.Testdata")
	}
}

func TestConfBlacklistings(t *testing.T) {
	a := &check.Args{
		Fname: dmsqdir + "/%v" + ".pre-configured" + global.Fext,
	}

	err := live.ConfBlacklistings(a)
	if err != nil {
		t.Errorf("check.ConfBlacklistings returned an error: %v ", err)
	}
}

func TestConfExclusions(t *testing.T) {
	a := &check.Args{
		Ex:  data.GetExcludes(*live.Blacklist),
		Dir: dmsqdir,
	}

	err := live.ConfExclusions(a)
	if err != nil {
		t.Errorf("Exclusions failure: %v", err)
	}
}

func TestConfExcludedDomains(t *testing.T) {
	a := &check.Args{
		Ex:  data.GetExcludes(*live.Blacklist),
		Dex: make(config.Dict),
		Dir: dmsqdir,
	}

	err := live.ConfExcludedDomains(a)
	if err != nil {
		t.Errorf("Excluded domains failure: %v", err)
	}
}

func TestConfFiles(t *testing.T) {
	a := &check.Args{
		Dir:   dmsqdir,
		Fname: dmsqdir + `/*` + global.Fext,
	}

	err := live.ConfFiles(a)
	if err != nil {
		t.Errorf("Problems with dnsmasq configuration files: %v", err)
	}
}

//TODO find and fix the bug check.ExtractIP
// func TestConfIP(t *testing.T) {
// 	a := &check.Args{
// 		Dir: dmsqdir,
// 	}
//
// 	if err := live.ConfIP(a); err != nil {
// 		t.Errorf("Problems with IP: %v", err)
// 	}
// }

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

//TODO find and fix the bug check.ExtractFQDN
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

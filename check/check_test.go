package check_test

import (
	"testing"

	"github.com/britannic/blacklist/check"
	"github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/data"
	"github.com/britannic/blacklist/global"
)

var (
	blacklist, err = config.Get(config.Testdata, global.Root)
	live           = &check.Cfg{Blacklist: blacklist}
)

func init() {
	if global.WhatOS == "darwin" {
		global.DmsqDir = "../testdata"
	}
}

func TestConfBlacklistings(t *testing.T) {
	if err != nil {
		t.Error("Couldn't load config.Testdata")
	}

	a := &check.Args{
		Fname: global.DmsqDir + "/%v" + ".pre-configured" + global.Fext,
	}

	err := live.ConfBlacklistings(a)
	if err != nil {
		t.Errorf("check.ConfBlacklistings returned an error: %v ", err)
	}
}

func TestConfExclusions(t *testing.T) {
	if err != nil {
		t.Error("Couldn't load config.Testdata")
	}

	a := &check.Args{
		Ex: data.GetExcludes(*live.Blacklist),
	}

	err := live.ConfExclusions(a)
	if err != nil {
		t.Errorf("Exclusions failure: %v", err)
	}
}

func TestConfExcludedDomains(t *testing.T) {
	if err != nil {
		t.Error("Couldn't load config.Testdata")
	}

	a := &check.Args{
		Ex:  data.GetExcludes(*live.Blacklist),
		Dex: make(config.Dict),
	}

	err := live.ConfExcludedDomains(a)
	if err != nil {
		t.Errorf("Excluded domains failure: %v", err)
	}
}

func TestConfFiles(t *testing.T) {
	if err != nil {
		t.Error("Couldn't load config.Testdata")
	}

	a := &check.Args{
		Fname: global.DmsqDir + `/*` + global.Fext,
	}

	err := live.ConfFiles(a)
	if err != nil {
		t.Errorf("Problems with dnsmasq configuration files: %v", err)
	}
}

func TestConfIP(t *testing.T) {
	if err != nil {
		t.Error("Couldn't load config.Testdata")
	}

	err := live.ConfIP()
	if err == nil {
		t.Errorf("Problems with IP: %v", err)
	}
}

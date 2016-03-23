package main

import (
	"testing"
	"time"

	"github.com/britannic/blacklist/check"
	"github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/data"
	"github.com/britannic/blacklist/global"
)

// func TestBuild(t *testing.T) {
// 	build := map[string]string{
// 		"build":   build,
// 		"githash": githash,
// 		"version": version,
// 	}
//
// 	for k := range build {
// 		if build[k] == "UNKNOWN" {
// 			t.Errorf("k is %v", build[k])
// 		}
// 	}
// }

func TestVarSrc(t *testing.T) {
	if len(src) == 0 {
		t.Errorf("Src should not be empty: %v", src)
	}
}

func TestGetBlacklists(t *testing.T) {
	timeout := time.Minute * 30
	b, err := config.Get(config.Testdata, global.Root)
	if err != nil {
		t.Errorf("unable to get configuration data, error code: %v\n", err)
	}

	if !data.IsDisabled(*b, global.Root) {
		areas := data.GetURLs(*b)
		ex := data.GetExcludes(*b)
		dex := make(config.Dict)
		getBlacklists(timeout, dex, ex, areas)

		var (
			blacklist = b
			live      = &check.Cfg{Blacklist: blacklist}
			a         = &check.Args{
				Fname: global.DmsqDir + "/%v" + ".pre-configured" + global.Fext,
				Dex:   make(config.Dict),
				Dir:   global.DmsqDir,
				Ex:    data.GetExcludes(*live.Blacklist),
				Data:  "",
				IP:    "",
			}
		)

		err = live.ConfBlacklistings(a)
		if err != nil {
			t.Errorf("check.ConfBlacklistings returned an error: %v ", err)
		}

		if err = live.ConfExclusions(a); err.Error() != "" {
			t.Errorf("Exclusions failure: %v", err)
		}

		if err = live.ConfExcludedDomains(a); err.Error() != "" {
			t.Errorf("Excluded domains failure: %#v", err)
		}

		a.Fname = global.DmsqDir + `/*` + global.Fext
		if err = live.ConfFiles(a); err != nil {
			t.Errorf("Problems with dnsmasq configuration files: %v", err)
		}

		if err = live.ConfIP(a); err.Error() != "" {
			t.Errorf("Problems with IP: %v", err)
		}
	}
}

// http://play.golang.org/p/KAwluDqGIl
var src = []*config.Src{
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "pre-configured",
		Type:    "domains",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "malc0de",
		Prfx:    "zone ",
		Type:    "domains",
		URL:     "http://malc0de.com/bl/ZONES",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "pre-configured",
		Type:    "hosts",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "adaway",
		Prfx:    "127.0.0.1 ",
		Type:    "hosts",
		URL:     "http://adaway.org/hosts.txt",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "malwaredomainlist",
		Prfx:    "127.0.0.1 ",
		Type:    "hosts",
		URL:     "http://www.malwaredomainlist.com/hostslist/hosts.txt",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "openphish",
		Prfx:    "http",
		Type:    "hosts",
		URL:     "https://openphish.com/feed.txt",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "someonewhocares",
		Prfx:    "0.0.0.0",
		Type:    "hosts",
		URL:     "http://someonewhocares.org/hosts/zero/",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "volkerschatz",
		Prfx:    "http",
		Type:    "hosts",
		URL:     "http://www.volkerschatz.com/net/adpaths",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "winhelp2002",
		Prfx:    "0.0.0.0 ",
		Type:    "hosts",
		URL:     "http://winhelp2002.mvps.org/hosts.txt",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "yoyo",
		Type:    "hosts",
		URL:     "http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext",
	},
}

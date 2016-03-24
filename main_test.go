package main

import (
	"os"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
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

func init() {
	f, err := os.OpenFile(global.Logfile, os.O_WRONLY|os.O_APPEND, 0755)
	if err == nil {
		log.SetFormatter(&log.TextFormatter{DisableColors: true})
		log.SetOutput(f)
	}
}

func TestVarSrc(t *testing.T) {
	if len(src) == 0 {
		t.Errorf("Src should not be empty: %v", src)
	}
}

func TestGetBlacklists(t *testing.T) {
	timeout := time.Minute * 30
	b, err := config.Get(config.Testdata, global.Area.Root)
	if err != nil {
		t.Errorf("unable to get configuration data, error code: %v\n", err)
	}

	if !data.IsDisabled(*b, global.Area.Root) {
		areas := data.GetURLs(*b)
		ex := data.GetExcludes(*b)
		dex := make(config.Dict)
		getBlacklists(timeout, dex, ex, areas)

		var (
			blacklist *config.Blacklist
			live      = &check.Cfg{Blacklist: blacklist}
			a         = &check.Args{
				Fname: global.DmsqDir + "/%v" + ".pre-configured" + global.Fext,
				Dex:   make(config.Dict),
				Dir:   global.DmsqDir,
				Data:  "",
				IP:    "",
			}
		)

		live.Blacklist, err = config.Get(config.Testdata, global.Area.Root)
		if err != nil {
			log.Fatal("Couldn't load config.Testdata")
		}

		a.Ex = data.GetExcludes(*live.Blacklist)

		err = live.Blacklistings(a)
		if err != nil {
			t.Errorf("check.Blacklistings returned an error: %v ", err)
		}

		if err = live.Exclusions(a); err.Error() != "" {
			t.Errorf("Exclusions failure: %v", err)
		}

		if err = live.ExcludedDomains(a); err.Error() != "" {
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

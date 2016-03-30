package main

import (
	"flag"
	"os"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/britannic/blacklist/check"
	"github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/data"
	"github.com/britannic/blacklist/global"
	. "github.com/britannic/blacklist/testutils"
)

var systemTest *bool

func TestBuild(t *testing.T) {
	build := map[string]string{
		"build":   build,
		"githash": githash,
		"version": version,
	}

	for k := range build {
		if build[k] != "UNKNOWN" {
			t.Errorf("%v is %v", k, build[k])
		}
	}
}

func init() {
	systemTest = flag.Bool("systemTest", false, "Set to true when running system tests")
	f, err := os.OpenFile(global.Logfile, os.O_WRONLY|os.O_APPEND, 0755)
	if err == nil {
		log.SetFormatter(&log.TextFormatter{DisableColors: true})
		log.SetOutput(f)
	}
}

// Test started when the test binary is started. Only calls main.
func TestSystem(t *testing.T) {
	if *systemTest {
		global.Args = []string{"-version"}
		main()
	}
}

func TestVarSrc(t *testing.T) {
	NotEquals(t, 0, len(src))
}

func TestGetOpts(t *testing.T) {
	// global.Args = []string{"-debug", "-i", "8", "-test", "-v", "-version", "-help"}

	o := getopts()

	want := `Debug
Poll
Test
Verb
Version
`
	Equals(t, want, o.String())

	tests := []struct {
		name string
		test interface{}
		exp  interface{}
	}{
		{
			name: "o.Debug",
			test: o.Debug,
			exp:  true,
		},
		// {
		// 	name:   "o.file",
		// 	test:   o.file,
		// 	exp: "",
		// },
		{
			name: "o.Poll",
			test: o.Poll,
			exp:  8,
		},
		{
			name: "o.Test",
			test: o.Test,
			exp:  true,
		},
		{
			name: "o.Verb",
			test: o.Verb,
			exp:  true,
		},
		{
			name: "o.Version",
			test: o.Version,
			exp:  true,
		},
	}

	for _, run := range tests {
		switch run.test.(type) {
		case bool:
			Equals(t, run.exp.(bool), run.test.(bool))

		case string:
			Equals(t, run.exp.(string), run.test.(string))

		case int:
			Equals(t, run.exp.(int), run.test.(int))
		}
	}
}

func TestGetBlacklists(t *testing.T) {
	timeout := time.Minute * 30
	b, err := config.Get(config.Testdata, global.Area.Root)
	OK(t, err)

	if !data.IsDisabled(*b, global.Area.Root) {
		areas := data.GetURLs(*b)
		ex := data.GetExcludes(*b)
		dex := make(config.Dict)

		for _, k := range []string{global.Area.Domains, global.Area.Hosts} {
			getBlacklists(timeout, dex, ex, areas[k])
		}

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
		OK(t, err)

		a.Ex = data.GetExcludes(*live.Blacklist)

		Assert(t, true, "Blacklist failure", live.Blacklistings(a))

		Assert(t, true, "Exclusions failure", live.Exclusions(a))

		Assert(t, true, "Excluded domains failure", live.ExcludedDomains(a))

		a.Fname = global.DmsqDir + `/*` + global.Fext
		Assert(t, true, "Problems with dnsmasq configuration files", live.ConfFiles(a))

		Assert(t, true, "Problems with dnsmasq configuration files content", live.ConfFilesContent(a))

		Assert(t, true, "Problems with IP redirection", live.ConfIP(a))
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

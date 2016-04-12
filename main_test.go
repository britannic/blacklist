package main

import (
	"flag"
	"fmt"
	"testing"
	"time"

	tlogger "github.com/Sirupsen/logrus/hooks/test"
	"github.com/britannic/blacklist/check"
	"github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/data"
	"github.com/britannic/blacklist/global"
	"github.com/britannic/blacklist/tdata"
	. "github.com/britannic/testutils"
)

var (
	tlog, hook = tlogger.NewNullLogger()
	systemTest *bool
)

func init() {
	global.Log = tlog
	systemTest = flag.Bool("systemTest", false, "Set to true when running system tests")
}

func TestBuild(t *testing.T) {
	pbuild := map[string]string{
		"build":   build,
		"githash": githash,
		"version": version,
	}

	for k := range pbuild {
		if pbuild[k] != "UNKNOWN" {
			t.Errorf("%v is %v", k, pbuild[k])
		}
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
	b, err := config.Get(tdata.Cfg, global.Area.Root)
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
				Log:   tlog,
			}
		)

		live.Blacklist, err = config.Get(tdata.Cfg, global.Area.Root)
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

func TestGetConfig(t *testing.T) {
	blist, err := getConfig(tdata.Cfg)
	OK(t, err)

	want := blacklist
	got := fmt.Sprint(blist)
	Equals(t, want, got)

	_, err = getConfig("")
	NotOK(t, err)
	Equals(t, fmt.Errorf("unable to get configuration data, error code: Configuration data is empty, cannot continue\n"), err)
}

// http://play.golang.org/p/KAwluDqGIl
var (
	src = []*config.Src{
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

	blacklist = "Node: blacklist\n\tDisabled: false\n\tRedirect IP: 0.0.0.0\n\tExcludes:\n\t\t122.2o7.net\n\t\t1e100.net\n\t\tadobedtm.com\n\t\takamai.net\n\t\tamazon.com\n\t\tamazonaws.com\n\t\tapple.com\n\t\task.com\n\t\tavast.com\n\t\tbitdefender.com\n\t\tcdn.visiblemeasures.com\n\t\tcloudfront.net\n\t\tcoremetrics.com\n\t\tedgesuite.net\n\t\tfreedns.afraid.org\n\t\tgithub.com\n\t\tgithubusercontent.com\n\t\tgoogle.com\n\t\tgoogleadservices.com\n\t\tgoogleapis.com\n\t\tgoogleusercontent.com\n\t\tgstatic.com\n\t\tgvt1.com\n\t\tgvt1.net\n\t\thb.disney.go.com\n\t\thp.com\n\t\thulu.com\n\t\timages-amazon.com\n\t\tmsdn.com\n\t\tpaypal.com\n\t\trackcdn.com\n\t\tschema.org\n\t\tskype.com\n\t\tsmacargo.com\n\t\tsourceforge.net\n\t\tssl-on9.com\n\t\tssl-on9.net\n\t\tstatic.chartbeat.com\n\t\tstorage.googleapis.com\n\t\twindows.net\n\t\tyimg.com\n\t\tytimg.com\n\nNode: domains\n\tDisabled: false\n\tIncludes:\n\t\tadsrvr.org\n\t\tadtechus.net\n\t\tadvertising.com\n\t\tcentade.com\n\t\tdoubleclick.net\n\t\tfree-counter.co.uk\n\t\tintellitxt.com\n\t\tkiosked.com\n\tSource: malc0de\n\t\tDisabled: false\n\t\tDescription: List of zones serving malicious executables observed by malc0de.com/database/\n\t\tPrefix: \"zone \"\n\t\tURL: http://malc0de.com/bl/ZONES\n\nNode: hosts\n\tDisabled: false\n\tInclude:\n\t\tbeap.gemini.yahoo.com\n\tSource: adaway\n\t\tDisabled: false\n\t\tDescription: Blocking mobile ad providers and some analytics providers\n\t\tPrefix: \"127.0.0.1 \"\n\t\tURL: http://adaway.org/hosts.txt\n\tSource: malwaredomainlist\n\t\tDisabled: false\n\t\tDescription: 127.0.0.1 based host and domain list\n\t\tPrefix: \"127.0.0.1 \"\n\t\tURL: http://www.malwaredomainlist.com/hostslist/hosts.txt\n\tSource: openphish\n\t\tDisabled: false\n\t\tDescription: OpenPhish automatic phishing detection\n\t\tPrefix: \"http\"\n\t\tURL: https://openphish.com/feed.txt\n\tSource: someonewhocares\n\t\tDisabled: false\n\t\tDescription: Zero based host and domain list\n\t\tPrefix: \"0.0.0.0\"\n\t\tURL: http://someonewhocares.org/hosts/zero/\n\tSource: volkerschatz\n\t\tDisabled: false\n\t\tDescription: Ad server blacklists\n\t\tPrefix: \"http\"\n\t\tURL: http://www.volkerschatz.com/net/adpaths\n\tSource: winhelp2002\n\t\tDisabled: false\n\t\tDescription: Zero based host and domain list\n\t\tPrefix: \"0.0.0.0 \"\n\t\tURL: http://winhelp2002.mvps.org/hosts.txt\n\tSource: yoyo\n\t\tDisabled: false\n\t\tDescription: Fully Qualified Domain Names only - no prefix to strip\n\t\tPrefix: \"\"\n\t\tURL: http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\n\n"
)

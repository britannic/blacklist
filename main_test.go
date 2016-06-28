package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/britannic/blacklist/internal/edgeos"
	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/britannic/testutils"
)

func TestArgs(t *testing.T) {
	tests := []struct {
		exp string
		opt []string
	}{
		{
			exp: "",
			opt: []string{},
		},
		{
			exp: "&{test Run config and data validation tests true false}\n",
			opt: []string{"-test"},
		},
		{
			exp: "&{version Show program version number true false}\n",
			opt: []string{"-version"},
		},
	}

	for _, test := range tests {
		o := getOpts()
		o.Init("blacklist", flag.ContinueOnError)
		out := new(bytes.Buffer)
		o.SetOutput(out)
		o.Parse(test.opt)
		o.setArgs(func(code int) {
			_ = code
			return
		})
		var act string
		o.Visit(func(flag *flag.Flag) {
			act += fmt.Sprintln(flag)
		})
		Equals(t, test.exp, act)
	}
}

func TestBasename(t *testing.T) {
	tests := []struct {
		s   string
		exp string
	}{
		{s: "e.txt", exp: "e"},
		{s: "/github.com/britannic/blacklist/internal/edgeos", exp: "edgeos"},
	}
	for _, tt := range tests {
		Equals(t, tt.exp, basename(tt.s))
	}
}

func TestBuild(t *testing.T) {
	want := map[string]string{
		"build":   build,
		"githash": githash,
		"version": version,
	}

	for k := range want {
		Equals(t, "UNKNOWN", want[k])
	}
}

func TestCommandLineArgs(t *testing.T) {
	out := new(bytes.Buffer)
	want := "  -arch string\n    \tSet EdgeOS CPU architecture (default \"amd64\")\n  -debug\n    \tEnable debug mode\n  -dir string\n    \tOverride dnsmasq directory (default \"/etc/dnsmasq.d\")\n  -f string\n    \t<file> # Load a configuration file\n  -i int\n    \tPolling interval (default 5)\n  -mips64 string\n    \tOverride target EdgeOS CPU architecture (default \"mips64\")\n  -os string\n    \tOverride native EdgeOS OS (default \"darwin\")\n  -test\n    \tRun config and data validation tests\n  -tmp string\n    \tOverride dnsmasq temporary directory (default \"/tmp\")\n  -v\tVerbose display\n  -version\n    \tShow program version number\n"

	var args string
	if len(os.Args) > 1 {
		args = os.Args[1]
	}
	switch {
	case runtime.GOOS == "linux" && args != "-test.coverprofile=profile.out":
		want = "  -arch=\"mips64\": Set EdgeOS CPU architecture\n  -debug=false: Enable debug mode\n  -dir=\"/etc/dnsmasq.d\": Override dnsmasq directory\n  -f=\"\": <file> # Load a configuration file\n  -i=5: Polling interval\n  -mips64=\"mips64\": Override target EdgeOS CPU architecture\n  -os=\"linux\": Override native EdgeOS OS\n  -test=false: Run config and data validation tests\n  -tmp=\"/tmp\": Override dnsmasq temporary directory\n  -v=false: Verbose display\n  -version=false: Show program version number\n"

	case runtime.GOOS == "linux" && args == "-test.coverprofile=profile.out":
		want = "  -arch=\"amd64\": Set EdgeOS CPU architecture\n  -debug=false: Enable debug mode\n  -dir=\"/etc/dnsmasq.d\": Override dnsmasq directory\n  -f=\"\": <file> # Load a configuration file\n  -i=5: Polling interval\n  -mips64=\"mips64\": Override target EdgeOS CPU architecture\n  -os=\"linux\": Override native EdgeOS OS\n  -test=false: Run config and data validation tests\n  -tmp=\"/tmp\": Override dnsmasq temporary directory\n  -v=false: Verbose display\n  -version=false: Show program version number\nflag provided but not defined: -test.coverprofile\n  -arch=\"amd64\": Set EdgeOS CPU architecture\n  -debug=false: Enable debug mode\n  -dir=\"/etc/dnsmasq.d\": Override dnsmasq directory\n  -f=\"\": <file> # Load a configuration file\n  -i=5: Polling interval\n  -mips64=\"mips64\": Override target EdgeOS CPU architecture\n  -os=\"linux\": Override native EdgeOS OS\n  -test=false: Run config and data validation tests\n  -tmp=\"/tmp\": Override dnsmasq temporary directory\n  -v=false: Verbose display\n  -version=false: Show program version number\n  -arch=\"amd64\": Set EdgeOS CPU architecture\n  -debug=false: Enable debug mode\n  -dir=\"/etc/dnsmasq.d\": Override dnsmasq directory\n  -f=\"\": <file> # Load a configuration file\n  -i=5: Polling interval\n  -mips64=\"mips64\": Override target EdgeOS CPU architecture\n  -os=\"linux\": Override native EdgeOS OS\n  -test=false: Run config and data validation tests\n  -tmp=\"/tmp\": Override dnsmasq temporary directory\n  -v=false: Verbose display\n  -version=false: Show program version number\n"

	case args == "-test.coverprofile=profile.out":
		want = "  -arch string\n    \tSet EdgeOS CPU architecture (default \"amd64\")\n  -debug\n    \tEnable debug mode\n  -dir string\n    \tOverride dnsmasq directory (default \"/etc/dnsmasq.d\")\n  -f string\n    \t<file> # Load a configuration file\n  -i int\n    \tPolling interval (default 5)\n  -mips64 string\n    \tOverride target EdgeOS CPU architecture (default \"mips64\")\n  -os string\n    \tOverride native EdgeOS OS (default \"darwin\")\n  -test\n    \tRun config and data validation tests\n  -tmp string\n    \tOverride dnsmasq temporary directory (default \"/tmp\")\n  -v\tVerbose display\n  -version\n    \tShow program version number\nflag provided but not defined: -test.coverprofile\n  -arch string\n    \tSet EdgeOS CPU architecture (default \"amd64\")\n  -debug\n    \tEnable debug mode\n  -dir string\n    \tOverride dnsmasq directory (default \"/etc/dnsmasq.d\")\n  -f string\n    \t<file> # Load a configuration file\n  -i int\n    \tPolling interval (default 5)\n  -mips64 string\n    \tOverride target EdgeOS CPU architecture (default \"mips64\")\n  -os string\n    \tOverride native EdgeOS OS (default \"darwin\")\n  -test\n    \tRun config and data validation tests\n  -tmp string\n    \tOverride dnsmasq temporary directory (default \"/tmp\")\n  -v\tVerbose display\n  -version\n    \tShow program version number\n  -arch string\n    \tSet EdgeOS CPU architecture (default \"amd64\")\n  -debug\n    \tEnable debug mode\n  -dir string\n    \tOverride dnsmasq directory (default \"/etc/dnsmasq.d\")\n  -f string\n    \t<file> # Load a configuration file\n  -i int\n    \tPolling interval (default 5)\n  -mips64 string\n    \tOverride target EdgeOS CPU architecture (default \"mips64\")\n  -os string\n    \tOverride native EdgeOS OS (default \"darwin\")\n  -test\n    \tRun config and data validation tests\n  -tmp string\n    \tOverride dnsmasq temporary directory (default \"/tmp\")\n  -v\tVerbose display\n  -version\n    \tShow program version number\n"
	}

	o := getOpts()
	o.Init("blacklist", flag.ContinueOnError)

	o.SetOutput(out)
	o.Parse([]string{"-h"})
	o.setArgs(func(code int) {
		_ = code
		return
	})
	got, err := ioutil.ReadAll(out)
	OK(t, err)
	Equals(t, want, string(got))
}

func TestGetCFG(t *testing.T) {
	exp := mainGetConfig
	o := getOpts()
	c := o.initEdgeOS()

	r := o.getCFG(c)
	Equals(t, exp, fmt.Sprint(r.Load()))
	*o.MIPS64 = "amd64"

	r = o.getCFG(c)
	Equals(t, "&{[] 0 -1}", fmt.Sprint(r.Load()))
}

func TestGetOpts(t *testing.T) {
	// global.Args = []string{"-debug", "-i", "8", "-test", "-v", "-version", "-help"}
	o := getOpts()
	want := "FlagSet\nARCH\nDebug\nDNSdir\nDNStmp\nFile\nMIPS64\nOS\nPoll\nTest\nVerb\nVersion\n"
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
		{
			name: "o.File",
			test: o.File,
			exp:  "",
		},
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

func TestSetArch(t *testing.T) {
	o := getOpts()

	tests := []struct {
		arch string
		exp  string
	}{
		{arch: "mips64", exp: "/etc/dnsmasq.d"},
		{arch: "linux", exp: "/tmp"},
		{arch: "darwin", exp: "/tmp"},
	}

	for _, test := range tests {
		Equals(t, test.exp, o.SetDir(test.arch))
	}
}

type cfgCLI struct {
	edgeos.CFGcli
}

func (c cfgCLI) Load() io.Reader {
	return strings.NewReader(tdata.Cfg)
}

func TestInitEdgeOS(t *testing.T) {
	o := getOpts()
	p := o.initEdgeOS()
	exp := []string{"edgeos.Parms{\nWildcard:  \"{*s *}\"\nAPI:       \"/bin/cli-shell-api\"\nArch:      \"amd64\"\nBash:      \"/bin/bash\"\nCores:     \"2\"\nDebug:     \"false\"\nDex:       \"**not initialized**\"\nDir:       \"/tmp\"\nDNSsvc:    \"service dnsmasq restart\"\nExc:       \"**not initialized**\"\nExt:       \"blacklist.conf\"\nFile:      \"**not initialized**\"\nFnFmt:     \"%v/%v.%v.%v\"\nInCLI:     \"inSession\"\nLevel:     \"service dns forwarding\"\nMethod:    \"GET\"\nNodes:     \"[domains hosts]\"\nPfx:       \"address=\"\nPoll:      \"5\"\nLtypes:    \"[file pre-configured-domain pre-configured-host url]\"\nTest:      \"false\"\nTimeout:   \"30s\"\nVerbosity: \"0\"\n}\n", "edgeos.Parms{\nWildcard:  \"{*s *}\"\nAPI:       \"/bin/cli-shell-api\"\nArch:      \"amd64\"\nBash:      \"/bin/bash\"\nCores:     \"2\"\nDebug:     \"false\"\nDex:       \"**not initialized**\"\nDir:       \"/tmp\"\nDNSsvc:    \"service dnsmasq restart\"\nExc:       \"**not initialized**\"\nExt:       \"blacklist.conf\"\nFile:      \"**not initialized**\"\nFnFmt:     \"%v/%v.%v.%v\"\nInCLI:     \"inSession\"\nLevel:     \"service dns forwarding\"\nMethod:    \"GET\"\nNodes:     \"[domains hosts]\"\nPfx:       \"address=\"\nPoll:      \"5\"\nLTypes:    \"[file pre-configured url]\"\nTest:      \"false\"\nTimeout:   \"30s\"\nVerbosity: \"0\"\n}\n"}

	switch runtime.GOOS {
	case "linux":
		Equals(t, exp[1], fmt.Sprint(p.Parms))
	default:
		Equals(t, exp[0], fmt.Sprint(p.Parms))
	}
}

var (
	JSONcfg = "{\n  \"nodes\": [{\n    \"blacklist\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [\n        \"122.2o7.net\",\n        \"1e100.net\",\n        \"adobedtm.com\",\n        \"akamai.net\",\n        \"amazon.com\",\n        \"amazonaws.com\",\n        \"apple.com\",\n        \"ask.com\",\n        \"avast.com\",\n        \"bitdefender.com\",\n        \"cdn.visiblemeasures.com\",\n        \"cloudfront.net\",\n        \"coremetrics.com\",\n        \"edgesuite.net\",\n        \"freedns.afraid.org\",\n        \"github.com\",\n        \"githubusercontent.com\",\n        \"google.com\",\n        \"googleadservices.com\",\n        \"googleapis.com\",\n        \"googleusercontent.com\",\n        \"gstatic.com\",\n        \"gvt1.com\",\n        \"gvt1.net\",\n        \"hb.disney.go.com\",\n        \"hp.com\",\n        \"hulu.com\",\n        \"images-amazon.com\",\n        \"msdn.com\",\n        \"paypal.com\",\n        \"rackcdn.com\",\n        \"schema.org\",\n        \"skype.com\",\n        \"smacargo.com\",\n        \"sourceforge.net\",\n        \"ssl-on9.com\",\n        \"ssl-on9.net\",\n        \"static.chartbeat.com\",\n        \"storage.googleapis.com\",\n        \"windows.net\",\n        \"yimg.com\",\n        \"ytimg.com\"\n        ]\n    },\n    \"domains\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [],\n      \"includes\": [\n        \"adsrvr.org\",\n        \"adtechus.net\",\n        \"advertising.com\",\n        \"centade.com\",\n        \"doubleclick.net\",\n        \"free-counter.co.uk\",\n        \"intellitxt.com\",\n        \"kiosked.com\"\n        ],\n      \"sources\": [{\n        \"malc0de\": {\n          \"disabled\": \"false\",\n          \"description\": \"List of zones serving malicious executables observed by malc0de.com/database/\",\n          \"prefix\": \"zone \",\n          \"file\": \"\",\n          \"url\": \"http://malc0de.com/bl/ZONES\"\n        }\n    }]\n    },\n    \"hosts\": {\n      \"disabled\": \"false\",\n      \"ip\": \"192.168.168.1\",\n      \"excludes\": [],\n      \"includes\": [\"beap.gemini.yahoo.com\"],\n      \"sources\": [{\n        \"adaway\": {\n          \"disabled\": \"false\",\n          \"description\": \"Blocking mobile ad providers and some analytics providers\",\n          \"prefix\": \"127.0.0.1 \",\n          \"file\": \"\",\n          \"url\": \"http://adaway.org/hosts.txt\"\n        },\n        \"malwaredomainlist\": {\n          \"disabled\": \"false\",\n          \"description\": \"127.0.0.1 based host and domain list\",\n          \"prefix\": \"127.0.0.1 \",\n          \"file\": \"\",\n          \"url\": \"http://www.malwaredomainlist.com/hostslist/hosts.txt\"\n        },\n        \"openphish\": {\n          \"disabled\": \"false\",\n          \"description\": \"OpenPhish automatic phishing detection\",\n          \"prefix\": \"http\",\n          \"file\": \"\",\n          \"url\": \"https://openphish.com/feed.txt\"\n        },\n        \"someonewhocares\": {\n          \"disabled\": \"false\",\n          \"description\": \"Zero based host and domain list\",\n          \"prefix\": \"0.0.0.0\",\n          \"file\": \"\",\n          \"url\": \"http://someonewhocares.org/hosts/zero/\"\n        },\n        \"tasty\": {\n          \"disabled\": \"false\",\n          \"description\": \"File source\",\n          \"prefix\": \"\",\n          \"file\": \"../testdata/blist.hosts.src\",\n          \"url\": \"\"\n        },\n        \"volkerschatz\": {\n          \"disabled\": \"false\",\n          \"description\": \"Ad server blacklists\",\n          \"prefix\": \"http\",\n          \"file\": \"\",\n          \"url\": \"http://www.volkerschatz.com/net/adpaths\"\n        },\n        \"winhelp2002\": {\n          \"disabled\": \"false\",\n          \"description\": \"Zero based host and domain list\",\n          \"prefix\": \"0.0.0.0 \",\n          \"file\": \"\",\n          \"url\": \"http://winhelp2002.mvps.org/hosts.txt\"\n        },\n        \"yoyo\": {\n          \"disabled\": \"false\",\n          \"description\": \"Fully Qualified Domain Names only - no prefix to strip\",\n          \"prefix\": \"\",\n          \"file\": \"\",\n          \"url\": \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n        }\n    }]\n    }\n  }]\n}"

	mainGetConfig = "&{blacklist {\n    disabled false\n    dns-redirect-ip 0.0.0.0\n    domains {\n        dns-redirect-ip 192.168.100.1\n        include adsrvr.org\n        include adtechus.net\n        include advertising.com\n        include centade.com\n        include doubleclick.net\n        include free-counter.co.uk\n        include intellitxt.com\n        include kiosked.com\n        include patoghee.in\n        source malc0de {\n            dns-redirect-ip 192.168.168.1\n            description \"List of zones serving malicious executables observed by malc0de.com/database/\"\n            prefix \"zone \"\n            url http://malc0de.com/bl/ZONES\n        }\n        source malwaredomains.com {\n            dns-redirect-ip 10.0.0.1\n            description \"Just domains\"\n            prefix \"\"\n            url http://mirror1.malwaredomains.com/files/justdomains\n        }\n        source simple_tracking {\n            description \"Basic tracking list by Disconnect\"\n            prefix \"\"\n            url https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt\n        }\n        source zeus {\n            description \"abuse.ch ZeuS domain blocklist\"\n            prefix \"\"\n            url https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist\n        }\n    }\n    exclude 1e100.net\n    exclude 2o7.net\n    exclude adobedtm.com\n    exclude akamai.net\n    exclude akamaihd.net\n    exclude amazon.com\n    exclude amazonaws.com\n    exclude apple.com\n    exclude ask.com\n    exclude avast.com\n    exclude bitdefender.com\n    exclude cdn.visiblemeasures.com\n    exclude cloudfront.net\n    exclude coremetrics.com\n    exclude edgesuite.net\n    exclude freedns.afraid.org\n    exclude github.com\n    exclude githubusercontent.com\n    exclude google.com\n    exclude googleadservices.com\n    exclude googleapis.com\n    exclude googletagmanager.com\n    exclude googleusercontent.com\n    exclude gstatic.com\n    exclude gvt1.com\n    exclude gvt1.net\n    exclude hb.disney.go.com\n    exclude hp.com\n    exclude hulu.com\n    exclude images-amazon.com\n    exclude live.com\n    exclude microsoft.com\n    exclude msdn.com\n    exclude msecnd.net\n    exclude paypal.com\n    exclude rackcdn.com\n    exclude schema.org\n    exclude shopify.com\n    exclude skype.com\n    exclude smacargo.com\n    exclude sourceforge.net\n    exclude ssl-on9.com\n    exclude ssl-on9.net\n    exclude sstatic.net\n    exclude static.chartbeat.com\n    exclude storage.googleapis.com\n    exclude windows.net\n    exclude xboxlive.com\n    exclude yimg.com\n    exclude ytimg.com\n    hosts {\n        include beap.gemini.yahoo.com\n        source openphish {\n            description \"OpenPhish automatic phishing detection\"\n            prefix http\n            url https://openphish.com/feed.txt\n        }\n        source raw.github.com {\n            description \"This hosts file is a merged collection of hosts from reputable sources\"\n            prefix \"0.0.0.0 \"\n            url https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts\n        }\n        source sysctl.org {\n            dns-redirect-ip 172.16.16.1\n            description \"This hosts file is a merged collection of hosts from cameleon\"\n            prefix \"127.0.0.1\t \"\n            url http://sysctl.org/cameleon/hosts\n        }\n        source tasty {\n            description \"File source\"\n            dns-redirect-ip 10.10.10.10\n            file ../testdata/blist.hosts.src\n          }\n        source volkerschatz {\n            description \"Ad server blacklists\"\n            prefix http\n            url http://www.volkerschatz.com/net/adpaths\n        }\n        source yoyo {\n            description \"Fully Qualified Domain Names only - no prefix to strip\"\n            prefix \"\"\n            url http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\n        }\n    }\n}\n\n\t/* Warning: Do not remove the following line. */\n\t/* === vyatta-config-version: \"config-management@1:conntrack@1:cron@1:dhcp-relay@1:dhcp-server@4:firewall@5:ipsec@5:nat@3:qos@1:quagga@2:system@4:ubnt-pptp@1:ubnt-util@1:vrrp@1:webgui@1:webproxy@1:zone-policy@1\" === */\n\t/* Release version: v1.8.5.4884695.160608.1057 */\n} 0 -1}"
)

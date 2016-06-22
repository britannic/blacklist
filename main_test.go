package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	. "github.com/britannic/testutils"
)

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

func TestCommandLineArgs(t *testing.T) {
	// fmt.Println(os.Args[1])
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

func TestGetCFG(t *testing.T) {
	exp := `blacklist {
    disabled false
    dns-redirect-ip 0.0.0.0
    domains {
        dns-redirect-ip 0.0.0.0
        include adsrvr.org
        include adtechus.net
        include advertising.com
        include centade.com
        include doubleclick.net
        include free-counter.co.uk
        include intellitxt.com
        include kiosked.com
        source malc0de {
            description "List of zones serving malicious executables observed by malc0de.com/database/"
            prefix "zone "
            url http://malc0de.com/bl/ZONES
        }
    }
    exclude 122.2o7.net
    exclude 1e100.net
    exclude adobedtm.com
    exclude akamai.net
    exclude amazon.com
    exclude amazonaws.com
    exclude apple.com
    exclude ask.com
    exclude avast.com
    exclude bitdefender.com
    exclude cdn.visiblemeasures.com
    exclude cloudfront.net
    exclude coremetrics.com
    exclude edgesuite.net
    exclude freedns.afraid.org
    exclude github.com
    exclude githubusercontent.com
    exclude google.com
    exclude googleadservices.com
    exclude googleapis.com
    exclude googleusercontent.com
    exclude gstatic.com
    exclude gvt1.com
    exclude gvt1.net
    exclude hb.disney.go.com
    exclude hp.com
    exclude hulu.com
    exclude images-amazon.com
    exclude msdn.com
    exclude paypal.com
    exclude rackcdn.com
    exclude schema.org
    exclude skype.com
    exclude smacargo.com
    exclude sourceforge.net
    exclude ssl-on9.com
    exclude ssl-on9.net
    exclude static.chartbeat.com
    exclude storage.googleapis.com
    exclude windows.net
    exclude yimg.com
    exclude ytimg.com
    hosts {
        dns-redirect-ip 192.168.168.1
        include beap.gemini.yahoo.com
        source adaway {
            description "Blocking mobile ad providers and some analytics providers"
            prefix "127.0.0.1 "
            url http://adaway.org/hosts.txt
        }
        source malwaredomainlist {
            description "127.0.0.1 based host and domain list"
            prefix "127.0.0.1 "
            url http://www.malwaredomainlist.com/hostslist/hosts.txt
        }
        source openphish {
            description "OpenPhish automatic phishing detection"
            prefix http
            url https://openphish.com/feed.txt
        }
        source someonewhocares {
            description "Zero based host and domain list"
            prefix 0.0.0.0
            url http://someonewhocares.org/hosts/zero/
        }
				source tasty {
						description "File source"
						dns-redirect-ip 0.0.0.0
						file /config/user-data/blist.hosts.src
				}
        source volkerschatz {
            description "Ad server blacklists"
            prefix http
            url http://www.volkerschatz.com/net/adpaths
        }
        source winhelp2002 {
            description "Zero based host and domain list"
						dns-redirect-ip 0.0.0.0
            prefix "0.0.0.0 "
            url http://winhelp2002.mvps.org/hosts.txt
        }
        source yoyo {
            description "Fully Qualified Domain Names only - no prefix to strip"
            prefix ""
            url http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext
        }
    }
	/* Warning: Do not remove the following line. */
	/* === vyatta-config-version: "config-management@1:conntrack@1:cron@1:dhcp-relay@1:dhcp-server@4:firewall@5:ipsec@5:nat@3:qos@1:quagga@2:system@4:ubnt-pptp@1:ubnt-util@1:vrrp@1:webgui@1:webproxy@1:zone-policy@1" === */
	/* Release version: v1.8.0.4853089.160219.1607 */
}`
	o := getOpts()
	r := o.getCFG()
	Equals(t, exp, fmt.Sprint(r.Load()))
}

var JSONcfg = "{\n  \"nodes\": [{\n    \"blacklist\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [\n        \"122.2o7.net\",\n        \"1e100.net\",\n        \"adobedtm.com\",\n        \"akamai.net\",\n        \"amazon.com\",\n        \"amazonaws.com\",\n        \"apple.com\",\n        \"ask.com\",\n        \"avast.com\",\n        \"bitdefender.com\",\n        \"cdn.visiblemeasures.com\",\n        \"cloudfront.net\",\n        \"coremetrics.com\",\n        \"edgesuite.net\",\n        \"freedns.afraid.org\",\n        \"github.com\",\n        \"githubusercontent.com\",\n        \"google.com\",\n        \"googleadservices.com\",\n        \"googleapis.com\",\n        \"googleusercontent.com\",\n        \"gstatic.com\",\n        \"gvt1.com\",\n        \"gvt1.net\",\n        \"hb.disney.go.com\",\n        \"hp.com\",\n        \"hulu.com\",\n        \"images-amazon.com\",\n        \"msdn.com\",\n        \"paypal.com\",\n        \"rackcdn.com\",\n        \"schema.org\",\n        \"skype.com\",\n        \"smacargo.com\",\n        \"sourceforge.net\",\n        \"ssl-on9.com\",\n        \"ssl-on9.net\",\n        \"static.chartbeat.com\",\n        \"storage.googleapis.com\",\n        \"windows.net\",\n        \"yimg.com\",\n        \"ytimg.com\"\n        ]\n    },\n    \"domains\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [],\n      \"includes\": [\n        \"adsrvr.org\",\n        \"adtechus.net\",\n        \"advertising.com\",\n        \"centade.com\",\n        \"doubleclick.net\",\n        \"free-counter.co.uk\",\n        \"intellitxt.com\",\n        \"kiosked.com\"\n        ],\n      \"sources\": [{\n        \"malc0de\": {\n          \"disabled\": \"false\",\n          \"description\": \"List of zones serving malicious executables observed by malc0de.com/database/\",\n          \"prefix\": \"zone \",\n          \"file\": \"\",\n          \"url\": \"http://malc0de.com/bl/ZONES\"\n        }\n    }]\n    },\n    \"hosts\": {\n      \"disabled\": \"false\",\n      \"ip\": \"192.168.168.1\",\n      \"excludes\": [],\n      \"includes\": [\"beap.gemini.yahoo.com\"],\n      \"sources\": [{\n        \"adaway\": {\n          \"disabled\": \"false\",\n          \"description\": \"Blocking mobile ad providers and some analytics providers\",\n          \"prefix\": \"127.0.0.1 \",\n          \"file\": \"\",\n          \"url\": \"http://adaway.org/hosts.txt\"\n        },\n        \"malwaredomainlist\": {\n          \"disabled\": \"false\",\n          \"description\": \"127.0.0.1 based host and domain list\",\n          \"prefix\": \"127.0.0.1 \",\n          \"file\": \"\",\n          \"url\": \"http://www.malwaredomainlist.com/hostslist/hosts.txt\"\n        },\n        \"openphish\": {\n          \"disabled\": \"false\",\n          \"description\": \"OpenPhish automatic phishing detection\",\n          \"prefix\": \"http\",\n          \"file\": \"\",\n          \"url\": \"https://openphish.com/feed.txt\"\n        },\n        \"someonewhocares\": {\n          \"disabled\": \"false\",\n          \"description\": \"Zero based host and domain list\",\n          \"prefix\": \"0.0.0.0\",\n          \"file\": \"\",\n          \"url\": \"http://someonewhocares.org/hosts/zero/\"\n        },\n        \"tasty\": {\n          \"disabled\": \"false\",\n          \"description\": \"File source\",\n          \"prefix\": \"\",\n          \"file\": \"/config/user-data/blist.hosts.src\",\n          \"url\": \"\"\n        },\n        \"volkerschatz\": {\n          \"disabled\": \"false\",\n          \"description\": \"Ad server blacklists\",\n          \"prefix\": \"http\",\n          \"file\": \"\",\n          \"url\": \"http://www.volkerschatz.com/net/adpaths\"\n        },\n        \"winhelp2002\": {\n          \"disabled\": \"false\",\n          \"description\": \"Zero based host and domain list\",\n          \"prefix\": \"0.0.0.0 \",\n          \"file\": \"\",\n          \"url\": \"http://winhelp2002.mvps.org/hosts.txt\"\n        },\n        \"yoyo\": {\n          \"disabled\": \"false\",\n          \"description\": \"Fully Qualified Domain Names only - no prefix to strip\",\n          \"prefix\": \"\",\n          \"file\": \"\",\n          \"url\": \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n        }\n    }]\n    }\n  }]\n}"

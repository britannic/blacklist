package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/britannic/blacklist/internal/edgeos"
	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/britannic/testutils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(t *testing.T) {
	Convey("Testing main()", t, func() {
		var (
			act          []error
			actReloadDNS string
		)

		exitCmd = func(int) { return }

		logFatalln = func(vals ...interface{}) {
			for _, v := range vals {
				if v != nil {
					act = append(act, v.(error))
				}
			}
		}

		logPrintf = func(s string, vals ...interface{}) {
			actReloadDNS = fmt.Sprintf(s, vals)
		}

		objex = []edgeos.IFace{edgeos.ExRtObj}

		main()
		So(act, ShouldBeNil)
		So(actReloadDNS, ShouldNotBeNil)
	})
}

func TestProcessObjects(t *testing.T) {
	c := setUpEnv()
	Convey("Testing processObjects", t, func() {
		Convey("Testing that the config is correctly loaded ", func() {
			So(c.String(), ShouldEqual, mainGetConfig)
			err := processObjects(c,
				[]edgeos.IFace{
					edgeos.ExRtObj,
					edgeos.ExDmObj,
					edgeos.ExHtObj,
				})
			So(err, ShouldBeNil)
		})
		Convey("Testing that c.Dex is correct after the load ", func() {
			So(c.Dex.String(), ShouldEqual, expMap)
		})
		Convey("Testing that c.Exc is correct after the load ", func() {
			So(c.Exc.String(), ShouldEqual, expMap)
		})

		Convey("Forcing processObjects to fail ", func() {
			So(processObjects(c, []edgeos.IFace{100}), ShouldNotBeNil)
		})

		Convey("Testing processObjects() with a non-existent directory ", func() {
			c.Dir = "EinenSieAugenBlick"
			So(processObjects(c, []edgeos.IFace{edgeos.FileObj}), ShouldNotBeNil)
		})
	})
}

func TestGetOpts(t *testing.T) {
	exitCmd = func(int) { return }
	origArgs := os.Args
	defer func() { os.Args = origArgs; return }()

	Convey("Testing commandline output", t, func() {
		act := new(bytes.Buffer)
		exp := vanillaArgs
		prog := path.Base(os.Args[0])
		os.Args = []string{prog, "-convey-json", "-h"}

		Convey("Testing getOpts() with vanilla arguments", func() {
			o := getOpts()
			o.Init("blacklist", flag.ContinueOnError)
			o.SetOutput(act)
			o.setArgs()

			So(act.String(), ShouldEqual, exp)

			exp = optsString
			So(o.String(), ShouldEqual, exp)
		})

		Convey("Testing getOpts() with -test", func() {
			os.Args = []string{prog, "-t"}

			o := getOpts()
			o.Init("blacklist", flag.ContinueOnError)
			o.SetOutput(act)
			o.setArgs()

			So(act.String(), ShouldEqual, "")
		})

		Convey("Testing getOpts() with -version", func() {
			os.Args = []string{prog, "-version"}

			o := getOpts()
			o.Init("blacklist", flag.ContinueOnError)
			o.SetOutput(act)
			o.setArgs()

			So(act.String(), ShouldEqual, "")
		})

		Convey("Now lets test with an invalid flag", func() {
			os.Args = []string{prog, "-z"}
			o := getOpts()
			o.Init("pixelserv", flag.ContinueOnError)
			o.SetOutput(act)
			o.setArgs()

			exp = "flag provided but not defined: -z\n" + vanillaArgs + vanillaArgs
			So(fmt.Sprint(act), ShouldEqual, exp)
		})
	})
}

func TestBasename(t *testing.T) {
	Convey("Testing basename()", t, func() {
		tests := []struct {
			s   string
			exp string
		}{
			{s: "e.txt", exp: "e"},
			{s: "/github.com/britannic/blacklist/internal/edgeos", exp: "edgeos"},
		}

		for _, tt := range tests {
			So(basename(tt.s), ShouldEqual, tt.exp)
		}
	})
}

func TestBuild(t *testing.T) {
	Convey("Testing Build() variables", t, func() {
		want := map[string]string{
			"build":   build,
			"githash": githash,
			"version": version,
		}

		for k := range want {
			So(want[k], ShouldEqual, "UNKNOWN")
		}
	})
}

func TestCommandLineArgs(t *testing.T) {
	Convey("Testing command line arguments", t, func() {
		origArgs := os.Args
		defer func() { os.Args = origArgs; return }()
		act := new(bytes.Buffer)
		exitCmd = func(int) { return }
		exp := vanillaArgs
		prog := path.Base(os.Args[0])
		os.Args = []string{prog, "-convey-json", "-h"}

		o := getOpts()
		o.Init("blacklist", flag.ContinueOnError)
		o.SetOutput(act)
		o.Parse(cleanArgs(os.Args[1:]))
		o.setArgs()

		So(act.String(), ShouldEqual, exp)
	})
}

func TestGetCFG(t *testing.T) {
	Convey("Testing getCFG()", t, func() {
		exitCmd = func(int) { return }
		o := getOpts()
		c := o.initEdgeOS()

		c.ReadCfg(o.getCFG(c))
		So(c.String(), ShouldEqual, mainGetConfig)

		*o.MIPS64 = "amd64"
		c = o.initEdgeOS()
		c.ReadCfg(o.getCFG(c))
		So(c.String(), ShouldEqual, "{\n  \"nodes\": [{\n  }]\n}")
		Equals(t, "{\n  \"nodes\": [{\n  }]\n}", c.String())
	})
}

func TestReloadDNS(t *testing.T) {
	Convey("Testing ReloadDNS()", t, func() {
		var exp string
		c := setUpEnv()
		exitCmd = func(int) { return }
		logPrintf = func(s string, vals ...interface{}) {
			exp = fmt.Sprintf(s, vals)
		}

		reloadDNS(c)
		So(exp, ShouldEqual, "ReloadDNS(): [/bin/bash: line 1: service: command not found\n]\n")
	})
}

func TestRemoveStaleFiles(t *testing.T) {
	Convey("Testing removeStaleFiles()", t, func() {
		c := setUpEnv()
		So(removeStaleFiles(c), ShouldBeNil)
		c.SetOpt(edgeos.Dir("EinenSieAugenBlick"), edgeos.Ext("[]a]"), edgeos.FileNameFmt("[]a]"), edgeos.WCard(edgeos.Wildcard{Node: "[]a]", Name: "]"}))
		So(removeStaleFiles(c), ShouldNotBeNil)
	})
}

func TestSetArch(t *testing.T) {
	Convey("Testing getCFG()", t, func() {
		exitCmd = func(int) { return }
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
			So(o.setDir(test.arch), ShouldEqual, test.exp)
		}
	})
}

type cfgCLI struct {
	edgeos.CFGcli
}

func (c cfgCLI) Load() io.Reader {
	return strings.NewReader(tdata.Cfg)
}

func TestInitEdgeOS(t *testing.T) {
	Convey("Testing initEdgeOS", t, func() {
		exitCmd = func(int) { return }
		o := getOpts()
		p := o.initEdgeOS()
		exp := "edgeos.Parms{\nWildcard:  \"{*s *}\"\nAPI:       \"/bin/cli-shell-api\"\nArch:      \"amd64\"\nBash:      \"/bin/bash\"\nCores:     \"2\"\nDebug:     \"false\"\nDex:       \"**not initialized**\"\nDir:       \"/tmp\"\nDNSsvc:    \"service dnsmasq restart\"\nExc:       \"**not initialized**\"\nExt:       \"blacklist.conf\"\nFile:      \"**not initialized**\"\nFnFmt:     \"%v/%v.%v.%v\"\nInCLI:     \"inSession\"\nLevel:     \"service dns forwarding\"\nMethod:    \"GET\"\nNodes:     \"[domains hosts]\"\nPfx:       \"address=\"\nPoll:      \"5\"\nLtypes:    \"[file pre-configured-domain pre-configured-host url]\"\nTest:      \"false\"\nTimeout:   \"30s\"\nVerbosity: \"0\"\n}\n"
		So(fmt.Sprint(p.Parms), ShouldEqual, exp)
	})
}

var (
	JSONcfg = "{\n  \"nodes\": [{\n    \"blacklist\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [\n        \"122.2o7.net\",\n        \"1e100.net\",\n        \"adobedtm.com\",\n        \"akamai.net\",\n        \"amazon.com\",\n        \"amazonaws.com\",\n        \"apple.com\",\n        \"ask.com\",\n        \"avast.com\",\n        \"bitdefender.com\",\n        \"cdn.visiblemeasures.com\",\n        \"cloudfront.net\",\n        \"coremetrics.com\",\n        \"edgesuite.net\",\n        \"freedns.afraid.org\",\n        \"github.com\",\n        \"githubusercontent.com\",\n        \"google.com\",\n        \"googleadservices.com\",\n        \"googleapis.com\",\n        \"googleusercontent.com\",\n        \"gstatic.com\",\n        \"gvt1.com\",\n        \"gvt1.net\",\n        \"hb.disney.go.com\",\n        \"hp.com\",\n        \"hulu.com\",\n        \"images-amazon.com\",\n        \"msdn.com\",\n        \"paypal.com\",\n        \"rackcdn.com\",\n        \"schema.org\",\n        \"skype.com\",\n        \"smacargo.com\",\n        \"sourceforge.net\",\n        \"ssl-on9.com\",\n        \"ssl-on9.net\",\n        \"static.chartbeat.com\",\n        \"storage.googleapis.com\",\n        \"windows.net\",\n        \"yimg.com\",\n        \"ytimg.com\"\n        ]\n    },\n    \"domains\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [],\n      \"includes\": [\n        \"adsrvr.org\",\n        \"adtechus.net\",\n        \"advertising.com\",\n        \"centade.com\",\n        \"doubleclick.net\",\n        \"free-counter.co.uk\",\n        \"intellitxt.com\",\n        \"kiosked.com\"\n        ],\n      \"sources\": [{\n        \"malc0de\": {\n          \"disabled\": \"false\",\n          \"description\": \"List of zones serving malicious executables observed by malc0de.com/database/\",\n          \"prefix\": \"zone \",\n          \"file\": \"\",\n          \"url\": \"http://malc0de.com/bl/ZONES\"\n        }\n    }]\n    },\n    \"hosts\": {\n      \"disabled\": \"false\",\n      \"ip\": \"192.168.168.1\",\n      \"excludes\": [],\n      \"includes\": [\"beap.gemini.yahoo.com\"],\n      \"sources\": [{\n        \"adaway\": {\n          \"disabled\": \"false\",\n          \"description\": \"Blocking mobile ad providers and some analytics providers\",\n          \"prefix\": \"127.0.0.1 \",\n          \"file\": \"\",\n          \"url\": \"http://adaway.org/hosts.txt\"\n        },\n        \"malwaredomainlist\": {\n          \"disabled\": \"false\",\n          \"description\": \"127.0.0.1 based host and domain list\",\n          \"prefix\": \"127.0.0.1 \",\n          \"file\": \"\",\n          \"url\": \"http://www.malwaredomainlist.com/hostslist/hosts.txt\"\n        },\n        \"openphish\": {\n          \"disabled\": \"false\",\n          \"description\": \"OpenPhish automatic phishing detection\",\n          \"prefix\": \"http\",\n          \"file\": \"\",\n          \"url\": \"https://openphish.com/feed.txt\"\n        },\n        \"someonewhocares\": {\n          \"disabled\": \"false\",\n          \"description\": \"Zero based host and domain list\",\n          \"prefix\": \"0.0.0.0\",\n          \"file\": \"\",\n          \"url\": \"http://someonewhocares.org/hosts/zero/\"\n        },\n        \"tasty\": {\n          \"disabled\": \"false\",\n          \"description\": \"File source\",\n          \"prefix\": \"\",\n          \"file\": \"../testdata/blist.hosts.src\",\n          \"url\": \"\"\n        },\n        \"volkerschatz\": {\n          \"disabled\": \"false\",\n          \"description\": \"Ad server blacklists\",\n          \"prefix\": \"http\",\n          \"file\": \"\",\n          \"url\": \"http://www.volkerschatz.com/net/adpaths\"\n        },\n        \"winhelp2002\": {\n          \"disabled\": \"false\",\n          \"description\": \"Zero based host and domain list\",\n          \"prefix\": \"0.0.0.0 \",\n          \"file\": \"\",\n          \"url\": \"http://winhelp2002.mvps.org/hosts.txt\"\n        },\n        \"yoyo\": {\n          \"disabled\": \"false\",\n          \"description\": \"Fully Qualified Domain Names only - no prefix to strip\",\n          \"prefix\": \"\",\n          \"file\": \"\",\n          \"url\": \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n        }\n    }]\n    }\n  }]\n}"

	mainGetConfig = "{\n  \"nodes\": [{\n    \"blacklist\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [\n        \"1e100.net\",\n        \"2o7.net\",\n        \"adobedtm.com\",\n        \"akamai.net\",\n        \"akamaihd.net\",\n        \"amazon.com\",\n        \"amazonaws.com\",\n        \"apple.com\",\n        \"ask.com\",\n        \"avast.com\",\n        \"bitdefender.com\",\n        \"cdn.visiblemeasures.com\",\n        \"cloudfront.net\",\n        \"coremetrics.com\",\n        \"edgesuite.net\",\n        \"freedns.afraid.org\",\n        \"github.com\",\n        \"githubusercontent.com\",\n        \"google.com\",\n        \"googleadservices.com\",\n        \"googleapis.com\",\n        \"googletagmanager.com\",\n        \"googleusercontent.com\",\n        \"gstatic.com\",\n        \"gvt1.com\",\n        \"gvt1.net\",\n        \"hb.disney.go.com\",\n        \"hp.com\",\n        \"hulu.com\",\n        \"images-amazon.com\",\n        \"live.com\",\n        \"microsoft.com\",\n        \"msdn.com\",\n        \"msecnd.net\",\n        \"paypal.com\",\n        \"rackcdn.com\",\n        \"schema.org\",\n        \"shopify.com\",\n        \"skype.com\",\n        \"smacargo.com\",\n        \"sourceforge.net\",\n        \"ssl-on9.com\",\n        \"ssl-on9.net\",\n        \"sstatic.net\",\n        \"static.chartbeat.com\",\n        \"storage.googleapis.com\",\n        \"windows.net\",\n        \"xboxlive.com\",\n        \"yimg.com\",\n        \"ytimg.com\"\n        ]\n    },\n    \"domains\": {\n      \"disabled\": \"false\",\n      \"ip\": \"192.168.100.1\",\n      \"excludes\": [],\n      \"includes\": [\n        \"adsrvr.org\",\n        \"adtechus.net\",\n        \"advertising.com\",\n        \"centade.com\",\n        \"doubleclick.net\",\n        \"free-counter.co.uk\",\n        \"intellitxt.com\",\n        \"kiosked.com\",\n        \"patoghee.in\"\n        ],\n      \"sources\": [{\n        \"malc0de\": {\n          \"disabled\": \"false\",\n          \"description\": \"List of zones serving malicious executables observed by malc0de.com/database/\",\n          \"ip\": \"192.168.168.1\",\n          \"prefix\": \"zone \",\n          \"url\": \"http://malc0de.com/bl/ZONES\",\n        },\n        \"malwaredomains.com\": {\n          \"disabled\": \"false\",\n          \"description\": \"Just domains\",\n          \"ip\": \"10.0.0.1\",\n          \"url\": \"http://mirror1.malwaredomains.com/files/justdomains\",\n        },\n        \"simple_tracking\": {\n          \"disabled\": \"false\",\n          \"description\": \"Basic tracking list by Disconnect\",\n          \"url\": \"https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt\",\n        },\n        \"zeus\": {\n          \"disabled\": \"false\",\n          \"description\": \"abuse.ch ZeuS domain blocklist\",\n          \"url\": \"https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist\",\n        }\n    }]\n    },\n    \"hosts\": {\n      \"disabled\": \"false\",\n      \"excludes\": [],\n      \"includes\": [\"beap.gemini.yahoo.com\"],\n      \"sources\": [{\n        \"openphish\": {\n          \"disabled\": \"false\",\n          \"description\": \"OpenPhish automatic phishing detection\",\n          \"prefix\": \"http\",\n          \"url\": \"https://openphish.com/feed.txt\",\n        },\n        \"raw.github.com\": {\n          \"disabled\": \"false\",\n          \"description\": \"This hosts file is a merged collection of hosts from reputable sources\",\n          \"prefix\": \"0.0.0.0 \",\n          \"url\": \"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts\",\n        },\n        \"sysctl.org\": {\n          \"disabled\": \"false\",\n          \"description\": \"This hosts file is a merged collection of hosts from cameleon\",\n          \"ip\": \"172.16.16.1\",\n          \"prefix\": \"127.0.0.1\\t \",\n          \"url\": \"http://sysctl.org/cameleon/hosts\",\n        },\n        \"tasty\": {\n          \"disabled\": \"false\",\n          \"description\": \"File source\",\n          \"ip\": \"10.10.10.10\",\n          \"file\": \"../testdata/blist.hosts.src\",\n        },\n        \"volkerschatz\": {\n          \"disabled\": \"false\",\n          \"description\": \"Ad server blacklists\",\n          \"prefix\": \"http\",\n          \"url\": \"http://www.volkerschatz.com/net/adpaths\",\n        },\n        \"yoyo\": {\n          \"disabled\": \"false\",\n          \"description\": \"Fully Qualified Domain Names only - no prefix to strip\",\n          \"url\": \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\",\n        }\n    }]\n    }\n  }]\n}"

	vanillaArgs = `  -arch string
    	Set EdgeOS CPU architecture (default "amd64")
  -debug
    	Enable debug mode
  -dir string
    	Override dnsmasq directory (default "/etc/dnsmasq.d")
  -f <file>
    	<file> # Load a configuration file
  -h	Display help
  -i int
    	Polling interval (default 5)
  -mips64 string
    	Override target EdgeOS CPU architecture (default "mips64")
  -os string
    	Override native EdgeOS OS (default "` + runtime.GOOS + `")
  -t	Run config and data validation tests
  -tmp string
    	Override dnsmasq temporary directory (default "/tmp")
  -v	Verbose display
  -version
    	Show version
`
	expMap = `"1e100.net":0,
"2o7.net":0,
"adobedtm.com":0,
"akamai.net":0,
"akamaihd.net":0,
"amazon.com":0,
"amazonaws.com":0,
"apple.com":0,
"ask.com":0,
"avast.com":0,
"bitdefender.com":0,
"cdn.visiblemeasures.com":0,
"cloudfront.net":0,
"coremetrics.com":0,
"edgesuite.net":0,
"freedns.afraid.org":0,
"github.com":0,
"githubusercontent.com":0,
"google.com":0,
"googleadservices.com":0,
"googleapis.com":0,
"googletagmanager.com":0,
"googleusercontent.com":0,
"gstatic.com":0,
"gvt1.com":0,
"gvt1.net":0,
"hb.disney.go.com":0,
"hp.com":0,
"hulu.com":0,
"images-amazon.com":0,
"live.com":0,
"microsoft.com":0,
"msdn.com":0,
"msecnd.net":0,
"paypal.com":0,
"rackcdn.com":0,
"schema.org":0,
"shopify.com":0,
"skype.com":0,
"smacargo.com":0,
"sourceforge.net":0,
"ssl-on9.com":0,
"ssl-on9.net":0,
"sstatic.net":0,
"static.chartbeat.com":0,
"storage.googleapis.com":0,
"windows.net":0,
"xboxlive.com":0,
"yimg.com":0,
"ytimg.com":0,
`

	optsString = `FlagSet
ARCH:    "amd64"
DEBUG:   "false"
DIR:     "/etc/dnsmasq.d"
F:       "**not initialized**"
H:       "true"
I:       "5"
MIPS64:  "mips64"
OS:      "` + runtime.GOOS + `"
T:       "false"
TMP:     "/tmp"
V:       "false"
VERSION: "false"
`
)

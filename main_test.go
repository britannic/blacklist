package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"

	e "github.com/britannic/blacklist/internal/edgeos"
	"github.com/britannic/mflag"
	. "github.com/britannic/testutils"
	. "github.com/smartystreets/goconvey/convey"
)

func (o *opts) String() string {
	var s string
	o.VisitAll(func(f *mflag.Flag) {
		s += fmt.Sprintf("  -%s", f.Name) // Two spaces before -; see next two comments.

		name, usage := mflag.UnquoteUsage(f)
		if len(name) > 0 {
			s += " " + name
		}
		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if len(s) <= 4 { // space, space, '-', 'x'.
			s += "\t"
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			s += "\n    \t"
		}
		s += usage
		if !mflag.IsZeroValue(f, f.DefValue) {
			if _, ok := f.Value.(*mflag.StringValue); ok {
				// put quotes on the value
				s += fmt.Sprintf(" (default %q)", f.DefValue)
			} else {
				s += fmt.Sprintf(" (default %v)", f.DefValue)
			}
		}
		s = fmt.Sprint(s, "\n")

	})

	return s
}

func TestLogFatalf(t *testing.T) {
	var (
		act string
		exp = "Something fatal happened!"
	)

	exitCmd = func(int) {}
	logCritf = func(f string, args ...interface{}) {
		act = fmt.Sprintf(f, args...)
	}

	Convey("Testing LogFatalf", t, func() {
		logFatalf("%v", exp)
		So(act, ShouldEqual, exp)
	})
}

func TestMain(t *testing.T) {
	origArgs := os.Args

	Convey("Testing main()", t, func() {
		var (
			act          string
			actReloadDNS string
			prog         = path.Base(os.Args[0])
		)

		exitCmd = func(int) {}

		logFatalf = func(f string, args ...interface{}) {
			act = fmt.Sprintf(f, args...)
		}

		logPrintf = func(f string, vals ...interface{}) {
			actReloadDNS = fmt.Sprintf(f, vals...)
		}

		screenLog()
		main()
		So(act, ShouldNotBeNil)
		So(actReloadDNS, ShouldNotBeNil)

		Convey("Testing main() with configuration file load", func() {
			act = ""
			os.Args = []string{prog, "-convey-json", "-f", "internal/testdata/config.erx.boot"}
			main()
			So(act, ShouldBeEmpty)
			os.Args = origArgs
		})

		Convey("Testing main() with non-existent configuration file load", func() {
			var s string
			os.Args = []string{prog, "-convey-json", "-f", "internal/testdata/config.bad.boot"}
			logFatalf = func(f string, args ...interface{}) {
				s = fmt.Sprintf(f, args...)
			}
			main()
			So(s, ShouldEqual, "Cannot open configuration file internal/testdata/config.bad.boot!")
			os.Args = origArgs
		})

		Convey("Testing main() with failed initEnv()", func() {
			var (
				act = new(bytes.Buffer)
				exp = ""
			)

			initEnvirons = func() (env *e.Config, err error) {
				env, _ = setUpEnv()
				err = fmt.Errorf("initEnvirons failed.")
				return env, err
			}

			os.Args = []string{prog, "-convey-json"}
			o := getOpts()
			o.Init("blacklist", mflag.ContinueOnError)
			o.SetOutput(act)
			main()
			So(act.String(), ShouldEqual, exp)
			os.Args = origArgs
		})
	})
}

func TestExitCmd(t *testing.T) {
	Convey("Testing exitCmd", t, func() {
		var (
			act int
		)
		exitCmd = func(i int) {
			act = i
		}

		exitCmd(0)
		So(act, ShouldEqual, 0)
	})
}

func TestInitEnv(t *testing.T) {
	Convey("Testing initEnv", t, func() {
		initEnv := func() (*e.Config, error) {
			return &e.Config{
				Parms: &e.Parms{Arch: "MegaOS"},
			}, nil
		}
		act, _ := initEnv()
		exp := "MegaOS"
		So(act.Arch, ShouldEqual, exp)
	})
}

func TestProcessObjects(t *testing.T) {
	c, _ := setUpEnv()
	badFileError := `open EinenSieAugenBlick/hosts.tasty.blacklist.conf: no such file or directory`
	Convey("Testing processObjects", t, func() {
		Convey("Testing that the config is correctly loaded ", func() {
			So(c.String(), ShouldEqual, mainGetConfig)
			err := processObjects(c,
				[]e.IFace{
					e.ExRtObj,
					e.ExDmObj,
					e.ExHtObj,
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
			So(processObjects(c, []e.IFace{100}), ShouldNotBeNil)
		})

		Convey("Testing processObjects() with a non-existent directory ", func() {
			c.Dir = "EinenSieAugenBlick"
			So(
				processObjects(c, []e.IFace{e.FileObj}),
				ShouldResemble,
				fmt.Errorf("%v", badFileError),
			)
		})
	})
}

func TestGetOpts(t *testing.T) {
	exitCmd = func(int) {}
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	Convey("Testing commandline output", t, func() {
		act := new(bytes.Buffer)
		exp := vanillaArgs
		prog := path.Base(os.Args[0])
		os.Args = []string{prog, "-convey-json", "-h"}

		Convey("Testing getOpts() with public arguments", func() {
			o := getOpts()
			o.Init("blacklist", mflag.ContinueOnError)
			o.SetOutput(act)
			o.setArgs()
			if IsDrone() {
				exp = vanillaArgsOnDrone
			}
			So(act.String(), ShouldEqual, exp)
		})

		Convey("Testing getOpts() with all arguments", func() {
			o := getOpts()
			o.Init("blacklist", mflag.ContinueOnError)
			o.SetOutput(act)
			o.setArgs()
			exp = allArgs
			So(o.String(), ShouldEqual, exp)
		})

		Convey("Testing getOpts() with -debug", func() {
			os.Args = []string{prog, "-debug"}

			o := getOpts()
			o.Init("blacklist", mflag.ContinueOnError)
			o.SetOutput(act)
			o.setArgs()

			So(act.String(), ShouldEqual, "")
		})

		Convey("Testing getOpts() with -t", func() {
			os.Args = []string{prog, "-t"}
			o := getOpts()
			o.Init("blacklist", mflag.ContinueOnError)
			o.SetOutput(act)
			o.setArgs()

			So(act.String(), ShouldEqual, "")
		})

		Convey("Testing getOpts() with -version", func() {
			os.Args = []string{prog, "-version"}

			o := getOpts()
			o.Init("blacklist", mflag.ContinueOnError)
			o.SetOutput(act)
			o.setArgs()

			So(act.String(), ShouldEqual, "")
		})

		Convey("Testing getOpts() with -v", func() {
			os.Args = []string{prog, "-v"}

			o := getOpts()
			o.Init("blacklist", mflag.ContinueOnError)
			o.SetOutput(act)
			o.setArgs()

			So(act.String(), ShouldEqual, "")
		})

		Convey("Now lets test with an invalid flag", func() {
			os.Args = []string{prog, "-z"}
			o := getOpts()
			o.Init("pixelserv", mflag.ContinueOnError)
			o.SetOutput(act)
			o.setArgs()

			exp = "flag provided but not defined: -z\n" + exp + exp
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
		defer func() { os.Args = origArgs }()
		act := new(bytes.Buffer)
		exitCmd = func(int) {}
		exp := vanillaArgs
		if IsDrone() {
			exp = vanillaArgsOnDrone
		}

		prog := path.Base(os.Args[0])
		os.Args = []string{prog, "-convey-json", "-h"}

		o := getOpts()
		o.Init("blacklist", mflag.ContinueOnError)
		o.SetOutput(act)
		o.Parse(cleanArgs(os.Args[1:]))
		o.setArgs()

		So(act.String(), ShouldEqual, exp)
	})
}

func TestGetCFG(t *testing.T) {
	Convey("Testing getCFG()", t, func() {
		exitCmd = func(int) {}
		o := getOpts()
		c := o.initEdgeOS()

		c.ReadCfg(o.getCFG(c))
		So(c.String(), ShouldEqual, mainGetConfig)

		*o.MIPS64 = "amd64"
		c = o.initEdgeOS()
		c.ReadCfg(o.getCFG(c))
		So(c.String(), ShouldEqual, "{\n  \"nodes\": [{\n  }]\n}")
	})
}

func TestKillFiles(t *testing.T) {
	Convey("Testing killFiles()", t, func() {
		exp := ""
		env, _ := setUpEnv()
		act := killFiles(env)
		So(fmt.Sprintf("%v", act), ShouldEqual, fmt.Sprintf("%v", exp))
	})
}

func TestReloadDNS(t *testing.T) {
	Convey("Testing ReloadDNS()", t, func() {
		var (
			act string
			exp = "ReloadDNS(): [/bin/bash: line 1: /etc/init.d/dnsmasq: No such file or directory\n]\n"
		)

		if IsDrone() {
			exp = "ReloadDNS(): [dnsmasq: unrecognized service\n]\n"
		}

		c, _ := setUpEnv()
		exitCmd = func(int) {}
		logPrintf = func(s string, v ...interface{}) {
			act = fmt.Sprintf(s, v)
		}

		reloadDNS(c)
		So(act, ShouldEqual, exp)
	})
}

func TestRemoveStaleFiles(t *testing.T) {
	Convey("Testing removeStaleFiles()", t, func() {
		c, _ := setUpEnv()
		So(removeStaleFiles(c), ShouldBeNil)
		_ = c.SetOpt(e.Dir("EinenSieAugenBlick"), e.Ext("[]a]"), e.FileNameFmt("[]a]"), e.WCard(e.Wildcard{Node: "[]a]", Name: "]"}))
		So(removeStaleFiles(c), ShouldNotBeNil)
	})
}

func TestSetArch(t *testing.T) {
	Convey("Testing getCFG()", t, func() {
		exitCmd = func(int) {}
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

func TestInitEdgeOS(t *testing.T) {
	Convey("Testing initEdgeOS", t, func() {
		exitCmd = func(int) {}
		o := getOpts()
		p := o.initEdgeOS()
		exp := `{
	"Log": {
		"Module": "blacklist",
		"ExtraCalldepth": 0
	},
	"API": "/bin/cli-shell-api",
	"Arch": "amd64",
	"Bash": "/bin/bash",
	"Cores": 2,
	"Disabled": false,
	"Dex": {},
	"Dir": "/tmp",
	"dnsmasq service": "/etc/init.d/dnsmasq restart",
	"Exc": {},
	"dnsmasq fileExt.": "blacklist.conf",
	"File name fmt": "%v/%v.%v.%v",
	"CLI Path": "service dns forwarding",
	"Leaf nodes": [
		"file",
		"domains.pre-configured",
		"hosts.pre-configured",
		"url"
	],
	"HTTP method": "GET",
	"Prefix": "address=",
	"Timeout": 30000000000,
	"Wildcard": {
		"Node": "*s",
		"Name": "*"
	}
}`
		So(fmt.Sprint(p.Parms), ShouldEqual, exp)
	})
}

var (
	// JSONcfg = "{\n  \"nodes\": [{\n    \"blacklist\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [\n        \"122.2o7.net\",\n        \"1e100.net\",\n        \"adobedtm.com\",\n        \"akamai.net\",\n        \"amazon.com\",\n        \"amazonaws.com\",\n        \"apple.com\",\n        \"ask.com\",\n        \"avast.com\",\n        \"bitdefender.com\",\n        \"cdn.visiblemeasures.com\",\n        \"cloudfront.net\",\n        \"coremetrics.com\",\n        \"edgesuite.net\",\n        \"freedns.afraid.org\",\n        \"github.com\",\n        \"githubusercontent.com\",\n        \"google.com\",\n        \"googleadservices.com\",\n        \"googleapis.com\",\n        \"googleusercontent.com\",\n        \"gstatic.com\",\n        \"gvt1.com\",\n        \"gvt1.net\",\n        \"hb.disney.go.com\",\n        \"hp.com\",\n        \"hulu.com\",\n        \"images-amazon.com\",\n        \"msdn.com\",\n        \"paypal.com\",\n        \"rackcdn.com\",\n        \"schema.org\",\n        \"skype.com\",\n        \"smacargo.com\",\n        \"sourceforge.net\",\n        \"ssl-on9.com\",\n        \"ssl-on9.net\",\n        \"static.chartbeat.com\",\n        \"storage.googleapis.com\",\n        \"windows.net\",\n        \"yimg.com\",\n        \"ytimg.com\"\n        ]\n    },\n    \"domains\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [],\n      \"includes\": [\n        \"adsrvr.org\",\n        \"adtechus.net\",\n        \"advertising.com\",\n        \"centade.com\",\n        \"doubleclick.net\",\n        \"free-counter.co.uk\",\n        \"intellitxt.com\",\n        \"kiosked.com\"\n        ],\n      \"sources\": [{\n        \"malc0de\": {\n          \"disabled\": \"false\",\n          \"description\": \"List of zones serving malicious executables observed by malc0de.com/database/\",\n          \"prefix\": \"zone \",\n          \"file\": \"\",\n          \"url\": \"http://malc0de.com/bl/ZONES\"\n        }\n    }]\n    },\n    \"hosts\": {\n      \"disabled\": \"false\",\n      \"ip\": \"192.168.168.1\",\n      \"excludes\": [],\n      \"includes\": [\"beap.gemini.yahoo.com\"],\n      \"sources\": [{\n        \"adaway\": {\n          \"disabled\": \"false\",\n          \"description\": \"Blocking mobile ad providers and some analytics providers\",\n          \"prefix\": \"127.0.0.1 \",\n          \"file\": \"\",\n          \"url\": \"http://adaway.org/hosts.txt\"\n        },\n        \"malwaredomainlist\": {\n          \"disabled\": \"false\",\n          \"description\": \"127.0.0.1 based host and domain list\",\n          \"prefix\": \"127.0.0.1 \",\n          \"file\": \"\",\n          \"url\": \"http://www.malwaredomainlist.com/hostslist/hosts.txt\"\n        },\n        \"openphish\": {\n          \"disabled\": \"false\",\n          \"description\": \"OpenPhish automatic phishing detection\",\n          \"prefix\": \"http\",\n          \"file\": \"\",\n          \"url\": \"https://openphish.com/feed.txt\"\n        },\n        \"someonewhocares\": {\n          \"disabled\": \"false\",\n          \"description\": \"Zero based host and domain list\",\n          \"prefix\": \"0.0.0.0\",\n          \"file\": \"\",\n          \"url\": \"http://someonewhocares.org/hosts/zero/\"\n        },\n        \"tasty\": {\n          \"disabled\": \"false\",\n          \"description\": \"File source\",\n          \"prefix\": \"\",\n          \"file\": \"./internal/testdata/blist.hosts.src\",\n          \"url\": \"\"\n        },\n        \"volkerschatz\": {\n          \"disabled\": \"false\",\n          \"description\": \"Ad server blacklists\",\n          \"prefix\": \"http\",\n          \"file\": \"\",\n          \"url\": \"http://www.volkerschatz.com/net/adpaths\"\n        },\n        \"winhelp2002\": {\n          \"disabled\": \"false\",\n          \"description\": \"Zero based host and domain list\",\n          \"prefix\": \"0.0.0.0 \",\n          \"file\": \"\",\n          \"url\": \"http://winhelp2002.mvps.org/hosts.txt\"\n        },\n        \"yoyo\": {\n          \"disabled\": \"false\",\n          \"description\": \"Fully Qualified Domain Names only - no prefix to strip\",\n          \"prefix\": \"\",\n          \"file\": \"\",\n          \"url\": \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n        }\n    }]\n    }\n  }]\n}"

	mainGetConfig = `{
  "nodes": [{
    "blacklist": {
      "disabled": "false",
      "ip": "192.168.168.1",
      "excludes": [
        "1e100.net",
        "2o7.net",
        "adobedtm.com",
        "akamai.net",
        "akamaihd.net",
        "amazon.com",
        "amazonaws.com",
        "apple.com",
        "ask.com",
        "avast.com",
        "avira-update.com",
        "bannerbank.com",
        "bing.com",
        "bit.ly",
        "bitdefender.com",
        "cdn.ravenjs.com",
        "cdn.visiblemeasures.com",
        "cloudfront.net",
        "coremetrics.com",
        "ebay.com",
        "edgesuite.net",
        "freedns.afraid.org",
        "github.com",
        "githubusercontent.com",
        "global.ssl.fastly.net",
        "google.com",
        "googleadservices.com",
        "googleapis.com",
        "googletagmanager.com",
        "googleusercontent.com",
        "gstatic.com",
        "gvt1.com",
        "gvt1.net",
        "hb.disney.go.com",
        "help.evernote.com",
        "herokuapp.com",
        "hp.com",
        "hulu.com",
        "images-amazon.com",
        "live.com",
        "microsoft.com",
        "msdn.com",
        "msecnd.net",
        "msftncsi.com",
        "paypal.com",
        "pop.h-cdn.co",
        "rackcdn.com",
        "rarlab.com",
        "schema.org",
        "shopify.com",
        "skype.com",
        "smacargo.com",
        "sourceforge.net",
        "spotify.com",
        "spotify.edgekey.net",
        "spotilocal.com",
        "ssl-on9.com",
        "ssl-on9.net",
        "sstatic.net",
        "static.chartbeat.com",
        "storage.googleapis.com",
        "windows.net",
        "xboxlive.com",
        "yimg.com",
        "ytimg.com"
        ]
    },
    "domains": {
      "disabled": "false",
      "excludes": ["bing.com"],
      "includes": [
        "adk2x.com",
        "adsrvr.org",
        "adtechus.net",
        "advertising.com",
        "centade.com",
        "doubleclick.net",
        "fastplayz.com",
        "free-counter.co.uk",
        "hilltopads.net",
        "intellitxt.com",
        "kiosked.com",
        "patoghee.in",
        "themillionaireinpjs.com",
        "traktrafficflow.com",
        "wwwpromoter.com"
        ],
      "sources": [{
        "malc0de": {
          "disabled": "false",
          "description": "List of zones serving malicious executables observed by malc0de.com/database/",
          "prefix": "zone ",
          "url": "http://malc0de.com/bl/ZONES",
        },
        "malwaredomains.com": {
          "disabled": "false",
          "description": "Just domains",
          "url": "http://mirror1.malwaredomains.com/files/justdomains",
        },
        "simple_tracking": {
          "disabled": "false",
          "description": "Basic tracking list by Disconnect",
          "url": "https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt",
        },
        "zeus": {
          "disabled": "false",
          "description": "abuse.ch ZeuS domain blocklist",
          "url": "https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist",
        }
    }]
    },
    "hosts": {
      "disabled": "false",
      "excludes": [],
      "includes": ["beap.gemini.yahoo.com"],
      "sources": [{
        "openphish": {
          "disabled": "false",
          "description": "OpenPhish automatic phishing detection",
          "prefix": "http",
          "url": "https://openphish.com/feed.txt",
        },
        "raw.github.com": {
          "disabled": "false",
          "description": "This hosts file is a merged collection of hosts from reputable sources",
          "prefix": "0.0.0.0 ",
          "url": "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
        },
        "sysctl.org": {
          "disabled": "false",
          "description": "This hosts file is a merged collection of hosts from cameleon",
          "prefix": "127.0.0.1\t ",
          "url": "http://sysctl.org/cameleon/hosts",
        },
        "yoyo": {
          "disabled": "false",
          "description": "Fully Qualified Domain Names only - no prefix to strip",
          "url": "http://pgl.yoyo.org/as/serverlist.phphostformat=nohtml&showintro=1&mimetype=plaintext",
        },
        "tasty": {
          "disabled": "false",
          "description": "File source",
          "ip": "10.10.10.10",
          "file": "./internal/testdata/blist.hosts.src",
        }
    }]
    }
  }]
}`
	vanillaArgs = `  -dir string
    	Override dnsmasq directory (default "/etc/dnsmasq.d")
  -f <file>
    	<file> # Load a config.boot file
  -h	Display help
  -v	Verbose display
  -version
    	Show version
`
	allArgs = `  -arch string
    	Set EdgeOS CPU architecture (default "amd64")
  -debug
    	Enable debug mode
  -dir string
    	Override dnsmasq directory (default "/etc/dnsmasq.d")
  -f <file>
    	<file> # Load a config.boot file
  -h
    	Display help
  -mips64 string
    	Override target EdgeOS CPU architecture (default "mips64")
  -mipsle string
    	Override target EdgeOS CPU architecture (default "mipsle")
  -os string
    	Override native EdgeOS OS (default "darwin")
  -t
    	Run config and data validation tests
  -tmp string
    	Override dnsmasq temporary directory (default "/tmp")
  -v
    	Verbose display
  -version
    	Show version
`

	vanillaArgsOnDrone = "  -dir=\"/etc/dnsmasq.d\": Override dnsmasq directory\n  -h=false: Display help\n  -v=false: Verbose display\n  -version=false: Show version\n"

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
"avira-update.com":0,
"bannerbank.com":0,
"bing.com":0,
"bit.ly":0,
"bitdefender.com":0,
"cdn.ravenjs.com":0,
"cdn.visiblemeasures.com":0,
"cloudfront.net":0,
"coremetrics.com":0,
"ebay.com":0,
"edgesuite.net":0,
"freedns.afraid.org":0,
"github.com":0,
"githubusercontent.com":0,
"global.ssl.fastly.net":0,
"google.com":0,
"googleadservices.com":0,
"googleapis.com":0,
"googletagmanager.com":0,
"googleusercontent.com":0,
"gstatic.com":0,
"gvt1.com":0,
"gvt1.net":0,
"hb.disney.go.com":0,
"help.evernote.com":0,
"herokuapp.com":0,
"hp.com":0,
"hulu.com":0,
"images-amazon.com":0,
"live.com":0,
"microsoft.com":0,
"msdn.com":0,
"msecnd.net":0,
"msftncsi.com":0,
"paypal.com":0,
"pop.h-cdn.co":0,
"rackcdn.com":0,
"rarlab.com":0,
"schema.org":0,
"shopify.com":0,
"skype.com":0,
"smacargo.com":0,
"sourceforge.net":0,
"spotify.com":0,
"spotify.edgekey.net":0,
"spotilocal.com":0,
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
)

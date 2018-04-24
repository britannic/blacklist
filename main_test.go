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

func init() {
	/*
	   The default failure mode is FailureHalts, which causes test execution
	   within a `Convey` block to halt at the first failure. You could use
	   that mode if the test were re-worked to aggregate all results into
	   a collection that was verified after all goroutines have finished.
	   But, as the code stands, you need to use the FailureContinues mode.

	   The following line sets the failure mode for all tests in the package:
	*/

	SetDefaultFailureMode(FailureContinues)
}

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
			prfx         = fmt.Sprintf("%s: ", prog)
		)

		exitCmd = func(int) {}

		logFatalf = func(f string, args ...interface{}) {
			act = fmt.Sprintf(f, args...)
		}

		logPrintf = func(f string, vals ...interface{}) {
			actReloadDNS = fmt.Sprintf(f, vals...)
		}

		screenLog(prfx)
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
			So(s, ShouldEqual, "cannot open configuration file internal/testdata/config.bad.boot!")
			os.Args = origArgs
		})

		Convey("Testing main() with failed initEnv()", func() {
			var (
				act = new(bytes.Buffer)
				exp = ""
			)

			initEnvirons = func() (env *e.Config, err error) {
				env, _ = setUpEnv()
				err = fmt.Errorf("initEnvirons failed")
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

func TestScreenLog(t *testing.T) {
	Convey("Testing ScreenLog(prefix)", t, func() {
		haveTerm = func() bool {
			return true
		}

		So(screenLog("prefix"), ShouldNotBeNil)
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
				Env: &e.Env{Arch: "MegaOS"},
			}, nil
		}
		act, _ := initEnv()
		exp := "MegaOS"
		So(act.Arch, ShouldEqual, exp)
	})
}

func TestProcessObjects(t *testing.T) {
	c, _ := setUpEnv()
	badFileError := `open EinenSieAugenBlick/domains.tasty.blacklist.conf: no such file or directory`
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

func TestFiles(t *testing.T) {
	Convey("Testing files()", t, func() {
		exp := ""
		env, _ := setUpEnv()
		act := files(env)
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

func TestNewScreenLogBackend(t *testing.T) {
	tests := []struct {
		exp    bool
		colors []string
		prefix string
	}{
		{exp: true, colors: boldcolors, prefix: "test"},
		{exp: false, colors: []string{}, prefix: "test"},
	}

	Convey("Testing newScreenLogBackend()", t, func() {
		for _, test := range tests {
			act := newScreenLogBackend(test.colors, test.prefix)
			So(act.Color, ShouldEqual, test.exp)
		}
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
	"HTTP method": "GET",
	"Prefix": {},
	"Timeout": 30000000000,
	"Wildcard": {
		"Node": "*s",
		"Name": "*"
	}
}`
		So(fmt.Sprint(p.Env), ShouldEqual, exp)
	})
}

var (
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
        "dropbox.com",
        "ebay.com",
        "edgesuite.net",
        "evernote.com",
        "express.co.uk",
        "feedly.com",
        "freedns.afraid.org",
        "github.com",
        "githubusercontent.com",
        "global.ssl.fastly.net",
        "google.com",
        "googleads.g.doubleclick.net",
        "googleadservices.com",
        "googleapis.com",
        "googletagmanager.com",
        "googleusercontent.com",
        "gstatic.com",
        "gvt1.com",
        "gvt1.net",
        "hb.disney.go.com",
        "herokuapp.com",
        "hp.com",
        "hulu.com",
        "images-amazon.com",
        "live.com",
        "magnetmail1.net",
        "microsoft.com",
        "microsoftonline.com",
        "msdn.com",
        "msecnd.net",
        "msftncsi.com",
        "mywot.com",
        "nsatc.net",
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
        "twimg.com",
        "viewpoint.com",
        "windows.net",
        "xboxlive.com",
        "yimg.com",
        "ytimg.com"
        ],
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
      "sources": [{}]
    },
    "domains": {
      "disabled": "false",
      "excludes": [],
      "includes": [],
      "sources": [{
        "NoBitCoin": {
          "disabled": "false",
          "description": "Blocking Web Browser Bitcoin Mining",
          "prefix": "0.0.0.0",
          "url": "https://raw.githubusercontent.com/hoshsadiq/adblock-nocoin-list/master/hosts.txt",
        },
        "malc0de": {
          "disabled": "false",
          "description": "List of zones serving malicious executables observed by malc0de.com/database/",
          "prefix": "zone",
          "url": "http://malc0de.com/bl/ZONES",
        },
        "malwaredomains.com": {
          "disabled": "false",
          "description": "Just Domains",
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
        },
        "tasty": {
          "disabled": "false",
          "description": "File source",
          "ip": "10.10.10.10",
          "file": "./internal/testdata/blist.hosts.src",
        }
    }]
    },
    "hosts": {
      "disabled": "false",
      "excludes": [],
      "includes": [
        "ads.feedly.com",
        "beap.gemini.yahoo.com"
        ],
      "sources": [{
        "githubSteveBlack": {
          "disabled": "false",
          "description": "Blacklists adware and malware websites",
          "prefix": "0.0.0.0",
          "url": "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
        },
        "hostsfile.org": {
          "disabled": "false",
          "description": "hostsfile.org bad hosts blacklist",
          "prefix": "127.0.0.1",
          "url": "http://www.hostsfile.org/Downloads/hosts.txt",
        },
        "openphish": {
          "disabled": "false",
          "description": "OpenPhish automatic phishing detection",
          "prefix": "http",
          "url": "https://openphish.com/feed.txt",
        },
        "sysctl.org": {
          "disabled": "false",
          "description": "This hosts file is a merged collection of hosts from Cameleon",
          "prefix": "127.0.0.1",
          "url": "http://sysctl.org/cameleon/hosts",
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
    	Enable Debug mode
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
"dropbox.com":0,
"ebay.com":0,
"edgesuite.net":0,
"evernote.com":0,
"express.co.uk":0,
"feedly.com":0,
"freedns.afraid.org":0,
"github.com":0,
"githubusercontent.com":0,
"global.ssl.fastly.net":0,
"google.com":0,
"googleads.g.doubleclick.net":0,
"googleadservices.com":0,
"googleapis.com":0,
"googletagmanager.com":0,
"googleusercontent.com":0,
"gstatic.com":0,
"gvt1.com":0,
"gvt1.net":0,
"hb.disney.go.com":0,
"herokuapp.com":0,
"hp.com":0,
"hulu.com":0,
"images-amazon.com":0,
"live.com":0,
"magnetmail1.net":0,
"microsoft.com":0,
"microsoftonline.com":0,
"msdn.com":0,
"msecnd.net":0,
"msftncsi.com":0,
"mywot.com":0,
"nsatc.net":0,
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
"twimg.com":0,
"viewpoint.com":0,
"windows.net":0,
"xboxlive.com":0,
"yimg.com":0,
"ytimg.com":0,
`
)

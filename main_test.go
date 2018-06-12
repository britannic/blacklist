package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	e "github.com/britannic/blacklist/internal/edgeos"
	"github.com/britannic/mflag"
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

var update = flag.Bool("update", false, "update .golden files")

func readGolden(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name+".golden") // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

func writeGolden(t *testing.T, actual []byte, name string) error {
	golden := filepath.Join("testdata", name+".golden")
	if *update {
		return ioutil.WriteFile(golden, actual, 0644)
	}
	return nil
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
				env, _ = initEnv()
				err = errors.New("initEnvirons failed")
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
	Convey("Testing ScreenLog()", t, func() {
		haveTerm = func() bool {
			return true
		}

		So(screenLog(""), ShouldNotBeNil)
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
	c, _ := initEnv()
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
				errors.New(badFileError),
			)
		})
	})
}

func TestSetArgs(t *testing.T) {
	var (
		origArgs = os.Args
		prog     = path.Base(os.Args[0])
	)

	exitCmd = func(int) {}
	defer func() { os.Args = origArgs }()

	tests := []struct {
		name string
		args []string
		exp  interface{}
	}{
		{
			name: "h",
			args: []string{prog, "-convey-json", "-h"},
			exp:  true,
		},
		{
			name: "debug",
			args: []string{prog, "-debug"},
			exp:  true,
		},
		{
			name: "dryrun",
			args: []string{prog, "-dryrun"},
			exp:  true,
		},
		{
			name: "version",
			args: []string{prog, "-version"},
			exp:  true,
		},
		{
			name: "v",
			args: []string{prog, "-v"},
			exp:  true,
		},
		{
			name: "invalid flag",
			args: []string{prog, "-z"},
			exp:  readGolden(t, "testInvalidArgs"),
		},
	}

	for _, tt := range tests {
		os.Args = nil
		if tt.args != nil {
			os.Args = tt.args
		}

		env := getOpts()
		env.Init(prog, mflag.ContinueOnError)

		Convey("Testing commandline output", t, func() {
			Convey("Testing setArgs() with "+tt.name+"\n", func() {
				switch {
				case tt.name == "invalid flag":
					act := new(bytes.Buffer)
					env.SetOutput(act)
					env.setArgs()
					// *update = true
					writeGolden(t, act.Bytes(), "testInvalidArgs")
					So(act.Bytes(), ShouldResemble, tt.exp.([]byte))
				default:
					env.setArgs()
					So(fmt.Sprint(env.Lookup(tt.name).Value.String()), ShouldEqual, fmt.Sprint(tt.exp))
				}
			})
		})
	}
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

// func TestCommandLineArgs(t *testing.T) {
// 	Convey("Testing command line arguments", t, func() {
// 		origArgs := os.Args
// 		defer func() { os.Args = origArgs }()
// 		act := new(bytes.Buffer)
// 		exitCmd = func(int) {}
// 		exp := vanillaArgs
// 		if IsDrone() {
// 			exp = vanillaArgsOnDrone
// 		}

// 		prog := path.Base(os.Args[0])
// 		os.Args = []string{prog, "-convey-json", "-h"}

// 		o := getOpts()
// 		o.Init("blacklist", mflag.ContinueOnError)
// 		o.SetOutput(act)
// 		o.Parse(cleanArgs(os.Args[1:]))
// 		o.setArgs()

// 		So(act.String(), ShouldEqual, exp)
// 	})
// }

func TestGetCFG(t *testing.T) {
	Convey("Testing getCFG()", t, func() {
		exitCmd = func(int) {}
		o := getOpts()
		c := o.initEdgeOS()

		c.Blacklist(o.getCFG(c))
		So(c.String(), ShouldEqual, mainGetConfig)

		*o.MIPS64 = "amd64"
		c = o.initEdgeOS()
		c.Blacklist(o.getCFG(c))
		So(c.String(), ShouldEqual, "{\n  \"nodes\": [{\n  }]\n}")
	})
}

func TestFiles(t *testing.T) {
	Convey("Testing files()", t, func() {
		exp := ""
		env, _ := initEnv()
		act := files(env)
		So(fmt.Sprintf("%v", act), ShouldEqual, fmt.Sprintf("%v", exp))
	})
}

func TestReloadDNS(t *testing.T) {
	Convey("Testing ReloadDNS()", t, func() {
		var (
			act string
			exp = "[Successfully restarted dnsmasq]"
		)

		// if IsDrone() {
		// 	exp = "ReloadDNS(): [dnsmasq: unrecognized service\n]\n"
		// }

		c, _ := initEnv()
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
		c, _ := initEnv()
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

func TestSetLogFile(t *testing.T) {
	oldprog := prog
	prog = "update-dnsmasq"
	tests := []struct {
		os  string
		exp string
	}{
		{os: "darwin", exp: fmt.Sprintf("/tmp/%s.log", prog)},
		{os: "linux", exp: fmt.Sprintf("/var/log/%s.log", prog)},
	}

	Convey("Testing setLogFile", t, func() {
		for _, tt := range tests {
			Convey("with OS: "+tt.os, func() {
				So(setLogFile(tt.os), ShouldEqual, tt.exp)
			})
		}
	})
	prog = oldprog
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

	expMap = `"1e100.net":{},
"2o7.net":{},
"adobedtm.com":{},
"akamai.net":{},
"akamaihd.net":{},
"amazon.com":{},
"amazonaws.com":{},
"apple.com":{},
"ask.com":{},
"avast.com":{},
"avira-update.com":{},
"bannerbank.com":{},
"bing.com":{},
"bit.ly":{},
"bitdefender.com":{},
"cdn.ravenjs.com":{},
"cdn.visiblemeasures.com":{},
"cloudfront.net":{},
"coremetrics.com":{},
"dropbox.com":{},
"ebay.com":{},
"edgesuite.net":{},
"evernote.com":{},
"express.co.uk":{},
"feedly.com":{},
"freedns.afraid.org":{},
"github.com":{},
"githubusercontent.com":{},
"global.ssl.fastly.net":{},
"google.com":{},
"googleads.g.doubleclick.net":{},
"googleadservices.com":{},
"googleapis.com":{},
"googletagmanager.com":{},
"googleusercontent.com":{},
"gstatic.com":{},
"gvt1.com":{},
"gvt1.net":{},
"hb.disney.go.com":{},
"herokuapp.com":{},
"hp.com":{},
"hulu.com":{},
"images-amazon.com":{},
"live.com":{},
"magnetmail1.net":{},
"microsoft.com":{},
"microsoftonline.com":{},
"msdn.com":{},
"msecnd.net":{},
"msftncsi.com":{},
"mywot.com":{},
"nsatc.net":{},
"paypal.com":{},
"pop.h-cdn.co":{},
"rackcdn.com":{},
"rarlab.com":{},
"schema.org":{},
"shopify.com":{},
"skype.com":{},
"smacargo.com":{},
"sourceforge.net":{},
"spotify.com":{},
"spotify.edgekey.net":{},
"spotilocal.com":{},
"ssl-on9.com":{},
"ssl-on9.net":{},
"sstatic.net":{},
"static.chartbeat.com":{},
"storage.googleapis.com":{},
"twimg.com":{},
"viewpoint.com":{},
"windows.net":{},
"xboxlive.com":{},
"yimg.com":{},
"ytimg.com":{},
`
)

package edgeos

import (
	"bytes"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	logging "github.com/britannic/go-logging"
	. "github.com/smartystreets/goconvey/convey"
)

func newLog() *logging.Logger {
	scrFmt := logging.MustStringFormatter(`%{level:.4s}[%{id:03x}]: %{message}`)
	scr := logging.NewLogBackend(os.Stdout, "", 0)
	scrFmttr := logging.NewBackendFormatter(scr, scrFmt)

	logging.SetBackend(scr, scrFmttr)

	return logging.MustGetLogger("blacklist")
}

func TestEnvLog(t *testing.T) {
	Convey("Testing log()", t, func() {
		tests := []struct {
			dbug bool
			name string
			str  string
		}{
			{name: "Info", str: "This is a log.Info test", dbug: false},
			{name: "Debug", str: "This is a log.Debug test", dbug: true},
			{name: "Error", str: "This is a log.Error test", dbug: true},
			{name: "Warning", str: "This is a log.Warning test", dbug: true},
			{name: "Not Debug", str: "This is a log.Debug test and there shouldn't be any output", dbug: false},
		}

		var (
			scrFmt   = logging.MustStringFormatter(`%{message}`)
			act      = &bytes.Buffer{}
			p        = &Env{Log: logging.MustGetLogger("TestEnvLog"), Verb: true}
			scr      = logging.NewLogBackend(act, "", 0)
			scrFmttr = logging.NewBackendFormatter(scr, scrFmt)
		)

		logging.SetBackend(scrFmttr)

		for _, tt := range tests {
			Convey("Testing "+tt.name, func() {
				p.Dbug = tt.dbug

				switch {
				case tt.dbug:
					p.Debug(tt.str)
					So(act.String(), ShouldEqual, tt.str+"\n")

				case tt.name == "Info":
					p.Log.Info(tt.str)
					So(act.String(), ShouldEqual, tt.str+"\n")

				case tt.name == "Warning":
					exp := tt.str
					exp += tt.str
					p.Log.Warning(tt.str)
					So(act.String(), ShouldEqual, exp)

				case tt.name == "Error":
					exp := tt.str
					exp += tt.str
					p.Log.Error(tt.str)
					So(act.String(), ShouldEqual, exp)

				default:
					p.Debug(tt.str)
					So(act.String(), ShouldEqual, "")
				}
				act.Reset()
			})
		}
	})
}

func TestOption(t *testing.T) {
	Convey("Testing Option()", t, func() {
		vanilla := Env{ctr: ctr{RWMutex: &sync.RWMutex{}, stat: make(stat)}}
		exp := `{
	"Log": null,
	"API": "/bin/cli-shell-api",
	"Arch": "amd64",
	"Bash": "/bin/bash",
	"Cores": 2,
	"Disabled": false,
	"Dbug": true,
	"Dex": {},
	"Dir": "/tmp",
	"dnsmasq service": "service dnsmasq restart",
	"Exc": {},
	"dnsmasq fileExt.": "blacklist.conf",
	"File": "/config/config.boot",
	"File name fmt": "%v/%v.%v.%v",
	"HTTP method": "GET",
	"Prefix": {},
	"Test": true,
	"Timeout": 30000000000,
	"Wildcard": {
		"Node": "*s",
		"Name": "*"
	}
}`

		expRaw := Env{
			ctr:      ctr{RWMutex: &sync.RWMutex{}, stat: make(stat)},
			API:      "/bin/cli-shell-api",
			Arch:     "amd64",
			Bash:     "/bin/bash",
			Cores:    2,
			Disabled: false,
			Dbug:     true,
			Dex:      &list{entry: entry{}},
			Dir:      "/tmp",
			DNSsvc:   "service dnsmasq restart",
			Exc:      &list{entry: entry{}},
			Ext:      "blacklist.conf",
			File:     "/config/config.boot",
			FnFmt:    "%v/%v.%v.%v",
			InCLI:    "inSession",
			Method:   "GET",
			Pfx:      dnsPfx{domain: "address=", host: "server="},
			Test:     true,
			Timeout:  30000000000,
			Wildcard: Wildcard{Node: "*s", Name: "*"},
		}

		c := NewConfig()
		vanilla.Dex = c.Dex
		vanilla.Exc = c.Exc
		So(c.Env, ShouldResemble, &vanilla)

		c = NewConfig(
			Arch(runtime.GOARCH),
			API("/bin/cli-shell-api"),
			Bash("/bin/bash"),
			Cores(2),
			Dbug(true),
			Dir("/tmp"),
			DNSsvc("service dnsmasq restart"),
			Ext("blacklist.conf"),
			File("/config/config.boot"),
			FileNameFmt("%v/%v.%v.%v"),
			InCLI("inSession"),
			Logger(nil),
			Method("GET"),
			Prefix("address=", "server="),
			Test(true),
			Timeout(30*time.Second),
			Verb(false),
			WCard(Wildcard{Node: "*s", Name: "*"}),
			// Writer(nil),
		)

		expRaw.Dex.RWMutex = c.Dex.RWMutex
		expRaw.Exc.RWMutex = c.Exc.RWMutex

		So(*c.Env, ShouldResemble, expRaw)
		So(c.Env.String(), ShouldEqual, exp)
	})
}

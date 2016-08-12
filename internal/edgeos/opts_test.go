package edgeos

import (
	"runtime"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOption(t *testing.T) {
	Convey("Testing Option()", t, func() {
		vanilla := Parms{}

		exp := `{
	"API": "/bin/cli-shell-api",
	"Arch": "amd64",
	"Bash": "/bin/bash",
	"Cores": 2,
	"Debug": true,
	"Dex": {
		"entry": {}
	},
	"Dir": "/tmp",
	"dnsmasq service": "service dnsmasq restart",
	"Exc": {
		"entry": {}
	},
	"dnsmasq fileExt.": "blacklist.conf",
	"File": "/config/config.boot",
	"File name fmt": "%v/%v.%v.%v",
	"CLI Path": "service dns forwarding",
	"Leaf nodes": [
		"file",
		"pre-configured-domain",
		"pre-configured-host",
		"url"
	],
	"HTTP method": "GET",
	"Nodes": [
		"domains",
		"hosts"
	],
	"Prefix": "address=",
	"Poll": 10,
	"Test": true,
	"Timeout": 30000000000,
	"Verbosity": false,
	"Wildcard": {}
}`

		expRaw := Parms{
			API:      "/bin/cli-shell-api",
			Arch:     "amd64",
			Bash:     "/bin/bash",
			Cores:    2,
			Debug:    true,
			Dex:      list{entry: entry{}},
			Dir:      "/tmp",
			DNSsvc:   "service dnsmasq restart",
			Exc:      list{entry: entry{}},
			Ext:      "blacklist.conf",
			File:     "/config/config.boot",
			FnFmt:    "%v/%v.%v.%v",
			InCLI:    "inSession",
			Level:    "service dns forwarding",
			Ltypes:   []string{files, PreDomns, PreHosts, urls},
			Method:   "GET",
			Nodes:    []string{domains, hosts},
			Pfx:      "address=",
			Poll:     10,
			Test:     true,
			Timeout:  30000000000,
			Wildcard: Wildcard{Node: "*s", Name: "*"},
		}

		c := NewConfig()
		vanilla.Dex = c.Dex
		vanilla.Exc = c.Exc
		So(c.Parms, ShouldResemble, &vanilla)

		c = NewConfig(
			Arch(runtime.GOARCH),
			API("/bin/cli-shell-api"),
			Bash("/bin/bash"),
			Cores(2),
			Debug(true),
			Dir("/tmp"),
			DNSsvc("service dnsmasq restart"),
			Ext("blacklist.conf"),
			File("/config/config.boot"),
			FileNameFmt("%v/%v.%v.%v"),
			InCLI("inSession"),
			Logger(nil),
			Method("GET"),
			Nodes([]string{domains, hosts}),
			Poll(10),
			Prefix("address="),
			Level("service dns forwarding"),
			LTypes([]string{"file", PreDomns, PreHosts, urls}),
			Test(true),
			Timeout(30*time.Second),
			Verb(false),
			WCard(Wildcard{Node: "*s", Name: "*"}),
			Writer(nil),
		)

		expRaw.Dex.RWMutex = c.Dex.RWMutex
		expRaw.Exc.RWMutex = c.Exc.RWMutex

		So(*c.Parms, ShouldResemble, expRaw)
		So(c.Parms.String(), ShouldEqual, exp)
	})
}

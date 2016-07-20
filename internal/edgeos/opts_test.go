package edgeos

import (
	"runtime"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOption(t *testing.T) {
	Convey("Testing Option()", t, func() {
		vanilla := Parms{
			API:       "",
			Arch:      "",
			Cores:     0,
			Debug:     false,
			Dex:       list{entry: entry(nil)},
			Dir:       "",
			DNSsvc:    "",
			Exc:       list{entry: entry(nil)},
			Ext:       "",
			File:      "",
			FnFmt:     "",
			InCLI:     "",
			IOWriter:  nil,
			Level:     "",
			Method:    "",
			Nodes:     []string(nil),
			Pfx:       "",
			Poll:      0,
			Ltypes:    []string(nil),
			Test:      false,
			Timeout:   0,
			Verbosity: 0,
			Wildcard:  Wildcard{},
		}

		exp := "edgeos.Parms{\nWildcard:  \"{*s *}\"\nAPI:       \"/bin/cli-shell-api\"\nArch:      \"amd64\"\nBash:      \"/bin/bash\"\nCores:     \"2\"\nDebug:     \"true\"\nDex:       \"**not initialized**\"\nDir:       \"/tmp\"\nDNSsvc:    \"service dnsmasq restart\"\nExc:       \"**not initialized**\"\nExt:       \"blacklist.conf\"\nFile:      \"/config/config.boot\"\nFnFmt:     \"%v/%v.%v.%v\"\nInCLI:     \"inSession\"\nIOWriter:  \"<nil>\"\nLevel:     \"service dns forwarding\"\nLtypes:    \"[file pre-configured-domain pre-configured-host url]\"\nMethod:    \"GET\"\nNodes:     \"[domains hosts]\"\nPfx:       \"address=\"\nPoll:      \"10\"\nTest:      \"true\"\nTimeout:   \"30s\"\nVerbosity: \"2\"\n}\n"

		expRaw := Parms{
			API:       "/bin/cli-shell-api",
			Arch:      "amd64",
			Bash:      "/bin/bash",
			Cores:     2,
			Debug:     true,
			Dex:       list{entry: entry{}},
			Dir:       "/tmp",
			DNSsvc:    "service dnsmasq restart",
			Exc:       list{entry: entry{}},
			Ext:       "blacklist.conf",
			File:      "/config/config.boot",
			FnFmt:     "%v/%v.%v.%v",
			InCLI:     "inSession",
			IOWriter:  nil,
			Level:     "service dns forwarding",
			Ltypes:    []string{files, PreDomns, PreHosts, urls},
			Method:    "GET",
			Nodes:     []string{domains, hosts},
			Pfx:       "address=",
			Poll:      10,
			Test:      true,
			Timeout:   30000000000,
			Verbosity: 2,
			Wildcard:  Wildcard{Node: "*s", Name: "*"},
		}

		So(Parms{}, ShouldResemble, vanilla)

		c := NewConfig(
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
			Method("GET"),
			Nodes([]string{domains, hosts}),
			Poll(10),
			Prefix("address="),
			Level("service dns forwarding"),
			LTypes([]string{"file", PreDomns, PreHosts, urls}),
			Test(true),
			Timeout(30*time.Second),
			Verbosity(2),
			WCard(Wildcard{Node: "*s", Name: "*"}),
			Writer(nil),
		)

		expRaw.Dex.RWMutex = c.Dex.RWMutex
		expRaw.Exc.RWMutex = c.Exc.RWMutex

		So(*c.Parms, ShouldResemble, expRaw)
		So(c.Parms.String(), ShouldEqual, exp)
	})
}

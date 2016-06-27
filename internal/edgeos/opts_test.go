package edgeos

import (
	"runtime"
	"testing"
	"time"

	. "github.com/britannic/testutils"
)

func TestOption(t *testing.T) {
	vanilla := Parms{
		API:       "",
		Arch:      "",
		Cores:     0,
		Debug:     false,
		Dex:       List{},
		Dir:       "",
		DNSsvc:    "",
		Exc:       List{},
		Ext:       "",
		File:      "",
		FnFmt:     "",
		InCLI:     "",
		Level:     "",
		Method:    "",
		Nodes:     []string(nil),
		Pfx:       "",
		Poll:      0,
		Stypes:    []string(nil),
		Test:      false,
		Timeout:   0,
		Verbosity: 0,
		Wildcard:  Wildcard{},
	}

	want := "edgeos.Parms{\nWildcard:  \"{*s *}\"\nAPI:       \"/bin/cli-shell-api\"\nArch:      \"amd64\"\nBash:      \"/bin/bash\"\nCores:     \"2\"\nDebug:     \"true\"\nDex:       \"\\\"synodal.com\\\":0,\"\nDir:       \"/tmp\"\nDNSsvc:    \"service dnsmasq restart\"\nExc:       \"\\\"goodactor.com\\\":0,\"\nExt:       \"blacklist.conf\"\nFile:      \"/config/config.boot\"\nFnFmt:     \"%v/%v.%v.%v\"\nInCLI:     \"inSession\"\nLevel:     \"service dns forwarding\"\nMethod:    \"GET\"\nNodes:     \"[domains hosts]\"\nPfx:       \"address=\"\nPoll:      \"10\"\nStypes:    \"[file pre-configured url]\"\nTest:      \"true\"\nTimeout:   \"30s\"\nVerbosity: \"2\"\n}\n"

	wantRaw := Parms{
		API:       "/bin/cli-shell-api",
		Arch:      "amd64",
		Bash:      "/bin/bash",
		Cores:     2,
		Debug:     true,
		Dex:       List{"synodal.com": 0},
		Dir:       "/tmp",
		DNSsvc:    "service dnsmasq restart",
		Exc:       List{"goodactor.com": 0},
		Ext:       "blacklist.conf",
		File:      "/config/config.boot",
		FnFmt:     "%v/%v.%v.%v",
		InCLI:     "inSession",
		Level:     "service dns forwarding",
		Method:    "GET",
		Nodes:     []string{domains, hosts},
		Pfx:       "address=",
		Poll:      10,
		Stypes:    []string{files, preConf, urls},
		Test:      true,
		Timeout:   30000000000,
		Verbosity: 2,
		Wildcard:  Wildcard{Node: "*s", Name: "*"},
	}

	c := NewConfig()
	Equals(t, vanilla, *c.Parms)

	c = NewConfig(
		Arch(runtime.GOARCH),
		API("/bin/cli-shell-api"),
		Bash("/bin/bash"),
		Cores(2),
		Debug(true),
		Dexcludes(List{"synodal.com": 0}),
		Dir("/tmp"),
		DNSsvc("service dnsmasq restart"),
		Excludes(List{"goodactor.com": 0}),
		Ext("blacklist.conf"),
		File("/config/config.boot"),
		FileNameFmt("%v/%v.%v.%v"),
		InCLI("inSession"),
		Method("GET"),
		Nodes([]string{"domains", "hosts"}),
		Poll(10),
		Prefix("address="),
		Level("service dns forwarding"),
		STypes([]string{"file", preConf, urls}),
		Test(true),
		Timeout(30*time.Second),
		Verbosity(2),
		WCard(Wildcard{Node: "*s", Name: "*"}),
	)

	Equals(t, wantRaw, *c.Parms)

	Equals(t, want, c.Parms.String())
}

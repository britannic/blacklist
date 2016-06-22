package edgeos

import (
	"runtime"
	"testing"

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
		Exc:       List{},
		Ext:       "",
		File:      "",
		FnFmt:     "",
		Level:     "",
		Method:    "",
		Nodes:     []string(nil),
		Pfx:       "",
		Poll:      0,
		Stypes:    []string(nil),
		Test:      false,
		Verbosity: 0,
		Wildcard:  Wildcard{},
	}

	want := "edgeos.Parms{\nAPI:       /bin/cli-shell-api\nArch:      amd64\nCores:     2\nDebug:     true\nDex:       \nDir:       /tmp\nExc:       \"badactor.com\":0,\nExt:       blacklist.conf\nFile:      /config/config.boot\nFnFmt:     %v/%v.%v.%v\nLevel:     service dns forwarding\nMethod:    GET\nNodes:     [domains hosts]\nPfx:       address=\nPoll:      10\nStypes:    [file pre-configured url]\nTest:      true\nVerbosity: 2\nWildcard:  {*s *}\n}\n"

	wantRaw := Parms{
		API:       "/bin/cli-shell-api",
		Arch:      "amd64",
		Cores:     2,
		Debug:     true,
		Dex:       List{},
		Dir:       "/tmp",
		Exc:       List{"badactor.com": 0},
		Ext:       "blacklist.conf",
		File:      "/config/config.boot",
		FnFmt:     "%v/%v.%v.%v",
		Level:     "service dns forwarding",
		Method:    "GET",
		Nodes:     []string{domains, hosts},
		Pfx:       "address=",
		Poll:      10,
		Stypes:    []string{files, preConf, urls},
		Test:      true,
		Verbosity: 2,
		Wildcard:  Wildcard{Node: "*s", Name: "*"},
	}

	c := NewConfig()
	Equals(t, vanilla, *c.Parms)

	c = NewConfig(
		Arch(runtime.GOARCH),
		API("/bin/cli-shell-api"),
		Cores(2),
		Debug(true),
		Dir("/tmp"),
		Excludes(List{"badactor.com": 0}),
		Ext("blacklist.conf"),
		File("/config/config.boot"),
		FileNameFmt("%v/%v.%v.%v"),
		Method("GET"),
		Nodes([]string{"domains", "hosts"}),
		Poll(10),
		Prefix("address="),
		Level("service dns forwarding"),
		STypes([]string{"file", preConf, urls}),
		Test(true),
		Verbosity(2),
		WCard(Wildcard{Node: "*s", Name: "*"}),
	)

	Equals(t, wantRaw, *c.Parms)

	Equals(t, want, c.Parms.String())
}

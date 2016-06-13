package edgeos

import (
	"runtime"
	"testing"

	. "github.com/britannic/testutils"
)

func TestOption(t *testing.T) {
	vanilla := Parms{Arch: "", Cores: 0, Dir: "", Dex: List{}, Debug: false, Exc: List{}, Ext: "", FnFmt: "", File: "", Pfx: "", Method: "", Nodes: []string(nil), Poll: 0, Stypes: []string(nil), Test: false, Verbosity: 0}

	want := "edgeos.Parms{\nArch:      amd64\nCores:     2\nDebug:     true\nDex:       \nDir:       /tmp\nExc:       \"badactor.com\":0,\nExt:       blacklist.conf\nFile:      /config/config.boot\nFnFmt:     %v/%v.%v.%v\nMethod:    GET\nNodes:     [domains hosts]\nPfx:       address=\nPoll:      10\nStypes:    [files pre-configured urls]\nTest:      true\nVerbosity: 2\n}\n"

	wantRaw := Parms{Arch: "amd64", Cores: 2, Dir: "/tmp", Dex: List{}, Debug: true, Exc: List{"badactor.com": 0}, Ext: "blacklist.conf", FnFmt: "%v/%v.%v.%v", File: "/config/config.boot", Pfx: "address=", Method: "GET", Nodes: []string{"domains", "hosts"}, Poll: 10, Stypes: []string{"files", "pre-configured", "urls"}, Test: true, Verbosity: 2}

	c := NewParms()
	Equals(t, vanilla, *c)

	prev := c.SetOpt(
		Arch(runtime.GOARCH),
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
		STypes([]string{"files", preConf, "urls"}),
		Test(true),
		Verbosity(2),
	)

	Equals(t, want, c.String())

	Equals(t, wantRaw, *c)

	c.SetOpt(prev)
	Equals(t, vanilla, *c)
}

package edgeos

import (
	"testing"

	. "github.com/britannic/testutils"
)

func TestOption(t *testing.T) {
	vanilla := parms{cores: 0, debug: false, dir: "", exc: nil, ext: "", file: "", method: "", nodes: nil, poll: 0, stypes: nil, test: false, verbosity: 0}

	want := "edgeos.parms{\ncores:\t\t2\ndir:\t\t\"/tmp\"\ndebug:\t\ttrue\nexc:\t\t\"badactor.com\": 0\next:\t\t\"blacklist.conf\"\nfile:\t\t\"/config/config.boot\"\nmethod:\t\t\"GET\"\nnode:\t\t\"domains\"\nnode:\t\t\"hosts\"\npoll:\t\t10\nstypes:\t\t[files pre-configured urls]\ntest:\t\ttrue\nverbosity:\t'\\x02'\n}\n"

	wantRaw := parms{cores: 2, debug: true, dir: "/tmp", exc: List{"badactor.com": 0}, ext: "blacklist.conf", file: "/config/config.boot", method: "GET", nodes: []string{"domains", "hosts"}, poll: 10, stypes: []string{"files", preConf, "urls"}, test: true, verbosity: 2}

	c := Config{}
	p := NewParms(&c)
	Equals(t, vanilla, *p)

	prev := p.SetOpt(
		Cores(2),
		Debug(true),
		Dir("/tmp"),
		Excludes(List{"badactor.com": 0}),
		Ext("blacklist.conf"),
		File("/config/config.boot"),
		Method("GET"),
		Nodes([]string{"domains", "hosts"}),
		Poll(10),
		STypes([]string{"files", preConf, "urls"}),
		Test(true),
		Verbosity(2),
	)

	Equals(t, want, p.String())

	Equals(t, wantRaw, *p)
	// fmt.Println(want, "\n\n", p.String())

	p.SetOpt(prev)
	Equals(t, vanilla, *p)
}

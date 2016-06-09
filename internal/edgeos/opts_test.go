package edgeos

import (
	"runtime"
	"testing"

	. "github.com/britannic/testutils"
)

func TestOption(t *testing.T) {
	vanilla := Parms{arch: "", cores: 0, debug: false, dir: "", exc: nil, ext: "", file: "", method: "", nodes: nil, poll: 0, stypes: nil, test: false, verbosity: 0}

	want := "edgeos.Parms{\narch:      " + runtime.GOARCH + "\ncores:     2\ndir:       /tmp\ndebug:     true\nexc:       \"badactor.com\":0,\next:       blacklist.conf\nfile:      /config/config.boot\nmethod:    GET\nnodes:     [domains hosts]\npoll:      10\nstypes:    [files pre-configured urls]\ntest:      true\nverbosity: 2\n}\n"

	wantRaw := Parms{arch: runtime.GOARCH, cores: 2, debug: true, dir: "/tmp", exc: List{"badactor.com": 0}, ext: "blacklist.conf", file: "/config/config.boot", method: "GET", nodes: []string{"domains", "hosts"}, poll: 10, stypes: []string{"files", preConf, "urls"}, test: true, verbosity: 2}

	c := Config{}
	p := NewParms(&c)
	Equals(t, vanilla, *p)

	prev := p.SetOpt(
		Arch(runtime.GOARCH),
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
	// fmt.Println(p.String())
	Equals(t, wantRaw, *p)

	p.SetOpt(prev)
	Equals(t, vanilla, *p)
}

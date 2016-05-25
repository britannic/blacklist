package edgeos

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/britannic/blacklist/tdata"
	. "github.com/britannic/testutils"
)

// func check(t *testing.T) {
// 	if !true {
// 		t.Skip("Not implemented; skipping tests")
// 	}
// }

// func uDiff(a, b string) string {
// 	diff := difflib.ContextDiff{
// 		A:        difflib.SplitLines(a),
// 		B:        difflib.SplitLines(b),
// 		FromFile: "Want",
// 		ToFile:   "Got",
// 		Context:  3,
// 		Eol:      "\n",
// 	}

// 	result, _ := difflib.GetContextDiffString(diff)
// 	return fmt.Sprintf(strings.Replace(result, "\t", " ", -1))
// }

func TestGetSubdomains(t *testing.T) {
	d := GetSubdomains("top.one.two.three.four.five.six.intellitxt.com")

	for key := range d {
		Assert(t, d.KeyExists(key), fmt.Sprintf("%v key doesn't exist", key), d)
	}
}

func TestGetType(t *testing.T) {
	tests := []struct {
		typeint int
		typestr string
	}{
		{typeint: unknown, typestr: Unknown},
		{typeint: root, typestr: blacklist},
		{typeint: pre, typestr: PreConf},
		{typeint: domain, typestr: Domains},
		{typeint: host, typestr: Hosts},
	}

	for _, test := range tests {
		Equals(t, test.typeint, typeStr(test.typestr))
		Equals(t, test.typestr, typeInt(test.typeint))
		Equals(t, test.typestr, getType(test.typeint))
		Equals(t, test.typeint, getType(test.typestr))
	}
}

func TestLoad(t *testing.T) {
	cfg, err := Load("zBroken", "service dns forwarding")
	NotOK(t, err)

	cfg, err = Load("showConfig", "")
	NotOK(t, err)

	pwd, err := os.Getwd()
	OK(t, err)

	cfg, err = Load("pwd", "")
	OK(t, err)

	got := new(bytes.Buffer)
	got.ReadFrom(cfg)
	Equals(t, pwd+"\n", got.String())
}

func TestReadCfg(t *testing.T) {
	b, err := ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)

	Equals(t, tdata.JSONcfg, b.String())
	// fmt.Println(uDiff(tdata.JSONcfg, b.String()))

	Equals(t, tdata.JSONrawcfg, b.JSON())
	// fmt.Println(uDiff(tdata.JSONrawcfg, b.JSON()))
	// fmt.Println(b.JSON())
	b, err = ReadCfg(bytes.NewBufferString(tdata.ZeroHostSourcesCfg))
	OK(t, err)
	// fmt.Println(b)

	Equals(t, tdata.JSONcfgZeroHostSources, b.String())
	// fmt.Println(uDiff(tdata.JSONcfgZeroHostSources, b.String()))

	b, got := ReadCfg(bytes.NewBufferString(""))
	NotEquals(t, nil, got)

	want := errors.New("Configuration data is empty, cannot continue")
	Equals(t, want.Error(), got.Error())
	Equals(t, Nodes{}, b)
	// fmt.Println(b)

	b, err = ReadCfg(bytes.NewBufferString(strippedCfg))
	OK(t, err)

	type dataSrc struct {
		desc   string
		file   string
		ip     string
		name   string
		prefix string
		run    bool
		url    string
	}

	type testSrc struct {
		disabled bool
		ip       string
		node     string
		s        dataSrc
	}

	tests := []testSrc{
		{
			disabled: false,
			ip:       "1.1.1.1",
			node:     blacklist,
			s: dataSrc{
				run: false,
			},
		},
		{
			disabled: false,
			ip:       "2.2.2.2",
			node:     Domains,
			s: dataSrc{
				desc:   "List of zones serving malicious executables observed by malc0de.com/database/",
				ip:     "4.4.4.4",
				name:   "malc0de",
				prefix: "zone ",
				run:    true,
				url:    "http://malc0de.com/bl/ZONES",
			},
		},
		{
			disabled: true,
			ip:       "3.3.3.3",
			node:     Hosts,
			s: dataSrc{
				desc: "File test",
				file: "/test/file",
				ip:   "5.5.5.5",
				name: "file",
				run:  true,
			},
		},
	}

	for _, test := range tests {
		Equals(t, test.disabled, b[test.node].Disabled)
		Equals(t, test.ip, b[test.node].IP)
		if test.s.run {
			Equals(t, test.s.ip, b[test.node].Data[test.s.name].IP)
			Equals(t, test.s.desc, b[test.node].Data[test.s.name].Desc)
			Equals(t, test.s.prefix, b[test.node].Data[test.s.name].Prefix)
			Equals(t, test.s.url, b[test.node].Data[test.s.name].URL)
		}
	}
}

// type TestRunCMDline struct{}
//
// func (r TestRunCMDline) Command(command string, args ...string) *exec.Cmd {
// 	cs := []string{"-test.run=TestCommand", "--"}
// 	cs = append(cs, args...)
// 	out := exec.Command(os.Args[0], cs...)
// 	out.Env = []string{"GO_WANT_COMMAND=1"}
// 	return out
// }

func TestToBool(t *testing.T) {
	tests := map[string]bool{"false": false, "true": true, "fail": false}

	for k := range tests {
		Equals(t, tests[k], ToBool(k))
	}
}

var (
	keys = []string{
		"six.com",
		"five.six.com",
		"four.five.six.com",
		"three.four.five.six.com",
		"two.three.four.five.six.com",
		"one.two.three.four.five.six.com",
		"top.one.two.three.four.five.six.com",
	}

	strippedCfg = `blacklist {
	disabled false
	dns-redirect-ip 1.1.1.1
	domains {
	disabled false
			dns-redirect-ip 2.2.2.2
			source malc0de {
					description "List of zones serving malicious executables observed by malc0de.com/database/"
					dns-redirect-ip 4.4.4.4
					prefix "zone "
					url http://malc0de.com/bl/ZONES
			}
	}
	hosts {
			disabled true
			dns-redirect-ip 3.3.3.3
			source file {
				description "File test"
				dns-redirect-ip 5.5.5.5
				file /test/file
		}
	}
}`
)

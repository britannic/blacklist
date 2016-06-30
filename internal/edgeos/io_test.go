package edgeos

import (
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/britannic/testutils"
)

func TestLoad(t *testing.T) {
	c := NewConfig(
		API("/bin/cli-shell-api"),
		Bash("/bin/bash"),
		InCLI("inSession"),
		Level("service dns forwarding"),
	)

	cfg, err := c.load("zBroken", "service dns forwarding")
	NotOK(t, err)

	cfg, err = c.load("showConfig", "")
	NotOK(t, err)

	r := CFGcli{Config: c}
	got, err := ioutil.ReadAll(r.read())
	OK(t, err)
	Equals(t, "", string(got))

	cfg, err = c.load("echo", "true")
	NotOK(t, err)

	Equals(t, []byte{}, cfg)
}

func TestPurgeFiles(t *testing.T) {
	var (
		dir       = "/tmp"
		ext       = ".delete"
		purgeList []string
		want      error
	)

	for i := 0; i < 10; i++ {
		fname := fmt.Sprintf("%v%v", i, ext)
		f, err := ioutil.TempFile(dir, fname)
		OK(t, err)
		purgeList = append(purgeList, f.Name())
	}

	err := purgeFiles(purgeList)
	OK(t, err)

	got := purgeFiles(purgeList)
	Equals(t, want, got)

	got = purgeFiles([]string{"/dev/null"})
	want = fmt.Errorf(`could not remove "/dev/null"`)
	Equals(t, want, got)
}

func TestAPICMD(t *testing.T) {
	type query struct {
		b    bool
		q, r string
	}
	testSrc := []*query{
		{
			b: false,
			q: "listNodes",
			r: "listNodes",
		},
		{
			b: true,
			q: "listNodes",
			r: "listActiveNodes",
		},
		{
			b: false,
			q: "listActiveNodes",
			r: "listNodes",
		},
		{
			b: false,
			q: "returnValue",
			r: "returnValue",
		},
		{
			b: true,
			q: "returnValue",
			r: "returnActiveValue",
		},
		{
			b: false,
			q: "returnActiveValue",
			r: "returnValue",
		},
		{
			b: false,
			q: "returnValues",
			r: "returnValues",
		},
		{
			b: true,
			q: "returnValues",
			r: "returnActiveValues",
		},
		{
			b: false,
			q: "returnActiveValues",
			r: "returnValues",
		},
		{
			b: false,
			q: "exists",
			r: "exists",
		},
		{
			b: true,
			q: "exists",
			r: "existsActive",
		},
		{
			b: false,
			q: "existsActive",
			r: "exists",
		},
		{
			b: false,
			q: "showCfg",
			r: "showCfg",
		},
		{
			b: true,
			q: "showCfg",
			r: "showConfig",
		},
		{
			b: false,
			q: "showConfig",
			r: "showCfg",
		},
	}

	for _, rq := range testSrc {
		Equals(t, rq.r, apiCMD(rq.q, rq.b))
	}

	c := NewConfig(
		API("/bin/cli-shell-api"),
		InCLI("inSession"),
		Level("service dns forwarding"),
	)
	act := fmt.Sprintf("%v %v", apiCMD("showConfig", c.InSession()), c.Level)
	exp := "showCfg service dns forwarding"
	Equals(t, exp, act)

}

func TestDeleteFile(t *testing.T) {
	dir := "../testdata"
	ext := "delete.me"

	tests := []struct {
		name string
		f    string
		exp  bool
	}{
		{
			name: "exists",
			f:    fmt.Sprintf("%v%v", "goodFile", ext),
			exp:  true,
		},
		{
			name: "non-existent",
			f:    fmt.Sprintf("%v%v", "badFile", ext),
			exp:  true,
		},
	}
	for _, tt := range tests {
		switch tt.name {
		case "exists":
			f, err := ioutil.TempFile(dir, tt.f)
			OK(t, err)
			Equals(t, tt.exp, deleteFile(f.Name()))
		default:
			Equals(t, tt.exp, deleteFile(tt.f))
		}
	}
}

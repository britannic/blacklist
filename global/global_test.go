package global_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	g "github.com/britannic/blacklist/global"
	. "github.com/britannic/testutils"
)

var (
	log, hook = test.NewNullLogger()
)

type testTable struct {
	test interface{}
	exp  interface{}
}

type cpu struct {
	arch  string
	check []testTable
}

func TestLog2Stdout(t *testing.T) {
	s := &g.Set{
		Level:  logrus.InfoLevel,
		Output: "screen",
	}

	g.LogInit(s)
	log.Info("TestLog2Stdout")

	Equals(t, "TestLog2Stdout", hook.LastEntry().Message)
	Equals(t, log.Level.String(), hook.LastEntry().Level.String())
}

func TestLog2File(t *testing.T) {
	s := &g.Set{
		File:   "/tmp/log_test.log",
		Level:  logrus.DebugLevel,
		Output: "file",
	}

	g.LogInit(s)

	if _, err := os.Stat(s.File); os.IsNotExist(err) {
		OK(t, err)
	}

	log.Info("TestLog2Stdout")

	Equals(t, "TestLog2Stdout", hook.LastEntry().Message)
	Equals(t, log.Level.String(), hook.LastEntry().Level.String())

	_ = os.Remove(s.File)
}

func TestSetVars(t *testing.T) {
	cwd, err := os.Getwd()
	OK(t, err)

	platforms := []cpu{
		{arch: "amd64",
			check: []testTable{
				{test: &g.Area.Domains, exp: "domains"},
				{test: &g.Area.Hosts, exp: "hosts"},
				{test: &g.Area.Root, exp: "blacklist"},
				{test: &g.Args, exp: []string{""}},
				{test: &g.Dbg, exp: false},
				{test: &g.DmsqDir, exp: cwd + "/testdata"},
				{test: &g.DNSRestart, exp: fmt.Sprintf("echo -n dnsmasq not implemented on %v", g.TestOS)},
				{test: g.Fext, exp: ".blacklist.conf"},
				{test: &g.FStr, exp: `%v/%v.%v` + g.Fext},
				{test: &g.LogFile, exp: fmt.Sprintf("%v/blacklist.log", g.DmsqDir)},
				{test: g.TestOS, exp: "darwin"},
			},
		},
		{arch: "mips64",
			check: []testTable{
				{test: &g.Area.Domains, exp: "domains"},
				{test: &g.Area.Hosts, exp: "hosts"},
				{test: &g.Area.Root, exp: "blacklist"},
				{test: &g.Args, exp: []string{""}},
				{test: &g.Dbg, exp: false},
				{test: &g.DmsqDir, exp: "/etc/dnsmasq.d"},
				{test: &g.DNSRestart, exp: "service dnsmasq restart"},
				{test: g.Fext, exp: ".blacklist.conf"},
				{test: &g.FStr, exp: `%v/%v.%v` + g.Fext},
				{test: &g.LogFile, exp: "/var/log/blacklist.log"},
				{test: g.TargetOS, exp: "linux"},
			},
		},
	}

	for _, k := range platforms {
		g.SetVars(k.arch)

		for _, run := range k.check {
			switch run.test.(type) {
			case bool:
				Equals(t, run.exp.(bool), run.test.(bool))

			case string:
				Equals(t, run.exp.(string), run.test.(string))

			case int:
				Equals(t, run.exp.(int), run.test.(int))

			case nil:
				fmt.Println("Test not properly defined! ", run)
			}
		}
	}
}

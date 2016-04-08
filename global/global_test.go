package global_test

import (
	"fmt"
	"os"
	"runtime"
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
	alt  interface{}
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
	platforms := []string{"amd64", "mips64"}

	tests := []testTable{
		{test: &g.Area.Domains, exp: "domains", alt: "domains"},
		{test: &g.Area.Hosts, exp: "hosts", alt: "hosts"},
		{test: &g.Area.Root, exp: "blacklist", alt: "blacklist"},
		{test: &g.Args, exp: []string{""}, alt: []string{""}},
		{test: &g.Dbg, exp: false, alt: false},
		{test: &g.DmsqDir, exp: "/etc/dnsmasq.d", alt: cwd + "/testdata"},
		{test: &g.DNSRestart, exp: "service dnsmasq restart", alt: fmt.Sprintf("echo -n dnsmasq not implemented on %v", g.TestOS)},
		{test: g.Fext, exp: ".blacklist.conf", alt: ".blacklist.conf"},
		{test: &g.FStr, exp: `%v/%v.%v` + g.Fext, alt: `%v/%v.%v` + g.Fext},
		{test: &g.LogFile, exp: "/var/log/blacklist.log", alt: fmt.Sprintf("%v/blacklist.log", g.DmsqDir)},
		{test: g.TestOS, exp: "darwin", alt: "darwin"},
	}

	g.WhatOS = runtime.GOOS
	Arch := g.WhatArch
	for _, Arch = range platforms {
		g.SetVars(Arch)

		for _, run := range tests {
			expect := run.exp
			if Arch == g.TargetArch {
				expect = run.alt
			}

			switch run.test.(type) {
			case bool:
				Equals(t, expect.(bool), run.test.(bool))

			case string:
				Equals(t, expect.(string), run.test.(string))

			case int:
				Equals(t, expect.(int), run.test.(int))

			case nil:
				fmt.Println("Test not properly defined! ", run)
			}
		}
	}
}

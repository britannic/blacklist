package utils_test

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"testing"

	"github.com/Sirupsen/logrus"
	tlog "github.com/Sirupsen/logrus/hooks/test"
	"github.com/britannic/blacklist/global"
	"github.com/britannic/blacklist/utils"
	. "github.com/britannic/testutils"
)

var (
	// d.String() = "The rest is history!"
	d         = []byte{84, 104, 101, 32, 114, 101, 115, 116, 32, 105, 115, 32, 104, 105, 115, 116, 111, 114, 121, 33}
	dmsqDir   string
	log, hook = tlog.NewNullLogger()
)

func init() {
	global.Log = log
	switch global.WhatArch {
	case global.TargetArch:
		dmsqDir = global.DmsqDir
	default:
		dmsqDir = "../testdata"
	}
	s := &utils.Set{
		Output: global.LogOutput,
		Level:  logrus.DebugLevel,
		Log:    log,
	}
	utils.Log2File(s)
}

// cmpID compares two UIDS
func cmpID(t *testing.T, want, got *user.User) {
	if want.Uid != got.Uid {
		t.Errorf("got Uid=%q; want %q", got.Uid, want.Uid)
	}
}

func TestBasename(t *testing.T) {
	dir := utils.Basename("/usr/blacklist/testing.txt")
	Equals(t, "testing", dir)
}

func TestCmpHash(t *testing.T) {
	//MD5 ("The rest is history!") = 0ba11c67af902879d20130d9ab414771
	want := d
	got := want

	Assert(t, utils.CmpHash(want, got), "Cmphash() failed", got, want)

	got = append(got, "This is different!"...)
	Assert(t, !utils.CmpHash(want, got), "Cmphash() failed", got, want)
}

func TestByteArray(t *testing.T) {
	var (
		want = d
		got  []byte
	)

	f := fmt.Sprintf("%v/delete.txt", dmsqDir)
	err := utils.WriteFile(f, d)
	OK(t, err)

	b, err := utils.GetFile(f)
	OK(t, err)

	_ = os.Remove(f)

	got = utils.GetByteArray(b, got)
	Equals(t, want, got)
}

func TestGetFile(t *testing.T) {
	var (
		want = d
		got  []byte
	)

	f := fmt.Sprintf("%v/delete.txt", dmsqDir)
	err := utils.WriteFile(f, d)
	OK(t, err)

	b, err := utils.GetFile(f)
	OK(t, err)

	_ = os.Remove(f)

	got = utils.GetByteArray(b, got)
	Equals(t, want, got)
}

func TestLog2Stdout(t *testing.T) {
	// logger, hook := test.NewNullLogger()
	s := &utils.Set{
		Level:  logrus.InfoLevel,
		Log:    log,
		Output: "screen",
	}

	utils.LogInit(s)
	log.Info("TestLog2Stdout")

	Equals(t, "TestLog2Stdout", hook.LastEntry().Message)
	Equals(t, log.Level.String(), hook.LastEntry().Level.String())
	// hook.Reset()
}

func TestLog2File(t *testing.T) {
	s := &utils.Set{
		File:   "/tmp/log_test.log",
		Level:  logrus.DebugLevel,
		Log:    log,
		Output: "file",
	}

	utils.LogInit(s)

	if _, err := os.Stat(s.File); os.IsNotExist(err) {
		OK(t, err)
	}

	log.Info("TestLog2Stdout")

	Equals(t, "TestLog2Stdout", hook.LastEntry().Message)
	Equals(t, log.Level.String(), hook.LastEntry().Level.String())

	_ = os.Remove(s.File)
}

func TestReloadDNS(t *testing.T) {
	tests := []struct {
		test   string
		expect bool
		want   string
	}{
		{
			test:   "echo Testing",
			expect: true,
			want:   "Testing\n",
		},
		{
			test:   "ztaswerkjlkjsdflkjsdf Testing",
			expect: false,
			want:   "/bin/bash: line 1: ztaswerkjlkjsdflkjsdf: command not found\n",
		},
	}

	for _, run := range tests {
		s, err := utils.ReloadDNS(run.test)
		switch run.expect {
		case false:
			Assert(t, err != nil, fmt.Sprint("Test should fail, so ReloadDNS() error shouldn't be nil!"), err)

		case true:
			log.Infof("testing: %v", run)
			Assert(t, err == nil, fmt.Sprint("Test should pass, so ReloadDNS() error should be nil!"), err)
		}

		Equals(t, s, run.want)
	}

}

func TestStringArray(t *testing.T) {
	var (
		b    *bufio.Scanner
		got  []string
		want = []string{"The rest is history!"}
	)

	f := fmt.Sprintf("%v/delete.txt", dmsqDir)
	err := utils.WriteFile(f, d)
	OK(t, err)

	b, err = utils.GetFile(f)
	OK(t, err)

	_ = os.Remove(f)

	got = utils.GetStringArray(b, got)
	Equals(t, want, got)
}

func TestIsAdmin(t *testing.T) {
	want, err := user.Current()
	OK(t, err)

	got, err := user.Lookup(want.Username)
	OK(t, err)

	Equals(t, want, got)

	osAdmin := false
	if got.Uid == "0" {
		osAdmin = true
	}

	switch {
	case !osAdmin:
		Assert(t, !utils.IsAdmin(), fmt.Sprintf("Should be standard user, got: %v", got.Uid), want)

	case osAdmin:
		Assert(t, utils.IsAdmin(), fmt.Sprintf("Should be root user, got: %v", got.Uid), want)
	}
}

func TestWriteFile(t *testing.T) {
	tFile := struct {
		badfile string
		tdata   []byte
		tdir    string
		tfile   string
	}{
		badfile: "/tmp/z/d/c/r/c:reallybadfile.zfts",
		tdata:   d,
		tdir:    "/tmp",
		tfile:   "Test.util.WriteFile",
	}

	f, err := ioutil.TempFile(tFile.tdir, tFile.tfile)
	OK(t, err)

	defer os.Remove(f.Name())
	defer f.Close()

	err = utils.WriteFile(f.Name(), d)
	OK(t, err)

	err = utils.WriteFile(tFile.badfile, d)
	NotOK(t, err)
}

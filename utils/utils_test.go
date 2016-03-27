package utils_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"reflect"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/britannic/blacklist/global"
	"github.com/britannic/blacklist/utils"
)

var (
	// d.String() = "The rest is history!"
	d       = []byte{84, 104, 101, 32, 114, 101, 115, 116, 32, 105, 115, 32, 104, 105, 115, 116, 111, 114, 121, 33}
	dmsqDir string
)

func init() {
	switch global.WhatOS {
	case global.TestOS:
		dmsqDir = "../testdata"
	default:
		dmsqDir = global.DmsqDir
	}
}

// cmpID compares two UIDS
func cmpID(t *testing.T, want, got *user.User) {
	if want.Uid != got.Uid {
		t.Errorf("got Uid=%q; want %q", got.Uid, want.Uid)
	}
}

func TestBasename(t *testing.T) {
	dir := utils.Basename("/usr/blacklist/testing.txt")
	if dir != "testing" {
		t.Error(dir)
	}
}

func TestCmpHash(t *testing.T) {
	//MD5 ("The rest is history!") = 0ba11c67af902879d20130d9ab414771
	want := d

	got := want
	if !utils.CmpHash(want, got) {
		t.Error("Failed!")
	}

	got = append(got, "This is different!"...)
	if utils.CmpHash(want, got) {
		t.Error("Failed!")
	}
}

func TestByteArray(t *testing.T) {
	var (
		want = d
		got  []byte
	)

	f := fmt.Sprintf("%v/delete.txt", dmsqDir)
	if err := utils.WriteFile(f, d); err != nil {
		t.Errorf("Error writing file: %v", err)
	}

	b, err := utils.GetFile(f)
	if err != nil {
		t.Errorf("Failed with error: %v", err)
	}

	_ = os.Remove(f)

	got = utils.GetByteArray(b, got)

	if !utils.CmpHash(want, got) {
		t.Errorf("Data not the same - Got: %v\nWant: %v\n", string(got[:]), string(want[:]))
	}
}

func TestGetFile(t *testing.T) {
	var (
		want = d
		got  []byte
	)

	f := fmt.Sprintf("%v/delete.txt", dmsqDir)
	if err := utils.WriteFile(f, d); err != nil {
		t.Errorf("Error writing file: %v", err)
	}

	b, err := utils.GetFile(f)
	if err != nil {
		t.Errorf("Failed with error: %v", err)
	}

	_ = os.Remove(f)

	got = utils.GetByteArray(b, got)
	if !utils.CmpHash(want, got) {
		t.Errorf("Data not the same - Got: %v\nWant: %v\n", string(got[:]), string(want[:]))
	}
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
		{
			test:   "which cd",
			expect: true,
			want:   "/usr/bin/cd\n",
		},
		{
			test:   "file /etc/services",
			expect: true,
			want:   "/etc/services: ASCII English text\n",
		},
	}

	for _, run := range tests {
		s, err := utils.ReloadDNS(run.test)
		switch run.expect {
		case false:
			if err == nil {
				t.Errorf("Test should fail, so ReloadDNS() error shouldn't be nil!")
			}
		case true:
			if err != nil {
				t.Errorf("Test should pass, so ReloadDNS() error should be nil! Error: %v", err)
			}
		}

		if s != run.want {
			t.Errorf("Want: %q, Got: %q", run.want, s)
		}
	}

}

func TestStringArray(t *testing.T) {
	var (
		got  []string
		want = []string{"The rest is history!"}
	)

	f := fmt.Sprintf("%v/delete.txt", dmsqDir)
	if err := utils.WriteFile(f, d); err != nil {
		t.Errorf("Error writing file: %v", err)
	}

	b, err := utils.GetFile(f)
	if err != nil {
		t.Errorf("Failed with error: %v", err)
	}

	_ = os.Remove(f)

	got = utils.GetStringArray(b, got)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Data not the same - Got: %v\nWant: %v\n", got, want)
	}
}

func TestIsAdmin(t *testing.T) {
	want, err := user.Current()
	if err != nil {
		t.Errorf("Current: %v", err)
	}

	got, err := user.Lookup(want.Username)
	if err != nil {
		t.Errorf("Lookup: %v", err)
	}

	cmpID(t, want, got)

	osAdmin := false
	if got.Uid == "0" {
		osAdmin = true
	}

	switch {
	case !osAdmin && utils.IsAdmin():
		t.Errorf("Standard user: %v", got)
	case osAdmin && !utils.IsAdmin():
		t.Errorf("Root: %v", got)
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
	if err != nil {
		log.Errorf("open %s file: %s", f.Name(), err)
	}

	defer os.Remove(f.Name())
	defer f.Close()

	err = utils.WriteFile(f.Name(), d)
	if err != nil {
		t.Errorf("Error writing %s file: %s", f.Name(), err)
	}

	err = utils.WriteFile(tFile.badfile, d)
	if err == nil {
		t.Errorf("Should not be able to write %s file: %v", tFile.badfile, err)
	}

}

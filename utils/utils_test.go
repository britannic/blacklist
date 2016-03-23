package utils_test

import (
	"fmt"
	"os"
	"os/user"
	"reflect"
	"testing"

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

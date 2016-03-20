package utils_test

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/britannic/blacklist/data"
	"github.com/britannic/blacklist/global"
	"github.com/britannic/blacklist/utils"
)

var (
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

func TestGetFile(t *testing.T) {
	want := []string{"The rest is history!"}

	f := fmt.Sprintf("%v/delete.txt", dmsqDir)
	if err := utils.WriteFile(f, d); err != nil {
		t.Errorf("Error writing file: %v", err)
	}

	got, err := utils.Getfile(f)
	if err != nil {
		t.Errorf("Failed with error: %v", err)
	}

	if delta := data.DiffArray(want, got); len(delta) > 0 {
		t.Errorf("Data not the same - difference: %v", delta)
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

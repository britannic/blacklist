package utils_test

import (
	"os/user"
	"testing"

	"github.com/britannic/blacklist/utils"
)

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

	if utils.IsAdmin() != osAdmin {
		t.Error(osAdmin)
	}
}

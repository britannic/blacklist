package config_test

import (
	"testing"

	"github.com/britannic/blacklist/config"
)

func TestBlacklistCfg(t *testing.T) {
	b, e := config.Get(config.Testdata, "blacklist")
	if e != nil {
		t.Error(b)
	}
}

func TestGetSubdomains(t *testing.T) {
	d := config.GetSubdomains("top.one.two.three.four.five.six.com")

	for _, key := range keys {
		if !d.KeyExists(key) {
			t.Errorf("%v key doesn't exist", key)
		}
	}
}

func TestKeyExists(t *testing.T) {
	d := config.GetSubdomains("top.one.two.three.four.five.six.com")

	for _, key := range keys {
		if !d.KeyExists(key) {
			t.Errorf("%v key doesn't exist", key)
		}
	}
}

func TestSubKeyExists(t *testing.T) {
	d := config.GetSubdomains("top.one.two.three.four.five.six.com")
	dLen := len(d)
	for i, key := range keys {
		if i < dLen {
			if !d.SubKeyExists(keys[i]) {
				t.Errorf("%v sub key doesn't exist", key)
			}
		}
	}
}

func TestToBool(t *testing.T) {
	if b := config.ToBool("false"); b {
		t.Errorf(`ToBool("false") `+"failed with %v\n", b)
	}
	if b := config.ToBool("true"); !b {
		t.Errorf(`ToBool("true") `+"failed with %v\n", b)
	}
}

var keys = []string{
	"six.com",
	"five.six.com",
	"four.five.six.com",
	"three.four.five.six.com",
	"two.three.four.five.six.com",
	"one.two.three.four.five.six.com",
	"top.one.two.three.four.five.six.com",
}

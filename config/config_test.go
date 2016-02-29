// Package config_test
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

func TestToBool(t *testing.T) {
	if b := config.ToBool("false"); b {
		t.Errorf(`ToBool("false") `+"failed with %v\n", b)
	}
	if b := config.ToBool("true"); !b {
		t.Errorf(`ToBool("true") `+"failed with %v\n", b)
	}
}

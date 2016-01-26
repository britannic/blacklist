// Package config_test
package config_test

import (
	"fmt"
	"testing"

	"github.com/britannic/blacklist/config"
)

func TestBlacklistCfg(t *testing.T) {
	b, e := config.Get(config.Testdata, "blacklist")
	if e != nil {
		fmt.Println(b)
	}
}

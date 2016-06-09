package edgeos

import (
	"bytes"

	"github.com/britannic/blacklist/internal/tdata"
)

// Mvars is a struct of initial
type Mvars struct {
	DNSdir   string
	DNStmp   string
	MIPS64   string
	WhatOS   string
	WhatArch string
}

// SetDir sets the directory according to the host CPU arch
func (m *Mvars) SetDir(arch string) (dir string) {
	switch arch {
	case m.MIPS64:
		dir = m.DNSdir
	default:
		dir = m.DNStmp
	}
	return dir
}

// GetCFG returns a *Config
func (m *Mvars) GetCFG(arch string) (c *Config, err error) {
	var cfg string
	c = &Config{}
	switch arch {
	case m.MIPS64:
		if cfg, err = LoadCfg(); err != nil {
			return c, err
		}
		c, err = ReadCfg(bytes.NewBufferString(cfg))
	default:
		c, err = ReadCfg(bytes.NewBufferString(tdata.Cfg))
	}
	return c, err
}

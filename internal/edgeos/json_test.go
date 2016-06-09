package edgeos

import (
	"bytes"
	"testing"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/britannic/testutils"
)

func TestConfigString(t *testing.T) {
	c, err := ReadCfg(bytes.NewBufferString(Cfg))
	OK(t, err)
	NewParms(c).SetOpt(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Method("GET"),
		Nodes([]string{"domains", "hosts"}),
		STypes([]string{"pre-configured", "files", "urls"}),
	)
	Equals(t, tdata.JSONcfg, c.String())

	c, err = ReadCfg(bytes.NewBufferString(tdata.ZeroHostSourcesCfg))
	OK(t, err)
	NewParms(c).SetOpt(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Method("GET"),
		Nodes([]string{"domains", "hosts"}),
		STypes([]string{"pre-configured", "files", "urls"}),
	)
	Equals(t, tdata.JSONcfgZeroHostSources, c.String())
}

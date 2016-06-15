package edgeos

import (
	"testing"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/britannic/testutils"
)

func TestConfigString(t *testing.T) {
	r := &CFGstatic{Cfg: tdata.Cfg}
	c, err := ReadCfg(r)
	OK(t, err)
	c.Parms = NewParms()
	c.SetOpt(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Method("GET"),
		Nodes([]string{"domains", "hosts"}),
		STypes([]string{"pre-configured", "files", "urls"}),
	)
	Equals(t, tdata.JSONcfg, c.String())

	r = &CFGstatic{Cfg: tdata.ZeroHostSourcesCfg}
	c, err = ReadCfg(r)
	OK(t, err)
	c.Parms = NewParms()
	c.SetOpt(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Method("GET"),
		Nodes([]string{"domains", "hosts"}),
		STypes([]string{"pre-configured", "files", "urls"}),
	)
	Equals(t, tdata.JSONcfgZeroHostSources, c.String())
}

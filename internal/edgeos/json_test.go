package edgeos

import (
	"testing"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/britannic/testutils"
)

func TestConfigString(t *testing.T) {
	c := NewConfig(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Method("GET"),
		Nodes([]string{"domains", "hosts"}),
		STypes([]string{"pre-configured", "file", urls}),
	)

	r := &CFGstatic{Cfg: tdata.Cfg}
	err := c.ReadCfg(r)
	OK(t, err)

	Equals(t, tdata.JSONcfg, c.String())

	r = &CFGstatic{Cfg: tdata.ZeroHostSourcesCfg}
	c = NewConfig(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Method("GET"),
		Nodes([]string{"domains", "hosts"}),
		STypes([]string{"pre-configured", "file", urls}),
	)

	err = c.ReadCfg(r)
	OK(t, err)

	Equals(t, tdata.JSONcfgZeroHostSources, c.String())
}

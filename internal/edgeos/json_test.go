package edgeos

import (
	"testing"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConfigString(t *testing.T) {
	Convey("Testing ConfigString()", t, func() {
		c := NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
			Method("GET"),
			LTypes([]string{files, PreDomns, PreHosts, urls}),
		)

		So(c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)
		So(c.String(), ShouldEqual, tdata.JSONcfg)

		c = NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
			Method("GET"),
			LTypes([]string{files, PreDomns, PreHosts, urls}),
		)

		So(c.ReadCfg(&CFGstatic{Cfg: tdata.ZeroHostSourcesCfg}), ShouldBeNil)
		So(c.String(), ShouldEqual, tdata.JSONcfgZeroHostSources)
	})
}

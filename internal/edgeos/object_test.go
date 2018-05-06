package edgeos

import (
	"sort"
	"testing"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/smartystreets/goconvey/convey"
)

func TestObjectsAddObj(t *testing.T) {
	Convey("Testing ObjectsAddObj()", t, func() {
		c := NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
		)

		So(c.Blacklist(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)

		o, err := c.NewContent(FileObj)
		So(err, ShouldBeNil)

		exp := o

		o.GetList().addObj(c, rootNode)

		So(o, ShouldResemble, exp)
		// tests := []struct {
		// 	name string
		// 	rEnv *Env
		// 	rx     []*source
		// 	c    *Config
		// 	node string
		// }{
		// // TODO: Add test cases.
		// }
		// for _, tt := range tests {
		// 	o := &Objects{
		// 		Env: tt.rEnv,
		// 	src:     tt.rx,
		// 	}
		// 	o.addObj(tt.c, tt.node)
		// }
	})
}

func TestObjectString(t *testing.T) {
	Convey("Testing ObjectString()", t, func() {
		c := NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
		)

		So(c.Blacklist(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)

		act := c.GetAll()
		So(act.Find("sysctl.org"), ShouldEqual, 9)
		So(act.Find("@#$%"), ShouldEqual, -1)
	})
}

func TestSortObject(t *testing.T) {
	Convey("Testing SortObject()", t, func() {
		act := &Objects{
			src: []*source{
				{name: "eagle"},
				{name: "aardvark"},
				{name: "dog"},
				{name: "crab"},
				{name: "beetle"},
			},
		}

		exp := &Objects{
			src: []*source{
				{name: "aardvark"},
				{name: "beetle"},
				{name: "crab"},
				{name: "dog"},
				{name: "eagle"},
			},
		}

		sort.Sort(act)
		So(act, ShouldResemble, exp)
	})
}

func TestFilter(t *testing.T) {
	Convey("Testing SortObject()", t, func() {
		tests := []struct {
			ltype string
			exp   sort.StringSlice
		}{
			{ltype: urls, exp: urlsOnly},
			{ltype: files, exp: filesOnly},
			{ltype: hosts, exp: sort.StringSlice(nil)},
		}

		c := NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
		)

		So(c.Blacklist(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)

		for _, tt := range tests {
			Convey("Testing "+tt.ltype, func() {
				act := c.GetAll().Filter(tt.ltype)
				So(act.Names(), ShouldResemble, tt.exp)
			})
		}
	})
}

func TestGetLtypeDesc(t *testing.T) {
	Convey("Testing getLtypeDesc()", t, func() {
		So(getLtypeDesc(""), ShouldEqual, "pre-configured unknown ltype")
		So(getLtypeDesc("Hyperbolic-frisbee-throwers"), ShouldEqual, "pre-configured Hyperbolic frisbee throwers")
	})
}

var (
	filesOnly = sort.StringSlice{"tasty"}
	urlsOnly  = sort.StringSlice{"malc0de", "malwaredomains.com", "openphish", "raw.github.com", "simple_tracking", "sysctl.org", "volkerschatz", "yoyo", "zeus"}
)

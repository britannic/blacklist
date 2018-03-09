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
			LTypes([]string{files, PreDomns, PreHosts, urls}),
		)

		So(c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)

		o, err := c.NewContent(FileObj)
		So(err, ShouldBeNil)

		exp := o

		o.GetList().addObj(c, rootNode)

		So(o, ShouldResemble, exp)
		// tests := []struct {
		// 	name string
		// 	rParms *Parms
		// 	rx     []*object
		// 	c    *Config
		// 	node string
		// }{
		// // TODO: Add test cases.
		// }
		// for _, tt := range tests {
		// 	o := &Objects{
		// 		Parms: tt.rParms,
		// 		x:     tt.rx,
		// 	}
		// 	o.addObj(tt.c, tt.node)
		// }
	})
}

func TestObjectString(t *testing.T) {
	Convey("Testing ObjectString()", t, func() {
		exp := "[\nDesc:\t \"pre-configured blacklisted domains\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"blacklisted-subdomains\"\nName:\t \"blacklisted-subdomains\"\nnType:\t \"preDomn\"\nPrefix:\t \"\"\nType:\t \"blacklisted-subdomains\"\nURL:\t \"\"\n \nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"malc0de\"\nnType:\t \"domn\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n \nDesc:\t \"Just domains\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"10.0.0.1\"\nLtype:\t \"url\"\nName:\t \"malwaredomains.com\"\nnType:\t \"domn\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"http://mirror1.malwaredomains.com/files/justdomains\"\n \nDesc:\t \"Basic tracking list by Disconnect\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"url\"\nName:\t \"simple_tracking\"\nnType:\t \"domn\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt\"\n \nDesc:\t \"abuse.ch ZeuS domain blocklist\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"url\"\nName:\t \"zeus\"\nnType:\t \"domn\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist\"\n \nDesc:\t \"pre-configured blacklisted hosts\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"blacklisted-servers\"\nName:\t \"blacklisted-servers\"\nnType:\t \"preHost\"\nPrefix:\t \"\"\nType:\t \"blacklisted-servers\"\nURL:\t \"\"\n \nDesc:\t \"OpenPhish automatic phishing detection\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"openphish\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"https://openphish.com/feed.txt\"\n \nDesc:\t \"This hosts file is a merged collection of hosts from reputable sources\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"raw.github.com\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0 \"\nType:\t \"hosts\"\nURL:\t \"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts\"\n \nDesc:\t \"This hosts file is a merged collection of hosts from cameleon\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"172.16.16.1\"\nLtype:\t \"url\"\nName:\t \"sysctl.org\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1\\t \"\nType:\t \"hosts\"\nURL:\t \"http://sysctl.org/cameleon/hosts\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../internal/testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"Ad server blacklists\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"volkerschatz\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"http://www.volkerschatz.com/net/adpaths\"\n \nDesc:\t \"Fully Qualified Domain Names only - no prefix to strip\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"yoyo\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"https://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n]"

		c := NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
			LTypes([]string{files, PreDomns, PreHosts, urls}),
		)

		So(c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)

		act := c.GetAll()
		So(act.String(), ShouldEqual, exp)
		So(act.Find("sysctl.org"), ShouldEqual, 8)
		So(act.Find("@#$%"), ShouldEqual, -1)
	})
}

func TestSortObject(t *testing.T) {
	Convey("Testing SortObject()", t, func() {
		act := &Objects{
			x: []*object{
				{name: "eagle"},
				{name: "aardvark"},
				{name: "dog"},
				{name: "crab"},
				{name: "beetle"},
			},
		}

		exp := &Objects{
			x: []*object{
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
			{ltype: "-" + urls, exp: urlsNone},
			{ltype: files, exp: filesOnly},
			{ltype: "-" + files, exp: filesNone},
			{ltype: "zones", exp: nil},
		}

		c := NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
			LTypes([]string{files, PreDomns, PreHosts, urls}),
		)

		So(c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)

		for _, tt := range tests {
			Convey("Testing "+tt.ltype, func() {
				act := c.GetAll().Filter(tt.ltype)
				So(act.Names(), ShouldResemble, tt.exp)
			})
		}
	})
}

var (
	filesNone = sort.StringSlice{"blacklisted-servers", "blacklisted-subdomains", "malc0de", "malwaredomains.com", "openphish", "raw.github.com", "simple_tracking", "sysctl.org", "volkerschatz", "yoyo", "zeus"}
	filesOnly = sort.StringSlice{"tasty"}
	urlsNone  = sort.StringSlice{"blacklisted-servers", "blacklisted-subdomains", "tasty"}
	urlsOnly  = sort.StringSlice{"malc0de", "malwaredomains.com", "openphish", "raw.github.com", "simple_tracking", "sysctl.org", "volkerschatz", "yoyo", "zeus"}
)

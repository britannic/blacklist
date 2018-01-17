package edgeos

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAddInc(t *testing.T) {
	Convey("Testing addInc()", t, func() {
		c := NewConfig(Nodes([]string{rootNode, domains, hosts}))
		err := c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg})
		So(err, ShouldBeNil)

		tests := []struct {
			name string
			exp  *object
			node string
		}{
			{name: "rootNode", node: rootNode, exp: nil},
			{
				name: "domains",
				node: domains,
				exp: &object{
					Parms: &Parms{
						Wildcard: Wildcard{
							Node: "",
							Name: "",
						},
						API:   "",
						Arch:  "",
						Bash:  "",
						Cores: 0,
						Dbug:  false,
						Dex: list{
							RWMutex: &sync.RWMutex{},
							entry:   entry{},
						},
						Dir:    "",
						DNSsvc: "",
						Exc: list{
							RWMutex: &sync.RWMutex{},
							entry:   entry{},
						},
						Ext:      "",
						File:     "",
						FnFmt:    "",
						InCLI:    "",
						ioWriter: nil,
						Level:    "",
						Ltypes:   nil,
						Method:   "",
						Nodes:    []string{"blacklist", "domains", "hosts"},
						Pfx:      "",
						Test:     false,
						Timeout:  time.Duration(0),
						Verb:     false},
					desc:     "pre-configured-domain blacklist content",
					disabled: false,
					err:      nil,
					exc:      nil,
					file:     "",
					inc:      []string{"adsrvr.org", "adtechus.net", "advertising.com", "centade.com", "doubleclick.net", "free-counter.co.uk", "intellitxt.com", "kiosked.com", "patoghee.in"},
					ip:       "192.168.100.1",
					ltype:    "pre-configured-domain",
					name:     "includes.[9]",
					nType:    ntype(6),
					Objects: Objects{
						Parms: nil,
						x:     nil,
					},
					prefix: "",
					r:      nil,
					url:    "",
				},
			},
			{name: "hosts",
				node: hosts,
				exp: &object{
					Parms: &Parms{
						Wildcard: Wildcard{
							Node: "",
							Name: "",
						},
						API:   "",
						Arch:  "",
						Bash:  "",
						Cores: 0,
						Dbug:  false,
						Dex: list{
							RWMutex: &sync.RWMutex{},
							entry:   entry{},
						},
						Dir:    "",
						DNSsvc: "",
						Exc: list{
							RWMutex: &sync.RWMutex{},
							entry:   entry{},
						},
						Ext:      "",
						File:     "",
						FnFmt:    "",
						InCLI:    "",
						ioWriter: nil,
						Level:    "",
						Ltypes:   nil,
						Method:   "",
						Nodes:    []string{"blacklist", "domains", "hosts"},
						Pfx:      "",
						Test:     false,
						Timeout:  time.Duration(0),
						Verb:     false},
					desc:     "pre-configured-host blacklist content",
					disabled: false,
					err:      nil,
					exc:      nil,
					file:     "",
					inc:      []string{"beap.gemini.yahoo.com"},
					ip:       "0.0.0.0",
					ltype:    "pre-configured-host",
					name:     "includes.[1]",
					nType:    ntype(7),
					Objects: Objects{
						Parms: nil,
						x:     nil,
					},
					prefix: "",
					r:      nil,
					url:    "",
				},
			},
		}

		for _, tt := range tests {
			Convey("Testing "+tt.name, func() {
				Printf("%#v", c.addInc(tt.node))
				So(c.addInc(tt.node), ShouldResemble, tt.exp)
			})
		}
	})
}

func TestExcludes(t *testing.T) {
	Convey("Testing excludes()", t, func() {
		c := NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
			Nodes([]string{domains, hosts}),
		)

		err := c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg})
		So(err, ShouldBeNil)

		excludes := list{
			entry: entry{"sstatic.net": 0, "yimg.com": 0, "ytimg.com": 0, "google.com": 0, "images-amazon.com": 0, "msdn.com": 0, "schema.org": 0, "skype.com": 0, "avast.com": 0, "bitdefender.com": 0, "cdn.visiblemeasures.com": 0, "cloudfront.net": 0, "microsoft.com": 0, "akamaihd.net": 0, "amazon.com": 0, "apple.com": 0, "shopify.com": 0, "storage.googleapis.com": 0, "msecnd.net": 0, "ssl-on9.com": 0, "windows.net": 0, "1e100.net": 0, "akamai.net": 0, "coremetrics.com": 0, "gstatic.com": 0, "gvt1.com": 0, "freedns.afraid.org": 0, "hb.disney.go.com": 0, "hp.com": 0, "live.com": 0, "rackcdn.com": 0, "edgesuite.net": 0, "googleapis.com": 0, "smacargo.com": 0, "static.chartbeat.com": 0, "gvt1.net": 0, "hulu.com": 0, "paypal.com": 0, "amazonaws.com": 0, "ask.com": 0, "github.com": 0, "githubusercontent.com": 0, "googletagmanager.com": 0, "sourceforge.net": 0, "xboxlive.com": 0, "2o7.net": 0, "adobedtm.com": 0, "googleadservices.com": 0, "googleusercontent.com": 0, "ssl-on9.net": 0},
		}
		tests := []struct {
			get  list
			list list
			// raw  []string
			name string
			node string
		}{
			{name: "c.excludes(rootNode)", get: c.excludes(rootNode), list: excludes, node: rootNode},
			{name: "c.excludes()", get: c.excludes(), list: excludes},
			{name: "c.excludes(domains)", get: c.excludes(domains), list: list{RWMutex: (*sync.RWMutex)(nil), entry: entry{}}, node: domains},
			{name: "c.excludes(hosts)", get: c.excludes(hosts), list: list{RWMutex: (*sync.RWMutex)(nil), entry: entry{}}, node: hosts},
		}

		for _, tt := range tests {
			Convey("Testing "+tt.name, func() {
				switch tt.node {
				case "":
					So(c.excludes(), ShouldResemble, tt.list)
				default:
					So(c.excludes(tt.node), ShouldResemble, tt.list)
				}
			})
		}
	})
}

func TestFiles(t *testing.T) {
	Convey("Testing c.GetAll().Files()", t, func() {
		r := &CFGstatic{Cfg: tdata.Cfg}
		c := NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
			Nodes([]string{domains, hosts}),
			LTypes([]string{files, PreDomns, PreHosts, urls}),
		)

		So(c.ReadCfg(r), ShouldBeNil)

		exp := "/tmp/domains.malc0de.blacklist.conf\n/tmp/domains.malwaredomains.com.blacklist.conf\n/tmp/domains.simple_tracking.blacklist.conf\n/tmp/domains.zeus.blacklist.conf\n/tmp/hosts.openphish.blacklist.conf\n/tmp/hosts.raw.github.com.blacklist.conf\n/tmp/hosts.sysctl.org.blacklist.conf\n/tmp/hosts.tasty.blacklist.conf\n/tmp/hosts.volkerschatz.blacklist.conf\n/tmp/hosts.yoyo.blacklist.conf\n/tmp/pre-configured-domain.includes.[9].blacklist.conf\n/tmp/pre-configured-host.includes.[1].blacklist.conf"

		So(c.GetAll().Files().String(), ShouldEqual, exp)
	})
}

func TestInSession(t *testing.T) {
	Convey("Testing InSession()", t, func() {
		c := NewConfig()
		So(c.InSession(), ShouldBeFalse)

		So(os.Setenv("_OFR_CONFIGURE", "ok"), ShouldBeNil)
		So(c.InSession(), ShouldBeTrue)

		So(os.Unsetenv("_OFR_CONFIGURE"), ShouldBeNil)
		So(c.InSession(), ShouldBeFalse)
	})
}

func TestNodes(t *testing.T) {
	Convey("Testing Nodes()", t, func() {
		c := NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
			Nodes([]string{rootNode, domains, hosts}),
			LTypes([]string{files, PreDomns, PreHosts, urls}),
		)

		So(c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)
		So(c.Nodes(), ShouldResemble, []string{"blacklist", "domains", "hosts"})
	})
}

func TestReadCfg(t *testing.T) {
	Convey("Testing ReadCfg()", t, func() {
		Convey("Testing with an empty configuration", func() {
			exp := errors.New("Configuration data is empty, cannot continue")
			act := NewConfig().ReadCfg(&CFGstatic{Cfg: ""})
			So(act, ShouldResemble, exp)
		})
		Convey("Testing testing with a disabled configuration", func() {
			act := NewConfig().ReadCfg(&CFGstatic{Cfg: tdata.DisabledCfg})
			So(act, ShouldBeEmpty)
		})
	})
}

func TestReadUnconfiguredCfg(t *testing.T) {
	Convey("Testing ReadCfg()", t, func() {
		exp := errors.New("Configuration data is empty, cannot continue")
		act := NewConfig().ReadCfg(&CFGstatic{Cfg: tdata.NoBlacklist})
		So(act, ShouldResemble, exp)
	})
}

func TestReloadDNS(t *testing.T) {
	Convey("Testing ReloadDNS()", t, func() {
		act, err := NewConfig(Bash("/bin/bash"), DNSsvc("true")).ReloadDNS()
		So(err, ShouldBeNil)
		So(string(act), ShouldEqual, "")
	})
}

func TestRemove(t *testing.T) {
	Convey("Testing c.GetAll().Files().Remove()", t, func() {
		dir, _ := ioutil.TempDir("/tmp", "testBlacklist")
		defer os.RemoveAll(dir)

		var (
			c = NewConfig(
				Dir(dir),
				Ext("blacklist.conf"),
				FileNameFmt("%v/%v.%v.%v"),
				Nodes([]string{domains, hosts}),
				LTypes([]string{files, PreDomns, PreHosts, urls}),
				WCard(Wildcard{Node: "*s", Name: "*"}),
			)
			exp []string
		)

		So(c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)

		Convey("Creating special case file", func() {
			f, err := os.Create(fmt.Sprintf("%v/hosts.raw.github.com.blacklist.conf", dir))
			So(err, ShouldBeNil)
			f.Close()
		})

		for _, node := range c.Nodes() {
			for i := range Iter(10) {
				fname := fmt.Sprintf("%v/%v.%v.%v", dir, node, i, c.Ext)
				f, err := os.Create(fname)
				So(err, ShouldBeNil)
				f.Close()
			}
		}

		exp = append(exp, c.GetAll().Files().Strings()...)

		for _, fname := range exp {
			f, err := os.Create(fname)
			So(err, ShouldBeNil)
			f.Close()
		}

		c.GetAll().Files().Remove()

		cf := &CFile{Parms: c.Parms}
		pattern := fmt.Sprintf(c.FnFmt, c.Dir, "*s", "*", c.Parms.Ext)
		act, err := cf.readDir(pattern)
		So(err, ShouldBeNil)
		So(act, ShouldResemble, exp)

		prev := c.SetOpt(WCard(Wildcard{Node: "[]a]", Name: "]"}))

		So(cf.Remove(), ShouldNotBeNil)
		c.SetOpt(prev)
	})
}

func TestLTypes(t *testing.T) {
	Convey("Testing LTypes()", t, func() {
		exp := []string{files, PreDomns, PreHosts, urls}
		c := NewConfig(Dir("/tmp"),
			Ext("blacklist.conf"),
			Nodes([]string{rootNode, domains, hosts}),
			LTypes(exp),
		)
		So(c.LTypes(), ShouldResemble, exp)
	})
}

func TestBooltoString(t *testing.T) {
	Convey("Testing booltoString()", t, func() {
		So(booltoStr(true), ShouldEqual, True)
		So(booltoStr(false), ShouldEqual, False)
	})
}

func TestToBool(t *testing.T) {
	Convey("Testing strToBool()", t, func() {
		So(strToBool(True), ShouldBeTrue)
		So(strToBool(False), ShouldBeFalse)
	})
}

func TestGetAll(t *testing.T) {
	Convey("Testing GetAll() objects", t, func() {

		c := NewConfig(
			Dir("/tmp"),
			Ext(".blacklist.conf"),
			Nodes([]string{domains, hosts}),
			LTypes([]string{files, PreDomns, PreHosts, urls}),
		)

		So(c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)

		tests := []struct {
			name  string
			ltype string
			exp   string
		}{
			{name: "GetAll().x", ltype: "", exp: expGetAll},
			{name: "GetAll(url).x", ltype: urls, exp: expURLS},
			{name: "GetAll(files).x", ltype: files, exp: expFiles},
			{name: "GetAll(PreDomns, PreHosts).x", ltype: PreDomns, exp: expPre},
			{name: "GetAll().String()", ltype: all, exp: c.Get(all).String()},
			{name: "c.Get(hosts).String()", ltype: hosts, exp: expHostObj},
		}

		for _, tt := range tests {
			Convey("Testing "+tt.name, func() {
				switch tt.ltype {
				case "":
					So(fmt.Sprint(c.GetAll().x), ShouldEqual, tt.exp)
				case all:
					So(c.GetAll().String(), ShouldEqual, tt.exp)
				case hosts:
					So(c.Get(hosts).String(), ShouldEqual, tt.exp)
				case PreDomns:
					So(fmt.Sprint(c.GetAll(PreDomns, PreHosts).x), ShouldEqual, tt.exp)
				default:
					So(fmt.Sprint(c.GetAll(tt.ltype).x), ShouldResemble, tt.exp)
				}
			})
		}

	})
}

var (
	expFiles = "[\nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../internal/testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n]"

	expGetAll = "[\nDesc:\t \"pre-configured-domain blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"pre-configured-domain\"\nName:\t \"includes.[9]\"\nnType:\t \"preDomn\"\nPrefix:\t \"\"\nType:\t \"pre-configured-domain\"\nURL:\t \"\"\n \nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"malc0de\"\nnType:\t \"domn\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n \nDesc:\t \"Just domains\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"10.0.0.1\"\nLtype:\t \"url\"\nName:\t \"malwaredomains.com\"\nnType:\t \"domn\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"http://mirror1.malwaredomains.com/files/justdomains\"\n \nDesc:\t \"Basic tracking list by Disconnect\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"url\"\nName:\t \"simple_tracking\"\nnType:\t \"domn\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt\"\n \nDesc:\t \"abuse.ch ZeuS domain blocklist\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"url\"\nName:\t \"zeus\"\nnType:\t \"domn\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist\"\n \nDesc:\t \"pre-configured-host blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured-host\"\nName:\t \"includes.[1]\"\nnType:\t \"preHost\"\nPrefix:\t \"\"\nType:\t \"pre-configured-host\"\nURL:\t \"\"\n \nDesc:\t \"OpenPhish automatic phishing detection\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"openphish\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"https://openphish.com/feed.txt\"\n \nDesc:\t \"This hosts file is a merged collection of hosts from reputable sources\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"raw.github.com\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0 \"\nType:\t \"hosts\"\nURL:\t \"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts\"\n \nDesc:\t \"This hosts file is a merged collection of hosts from cameleon\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"172.16.16.1\"\nLtype:\t \"url\"\nName:\t \"sysctl.org\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1\\t \"\nType:\t \"hosts\"\nURL:\t \"http://sysctl.org/cameleon/hosts\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../internal/testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"Ad server blacklists\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"volkerschatz\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"http://www.volkerschatz.com/net/adpaths\"\n \nDesc:\t \"Fully Qualified Domain Names only - no prefix to strip\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"yoyo\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n]"

	expHostObj = "[\nDesc:\t \"pre-configured-host blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured-host\"\nName:\t \"includes.[1]\"\nnType:\t \"preHost\"\nPrefix:\t \"\"\nType:\t \"pre-configured-host\"\nURL:\t \"\"\n \nDesc:\t \"OpenPhish automatic phishing detection\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"openphish\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"https://openphish.com/feed.txt\"\n \nDesc:\t \"This hosts file is a merged collection of hosts from reputable sources\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"raw.github.com\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0 \"\nType:\t \"hosts\"\nURL:\t \"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts\"\n \nDesc:\t \"This hosts file is a merged collection of hosts from cameleon\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"172.16.16.1\"\nLtype:\t \"url\"\nName:\t \"sysctl.org\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1\\t \"\nType:\t \"hosts\"\nURL:\t \"http://sysctl.org/cameleon/hosts\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../internal/testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"Ad server blacklists\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"volkerschatz\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"http://www.volkerschatz.com/net/adpaths\"\n \nDesc:\t \"Fully Qualified Domain Names only - no prefix to strip\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"yoyo\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n]"

	expPre = "[\nDesc:\t \"pre-configured-domain blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"pre-configured-domain\"\nName:\t \"includes.[9]\"\nnType:\t \"preDomn\"\nPrefix:\t \"\"\nType:\t \"pre-configured-domain\"\nURL:\t \"\"\n \nDesc:\t \"pre-configured-host blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured-host\"\nName:\t \"includes.[1]\"\nnType:\t \"preHost\"\nPrefix:\t \"\"\nType:\t \"pre-configured-host\"\nURL:\t \"\"\n]"

	expURLS = "[\nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"malc0de\"\nnType:\t \"domn\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n \nDesc:\t \"Just domains\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"10.0.0.1\"\nLtype:\t \"url\"\nName:\t \"malwaredomains.com\"\nnType:\t \"domn\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"http://mirror1.malwaredomains.com/files/justdomains\"\n \nDesc:\t \"Basic tracking list by Disconnect\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"url\"\nName:\t \"simple_tracking\"\nnType:\t \"domn\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt\"\n \nDesc:\t \"abuse.ch ZeuS domain blocklist\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"url\"\nName:\t \"zeus\"\nnType:\t \"domn\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist\"\n \nDesc:\t \"OpenPhish automatic phishing detection\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"openphish\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"https://openphish.com/feed.txt\"\n \nDesc:\t \"This hosts file is a merged collection of hosts from reputable sources\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"raw.github.com\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0 \"\nType:\t \"hosts\"\nURL:\t \"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts\"\n \nDesc:\t \"This hosts file is a merged collection of hosts from cameleon\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"172.16.16.1\"\nLtype:\t \"url\"\nName:\t \"sysctl.org\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1\\t \"\nType:\t \"hosts\"\nURL:\t \"http://sysctl.org/cameleon/hosts\"\n \nDesc:\t \"Ad server blacklists\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"volkerschatz\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"http://www.volkerschatz.com/net/adpaths\"\n \nDesc:\t \"Fully Qualified Domain Names only - no prefix to strip\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"yoyo\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n]"
)

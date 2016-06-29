package edgeos

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/britannic/testutils"
	"github.com/pmezard/go-difflib/difflib"
)

func uDiff(a, b string) string {
	diff := difflib.ContextDiff{
		A:        difflib.SplitLines(a),
		B:        difflib.SplitLines(b),
		FromFile: "Want",
		ToFile:   "Got",
		Context:  3,
		Eol:      "\n",
	}

	result, _ := difflib.GetContextDiffString(diff)
	return fmt.Sprintf(strings.Replace(result, "\t", " ", -1))
}

func TestAddInc(t *testing.T) {
	c := NewConfig(
		Nodes([]string{rootNode, domains, hosts}),
	)
	err := c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg})
	OK(t, err)

	tests := []struct {
		exp  string
		node string
	}{
		{node: rootNode, exp: "<nil>"},
		{node: domains, exp: "\nDesc:\t \"pre-configured-domain blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"pre-configured-domain\"\nName:\t \"includes.[9]\"\nnType:\t \"preDomn\"\nPrefix:\t \"\"\nType:\t \"pre-configured-domain\"\nURL:\t \"\"\n"},
		{node: hosts, exp: "\nDesc:\t \"pre-configured-host blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured-host\"\nName:\t \"includes.[1]\"\nnType:\t \"preHost\"\nPrefix:\t \"\"\nType:\t \"pre-configured-host\"\nURL:\t \"\"\n"},
	}

	for _, tt := range tests {
		Equals(t, tt.exp, fmt.Sprint(c.addInc(tt.node)))
	}
}

func TestExcludes(t *testing.T) {
	c := NewConfig(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Nodes([]string{domains, hosts}),
	)

	err := c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg})
	OK(t, err)

	excludes := List{"sstatic.net": 0, "yimg.com": 0, "ytimg.com": 0, "google.com": 0, "images-amazon.com": 0, "msdn.com": 0, "schema.org": 0, "skype.com": 0, "avast.com": 0, "bitdefender.com": 0, "cdn.visiblemeasures.com": 0, "cloudfront.net": 0, "microsoft.com": 0, "akamaihd.net": 0, "amazon.com": 0, "apple.com": 0, "shopify.com": 0, "storage.googleapis.com": 0, "msecnd.net": 0, "ssl-on9.com": 0, "windows.net": 0, "1e100.net": 0, "akamai.net": 0, "coremetrics.com": 0, "gstatic.com": 0, "gvt1.com": 0, "freedns.afraid.org": 0, "hb.disney.go.com": 0, "hp.com": 0, "live.com": 0, "rackcdn.com": 0, "edgesuite.net": 0, "googleapis.com": 0, "smacargo.com": 0, "static.chartbeat.com": 0, "gvt1.net": 0, "hulu.com": 0, "paypal.com": 0, "amazonaws.com": 0, "ask.com": 0, "github.com": 0, "githubusercontent.com": 0, "googletagmanager.com": 0, "sourceforge.net": 0, "xboxlive.com": 0, "2o7.net": 0, "adobedtm.com": 0, "googleadservices.com": 0, "googleusercontent.com": 0, "ssl-on9.net": 0}

	tests := []struct {
		get  List
		list List
		raw  []string
		node string
	}{
		{get: c.excludes(rootNode), list: excludes, node: rootNode},
		{get: c.excludes(), list: excludes},
		{get: c.excludes(domains), list: List{}, node: domains},
		{get: c.excludes(hosts), list: List{}, node: hosts},
	}

	for _, test := range tests {
		switch test.node {
		case "":
			Equals(t, test.list, c.excludes())
		default:
			Equals(t, test.list, c.excludes(test.node))
		}

	}
}

func TestFiles(t *testing.T) {
	r := &CFGstatic{Cfg: tdata.Cfg}
	c := NewConfig(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Nodes([]string{domains, hosts}),
		LTypes([]string{files, PreDomns, PreHosts, urls}),
	)
	err := c.ReadCfg(r)
	OK(t, err)

	want := "/tmp/domains.malc0de.blacklist.conf\n/tmp/domains.malwaredomains.com.blacklist.conf\n/tmp/domains.simple_tracking.blacklist.conf\n/tmp/domains.zeus.blacklist.conf\n/tmp/hosts.openphish.blacklist.conf\n/tmp/hosts.raw.github.com.blacklist.conf\n/tmp/hosts.sysctl.org.blacklist.conf\n/tmp/hosts.tasty.blacklist.conf\n/tmp/hosts.volkerschatz.blacklist.conf\n/tmp/hosts.yoyo.blacklist.conf\n/tmp/pre-configured-domain.includes.[9].blacklist.conf\n/tmp/pre-configured-host.includes.[1].blacklist.conf"

	got := c.GetAll().Files().String()
	Equals(t, want, got)
}

func TestInSession(t *testing.T) {
	c := NewConfig()
	Equals(t, false, c.InSession())
	err := os.Setenv("_OFR_CONFIGURE", "ok")
	OK(t, err)
	Equals(t, true, c.InSession())
	err = os.Unsetenv("_OFR_CONFIGURE")
	OK(t, err)
}

func TestNodes(t *testing.T) {
	c := NewConfig(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Nodes([]string{rootNode, domains, hosts}),
		LTypes([]string{files, PreDomns, PreHosts, urls}),
	)
	r := &CFGstatic{Cfg: tdata.Cfg}
	nodes := []string{"blacklist", "domains", "hosts"}

	err := c.ReadCfg(r)
	OK(t, err)

	Equals(t, nodes, c.Nodes())
}

func TestReadCfg(t *testing.T) {
	want := errors.New("Configuration data is empty, cannot continue")
	l := &CFGstatic{Cfg: ""}
	c := NewConfig()
	got := c.ReadCfg(l)
	Equals(t, want, got)
}

// type dummy struct{}
//
// func (d *dummy) readDir(s string) ([]os.FileInfo, error) {
// 	return []os.FileInfo{}, fmt.Errorf("%v totally failed!", s)
// }
//
// func (d *dummy) Remove() error {
// 	return nil
// }

func TestReloadDNS(t *testing.T) {
	c := NewConfig(DNSsvc("true"))
	act, err := c.ReloadDNS()
	OK(t, err)
	Equals(t, "", string(act))
}

func TestRemove(t *testing.T) {
	var (
		dir    = "../testdata"
		ext    = "blacklist.conf"
		nodes  = []string{rootNode, domains, hosts}
		r      = &CFGstatic{Cfg: tdata.Cfg}
		Ltypes = []string{files, PreDomns, PreHosts, urls}
		c      = NewConfig(
			Dir(dir),
			Ext(ext),
			FileNameFmt("%v/%v.%v.%v"),
			Nodes(nodes),
			LTypes(Ltypes),
			WCard(Wildcard{Node: "*s", Name: "*"}),
		)
	)

	err := c.ReadCfg(r)
	OK(t, err)

	// Special case file Name
	f, err := os.Create(fmt.Sprintf("%v/hosts.raw.github.com.blacklist.conf", dir))
	f.Close()

	for _, node := range nodes {
		for i := 1; i < 10; i++ {
			fname := fmt.Sprintf("%v/%v.%v.%v", dir, node, i, ext)
			f, err = os.Create(fname)
			OK(t, err)
			f.Close()
		}
	}

	var want []string

	want = append(want, c.GetAll().Files().Strings()...)

	for _, fname := range want {
		f, err = os.Create(fname)
		OK(t, err)
		f.Close()
	}

	c.GetAll().Files().Remove()

	cf := &CFile{Parms: c.Parms}
	pattern := fmt.Sprintf(c.FnFmt, c.Dir, "*s", "*", c.Parms.Ext)
	got, err := cf.readDir(pattern)
	OK(t, err)

	Equals(t, want, got)

	prev := c.SetOpt(WCard(Wildcard{Node: "[]a]", Name: "]"}))

	err = cf.Remove()
	NotOK(t, err)
	c.SetOpt(prev)
}

func TestLTypes(t *testing.T) {
	want := []string{files, PreDomns, PreHosts, urls}
	c := &Config{}
	c.Parms = NewParms()
	c.SetOpt(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Nodes([]string{rootNode, domains, hosts}),
		LTypes(want),
	)
	Equals(t, want, c.LTypes())
}

func TestBoolToString(t *testing.T) {
	Equals(t, True, BooltoStr(true))
	Equals(t, False, BooltoStr(false))
}

func TestToBool(t *testing.T) {
	Equals(t, true, StrToBool(True))
	Equals(t, false, StrToBool(False))
}

func TestGetAll(t *testing.T) {
	want := "[\nDesc:\t \"pre-configured-domain blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"pre-configured-domain\"\nName:\t \"includes.[9]\"\nnType:\t \"preDomn\"\nPrefix:\t \"\"\nType:\t \"pre-configured-domain\"\nURL:\t \"\"\n \nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"malc0de\"\nnType:\t \"domain\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n \nDesc:\t \"Just domains\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"10.0.0.1\"\nLtype:\t \"url\"\nName:\t \"malwaredomains.com\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"http://mirror1.malwaredomains.com/files/justdomains\"\n \nDesc:\t \"Basic tracking list by Disconnect\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"url\"\nName:\t \"simple_tracking\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt\"\n \nDesc:\t \"abuse.ch ZeuS domain blocklist\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"url\"\nName:\t \"zeus\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist\"\n \nDesc:\t \"pre-configured-host blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured-host\"\nName:\t \"includes.[1]\"\nnType:\t \"preHost\"\nPrefix:\t \"\"\nType:\t \"pre-configured-host\"\nURL:\t \"\"\n \nDesc:\t \"OpenPhish automatic phishing detection\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"openphish\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"https://openphish.com/feed.txt\"\n \nDesc:\t \"This hosts file is a merged collection of hosts from reputable sources\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"raw.github.com\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0 \"\nType:\t \"hosts\"\nURL:\t \"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts\"\n \nDesc:\t \"This hosts file is a merged collection of hosts from cameleon\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"172.16.16.1\"\nLtype:\t \"url\"\nName:\t \"sysctl.org\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1\\t \"\nType:\t \"hosts\"\nURL:\t \"http://sysctl.org/cameleon/hosts\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"Ad server blacklists\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"volkerschatz\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"http://www.volkerschatz.com/net/adpaths\"\n \nDesc:\t \"Fully Qualified Domain Names only - no prefix to strip\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"yoyo\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n]"

	wantHostObj := "[\nDesc:\t \"pre-configured-host blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured-host\"\nName:\t \"includes.[1]\"\nnType:\t \"preHost\"\nPrefix:\t \"\"\nType:\t \"pre-configured-host\"\nURL:\t \"\"\n \nDesc:\t \"OpenPhish automatic phishing detection\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"openphish\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"https://openphish.com/feed.txt\"\n \nDesc:\t \"This hosts file is a merged collection of hosts from reputable sources\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"raw.github.com\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0 \"\nType:\t \"hosts\"\nURL:\t \"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts\"\n \nDesc:\t \"This hosts file is a merged collection of hosts from cameleon\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"172.16.16.1\"\nLtype:\t \"url\"\nName:\t \"sysctl.org\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1\\t \"\nType:\t \"hosts\"\nURL:\t \"http://sysctl.org/cameleon/hosts\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"Ad server blacklists\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"volkerschatz\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"http://www.volkerschatz.com/net/adpaths\"\n \nDesc:\t \"Fully Qualified Domain Names only - no prefix to strip\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"yoyo\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n]"

	wantURLS := "[\nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"malc0de\"\nnType:\t \"domain\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n \nDesc:\t \"Just domains\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"10.0.0.1\"\nLtype:\t \"url\"\nName:\t \"malwaredomains.com\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"http://mirror1.malwaredomains.com/files/justdomains\"\n \nDesc:\t \"Basic tracking list by Disconnect\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"url\"\nName:\t \"simple_tracking\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt\"\n \nDesc:\t \"abuse.ch ZeuS domain blocklist\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"url\"\nName:\t \"zeus\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist\"\n \nDesc:\t \"OpenPhish automatic phishing detection\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"openphish\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"https://openphish.com/feed.txt\"\n \nDesc:\t \"This hosts file is a merged collection of hosts from reputable sources\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"raw.github.com\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0 \"\nType:\t \"hosts\"\nURL:\t \"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts\"\n \nDesc:\t \"This hosts file is a merged collection of hosts from cameleon\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"172.16.16.1\"\nLtype:\t \"url\"\nName:\t \"sysctl.org\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1\\t \"\nType:\t \"hosts\"\nURL:\t \"http://sysctl.org/cameleon/hosts\"\n \nDesc:\t \"Ad server blacklists\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"volkerschatz\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"http://www.volkerschatz.com/net/adpaths\"\n \nDesc:\t \"Fully Qualified Domain Names only - no prefix to strip\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"yoyo\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n]"

	wantFiles := "[\nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n]"

	wantPre := "[\nDesc:\t \"pre-configured-domain blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.100.1\"\nLtype:\t \"pre-configured-domain\"\nName:\t \"includes.[9]\"\nnType:\t \"preDomn\"\nPrefix:\t \"\"\nType:\t \"pre-configured-domain\"\nURL:\t \"\"\n \nDesc:\t \"pre-configured-host blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured-host\"\nName:\t \"includes.[1]\"\nnType:\t \"preHost\"\nPrefix:\t \"\"\nType:\t \"pre-configured-host\"\nURL:\t \"\"\n]"

	c := NewConfig(
		Dir("/tmp"),
		Ext(".blacklist.conf"),
		Nodes([]string{domains, hosts}),
		LTypes([]string{files, PreDomns, PreHosts, urls}),
	)

	err := c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg})
	OK(t, err)

	Equals(t, want, fmt.Sprint(c.GetAll().S))
	Equals(t, wantURLS, fmt.Sprint(c.GetAll(urls).S))
	Equals(t, wantFiles, fmt.Sprint(c.GetAll(files).S))
	Equals(t, wantPre, fmt.Sprint(c.GetAll(PreDomns, PreHosts).S))
	Equals(t, c.Get(all).String(), c.GetAll().String())
	Equals(t, wantHostObj, c.Get(hosts).String())
}

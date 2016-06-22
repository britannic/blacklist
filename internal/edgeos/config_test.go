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

func TestExcludes(t *testing.T) {
	r := &CFGstatic{Cfg: tdata.Cfg}

	c := NewConfig(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Nodes([]string{"domains", "hosts"}),
	)

	err := c.ReadCfg(r)
	OK(t, err)

	excludes := List{"msdn.com": 0, "paypal.com": 0, "smacargo.com": 0, "122.2o7.net": 0, "freedns.afraid.org": 0, "githubusercontent.com": 0, "gvt1.com": 0, "hulu.com": 0, "sourceforge.net": 0, "coremetrics.com": 0, "google.com": 0, "googleadservices.com": 0, "gvt1.net": 0, "skype.com": 0, "adobedtm.com": 0, "cloudfront.net": 0, "gstatic.com": 0, "akamai.net": 0, "bitdefender.com": 0, "images-amazon.com": 0, "windows.net": 0, "ytimg.com": 0, "cdn.visiblemeasures.com": 0, "hp.com": 0, "static.chartbeat.com": 0, "amazonaws.com": 0, "apple.com": 0, "googleapis.com": 0, "storage.googleapis.com": 0, "yimg.com": 0, "1e100.net": 0, "schema.org": 0, "ssl-on9.net": 0, "googleusercontent.com": 0, "hb.disney.go.com": 0, "rackcdn.com": 0, "amazon.com": 0, "ask.com": 0, "avast.com": 0, "edgesuite.net": 0, "github.com": 0, "ssl-on9.com": 0}

	tests := []struct {
		get  []string
		list List
		raw  []string
		node string
	}{
		{get: c.excludes("blacklist"), list: excludes, node: "blacklist"},
		{get: c.excludes("all"), list: excludes, node: "all"},
		{get: c.excludes("domains"), list: List{}, node: "domains"},
		{get: c.excludes("hosts"), list: List{}, node: "hosts"},
	}

	for _, test := range tests {
		Equals(t, test.list, c.Excludes(test.node))
	}
}

func TestFiles(t *testing.T) {
	r := &CFGstatic{Cfg: tdata.Cfg}
	c := NewConfig(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Nodes([]string{"domains", "hosts"}),
		STypes([]string{"file", "pre-configured", "url"}),
	)
	err := c.ReadCfg(r)
	OK(t, err)

	want := "/tmp/domains.malc0de.blacklist.conf\n/tmp/domains.pre-configured.blacklist.conf\n/tmp/hosts.adaway.blacklist.conf\n/tmp/hosts.malwaredomainlist.blacklist.conf\n/tmp/hosts.openphish.blacklist.conf\n/tmp/hosts.pre-configured.blacklist.conf\n/tmp/hosts.someonewhocares.blacklist.conf\n/tmp/hosts.tasty.blacklist.conf\n/tmp/hosts.volkerschatz.blacklist.conf\n/tmp/hosts.winhelp2002.blacklist.conf\n/tmp/hosts.yoyo.blacklist.conf"

	got := c.GetAll().Files().String()
	Equals(t, want, got)
}

func TestNodes(t *testing.T) {
	c := NewConfig(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Nodes([]string{"domains", "hosts"}),
		STypes([]string{"file", "pre-configured", "url"}),
	)
	r := &CFGstatic{Cfg: tdata.Cfg}
	nodes := []string{"blacklist", "domains", "hosts"}

	err := c.ReadCfg(r)
	OK(t, err)

	Equals(t, nodes, c.Nodes())
	// fmt.Println(c.Nodes())
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
// func (d *dummy) ReadDir(s string) ([]os.FileInfo, error) {
// 	return []os.FileInfo{}, fmt.Errorf("%v totally failed!", s)
// }
//
// func (d *dummy) Remove() error {
// 	return nil
// }

func TestRemove(t *testing.T) {
	var (
		dir    = "../testdata"
		ext    = "blacklist.conf"
		nodes  = []string{"domains", "hosts"}
		r      = &CFGstatic{Cfg: tdata.Cfg}
		stypes = []string{"file", "pre-configured", "url"}
		c      = NewConfig(
			Dir(dir),
			Ext(ext),
			FileNameFmt("%v/%v.%v.%v"),
			Nodes(nodes),
			STypes(stypes),
		)
	)

	err := c.ReadCfg(r)
	OK(t, err)

	for i := 1; i < 10; i++ {
		for _, node := range nodes {
			fname := fmt.Sprintf("%v/%v.%v.%v", dir, node, i, ext)
			f, err := os.Create(fname)
			OK(t, err)
			f.Close()
		}
	}

	var want []string

	want = append(want, c.GetAll().Files().Strings()...)

	for _, fname := range want {
		f, err := os.Create(fname)
		OK(t, err)
		f.Close()
	}

	c.GetAll().Files().Remove()

	cf := &CFile{Parms: c.Parms}
	pattern := fmt.Sprintf(c.FnFmt, c.Dir, "*s", "*", c.Parms.Ext)
	got, err := cf.ReadDir(pattern)
	OK(t, err)

	Equals(t, want, got)
	// fmt.Println(DiffArray(want, got))
	prev := c.SetOpt(WCard(Wildcard{node: "[]a]", name: "]"}))
	// pattern = fmt.Sprintf(c.FnFmt, c.Dir, "[]a]", "]", c.Parms.Ext)
	// _, err = cf.ReadDir(pattern)
	// NotOK(t, err)

	err = cf.Remove()
	NotOK(t, err)
	c.SetOpt(prev)
}

func TestSource(t *testing.T) {
	c := NewConfig(
		Dir("/tmp"),
		Ext(".blacklist.conf"),
		Nodes([]string{"domains", "hosts"}),
		STypes([]string{"file", "pre-configured", "url"}),
	)

	r := &CFGstatic{Cfg: tdata.Cfg}
	err := c.ReadCfg(r)
	OK(t, err)

	var (
		wdomains = "\nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"\"\n\nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"malc0de\"\nnType:\t \"domain\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n"

		whosts = "\nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"\"\n\nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"malc0de\"\nnType:\t \"domain\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n\nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n\nDesc:\t \"Blocking mobile ad providers and some analytics providers\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"adaway\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://adaway.org/hosts.txt\"\n\nDesc:\t \"127.0.0.1 based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"malwaredomainlist\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://www.malwaredomainlist.com/hostslist/hosts.txt\"\n\nDesc:\t \"OpenPhish automatic phishing detection\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"openphish\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"https://openphish.com/feed.txt\"\n\nDesc:\t \"Zero based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"someonewhocares\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0\"\nType:\t \"hosts\"\nURL:\t \"http://someonewhocares.org/hosts/zero/\"\n\nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"/config/user-data/blist.hosts.src\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n\nDesc:\t \"Ad server blacklists\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"volkerschatz\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"http://www.volkerschatz.com/net/adpaths\"\n\nDesc:\t \"Zero based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"winhelp2002\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0 \"\nType:\t \"hosts\"\nURL:\t \"http://winhelp2002.mvps.org/hosts.txt\"\n\nDesc:\t \"Fully Qualified Domain Names only - no prefix to strip\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"yoyo\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n"
		got    string
		want   = map[string]string{
			"domains": wdomains,
			"hosts":   whosts,
		}
	)
	for _, node := range []string{"domains", "hosts"} {
		srcs := *c.Get(node).Source("all")
		for _, src := range srcs.S {
			got += src.String()
		}
		Equals(t, want[node], got)
		// fmt.Println(uDiff(want[node], got))
	}
}

func TestSTypes(t *testing.T) {
	want := []string{"file", "pre-configured", "url"}
	c := &Config{}
	c.Parms = NewParms()
	c.SetOpt(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Nodes([]string{"domains", "hosts"}),
		STypes(want),
	)
	Equals(t, want, c.STypes())
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
	r := &CFGstatic{Cfg: tdata.Cfg}
	want := "[\nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"\"\n \nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"malc0de\"\nnType:\t \"domain\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n \nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"Blocking mobile ad providers and some analytics providers\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"adaway\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://adaway.org/hosts.txt\"\n \nDesc:\t \"127.0.0.1 based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"malwaredomainlist\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://www.malwaredomainlist.com/hostslist/hosts.txt\"\n \nDesc:\t \"OpenPhish automatic phishing detection\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"openphish\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"https://openphish.com/feed.txt\"\n \nDesc:\t \"Zero based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"someonewhocares\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0\"\nType:\t \"hosts\"\nURL:\t \"http://someonewhocares.org/hosts/zero/\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"/config/user-data/blist.hosts.src\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"Ad server blacklists\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"volkerschatz\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"http://www.volkerschatz.com/net/adpaths\"\n \nDesc:\t \"Zero based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"winhelp2002\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0 \"\nType:\t \"hosts\"\nURL:\t \"http://winhelp2002.mvps.org/hosts.txt\"\n \nDesc:\t \"Fully Qualified Domain Names only - no prefix to strip\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"yoyo\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n]"

	wantURLS := "[\nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"malc0de\"\nnType:\t \"domain\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n \nDesc:\t \"Blocking mobile ad providers and some analytics providers\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"adaway\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://adaway.org/hosts.txt\"\n \nDesc:\t \"127.0.0.1 based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"malwaredomainlist\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://www.malwaredomainlist.com/hostslist/hosts.txt\"\n \nDesc:\t \"OpenPhish automatic phishing detection\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"openphish\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"https://openphish.com/feed.txt\"\n \nDesc:\t \"Zero based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"someonewhocares\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0\"\nType:\t \"hosts\"\nURL:\t \"http://someonewhocares.org/hosts/zero/\"\n \nDesc:\t \"Ad server blacklists\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"volkerschatz\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"http://www.volkerschatz.com/net/adpaths\"\n \nDesc:\t \"Zero based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"winhelp2002\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0 \"\nType:\t \"hosts\"\nURL:\t \"http://winhelp2002.mvps.org/hosts.txt\"\n \nDesc:\t \"Fully Qualified Domain Names only - no prefix to strip\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"yoyo\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n]"

	wantFiles := "[\nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"/config/user-data/blist.hosts.src\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n]"

	wantPre := "[\nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"\"\n \nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n]"

	c := NewConfig(
		Dir("/tmp"),
		Ext(".blacklist.conf"),
		Nodes([]string{"domains", "hosts"}),
		STypes([]string{"file", "pre-configured", "url"}),
	)

	err := c.ReadCfg(r)
	OK(t, err)

	Equals(t, want, fmt.Sprint(c.GetAll().S))
	Equals(t, wantURLS, fmt.Sprint(c.GetAll(urls).S))
	Equals(t, wantFiles, fmt.Sprint(c.GetAll(files).S))
	Equals(t, wantPre, fmt.Sprint(c.GetAll(preConf).S))
	Equals(t, c.Get(all).String(), c.GetAll().String())
}

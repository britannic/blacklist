package edgeos_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/britannic/blacklist/internal/edgeos"
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
	r := &edgeos.CFGstatic{Cfg: tdata.Cfg}
	c, err := edgeos.ReadCfg(r)
	OK(t, err)

	c.Parms = edgeos.NewParms()
	c.SetOpt(
		edgeos.Dir("/tmp"),
		edgeos.Ext("blacklist.conf"),
		edgeos.Nodes([]string{"domains", "hosts"}),
	)

	excludes := []string{"122.2o7.net", "1e100.net", "adobedtm.com", "akamai.net", "amazon.com", "amazonaws.com", "apple.com", "ask.com", "avast.com", "bitdefender.com", "cdn.visiblemeasures.com", "cloudfront.net", "coremetrics.com", "edgesuite.net", "freedns.afraid.org", "github.com", "githubusercontent.com", "google.com", "googleadservices.com", "googleapis.com", "googleusercontent.com", "gstatic.com", "gvt1.com", "gvt1.net", "hb.disney.go.com", "hp.com", "hulu.com", "images-amazon.com", "msdn.com", "paypal.com", "rackcdn.com", "schema.org", "skype.com", "smacargo.com", "sourceforge.net", "ssl-on9.com", "ssl-on9.net", "static.chartbeat.com", "storage.googleapis.com", "windows.net", "yimg.com", "ytimg.com"}

	tests := []struct {
		get  string
		raw  []string
		node string
	}{
		{get: "", raw: []string{"122.2o7.net", "1e100.net", "adobedtm.com", "akamai.net", "amazon.com", "amazonaws.com", "apple.com", "ask.com", "avast.com", "bitdefender.com", "cdn.visiblemeasures.com", "cloudfront.net", "coremetrics.com", "edgesuite.net", "freedns.afraid.org", "github.com", "githubusercontent.com", "google.com", "googleadservices.com", "googleapis.com", "googleusercontent.com", "gstatic.com", "gvt1.com", "gvt1.net", "hb.disney.go.com", "hp.com", "hulu.com", "images-amazon.com", "msdn.com", "paypal.com", "rackcdn.com", "schema.org", "skype.com", "smacargo.com", "sourceforge.net", "ssl-on9.com", "ssl-on9.net", "static.chartbeat.com", "storage.googleapis.com", "windows.net", "yimg.com", "ytimg.com"}, node: "all"},
		{get: edgeos.UpdateList(excludes).String(), raw: excludes, node: "blacklist"},
		{get: "", raw: []string{}, node: "domains"},
		{get: "", raw: []string{}, node: "hosts"},
	}

	for _, test := range tests {
		Equals(t, test.get, c.Get(test.node).Excludes().String())
		Equals(t, test.raw, c.Excludes(test.node))
	}
}

func TestFiles(t *testing.T) {
	r := &edgeos.CFGstatic{Cfg: tdata.Cfg}
	c, err := edgeos.ReadCfg(r)
	OK(t, err)
	c.Parms = edgeos.NewParms()
	c.SetOpt(
		edgeos.Dir("/tmp"),
		edgeos.Ext("blacklist.conf"),
		edgeos.Nodes([]string{"domains", "hosts"}),
		edgeos.STypes([]string{"files", "pre-configured", "urls"}),
	)

	want := "/tmp/domains.malc0de.blacklist.conf\n/tmp/domains.pre-configured.blacklist.conf\n/tmp/hosts.adaway.blacklist.conf\n/tmp/hosts.malwaredomainlist.blacklist.conf\n/tmp/hosts.openphish.blacklist.conf\n/tmp/hosts.pre-configured.blacklist.conf\n/tmp/hosts.someonewhocares.blacklist.conf\n/tmp/hosts.tasty.blacklist.conf\n/tmp/hosts.volkerschatz.blacklist.conf\n/tmp/hosts.winhelp2002.blacklist.conf\n/tmp/hosts.yoyo.blacklist.conf"

	got := c.GetAll().Files().String()
	Equals(t, want, got)
}

func TestNodes(t *testing.T) {
	r := &edgeos.CFGstatic{Cfg: tdata.Cfg}
	c, err := edgeos.ReadCfg(r)
	OK(t, err)
	c.Parms = edgeos.NewParms()
	c.SetOpt(
		edgeos.Dir("/tmp"),
		edgeos.Ext("blacklist.conf"),
		edgeos.Nodes([]string{"domains", "hosts"}),
		edgeos.STypes([]string{"files", "pre-configured", "urls"}),
	)

	nodes := []string{"blacklist", "domains", "hosts"}
	Equals(t, nodes, c.Nodes())
	// fmt.Println(c.Nodes())
}

func TestReadCfg(t *testing.T) {
	want := errors.New("Configuration data is empty, cannot continue")
	l := &edgeos.CFGstatic{Cfg: ""}
	_, got := edgeos.ReadCfg(l)
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
		stypes = []string{"files", "pre-configured", "urls"}
	)
	r := &edgeos.CFGstatic{Cfg: tdata.Cfg}
	c, err := edgeos.ReadCfg(r)
	OK(t, err)
	c.Parms = edgeos.NewParms()
	c.SetOpt(
		edgeos.Dir(dir),
		edgeos.Ext(ext),
		edgeos.Nodes(nodes),
		edgeos.STypes(stypes),
	)

	for i := 1; i < 10; i++ {
		for _, node := range nodes {
			os.Open(fmt.Sprintf("%v.%v.%v", node, i, ext))
		}
	}

	var want []string
	for _, node := range nodes {
		want = append(want, c.Get(node).Source("all").Files().Strings()...)
		c.Get(node).Source("all").Files().Remove()
	}

	cf := &edgeos.CFile{}
	dlist, err := cf.ReadDir(dir)
	OK(t, err)

	var got []string
	for _, f := range dlist {
		if strings.Contains(f.Name(), ext) {
			got = append(got, dir+"/"+f.Name())
		}
	}

	Equals(t, want, got)
	// fmt.Println(edgeos.DiffArray(want, got))

	c.SetOpt(edgeos.Dir("re-#-Zzzzz"))
	err = c.Get("domains").Source("all").Files().Remove()
	NotOK(t, err)
}

func TestSource(t *testing.T) {
	r := &edgeos.CFGstatic{Cfg: tdata.Cfg}
	c, err := edgeos.ReadCfg(r)
	OK(t, err)
	c.Parms = edgeos.NewParms()
	c.Parms.SetOpt(
		edgeos.Dir("/tmp"),
		edgeos.Ext(".blacklist.conf"),
		edgeos.Nodes([]string{"domains", "hosts"}),
		edgeos.STypes([]string{"files", "pre-configured", "urls"}),
	)

	var (
		wdomains = "\nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"\"\n\nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"urls\"\nName:\t \"malc0de\"\nnType:\t \"domain\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n"

		whosts = "\nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"\"\n\nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"urls\"\nName:\t \"malc0de\"\nnType:\t \"domain\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n\nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n\nDesc:\t \"Blocking mobile ad providers and some analytics providers\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"adaway\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://adaway.org/hosts.txt\"\n\nDesc:\t \"127.0.0.1 based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"malwaredomainlist\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://www.malwaredomainlist.com/hostslist/hosts.txt\"\n\nDesc:\t \"OpenPhish automatic phishing detection\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"openphish\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"https://openphish.com/feed.txt\"\n\nDesc:\t \"Zero based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"someonewhocares\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0\"\nType:\t \"hosts\"\nURL:\t \"http://someonewhocares.org/hosts/zero/\"\n\nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"/config/user-data/blist.hosts.src\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"files\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n\nDesc:\t \"Ad server blacklists\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"volkerschatz\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"http://www.volkerschatz.com/net/adpaths\"\n\nDesc:\t \"Zero based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"urls\"\nName:\t \"winhelp2002\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0 \"\nType:\t \"hosts\"\nURL:\t \"http://winhelp2002.mvps.org/hosts.txt\"\n\nDesc:\t \"Fully Qualified Domain Names only - no prefix to strip\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"yoyo\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n"
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
	want := []string{"files", "pre-configured", "urls"}
	c := &edgeos.Config{}
	c.Parms = edgeos.NewParms()
	c.SetOpt(
		edgeos.Dir("/tmp"),
		edgeos.Ext("blacklist.conf"),
		edgeos.Nodes([]string{"domains", "hosts"}),
		edgeos.STypes(want),
	)
	Equals(t, want, c.STypes())
}

func TestBoolToString(t *testing.T) {
	Equals(t, edgeos.True, edgeos.BooltoStr(true))
	Equals(t, edgeos.False, edgeos.BooltoStr(false))
}

func TestToBool(t *testing.T) {
	Equals(t, true, edgeos.StrToBool(edgeos.True))
	Equals(t, false, edgeos.StrToBool(edgeos.False))
}

func TestGetAll(t *testing.T) {
	r := &edgeos.CFGstatic{Cfg: tdata.Cfg}
	c, err := edgeos.ReadCfg(r)
	OK(t, err)
	c.Parms = edgeos.NewParms()
	c.Parms.SetOpt(
		edgeos.Dir("/tmp"),
		edgeos.Ext(".blacklist.conf"),
		edgeos.Nodes([]string{"domains", "hosts"}),
		edgeos.STypes([]string{"files", "pre-configured", "urls"}),
	)

	Equals(t, "[\nDesc:\t \"Blocking mobile ad providers and some analytics providers\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"adaway\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://adaway.org/hosts.txt\"\n \nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"urls\"\nName:\t \"malc0de\"\nnType:\t \"domain\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n \nDesc:\t \"127.0.0.1 based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"malwaredomainlist\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://www.malwaredomainlist.com/hostslist/hosts.txt\"\n \nDesc:\t \"OpenPhish automatic phishing detection\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"openphish\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"https://openphish.com/feed.txt\"\n \nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"\"\n \nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"Zero based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"someonewhocares\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0\"\nType:\t \"hosts\"\nURL:\t \"http://someonewhocares.org/hosts/zero/\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"/config/user-data/blist.hosts.src\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"files\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"Ad server blacklists\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"volkerschatz\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"http://www.volkerschatz.com/net/adpaths\"\n \nDesc:\t \"Zero based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"urls\"\nName:\t \"winhelp2002\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0 \"\nType:\t \"hosts\"\nURL:\t \"http://winhelp2002.mvps.org/hosts.txt\"\n \nDesc:\t \"Fully Qualified Domain Names only - no prefix to strip\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"yoyo\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n]", fmt.Sprint(c.GetAll().S))
	// fmt.Println(c.GetAll().S)
}

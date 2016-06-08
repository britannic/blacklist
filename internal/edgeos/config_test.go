package edgeos_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
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
	c, err := edgeos.ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)
	edgeos.NewParms(c).SetOpt(
		edgeos.Dir("/tmp"),
		edgeos.Ext("blacklist.conf"),
		edgeos.Nodes([]string{"domains", "hosts"}),
	)

	excludes := []string{
		"coremetrics.com",
		"hulu.com",
		"images-amazon.com",
		"skype.com",
		"akamai.net",
		"schema.org",
		"googleusercontent.com",
		"avast.com",
		"google.com",
		"smacargo.com",
		"windows.net",
		"adobedtm.com",
		"gvt1.com",
		"ytimg.com",
		"apple.com",
		"bitdefender.com",
		"freedns.afraid.org",
		"githubusercontent.com",
		"gstatic.com",
		"gvt1.net",
		"ask.com",
		"msdn.com",
		"paypal.com",
		"ssl-on9.com",
		"amazonaws.com",
		"1e100.net",
		"googleadservices.com",
		"googleapis.com",
		"hb.disney.go.com",
		"rackcdn.com",
		"122.2o7.net",
		"cdn.visiblemeasures.com",
		"cloudfront.net",
		"edgesuite.net",
		"github.com",
		"hp.com",
		"sourceforge.net",
		"ssl-on9.net",
		"amazon.com",
		"storage.googleapis.com",
		"yimg.com",
		"static.chartbeat.com",
	}
	want := edgeos.UpdateList(excludes)
	Equals(t, want, c.Get("blacklist").Excludes())
}

func TestFiles(t *testing.T) {
	c, err := edgeos.ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)
	edgeos.NewParms(c).SetOpt(
		edgeos.Dir("/tmp"),
		edgeos.Ext("blacklist.conf"),
		edgeos.Nodes([]string{"domains", "hosts"}),
		edgeos.STypes([]string{"files", "pre-configured", "urls"}),
	)

	want := map[string]string{"domains": "/tmp/domains.malc0de.blacklist.conf\n/tmp/domains.pre-configured.blacklist.conf", "hosts": "/tmp/hosts.adaway.blacklist.conf\n/tmp/hosts.malwaredomainlist.blacklist.conf\n/tmp/hosts.openphish.blacklist.conf\n/tmp/hosts.pre-configured.blacklist.conf\n/tmp/hosts.someonewhocares.blacklist.conf\n/tmp/hosts.tasty.blacklist.conf\n/tmp/hosts.volkerschatz.blacklist.conf\n/tmp/hosts.winhelp2002.blacklist.conf\n/tmp/hosts.yoyo.blacklist.conf"}

	for _, node := range []string{"domains", "hosts"} {
		got := c.Get(node).Source("all").Files().String()
		Equals(t, want[node], got)
	}
}

func TestNodes(t *testing.T) {
	c, err := edgeos.ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)
	edgeos.NewParms(c).SetOpt(
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
	_, got := edgeos.ReadCfg(bytes.NewBufferString(""))
	Equals(t, want, got)
}

func TestRemove(t *testing.T) {
	var (
		dir    = "../testdata"
		ext    = "blacklist.conf"
		nodes  = []string{"domains", "hosts"}
		stypes = []string{"files", "pre-configured", "urls"}
	)
	c, err := edgeos.ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)
	edgeos.NewParms(c).SetOpt(
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

	dlist, err := ioutil.ReadDir(dir)
	OK(t, err)
	var got []string
	for _, f := range dlist {
		if strings.Contains(f.Name(), ext) {
			got = append(got, dir+"/"+f.Name())
		}
	}
	Equals(t, want, got)
	// fmt.Println(edgeos.DiffArray(want, got))
}

func TestSource(t *testing.T) {
	c, err := edgeos.ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)
	edgeos.NewParms(c).SetOpt(
		edgeos.Dir("/tmp"),
		edgeos.Ext("blacklist.conf"),
		edgeos.Nodes([]string{"domains", "hosts"}),
		edgeos.STypes([]string{"files", "pre-configured", "urls"}),
	)

	var (
		wdomains = "\nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"\"\n\nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"urls\"\nName:\t \"malc0de\"\nnType:\t \"domain\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n"

		whosts = "\nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"\"\n\nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"urls\"\nName:\t \"malc0de\"\nnType:\t \"domain\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n\nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n\nDesc:\t \"Blocking mobile ad providers and some analytics providers\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"adaway\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://adaway.org/hosts.txt\"\n\nDesc:\t \"127.0.0.1 based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"\"\nLtype:\t \"urls\"\nName:\t \"malwaredomainlist\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://www.malwaredomainlist.com/hostslist/hosts.txt\"\n\nDesc:\t \"OpenPhish automatic phishing detection\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"\"\nLtype:\t \"urls\"\nName:\t \"openphish\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"https://openphish.com/feed.txt\"\n\nDesc:\t \"Zero based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"\"\nLtype:\t \"urls\"\nName:\t \"someonewhocares\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0\"\nType:\t \"hosts\"\nURL:\t \"http://someonewhocares.org/hosts/zero/\"\n\nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"/config/user-data/blist.hosts.src\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"files\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n\nDesc:\t \"Ad server blacklists\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"\"\nLtype:\t \"urls\"\nName:\t \"volkerschatz\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"http://www.volkerschatz.com/net/adpaths\"\n\nDesc:\t \"Zero based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"urls\"\nName:\t \"winhelp2002\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0 \"\nType:\t \"hosts\"\nURL:\t \"http://winhelp2002.mvps.org/hosts.txt\"\n\nDesc:\t \"Fully Qualified Domain Names only - no prefix to strip\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"\"\nLtype:\t \"urls\"\nName:\t \"yoyo\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n"
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
	edgeos.NewParms(c).SetOpt(
		edgeos.Dir("/tmp"),
		edgeos.Ext("blacklist.conf"),
		edgeos.Nodes([]string{"domains", "hosts"}),
		edgeos.STypes(want),
	)
	Equals(t, want, c.STypes())
}

func TestToBool(t *testing.T) {
	Equals(t, true, edgeos.StrToBool("true"))
}

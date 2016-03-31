package data_test

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/data"
	g "github.com/britannic/blacklist/global"
	"github.com/britannic/blacklist/regx"
	"github.com/britannic/blacklist/utils"
	. "github.com/britannic/testutils"
)

var dmsqdir, logfile string

func init() {
	switch g.WhatOS {
	case g.TestOS:
		dmsqdir = "../testdata"
		logfile = dmsqdir + "/blacklist.log"

	default:
		dmsqdir = g.DmsqDir
	}
	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND, 0755)
	if err == nil {
		log.SetFormatter(&log.TextFormatter{DisableColors: true})
		log.SetOutput(f)
	}
}

func shuffleArray(slice []string) {
	rand.Seed(time.Now().UnixNano())
	n := len(slice)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func TestDiffArray(t *testing.T) {
	biggest := []string{"one", "two", "three", "four", "five", "six"}
	smallest := []string{"one", "two", "three"}
	want := []string{"five", "four", "six"}

	got := data.DiffArray(biggest, smallest)
	Equals(t, want, got)

	got = data.DiffArray(smallest, biggest)
	Equals(t, want, got)

	shuffleArray(biggest)
	got = data.DiffArray(smallest, biggest)
	Equals(t, want, got)

	shuffleArray(smallest)
	got = data.DiffArray(smallest, biggest)
	Equals(t, want, got)
}

func TestExclusions(t *testing.T) {
	b, err := config.Get(config.Testdata, g.Area.Root)
	OK(t, err)

	var (
		// tdata  string
		dex    = make(config.Dict)
		ex     = data.GetExcludes(*b)
		globex = data.GetExcludes(*b)
	)

	for _, s := range src {
		f := fmt.Sprintf("../testdata/tdata.%v.%v", s.Type, s.Name)
		testdata, err := utils.GetFile(f)
		OK(t, err)

		gdata := data.Process(s, globex, dex, testdata)

		for k := range gdata.List {
			i := strings.Count(k, ".")
			Assert(t, i >= 1, fmt.Sprintf(`key: %v has "." count of %v`, k, i), k)

			switch {
			case i == 1:
				Assert(t, !ex.KeyExists(k), fmt.Sprintf("Exclusion failure, found matching key: %v", k), k)

			case i > 1:
				Assert(t, !ex.SubKeyExists(k), fmt.Sprintf("Exclusion failure, found submatch for key: %v", k), k)
			}
		}
	}
}

func TestGetHTTP(t *testing.T) {
	type tdata struct {
		body  []byte
		err   error
		prcsd *config.Src
	}

	h := &tdata{}
	d := []*tdata{}
	rx := regx.Regex

	b, err := config.Get(config.Testdata, g.Area.Root)
	OK(t, err)

	a := data.GetURLs(*b)
	ex := make(config.Dict)
	dex := make(config.Dict)

	for k := range a {
		for _, u := range a[k] {
			if len(u.URL) > 0 {
				h.body, h.err = data.GetHTTP(u.URL)
				d = append(d, h)
				bdata := bufio.NewScanner(strings.NewReader(string(h.body[:])))
				h.prcsd = data.Process(u, ex, dex, bdata)
			}
		}
	}

	for _, z := range d {
		for got := range z.prcsd.List {
			want := rx.FQDN.FindStringSubmatch(got)[1]
			Equals(t, want, got)
		}
	}
}

func TestGetUrls(t *testing.T) {
	blist, err := config.Get(config.Testdata, g.Area.Root)
	OK(t, err)

	b := *blist
	a := data.GetURLs(b)

	for k := range b {
		for _, url := range a[k] {
			if g, ok := b[k].Source[url.Name]; ok {
				want := g.URL
				Assert(t, want == url.URL, fmt.Sprintf("%v URL mismatch:", url.Name), url)
			}
		}
	}
}

func TestIsDisabled(t *testing.T) {
	c := make(config.Blacklist)
	c[g.Area.Root] = &config.Node{}

	c[g.Area.Root].Disable = true
	Equals(t, true, data.IsDisabled(c, g.Area.Root))

	c[g.Area.Root].Disable = false
	Equals(t, false, data.IsDisabled(c, g.Area.Root))
}

func TestListFiles(t *testing.T) {
	b, err := config.Get(config.Testdata, g.Area.Root)
	OK(t, err)

	urls := data.GetURLs(*b)

	var want []string

	for _, k := range data.SortKeys(urls) {
		for _, s := range urls[k] {
			want = append(want, fmt.Sprintf(g.FStr, dmsqdir, s.Type, s.Name))
		}
	}

	sort.Strings(want)
	got, err := data.ListFiles(dmsqdir)
	OK(t, err)
	Equals(t, want, got)

}

func TestProcess(t *testing.T) {
	for _, s := range src {
		ex := make(config.Dict)
		dex := make(config.Dict)
		f := fmt.Sprintf("%v/tdata.%v.%v", dmsqdir, s.Type, s.Name)
		testdata, err := utils.GetFile(f)
		OK(t, err)

		f = fmt.Sprintf("%v/sdata.%v.%v", dmsqdir, s.Type, s.Name)
		staticdata, err := utils.GetFile(f)
		OK(t, err)

		var wdata string
		for staticdata.Scan() {
			wdata += staticdata.Text() + "\n"
		}

		gdata := string(data.GetList(data.Process(s, ex, dex, testdata))[:])

		Equals(t, wdata[:], gdata)
	}
}

func TestPurgeFiles(t *testing.T) {

	b, err := config.Get(config.Testdata, g.Area.Root)
	OK(t, err)

	urls := data.GetURLs(*b)

	var want []string

	for _, k := range data.SortKeys(urls) {
		for _, s := range urls[k] {
			want = append(want, fmt.Sprintf(g.FStr, dmsqdir, s.Type, s.Name))
			for i := 0; i < 5; i++ {
				data := []byte{84, 104, 101, 32, 114, 101, 115, 116, 32, 105, 115, 32, 104, 105, 115, 116, 111, 114, 121, 33}
				err := utils.WriteFile(fmt.Sprintf("%v/%v.%v[%v]%v", dmsqdir, s.Type, s.Name, i, g.Fext), data)
				OK(t, err)
			}
		}
	}

	err = data.PurgeFiles(urls, dmsqdir)
	OK(t, err)

	got, err := data.ListFiles(dmsqdir)
	OK(t, err)

	sort.Strings(want)
	Equals(t, want, got)

	err = data.PurgeFiles(urls, "Break it Z...")
	NotOK(t, err)
}

func TestStripPrefixAndSuffix(t *testing.T) {
	rx := regx.Regex
	tline := `This is a complete sentence and should not have a comment.`

	for _, s := range src {
		var l string
		switch s.Prfx {
		case "http":
			l = s.Prfx + "://" + tline
		default:
			l = s.Prfx + tline
		}

		l += ` # Comment.`

		r, ok := data.StripPrefixAndSuffix(l, s.Prfx, rx)
		Assert(t, tline == r, fmt.Sprintf("stripPrefix() failed for %v", s.Name), r)
		Assert(t, ok, fmt.Sprintf("stripPrefix() failed for %v", s.Name), ok)

		r, ok = data.StripPrefixAndSuffix(tline, "http", rx)
		Assert(t, tline == r, fmt.Sprintf("stripPrefix() failed for %v", s.Name), r)
		Assert(t, ok == false, fmt.Sprintf("stripPrefix() failed for %v", s.Name), ok)
	}
}

// http://play.golang.org/p/KAwluDqGIl
var src = []*config.Src{
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "pre-configured",
		Type:    "domains",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "malc0de",
		Prfx:    "zone ",
		Type:    "domains",
		URL:     "http://malc0de.com/bl/ZONES",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "pre-configured",
		Type:    "hosts",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "adaway",
		Prfx:    "127.0.0.1 ",
		Type:    "hosts",
		URL:     "http://adaway.org/hosts.txt",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "malwaredomainlist",
		Prfx:    "127.0.0.1 ",
		Type:    "hosts",
		URL:     "http://www.malwaredomainlist.com/hostslist/hosts.txt",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "openphish",
		Prfx:    "http",
		Type:    "hosts",
		URL:     "https://openphish.com/feed.txt",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "someonewhocares",
		Prfx:    "0.0.0.0",
		Type:    "hosts",
		URL:     "http://someonewhocares.org/hosts/zero/",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "volkerschatz",
		Prfx:    "http",
		Type:    "hosts",
		URL:     "http://www.volkerschatz.com/net/adpaths",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "winhelp2002",
		Prfx:    "0.0.0.0 ",
		Type:    "hosts",
		URL:     "http://winhelp2002.mvps.org/hosts.txt",
	},
	{
		Disable: false,
		IP:      "0.0.0.0",
		Name:    "yoyo",
		Type:    "hosts",
		URL:     "http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext",
	},
}

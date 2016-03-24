package data_test

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/data"
	g "github.com/britannic/blacklist/global"
	"github.com/britannic/blacklist/regx"
	"github.com/britannic/blacklist/utils"
)

var dmsqDir string

func init() {
	switch g.WhatOS {
	case g.TestOS:
		dmsqDir = "../testdata"
	default:
		dmsqDir = g.DmsqDir
	}
}

func TestExclusions(t *testing.T) {
	b, err := config.Get(config.Testdata, g.Area.Root)
	if err != nil {
		t.Error("Couldn't load config.Testdata")
	}

	var (
		// tdata  string
		dex    = make(config.Dict)
		ex     = data.GetExcludes(*b)
		globex = data.GetExcludes(*b)
	)

	for _, s := range src {
		f := fmt.Sprintf("../testdata/tdata.%v.%v", s.Type, s.Name)
		testdata, err := utils.GetFile(f)
		if err != nil {
			t.Errorf("Cannot open %v", f)
		}

		gdata := data.Process(s, globex, dex, testdata)

		for k := range gdata.List {
			i := strings.Count(k, ".")
			if i < 1 {
				t.Errorf("key: %v has . count of %v", k, i)
			}

			switch {
			case i == 1:
				if ex.KeyExists(k) {
					t.Errorf("Exclusion failure, found matching key: %v", k)
				}
			case i > 1:
				if ex.SubKeyExists(k) {
					t.Errorf("Exclusion failure, found submatch for key: %v", k)
				}
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
	if err != nil {
		t.Errorf("unable to get configuration data, error code: %v\n", err)
	}

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
			if strings.Compare(got, want) != 0 {
				t.Errorf("wanted: %q, got: %q", want, got)
			}
		}
	}
}

func TestGetUrls(t *testing.T) {
	blist, err := config.Get(config.Testdata, g.Area.Root)
	if err != nil {
		t.Errorf("unable to get configuration data, error code: %v\n", err)
	}

	b := *blist
	a := data.GetURLs(b)

	for k := range b {
		for _, url := range a[k] {
			if g, ok := b[k].Source[url.Name]; ok {
				want := g.URL
				got := url.URL
				if want != url.URL {
					t.Errorf("%v URL mismatch:", url.Name)
					fmt.Printf("Wanted %v\nGot: %v", want, got)
				}
			}
		}
	}
}

func TestIsDisabled(t *testing.T) {
	c := make(config.Blacklist)
	c[g.Area.Root] = &config.Node{}
	c[g.Area.Root].Disable = true
	if data.IsDisabled(c, g.Area.Root) != true {
		t.Error("Should be true")
	}
	c[g.Area.Root].Disable = false
	if data.IsDisabled(c, g.Area.Root) != false {
		t.Error("Should be false")
	}
}

func TestProcess(t *testing.T) {
	for _, s := range src {
		ex := make(config.Dict)
		dex := make(config.Dict)
		f := fmt.Sprintf("%v/tdata.%v.%v", dmsqDir, s.Type, s.Name)
		testdata, err := utils.GetFile(f)
		if err != nil {
			t.Errorf("Cannot open %v", f)
		}

		f = fmt.Sprintf("%v/sdata.%v.%v", dmsqDir, s.Type, s.Name)
		staticdata, err := utils.GetFile(f)
		if err != nil {
			t.Errorf("Cannot open %v", f)
		}

		var wdata string
		for staticdata.Scan() {
			wdata += staticdata.Text() + "\n"
		}
		// wdata = utils.GetByteArray(staticdata, wdata)
		gdata := string(data.GetList(data.Process(s, ex, dex, testdata))[:])

		if !utils.CmpHash([]byte(wdata), []byte(gdata)) {
			mismatch := []*struct {
				d string
				f string
			}{
				{
					d: string(wdata[:]),
					f: fmt.Sprintf("/tmp/want.%v.%v", s.Type, s.Name),
				},
				{
					d: gdata,
					f: fmt.Sprintf("/tmp/got.%v.%v", s.Type, s.Name),
				},
			}

			for _, m := range mismatch {
				utils.WriteFile(m.f, []byte(m.d))
			}
			t.Errorf("data mismatch between standard and data.Processed data for %q.", s.Name)
		}
	}
}

func TestPurgeFiles(t *testing.T) {

	b, err := config.Get(config.Testdata, g.Area.Root)
	if err != nil {
		t.Errorf("unable to get configuration data, error code: %v\n", err)
	}

	urls := data.GetURLs(*b)

	var want, got []string

	for k := range urls {
		for _, s := range urls[k] {
			want = append(want, fmt.Sprintf(g.FStr, dmsqDir, s.Type, s.Name))
			for i := 0; i < 5; i++ {
				data := []byte{84, 104, 101, 32, 114, 101, 115, 116, 32, 105, 115, 32, 104, 105, 115, 116, 111, 114, 121, 33}
				utils.WriteFile(fmt.Sprintf("%v/%v.%v[%v]%v", dmsqDir, s.Type, s.Name, i, g.Fext), data)
			}
		}
	}

	if err := data.PurgeFiles(urls, dmsqDir); err != nil {
		t.Errorf("Error removing unused conf files: %v", err)
	}

	got = data.ListFiles(dmsqDir)
	delta := data.DiffArray(want, got)
	if len(delta) > 0 {
		t.Errorf("Issue purging files, difference: %v\nGot: %v\nWant: %v\n", delta, got, want)
	}
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
		switch {
		case tline != r:
			t.Errorf("stripPrefix() failed for %v", s.Name)
			fmt.Printf("Want: %v\nGot: %v\n", tline, r)
		case !ok:
			t.Errorf("stripPrefix() failed for %v", s.Name)
		}
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

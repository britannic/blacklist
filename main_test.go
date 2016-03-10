package main

import (
	"crypto/md5"
	"fmt"
	"os/user"
	"runtime"
	"testing"

	"github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/regx"
)

func compare(t *testing.T, want, got *user.User) {
	if want.Uid != got.Uid {
		t.Errorf("got Uid=%q; want %q", got.Uid, want.Uid)
	}
}

func TestBasename(t *testing.T) {
	dirPath := "/usr/root/testing.txt"
	dir := basename(dirPath)
	if dir != "testing" {
		t.Error(dir)
	}
}

func TestIsAdmin(t *testing.T) {
	want, err := user.Current()
	if err != nil {
		t.Fatalf("Current: %v", err)
	}

	got, err := user.Lookup(want.Username)
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}

	compare(t, want, got)

	osAdmin := false
	if got.Uid == "0" {
		osAdmin = true
	}

	if isAdmin() != osAdmin {
		t.Error(osAdmin)
	}
}

// cmpHash compares the hashes of a to b and returns true if they're identical
func cmpHash(a, b []byte) bool {
	if md5.Sum(a) == md5.Sum(b) {
		return true
	}
	return false
}

func TestProcess(t *testing.T) {
	for _, s := range src {
		e := make(excludes)
		f := fmt.Sprintf("./tdata.%v.%v", s.Type, s.Name)
		testdata, err := getfile(f)
		if err != nil {
			t.Fatalf("Cannot open %v", f)
		}

		var tdata string
		for _, l := range testdata {
			if len(l) > 0 {
				tdata += l + "\n"
			}
		}

		f = fmt.Sprintf("./sdata.%v.%v", s.Type, s.Name)
		staticdata, err := getfile(f)
		if err != nil {
			t.Fatalf("Cannot open %v", f)
		}

		var wdata string
		for _, l := range staticdata {
			if len(l) > 0 {
				wdata += l + "\n"
			}
		}

		gdata := string(getList(process(s, e, tdata))[:])

		if !cmpHash([]byte(wdata), []byte(gdata)) {
			mismatch := []*struct {
				d string
				f string
			}{
				{
					d: wdata,
					f: fmt.Sprintf("/tmp/want.%v.%v", s.Type, s.Name),
				},
				{
					d: gdata,
					f: fmt.Sprintf("/tmp/got.%v.%v", s.Type, s.Name),
				},
			}

			for _, m := range mismatch {
				writeFile(m.f, []byte(m.d))
			}
			t.Errorf("data mismatch between standard and processed data for %q.", s.Name)
		}
	}
}

func TestWriteAndReadFile(t *testing.T) {
	fname := "/tmp/delete.me"
	data := []byte{}
	data = append(data, `This is a test file. Delete it!`...)
	err := writeFile(fname, data)
	if err != nil {
		t.Error(err)
	}

	fdata, err := getfile(fname)
	fdata = append(fdata, string(data[:]))
	switch {
	case err != nil:
		t.Error(err)
	case fdata[0] != fdata[1]:
		t.Error("data mismatch between writeFile and getFile")
	}
}

func TestStripPrefix(t *testing.T) {
	rx := regx.Regex()
	tline := `[This line should be delimited by "[]" only.]`

	for _, s := range src {
		var l string
		switch s.Prfx {
		case "http":
			l = s.Prfx + "://" + tline
		default:
			l = s.Prfx + tline
		}

		r, ok := stripPrefix(l, s.Prfx, rx)
		switch {
		case tline != r:
			t.Errorf("stripPrefix() failed for %v", s.Name)
			fmt.Printf("Want: %v\nGot: %v\n", tline, r)
		case !ok:
			t.Errorf("stripPrefix() failed for %v", s.Name)
		}
	}
}

func TestGetUrls(t *testing.T) {
	blist, err := config.Get(config.Testdata, root)
	if err != nil {
		t.Errorf("unable to get configuration data, error code: %v\n", err)
	}

	b := *blist
	a := getURLs(b)

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

func TestPurgeFiles(t *testing.T) {
	whatOS := runtime.GOOS
	if whatOS == "darwin" {
		dmsqDir = "/tmp"
		logfile = "/tmp/blacklist.log"
	}

	b, err := config.Get(config.Testdata, root)
	if err != nil {
		t.Errorf("unable to get configuration data, error code: %v\n", err)
	}

	urls := getURLs(*b)
	if err := purgeFiles(urls); err != nil {
		t.Errorf("Error removing unused conf files: %v", err)
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

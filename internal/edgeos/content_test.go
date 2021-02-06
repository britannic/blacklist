package edgeos

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/sync/errgroup"
)

type dummyConfig struct {
	*Env
	s []string
	t *testing.T
}

func (d *dummyConfig) ProcessContent(cts ...Contenter) error {
	var (
		a, b  int32
		area  string
		tally = &stats{dropped: a, kept: b}
	)

	for _, ct := range cts {
		o := ct.GetList().src
		for _, src := range o {
			area = typeInt(src.nType)
			src.ctr.Lock()
			src.ctr.stat[area] = tally
			src.ctr.Unlock()
			b, _ := ioutil.ReadAll(src.process().r)
			d.s = append(d.s, strings.TrimSuffix(string(b), "\n"))
		}
	}
	return nil
}

func TestConfigProcessContent(t *testing.T) {
	Convey("Testing ProcessContent()", t, func() {
		newCfg := func() *Config {
			return NewConfig(
				API("/bin/cli-shell-api"),
				Arch(runtime.GOARCH),
				Bash("/bin/bash"),
				Cores(runtime.NumCPU()),
				Dir("/tmp"),
				DNSsvc("service dnsmasq restart"),
				Ext("blacklist.conf"),
				FileNameFmt("%v/%v.%v.%v"),
				InCLI("inSession"),
				Logger(newLog()),
				Method("GET"),
				Prefix("address=", "addn-hosts="),
				Timeout(30*time.Second),
				WCard(Wildcard{Node: "*s", Name: "*"}),
			)
		}

		tests := []struct {
			c      *Config
			cfg    string
			ct     IFace
			err    error
			expErr bool
			name   string
		}{
			{
				c:      newCfg(),
				cfg:    testallCfg,
				ct:     URLhObj,
				err:    errors.New("Get \"http://127.0.0.1:8081/hosts/host.txt\": dial tcp 127.0.0.1:8081: connect: connection refused"),
				expErr: true,
				name:   "Hosts blacklist source",
			},
			{
				c:      newCfg(),
				cfg:    testallCfg,
				ct:     URLdObj,
				err:    errors.New("Get \"http://127.0.0.1:8081/domains/domain.txt\": dial tcp 127.0.0.1:8081: connect: connection refused"),
				expErr: true,
				name:   "Domains blacklist source",
			},
			{
				c:      newCfg(),
				cfg:    testallCfg,
				ct:     FileObj,
				err:    errors.New("open /:~//hosts.tasty.blacklist.conf: no such file or directory"),
				expErr: true,
				name:   "File source",
			},
			{
				c:      newCfg(),
				cfg:    testallCfg,
				ct:     ExHtObj,
				err:    errors.New("open /:~//hosts.whitelisted-servers.blacklist.conf: no such file or directory"),
				expErr: true,
				name:   "Whitelisted hosts",
			},
			{
				c:      newCfg(),
				cfg:    testallCfg,
				ct:     ExDmObj,
				err:    errors.New("open /:~//domains.whitelisted-subdomains.blacklist.conf: no such file or directory"),
				expErr: true,
				name:   "Whitelisted domains",
			},
			{
				c:      newCfg(),
				cfg:    testallCfg,
				ct:     PreHObj,
				err:    errors.New("open /:~//hosts.blacklisted-servers.blacklist.conf: no such file or directory"),
				expErr: true,
				name:   "Blacklisted hosts",
			},
			{
				c:      newCfg(),
				cfg:    testallCfg,
				ct:     PreDObj,
				err:    errors.New("open /:~//domains.blacklisted-subdomains.blacklist.conf: no such file or directory"),
				expErr: true,
				name:   "Blacklisted domains",
			},
			{
				c:      newCfg(),
				cfg:    testallCfg,
				ct:     ExRtObj,
				err:    errors.New("open /:~//roots.global-whitelisted-domains.blacklist.conf: no such file or directory"),
				expErr: true,
				name:   "Global whitelist",
			},
			{
				c:      newCfg(),
				cfg:    testCfg,
				ct:     FileObj,
				err:    fmt.Errorf("open /:~/=../../internal/testdata/blist.hosts.src: no such file or directory"),
				expErr: true,
				name:   "Non-existent File source",
			},
		}
		for _, tt := range tests {
			Convey("current test: "+tt.name, func() {
				if tt.name == "" {
					tt.c.Dir = "/:~/"
				}
				So(tt.c.Blacklist(&CFGstatic{Cfg: tt.cfg}), ShouldBeNil)

				obj, err := tt.c.NewContent(tt.ct)
				So(err, ShouldBeNil)

				err = tt.c.ProcessContent(obj)
				if (err != nil) == tt.expErr {
					So(err.Error(), ShouldEqual, tt.err.Error())
				}
			})
		}

		Convey("Testing ProcessContent() if no arguments ", func() {
			// var g errgroup.Group
			// g.Go(func() error { return newCfg().ProcessContent() })
			// err := g.Wait()
			So(newCfg().ProcessContent(), ShouldNotBeNil)
		})
	})
}

func TestNewContent(t *testing.T) {
	Convey("Testing NewContent()", t, func() {
		expFileObj := "addn-hosts=/0.really.bad.phishing.site.ru/0.0.0.0\naddn-hosts=/cw.bad.ultraadverts.site.eu/0.0.0.0\naddn-hosts=/really.bad.phishing.site.ru/0.0.0.0\naddn-hosts=/www.ytimg.com/0.0.0.0\naddn-hosts=/ytimg.com/0.0.0.0"
		tests := []struct {
			err       error
			exp       string
			fail      bool
			i         int
			ltype     string
			name      string
			obj       IFace
			page      string
			page2     string
			pageData  string
			pageData2 string
			pos       int
			svr       *HTTPserver
			svr2      *HTTPserver
		}{
			{
				i:     1,
				exp:   excRootContent,
				fail:  false,
				ltype: ExcRoots,
				name:  "z" + ExcRoots,
				obj:   ExRtObj,
				pos:   -1,
			},
			{
				i:     1,
				exp:   "address=/adinfuse.com/#",
				fail:  false,
				ltype: ExcDomns,
				name:  "z" + ExcDomns,
				obj:   ExDmObj,
				pos:   -1,
			},
			{
				i:     1,
				exp:   "address=/wv.inner-active.mobi/#",
				fail:  false,
				ltype: ExcHosts,
				name:  "z" + ExcHosts,
				obj:   ExHtObj,
				pos:   -1,
			},
			{
				i:     1,
				exp:   "address=/adsrvr.org/192.1.1.1\naddress=/adtechus.net/192.1.1.1\naddress=/advertising.com/192.1.1.1\naddress=/centade.com/192.1.1.1\naddress=/doubleclick.net/192.1.1.1\naddress=/free-counter.co.uk/192.1.1.1\naddress=/intellitxt.com/192.1.1.1\naddress=/kiosked.com/192.1.1.1",
				fail:  false,
				ltype: PreDomns,
				name:  "z" + PreDomns,
				obj:   PreDObj,
				pos:   -1,
			},
			{
				i:     1,
				exp:   "addn-hosts=/beap.gemini.yahoo.com/0.0.0.0",
				fail:  false,
				ltype: PreHosts,
				name:  "z" + PreHosts,
				obj:   PreHObj,
				pos:   -1,
			},
			{
				i:     1,
				exp:   expFileObj,
				fail:  false,
				ltype: files,
				name:  "ztasty",
				obj:   FileObj,
				pos:   -1,
			},
			{
				i:         1,
				exp:       domainsContent,
				fail:      false,
				ltype:     urls,
				name:      "zmalc0de",
				obj:       URLdObj,
				pos:       -1,
				page:      "/hosts.txt",
				page2:     "/domains.txt",
				pageData:  httpHostData,
				pageData2: HTTPDomainData,
				svr:       new(HTTPserver),
				svr2:      new(HTTPserver),
			},
			{
				i:         1,
				exp:       hostsContent,
				fail:      false,
				ltype:     urls,
				name:      "zadaway",
				obj:       URLhObj,
				pos:       -1,
				page:      "/hosts.txt",
				page2:     "/domains.txt",
				pageData:  httpHostData,
				pageData2: HTTPDomainData,
				svr:       new(HTTPserver),
				svr2:      new(HTTPserver),
			},
			{
				i:     1,
				exp:   excRootContent,
				fail:  false,
				ltype: ExcRoots,
				name:  ExcRoots,
				obj:   ExRtObj,
				pos:   0,
			},
			{
				i:     1,
				exp:   "",
				fail:  false,
				ltype: PreRoots,
				name:  "z" + PreRoots,
				obj:   PreRObj,
				pos:   -1,
			},
			{
				i:     1,
				exp:   "address=/adinfuse.com/#",
				fail:  false,
				ltype: ExcDomns,
				name:  ExcDomns,
				obj:   ExDmObj,
				pos:   0,
			},
			{
				i:     1,
				exp:   "address=/wv.inner-active.mobi/#",
				fail:  false,
				ltype: ExcHosts,
				name:  ExcHosts,
				obj:   ExHtObj,
				pos:   0,
			},
			{
				i:     1,
				exp:   "address=/adsrvr.org/192.1.1.1\naddress=/adtechus.net/192.1.1.1\naddress=/advertising.com/192.1.1.1\naddress=/centade.com/192.1.1.1\naddress=/doubleclick.net/192.1.1.1\naddress=/free-counter.co.uk/192.1.1.1\naddress=/intellitxt.com/192.1.1.1\naddress=/kiosked.com/192.1.1.1",
				fail:  false,
				ltype: PreDomns,
				name:  PreDomns,
				obj:   PreDObj,
				pos:   0,
			},
			{
				i:     1,
				exp:   "addn-hosts=/beap.gemini.yahoo.com/0.0.0.0",
				fail:  false,
				ltype: PreHosts,
				name:  PreHosts,
				obj:   PreHObj,
				pos:   0,
			},
			{
				i:     1,
				exp:   "",
				fail:  false,
				ltype: PreRoots,
				name:  PreRoots,
				obj:   PreRObj,
				pos:   0,
			},
			{
				i:     1,
				exp:   expFileObj,
				fail:  false,
				ltype: files,
				name:  "tasty",
				obj:   FileObj,
				pos:   0,
			},
			{
				i:         1,
				exp:       domainsContent,
				fail:      false,
				ltype:     urls,
				name:      "malc0de",
				obj:       URLdObj,
				pos:       0,
				page:      "/hosts.txt",
				page2:     "/domains.txt",
				pageData:  httpHostData,
				pageData2: HTTPDomainData,
				svr:       new(HTTPserver),
				svr2:      new(HTTPserver),
			},
			{
				i:         1,
				exp:       hostsContent,
				fail:      false,
				ltype:     urls,
				name:      "adaway",
				obj:       URLhObj,
				pos:       0,
				page:      "/hosts.txt",
				page2:     "/domains.txt",
				pageData:  httpHostData,
				pageData2: HTTPDomainData,
				svr:       new(HTTPserver),
				svr2:      new(HTTPserver),
			},
			{
				i:    0,
				err:  errors.New("invalid interface requested"),
				fail: true,
				obj:  Invalid,
				pos:  -1,
			},
		}

		c := NewConfig(
			API("/bin/cli-shell-api"),
			Arch(runtime.GOARCH),
			Bash("/bin/bash"),
			Cores(runtime.NumCPU()),
			Dir("/tmp"),
			Disabled(false),
			DNSsvc("service dnsmasq restart"),
			Ext("blacklist.conf"),
			FileNameFmt("%v/%v.%v.%v"),
			InCLI("inSession"),
			Logger(newLog()),
			Method("GET"),
			Prefix("address=", "addn-hosts="),
			Timeout(30*time.Second),
			WCard(Wildcard{Node: "*s", Name: "*"}),
		)

		So(c.Blacklist(&CFGstatic{Cfg: Cfg}), ShouldBeNil)

		c.Dex.merge(&list{RWMutex: &sync.RWMutex{}, entry: entry{"amazon-de.com": struct{}{}}})
		So(c.Dex.String(), ShouldEqual, `"amazon-de.com":{},
`)

		for _, tt := range tests {
			Convey("processing "+tt.name, func() {
				objs, err := c.NewContent(tt.obj)
				if tt.ltype == urls {
					uri1 := tt.svr.NewHTTPServer().String() + tt.page
					objs.SetURL("adaway", uri1)
					uri2 := tt.svr2.NewHTTPServer().String() + tt.page2
					objs.SetURL("malc0de", uri2)

					go tt.svr.Mux.HandleFunc(tt.page,
						func(w http.ResponseWriter, r *http.Request) {
							fmt.Fprint(w, tt.pageData)
						},
					)

					go tt.svr2.Mux.HandleFunc(tt.page2,
						func(w http.ResponseWriter, r *http.Request) {
							fmt.Fprint(w, tt.pageData2)
						},
					)
				}

				switch tt.fail {
				case false:
					So(err, ShouldBeNil)

					d := &dummyConfig{Env: c.Env, t: t}
					d.ProcessContent(objs)

					So(strings.Join(d.s, "\n"), ShouldEqual, tt.exp)

					objs.SetURL(tt.name, tt.name)

					So(objs.Find(tt.name), ShouldEqual, tt.pos)
					So(objs.Len(), ShouldEqual, tt.i)

				default:
					So(err.Error(), ShouldEqual, tt.err.Error())
				}
			})
		}
	})
}

func TestContenterString(t *testing.T) {
	Convey("Testing GetAllContent()", t, func() {
		c := NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
			FileNameFmt("%v/%v.%v.%v"),
			Method("GET"),
			Prefix("address=", "addn-hosts="),
		)

		So(c.Blacklist(&CFGstatic{Cfg: testallCfg}), ShouldBeNil)

		tests := []struct {
			iFace IFace
			exp   string
			name  string
		}{
			{name: "ExDmObj", iFace: ExDmObj, exp: "\nDesc:         \"pre-configured whitelisted subdomains\"\nDisabled:     \"false\"\nFile:         \"**Undefined**\"\nIP:           \"0.0.0.0\"\nLtype:        \"whitelisted-subdomains\"\nName:         \"whitelisted-subdomains\"\nnType:        \"excDomn\"\nPrefix:       \"**Undefined**\"\nType:         \"whitelisted-subdomains\"\nURL:          \"**Undefined**\"\nWhitelist:\n              \"**No entries found**\"\nBlacklist:\n              \"**No entries found**\"\n"},
			{name: "ExHtObj", iFace: ExHtObj, exp: "\nDesc:         \"pre-configured whitelisted servers\"\nDisabled:     \"false\"\nFile:         \"**Undefined**\"\nIP:           \"192.168.168.1\"\nLtype:        \"whitelisted-servers\"\nName:         \"whitelisted-servers\"\nnType:        \"excHost\"\nPrefix:       \"**Undefined**\"\nType:         \"whitelisted-servers\"\nURL:          \"**Undefined**\"\nWhitelist:\n              \"**No entries found**\"\nBlacklist:\n              \"**No entries found**\"\n"},
			{name: "ExRtObj", iFace: ExRtObj, exp: "\nDesc:         \"pre-configured global whitelisted domains\"\nDisabled:     \"false\"\nFile:         \"**Undefined**\"\nIP:           \"0.0.0.0\"\nLtype:        \"global-whitelisted-domains\"\nName:         \"global-whitelisted-domains\"\nnType:        \"excRoot\"\nPrefix:       \"**Undefined**\"\nType:         \"global-whitelisted-domains\"\nURL:          \"**Undefined**\"\nWhitelist:\n              \"122.2o7.net\"\n              \"1e100.net\"\n              \"adobedtm.com\"\n              \"akamai.net\"\n              \"amazon.com\"\n              \"amazonaws.com\"\n              \"apple.com\"\n              \"ask.com\"\n              \"avast.com\"\n              \"bitdefender.com\"\n              \"cdn.visiblemeasures.com\"\n              \"cloudfront.net\"\n              \"coremetrics.com\"\n              \"edgesuite.net\"\n              \"freedns.afraid.org\"\n              \"github.com\"\n              \"githubusercontent.com\"\n              \"google.com\"\n              \"googleadservices.com\"\n              \"googleapis.com\"\n              \"googleusercontent.com\"\n              \"gstatic.com\"\n              \"gvt1.com\"\n              \"gvt1.net\"\n              \"hb.disney.go.com\"\n              \"hp.com\"\n              \"hulu.com\"\n              \"images-amazon.com\"\n              \"msdn.com\"\n              \"paypal.com\"\n              \"rackcdn.com\"\n              \"schema.org\"\n              \"skype.com\"\n              \"smacargo.com\"\n              \"sourceforge.net\"\n              \"ssl-on9.com\"\n              \"ssl-on9.net\"\n              \"static.chartbeat.com\"\n              \"storage.googleapis.com\"\n              \"windows.net\"\n              \"yimg.com\"\n              \"ytimg.com\"\nBlacklist:\n              \"**No entries found**\"\n"},
			{name: "FileObj", iFace: FileObj, exp: "\nDesc:         \"File source\"\nDisabled:     \"false\"\nFile:         \"../../internal/testdata/blist.hosts.src\"\nIP:           \"0.0.0.0\"\nLtype:        \"file\"\nName:         \"tasty\"\nnType:        \"host\"\nPrefix:       \"**Undefined**\"\nType:         \"hosts\"\nURL:          \"**Undefined**\"\nWhitelist:\n              \"**No entries found**\"\nBlacklist:\n              \"**No entries found**\"\n"},
			{name: "PreDObj", iFace: PreDObj, exp: "\nDesc:         \"pre-configured blacklisted subdomains\"\nDisabled:     \"false\"\nFile:         \"**Undefined**\"\nIP:           \"0.0.0.0\"\nLtype:        \"blacklisted-subdomains\"\nName:         \"blacklisted-subdomains\"\nnType:        \"preDomn\"\nPrefix:       \"**Undefined**\"\nType:         \"blacklisted-subdomains\"\nURL:          \"**Undefined**\"\nWhitelist:\n              \"**No entries found**\"\nBlacklist:\n              \"adsrvr.org\"\n              \"adtechus.net\"\n              \"advertising.com\"\n              \"centade.com\"\n              \"doubleclick.net\"\n              \"free-counter.co.uk\"\n              \"intellitxt.com\"\n              \"kiosked.com\"\n"},
			{name: "PreHObj", iFace: PreHObj, exp: "\nDesc:         \"pre-configured blacklisted servers\"\nDisabled:     \"false\"\nFile:         \"**Undefined**\"\nIP:           \"192.168.168.1\"\nLtype:        \"blacklisted-servers\"\nName:         \"blacklisted-servers\"\nnType:        \"preHost\"\nPrefix:       \"**Undefined**\"\nType:         \"blacklisted-servers\"\nURL:          \"**Undefined**\"\nWhitelist:\n              \"**No entries found**\"\nBlacklist:\n              \"beap.gemini.yahoo.com\"\n"},
			{name: "PreRObj", iFace: PreRObj, exp: "\nDesc:         \"pre-configured global blacklisted domains\"\nDisabled:     \"false\"\nFile:         \"**Undefined**\"\nIP:           \"0.0.0.0\"\nLtype:        \"global-blacklisted-domains\"\nName:         \"global-blacklisted-domains\"\nnType:        \"preRoot\"\nPrefix:       \"**Undefined**\"\nType:         \"global-blacklisted-domains\"\nURL:          \"**Undefined**\"\nWhitelist:\n              \"**No entries found**\"\nBlacklist:\n              \"**No entries found**\"\n"},
			{name: "URLdObj", iFace: URLdObj, exp: "\nDesc:         \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled:     \"false\"\nFile:         \"**Undefined**\"\nIP:           \"0.0.0.0\"\nLtype:        \"url\"\nName:         \"malc0de\"\nnType:        \"domn\"\nPrefix:       \"zone \"\nType:         \"domains\"\nURL:          \"http://127.0.0.1:8081/domains/domain.txt\"\nWhitelist:\n              \"**No entries found**\"\nBlacklist:\n              \"**No entries found**\"\n"},
			{name: "URLhObj", iFace: URLhObj, exp: "\nDesc:         \"Blocking mobile ad providers and some analytics providers\"\nDisabled:     \"false\"\nFile:         \"**Undefined**\"\nIP:           \"192.168.168.1\"\nLtype:        \"url\"\nName:         \"adaway\"\nnType:        \"host\"\nPrefix:       \"127.0.0.1 \"\nType:         \"hosts\"\nURL:          \"http://127.0.0.1:8081/hosts/host.txt\"\nWhitelist:\n              \"**No entries found**\"\nBlacklist:\n              \"**No entries found**\"\n"},
		}

		for _, tt := range tests {
			Convey("Testing "+tt.name+" Contenter.String()", func() {
				ct, err := c.NewContent(tt.iFace)
				So(err, ShouldBeNil)
				So(ct.String(), ShouldEqual, tt.exp)
			})
		}
	})
}

func TestIFaceString(t *testing.T) {
	tests := []struct {
		iface IFace
		name  string
		exp   string
	}{
		{name: "ExDmObj", iface: ExDmObj, exp: ExcDomns},
		{name: "ExHtObj", iface: ExHtObj, exp: ExcHosts},
		{name: "ExRtObj", iface: ExRtObj, exp: ExcRoots},
		{name: "FileObj", iface: FileObj, exp: files},
		{name: "Invalid", iface: Invalid, exp: notknown},
		{name: "PreDObj", iface: PreDObj, exp: PreDomns},
		{name: "PreHObj", iface: PreHObj, exp: PreHosts},
		{name: "PreRObj", iface: PreRObj, exp: PreRoots},
		{name: "URLdObj", iface: URLdObj, exp: urls},
		{name: "URLhObj", iface: URLhObj, exp: urls},
	}

	Convey("Testing IFace.String()", t, func() {
		for _, tt := range tests {
			Convey("with "+tt.name, func() {
				s := tt.iface.String()
				fmt.Println(s)
				So(tt.iface.String(), ShouldEqual, tt.exp)
			})
		}
	})
}

func TestMultiObjNewContent(t *testing.T) {
	Convey("Testing Multi Object NewContent()", t, func() {
		dir, err := ioutil.TempDir("/tmp", "testBlacklist")
		So(err, ShouldBeNil)
		defer os.RemoveAll(dir)

		c := NewConfig(
			Dir(dir),
			Ext("blacklist.conf"),
			FileNameFmt("%v/%v.%v.%v"),
			Logger(newLog()),
			Method("GET"),
			Prefix("address=", "addn-hosts="),
		)

		So(c.Blacklist(&CFGstatic{Cfg: CfgMimimal}), ShouldBeNil)

		tests := []struct {
			iFace IFace
			exp   string
			name  string
		}{
			{name: "ExRtObj", iFace: ExRtObj, exp: "address=/ytimg.com/#"},
			{name: "ExDmObj", iFace: ExDmObj, exp: ""},
			{name: "ExHtObj", iFace: ExHtObj, exp: ""},
			{name: "PreDObj", iFace: PreDObj, exp: "address=/awfuladvertising.com/0.0.0.0\naddress=/badadsrvr.org/0.0.0.0\naddress=/badintellitxt.com/0.0.0.0\naddress=/disgusting.unkiosked.com/0.0.0.0\naddress=/filthydoubleclick.net/0.0.0.0\naddress=/iffyfree-counter.co.uk/0.0.0.0\naddress=/nastycentade.com/0.0.0.0\naddress=/worseadtechus.net/0.0.0.0"},
			{name: "PreHObj", iFace: PreHObj, exp: "addn-hosts=/beap.gemini.yahoo.com/192.168.168.1"},
			{name: "PreRObj", iFace: PreRObj, exp: "address=/adsrvr.org/0.0.0.0\naddress=/adtechus.net/0.0.0.0\naddress=/advertising.com/0.0.0.0\naddress=/centade.com/0.0.0.0\naddress=/doubleclick.net/0.0.0.0\naddress=/free-counter.co.uk/0.0.0.0\naddress=/intellitxt.com/0.0.0.0\naddress=/kiosked.com/0.0.0.0"},
			{name: "FileObj", iFace: FileObj, exp: expFileObj},
			{name: "URLdObj", iFace: URLdObj, exp: expURLdObj},
			{name: "URLhObj", iFace: URLhObj, exp: expURLhOBJ},
		}

		for _, tt := range tests {
			Convey("Testing "+tt.name+" ProcessContent()", func() {
				ct, err := c.NewContent(tt.iFace)
				So(err, ShouldBeNil)

				switch tt.iFace {
				case ExRtObj, ExDmObj, ExHtObj, PreDObj, PreHObj, PreRObj:
					d := &dummyConfig{Env: c.Env, t: t}
					d.ProcessContent(ct)
					So(strings.Join(d.s, "\n"), ShouldEqual, tt.exp)
				default:
					So(ct.String(), ShouldEqual, tt.exp)
				}
			})
		}
	})
}

func TestProcessContent(t *testing.T) {
	Convey("Testing ProcessContent(), setting up temporary directory in /tmp", t, func() {
		dir, err := ioutil.TempDir("/tmp", "testBlacklist")
		So(err, ShouldBeNil)
		defer os.RemoveAll(dir)

		Convey("Testing ProcessContent()", func() {
			c := NewConfig(
				Dir(dir),
				Ext("blacklist.conf"),
				FileNameFmt("%v/%v.%v.%v"),
				Logger(newLog()),
				Method("GET"),
				Prefix("address=", "addn-hosts="),
			)

			tests := []struct {
				dropped   int32
				extracted int32
				kept      int32
				err       error
				exp       string
				expDexMap list
				expExcMap list
				f         string
				fdata     string
				name      string
				obj       IFace
			}{
				{
					name:      "ExRtObj",
					dropped:   0,
					extracted: 1,
					kept:      1,
					err:       nil,
					exp: `
Desc:         "pre-configured global whitelisted domains"
Disabled:     "false"
File:         "**Undefined**"
IP:           "0.0.0.0"
Ltype:        "global-whitelisted-domains"
Name:         "global-whitelisted-domains"
nType:        "excRoot"
Prefix:       "**Undefined**"
Type:         "global-whitelisted-domains"
URL:          "**Undefined**"
Whitelist:
              "ytimg.com"
Blacklist:
              "**No entries found**"
`,
					expDexMap: list{entry: entry{"ytimg.com": struct{}{}}},
					expExcMap: list{entry: entry{"ytimg.com": struct{}{}}},
					obj:       ExRtObj,
				},
				{
					name:      "ExDmObj",
					dropped:   0,
					extracted: 0,
					kept:      0,
					err:       nil,
					exp: `
Desc:         "pre-configured whitelisted subdomains"
Disabled:     "false"
File:         "**Undefined**"
IP:           "0.0.0.0"
Ltype:        "whitelisted-subdomains"
Name:         "whitelisted-subdomains"
nType:        "excDomn"
Prefix:       "**Undefined**"
Type:         "whitelisted-subdomains"
URL:          "**Undefined**"
Whitelist:
              "**No entries found**"
Blacklist:
              "**No entries found**"
`,
					expDexMap: list{RWMutex: &sync.RWMutex{}, entry: make(entry)},
					expExcMap: list{RWMutex: &sync.RWMutex{}, entry: make(entry)},
					obj:       ExDmObj,
				},
				{
					name:      "ExHtObj",
					dropped:   0,
					extracted: 0,
					kept:      0,
					err:       nil,
					exp: `
Desc:         "pre-configured whitelisted servers"
Disabled:     "false"
File:         "**Undefined**"
IP:           "192.168.168.1"
Ltype:        "whitelisted-servers"
Name:         "whitelisted-servers"
nType:        "excHost"
Prefix:       "**Undefined**"
Type:         "whitelisted-servers"
URL:          "**Undefined**"
Whitelist:
              "**No entries found**"
Blacklist:
              "**No entries found**"
`,
					expDexMap: list{RWMutex: &sync.RWMutex{}, entry: make(entry)},
					expExcMap: list{RWMutex: &sync.RWMutex{}, entry: make(entry)},
					obj:       ExHtObj,
				},
				{
					name:      "PreDObj",
					dropped:   0,
					extracted: 8,
					kept:      8,
					err:       nil,
					exp: `
Desc:         "pre-configured blacklisted subdomains"
Disabled:     "false"
File:         "**Undefined**"
IP:           "0.0.0.0"
Ltype:        "blacklisted-subdomains"
Name:         "blacklisted-subdomains"
nType:        "preDomn"
Prefix:       "**Undefined**"
Type:         "blacklisted-subdomains"
URL:          "**Undefined**"
Whitelist:
              "**No entries found**"
Blacklist:
              "adsrvr.org"
              "adtechus.net"
              "advertising.com"
              "centade.com"
              "doubleclick.net"
              "free-counter.co.uk"
              "intellitxt.com"
              "kiosked.com"
`,
					expDexMap: list{
						entry: entry{
							"adsrvr.org":         struct{}{},
							"adtechus.net":       struct{}{},
							"advertising.com":    struct{}{},
							"centade.com":        struct{}{},
							"doubleclick.net":    struct{}{},
							"free-counter.co.uk": struct{}{},
							"intellitxt.com":     struct{}{},
							"kiosked.com":        struct{}{},
						},
					},
					expExcMap: list{entry: entry{"ytimg.com": struct{}{}}},
					f:         dir + "/domains.blacklisted-subdomains.blacklist.conf",
					fdata: `address=/awfuladvertising.com/0.0.0.0
address=/badadsrvr.org/0.0.0.0
address=/badintellitxt.com/0.0.0.0
address=/disgusting.unkiosked.com/0.0.0.0
address=/filthydoubleclick.net/0.0.0.0
address=/iffyfree-counter.co.uk/0.0.0.0
address=/nastycentade.com/0.0.0.0
address=/worseadtechus.net/0.0.0.0
`,
					obj: PreDObj,
				},
				{
					name:      "PreHObj",
					dropped:   0,
					extracted: 1,
					kept:      1,
					err:       nil,
					exp: `
Desc:         "pre-configured blacklisted servers"
Disabled:     "false"
File:         "**Undefined**"
IP:           "192.168.168.1"
Ltype:        "blacklisted-servers"
Name:         "blacklisted-servers"
nType:        "preHost"
Prefix:       "**Undefined**"
Type:         "blacklisted-servers"
URL:          "**Undefined**"
Whitelist:
              "**No entries found**"
Blacklist:
              "beap.gemini.yahoo.com"
`,
					expDexMap: list{entry: entry{"ytimg.com": struct{}{}}},
					expExcMap: list{entry: entry{"ytimg.com": struct{}{}}},
					f:         dir + "/hosts.blacklisted-servers.blacklist.conf",
					fdata:     "addn-hosts=/beap.gemini.yahoo.com/192.168.168.1\n",
					obj:       PreHObj,
				},
				{
					name:      "PreRObj",
					dropped:   0,
					extracted: 8,
					kept:      8,
					err:       nil,
					exp: `
Desc:         "pre-configured global blacklisted domains"
Disabled:     "false"
File:         "**Undefined**"
IP:           "0.0.0.0"
Ltype:        "global-blacklisted-domains"
Name:         "global-blacklisted-domains"
nType:        "preRoot"
Prefix:       "**Undefined**"
Type:         "global-blacklisted-domains"
URL:          "**Undefined**"
Whitelist:
              "**No entries found**"
Blacklist:
              "adsrvr.org"
              "adtechus.net"
              "advertising.com"
              "centade.com"
              "doubleclick.net"
              "free-counter.co.uk"
              "intellitxt.com"
              "kiosked.com"
`,
					expDexMap: list{entry: entry{}},
					expExcMap: list{
						entry: entry{
							"adsrvr.org":         struct{}{},
							"adtechus.net":       struct{}{},
							"advertising.com":    struct{}{},
							"centade.com":        struct{}{},
							"doubleclick.net":    struct{}{},
							"free-counter.co.uk": struct{}{},
							"intellitxt.com":     struct{}{},
							"kiosked.com":        struct{}{},
						},
					},
					obj: PreRObj,
				},
				{
					name:      "FileObj",
					dropped:   2,
					extracted: 23,
					kept:      21,
					err:       fmt.Errorf("open %v/hosts./tasty.blacklist.conf: no such file or directory", dir),
					exp:       filesMin,
					expDexMap: list{
						entry: entry{
							"cw.bad.ultraadverts.site.eu": struct{}{},
							"really.bad.phishing.site.ru": struct{}{},
						},
					},
					expExcMap: list{entry: entry{"ytimg.com": struct{}{}}},
					f:         dir + "/hosts.tasty.blacklist.conf",
					fdata: `addn-hosts=/0.really.bad.phishing.site.ru/10.10.10.10
addn-hosts=/cw.bad.ultraadverts.site.eu/10.10.10.10
addn-hosts=/really.bad.phishing.site.ru/10.10.10.10
`,
					obj: FileObj,
				},
			}

			So(c.Blacklist(&CFGstatic{Cfg: CfgMimimal}), ShouldBeNil)

			for _, tt := range tests {
				Convey("Testing "+tt.name+" ProcessContent()", func() {
					var (
						ct    Contenter
						objex []IFace
					)

					switch tt.obj {
					case FileObj, URLdObj, URLhObj:
						objex = []IFace{
							PreRObj,
							PreDObj,
							PreHObj,
							ExRtObj,
							ExDmObj,
							ExHtObj,
							tt.obj,
						}
					default:
						objex = []IFace{tt.obj}
					}

					var g errgroup.Group
					g.Go(
						func() (err error) {
							for _, o := range objex {
								ct, _ = c.NewContent(o)
								err = c.ProcessContent(ct)
							}
							return err
						})

					if g.Wait() != nil {
						Convey("Testing "+tt.name+" ProcessContent().Error():", func() {
							Convey("Error should match expected", func() {
								So(err, ShouldResemble, tt.err)
							})
						})
					}

					dropped, extracted, kept := c.GetTotalStats()

					Convey("Dropped entries should match", func() {
						So(dropped, ShouldEqual, tt.dropped)
					})

					Convey("Extracted entries should match", func() {
						So(extracted, ShouldEqual, tt.extracted)
					})

					Convey("Kept entries should match", func() {
						So(kept, ShouldEqual, tt.kept)
					})

					switch tt.f {
					case "":
						Convey("Testing "+tt.name+" ProcessContent(): Dex map should match expected", func() {
							So(c.Dex.entry, ShouldResemble, tt.expDexMap.entry)
						})

						Convey("Testing "+tt.name+" ProcessContent(): Exc map should match expected", func() {
							So(c.Exc.entry, ShouldResemble, tt.expExcMap.entry)
						})

						Convey("Testing "+tt.name+" ProcessContent(): ct should match expected", func() {
							So(ct.String(), ShouldEqual, tt.exp)
						})
					default:
						reader, err := GetFile(tt.f)
						So(err, ShouldBeNil)

						act, err := ioutil.ReadAll(reader)
						So(err, ShouldBeNil)

						Convey("Testing "+tt.name+" ProcessContent(): file data should match expected", func() {
							So(string(act), ShouldEqual, tt.fdata)
						})
					}
				})
			}
		})
	})
}

func TestProcessZeroContent(t *testing.T) {
	Convey("Testing ProcessZeroContent()", t, func() {
		dir, err := ioutil.TempDir("/tmp", "testBlacklist")
		So(err, ShouldBeNil)
		defer os.RemoveAll(dir)
		c := NewConfig(
			Dir(dir),
			Ext("blacklist.conf"),
			FileNameFmt("%v/%v.%v.%v"),
			Logger(newLog()),
			Method("GET"),
			Prefix("address=", "addn-hosts="),
		)

		err = c.Blacklist(&CFGstatic{Cfg: cfgRedundant})
		So(err, ShouldBeNil)

		for _, o := range []IFace{ExRtObj, FileObj} {
			ct, err := c.NewContent(o)
			So(err, ShouldBeNil)

			err = c.ProcessContent(ct)
			So(err, ShouldBeNil)
		}

		dropped, extracted, kept := c.GetTotalStats()

		Convey("Dropped entries should match", func() {
			So(dropped, ShouldEqual, 1)
		})

		Convey("Extracted entries should match", func() {
			So(extracted, ShouldEqual, 2)
		})

		Convey("Kept entries should match", func() {
			So(kept, ShouldEqual, 1)
		})
	})
}

func TestProcessBadFile(t *testing.T) {
	Convey("Testing ProcessBadFile()", t, func() {
		dir, err := ioutil.TempDir("/tmp", "testBlacklist")
		So(err, ShouldBeNil)
		defer os.RemoveAll(dir)
		c := NewConfig(
			Dir("/:~/"),
			Ext("blacklist.conf"),
			FileNameFmt("%v/%v.%v.%v"),
			Logger(newLog()),
			Method("GET"),
			Prefix("address=", "addn-hosts="),
		)

		err = c.Blacklist(&CFGstatic{Cfg: CfgMimimal})
		So(err, ShouldBeNil)

		ct, err := c.NewContent(FileObj)
		So(err, ShouldBeNil)

		err = c.ProcessContent(ct)
		So(err.Error(), ShouldEqual, "open /:~//hosts.tasty.blacklist.conf: no such file or directory")
	})
}

func TestWriteFile(t *testing.T) {
	Convey("Testing WriteFile()", t, func() {
		tests := []struct {
			data  io.Reader
			dir   string
			fname string
			ok    bool
			want  string
		}{
			{
				data:  strings.NewReader("The rest is history!"),
				dir:   "/tmp",
				fname: "Test.util.writeFile",
				ok:    true,
				want:  "",
			},
			{
				data:  bytes.NewBuffer([]byte{84, 104, 101, 32, 114, 101, 115, 116, 32, 105, 115, 32, 104, 105, 115, 116, 111, 114, 121, 33}),
				dir:   "/tmp",
				fname: "Test.util.writeFile",
				ok:    true,
				want:  "",
			},
			{
				data:  bytes.NewBufferString("This shouldn't be written!"),
				dir:   "",
				fname: "/",
				ok:    false,
				want:  "open /: is a directory",
			},
		}

		for _, tt := range tests {
			switch tt.ok {
			case true:
				f, err := ioutil.TempFile(tt.dir, tt.fname)
				So(err, ShouldBeNil)
				b := &bList{
					file: f.Name(),
					r:    tt.data,
					size: 20,
				}
				So(b.writeFile(), ShouldBeNil)
				os.Remove(f.Name())

			default:
				b := &bList{
					file: tt.dir + tt.fname,
					r:    tt.data,
					size: 20,
				}
				So(b.writeFile().Error(), ShouldResemble, tt.want)
			}
		}
	})
}

var (
	// Cfg contains a valid full EdgeOS blacklist configuration
	Cfg = `blacklist {
    disabled false
    dns-redirect-ip 0.0.0.0
    domains {
        dns-redirect-ip 192.1.1.1
		exclude adinfuse.com
        include adsrvr.org
        include adtechus.net
        include advertising.com
        include centade.com
        include doubleclick.net
        include free-counter.co.uk
        include intellitxt.com
        include kiosked.com
        source malc0de {
            description "List of zones serving malicious executables observed by malc0de.com/database/"
            prefix "zone "
            url http://malc0de.com/bl/ZONES
        }
    }
    exclude 122.2o7.net
    exclude 1e100.net
    exclude adobedtm.com
    exclude akamai.net
    exclude amazon.com
    exclude amazonaws.com
    exclude apple.com
    exclude ask.com
    exclude avast.com
    exclude bitdefender.com
    exclude cdn.visiblemeasures.com
    exclude cloudfront.net
    exclude coremetrics.com
    exclude edgesuite.net
    exclude freedns.afraid.org
    exclude github.com
    exclude githubusercontent.com
    exclude google.com
    exclude googleadservices.com
    exclude googleapis.com
    exclude googleusercontent.com
    exclude gstatic.com
    exclude gvt1.com
    exclude gvt1.net
    exclude hb.disney.go.com
    exclude hp.com
    exclude hulu.com
    exclude images-amazon.com
	exclude jumptap.com
    exclude msdn.com
    exclude paypal.com
    exclude rackcdn.com
    exclude schema.org
    exclude skype.com
    exclude smacargo.com
    exclude sourceforge.net
    exclude ssl-on9.com
    exclude ssl-on9.net
    exclude static.chartbeat.com
    exclude storage.googleapis.com
	exclude usemaxserver.de
    exclude windows.net
    exclude yimg.com
    exclude ytimg.com
    hosts {
		exclude wv.inner-active.mobi
        include beap.gemini.yahoo.com
        source adaway {
            description "Blocking mobile ad providers and some analytics providers"
			dns-redirect-ip 192.168.168.1
            prefix "127.0.0.1 "
            url http://adaway.org/hosts.txt
        }
				source tasty {
						description "File source"
						dns-redirect-ip 0.0.0.0
						file ../../internal/testdata/blist.hosts.src
				}
    }
}`

	cfgRedundant = `blacklist {
	disabled false
	dns-redirect-ip 0.0.0.0
	domains {
		source tasty {
			description "File source"
			dns-redirect-ip 10.10.10.10
			file ../../internal/testdata/blist.nohosts.src
	}
	}
	exclude ytimg.com
}`
	// CfgMimimal contains a valid minimal EdgeOS blacklist configuration
	CfgMimimal = `blacklist {
	disabled false
	dns-redirect-ip 0.0.0.0
	domains {
			include badadsrvr.org
			include worseadtechus.net
			include awfuladvertising.com
			include nastycentade.com
			include filthydoubleclick.net
			include iffyfree-counter.co.uk
			include badintellitxt.com
			include disgusting.unkiosked.com
			source malc0de {
					description "List of zones serving malicious executables observed by malc0de.com/database/"
					prefix "zone "
					url http://malc0de.com/bl/ZONES
			}
	}
	exclude ytimg.com
	include adsrvr.org
	include adtechus.net
	include advertising.com
	include centade.com
	include doubleclick.net
	include free-counter.co.uk
	include intellitxt.com
	include kiosked.com
	hosts {
			dns-redirect-ip 192.168.168.1
			include beap.gemini.yahoo.com
			source tasty {
					description "File source"
					dns-redirect-ip 10.10.10.10
					file ../../internal/testdata/blist.hosts.src
			}
			source adaway {
          description "Blocking mobile ad providers and some analytics providers"
			    dns-redirect-ip 192.168.168.1
          prefix "127.0.0.1 "
          url http://adaway.org/hosts.txt
      }
	}
}`

	// testallCfg contains a valid full EdgeOS blacklist configuration with localized URLs
	testallCfg = `blacklist {
	disabled false
	dns-redirect-ip 0.0.0.0
	domains {
			dns-redirect-ip 0.0.0.0
			include adsrvr.org
			include adtechus.net
			include advertising.com
			include centade.com
			include doubleclick.net
			include free-counter.co.uk
			include intellitxt.com
			include kiosked.com
			source malc0de {
					description "List of zones serving malicious executables observed by malc0de.com/database/"
					prefix "zone "
					url http://127.0.0.1:8081/domains/domain.txt
			}
	}
	exclude 122.2o7.net
	exclude 1e100.net
	exclude adobedtm.com
	exclude akamai.net
	exclude amazon.com
	exclude amazonaws.com
	exclude apple.com
	exclude ask.com
	exclude avast.com
	exclude bitdefender.com
	exclude cdn.visiblemeasures.com
	exclude cloudfront.net
	exclude coremetrics.com
	exclude edgesuite.net
	exclude freedns.afraid.org
	exclude github.com
	exclude githubusercontent.com
	exclude google.com
	exclude googleadservices.com
	exclude googleapis.com
	exclude googleusercontent.com
	exclude gstatic.com
	exclude gvt1.com
	exclude gvt1.net
	exclude hb.disney.go.com
	exclude hp.com
	exclude hulu.com
	exclude images-amazon.com
	exclude msdn.com
	exclude paypal.com
	exclude rackcdn.com
	exclude schema.org
	exclude skype.com
	exclude smacargo.com
	exclude sourceforge.net
	exclude ssl-on9.com
	exclude ssl-on9.net
	exclude static.chartbeat.com
	exclude storage.googleapis.com
	exclude windows.net
	exclude yimg.com
	exclude ytimg.com
	hosts {
			dns-redirect-ip 192.168.168.1
			include beap.gemini.yahoo.com
			source adaway {
					description "Blocking mobile ad providers and some analytics providers"
					prefix "127.0.0.1 "
					url http://127.0.0.1:8081/hosts/host.txt
			}
			source tasty {
					description "File source"
					dns-redirect-ip 0.0.0.0
					file ../../internal/testdata/blist.hosts.src
			}
	}
}`

	hostsContent = `addn-hosts=/a.applovin.com/192.168.168.1
addn-hosts=/a.glcdn.co/192.168.168.1
addn-hosts=/a.vserv.mobi/192.168.168.1
addn-hosts=/ad.leadboltapps.net/192.168.168.1
addn-hosts=/ad.madvertise.de/192.168.168.1
addn-hosts=/ad.where.com/192.168.168.1
addn-hosts=/ad1.adinfuse.com/192.168.168.1
addn-hosts=/ad2.adinfuse.com/192.168.168.1
addn-hosts=/adcontent.saymedia.com/192.168.168.1
addn-hosts=/adinfuse.com/192.168.168.1
addn-hosts=/admicro1.vcmedia.vn/192.168.168.1
addn-hosts=/admicro2.vcmedia.vn/192.168.168.1
addn-hosts=/admin.vserv.mobi/192.168.168.1
addn-hosts=/ads.adiquity.com/192.168.168.1
addn-hosts=/ads.admarvel.com/192.168.168.1
addn-hosts=/ads.admoda.com/192.168.168.1
addn-hosts=/ads.celtra.com/192.168.168.1
addn-hosts=/ads.flurry.com/192.168.168.1
addn-hosts=/ads.matomymobile.com/192.168.168.1
addn-hosts=/ads.mobgold.com/192.168.168.1
addn-hosts=/ads.mobilityware.com/192.168.168.1
addn-hosts=/ads.mopub.com/192.168.168.1
addn-hosts=/ads.n-ws.org/192.168.168.1
addn-hosts=/ads.ookla.com/192.168.168.1
addn-hosts=/ads.saymedia.com/192.168.168.1
addn-hosts=/ads.smartdevicemedia.com/192.168.168.1
addn-hosts=/ads.srcxad.net/192.168.168.1
addn-hosts=/ads.vserv.mobi/192.168.168.1
addn-hosts=/ads2.mediaarmor.com/192.168.168.1
addn-hosts=/adserver.ubiyoo.com/192.168.168.1
addn-hosts=/adultmoda.com/192.168.168.1
addn-hosts=/android-sdk31.transpera.com/192.168.168.1
addn-hosts=/android.bcfads.com/192.168.168.1
addn-hosts=/api.airpush.com/192.168.168.1
addn-hosts=/api.analytics.omgpop.com/192.168.168.1
addn-hosts=/api.yp.com/192.168.168.1
addn-hosts=/apps.buzzcity.net/192.168.168.1
addn-hosts=/apps.mobilityware.com/192.168.168.1
addn-hosts=/as.adfonic.net/192.168.168.1
addn-hosts=/asotrack1.fluentmobile.com/192.168.168.1
addn-hosts=/assets.cntdy.mobi/192.168.168.1
addn-hosts=/atti.velti.com/192.168.168.1
addn-hosts=/b.scorecardresearch.com/192.168.168.1
addn-hosts=/banners.bigmobileads.com/192.168.168.1
addn-hosts=/bigmobileads.com/192.168.168.1
addn-hosts=/bo.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq01.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq02.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq03.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq04.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq05.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq06.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq07.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq08.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq09.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq10.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq11.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq12.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq13.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq14.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq15.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq16.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq17.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq18.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq19.jumptap.com/192.168.168.1
addn-hosts=/bos-tapreq20.jumptap.com/192.168.168.1
addn-hosts=/c.vrvm.com/192.168.168.1
addn-hosts=/c.vserv.mobi/192.168.168.1
addn-hosts=/c753738.r38.cf2.rackcdn.com/192.168.168.1
addn-hosts=/cache-ssl.celtra.com/192.168.168.1
addn-hosts=/cache.celtra.com/192.168.168.1
addn-hosts=/cdn.celtra.com/192.168.168.1
addn-hosts=/cdn.nearbyad.com/192.168.168.1
addn-hosts=/cdn.trafficforce.com/192.168.168.1
addn-hosts=/cdn.us.goldspotmedia.com/192.168.168.1
addn-hosts=/cdn.vdopia.com/192.168.168.1
addn-hosts=/cdn1.crispadvertising.com/192.168.168.1
addn-hosts=/cdn1.inner-active.mobi/192.168.168.1
addn-hosts=/cdn2.crispadvertising.com/192.168.168.1
addn-hosts=/click.buzzcity.net/192.168.168.1
addn-hosts=/creative1cdn.mobfox.com/192.168.168.1
addn-hosts=/d.applovin.com/192.168.168.1
addn-hosts=/d2bgg7rjywcwsy.cloudfront.net/192.168.168.1
addn-hosts=/d3anogn3pbtk4v.cloudfront.net/192.168.168.1
addn-hosts=/d3oltyb66oj2v8.cloudfront.net/192.168.168.1
addn-hosts=/edge.reporo.net/192.168.168.1
addn-hosts=/ftpcontent.worldnow.com/192.168.168.1
addn-hosts=/funnel0.adinfuse.com/192.168.168.1
addn-hosts=/gemini.yahoo.com/192.168.168.1
addn-hosts=/go.adinfuse.com/192.168.168.1
addn-hosts=/go.mobpartner.mobi/192.168.168.1
addn-hosts=/go.vrvm.com/192.168.168.1
addn-hosts=/gsmtop.net/192.168.168.1
addn-hosts=/gts-ads.twistbox.com/192.168.168.1
addn-hosts=/hhbeksrcw5d9e.pflexads.com/192.168.168.1
addn-hosts=/hybl9bazbc35.pflexads.com/192.168.168.1
addn-hosts=/i.jumptap.com/192.168.168.1
addn-hosts=/i.tapit.com/192.168.168.1
addn-hosts=/images.millennialmedia.com/192.168.168.1
addn-hosts=/images.mpression.net/192.168.168.1
addn-hosts=/img.ads.huntmad.com/192.168.168.1
addn-hosts=/img.ads.mobilefuse.net/192.168.168.1
addn-hosts=/img.ads.mocean.mobi/192.168.168.1
addn-hosts=/img.ads.mojiva.com/192.168.168.1
addn-hosts=/img.ads.taptapnetworks.com/192.168.168.1
addn-hosts=/intouch.adinfuse.com/192.168.168.1
addn-hosts=/lb.usemaxserver.de/192.168.168.1
addn-hosts=/m.adsymptotic.com/192.168.168.1
addn-hosts=/m2m1.inner-active.mobi/192.168.168.1
addn-hosts=/media.mobpartner.mobi/192.168.168.1
addn-hosts=/medrx.sensis.com.au/192.168.168.1
addn-hosts=/mobile.banzai.it/192.168.168.1
addn-hosts=/mobiledl.adboe.com/192.168.168.1
addn-hosts=/mobpartner.mobi/192.168.168.1
addn-hosts=/mwc.velti.com/192.168.168.1
addn-hosts=/netdna.reporo.net/192.168.168.1
addn-hosts=/oasc04012.247realmedia.com/192.168.168.1
addn-hosts=/orange-fr.adinfuse.com/192.168.168.1
addn-hosts=/orangeuk-mc.adinfuse.com/192.168.168.1
addn-hosts=/orencia.pflexads.com/192.168.168.1
addn-hosts=/pdn.applovin.com/192.168.168.1
addn-hosts=/r.edge.inmobicdn.net/192.168.168.1
addn-hosts=/r.mobpartner.mobi/192.168.168.1
addn-hosts=/req.appads.com/192.168.168.1
addn-hosts=/rs-staticart.ybcdn.net/192.168.168.1
addn-hosts=/ru.velti.com/192.168.168.1
addn-hosts=/s0.2mdn.net/192.168.168.1
addn-hosts=/s3.phluant.com/192.168.168.1
addn-hosts=/sf.vserv.mobi/192.168.168.1
addn-hosts=/show.buzzcity.net/192.168.168.1
addn-hosts=/sky-connect.adinfuse.com/192.168.168.1
addn-hosts=/sky.adinfuse.com/192.168.168.1
addn-hosts=/static.cdn.gtsmobi.com/192.168.168.1
addn-hosts=/static.estebull.com/192.168.168.1
addn-hosts=/stats.pflexads.com/192.168.168.1
addn-hosts=/track.celtra.com/192.168.168.1
addn-hosts=/tracking.klickthru.com/192.168.168.1
addn-hosts=/uk-ad2.adinfuse.com/192.168.168.1
addn-hosts=/uk-go.adinfuse.com/192.168.168.1
addn-hosts=/web63.jumptap.com/192.168.168.1
addn-hosts=/web64.jumptap.com/192.168.168.1
addn-hosts=/web65.jumptap.com/192.168.168.1
addn-hosts=/wv.inner-active.mobi/192.168.168.1
addn-hosts=/www.eltrafiko.com/192.168.168.1
addn-hosts=/www.mmnetwork.mobi/192.168.168.1
addn-hosts=/www.pflexads.com/192.168.168.1
addn-hosts=/www.roastfiles2017.com/192.168.168.1
addn-hosts=/wwww.adleads.com/192.168.168.1`

	domainsContent = "address=/192-168-0-255.com/192.1.1.1\naddress=/asi-37.fr/192.1.1.1\naddress=/bagbackpack.com/192.1.1.1\naddress=/bitmeyenkartusistanbul.com/192.1.1.1\naddress=/byxon.com/192.1.1.1\naddress=/img001.com/192.1.1.1\naddress=/loadto.net/192.1.1.1\naddress=/roastfiles2017.com/192.1.1.1"

	// domainsPreContent = "address=/adsrvr.org/192.1.1.1\naddress=/adtechus.net/192.1.1.1\naddress=/advertising.com/192.1.1.1\naddress=/centade.com/192.1.1.1\naddress=/doubleclick.net/192.1.1.1\naddress=/free-counter.co.uk/192.1.1.1\naddress=/intellitxt.com/192.1.1.1\naddress=/kiosked.com/192.1.1.1\n"

	// expPreGetAll = "[\nDesc:\t \"pre-configured blacklisted subdomains\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"blacklisted-subdomains\"\nName:\t \"blacklisted-subdomains\"\nnType:\t \"preDomn\"\nPrefix:\t \"\"\nType:\t \"blacklisted-subdomains\"\nURL:\t \"\"\n \nDesc:\t \"pre-configured blacklisted servers\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"blacklisted-servers\"\nName:\t \"blacklisted-servers\"\nnType:\t \"preHost\"\nPrefix:\t \"\"\nType:\t \"blacklisted-servers\"\nURL:\t \"\"\n]"

	// expAll = "[\nDesc:\t \"pre-configured blacklisted subdomains\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"blacklisted-subdomains\"\nName:\t \"blacklisted-subdomains\"\nnType:\t \"preDomn\"\nPrefix:\t \"\"\nType:\t \"blacklisted-subdomains\"\nURL:\t \"\"\n \nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"malc0de\"\nnType:\t \"domn\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://127.0.0.1:8081/domains/domain.txt\"\n \nDesc:\t \"pre-configured blacklisted servers\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"blacklisted-servers\"\nName:\t \"blacklisted-servers\"\nnType:\t \"preHost\"\nPrefix:\t \"\"\nType:\t \"blacklisted-servers\"\nURL:\t \"\"\n \nDesc:\t \"Blocking mobile ad providers and some analytics providers\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"adaway\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://127.0.0.1:8081/hosts/host.txt\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../../internal/testdata/blist.hosts.src\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n]"

	expFileObj = `
Desc:         "File source"
Disabled:     "false"
File:         "../../internal/testdata/blist.hosts.src"
IP:           "10.10.10.10"
Ltype:        "file"
Name:         "tasty"
nType:        "host"
Prefix:       "**Undefined**"
Type:         "hosts"
URL:          "**Undefined**"
Whitelist:
              "**No entries found**"
Blacklist:
              "**No entries found**"
`

	expURLdObj = `
Desc:         "List of zones serving malicious executables observed by malc0de.com/database/"
Disabled:     "false"
File:         "**Undefined**"
IP:           "0.0.0.0"
Ltype:        "url"
Name:         "malc0de"
nType:        "domn"
Prefix:       "zone "
Type:         "domains"
URL:          "http://malc0de.com/bl/ZONES"
Whitelist:
              "**No entries found**"
Blacklist:
              "**No entries found**"
`

	expURLhOBJ = `
Desc:         "Blocking mobile ad providers and some analytics providers"
Disabled:     "false"
File:         "**Undefined**"
IP:           "192.168.168.1"
Ltype:        "url"
Name:         "adaway"
nType:        "host"
Prefix:       "127.0.0.1 "
Type:         "hosts"
URL:          "http://adaway.org/hosts.txt"
Whitelist:
              "**No entries found**"
Blacklist:
              "**No entries found**"
`

	filesMin = "[\nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../../internal/testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../../internal/testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"/tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n]"

	excRootContent = "address=/122.2o7.net/#\naddress=/1e100.net/#\naddress=/adobedtm.com/#\naddress=/akamai.net/#\naddress=/amazon.com/#\naddress=/amazonaws.com/#\naddress=/apple.com/#\naddress=/ask.com/#\naddress=/avast.com/#\naddress=/bitdefender.com/#\naddress=/cdn.visiblemeasures.com/#\naddress=/cloudfront.net/#\naddress=/coremetrics.com/#\naddress=/edgesuite.net/#\naddress=/freedns.afraid.org/#\naddress=/github.com/#\naddress=/githubusercontent.com/#\naddress=/google.com/#\naddress=/googleadservices.com/#\naddress=/googleapis.com/#\naddress=/googleusercontent.com/#\naddress=/gstatic.com/#\naddress=/gvt1.com/#\naddress=/gvt1.net/#\naddress=/hb.disney.go.com/#\naddress=/hp.com/#\naddress=/hulu.com/#\naddress=/images-amazon.com/#\naddress=/jumptap.com/#\naddress=/msdn.com/#\naddress=/paypal.com/#\naddress=/rackcdn.com/#\naddress=/schema.org/#\naddress=/skype.com/#\naddress=/smacargo.com/#\naddress=/sourceforge.net/#\naddress=/ssl-on9.com/#\naddress=/ssl-on9.net/#\naddress=/static.chartbeat.com/#\naddress=/storage.googleapis.com/#\naddress=/usemaxserver.de/#\naddress=/windows.net/#\naddress=/yimg.com/#\naddress=/ytimg.com/#"

	testCfg = `blacklist {
	disabled false
	dns-redirect-ip 0.0.0.0
	domains {
			include adsrvr.org
			include adtechus.net
			include advertising.com
			include centade.com
			include doubleclick.net
			include free-counter.co.uk
			include intellitxt.com
			include kiosked.com
	}
	exclude ytimg.com
	hosts {
		dns-redirect-ip 192.168.168.1
		include beap.gemini.yahoo.com
		source tasty {
			description "File source"
			dns-redirect-ip 10.10.10.10
			file /:~/=../../internal/testdata/blist.hosts.src
		}
	}
}`
)

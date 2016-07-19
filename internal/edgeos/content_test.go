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
)

type dummyConfig struct {
	s []string
	t *testing.T
}

func (d *dummyConfig) ProcessContent(cts ...Contenter) {
	for _, ct := range cts {
		o := ct.GetList().x
		for _, src := range o {
			b, _ := ioutil.ReadAll(src.process().r)
			d.s = append(d.s, strings.TrimSuffix(string(b), "\n"))
		}
	}
}

func TestConfigProcessContent(t *testing.T) {
	Convey("Testing ProcessContent()", t, func() {
		newCfg := func() *Config {
			return NewConfig(
				API("/bin/cli-shell-api"),
				Arch(runtime.GOARCH),
				Bash("/bin/bash"),
				Cores(runtime.NumCPU()),
				Dir("/:~/"),
				DNSsvc("service dnsmasq restart"),
				Ext("blacklist.conf"),
				FileNameFmt("%v/%v.%v.%v"),
				InCLI("inSession"),
				Level("service dns forwarding"),
				Method("GET"),
				Nodes([]string{domains, hosts}),
				Prefix("address="),
				LTypes([]string{files, PreDomns, PreHosts, urls}),
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
				cfg:    testCfg,
				ct:     FileObj,
				err:    fmt.Errorf("open /:~/=../testdata/blist.hosts.src: no such file or directory\nopen /:~//hosts.tasty.blacklist.conf: no such file or directory"),
				expErr: true,
				name:   "File",
			},
		}
		for _, tt := range tests {
			So(tt.c.ReadCfg(&CFGstatic{Cfg: tt.cfg}), ShouldBeNil)

			obj, err := tt.c.NewContent(tt.ct)
			So(err, ShouldBeNil)
			if err := tt.c.ProcessContent(obj); (err != nil) == tt.expErr {
				So(err.Error(), ShouldEqual, tt.err.Error())
			}
		}

		Convey("Testing ProcessContent() if no arguments ", func() {
			So(newCfg().ProcessContent().Error(), ShouldNotBeNil)
		})
	})
}

func TestNewContent(t *testing.T) {
	Convey("Testing NewContent()", t, func() {
		tests := []struct {
			err       error
			exp       string
			fail      bool
			i         int
			leaf      string
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
				name:  ExcRoots,
				obj:   ExRtObj,
				pos:   0,
			},
			{
				i:     1,
				exp:   "address=/adinfuse.com/192.1.1.1",
				fail:  false,
				ltype: ExcDomns,
				name:  ExcDomns,
				obj:   ExDmObj,
				pos:   0,
			},
			{
				i:     1,
				exp:   "address=/wv.inner-active.mobi/0.0.0.0",
				fail:  false,
				ltype: ExcHosts,
				name:  ExcHosts,
				obj:   ExHtObj,
				pos:   0,
			},
			{
				i:     1,
				exp:   "",
				fail:  false,
				ltype: ExcRoots,
				name:  "exclusive root domains",
				obj:   ExRtObj,
				pos:   -1,
			},
			{
				i:     1,
				exp:   "",
				fail:  false,
				ltype: ExcDomns,
				name:  "exclusive domains",
				obj:   ExDmObj,
				pos:   -1,
			},
			{
				i:     1,
				exp:   "",
				fail:  false,
				ltype: ExcHosts,
				name:  "exclusive hosts",
				obj:   ExHtObj,
				pos:   -1,
			},
			{
				i:     1,
				exp:   "address=/adsrvr.org/192.1.1.1\naddress=/adtechus.net/192.1.1.1\naddress=/advertising.com/192.1.1.1\naddress=/centade.com/192.1.1.1\naddress=/doubleclick.net/192.1.1.1\naddress=/free-counter.co.uk/192.1.1.1\naddress=/intellitxt.com/192.1.1.1\naddress=/kiosked.com/192.1.1.1",
				fail:  false,
				ltype: PreDomns,
				name:  fmt.Sprintf("includes.[8]"),
				obj:   PreDObj,
				pos:   0,
			},
			{
				i:     1,
				exp:   "",
				fail:  false,
				ltype: PreDomns,
				name:  "pre",
				obj:   PreDObj,
				pos:   -1,
			},
			{
				i:     1,
				exp:   "address=/beap.gemini.yahoo.com/0.0.0.0",
				fail:  false,
				ltype: PreHosts,
				name:  fmt.Sprintf("includes.[1]"),
				obj:   PreHObj,
				pos:   0,
			},
			{
				i:     1,
				exp:   "",
				fail:  false,
				ltype: PreHosts,
				name:  "pre",
				obj:   PreHObj,
				pos:   -1,
			},
			{
				i:     1,
				exp:   "address=/cw.bad.ultraadverts.site.eu/0.0.0.0\naddress=/really.bad.phishing.site.ru/0.0.0.0",
				fail:  false,
				ltype: files,
				name:  "tasty",
				obj:   FileObj,
				pos:   0,
			},
			{
				i:     1,
				exp:   "",
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
				exp:       "",
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
				i:         1,
				exp:       "",
				fail:      false,
				ltype:     urls,
				name:      "zadway",
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
				i:    0,
				err:  errors.New("Invalid interface requested"),
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
			DNSsvc("service dnsmasq restart"),
			Ext("blacklist.conf"),
			FileNameFmt("%v/%v.%v.%v"),
			InCLI("inSession"),
			Level("service dns forwarding"),
			Method("GET"),
			Nodes([]string{domains, hosts}),
			Prefix("address="),
			LTypes([]string{files, PreDomns, PreHosts, urls}),
			Timeout(30*time.Second),
			WCard(Wildcard{Node: "*s", Name: "*"}),
		)

		So(c.ReadCfg(&CFGstatic{Cfg: Cfg}), ShouldBeNil)

		for _, tt := range tests {
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

				d := &dummyConfig{t: t}
				d.ProcessContent(objs)
				So(strings.Join(d.s, "\n"), ShouldEqual, tt.exp)

				objs.SetURL(tt.name, tt.name)
				So(objs.Find(tt.name), ShouldEqual, tt.pos)
				So(objs.Len(), ShouldEqual, tt.i)

			default:
				So(err.Error(), ShouldEqual, tt.err.Error())
			}
		}
	})
}

func TestGetAllContent(t *testing.T) {
	Convey("Testing GetAllContent()", t, func() {
		c := NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
			FileNameFmt("%v/%v.%v.%v"),
			Method("GET"),
			Nodes([]string{"domains", "hosts"}),
			Prefix("address="),
			LTypes([]string{PreDomns, PreHosts, files, urls}),
		)

		So(c.ReadCfg(&CFGstatic{Cfg: testallCfg}), ShouldBeNil)
		So(fmt.Sprint(c.GetAll(PreDomns, PreHosts)), ShouldEqual, expPreGetAll)
		So(fmt.Sprint(c.GetAll()), ShouldEqual, expAll)
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
			Method("GET"),
			Nodes([]string{domains, hosts}),
			Prefix("address="),
			LTypes([]string{PreDomns, PreHosts, files, urls}),
		)

		So(c.ReadCfg(&CFGstatic{Cfg: CfgMimimal}), ShouldBeNil)

		tests := []struct {
			iFace IFace
			exp   string
			name  string
		}{
			{name: "ExRtObj", iFace: ExRtObj, exp: "address=/ytimg.com/0.0.0.0"},
			{name: "ExDmObj", iFace: ExDmObj, exp: ""},
			{name: "ExHtObj", iFace: ExHtObj, exp: ""},
			{name: "PreDObj", iFace: PreDObj, exp: "address=/adsrvr.org/0.0.0.0\naddress=/adtechus.net/0.0.0.0\naddress=/advertising.com/0.0.0.0\naddress=/centade.com/0.0.0.0\naddress=/doubleclick.net/0.0.0.0\naddress=/free-counter.co.uk/0.0.0.0\naddress=/intellitxt.com/0.0.0.0\naddress=/kiosked.com/0.0.0.0"},
			{name: "PreHObj", iFace: PreHObj, exp: "address=/beap.gemini.yahoo.com/192.168.168.1"},
			{name: "FileObj", iFace: FileObj, exp: "[\nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"/tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n]"},
			{name: "URLdObj", iFace: URLdObj, exp: "[\nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"malc0de\"\nnType:\t \"domn\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n]"},
			{name: "URLhObj", iFace: URLhObj, exp: "[\nDesc:\t \"Blocking mobile ad providers and some analytics providers\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"adaway\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://adaway.org/hosts.txt\"\n]"},
		}

		for _, tt := range tests {
			Convey("Testing "+tt.name+" ProcessContent()", func() {
				ct, err := c.NewContent(tt.iFace)
				So(err, ShouldBeNil)

				switch tt.iFace {
				case ExRtObj, ExDmObj, ExHtObj, PreDObj, PreHObj:
					d := &dummyConfig{t: t}
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
				Method("GET"),
				Nodes([]string{domains, hosts}),
				Prefix("address="),
				LTypes([]string{PreDomns, PreHosts, files, urls}),
			)

			tests := []struct {
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
					err:       nil,
					exp:       "[\nDesc:\t \"root-excludes exclusions\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"root-excludes\"\nName:\t \"root-excludes\"\nnType:\t \"excRoot\"\nPrefix:\t \"\"\nType:\t \"root-excludes\"\nURL:\t \"\"\n]",
					expDexMap: list{entry: entry{"ytimg.com": 0}},
					expExcMap: list{entry: entry{"ytimg.com": 0}},
					obj:       ExRtObj,
				},
				{
					name:      "ExDmObj",
					err:       nil,
					exp:       "[\nDesc:\t \"domn-excludes exclusions\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"domn-excludes\"\nName:\t \"domn-excludes\"\nnType:\t \"excDomn\"\nPrefix:\t \"\"\nType:\t \"domn-excludes\"\nURL:\t \"\"\n]",
					expDexMap: list{RWMutex: &sync.RWMutex{}, entry: make(entry)},
					expExcMap: list{RWMutex: &sync.RWMutex{}, entry: make(entry)},
					obj:       ExDmObj,
				},
				{
					name:      "ExHtObj",
					err:       nil,
					exp:       "[\nDesc:\t \"host-excludes exclusions\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"host-excludes\"\nName:\t \"host-excludes\"\nnType:\t \"excHost\"\nPrefix:\t \"\"\nType:\t \"host-excludes\"\nURL:\t \"\"\n]",
					expDexMap: list{RWMutex: &sync.RWMutex{}, entry: make(entry)},
					expExcMap: list{RWMutex: &sync.RWMutex{}, entry: make(entry)},
					obj:       ExHtObj,
				},
				{
					name: "PreDObj",
					err:  nil,
					exp:  "[\nDesc:\t \"pre-configured-domain blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured-domain\"\nName:\t \"includes.[8]\"\nnType:\t \"preDomn\"\nPrefix:\t \"\"\nType:\t \"pre-configured-domain\"\nURL:\t \"\"\n]",
					expDexMap: list{
						entry: entry{
							"adsrvr.org":         0,
							"adtechus.net":       0,
							"advertising.com":    0,
							"centade.com":        0,
							"doubleclick.net":    0,
							"free-counter.co.uk": 0,
							"intellitxt.com":     0,
							"kiosked.com":        0,
						},
					},
					expExcMap: list{entry: entry{"ytimg.com": 0}},
					f:         dir + "/pre-configured-domain.includes.[8].blacklist.conf",
					fdata:     "address=/adsrvr.org/0.0.0.0\naddress=/adtechus.net/0.0.0.0\naddress=/advertising.com/0.0.0.0\naddress=/centade.com/0.0.0.0\naddress=/doubleclick.net/0.0.0.0\naddress=/free-counter.co.uk/0.0.0.0\naddress=/intellitxt.com/0.0.0.0\naddress=/kiosked.com/0.0.0.0\n",
					obj:       PreDObj,
				},
				{
					name:      "PreHObj",
					err:       nil,
					exp:       "[\nDesc:\t \"pre-configured-host blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"pre-configured-host\"\nName:\t \"includes.[1]\"\nnType:\t \"preHost\"\nPrefix:\t \"\"\nType:\t \"pre-configured-host\"\nURL:\t \"\"\n]",
					expDexMap: list{entry: entry{"ytimg.com": 0}},
					expExcMap: list{entry: entry{"ytimg.com": 0}},
					f:         dir + "/pre-configured-host.includes.[1].blacklist.conf",
					fdata:     "address=/beap.gemini.yahoo.com/192.168.168.1\n",
					obj:       PreHObj,
				},
				{
					name: "FileObj",
					err:  fmt.Errorf("open %v/hosts./tasty.blacklist.conf: no such file or directory", dir),
					exp:  filesMin,
					expDexMap: list{
						entry: entry{
							"cw.bad.ultraadverts.site.eu": 1,
							"really.bad.phishing.site.ru": 1,
						},
					},
					expExcMap: list{entry: entry{"ytimg.com": 0}},
					f:         dir + "/hosts.tasty.blacklist.conf",
					fdata:     "address=/cw.bad.ultraadverts.site.eu/10.10.10.10\naddress=/really.bad.phishing.site.ru/10.10.10.10\n",
					obj:       FileObj,
				},
			}

			So(c.ReadCfg(&CFGstatic{Cfg: CfgMimimal}), ShouldBeNil)

			for _, tt := range tests {
				Convey("Testing "+tt.name+" ProcessContent()", func() {
					obj, err := c.NewContent(tt.obj)
					So(err, ShouldBeNil)

					if tt.f != "" {
						So(fmt.Sprint(obj), ShouldEqual, tt.exp)
					}

					if err = c.ProcessContent(obj); err != nil {
						Convey("Testing "+tt.name+" ProcessContent().Error():", func() {
							Convey("Error should match expected", func() {
								So(err, ShouldResemble, tt.err)
							})
						})
					}

					switch tt.f {
					default:
						reader, err := getFile(tt.f)
						So(err, ShouldBeNil)

						act, err := ioutil.ReadAll(reader)
						So(err, ShouldBeNil)

						Convey("Testing "+tt.name+" ProcessContent(): file data should match expected", func() {
							So(string(act), ShouldEqual, tt.fdata)
						})

					case "":
						Convey("Testing "+tt.name+" ProcessContent(): Dex map should match expected", func() {
							So(c.Dex.entry, ShouldResemble, tt.expDexMap.entry)
						})

						Convey("Testing "+tt.name+" ProcessContent(): Exc map should match expected", func() {
							So(c.Exc.entry, ShouldResemble, tt.expExcMap.entry)
						})

						Convey("Testing "+tt.name+" ProcessContent(): Obj should match expected", func() {
							So(obj.String(), ShouldEqual, tt.exp)
						})
					}
				})
			}
		})
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
				b := &blist{
					file: f.Name(),
					r:    tt.data,
				}
				So(b.writeFile(), ShouldBeNil)
				os.Remove(f.Name())

			default:
				b := &blist{
					file: tt.dir + tt.fname,
					r:    tt.data,
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
						file ../testdata/blist.hosts.src
				}
    }
}`

	// CfgMimimal contains a valid minimal EdgeOS blacklist configuration
	CfgMimimal = `blacklist {
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
			source malc0de {
					description "List of zones serving malicious executables observed by malc0de.com/database/"
					prefix "zone "
					url http://malc0de.com/bl/ZONES
			}
	}
	exclude ytimg.com
	hosts {
			dns-redirect-ip 192.168.168.1
			include beap.gemini.yahoo.com
			source tasty {
									description "File source"
									dns-redirect-ip 10.10.10.10
									file ../testdata/blist.hosts.src
							}
			source /tasty {
									description "File source"
									dns-redirect-ip 10.10.10.10
									file ../testdata/blist.hosts.src
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
					url http://localhost:8081/domains/domain.txt
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
					url http://localhost:8081/hosts/host.txt
			}
			source tasty {
					description "File source"
					dns-redirect-ip 0.0.0.0
					file ../testdata/blist.hosts.src
			}
	}
}`

	hostsContent = "address=/a.applovin.com/192.168.168.1\naddress=/a.glcdn.co/192.168.168.1\naddress=/a.vserv.mobi/192.168.168.1\naddress=/ad.leadboltapps.net/192.168.168.1\naddress=/ad.madvertise.de/192.168.168.1\naddress=/ad.where.com/192.168.168.1\naddress=/adcontent.saymedia.com/192.168.168.1\naddress=/admicro1.vcmedia.vn/192.168.168.1\naddress=/admicro2.vcmedia.vn/192.168.168.1\naddress=/admin.vserv.mobi/192.168.168.1\naddress=/ads.adiquity.com/192.168.168.1\naddress=/ads.admarvel.com/192.168.168.1\naddress=/ads.admoda.com/192.168.168.1\naddress=/ads.celtra.com/192.168.168.1\naddress=/ads.flurry.com/192.168.168.1\naddress=/ads.matomymobile.com/192.168.168.1\naddress=/ads.mobgold.com/192.168.168.1\naddress=/ads.mobilityware.com/192.168.168.1\naddress=/ads.mopub.com/192.168.168.1\naddress=/ads.n-ws.org/192.168.168.1\naddress=/ads.ookla.com/192.168.168.1\naddress=/ads.saymedia.com/192.168.168.1\naddress=/ads.smartdevicemedia.com/192.168.168.1\naddress=/ads.vserv.mobi/192.168.168.1\naddress=/ads.xxxad.net/192.168.168.1\naddress=/ads2.mediaarmor.com/192.168.168.1\naddress=/adserver.ubiyoo.com/192.168.168.1\naddress=/adultmoda.com/192.168.168.1\naddress=/android-sdk31.transpera.com/192.168.168.1\naddress=/android.bcfads.com/192.168.168.1\naddress=/api.airpush.com/192.168.168.1\naddress=/api.analytics.omgpop.com/192.168.168.1\naddress=/api.yp.com/192.168.168.1\naddress=/apps.buzzcity.net/192.168.168.1\naddress=/apps.mobilityware.com/192.168.168.1\naddress=/as.adfonic.net/192.168.168.1\naddress=/asotrack1.fluentmobile.com/192.168.168.1\naddress=/assets.cntdy.mobi/192.168.168.1\naddress=/atti.velti.com/192.168.168.1\naddress=/b.scorecardresearch.com/192.168.168.1\naddress=/banners.bigmobileads.com/192.168.168.1\naddress=/bigmobileads.com/192.168.168.1\naddress=/c.vrvm.com/192.168.168.1\naddress=/c.vserv.mobi/192.168.168.1\naddress=/cache-ssl.celtra.com/192.168.168.1\naddress=/cache.celtra.com/192.168.168.1\naddress=/cdn.celtra.com/192.168.168.1\naddress=/cdn.nearbyad.com/192.168.168.1\naddress=/cdn.trafficforce.com/192.168.168.1\naddress=/cdn.us.goldspotmedia.com/192.168.168.1\naddress=/cdn.vdopia.com/192.168.168.1\naddress=/cdn1.crispadvertising.com/192.168.168.1\naddress=/cdn1.inner-active.mobi/192.168.168.1\naddress=/cdn2.crispadvertising.com/192.168.168.1\naddress=/click.buzzcity.net/192.168.168.1\naddress=/creative1cdn.mobfox.com/192.168.168.1\naddress=/d.applovin.com/192.168.168.1\naddress=/edge.reporo.net/192.168.168.1\naddress=/ftpcontent.worldnow.com/192.168.168.1\naddress=/gemini.yahoo.com/192.168.168.1\naddress=/go.mobpartner.mobi/192.168.168.1\naddress=/go.vrvm.com/192.168.168.1\naddress=/gsmtop.net/192.168.168.1\naddress=/gts-ads.twistbox.com/192.168.168.1\naddress=/hhbekxxw5d9e.pflexads.com/192.168.168.1\naddress=/hybl9bazbc35.pflexads.com/192.168.168.1\naddress=/i.tapit.com/192.168.168.1\naddress=/images.millennialmedia.com/192.168.168.1\naddress=/images.mpression.net/192.168.168.1\naddress=/img.ads.huntmad.com/192.168.168.1\naddress=/img.ads.mobilefuse.net/192.168.168.1\naddress=/img.ads.mocean.mobi/192.168.168.1\naddress=/img.ads.mojiva.com/192.168.168.1\naddress=/img.ads.taptapnetworks.com/192.168.168.1\naddress=/m.adsymptotic.com/192.168.168.1\naddress=/m2m1.inner-active.mobi/192.168.168.1\naddress=/media.mobpartner.mobi/192.168.168.1\naddress=/medrx.sensis.com.au/192.168.168.1\naddress=/mobile.banzai.it/192.168.168.1\naddress=/mobiledl.adboe.com/192.168.168.1\naddress=/mobpartner.mobi/192.168.168.1\naddress=/mwc.velti.com/192.168.168.1\naddress=/netdna.reporo.net/192.168.168.1\naddress=/oasc04012.247realmedia.com/192.168.168.1\naddress=/orencia.pflexads.com/192.168.168.1\naddress=/pdn.applovin.com/192.168.168.1\naddress=/r.edge.inmobicdn.net/192.168.168.1\naddress=/r.mobpartner.mobi/192.168.168.1\naddress=/req.appads.com/192.168.168.1\naddress=/rs-staticart.ybcdn.net/192.168.168.1\naddress=/ru.velti.com/192.168.168.1\naddress=/s0.2mdn.net/192.168.168.1\naddress=/s3.phluant.com/192.168.168.1\naddress=/sf.vserv.mobi/192.168.168.1\naddress=/show.buzzcity.net/192.168.168.1\naddress=/static.cdn.gtsmobi.com/192.168.168.1\naddress=/static.estebull.com/192.168.168.1\naddress=/stats.pflexads.com/192.168.168.1\naddress=/track.celtra.com/192.168.168.1\naddress=/tracking.klickthru.com/192.168.168.1\naddress=/www.eltrafiko.com/192.168.168.1\naddress=/www.mmnetwork.mobi/192.168.168.1\naddress=/www.pflexads.com/192.168.168.1\naddress=/wwww.adleads.com/192.168.168.1"

	domainsContent = "address=/.192-168-0-255.com/192.1.1.1\naddress=/.asi-37.fr/192.1.1.1\naddress=/.bagbackpack.com/192.1.1.1\naddress=/.bitmeyenkartusistanbul.com/192.1.1.1\naddress=/.byxon.com/192.1.1.1\naddress=/.img001.com/192.1.1.1\naddress=/.loadto.net/192.1.1.1\naddress=/.roastfiles2017.com/192.1.1.1"

	domainsPreContent = "address=/.adsrvr.org/192.1.1.1\naddress=/.adtechus.net/192.1.1.1\naddress=/.advertising.com/192.1.1.1\naddress=/.centade.com/192.1.1.1\naddress=/.doubleclick.net/192.1.1.1\naddress=/.free-counter.co.uk/192.1.1.1\naddress=/.intellitxt.com/192.1.1.1\naddress=/.kiosked.com/192.1.1.1\n"

	expPreGetAll = "[\nDesc:\t \"pre-configured-domain blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured-domain\"\nName:\t \"includes.[8]\"\nnType:\t \"preDomn\"\nPrefix:\t \"\"\nType:\t \"pre-configured-domain\"\nURL:\t \"\"\n \nDesc:\t \"pre-configured-host blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"pre-configured-host\"\nName:\t \"includes.[1]\"\nnType:\t \"preHost\"\nPrefix:\t \"\"\nType:\t \"pre-configured-host\"\nURL:\t \"\"\n]"

	expAll = "[\nDesc:\t \"pre-configured-domain blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured-domain\"\nName:\t \"includes.[8]\"\nnType:\t \"preDomn\"\nPrefix:\t \"\"\nType:\t \"pre-configured-domain\"\nURL:\t \"\"\n \nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"malc0de\"\nnType:\t \"domn\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://localhost:8081/domains/domain.txt\"\n \nDesc:\t \"pre-configured-host blacklist content\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"pre-configured-host\"\nName:\t \"includes.[1]\"\nnType:\t \"preHost\"\nPrefix:\t \"\"\nType:\t \"pre-configured-host\"\nURL:\t \"\"\n \nDesc:\t \"Blocking mobile ad providers and some analytics providers\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"adaway\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://localhost:8081/hosts/host.txt\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../testdata/blist.hosts.src\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n]"

	domainhostContent = "address=/.192-168-0-255.com/192.1.1.1\naddress=/.asi-37.fr/192.1.1.1\naddress=/.bagbackpack.com/192.1.1.1\naddress=/.bitmeyenkartusistanbul.com/192.1.1.1\naddress=/.byxon.com/192.1.1.1\naddress=/.img001.com/192.1.1.1\naddress=/.loadto.net/192.1.1.1\naddress=/.roastfiles2017.com/192.1.1.1\naddress=/a.applovin.com/192.168.168.1\naddress=/a.glcdn.co/192.168.168.1\naddress=/a.vserv.mobi/192.168.168.1\naddress=/ad.leadboltapps.net/192.168.168.1\naddress=/ad.madvertise.de/192.168.168.1\naddress=/ad.where.com/192.168.168.1\naddress=/adcontent.saymedia.com/192.168.168.1\naddress=/admicro1.vcmedia.vn/192.168.168.1\naddress=/admicro2.vcmedia.vn/192.168.168.1\naddress=/admin.vserv.mobi/192.168.168.1\naddress=/ads.adiquity.com/192.168.168.1\naddress=/ads.admarvel.com/192.168.168.1\naddress=/ads.admoda.com/192.168.168.1\naddress=/ads.celtra.com/192.168.168.1\naddress=/ads.flurry.com/192.168.168.1\naddress=/ads.matomymobile.com/192.168.168.1\naddress=/ads.mobgold.com/192.168.168.1\naddress=/ads.mobilityware.com/192.168.168.1\naddress=/ads.mopub.com/192.168.168.1\naddress=/ads.n-ws.org/192.168.168.1\naddress=/ads.ookla.com/192.168.168.1\naddress=/ads.saymedia.com/192.168.168.1\naddress=/ads.smartdevicemedia.com/192.168.168.1\naddress=/ads.vserv.mobi/192.168.168.1\naddress=/ads.xxxad.net/192.168.168.1\naddress=/ads2.mediaarmor.com/192.168.168.1\naddress=/adserver.ubiyoo.com/192.168.168.1\naddress=/adultmoda.com/192.168.168.1\naddress=/android-sdk31.transpera.com/192.168.168.1\naddress=/android.bcfads.com/192.168.168.1\naddress=/api.airpush.com/192.168.168.1\naddress=/api.analytics.omgpop.com/192.168.168.1\naddress=/api.yp.com/192.168.168.1\naddress=/apps.buzzcity.net/192.168.168.1\naddress=/apps.mobilityware.com/192.168.168.1\naddress=/as.adfonic.net/192.168.168.1\naddress=/asotrack1.fluentmobile.com/192.168.168.1\naddress=/assets.cntdy.mobi/192.168.168.1\naddress=/atti.velti.com/192.168.168.1\naddress=/b.scorecardresearch.com/192.168.168.1\naddress=/banners.bigmobileads.com/192.168.168.1\naddress=/bigmobileads.com/192.168.168.1\naddress=/c.vrvm.com/192.168.168.1\naddress=/c.vserv.mobi/192.168.168.1\naddress=/cache-ssl.celtra.com/192.168.168.1\naddress=/cache.celtra.com/192.168.168.1\naddress=/cdn.celtra.com/192.168.168.1\naddress=/cdn.nearbyad.com/192.168.168.1\naddress=/cdn.trafficforce.com/192.168.168.1\naddress=/cdn.us.goldspotmedia.com/192.168.168.1\naddress=/cdn.vdopia.com/192.168.168.1\naddress=/cdn1.crispadvertising.com/192.168.168.1\naddress=/cdn1.inner-active.mobi/192.168.168.1\naddress=/cdn2.crispadvertising.com/192.168.168.1\naddress=/click.buzzcity.net/192.168.168.1\naddress=/creative1cdn.mobfox.com/192.168.168.1\naddress=/d.applovin.com/192.168.168.1\naddress=/edge.reporo.net/192.168.168.1\naddress=/ftpcontent.worldnow.com/192.168.168.1\naddress=/gemini.yahoo.com/192.168.168.1\naddress=/go.mobpartner.mobi/192.168.168.1\naddress=/go.vrvm.com/192.168.168.1\naddress=/gsmtop.net/192.168.168.1\naddress=/gts-ads.twistbox.com/192.168.168.1\naddress=/hhbekxxw5d9e.pflexads.com/192.168.168.1\naddress=/hybl9bazbc35.pflexads.com/192.168.168.1\naddress=/i.tapit.com/192.168.168.1\naddress=/images.millennialmedia.com/192.168.168.1\naddress=/images.mpression.net/192.168.168.1\naddress=/img.ads.huntmad.com/192.168.168.1\naddress=/img.ads.mobilefuse.net/192.168.168.1\naddress=/img.ads.mocean.mobi/192.168.168.1\naddress=/img.ads.mojiva.com/192.168.168.1\naddress=/img.ads.taptapnetworks.com/192.168.168.1\naddress=/m.adsymptotic.com/192.168.168.1\naddress=/m2m1.inner-active.mobi/192.168.168.1\naddress=/media.mobpartner.mobi/192.168.168.1\naddress=/medrx.sensis.com.au/192.168.168.1\naddress=/mobile.banzai.it/192.168.168.1\naddress=/mobiledl.adboe.com/192.168.168.1\naddress=/mobpartner.mobi/192.168.168.1\naddress=/mwc.velti.com/192.168.168.1\naddress=/netdna.reporo.net/192.168.168.1\naddress=/oasc04012.247realmedia.com/192.168.168.1\naddress=/orencia.pflexads.com/192.168.168.1\naddress=/pdn.applovin.com/192.168.168.1\naddress=/r.edge.inmobicdn.net/192.168.168.1\naddress=/r.mobpartner.mobi/192.168.168.1\naddress=/req.appads.com/192.168.168.1\naddress=/rs-staticart.ybcdn.net/192.168.168.1\naddress=/ru.velti.com/192.168.168.1\naddress=/s0.2mdn.net/192.168.168.1\naddress=/s3.phluant.com/192.168.168.1\naddress=/sf.vserv.mobi/192.168.168.1\naddress=/show.buzzcity.net/192.168.168.1\naddress=/static.cdn.gtsmobi.com/192.168.168.1\naddress=/static.estebull.com/192.168.168.1\naddress=/stats.pflexads.com/192.168.168.1\naddress=/track.celtra.com/192.168.168.1\naddress=/tracking.klickthru.com/192.168.168.1\naddress=/www.eltrafiko.com/192.168.168.1\naddress=/www.mmnetwork.mobi/192.168.168.1\naddress=/www.pflexads.com/192.168.168.1\naddress=/wwww.adleads.com/192.168.168.1"

	domainMin = "address=/.01lm.com/0.0.0.0\naddress=/.2biking.com/0.0.0.0\naddress=/.323trs.com/0.0.0.0\naddress=/.51jetso.com/0.0.0.0\naddress=/.52zsoft.com/0.0.0.0\naddress=/.54nb.com/0.0.0.0\naddress=/.9364.org/0.0.0.0\naddress=/.antalyanalburiye.com/0.0.0.0\naddress=/.bellefonte.net/0.0.0.0\naddress=/.bow-spell-effect1.ru/0.0.0.0\naddress=/.bplaced.net/0.0.0.0\naddress=/.cloudme.com/0.0.0.0\naddress=/.falcogames.com/0.0.0.0\naddress=/.freegamer.info/0.0.0.0\naddress=/.frizoupuzzles.org/0.0.0.0\naddress=/.fssblangenlois.ac.at/0.0.0.0\naddress=/.gamegogle.com/0.0.0.0\naddress=/.gasparini.com.br/0.0.0.0\naddress=/.getpics.net/0.0.0.0\naddress=/.gezila.com/0.0.0.0\naddress=/.glazeautocaremobile.com/0.0.0.0\naddress=/.goldenlifewomen.com/0.0.0.0\naddress=/.goosai.com/0.0.0.0\naddress=/.holidaysinkeralam.com/0.0.0.0\naddress=/.hotlaps.com.au/0.0.0.0\naddress=/.i2cchip.com/0.0.0.0\naddress=/.ibxdnl.com/0.0.0.0\naddress=/.igetmyservice.com/0.0.0.0\naddress=/.iprojhq.com/0.0.0.0\naddress=/.izmirhavaalaniarackiralama.net/0.0.0.0\naddress=/.jingshang.com.tw/0.0.0.0\naddress=/.justgetitfaster.com/0.0.0.0\naddress=/.kanberdemir.com/0.0.0.0\naddress=/.kpzip.com/0.0.0.0\naddress=/.kraonkelaere.com/0.0.0.0\naddress=/.laptopb4you.com/0.0.0.0\naddress=/.liftune.com/0.0.0.0\naddress=/.m-games.huu.cz/0.0.0.0\naddress=/.martiniracing.com.br/0.0.0.0\naddress=/.mireene.com/0.0.0.0\naddress=/.mixtrio.net/0.0.0.0\naddress=/.mstdls.com/0.0.0.0\naddress=/.mypcapp.com/0.0.0.0\naddress=/.perso.sfr.fr/0.0.0.0\naddress=/.pixelmon-world.com/0.0.0.0\naddress=/.plexcera.com/0.0.0.0\naddress=/.rd1994.com/0.0.0.0\naddress=/.sf-addon.com/0.0.0.0\naddress=/.skypedong.com/0.0.0.0\naddress=/.spirlymo.com/0.0.0.0\naddress=/.sportstherapy.net/0.0.0.0\naddress=/.talka-studios.com/0.0.0.0\naddress=/.thewitchez-cafe.co.uk/0.0.0.0\naddress=/.tirekoypazari.com/0.0.0.0\naddress=/.updatestar.net/0.0.0.0\naddress=/.urban-garden.net/0.0.0.0\naddress=/.utilbada.com/0.0.0.0\naddress=/.utilcom.net/0.0.0.0\naddress=/.utiljoy.com/0.0.0.0\naddress=/.vim6.com/0.0.0.0\naddress=/.windows.net/0.0.0.0\naddress=/.xiazai4.net/0.0.0.0\naddress=/.xunyou.com/0.0.0.0\n"

	filesMin = "[\nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"/tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n]"

	excRootContent = "address=/122.2o7.net/0.0.0.0\naddress=/1e100.net/0.0.0.0\naddress=/adobedtm.com/0.0.0.0\naddress=/akamai.net/0.0.0.0\naddress=/amazon.com/0.0.0.0\naddress=/amazonaws.com/0.0.0.0\naddress=/apple.com/0.0.0.0\naddress=/ask.com/0.0.0.0\naddress=/avast.com/0.0.0.0\naddress=/bitdefender.com/0.0.0.0\naddress=/cdn.visiblemeasures.com/0.0.0.0\naddress=/cloudfront.net/0.0.0.0\naddress=/coremetrics.com/0.0.0.0\naddress=/edgesuite.net/0.0.0.0\naddress=/freedns.afraid.org/0.0.0.0\naddress=/github.com/0.0.0.0\naddress=/githubusercontent.com/0.0.0.0\naddress=/google.com/0.0.0.0\naddress=/googleadservices.com/0.0.0.0\naddress=/googleapis.com/0.0.0.0\naddress=/googleusercontent.com/0.0.0.0\naddress=/gstatic.com/0.0.0.0\naddress=/gvt1.com/0.0.0.0\naddress=/gvt1.net/0.0.0.0\naddress=/hb.disney.go.com/0.0.0.0\naddress=/hp.com/0.0.0.0\naddress=/hulu.com/0.0.0.0\naddress=/images-amazon.com/0.0.0.0\naddress=/jumptap.com/0.0.0.0\naddress=/msdn.com/0.0.0.0\naddress=/paypal.com/0.0.0.0\naddress=/rackcdn.com/0.0.0.0\naddress=/schema.org/0.0.0.0\naddress=/skype.com/0.0.0.0\naddress=/smacargo.com/0.0.0.0\naddress=/sourceforge.net/0.0.0.0\naddress=/ssl-on9.com/0.0.0.0\naddress=/ssl-on9.net/0.0.0.0\naddress=/static.chartbeat.com/0.0.0.0\naddress=/storage.googleapis.com/0.0.0.0\naddress=/usemaxserver.de/0.0.0.0\naddress=/windows.net/0.0.0.0\naddress=/yimg.com/0.0.0.0\naddress=/ytimg.com/0.0.0.0"

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
									file /:~/=../testdata/blist.hosts.src
							}
	}
}`
)

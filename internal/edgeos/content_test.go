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
	"testing"
	"time"

	. "github.com/britannic/testutils"
)

type ProcessorContenter interface {
	ProcessContent()
}

type dummyConfig struct {
	s []string
	t *testing.T
}

func (d *dummyConfig) ProcessContent(ct Contenter) {
	for _, src := range *ct.GetBlacklist() {
		b, err := ioutil.ReadAll(src.process().r)
		OK(d.t, err)
		d.s = append(d.s, strings.TrimSuffix(string(b), "\n"))
	}
}

func TestCreateObject(t *testing.T) {

	tests := []struct {
		err       error
		fail      bool
		name      string
		obj       iFace
		exp       string
		leaf      string
		ltype     string
		page      string
		page2     string
		pageData  string
		pageData2 string
		pos       int
		svr       *HTTPserver
		svr2      *HTTPserver
	}{
		{
			exp:   "address=/.adsrvr.org/192.1.1.1\naddress=/.adtechus.net/192.1.1.1\naddress=/.advertising.com/192.1.1.1\naddress=/.centade.com/192.1.1.1\naddress=/.doubleclick.net/192.1.1.1\naddress=/.free-counter.co.uk/192.1.1.1\naddress=/.intellitxt.com/192.1.1.1\naddress=/.kiosked.com/192.1.1.1\naddress=/beap.gemini.yahoo.com/0.0.0.0",
			fail:  false,
			ltype: preConf,
			name:  preConf,
			obj:   PreObj,
			pos:   0,
		},
		{
			exp:   "\n",
			fail:  false,
			ltype: preConf,
			name:  "pre",
			obj:   PreObj,
			pos:   -1,
		},
		{
			exp:   "address=/really.bad.phishing.site.ru/0.0.0.0",
			fail:  false,
			ltype: files,
			name:  "tasty",
			obj:   FileObj,
			pos:   0,
		},
		{
			exp:   "",
			fail:  false,
			ltype: files,
			name:  "ztasty",
			obj:   FileObj,
			pos:   -1,
		},
		{
			exp:       domainhostContent,
			fail:      false,
			ltype:     urls,
			name:      "malc0de",
			obj:       URLObj,
			pos:       0,
			page:      "/hosts.txt",
			page2:     "/domains.txt",
			pageData:  httpHostData,
			pageData2: HTTPDomainData,
			svr:       new(HTTPserver),
			svr2:      new(HTTPserver),
		},
		{
			exp:       "\n",
			fail:      false,
			ltype:     urls,
			name:      "zmalc0de",
			obj:       URLObj,
			pos:       -1,
			page:      "/hosts.txt",
			page2:     "/domains.txt",
			pageData:  httpHostData,
			pageData2: HTTPDomainData,
			svr:       new(HTTPserver),
			svr2:      new(HTTPserver),
		},
		{
			exp:  "",
			fail: false,
			name: cntnt,
			obj:  ContObj,
			pos:  0,
		},
		{
			exp:  "",
			fail: false,
			name: "not contents",
			obj:  ContObj,
			pos:  -1,
		},
		{
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
		Nodes([]string{"domains", "hosts"}),
		Prefix("address="),
		STypes([]string{files, preConf, urls}),
		Timeout(30*time.Second),
		WCard(Wildcard{Node: "*s", Name: "*"}),
	)

	r := &CFGstatic{Cfg: Cfg}
	err := c.ReadCfg(r)
	OK(t, err)
	c.SetOpt(
		Dexcludes(c.Excludes(rootNode, domains)),
		Excludes(c.Excludes(hosts)),
	)

	for _, tt := range tests {
		objs, err := c.CreateObject(tt.obj)
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
			OK(t, err)
			d := &dummyConfig{}
			d.ProcessContent(objs)
			d.t = t
			Equals(t, tt.exp, strings.Join(d.s, "\n"))

			objs.SetURL(tt.name, tt.name)
			Equals(t, tt.pos, objs.Find(tt.name))

		default:
			Equals(t, tt.err, err)
		}
	}
}

func TestGetAllContent(t *testing.T) {
	var (
		r = &CFGstatic{Cfg: testallCfg}
		c = NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
			FileNameFmt("%v/%v.%v.%v"),
			Method("GET"),
			Nodes([]string{"domains", "hosts"}),
			Prefix("address="),
			STypes([]string{preConf, "file", urls}),
		)
	)

	err := c.ReadCfg(r)
	OK(t, err)

	act := fmt.Sprint(c.GetAll(preConf))
	Equals(t, wantPre, act)

	act = fmt.Sprint(c.GetAll())
	Equals(t, wantAll, act)
}

func TestProcessContent(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "testBlacklist")
	OK(t, err)
	defer os.RemoveAll(dir)

	var (
		c = NewConfig(
			Dir(dir),
			Ext("blacklist.conf"),
			FileNameFmt("%v/%v.%v.%v"),
			Method("GET"),
			Nodes([]string{"domains", "hosts"}),
			Prefix("address="),
			STypes([]string{preConf, "file", urls}),
		)
		r     = &CFGstatic{Cfg: CfgMimimal}
		tests = []struct {
			err   error
			exp   string
			f     []string
			fdata []string
			obj   iFace
		}{
			{
				exp:   wantPre,
				err:   nil,
				fdata: []string{"address=/.adsrvr.org/0.0.0.0\naddress=/.adtechus.net/0.0.0.0\naddress=/.advertising.com/0.0.0.0\naddress=/.centade.com/0.0.0.0\naddress=/.doubleclick.net/0.0.0.0\naddress=/.free-counter.co.uk/0.0.0.0\naddress=/.intellitxt.com/0.0.0.0\naddress=/.kiosked.com/0.0.0.0\n", "address=/beap.gemini.yahoo.com/0.0.0.0\n"},
				obj:   PreObj,
			},
			{
				err:   errors.New("open " + dir + "/hosts./tasty.blacklist.conf: no such file or directory"),
				exp:   filesMin,
				f:     []string{dir + "/hosts.tasty.blacklist.conf", dir + "/hosts.tasty.blacklist.conf"},
				fdata: []string{"address=/really.bad.phishing.site.ru/10.10.10.10\n", "address=/really.bad.phishing.site.ru/10.10.10.10\n"},
				obj:   FileObj,
			},
			{
				err:   nil,
				exp:   "[\nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"malc0de\"\nnType:\t \"domain\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n]",
				f:     []string{dir + "/domains.malc0de.blacklist.conf"},
				fdata: []string{domainMin},
				obj:   URLObj,
			},
		}
	)

	err = c.ReadCfg(r)
	OK(t, err)

	for _, tt := range tests {
		obj, err := c.CreateObject(tt.obj)
		OK(t, err)

		Equals(t, tt.exp, fmt.Sprint(obj))

		if err = c.ProcessContent(obj); err != nil {
			Equals(t, tt.err, err)
		}

		for i, f := range tt.f {
			reader, err := getFile(f)
			OK(t, err)

			act, err := ioutil.ReadAll(reader)
			OK(t, err)

			Equals(t, tt.fdata[i], string(act))
		}
	}
}

func TestWriteFile(t *testing.T) {
	writeFileTests := []struct {
		Content
		data  io.Reader
		dir   string
		fname string
		ok    bool
		want  string
	}{
		{
			data:  strings.NewReader("The rest is history!"),
			dir:   "/tmp",
			fname: "Test.util.WriteFile",
			ok:    true,
			want:  "",
		},
		{
			data:  bytes.NewBuffer([]byte{84, 104, 101, 32, 114, 101, 115, 116, 32, 105, 115, 32, 104, 105, 115, 116, 111, 114, 121, 33}),
			dir:   "/tmp",
			fname: "Test.util.WriteFile",
			ok:    true,
			want:  "",
		},
		{
			data:  bytes.NewBufferString("This shouldn't be written!"),
			dir:   "",
			fname: "/tmp/z/d/c/r/c:reallybadfile.zfts",
			ok:    false,
			want:  `unable to open file: /tmp/z/d/c/r/c:reallybadfile.zfts for writing, error: open /tmp/z/d/c/r/c:reallybadfile.zfts: no such file or directory`,
		},
	}

	c := Config{Parms: NewParms()}
	c.SetOpt(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		FileNameFmt("%v/%v.%v.%v"),
		Nodes([]string{"domains", "hosts"}),
	)

	for _, test := range writeFileTests {
		switch test.ok {
		case true:
			f, err := ioutil.TempFile(test.dir, test.fname)
			OK(t, err)
			b := &blist{
				file: f.Name(),
				r:    test.data,
			}
			err = b.WriteFile()
			OK(t, err)
			os.Remove(f.Name())

		default:
			b := &blist{
				file: test.dir + test.fname,
				r:    test.data,
			}
			err := b.WriteFile()
			NotOK(t, err)
			Equals(t, "open /tmp/z/d/c/r/c:reallybadfile.zfts: no such file or directory", err.Error())
		}
	}
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

	hostsContent = "address=/a.applovin.com/192.168.168.1\naddress=/a.glcdn.co/192.168.168.1\naddress=/a.vserv.mobi/192.168.168.1\naddress=/ad.leadboltapps.net/192.168.168.1\naddress=/ad.madvertise.de/192.168.168.1\naddress=/ad.where.com/192.168.168.1\naddress=/ad1.adinfuse.com/192.168.168.1\naddress=/ad2.adinfuse.com/192.168.168.1\naddress=/adcontent.saymedia.com/192.168.168.1\naddress=/adinfuse.com/192.168.168.1\naddress=/admicro1.vcmedia.vn/192.168.168.1\naddress=/admicro2.vcmedia.vn/192.168.168.1\naddress=/admin.vserv.mobi/192.168.168.1\naddress=/ads.adiquity.com/192.168.168.1\naddress=/ads.admarvel.com/192.168.168.1\naddress=/ads.admoda.com/192.168.168.1\naddress=/ads.celtra.com/192.168.168.1\naddress=/ads.flurry.com/192.168.168.1\naddress=/ads.matomymobile.com/192.168.168.1\naddress=/ads.mobgold.com/192.168.168.1\naddress=/ads.mobilityware.com/192.168.168.1\naddress=/ads.mopub.com/192.168.168.1\naddress=/ads.n-ws.org/192.168.168.1\naddress=/ads.ookla.com/192.168.168.1\naddress=/ads.saymedia.com/192.168.168.1\naddress=/ads.smartdevicemedia.com/192.168.168.1\naddress=/ads.vserv.mobi/192.168.168.1\naddress=/ads.xxxad.net/192.168.168.1\naddress=/ads2.mediaarmor.com/192.168.168.1\naddress=/adserver.ubiyoo.com/192.168.168.1\naddress=/adultmoda.com/192.168.168.1\naddress=/android-sdk31.transpera.com/192.168.168.1\naddress=/android.bcfads.com/192.168.168.1\naddress=/api.airpush.com/192.168.168.1\naddress=/api.analytics.omgpop.com/192.168.168.1\naddress=/api.yp.com/192.168.168.1\naddress=/apps.buzzcity.net/192.168.168.1\naddress=/apps.mobilityware.com/192.168.168.1\naddress=/as.adfonic.net/192.168.168.1\naddress=/asotrack1.fluentmobile.com/192.168.168.1\naddress=/assets.cntdy.mobi/192.168.168.1\naddress=/atti.velti.com/192.168.168.1\naddress=/b.scorecardresearch.com/192.168.168.1\naddress=/banners.bigmobileads.com/192.168.168.1\naddress=/bigmobileads.com/192.168.168.1\naddress=/bo.jumptap.com/192.168.168.1\naddress=/bos-tapreq01.jumptap.com/192.168.168.1\naddress=/bos-tapreq02.jumptap.com/192.168.168.1\naddress=/bos-tapreq03.jumptap.com/192.168.168.1\naddress=/bos-tapreq04.jumptap.com/192.168.168.1\naddress=/bos-tapreq05.jumptap.com/192.168.168.1\naddress=/bos-tapreq06.jumptap.com/192.168.168.1\naddress=/bos-tapreq07.jumptap.com/192.168.168.1\naddress=/bos-tapreq08.jumptap.com/192.168.168.1\naddress=/bos-tapreq09.jumptap.com/192.168.168.1\naddress=/bos-tapreq10.jumptap.com/192.168.168.1\naddress=/bos-tapreq11.jumptap.com/192.168.168.1\naddress=/bos-tapreq12.jumptap.com/192.168.168.1\naddress=/bos-tapreq13.jumptap.com/192.168.168.1\naddress=/bos-tapreq14.jumptap.com/192.168.168.1\naddress=/bos-tapreq15.jumptap.com/192.168.168.1\naddress=/bos-tapreq16.jumptap.com/192.168.168.1\naddress=/bos-tapreq17.jumptap.com/192.168.168.1\naddress=/bos-tapreq18.jumptap.com/192.168.168.1\naddress=/bos-tapreq19.jumptap.com/192.168.168.1\naddress=/bos-tapreq20.jumptap.com/192.168.168.1\naddress=/c.vrvm.com/192.168.168.1\naddress=/c.vserv.mobi/192.168.168.1\naddress=/c753738.r38.cf2.rackcdn.com/192.168.168.1\naddress=/cache-ssl.celtra.com/192.168.168.1\naddress=/cache.celtra.com/192.168.168.1\naddress=/cdn.celtra.com/192.168.168.1\naddress=/cdn.nearbyad.com/192.168.168.1\naddress=/cdn.trafficforce.com/192.168.168.1\naddress=/cdn.us.goldspotmedia.com/192.168.168.1\naddress=/cdn.vdopia.com/192.168.168.1\naddress=/cdn1.crispadvertising.com/192.168.168.1\naddress=/cdn1.inner-active.mobi/192.168.168.1\naddress=/cdn2.crispadvertising.com/192.168.168.1\naddress=/click.buzzcity.net/192.168.168.1\naddress=/creative1cdn.mobfox.com/192.168.168.1\naddress=/d.applovin.com/192.168.168.1\naddress=/d2bgg7rjywcwsy.cloudfront.net/192.168.168.1\naddress=/d3anogn3pbtk4v.cloudfront.net/192.168.168.1\naddress=/d3oltyb66oj2v8.cloudfront.net/192.168.168.1\naddress=/edge.reporo.net/192.168.168.1\naddress=/ftpcontent.worldnow.com/192.168.168.1\naddress=/funnel0.adinfuse.com/192.168.168.1\naddress=/gemini.yahoo.com/192.168.168.1\naddress=/go.adinfuse.com/192.168.168.1\naddress=/go.mobpartner.mobi/192.168.168.1\naddress=/go.vrvm.com/192.168.168.1\naddress=/gsmtop.net/192.168.168.1\naddress=/gts-ads.twistbox.com/192.168.168.1\naddress=/hhbekxxw5d9e.pflexads.com/192.168.168.1\naddress=/hybl9bazbc35.pflexads.com/192.168.168.1\naddress=/i.jumptap.com/192.168.168.1\naddress=/i.tapit.com/192.168.168.1\naddress=/images.millennialmedia.com/192.168.168.1\naddress=/images.mpression.net/192.168.168.1\naddress=/img.ads.huntmad.com/192.168.168.1\naddress=/img.ads.mobilefuse.net/192.168.168.1\naddress=/img.ads.mocean.mobi/192.168.168.1\naddress=/img.ads.mojiva.com/192.168.168.1\naddress=/img.ads.taptapnetworks.com/192.168.168.1\naddress=/intouch.adinfuse.com/192.168.168.1\naddress=/lb.usemaxserver.de/192.168.168.1\naddress=/m.adsymptotic.com/192.168.168.1\naddress=/m2m1.inner-active.mobi/192.168.168.1\naddress=/media.mobpartner.mobi/192.168.168.1\naddress=/medrx.sensis.com.au/192.168.168.1\naddress=/mobile.banzai.it/192.168.168.1\naddress=/mobiledl.adboe.com/192.168.168.1\naddress=/mobpartner.mobi/192.168.168.1\naddress=/mwc.velti.com/192.168.168.1\naddress=/netdna.reporo.net/192.168.168.1\naddress=/oasc04012.247realmedia.com/192.168.168.1\naddress=/orange-fr.adinfuse.com/192.168.168.1\naddress=/orangeuk-mc.adinfuse.com/192.168.168.1\naddress=/orencia.pflexads.com/192.168.168.1\naddress=/pdn.applovin.com/192.168.168.1\naddress=/r.edge.inmobicdn.net/192.168.168.1\naddress=/r.mobpartner.mobi/192.168.168.1\naddress=/req.appads.com/192.168.168.1\naddress=/rs-staticart.ybcdn.net/192.168.168.1\naddress=/ru.velti.com/192.168.168.1\naddress=/s0.2mdn.net/192.168.168.1\naddress=/s3.phluant.com/192.168.168.1\naddress=/sf.vserv.mobi/192.168.168.1\naddress=/show.buzzcity.net/192.168.168.1\naddress=/sky-connect.adinfuse.com/192.168.168.1\naddress=/sky.adinfuse.com/192.168.168.1\naddress=/static.cdn.gtsmobi.com/192.168.168.1\naddress=/static.estebull.com/192.168.168.1\naddress=/stats.pflexads.com/192.168.168.1\naddress=/track.celtra.com/192.168.168.1\naddress=/tracking.klickthru.com/192.168.168.1\naddress=/uk-ad2.adinfuse.com/192.168.168.1\naddress=/uk-go.adinfuse.com/192.168.168.1\naddress=/web63.jumptap.com/192.168.168.1\naddress=/web64.jumptap.com/192.168.168.1\naddress=/web65.jumptap.com/192.168.168.1\naddress=/wv.inner-active.mobi/192.168.168.1\naddress=/www.eltrafiko.com/192.168.168.1\naddress=/www.mmnetwork.mobi/192.168.168.1\naddress=/www.pflexads.com/192.168.168.1\naddress=/wwww.adleads.com/192.168.168.1\n"

	domainsContent = "address=/.192-168-0-255.com/192.1.1.1\naddress=/.asi-37.fr/192.1.1.1\naddress=/.bagbackpack.com/192.1.1.1\naddress=/.bitmeyenkartusistanbul.com/192.1.1.1\naddress=/.byxon.com/192.1.1.1\naddress=/.img001.com/192.1.1.1\naddress=/.loadto.net/192.1.1.1\naddress=/.roastfiles2017.com/192.1.1.1\n"

	domainsPreContent = "address=/.adsrvr.org/192.1.1.1\naddress=/.adtechus.net/192.1.1.1\naddress=/.advertising.com/192.1.1.1\naddress=/.centade.com/192.1.1.1\naddress=/.doubleclick.net/192.1.1.1\naddress=/.free-counter.co.uk/192.1.1.1\naddress=/.intellitxt.com/192.1.1.1\naddress=/.kiosked.com/192.1.1.1\n"

	wantPre = "[\nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"\"\n \nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n]"

	wantAll = "[\nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"\"\n \nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"url\"\nName:\t \"malc0de\"\nnType:\t \"domain\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://localhost:8081/domains/domain.txt\"\n \nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"Blocking mobile ad providers and some analytics providers\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"url\"\nName:\t \"adaway\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://localhost:8081/hosts/host.txt\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../testdata/blist.hosts.src\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n]"

	domainhostContent = "address=/.192-168-0-255.com/192.1.1.1\naddress=/.asi-37.fr/192.1.1.1\naddress=/.bagbackpack.com/192.1.1.1\naddress=/.bitmeyenkartusistanbul.com/192.1.1.1\naddress=/.byxon.com/192.1.1.1\naddress=/.img001.com/192.1.1.1\naddress=/.loadto.net/192.1.1.1\naddress=/.roastfiles2017.com/192.1.1.1\naddress=/a.applovin.com/192.168.168.1\naddress=/a.glcdn.co/192.168.168.1\naddress=/a.vserv.mobi/192.168.168.1\naddress=/ad.leadboltapps.net/192.168.168.1\naddress=/ad.madvertise.de/192.168.168.1\naddress=/ad.where.com/192.168.168.1\naddress=/adcontent.saymedia.com/192.168.168.1\naddress=/admicro1.vcmedia.vn/192.168.168.1\naddress=/admicro2.vcmedia.vn/192.168.168.1\naddress=/admin.vserv.mobi/192.168.168.1\naddress=/ads.adiquity.com/192.168.168.1\naddress=/ads.admarvel.com/192.168.168.1\naddress=/ads.admoda.com/192.168.168.1\naddress=/ads.celtra.com/192.168.168.1\naddress=/ads.flurry.com/192.168.168.1\naddress=/ads.matomymobile.com/192.168.168.1\naddress=/ads.mobgold.com/192.168.168.1\naddress=/ads.mobilityware.com/192.168.168.1\naddress=/ads.mopub.com/192.168.168.1\naddress=/ads.n-ws.org/192.168.168.1\naddress=/ads.ookla.com/192.168.168.1\naddress=/ads.saymedia.com/192.168.168.1\naddress=/ads.smartdevicemedia.com/192.168.168.1\naddress=/ads.vserv.mobi/192.168.168.1\naddress=/ads.xxxad.net/192.168.168.1\naddress=/ads2.mediaarmor.com/192.168.168.1\naddress=/adserver.ubiyoo.com/192.168.168.1\naddress=/adultmoda.com/192.168.168.1\naddress=/android-sdk31.transpera.com/192.168.168.1\naddress=/android.bcfads.com/192.168.168.1\naddress=/api.airpush.com/192.168.168.1\naddress=/api.analytics.omgpop.com/192.168.168.1\naddress=/api.yp.com/192.168.168.1\naddress=/apps.buzzcity.net/192.168.168.1\naddress=/apps.mobilityware.com/192.168.168.1\naddress=/as.adfonic.net/192.168.168.1\naddress=/asotrack1.fluentmobile.com/192.168.168.1\naddress=/assets.cntdy.mobi/192.168.168.1\naddress=/atti.velti.com/192.168.168.1\naddress=/b.scorecardresearch.com/192.168.168.1\naddress=/banners.bigmobileads.com/192.168.168.1\naddress=/bigmobileads.com/192.168.168.1\naddress=/c.vrvm.com/192.168.168.1\naddress=/c.vserv.mobi/192.168.168.1\naddress=/cache-ssl.celtra.com/192.168.168.1\naddress=/cache.celtra.com/192.168.168.1\naddress=/cdn.celtra.com/192.168.168.1\naddress=/cdn.nearbyad.com/192.168.168.1\naddress=/cdn.trafficforce.com/192.168.168.1\naddress=/cdn.us.goldspotmedia.com/192.168.168.1\naddress=/cdn.vdopia.com/192.168.168.1\naddress=/cdn1.crispadvertising.com/192.168.168.1\naddress=/cdn1.inner-active.mobi/192.168.168.1\naddress=/cdn2.crispadvertising.com/192.168.168.1\naddress=/click.buzzcity.net/192.168.168.1\naddress=/creative1cdn.mobfox.com/192.168.168.1\naddress=/d.applovin.com/192.168.168.1\naddress=/edge.reporo.net/192.168.168.1\naddress=/ftpcontent.worldnow.com/192.168.168.1\naddress=/gemini.yahoo.com/192.168.168.1\naddress=/go.mobpartner.mobi/192.168.168.1\naddress=/go.vrvm.com/192.168.168.1\naddress=/gsmtop.net/192.168.168.1\naddress=/gts-ads.twistbox.com/192.168.168.1\naddress=/hhbekxxw5d9e.pflexads.com/192.168.168.1\naddress=/hybl9bazbc35.pflexads.com/192.168.168.1\naddress=/i.tapit.com/192.168.168.1\naddress=/images.millennialmedia.com/192.168.168.1\naddress=/images.mpression.net/192.168.168.1\naddress=/img.ads.huntmad.com/192.168.168.1\naddress=/img.ads.mobilefuse.net/192.168.168.1\naddress=/img.ads.mocean.mobi/192.168.168.1\naddress=/img.ads.mojiva.com/192.168.168.1\naddress=/img.ads.taptapnetworks.com/192.168.168.1\naddress=/m.adsymptotic.com/192.168.168.1\naddress=/m2m1.inner-active.mobi/192.168.168.1\naddress=/media.mobpartner.mobi/192.168.168.1\naddress=/medrx.sensis.com.au/192.168.168.1\naddress=/mobile.banzai.it/192.168.168.1\naddress=/mobiledl.adboe.com/192.168.168.1\naddress=/mobpartner.mobi/192.168.168.1\naddress=/mwc.velti.com/192.168.168.1\naddress=/netdna.reporo.net/192.168.168.1\naddress=/oasc04012.247realmedia.com/192.168.168.1\naddress=/orencia.pflexads.com/192.168.168.1\naddress=/pdn.applovin.com/192.168.168.1\naddress=/r.edge.inmobicdn.net/192.168.168.1\naddress=/r.mobpartner.mobi/192.168.168.1\naddress=/req.appads.com/192.168.168.1\naddress=/rs-staticart.ybcdn.net/192.168.168.1\naddress=/ru.velti.com/192.168.168.1\naddress=/s0.2mdn.net/192.168.168.1\naddress=/s3.phluant.com/192.168.168.1\naddress=/sf.vserv.mobi/192.168.168.1\naddress=/show.buzzcity.net/192.168.168.1\naddress=/static.cdn.gtsmobi.com/192.168.168.1\naddress=/static.estebull.com/192.168.168.1\naddress=/stats.pflexads.com/192.168.168.1\naddress=/track.celtra.com/192.168.168.1\naddress=/tracking.klickthru.com/192.168.168.1\naddress=/www.eltrafiko.com/192.168.168.1\naddress=/www.mmnetwork.mobi/192.168.168.1\naddress=/www.pflexads.com/192.168.168.1\naddress=/wwww.adleads.com/192.168.168.1"

	domainMin = "address=/.01lm.com/0.0.0.0\naddress=/.2biking.com/0.0.0.0\naddress=/.323trs.com/0.0.0.0\naddress=/.51jetso.com/0.0.0.0\naddress=/.52zsoft.com/0.0.0.0\naddress=/.54nb.com/0.0.0.0\naddress=/.9364.org/0.0.0.0\naddress=/.antalyanalburiye.com/0.0.0.0\naddress=/.bellefonte.net/0.0.0.0\naddress=/.bow-spell-effect1.ru/0.0.0.0\naddress=/.bplaced.net/0.0.0.0\naddress=/.cloudme.com/0.0.0.0\naddress=/.falcogames.com/0.0.0.0\naddress=/.freegamer.info/0.0.0.0\naddress=/.frizoupuzzles.org/0.0.0.0\naddress=/.fssblangenlois.ac.at/0.0.0.0\naddress=/.gamegogle.com/0.0.0.0\naddress=/.gasparini.com.br/0.0.0.0\naddress=/.getpics.net/0.0.0.0\naddress=/.gezila.com/0.0.0.0\naddress=/.glazeautocaremobile.com/0.0.0.0\naddress=/.goldenlifewomen.com/0.0.0.0\naddress=/.goosai.com/0.0.0.0\naddress=/.holidaysinkeralam.com/0.0.0.0\naddress=/.hotlaps.com.au/0.0.0.0\naddress=/.i2cchip.com/0.0.0.0\naddress=/.ibxdnl.com/0.0.0.0\naddress=/.igetmyservice.com/0.0.0.0\naddress=/.iprojhq.com/0.0.0.0\naddress=/.izmirhavaalaniarackiralama.net/0.0.0.0\naddress=/.jingshang.com.tw/0.0.0.0\naddress=/.justgetitfaster.com/0.0.0.0\naddress=/.kanberdemir.com/0.0.0.0\naddress=/.kpzip.com/0.0.0.0\naddress=/.kraonkelaere.com/0.0.0.0\naddress=/.laptopb4you.com/0.0.0.0\naddress=/.liftune.com/0.0.0.0\naddress=/.m-games.huu.cz/0.0.0.0\naddress=/.martiniracing.com.br/0.0.0.0\naddress=/.mireene.com/0.0.0.0\naddress=/.mixtrio.net/0.0.0.0\naddress=/.mstdls.com/0.0.0.0\naddress=/.mypcapp.com/0.0.0.0\naddress=/.perso.sfr.fr/0.0.0.0\naddress=/.pixelmon-world.com/0.0.0.0\naddress=/.plexcera.com/0.0.0.0\naddress=/.rd1994.com/0.0.0.0\naddress=/.sf-addon.com/0.0.0.0\naddress=/.skypedong.com/0.0.0.0\naddress=/.spirlymo.com/0.0.0.0\naddress=/.sportstherapy.net/0.0.0.0\naddress=/.talka-studios.com/0.0.0.0\naddress=/.thewitchez-cafe.co.uk/0.0.0.0\naddress=/.tirekoypazari.com/0.0.0.0\naddress=/.updatestar.net/0.0.0.0\naddress=/.urban-garden.net/0.0.0.0\naddress=/.utilbada.com/0.0.0.0\naddress=/.utilcom.net/0.0.0.0\naddress=/.utiljoy.com/0.0.0.0\naddress=/.vim6.com/0.0.0.0\naddress=/.windows.net/0.0.0.0\naddress=/.xiazai4.net/0.0.0.0\naddress=/.xunyou.com/0.0.0.0\n"

	filesMin = "[\nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"../testdata/blist.hosts.src\"\nIP:\t \"10.10.10.10\"\nLtype:\t \"file\"\nName:\t \"/tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n]"
)

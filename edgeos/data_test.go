package edgeos

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"testing"

	"github.com/britannic/blacklist/regx"
	"github.com/britannic/blacklist/tdata"
	. "github.com/britannic/testutils"
)

func TestDiffArray(t *testing.T) {
	biggest := []string{"one", "two", "three", "four", "five", "six"}
	smallest := []string{"one", "two", "three"}
	want := []string{"five", "four", "six"}

	got := DiffArray(biggest, smallest)
	Equals(t, want, got)

	got = DiffArray(smallest, biggest)
	Equals(t, want, got)

	shuffleArray(biggest)
	got = DiffArray(smallest, biggest)
	Equals(t, want, got)

	shuffleArray(smallest)
	got = DiffArray(smallest, biggest)
	Equals(t, want, got)
}

func TestFormatData(t *testing.T) {
	n, err := ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)

	c := n.NewConfig()

	for node := range n {
		var (
			got       io.Reader
			gotList   = make(List)
			lines     []string
			wantBytes []byte
			wantList  = make(List)
		)
		eq := getSeparator(node)

		for _, k := range n[node].Includes {
			wantList[k] = 0
			lines = append(lines, fmt.Sprintf("address=%v%v/%v\n", eq, k, c.Get(node).IP))
		}

		sort.Strings(lines)
		for _, line := range lines {
			wantBytes = append(wantBytes, line...)
		}

		fmttr := "address=" + eq + "%v/" + c.Get(node).IP + "\n"
		got, gotList = c.FormatData(fmttr, c.Get(node).Inc)

		gotBytes, err := ioutil.ReadAll(got)
		OK(t, err)
		Equals(t, string(wantBytes[:]), string(gotBytes[:]))
		Equals(t, wantList, gotList)
	}
}

func TestGet(t *testing.T) {
	// reader, err := os.Open("/Users/Neil/EdgeOs/Config.boots/config.boot")
	// n, err := ReadCfg(reader)
	n, err := ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)

	c := n.NewConfig()
	for node := range n {
		Equals(t, n[node].IP, c.Get(node).IP)
		Equals(t, n[node].Excludes, c.Get(node).Exc)
		Equals(t, n[node].Includes, c.Get(node).Inc)
		Equals(t, n.getLeaves(node), c.Get(node).Nodes)
		Equals(t, n.getSrcs(node), c.Get(node).Sources)
	}

	n, err = ReadCfg(bytes.NewBufferString(noIPCfg))
	OK(t, err)

	c = n.NewConfig()
	for _, node := range []string{"domains", "hosts"} {
		Equals(t, "192.168.1.1", c.Get(node).IP)
	}
}

func TestGetExcludes(t *testing.T) {
	want := List{
		"122.2o7.net":             0,
		"1e100.net":               0,
		"adobedtm.com":            0,
		"akamai.net":              0,
		"amazon.com":              0,
		"amazonaws.com":           0,
		"apple.com":               0,
		"ask.com":                 0,
		"avast.com":               0,
		"bitdefender.com":         0,
		"cdn.visiblemeasures.com": 0,
		"cloudfront.net":          0,
		"coremetrics.com":         0,
		"edgesuite.net":           0,
		"freedns.afraid.org":      0,
		"github.com":              0,
		"githubusercontent.com":   0,
		"google.com":              0,
		"googleadservices.com":    0,
		"googleapis.com":          0,
		"googleusercontent.com":   0,
		"gstatic.com":             0,
		"gvt1.com":                0,
		"gvt1.net":                0,
		"hb.disney.go.com":        0,
		"hp.com":                  0,
		"hulu.com":                0,
		"images-amazon.com":       0,
		"msdn.com":                0,
		"paypal.com":              0,
		"rackcdn.com":             0,
		"schema.org":              0,
		"skype.com":               0,
		"smacargo.com":            0,
		"sourceforge.net":         0,
		"ssl-on9.com":             0,
		"ssl-on9.net":             0,
		"static.chartbeat.com":    0,
		"storage.googleapis.com":  0,
		"windows.net":             0,
		"yimg.com":                0,
		"ytimg.com":               0,
	}
	n, err := ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)

	domEx := "big.bopper.com"
	c := n.NewConfig()
	c[Domains].Exc = []string{
		domEx,
	}

	dex := List{domEx: 0}
	got := make(List)
	dex, got = c.GetExcludes(dex, got, []string{blacklist, Domains, Hosts})

	Equals(t, want, got)
	Equals(t, List{domEx: 0}, dex)
}

func TestFiles(t *testing.T) {
	n, err := ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)

	c := n.NewConfig()
	dir := "/tmp"
	got := c.Files(dir, []string{blacklist, Domains, Hosts})
	want := []string{
		dir + "/domains.malc0de.blacklist.conf",
		dir + "/domains.pre-configured.blacklist.conf",
		dir + "/hosts.adaway.blacklist.conf",
		dir + "/hosts.malwaredomainlist.blacklist.conf",
		dir + "/hosts.openphish.blacklist.conf",
		dir + "/hosts.pre-configured.blacklist.conf",
		dir + "/hosts.someonewhocares.blacklist.conf",
		dir + "/hosts.tasty.blacklist.conf",
		dir + "/hosts.volkerschatz.blacklist.conf",
		dir + "/hosts.winhelp2002.blacklist.conf",
		dir + "/hosts.yoyo.blacklist.conf",
	}
	Equals(t, want, got)
}

func TestGetLeaves(t *testing.T) {
	n, err := ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)
	got := n.getLeaves("all")
	want := []Leaf{}

	for _, k := range n.SortKeys() {
		want = append(want, (*n[k]))
	}

	Equals(t, want, got)
}

func TestProcess(t *testing.T) {
	for _, test := range ProcessTests {
		// if test.num != 5 {
		// 	continue
		// }
		h := new(HTTPserver)
		URL := h.NewHTTPServer().String()
		h.mux.HandleFunc(test.page,
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, test.data)
			},
		)
		reader, err := GetHTTP(test.method, URL+test.page)
		switch test.ok {
		case true:
			OK(t, err)
		default:
			NotOK(t, err)
		}
		// fmt.Println(test.num)
		got := Process(test.source, test.dex, test.ex, reader).List
		Equals(t, test.want, got)
		// fmt.Println(got)
	}
}

func TestStripPrefixAndSuffix(t *testing.T) {
	rx := regx.Objects
	want := "This is a complete sentence and should not have a comment."

	for _, s := range src {
		var (
			got string
			ok  bool
		)
		switch s.Prefix {
		case "http":
			got = s.Prefix + "://" + want
		default:
			got = s.Prefix + want
		}

		got += " # Comment."

		got, ok = stripPrefixAndSuffix(got, s.Prefix, rx)
		Equals(t, true, ok)
		Equals(t, want, got)
	}

	// Break "http"
	want = "http\r://really.not.correct.com"
	got, ok := stripPrefixAndSuffix(want, "http", rx)
	Equals(t, false, ok)
	Equals(t, want, got)
}

func TestWriteIncludes(t *testing.T) {
	wantDex := List{
		"adsrvr.org":         0,
		"adtechus.net":       0,
		"advertising.com":    0,
		"centade.com":        0,
		"doubleclick.net":    0,
		"free-counter.co.uk": 0,
		"intellitxt.com":     0,
		"kiosked.com":        0,
	}

	wantEx := List{
		"beap.gemini.yahoo.com": 0,
	}

	nodes, err := ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)

	c := nodes.NewConfig()

	gotDex, gotEx := c.WriteIncludes("/tmp", []string{blacklist, Domains, Hosts})

	Equals(t, wantEx, gotEx)
	Equals(t, wantDex, gotDex)
}

var (
	ProcessTests = []struct {
		data   string
		dex    List
		ex     List
		method string
		num    int
		ok     bool
		page   string
		source *Srcs
		want   List
	}{
		{ // 1
			data:   httpHostData,
			ex:     make(List),
			dex:    make(List),
			method: method,
			num:    1,
			ok:     true,
			page:   "/hosts.txt",
			source: &Srcs{
				List:   make(List),
				Prefix: "127.0.0.1 ",
			},
			want: List{
				"a.applovin.com":                0,
				"a.glcdn.co":                    0,
				"a.vserv.mobi":                  0,
				"ad.leadboltapps.net":           0,
				"ad.madvertise.de":              0,
				"ad.where.com":                  0,
				"ad1.adinfuse.com":              0,
				"ad2.adinfuse.com":              0,
				"adcontent.saymedia.com":        0,
				"adinfuse.com":                  0,
				"admicro1.vcmedia.vn":           0,
				"admicro2.vcmedia.vn":           0,
				"admin.vserv.mobi":              0,
				"ads.adiquity.com":              0,
				"ads.admarvel.com":              0,
				"ads.admoda.com":                0,
				"ads.celtra.com":                0,
				"ads.flurry.com":                0,
				"ads.matomymobile.com":          0,
				"ads.mobgold.com":               0,
				"ads.mobilityware.com":          0,
				"ads.mopub.com":                 0,
				"ads.n-ws.org":                  0,
				"ads.ookla.com":                 0,
				"ads.saymedia.com":              0,
				"ads.smartdevicemedia.com":      0,
				"ads.vserv.mobi":                0,
				"ads.xxxad.net":                 0,
				"ads2.mediaarmor.com":           0,
				"adserver.ubiyoo.com":           0,
				"adultmoda.com":                 0,
				"android-sdk31.transpera.com":   0,
				"android.bcfads.com":            0,
				"api.airpush.com":               0,
				"api.analytics.omgpop.com":      0,
				"api.yp.com":                    0,
				"apps.buzzcity.net":             0,
				"apps.mobilityware.com":         0,
				"as.adfonic.net":                0,
				"asotrack1.fluentmobile.com":    0,
				"assets.cntdy.mobi":             0,
				"atti.velti.com":                0,
				"b.scorecardresearch.com":       0,
				"banners.bigmobileads.com":      0,
				"bigmobileads.com":              0,
				"bo.jumptap.com":                0,
				"bos-tapreq01.jumptap.com":      0,
				"bos-tapreq02.jumptap.com":      0,
				"bos-tapreq03.jumptap.com":      0,
				"bos-tapreq04.jumptap.com":      0,
				"bos-tapreq05.jumptap.com":      0,
				"bos-tapreq06.jumptap.com":      0,
				"bos-tapreq07.jumptap.com":      0,
				"bos-tapreq08.jumptap.com":      0,
				"bos-tapreq09.jumptap.com":      0,
				"bos-tapreq10.jumptap.com":      0,
				"bos-tapreq11.jumptap.com":      0,
				"bos-tapreq12.jumptap.com":      0,
				"bos-tapreq13.jumptap.com":      0,
				"bos-tapreq14.jumptap.com":      0,
				"bos-tapreq15.jumptap.com":      0,
				"bos-tapreq16.jumptap.com":      0,
				"bos-tapreq17.jumptap.com":      0,
				"bos-tapreq18.jumptap.com":      0,
				"bos-tapreq19.jumptap.com":      0,
				"bos-tapreq20.jumptap.com":      0,
				"c.vrvm.com":                    0,
				"c.vserv.mobi":                  0,
				"c753738.r38.cf2.rackcdn.com":   0,
				"cache-ssl.celtra.com":          0,
				"cache.celtra.com":              0,
				"cdn.celtra.com":                0,
				"cdn.nearbyad.com":              0,
				"cdn.trafficforce.com":          0,
				"cdn.us.goldspotmedia.com":      0,
				"cdn.vdopia.com":                0,
				"cdn1.crispadvertising.com":     0,
				"cdn1.inner-active.mobi":        0,
				"cdn2.crispadvertising.com":     0,
				"click.buzzcity.net":            0,
				"creative1cdn.mobfox.com":       0,
				"d.applovin.com":                0,
				"d2bgg7rjywcwsy.cloudfront.net": 0,
				"d3anogn3pbtk4v.cloudfront.net": 0,
				"d3oltyb66oj2v8.cloudfront.net": 0,
				"edge.reporo.net":               0,
				"ftpcontent.worldnow.com":       0,
				"funnel0.adinfuse.com":          0,
				"gemini.yahoo.com":              0,
				"go.adinfuse.com":               0,
				"go.mobpartner.mobi":            0,
				"go.vrvm.com":                   0,
				"gsmtop.net":                    0,
				"gts-ads.twistbox.com":          0,
				"hhbekxxw5d9e.pflexads.com":     1,
				"hybl9bazbc35.pflexads.com":     0,
				"i.jumptap.com":                 0,
				"i.tapit.com":                   0,
				"images.millennialmedia.com":    0,
				"images.mpression.net":          0,
				"img.ads.huntmad.com":           0,
				"img.ads.mobilefuse.net":        0,
				"img.ads.mocean.mobi":           0,
				"img.ads.mojiva.com":            0,
				"img.ads.taptapnetworks.com":    0,
				"intouch.adinfuse.com":          0,
				"lb.usemaxserver.de":            0,
				"m.adsymptotic.com":             0,
				"m2m1.inner-active.mobi":        0,
				"media.mobpartner.mobi":         0,
				"medrx.sensis.com.au":           0,
				"mobile.banzai.it":              0,
				"mobiledl.adboe.com":            0,
				"mobpartner.mobi":               0,
				"mwc.velti.com":                 0,
				"netdna.reporo.net":             0,
				"oasc04012.247realmedia.com":    0,
				"orange-fr.adinfuse.com":        0,
				"orangeuk-mc.adinfuse.com":      0,
				"orencia.pflexads.com":          0,
				"pdn.applovin.com":              0,
				"r.edge.inmobicdn.net":          0,
				"r.mobpartner.mobi":             0,
				"req.appads.com":                0,
				"rs-staticart.ybcdn.net":        0,
				"ru.velti.com":                  0,
				"s0.2mdn.net":                   0,
				"s3.phluant.com":                0,
				"sf.vserv.mobi":                 0,
				"show.buzzcity.net":             0,
				"sky-connect.adinfuse.com":      0,
				"sky.adinfuse.com":              0,
				"static.cdn.gtsmobi.com":        0,
				"static.estebull.com":           0,
				"stats.pflexads.com":            0,
				"track.celtra.com":              0,
				"tracking.klickthru.com":        0,
				"uk-ad2.adinfuse.com":           0,
				"uk-go.adinfuse.com":            0,
				"web63.jumptap.com":             0,
				"web64.jumptap.com":             0,
				"web65.jumptap.com":             0,
				"wv.inner-active.mobi":          0,
				"www.eltrafiko.com":             0,
				"www.mmnetwork.mobi":            0,
				"www.pflexads.com":              0,
				"wwww.adleads.com":              0,
			},
		},
		{ // 2
			data:   httpDomainData,
			ex:     make(List),
			dex:    make(List),
			method: method,
			num:    2,
			ok:     true,
			page:   "/domains.txt",
			source: &Srcs{
				List:   make(List),
				Prefix: "zone ",
			},
			want: List{
				"192-168-0-255.com":          0,
				"asi-37.fr":                  0,
				"bagbackpack.com":            0,
				"bitmeyenkartusistanbul.com": 0,
				"byxon.com":                  0,
				"img001.com":                 0,
				"loadto.net":                 0,
				"roastfiles2017.com":         0,
			},
		},
		{ // 3
			data: `# [General]
  					127.0.0.1 a.glcdn.co
  					127.0.0.1 ad.madvertise.de
  					127.0.0.1 adcontent.saymedia.com
  					127.0.0.1 admicro1.vcmedia.vn
  					127.0.0.1 admicro2.vcmedia.vn
  					127.0.0.1 ads.admoda.com
  					127.0.0.1 ads.mobgold.com
  					127.0.0.1 ads.mopub.com
  					127.0.0.1 ads.saymedia.com
  					127.0.0.1 ads.xxxad.net
  					127.0.0.1 android.bcfads.com
  					127.0.0.1 api.analytics.omgpop.com
  					127.0.0.1 apps.buzzcity.net
  					127.0.0.1 assets.cntdy.mobi
  					127.0.0.1 banners.bigmobileads.com
  					127.0.0.1 bigmobileads.com
  					127.0.0.1 c.vrvm.com
  					127.0.0.1 click.buzzcity.net
  					127.0.0.1 creative1cdn.mobfox.com
  					127.0.0.1 ftpcontent.worldnow.com
  					127.0.0.1 go.vrvm.com
  					127.0.0.1 gsmtop.net
  					127.0.0.1 hhbekxxw5d9e.pflexads.com
  					127.0.0.1 images.millennialmedia.com
  					127.0.0.1 images.mpression.net
  					127.0.0.1 img.ads.huntmad.com
  					127.0.0.1 img.ads.mocean.mobi
  					127.0.0.1 img.ads.mojiva.com
  					127.0.0.1 lb.usemaxserver.de
  					127.0.0.1 mobile.Banzai.it
  					127.0.0.1 r.mobpartner.mobi,
  					127.0.0.1 oasc04012.247realmedia.com
  					127.0.0.1 orencia.pflexads.com
  					127.0.0.1 pdn.applovin.com
  					127.0.0.1 r.edge.inmobicdn.net
  					127.0.0.1 req.appads.com
  					127.0.0.1 s0.2mdn.net
  					127.0.0.1 s3.phluant.com
  					127.0.0.1 show.buzzcity.net
  					127.0.0.1 static.estebull.com
  					127.0.0.1 stats.pflexads.com
  					127.0.0.1 tracking.klickthru.com
  					127.0.0.1 www.mmnetwork.mobi
  					127.0.0.1 wwww.adleads.com
  					`,
			ex: List{
				"mwc.velti.com":              0,
				"netdna.reporo.net":          0,
				"oasc04012.247realmedia.com": 0,
				"orange-fr.adinfuse.com":     0,
				"orangeuk-mc.adinfuse.com":   0,
				"orencia.pflexads.com":       0,
				"pdn.applovin.com":           0,
				"r.edge.inmobicdn.net":       0,
				"req.appads.com":             0,
				"rs-staticart.ybcdn.net":     0,
				"ru.velti.com":               0,
				"s0.2mdn.net":                0,
				"s3.phluant.com":             0,
				"sf.vserv.mobi":              0,
				"show.buzzcity.net":          0,
				"sky-connect.adinfuse.com":   0,
				"sky.adinfuse.com":           0,
				"static.cdn.gtsmobi.com":     0,
				"static.estebull.com":        0,
			},
			dex: List{
				"adleads.com":     0,
				"mobpartner.mobi": 0,
				"phluant.com":     0,
			},
			method: method,
			num:    3,
			ok:     true,
			page:   "/hosts.txt",
			source: &Srcs{
				List: List{
					"pdn.applovin.com":       0,
					"r.edge.inmobicdn.net":   0,
					"r.mobpartner.mobi":      0,
					"req.appads.com":         0,
					"rs-staticart.ybcdn.net": 0,
					"ru.velti.com":           0,
					"s0.2mdn.net":            0,
				},
				Prefix: "127.0.0.1 ",
			},
			want: List{
				"a.glcdn.co":                 0,
				"ad.madvertise.de":           0,
				"adcontent.saymedia.com":     0,
				"admicro1.vcmedia.vn":        0,
				"admicro2.vcmedia.vn":        0,
				"ads.admoda.com":             0,
				"ads.mobgold.com":            0,
				"ads.mopub.com":              0,
				"ads.saymedia.com":           0,
				"ads.xxxad.net":              0,
				"android.bcfads.com":         0,
				"api.analytics.omgpop.com":   0,
				"apps.buzzcity.net":          0,
				"assets.cntdy.mobi":          0,
				"banners.bigmobileads.com":   0,
				"bigmobileads.com":           0,
				"c.vrvm.com":                 0,
				"click.buzzcity.net":         0,
				"creative1cdn.mobfox.com":    0,
				"ftpcontent.worldnow.com":    0,
				"go.vrvm.com":                0,
				"gsmtop.net":                 0,
				"hhbekxxw5d9e.pflexads.com":  0,
				"images.millennialmedia.com": 0,
				"images.mpression.net":       0,
				"img.ads.huntmad.com":        0,
				"img.ads.mocean.mobi":        0,
				"img.ads.mojiva.com":         0,
				"lb.usemaxserver.de":         0,
				"mobile.banzai.it":           0,
				"pdn.applovin.com":           1,
				"r.edge.inmobicdn.net":       1,
				"r.mobpartner.mobi":          0,
				"req.appads.com":             1,
				"rs-staticart.ybcdn.net":     0,
				"ru.velti.com":               0,
				"s0.2mdn.net":                1,
				"stats.pflexads.com":         0,
				"tracking.klickthru.com":     0,
				"www.mmnetwork.mobi":         0,
			},
		},
		{ // 4
			data: `zone "192-168-0-255.com"  {type master; file "/etc/namedb/blockeddomain.hosts";};
  				zone "asi-37.fr"  {type master; file "/etc/namedb/blockeddomain.hosts";};
  				zone "bagbackpack.com"  {type master; file "/etc/namedb/blockeddomain.hosts";};
  				zone "bitmeyenkartusistanbul.com"  {type master; file "/etc/namedb/blockeddomain.hosts";};
  				zone "bitmeyenkartusistanbul.com"  {type master; file "/etc/namedb/blockeddomain.hosts";};
  							zone "bitmeyenkartusistanbul.com"  {type master; file "/etc/namedb/blockeddomain.hosts";};`,
			ex: List{
				"bagbackpack.com": 0,
				"asi-37.fr":       0,
			},
			dex:    make(List),
			method: method,
			num:    4,
			ok:     true,
			page:   "/domains.txt",
			source: &Srcs{
				List:   make(List),
				Prefix: "zone ",
			},
			want: List{
				"192-168-0-255.com":          0,
				"bitmeyenkartusistanbul.com": 2,
			},
		},
		{ // 5
			data: `# [General]
  					http:/a.glcdn.co
  					http:/ad.madvertise.de
  					http:/adcontent.saymedia.com
  					http:/adcontent.saymedia.com
  					http:/img.ads.huntmad.com
  					http:/creative1cdn.mobfox.com
  					http://oasc04012.247realmedia.com
  					http:/smaato.net
  					http:/c47.smaato.net
  					http:/c48.smaato.net
  					http:/c49.smaato.net
  					http:/c50.smaato.net
  					http:/c51.smaato.net
  					http:/c52.smaato.net
  					http:/c53.smaato.net
  					http:/c54.smaato.net
  					http:/c55.smaato.net
  					http:/c56.smaato.net
  					http:/c57.smaato.net
  					http:/c58.smaato.net
  					http:/c59.smaato.net
  					http:/c60.smaato.net
  					http:/f03.smaato.net`,
			ex: List{"a.glcdn.co": 0},
			dex: List{
				"saymedia.com": 0,
				"smaato.net":   0,
			},
			method: method,
			num:    5,
			ok:     true,
			page:   "/hosts.txt",
			source: &Srcs{
				List: List{
					"creative1cdn.mobfox.com":    0,
					"oasc04012.247realmedia.com": 0,
				},
				Prefix: "http",
			},
			want: List{
				"ad.madvertise.de":           0,
				"creative1cdn.mobfox.com":    1,
				"img.ads.huntmad.com":        0,
				"oasc04012.247realmedia.com": 1,
			},
		},
	}

	src = []*Srcs{
		{
			Disabled: false,
			IP:       "0.0.0.0",
			Name:     "malc0de",
			Prefix:   "zone ",
			Type:     domain,
			URL:      "http://malc0de.com/bl/ZONES",
		},
		{
			Disabled: false,
			IP:       "0.0.0.0",
			Name:     "adaway",
			Prefix:   "127.0.0.1 ",
			Type:     host,
			URL:      "http://adaway.org/hosts.txt",
		},
		{
			Disabled: false,
			IP:       "0.0.0.0",
			Name:     "malwaredomainlist",
			Prefix:   "127.0.0.1 ",
			Type:     host,
			URL:      "http://www.malwaredomainlist.com/hostslist/hosts.txt",
		},
		{
			Disabled: false,
			IP:       "0.0.0.0",
			Name:     "openphish",
			Prefix:   "http",
			Type:     host,
			URL:      "https://openphish.com/feed.txt",
		},
		{
			Disabled: false,
			IP:       "0.0.0.0",
			Name:     "someonewhocares",
			Prefix:   "0.0.0.0",
			Type:     host,
			URL:      "http://someonewhocares.org/hosts/zero/",
		},
		{
			Name:     "tasty",
			Disabled: false,
			File:     "/config/user-data/blist.hosts.src",
			IP:       "",
			Prefix:   "",
			Type:     host,
			URL:      "",
		},
		{
			Disabled: false,
			IP:       "0.0.0.0",
			Name:     "volkerschatz",
			Prefix:   "http",
			Type:     host,
			URL:      "http://www.volkerschatz.com/net/adpaths",
		},
		{
			Disabled: false,
			IP:       "0.0.0.0",
			Name:     "winhelp2002",
			Prefix:   "0.0.0.0 ",
			Type:     host,
			URL:      "http://winhelp2002.mvps.org/hosts.txt",
		},
		{
			Disabled: false,
			IP:       "0.0.0.0",
			Name:     "yoyo",
			Type:     host,
			URL:      "http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext",
		},
	}

	noIPCfg = `blacklist {
	disabled false
	dns-redirect-ip 192.168.1.1
	domains {
	disabled false
			dns-redirect-ip
			source malc0de {
					description "List of zones serving malicious executables observed by malc0de.com/database/"
					dns-redirect-ip 4.4.4.4
					prefix "zone "
					url http://malc0de.com/bl/ZONES
			}
	}
	hosts {
			disabled true
			dns-redirect-ip
			source file {
				description "File test"
				dns-redirect-ip 5.5.5.5
				file /test/file
		}
	}
}`
)

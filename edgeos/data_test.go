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

	for _, node := range []string{Domains, Hosts} {
		var (
			got       io.Reader
			gotList   = make(List)
			lines     []string
			wantBytes []byte
		)
		eq := GetSeparator(node)

		for _, k := range n[node].Includes {
			lines = append(lines, fmt.Sprintf("address=%v%v/%v", eq, k, c.Get(node).IP)+"\n")
			gotList[k] = 0
		}

		sort.Strings(lines)
		for _, line := range lines {
			wantBytes = append(wantBytes, line...)
		}

		fmttr := "address=" + eq + "%v/" + c.Get(node).IP
		got = FormatData(fmttr, gotList)

		gotBytes, err := ioutil.ReadAll(got)
		OK(t, err)
		Equals(t, wantBytes[:], gotBytes[:])
		// fmt.Println(string(gotBytes[:]))
	}
}

func TestGet(t *testing.T) {
	// reader, err := os.Open("/Users/Neil/EdgeOs/Config.boots/config.boot")
	// n, err := ReadCfg(reader)
	n, err := ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)

	c := n.NewConfig()
	for node := range n {
		switch {
		case node != Root && n[node].IP != "":
			Equals(t, n[node].IP, c.IP(node))
			ip := c.o[node].IP
			c.o[node].IP = ""
			Equals(t, n[Root].IP, c.IP(node))
			c.o[node].IP = ip
		case node != Root:
			Equals(t, Srcs{Name: PreConf, Disabled: false, Type: getType(node).(int), No: 0}, c.Source(node, PreConf))
		}

		inc := []byte{}
		for _, k := range n[node].Includes {
			inc = append(inc, []byte(k+"\n")...)
		}

		switch node {
		case Domains, Hosts:
			Equals(t, bytes.NewBuffer(inc), c.GetIncludes(node))
		}

		Equals(t, n[node].Excludes, c.Get(node).Exc)
		Equals(t, c.Get(node).Exc, c.Excludes(node))
		Equals(t, n[node].Includes, c.Get(node).Inc)
		Equals(t, c.Get(node).Inc, c.Includes(node))
		Equals(t, n.getSrcs(node), c.Sources(node))
		Equals(t, n.getSrcs(node), c.Get(node).Sources)
		Equals(t, c.Get(node).Disabled, c.Disabled(node))
	}

	n, err = ReadCfg(bytes.NewBufferString(noIPCfg))
	OK(t, err)

	c = n.NewConfig()
	for _, node := range []string{"domains", "hosts"} {
		Equals(t, "192.168.1.1", c.Get(node).IP)
	}
}

func TestGetExcludes(t *testing.T) {
	want := testMap
	n, err := ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)

	domEx := "big.bopper.com"
	c := n.NewConfig()
	c.o[Domains].Exc = []string{
		domEx,
	}

	c.dex = List{domEx: 0}
	c.GetExcludes([]string{blacklist, Domains, Hosts})

	Equals(t, want, c.ex)
	Equals(t, List{domEx: 0}, c.dex)
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
	n, err := ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)
	c := n.NewConfig()
	for _, test := range ProcessTests {
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

		c.dex = test.dex
		c.ex = test.ex
		got := c.Process(test.source, reader)
		Equals(t, test.want, got)
	}
}

func TestSource(t *testing.T) {
	n, err := ReadCfg(bytes.NewBufferString(tdata.Cfg))
	OK(t, err)

	c := n.NewConfig()

	test := map[string]map[string][]Srcs{"domains": map[string][]Srcs{"pre-configured": []Srcs{Srcs{Desc: "", Disabled: false, File: "", IP: "0.0.0.0", List: List(nil), Name: "pre-configured", No: 0, Prefix: "", Type: 2, URL: ""}}, "urls": []Srcs{Srcs{Desc: "List of zones serving malicious executables observed by malc0de.com/database/", Disabled: false, File: "", IP: "", List: List(nil), Name: "malc0de", No: 0, Prefix: "zone ", Type: 2, URL: "http://malc0de.com/bl/ZONES"}, Srcs{Desc: "", Disabled: false, File: "", IP: "", List: List(nil), Name: "pre-configured", No: 0, Prefix: "", Type: 2, URL: ""}}, "files": []Srcs{}}, "hosts": map[string][]Srcs{"urls": []Srcs{Srcs{Desc: "Blocking mobile ad providers and some analytics providers", Disabled: false, File: "", IP: "", List: List(nil), Name: "adaway", No: 0, Prefix: "127.0.0.1 ", Type: 3, URL: "http://adaway.org/hosts.txt"}, Srcs{Desc: "127.0.0.1 based host and domain list", Disabled: false, File: "", IP: "", List: List(nil), Name: "malwaredomainlist", No: 0, Prefix: "127.0.0.1 ", Type: 3, URL: "http://www.malwaredomainlist.com/hostslist/hosts.txt"}, Srcs{Desc: "OpenPhish automatic phishing detection", Disabled: false, File: "", IP: "", List: List(nil), Name: "openphish", No: 0, Prefix: "http", Type: 3, URL: "https://openphish.com/feed.txt"}, Srcs{Desc: "", Disabled: false, File: "", IP: "", List: List(nil), Name: "pre-configured", No: 0, Prefix: "", Type: 3, URL: ""}, Srcs{Desc: "Zero based host and domain list", Disabled: false, File: "", IP: "", List: List(nil), Name: "someonewhocares", No: 0, Prefix: "0.0.0.0", Type: 3, URL: "http://someonewhocares.org/hosts/zero/"}, Srcs{Desc: "File source", Disabled: false, File: "/config/user-data/blist.hosts.src", IP: "0.0.0.0", List: List(nil), Name: "tasty", No: 0, Prefix: "", Type: 3, URL: ""}, Srcs{Desc: "Ad server blacklists", Disabled: false, File: "", IP: "", List: List(nil), Name: "volkerschatz", No: 0, Prefix: "http", Type: 3, URL: "http://www.volkerschatz.com/net/adpaths"}, Srcs{Desc: "Zero based host and domain list", Disabled: false, File: "", IP: "0.0.0.0", List: List(nil), Name: "winhelp2002", No: 0, Prefix: "0.0.0.0 ", Type: 3, URL: "http://winhelp2002.mvps.org/hosts.txt"}, Srcs{Desc: "Fully Qualified Domain Names only - no prefix to strip", Disabled: false, File: "", IP: "", List: List(nil), Name: "yoyo", No: 0, Prefix: "", Type: 3, URL: "http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext"}}, "files": []Srcs{Srcs{Desc: "File source", Disabled: false, File: "/config/user-data/blist.hosts.src", IP: "0.0.0.0", List: List(nil), Name: "tasty", No: 0, Prefix: "", Type: 3, URL: ""}}}}

	for node := range n {
		switch node {
		case Domains, Hosts:
			for _, stype := range []string{PreConf, "files", "urls"} {
				test[node][stype] = c.Source(node, stype)
				Equals(t, test[node][stype], c.Source(node, stype))
			}
		}
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

func TestUpdateList(t *testing.T) {
	var (
		got  = make(List)
		want = testMap
	)

	UpdateList(testArray, got)
	Equals(t, want, got)
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
		want   io.Reader
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
			want: bytes.NewBuffer([]byte(`address=/a.applovin.com/
address=/a.glcdn.co/
address=/a.vserv.mobi/
address=/ad.leadboltapps.net/
address=/ad.madvertise.de/
address=/ad.where.com/
address=/ad1.adinfuse.com/
address=/ad2.adinfuse.com/
address=/adcontent.saymedia.com/
address=/adinfuse.com/
address=/admicro1.vcmedia.vn/
address=/admicro2.vcmedia.vn/
address=/admin.vserv.mobi/
address=/ads.adiquity.com/
address=/ads.admarvel.com/
address=/ads.admoda.com/
address=/ads.celtra.com/
address=/ads.flurry.com/
address=/ads.matomymobile.com/
address=/ads.mobgold.com/
address=/ads.mobilityware.com/
address=/ads.mopub.com/
address=/ads.n-ws.org/
address=/ads.ookla.com/
address=/ads.saymedia.com/
address=/ads.smartdevicemedia.com/
address=/ads.vserv.mobi/
address=/ads.xxxad.net/
address=/ads2.mediaarmor.com/
address=/adserver.ubiyoo.com/
address=/adultmoda.com/
address=/android-sdk31.transpera.com/
address=/android.bcfads.com/
address=/api.airpush.com/
address=/api.analytics.omgpop.com/
address=/api.yp.com/
address=/apps.buzzcity.net/
address=/apps.mobilityware.com/
address=/as.adfonic.net/
address=/asotrack1.fluentmobile.com/
address=/assets.cntdy.mobi/
address=/atti.velti.com/
address=/b.scorecardresearch.com/
address=/banners.bigmobileads.com/
address=/bigmobileads.com/
address=/bo.jumptap.com/
address=/bos-tapreq01.jumptap.com/
address=/bos-tapreq02.jumptap.com/
address=/bos-tapreq03.jumptap.com/
address=/bos-tapreq04.jumptap.com/
address=/bos-tapreq05.jumptap.com/
address=/bos-tapreq06.jumptap.com/
address=/bos-tapreq07.jumptap.com/
address=/bos-tapreq08.jumptap.com/
address=/bos-tapreq09.jumptap.com/
address=/bos-tapreq10.jumptap.com/
address=/bos-tapreq11.jumptap.com/
address=/bos-tapreq12.jumptap.com/
address=/bos-tapreq13.jumptap.com/
address=/bos-tapreq14.jumptap.com/
address=/bos-tapreq15.jumptap.com/
address=/bos-tapreq16.jumptap.com/
address=/bos-tapreq17.jumptap.com/
address=/bos-tapreq18.jumptap.com/
address=/bos-tapreq19.jumptap.com/
address=/bos-tapreq20.jumptap.com/
address=/c.vrvm.com/
address=/c.vserv.mobi/
address=/c753738.r38.cf2.rackcdn.com/
address=/cache-ssl.celtra.com/
address=/cache.celtra.com/
address=/cdn.celtra.com/
address=/cdn.nearbyad.com/
address=/cdn.trafficforce.com/
address=/cdn.us.goldspotmedia.com/
address=/cdn.vdopia.com/
address=/cdn1.crispadvertising.com/
address=/cdn1.inner-active.mobi/
address=/cdn2.crispadvertising.com/
address=/click.buzzcity.net/
address=/creative1cdn.mobfox.com/
address=/d.applovin.com/
address=/d2bgg7rjywcwsy.cloudfront.net/
address=/d3anogn3pbtk4v.cloudfront.net/
address=/d3oltyb66oj2v8.cloudfront.net/
address=/edge.reporo.net/
address=/ftpcontent.worldnow.com/
address=/funnel0.adinfuse.com/
address=/gemini.yahoo.com/
address=/go.adinfuse.com/
address=/go.mobpartner.mobi/
address=/go.vrvm.com/
address=/gsmtop.net/
address=/gts-ads.twistbox.com/
address=/hhbekxxw5d9e.pflexads.com/
address=/hybl9bazbc35.pflexads.com/
address=/i.jumptap.com/
address=/i.tapit.com/
address=/images.millennialmedia.com/
address=/images.mpression.net/
address=/img.ads.huntmad.com/
address=/img.ads.mobilefuse.net/
address=/img.ads.mocean.mobi/
address=/img.ads.mojiva.com/
address=/img.ads.taptapnetworks.com/
address=/intouch.adinfuse.com/
address=/lb.usemaxserver.de/
address=/m.adsymptotic.com/
address=/m2m1.inner-active.mobi/
address=/media.mobpartner.mobi/
address=/medrx.sensis.com.au/
address=/mobile.banzai.it/
address=/mobiledl.adboe.com/
address=/mobpartner.mobi/
address=/mwc.velti.com/
address=/netdna.reporo.net/
address=/oasc04012.247realmedia.com/
address=/orange-fr.adinfuse.com/
address=/orangeuk-mc.adinfuse.com/
address=/orencia.pflexads.com/
address=/pdn.applovin.com/
address=/r.edge.inmobicdn.net/
address=/r.mobpartner.mobi/
address=/req.appads.com/
address=/rs-staticart.ybcdn.net/
address=/ru.velti.com/
address=/s0.2mdn.net/
address=/s3.phluant.com/
address=/sf.vserv.mobi/
address=/show.buzzcity.net/
address=/sky-connect.adinfuse.com/
address=/sky.adinfuse.com/
address=/static.cdn.gtsmobi.com/
address=/static.estebull.com/
address=/stats.pflexads.com/
address=/track.celtra.com/
address=/tracking.klickthru.com/
address=/uk-ad2.adinfuse.com/
address=/uk-go.adinfuse.com/
address=/web63.jumptap.com/
address=/web64.jumptap.com/
address=/web65.jumptap.com/
address=/wv.inner-active.mobi/
address=/www.eltrafiko.com/
address=/www.mmnetwork.mobi/
address=/www.pflexads.com/
address=/wwww.adleads.com/
`)),
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
			want: bytes.NewBuffer([]byte(`address=/192-168-0-255.com/
address=/asi-37.fr/
address=/bagbackpack.com/
address=/bitmeyenkartusistanbul.com/
address=/byxon.com/
address=/img001.com/
address=/loadto.net/
address=/roastfiles2017.com/
`)),
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
			want: bytes.NewBuffer([]byte(`address=/a.glcdn.co/
address=/ad.madvertise.de/
address=/adcontent.saymedia.com/
address=/admicro1.vcmedia.vn/
address=/admicro2.vcmedia.vn/
address=/ads.admoda.com/
address=/ads.mobgold.com/
address=/ads.mopub.com/
address=/ads.saymedia.com/
address=/ads.xxxad.net/
address=/android.bcfads.com/
address=/api.analytics.omgpop.com/
address=/apps.buzzcity.net/
address=/assets.cntdy.mobi/
address=/banners.bigmobileads.com/
address=/bigmobileads.com/
address=/c.vrvm.com/
address=/click.buzzcity.net/
address=/creative1cdn.mobfox.com/
address=/ftpcontent.worldnow.com/
address=/go.vrvm.com/
address=/gsmtop.net/
address=/hhbekxxw5d9e.pflexads.com/
address=/images.millennialmedia.com/
address=/images.mpression.net/
address=/img.ads.huntmad.com/
address=/img.ads.mocean.mobi/
address=/img.ads.mojiva.com/
address=/lb.usemaxserver.de/
address=/mobile.banzai.it/
address=/pdn.applovin.com/
address=/r.edge.inmobicdn.net/
address=/r.mobpartner.mobi/
address=/req.appads.com/
address=/rs-staticart.ybcdn.net/
address=/ru.velti.com/
address=/s0.2mdn.net/
address=/stats.pflexads.com/
address=/tracking.klickthru.com/
address=/www.mmnetwork.mobi/
`)),
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
			want: bytes.NewBuffer([]byte(`address=/192-168-0-255.com/
address=/bitmeyenkartusistanbul.com/
`)),
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
			want: bytes.NewBuffer([]byte(`address=/ad.madvertise.de/
address=/creative1cdn.mobfox.com/
address=/img.ads.huntmad.com/
address=/oasc04012.247realmedia.com/
`)),
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

	testArray = []string{
		"122.2o7.net",
		"1e100.net",
		"adobedtm.com",
		"akamai.net",
		"amazon.com",
		"amazonaws.com",
		"apple.com",
		"ask.com",
		"avast.com",
		"bitdefender.com",
		"cdn.visiblemeasures.com",
		"cloudfront.net",
		"coremetrics.com",
		"edgesuite.net",
		"freedns.afraid.org",
		"github.com",
		"githubusercontent.com",
		"google.com",
		"googleadservices.com",
		"googleapis.com",
		"googleusercontent.com",
		"gstatic.com",
		"gvt1.com",
		"gvt1.net",
		"hb.disney.go.com",
		"hp.com",
		"hulu.com",
		"images-amazon.com",
		"msdn.com",
		"paypal.com",
		"rackcdn.com",
		"schema.org",
		"skype.com",
		"smacargo.com",
		"sourceforge.net",
		"ssl-on9.com",
		"ssl-on9.net",
		"static.chartbeat.com",
		"storage.googleapis.com",
		"windows.net",
		"yimg.com",
		"ytimg.com",
	}

	testMap = List{
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
)

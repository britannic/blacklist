package edgeos

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/britannic/testutils"
)

type HTTPserver struct {
	mux    *http.ServeMux
	server *httptest.Server
}

func (h *HTTPserver) NewHTTPServer() *url.URL {
	// test server
	h.mux = http.NewServeMux()
	h.server = httptest.NewServer(h.mux)
	URL, _ := url.Parse(h.server.URL)
	return URL
}

func TestGetHTTP(t *testing.T) {
	var (
		got    []byte
		h      = new(HTTPserver)
		method = "GET"
		page   = "/domains.txt"
		want   = httpDomainData
	)

	URL := h.NewHTTPServer().String()
	h.mux.HandleFunc(page,
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, want)
		},
	)

	tests := []struct {
		ok     bool
		method string
		URL    string
		want   string
	}{
		{ok: true, method: method, URL: URL + page, want: want},
		{ok: false, method: method, URL: "bad url", want: "No data returned for bad url..."},
		{ok: false, method: "bad method", URL: "bad url", want: "Unable to form request for bad url..."},
	}

	for _, test := range tests {
		body, err := GetHTTP(test.method, test.URL)
		switch test.ok {
		case true:
			OK(t, err)
		default:
			NotOK(t, err)
		}

		got, err = ioutil.ReadAll(body)
		OK(t, err)
		Equals(t, test.want, string(got[:]))
	}
}

var (
	method = "GET"

	httpDomainData = `
// This bind zone is intended to be included in a running dns server for a local net
// It will return 127.0.0.1 for domains serving malicious executables observed by malc0de.com/database/
// This file will be automatically updated daily and populated with the last 30 days of malicious domains.
// Additional information to get this working can be found http://www.malwaredomains.com/wordpress/?page_id=6
// Last updated 2016-03-09

zone "192-168-0-255.com"  {type master; file "/etc/namedb/blockeddomain.hosts";};
zone "asi-37.fr"  {type master; file "/etc/namedb/blockeddomain.hosts";};
zone "bagbackpack.com"  {type master; file "/etc/namedb/blockeddomain.hosts";};
zone "bitmeyenkartusistanbul.com"  {type master; file "/etc/namedb/blockeddomain.hosts";};
zone "byxon.com"  {type master; file "/etc/namedb/blockeddomain.hosts";};
zone "img001.com"  {type master; file "/etc/namedb/blockeddomain.hosts";};
zone "loadto.net"  {type master; file "/etc/namedb/blockeddomain.hosts";};
zone "roastfiles2017.com"  {type master; file "/etc/namedb/blockeddomain.hosts";};`

	httpHostData = `# AdAway default blocklist
# Blocking mobile ad providers and some analytics providers
#
# Contribute:
# Create an issue at https://github.com/dschuermann/ad-away/issues
#
# Changelog:
# 2014-05-18 Now with a valid SSL certificate available at https://adaway.org/hosts.txt
# 2013-03-29 Integrated some hosts from
#            http://adblock.gjtech.net/?format=hostfile
# 2013-03-14 Back from the dead
#
# License:
# CC Attribution 3.0 (http://creativecommons.org/licenses/by/3.0/)
#
# Contributions by:
# Kicelo, Dominik Schuermann
#

127.0.0.1  localhost
::1  localhost

# [General]
127.0.0.1 lb.usemaxserver.de
127.0.0.1 tracking.klickthru.com
127.0.0.1 gsmtop.net
127.0.0.1 click.buzzcity.net
127.0.0.1 ads.admoda.com
127.0.0.1 stats.pflexads.com
127.0.0.1 a.glcdn.co
127.0.0.1 wwww.adleads.com
127.0.0.1 ad.madvertise.de
127.0.0.1 apps.buzzcity.net
127.0.0.1 ads.mobgold.com
127.0.0.1 android.bcfads.com
127.0.0.1 req.appads.com
127.0.0.1 show.buzzcity.net
127.0.0.1 api.analytics.omgpop.com
127.0.0.1 r.edge.inmobicdn.net
127.0.0.1 www.mmnetwork.mobi
127.0.0.1 img.ads.huntmad.com
127.0.0.1 creative1cdn.mobfox.com
127.0.0.1 admicro2.vcmedia.vn
127.0.0.1 admicro1.vcmedia.vn
127.0.0.1 s3.phluant.com
127.0.0.1 c.vrvm.com
127.0.0.1 go.vrvm.com
127.0.0.1 static.estebull.com
127.0.0.1 mobile.Banzai.it
127.0.0.1 ads.xxxad.net
127.0.0.1 hhbekxxw5d9e.pflexads.com
127.0.0.1 img.ads.mojiva.com
127.0.0.1 adcontent.saymedia.com
127.0.0.1 ads.saymedia.com
127.0.0.1 ftpcontent.worldnow.com
127.0.0.1 s0.2mdn.net
127.0.0.1 img.ads.mocean.mobi
127.0.0.1 bigmobileads.com
127.0.0.1 banners.bigmobileads.com
127.0.0.1 ads.mopub.com
127.0.0.1 images.mpression.net
127.0.0.1 images.millennialmedia.com
127.0.0.1 oasc04012.247realmedia.com
127.0.0.1 assets.cntdy.mobi
127.0.0.1 ad.leadboltapps.net ## another airpush style ad#
127.0.0.1 api.airpush.com ## hope this is all #
127.0.0.1 ad.where.com
127.0.0.1 i.tapit.com
127.0.0.1 cdn1.crispadvertising.com
127.0.0.1 cdn2.crispadvertising.com
127.0.0.1 medrx.sensis.com.au
127.0.0.1 rs-staticart.ybcdn.net
127.0.0.1 img.ads.taptapnetworks.com
127.0.0.1 adserver.ubiyoo.com
127.0.0.1 c753738.r38.cf2.rackcdn.com
127.0.0.1 edge.reporo.net
127.0.0.1 ads.n-ws.org
127.0.0.1 adultmoda.com
127.0.0.1 ads.smartdevicemedia.com
127.0.0.1 b.scorecardresearch.com
127.0.0.1 m.adsymptotic.com
127.0.0.1 cdn.vdopia.com
127.0.0.1 api.yp.com
127.0.0.1 asotrack1.fluentmobile.com
127.0.0.1 android-sdk31.transpera.com
127.0.0.1 apps.mobilityware.com
127.0.0.1 ads.mobilityware.com
127.0.0.1 ads.admarvel.com
127.0.0.1 netdna.reporo.net
127.0.0.1 www.eltrafiko.com
127.0.0.1 cdn.trafficforce.com
127.0.0.1 gts-ads.twistbox.com
127.0.0.1 static.cdn.gtsmobi.com
127.0.0.1 ads.matomymobile.com
127.0.0.1 ads.adiquity.com
127.0.0.1 img.ads.mobilefuse.net
127.0.0.1 as.adfonic.net
127.0.0.1 media.mobpartner.mobi
127.0.0.1 cdn.us.goldspotmedia.com
127.0.0.1 ads2.mediaarmor.com
127.0.0.1 cdn.nearbyad.com
127.0.0.1 ads.ookla.com
127.0.0.1 mobiledl.adboe.com
127.0.0.1 ads.flurry.com
127.0.0.1 gemini.yahoo.com

# [hosted on cloudfront]
127.0.0.1 d3anogn3pbtk4v.cloudfront.net
127.0.0.1 d3oltyb66oj2v8.cloudfront.net
127.0.0.1 d2bgg7rjywcwsy.cloudfront.net

# [vserv.mobi]
127.0.0.1 a.vserv.mobi
127.0.0.1 admin.vserv.mobi
127.0.0.1 c.vserv.mobi
127.0.0.1 ads.vserv.mobi
127.0.0.1 sf.vserv.mobi

# [pflexads.com]
127.0.0.1 hybl9bazbc35.pflexads.com
127.0.0.1 hhbekxxw5d9e.pflexads.com
127.0.0.1 www.pflexads.com
127.0.0.1 orencia.pflexads.com

# [velti.com]
127.0.0.1 atti.velti.com
127.0.0.1 ru.velti.com
127.0.0.1 mwc.velti.com

# [celtra.com]
127.0.0.1 cdn.celtra.com
127.0.0.1 ads.celtra.com
127.0.0.1 cache-ssl.celtra.com
127.0.0.1 cache.celtra.com
127.0.0.1 track.celtra.com

# [inner-active.mobi]
127.0.0.1 wv.inner-active.mobi
127.0.0.1 cdn1.inner-active.mobi
127.0.0.1 m2m1.inner-active.mobi

# [Jumptab]
127.0.0.1 bos-tapreq01.jumptap.com
127.0.0.1 bos-tapreq02.jumptap.com
127.0.0.1 bos-tapreq03.jumptap.com
127.0.0.1 bos-tapreq04.jumptap.com
127.0.0.1 bos-tapreq05.jumptap.com
127.0.0.1 bos-tapreq06.jumptap.com
127.0.0.1 bos-tapreq07.jumptap.com
127.0.0.1 bos-tapreq08.jumptap.com
127.0.0.1 bos-tapreq09.jumptap.com
127.0.0.1 bos-tapreq10.jumptap.com
127.0.0.1 bos-tapreq11.jumptap.com
127.0.0.1 bos-tapreq12.jumptap.com
127.0.0.1 bos-tapreq13.jumptap.com
127.0.0.1 bos-tapreq14.jumptap.com
127.0.0.1 bos-tapreq15.jumptap.com
127.0.0.1 bos-tapreq16.jumptap.com
127.0.0.1 bos-tapreq17.jumptap.com
127.0.0.1 bos-tapreq18.jumptap.com
127.0.0.1 bos-tapreq19.jumptap.com
127.0.0.1 bos-tapreq20.jumptap.com
127.0.0.1 web64.jumptap.com
127.0.0.1 web63.jumptap.com
127.0.0.1 web65.jumptap.com
127.0.0.1 bo.jumptap.com
127.0.0.1 i.jumptap.com

# [applovin]
127.0.0.1 a.applovin.com
127.0.0.1 d.applovin.com
127.0.0.1 pdn.applovin.com

# [Mobpartner]
127.0.0.1 mobpartner.mobi
127.0.0.1 go.mobpartner.mobi
127.0.0.1 r.mobpartner.mobi

# [Adinfuse]
127.0.0.1 uk-ad2.adinfuse.com
127.0.0.1 adinfuse.com
127.0.0.1 go.adinfuse.com
127.0.0.1 ad1.adinfuse.com
127.0.0.1 ad2.adinfuse.com
127.0.0.1 sky.adinfuse.com
127.0.0.1 orange-fr.adinfuse.com
127.0.0.1 sky-connect.adinfuse.com
127.0.0.1 uk-go.adinfuse.com
127.0.0.1 orangeuk-mc.adinfuse.com
127.0.0.1 intouch.adinfuse.com
127.0.0.1 funnel0.adinfuse.com
`
)

package edgeos

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	. "github.com/britannic/testutils"
)

func TestGetContent(t *testing.T) {
	tests := []struct {
		exp      string
		leaf     string
		ltype    string
		node     string
		page     string
		pageData string
		svr      *HTTPserver
		url      string
		want     string
	}{
		{
			exp:   "address=/.adsrvr.org/0.0.0.0\naddress=/.adtechus.net/0.0.0.0\naddress=/.advertising.com/0.0.0.0\naddress=/.centade.com/0.0.0.0\naddress=/.doubleclick.net/0.0.0.0\naddress=/.free-counter.co.uk/0.0.0.0\naddress=/.intellitxt.com/0.0.0.0\naddress=/.kiosked.com/0.0.0.0\n",
			ltype: preConf,
			node:  domains,
			want:  "adsrvr.org\nadtechus.net\nadvertising.com\ncentade.com\ndoubleclick.net\nfree-counter.co.uk\nintellitxt.com\nkiosked.com",
		},
		{
			exp:   "address=/beap.gemini.yahoo.com/192.168.168.1\n",
			ltype: preConf,
			node:  hosts,
			want:  "beap.gemini.yahoo.com",
		},
		{
			exp:   "address=/really.bad.phishing.site.ru/0.0.0.0\n",
			leaf:  "tasty",
			ltype: files,
			node:  hosts,
			url:   "../testdata/blist.hosts.src",
			want:  "really.bad.phishing.site.ru\n",
		},
		{
			exp:      "address=/.192-168-0-255.com/0.0.0.0\naddress=/.adsrvr.org/0.0.0.0\naddress=/.adtechus.net/0.0.0.0\naddress=/.advertising.com/0.0.0.0\naddress=/.asi-37.fr/0.0.0.0\naddress=/.bagbackpack.com/0.0.0.0\naddress=/.bitmeyenkartusistanbul.com/0.0.0.0\naddress=/.byxon.com/0.0.0.0\naddress=/.centade.com/0.0.0.0\naddress=/.doubleclick.net/0.0.0.0\naddress=/.free-counter.co.uk/0.0.0.0\naddress=/.img001.com/0.0.0.0\naddress=/.intellitxt.com/0.0.0.0\naddress=/.kiosked.com/0.0.0.0\naddress=/.loadto.net/0.0.0.0\naddress=/.roastfiles2017.com/0.0.0.0\n",
			leaf:     "malc0de",
			ltype:    urls,
			node:     domains,
			page:     "/domains.txt",
			pageData: HTTPDomainData,
			svr:      new(HTTPserver),
			want:     HTTPDomainData,
		},
		{
			exp:      "address=/a.applovin.com/192.168.168.1\naddress=/a.glcdn.co/192.168.168.1\naddress=/a.vserv.mobi/192.168.168.1\naddress=/ad.leadboltapps.net/192.168.168.1\naddress=/ad.madvertise.de/192.168.168.1\naddress=/ad.where.com/192.168.168.1\naddress=/ad1.adinfuse.com/192.168.168.1\naddress=/ad2.adinfuse.com/192.168.168.1\naddress=/adcontent.saymedia.com/192.168.168.1\naddress=/adinfuse.com/192.168.168.1\naddress=/admicro1.vcmedia.vn/192.168.168.1\naddress=/admicro2.vcmedia.vn/192.168.168.1\naddress=/admin.vserv.mobi/192.168.168.1\naddress=/ads.adiquity.com/192.168.168.1\naddress=/ads.admarvel.com/192.168.168.1\naddress=/ads.admoda.com/192.168.168.1\naddress=/ads.celtra.com/192.168.168.1\naddress=/ads.flurry.com/192.168.168.1\naddress=/ads.matomymobile.com/192.168.168.1\naddress=/ads.mobgold.com/192.168.168.1\naddress=/ads.mobilityware.com/192.168.168.1\naddress=/ads.mopub.com/192.168.168.1\naddress=/ads.n-ws.org/192.168.168.1\naddress=/ads.ookla.com/192.168.168.1\naddress=/ads.saymedia.com/192.168.168.1\naddress=/ads.smartdevicemedia.com/192.168.168.1\naddress=/ads.vserv.mobi/192.168.168.1\naddress=/ads.xxxad.net/192.168.168.1\naddress=/ads2.mediaarmor.com/192.168.168.1\naddress=/adserver.ubiyoo.com/192.168.168.1\naddress=/adultmoda.com/192.168.168.1\naddress=/android-sdk31.transpera.com/192.168.168.1\naddress=/android.bcfads.com/192.168.168.1\naddress=/api.airpush.com/192.168.168.1\naddress=/api.analytics.omgpop.com/192.168.168.1\naddress=/api.yp.com/192.168.168.1\naddress=/apps.buzzcity.net/192.168.168.1\naddress=/apps.mobilityware.com/192.168.168.1\naddress=/as.adfonic.net/192.168.168.1\naddress=/asotrack1.fluentmobile.com/192.168.168.1\naddress=/assets.cntdy.mobi/192.168.168.1\naddress=/atti.velti.com/192.168.168.1\naddress=/b.scorecardresearch.com/192.168.168.1\naddress=/banners.bigmobileads.com/192.168.168.1\naddress=/bigmobileads.com/192.168.168.1\naddress=/bo.jumptap.com/192.168.168.1\naddress=/bos-tapreq01.jumptap.com/192.168.168.1\naddress=/bos-tapreq02.jumptap.com/192.168.168.1\naddress=/bos-tapreq03.jumptap.com/192.168.168.1\naddress=/bos-tapreq04.jumptap.com/192.168.168.1\naddress=/bos-tapreq05.jumptap.com/192.168.168.1\naddress=/bos-tapreq06.jumptap.com/192.168.168.1\naddress=/bos-tapreq07.jumptap.com/192.168.168.1\naddress=/bos-tapreq08.jumptap.com/192.168.168.1\naddress=/bos-tapreq09.jumptap.com/192.168.168.1\naddress=/bos-tapreq10.jumptap.com/192.168.168.1\naddress=/bos-tapreq11.jumptap.com/192.168.168.1\naddress=/bos-tapreq12.jumptap.com/192.168.168.1\naddress=/bos-tapreq13.jumptap.com/192.168.168.1\naddress=/bos-tapreq14.jumptap.com/192.168.168.1\naddress=/bos-tapreq15.jumptap.com/192.168.168.1\naddress=/bos-tapreq16.jumptap.com/192.168.168.1\naddress=/bos-tapreq17.jumptap.com/192.168.168.1\naddress=/bos-tapreq18.jumptap.com/192.168.168.1\naddress=/bos-tapreq19.jumptap.com/192.168.168.1\naddress=/bos-tapreq20.jumptap.com/192.168.168.1\naddress=/c.vrvm.com/192.168.168.1\naddress=/c.vserv.mobi/192.168.168.1\naddress=/c753738.r38.cf2.rackcdn.com/192.168.168.1\naddress=/cache-ssl.celtra.com/192.168.168.1\naddress=/cache.celtra.com/192.168.168.1\naddress=/cdn.celtra.com/192.168.168.1\naddress=/cdn.nearbyad.com/192.168.168.1\naddress=/cdn.trafficforce.com/192.168.168.1\naddress=/cdn.us.goldspotmedia.com/192.168.168.1\naddress=/cdn.vdopia.com/192.168.168.1\naddress=/cdn1.crispadvertising.com/192.168.168.1\naddress=/cdn1.inner-active.mobi/192.168.168.1\naddress=/cdn2.crispadvertising.com/192.168.168.1\naddress=/click.buzzcity.net/192.168.168.1\naddress=/creative1cdn.mobfox.com/192.168.168.1\naddress=/d.applovin.com/192.168.168.1\naddress=/d2bgg7rjywcwsy.cloudfront.net/192.168.168.1\naddress=/d3anogn3pbtk4v.cloudfront.net/192.168.168.1\naddress=/d3oltyb66oj2v8.cloudfront.net/192.168.168.1\naddress=/edge.reporo.net/192.168.168.1\naddress=/ftpcontent.worldnow.com/192.168.168.1\naddress=/funnel0.adinfuse.com/192.168.168.1\naddress=/gemini.yahoo.com/192.168.168.1\naddress=/go.adinfuse.com/192.168.168.1\naddress=/go.mobpartner.mobi/192.168.168.1\naddress=/go.vrvm.com/192.168.168.1\naddress=/gsmtop.net/192.168.168.1\naddress=/gts-ads.twistbox.com/192.168.168.1\naddress=/hhbekxxw5d9e.pflexads.com/192.168.168.1\naddress=/hybl9bazbc35.pflexads.com/192.168.168.1\naddress=/i.jumptap.com/192.168.168.1\naddress=/i.tapit.com/192.168.168.1\naddress=/images.millennialmedia.com/192.168.168.1\naddress=/images.mpression.net/192.168.168.1\naddress=/img.ads.huntmad.com/192.168.168.1\naddress=/img.ads.mobilefuse.net/192.168.168.1\naddress=/img.ads.mocean.mobi/192.168.168.1\naddress=/img.ads.mojiva.com/192.168.168.1\naddress=/img.ads.taptapnetworks.com/192.168.168.1\naddress=/intouch.adinfuse.com/192.168.168.1\naddress=/lb.usemaxserver.de/192.168.168.1\naddress=/m.adsymptotic.com/192.168.168.1\naddress=/m2m1.inner-active.mobi/192.168.168.1\naddress=/media.mobpartner.mobi/192.168.168.1\naddress=/medrx.sensis.com.au/192.168.168.1\naddress=/mobile.banzai.it/192.168.168.1\naddress=/mobiledl.adboe.com/192.168.168.1\naddress=/mobpartner.mobi/192.168.168.1\naddress=/mwc.velti.com/192.168.168.1\naddress=/netdna.reporo.net/192.168.168.1\naddress=/oasc04012.247realmedia.com/192.168.168.1\naddress=/orange-fr.adinfuse.com/192.168.168.1\naddress=/orangeuk-mc.adinfuse.com/192.168.168.1\naddress=/orencia.pflexads.com/192.168.168.1\naddress=/pdn.applovin.com/192.168.168.1\naddress=/r.edge.inmobicdn.net/192.168.168.1\naddress=/r.mobpartner.mobi/192.168.168.1\naddress=/req.appads.com/192.168.168.1\naddress=/rs-staticart.ybcdn.net/192.168.168.1\naddress=/ru.velti.com/192.168.168.1\naddress=/s0.2mdn.net/192.168.168.1\naddress=/s3.phluant.com/192.168.168.1\naddress=/sf.vserv.mobi/192.168.168.1\naddress=/show.buzzcity.net/192.168.168.1\naddress=/sky-connect.adinfuse.com/192.168.168.1\naddress=/sky.adinfuse.com/192.168.168.1\naddress=/static.cdn.gtsmobi.com/192.168.168.1\naddress=/static.estebull.com/192.168.168.1\naddress=/stats.pflexads.com/192.168.168.1\naddress=/track.celtra.com/192.168.168.1\naddress=/tracking.klickthru.com/192.168.168.1\naddress=/uk-ad2.adinfuse.com/192.168.168.1\naddress=/uk-go.adinfuse.com/192.168.168.1\naddress=/web63.jumptap.com/192.168.168.1\naddress=/web64.jumptap.com/192.168.168.1\naddress=/web65.jumptap.com/192.168.168.1\naddress=/wv.inner-active.mobi/192.168.168.1\naddress=/www.eltrafiko.com/192.168.168.1\naddress=/www.mmnetwork.mobi/192.168.168.1\naddress=/www.pflexads.com/192.168.168.1\naddress=/wwww.adleads.com/192.168.168.1\n",
			leaf:     "adaway",
			ltype:    urls,
			node:     hosts,
			page:     "/hosts.txt",
			pageData: httpHostData,
			svr:      new(HTTPserver),
			want:     httpHostData,
		},
	}
	r := &CFGstatic{Cfg: Cfg}
	c := NewConfig(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		FileNameFmt("%v/%v.%v.%v"),
		Method("GET"),
		Nodes([]string{domains, hosts}),
		Prefix("address="),
		STypes([]string{preConf, files, urls}),
	)

	err := c.ReadCfg(r)
	OK(t, err)

	for _, test := range tests {
		objs := c.bNodes[test.node].Objects
		switch test.ltype {
		case urls:
			test.url = test.svr.NewHTTPServer().String() + test.page
			go test.svr.Mux.HandleFunc(test.page,
				func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, test.pageData)
				},
			)

			if i := objs.Find(test.leaf); i > -1 {
				objs.S[i].url = test.url
			}
		case files:
			if i := objs.Find(test.leaf); i > -1 {
				objs.S[i].file = test.url
			}
		}

		for _, src := range *c.Get(test.node).Source(test.ltype).GetContent() {
			got, err := ioutil.ReadAll(src.r)
			OK(t, err)
			Equals(t, test.want, string(got))
			src.r = io.MultiReader(bytes.NewReader(got))
			got, err = ioutil.ReadAll(src.Process().r)
			OK(t, err)
			Equals(t, test.exp, string(got))
		}
	}
}

func TestGetAllContent(t *testing.T) {
	var (
		wantPre = "File: \"/tmp/domains.pre-configured.blacklist.conf\"\nData:\naddress=/.adsrvr.org/0.0.0.0\naddress=/.adtechus.net/0.0.0.0\naddress=/.advertising.com/0.0.0.0\naddress=/.centade.com/0.0.0.0\naddress=/.doubleclick.net/0.0.0.0\naddress=/.free-counter.co.uk/0.0.0.0\naddress=/.intellitxt.com/0.0.0.0\naddress=/.kiosked.com/0.0.0.0\n\nFile: \"/tmp/hosts.pre-configured.blacklist.conf\"\nData:\naddress=/beap.gemini.yahoo.com/192.168.168.1\n\n"

		wantAll = "File: \"/tmp/domains.pre-configured.blacklist.conf\"\nData:\naddress=/.adsrvr.org/0.0.0.0\naddress=/.adtechus.net/0.0.0.0\naddress=/.advertising.com/0.0.0.0\naddress=/.centade.com/0.0.0.0\naddress=/.doubleclick.net/0.0.0.0\naddress=/.free-counter.co.uk/0.0.0.0\naddress=/.intellitxt.com/0.0.0.0\naddress=/.kiosked.com/0.0.0.0\n\nFile: \"/tmp/domains.malc0de.blacklist.conf\"\nData:\naddress=/.adsrvr.org/0.0.0.0\naddress=/.adtechus.net/0.0.0.0\naddress=/.advertising.com/0.0.0.0\naddress=/.centade.com/0.0.0.0\naddress=/.doubleclick.net/0.0.0.0\naddress=/.free-counter.co.uk/0.0.0.0\naddress=/.intellitxt.com/0.0.0.0\naddress=/.kiosked.com/0.0.0.0\n\nFile: \"/tmp/hosts.pre-configured.blacklist.conf\"\nData:\n\nFile: \"/tmp/hosts.adaway.blacklist.conf\"\nData:\naddress=/a.applovin.com/192.168.168.1\naddress=/a.glcdn.co/192.168.168.1\naddress=/a.vserv.mobi/192.168.168.1\naddress=/ad.leadboltapps.net/192.168.168.1\naddress=/ad.madvertise.de/192.168.168.1\naddress=/ad.where.com/192.168.168.1\naddress=/ad1.adinfuse.com/192.168.168.1\naddress=/ad2.adinfuse.com/192.168.168.1\naddress=/adcontent.saymedia.com/192.168.168.1\naddress=/adinfuse.com/192.168.168.1\naddress=/admicro1.vcmedia.vn/192.168.168.1\naddress=/admicro2.vcmedia.vn/192.168.168.1\naddress=/admin.vserv.mobi/192.168.168.1\naddress=/ads.adiquity.com/192.168.168.1\naddress=/ads.admarvel.com/192.168.168.1\naddress=/ads.admoda.com/192.168.168.1\naddress=/ads.celtra.com/192.168.168.1\naddress=/ads.flurry.com/192.168.168.1\naddress=/ads.matomymobile.com/192.168.168.1\naddress=/ads.mobgold.com/192.168.168.1\naddress=/ads.mobilityware.com/192.168.168.1\naddress=/ads.mopub.com/192.168.168.1\naddress=/ads.n-ws.org/192.168.168.1\naddress=/ads.ookla.com/192.168.168.1\naddress=/ads.saymedia.com/192.168.168.1\naddress=/ads.smartdevicemedia.com/192.168.168.1\naddress=/ads.vserv.mobi/192.168.168.1\naddress=/ads.xxxad.net/192.168.168.1\naddress=/ads2.mediaarmor.com/192.168.168.1\naddress=/adserver.ubiyoo.com/192.168.168.1\naddress=/adultmoda.com/192.168.168.1\naddress=/android-sdk31.transpera.com/192.168.168.1\naddress=/android.bcfads.com/192.168.168.1\naddress=/api.airpush.com/192.168.168.1\naddress=/api.analytics.omgpop.com/192.168.168.1\naddress=/api.yp.com/192.168.168.1\naddress=/apps.buzzcity.net/192.168.168.1\naddress=/apps.mobilityware.com/192.168.168.1\naddress=/as.adfonic.net/192.168.168.1\naddress=/asotrack1.fluentmobile.com/192.168.168.1\naddress=/assets.cntdy.mobi/192.168.168.1\naddress=/atti.velti.com/192.168.168.1\naddress=/b.scorecardresearch.com/192.168.168.1\naddress=/banners.bigmobileads.com/192.168.168.1\naddress=/bigmobileads.com/192.168.168.1\naddress=/bo.jumptap.com/192.168.168.1\naddress=/bos-tapreq01.jumptap.com/192.168.168.1\naddress=/bos-tapreq02.jumptap.com/192.168.168.1\naddress=/bos-tapreq03.jumptap.com/192.168.168.1\naddress=/bos-tapreq04.jumptap.com/192.168.168.1\naddress=/bos-tapreq05.jumptap.com/192.168.168.1\naddress=/bos-tapreq06.jumptap.com/192.168.168.1\naddress=/bos-tapreq07.jumptap.com/192.168.168.1\naddress=/bos-tapreq08.jumptap.com/192.168.168.1\naddress=/bos-tapreq09.jumptap.com/192.168.168.1\naddress=/bos-tapreq10.jumptap.com/192.168.168.1\naddress=/bos-tapreq11.jumptap.com/192.168.168.1\naddress=/bos-tapreq12.jumptap.com/192.168.168.1\naddress=/bos-tapreq13.jumptap.com/192.168.168.1\naddress=/bos-tapreq14.jumptap.com/192.168.168.1\naddress=/bos-tapreq15.jumptap.com/192.168.168.1\naddress=/bos-tapreq16.jumptap.com/192.168.168.1\naddress=/bos-tapreq17.jumptap.com/192.168.168.1\naddress=/bos-tapreq18.jumptap.com/192.168.168.1\naddress=/bos-tapreq19.jumptap.com/192.168.168.1\naddress=/bos-tapreq20.jumptap.com/192.168.168.1\naddress=/c.vrvm.com/192.168.168.1\naddress=/c.vserv.mobi/192.168.168.1\naddress=/c753738.r38.cf2.rackcdn.com/192.168.168.1\naddress=/cache-ssl.celtra.com/192.168.168.1\naddress=/cache.celtra.com/192.168.168.1\naddress=/cdn.celtra.com/192.168.168.1\naddress=/cdn.nearbyad.com/192.168.168.1\naddress=/cdn.trafficforce.com/192.168.168.1\naddress=/cdn.us.goldspotmedia.com/192.168.168.1\naddress=/cdn.vdopia.com/192.168.168.1\naddress=/cdn1.crispadvertising.com/192.168.168.1\naddress=/cdn1.inner-active.mobi/192.168.168.1\naddress=/cdn2.crispadvertising.com/192.168.168.1\naddress=/click.buzzcity.net/192.168.168.1\naddress=/creative1cdn.mobfox.com/192.168.168.1\naddress=/d.applovin.com/192.168.168.1\naddress=/d2bgg7rjywcwsy.cloudfront.net/192.168.168.1\naddress=/d3anogn3pbtk4v.cloudfront.net/192.168.168.1\naddress=/d3oltyb66oj2v8.cloudfront.net/192.168.168.1\naddress=/edge.reporo.net/192.168.168.1\naddress=/ftpcontent.worldnow.com/192.168.168.1\naddress=/funnel0.adinfuse.com/192.168.168.1\naddress=/gemini.yahoo.com/192.168.168.1\naddress=/go.adinfuse.com/192.168.168.1\naddress=/go.mobpartner.mobi/192.168.168.1\naddress=/go.vrvm.com/192.168.168.1\naddress=/gsmtop.net/192.168.168.1\naddress=/gts-ads.twistbox.com/192.168.168.1\naddress=/hhbekxxw5d9e.pflexads.com/192.168.168.1\naddress=/hybl9bazbc35.pflexads.com/192.168.168.1\naddress=/i.jumptap.com/192.168.168.1\naddress=/i.tapit.com/192.168.168.1\naddress=/images.millennialmedia.com/192.168.168.1\naddress=/images.mpression.net/192.168.168.1\naddress=/img.ads.huntmad.com/192.168.168.1\naddress=/img.ads.mobilefuse.net/192.168.168.1\naddress=/img.ads.mocean.mobi/192.168.168.1\naddress=/img.ads.mojiva.com/192.168.168.1\naddress=/img.ads.taptapnetworks.com/192.168.168.1\naddress=/intouch.adinfuse.com/192.168.168.1\naddress=/lb.usemaxserver.de/192.168.168.1\naddress=/m.adsymptotic.com/192.168.168.1\naddress=/m2m1.inner-active.mobi/192.168.168.1\naddress=/media.mobpartner.mobi/192.168.168.1\naddress=/medrx.sensis.com.au/192.168.168.1\naddress=/mobile.banzai.it/192.168.168.1\naddress=/mobiledl.adboe.com/192.168.168.1\naddress=/mobpartner.mobi/192.168.168.1\naddress=/mwc.velti.com/192.168.168.1\naddress=/netdna.reporo.net/192.168.168.1\naddress=/oasc04012.247realmedia.com/192.168.168.1\naddress=/orange-fr.adinfuse.com/192.168.168.1\naddress=/orangeuk-mc.adinfuse.com/192.168.168.1\naddress=/orencia.pflexads.com/192.168.168.1\naddress=/pdn.applovin.com/192.168.168.1\naddress=/r.edge.inmobicdn.net/192.168.168.1\naddress=/r.mobpartner.mobi/192.168.168.1\naddress=/req.appads.com/192.168.168.1\naddress=/rs-staticart.ybcdn.net/192.168.168.1\naddress=/ru.velti.com/192.168.168.1\naddress=/s0.2mdn.net/192.168.168.1\naddress=/s3.phluant.com/192.168.168.1\naddress=/sf.vserv.mobi/192.168.168.1\naddress=/show.buzzcity.net/192.168.168.1\naddress=/sky-connect.adinfuse.com/192.168.168.1\naddress=/sky.adinfuse.com/192.168.168.1\naddress=/static.cdn.gtsmobi.com/192.168.168.1\naddress=/static.estebull.com/192.168.168.1\naddress=/stats.pflexads.com/192.168.168.1\naddress=/track.celtra.com/192.168.168.1\naddress=/tracking.klickthru.com/192.168.168.1\naddress=/uk-ad2.adinfuse.com/192.168.168.1\naddress=/uk-go.adinfuse.com/192.168.168.1\naddress=/web63.jumptap.com/192.168.168.1\naddress=/web64.jumptap.com/192.168.168.1\naddress=/web65.jumptap.com/192.168.168.1\naddress=/wv.inner-active.mobi/192.168.168.1\naddress=/www.eltrafiko.com/192.168.168.1\naddress=/www.mmnetwork.mobi/192.168.168.1\naddress=/www.pflexads.com/192.168.168.1\naddress=/wwww.adleads.com/192.168.168.1\n\nFile: \"/tmp/hosts.tasty.blacklist.conf\"\nData:\naddress=/really.bad.phishing.site.ru/0.0.0.0\n\n"

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
		dPage = "/hosts/host.txt"
		hPage = "/domains/domain.txt"
	)

	svr := new(HTTPserver)
	domainsHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, HTTPDomainData)
	}

	hostsHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, httpHostData)
	}
	err := c.ReadCfg(r)
	OK(t, err)

	for _, node := range c.Nodes() {
		objs := c.bNodes[node].Objects
		if i := objs.Find("malc0de"); i > -1 {
			objs.S[i].url = svr.NewHTTPServer().String() + dPage
		}
		// c.bNodes[domains].data["malc0de"].url = svr.NewHTTPServer().String() + dPage
		if i := objs.Find("adaway"); i > -1 {
			objs.S[i].url = svr.NewHTTPServer().String() + hPage
		}
		// c.bNodes[hosts].data["adaway"].url = svr.NewHTTPServer().String() + hPage
	}

	svr.Mux.HandleFunc(dPage, domainsHandler)
	svr.Mux.HandleFunc(hPage, hostsHandler)

	act := fmt.Sprint(c.GetAll(preConf).GetContent())
	// fmt.Println(z)
	Equals(t, wantPre, act)

	act = fmt.Sprint(c.GetAll().GetContent())
	Equals(t, wantAll, act)
	// fmt.Println(z)
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
			data:  bytes.NewBufferString("The rest is history!"),
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
            url http://adaway.org/hosts.txt
        }
				source tasty {
						description "File source"
						dns-redirect-ip 0.0.0.0
						file /config/user-data/blist.hosts.src
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
)

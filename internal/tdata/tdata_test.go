package tdata_test

import (
	"fmt"
	"testing"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTdataCfg(t *testing.T) {
	Convey("Testing TdataCfg()", t, func() {
		exp := map[string]string{
			"cfg":  cfg,
			"cfg2": cfg2,
			"cfg3": cfg3,
		}

		for k := range exp {
			act, err := tdata.Get(k)
			So(err, ShouldBeNil)
			So(act, ShouldEqual, exp[k])
		}
	})
}

func TestTdataCfg2(t *testing.T) {
	Convey("Testing TdataCfg2()", t, func() {
		act, err := tdata.Get("cfg2")
		So(err, ShouldBeNil)
		So(act, ShouldEqual, cfg2)
	})
}

func TestTdataCfg3(t *testing.T) {
	Convey("Testing TdataCfg3()", t, func() {
		act, err := tdata.Get("cfg3")
		So(err, ShouldBeNil)
		So(act, ShouldEqual, cfg3)
	})
}

func TestTdataFileManifest(t *testing.T) {
	Convey("Testing TdataFileManifest()", t, func() {
		act, err := tdata.Get("fileManifest")
		So(err, ShouldBeNil)
		So(act, ShouldEqual, fileManifest)
	})
}

func TestTdataDefault(t *testing.T) {
	Convey("Testing TdataDefault()", t, func() {
		_, err := tdata.Get("default")
		So(err, ShouldNotBeNil)
		So(err, ShouldResemble, fmt.Errorf("Get(default) is unknown!"))
	})
}

var (
	// Cfg contains a valid full EdgeOS blacklist configuration
	cfg = "blacklist {\n    disabled false\n    dns-redirect-ip 0.0.0.0\n    domains {\n        dns-redirect-ip 192.168.100.1\n        include adsrvr.org\n        include adtechus.net\n        include advertising.com\n        include centade.com\n        include doubleclick.net\n        include free-counter.co.uk\n        include intellitxt.com\n        include kiosked.com\n        include patoghee.in\n        source malc0de {\n            dns-redirect-ip 192.168.168.1\n            description \"List of zones serving malicious executables observed by malc0de.com/database/\"\n            prefix \"zone \"\n            url http://malc0de.com/bl/ZONES\n        }\n        source malwaredomains.com {\n            dns-redirect-ip 10.0.0.1\n            description \"Just domains\"\n            prefix \"\"\n            url http://mirror1.malwaredomains.com/files/justdomains\n        }\n        source simple_tracking {\n            description \"Basic tracking list by Disconnect\"\n            prefix \"\"\n            url https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt\n        }\n        source zeus {\n            description \"abuse.ch ZeuS domain blocklist\"\n            prefix \"\"\n            url https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist\n        }\n    }\n    exclude 1e100.net\n    exclude 2o7.net\n    exclude adobedtm.com\n    exclude akamai.net\n    exclude akamaihd.net\n    exclude amazon.com\n    exclude amazonaws.com\n    exclude apple.com\n    exclude ask.com\n    exclude avast.com\n    exclude bitdefender.com\n    exclude cdn.visiblemeasures.com\n    exclude cloudfront.net\n    exclude coremetrics.com\n    exclude edgesuite.net\n    exclude freedns.afraid.org\n    exclude github.com\n    exclude githubusercontent.com\n    exclude google.com\n    exclude googleadservices.com\n    exclude googleapis.com\n    exclude googletagmanager.com\n    exclude googleusercontent.com\n    exclude gstatic.com\n    exclude gvt1.com\n    exclude gvt1.net\n    exclude hb.disney.go.com\n    exclude hp.com\n    exclude hulu.com\n    exclude images-amazon.com\n    exclude live.com\n    exclude microsoft.com\n    exclude msdn.com\n    exclude msecnd.net\n    exclude paypal.com\n    exclude rackcdn.com\n    exclude schema.org\n    exclude shopify.com\n    exclude skype.com\n    exclude smacargo.com\n    exclude sourceforge.net\n    exclude ssl-on9.com\n    exclude ssl-on9.net\n    exclude sstatic.net\n    exclude static.chartbeat.com\n    exclude storage.googleapis.com\n    exclude windows.net\n    exclude xboxlive.com\n    exclude yimg.com\n    exclude ytimg.com\n    hosts {\n        include beap.gemini.yahoo.com\n        source openphish {\n            description \"OpenPhish automatic phishing detection\"\n            prefix http\n            url https://openphish.com/feed.txt\n        }\n        source raw.github.com {\n            description \"This hosts file is a merged collection of hosts from reputable sources\"\n            prefix \"0.0.0.0 \"\n            url https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts\n        }\n        source sysctl.org {\n            dns-redirect-ip 172.16.16.1\n            description \"This hosts file is a merged collection of hosts from cameleon\"\n            prefix \"127.0.0.1\t \"\n            url http://sysctl.org/cameleon/hosts\n        }\n        source tasty {\n            description \"File source\"\n            dns-redirect-ip 10.10.10.10\n            file ../internal/testdata/blist.hosts.src\n        }\n        source volkerschatz {\n            description \"Ad server blacklists\"\n            prefix http\n            url http://www.volkerschatz.com/net/adpaths\n        }\n        source yoyo {\n            description \"Fully Qualified Domain Names only - no prefix to strip\"\n            prefix \"\"\n            url http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\n        }\n    }\n}\n\n\t/* Warning: Do not remove the following line. */\n\t/* === vyatta-config-version: \"config-management@1:conntrack@1:cron@1:dhcp-relay@1:dhcp-server@4:firewall@5:ipsec@5:nat@3:qos@1:quagga@2:system@4:ubnt-pptp@1:ubnt-util@1:vrrp@1:webgui@1:webproxy@1:zone-policy@1\" === */\n\t/* Release version: v1.8.5.4884695.160608.1057 */\n}"

	// Cfg2 contains a valid partial EdgeOS blacklist configuration
	cfg2 = `blacklist {
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
        include beap.gemini.yahoo.com
    }
}`

	// Cfg3 contains a valid partial EdgeOS blacklist configuration
	cfg3 = "blacklist {\n    disabled false\n    dns-redirect-ip 0.0.0.0\n    domains {\n        include adsrvr.org\n        include adtechus.net\n        include advertising.com\n        include centade.com\n        include doubleclick.net\n        include free-counter.co.uk\n        include intellitxt.com\n        include kiosked.com\n        source malc0de {\n            description \"List of zones serving malicious executables observed by malc0de.com/database/\"\n            prefix \"zone \"\n            url http://malc0de.com/bl/ZONES\n        }\n    }\n    exclude ytimg.com\n    hosts {\n        include beap.gemini.yahoo.com\n        source tasty {\n            description \"File source\"\n            dns-redirect-ip 10.10.10.10\n            file ../internal/testdata/blist.hosts.src\n        }\n    }\n}"

	// FileManifest is complete list of the blacklist config node templates
	fileManifest = `../payload/blacklist
../payload/blacklist/disabled
../payload/blacklist/disabled/node.def
../payload/blacklist/dns-redirect-ip
../payload/blacklist/dns-redirect-ip/node.def
../payload/blacklist/domains
../payload/blacklist/domains/dns-redirect-ip
../payload/blacklist/domains/dns-redirect-ip/node.def
../payload/blacklist/domains/exclude
../payload/blacklist/domains/exclude/node.def
../payload/blacklist/domains/include
../payload/blacklist/domains/include/node.def
../payload/blacklist/domains/node.def
../payload/blacklist/domains/source
../payload/blacklist/domains/source/node.def
../payload/blacklist/domains/source/node.tag
../payload/blacklist/domains/source/node.tag/description
../payload/blacklist/domains/source/node.tag/description/node.def
../payload/blacklist/domains/source/node.tag/prefix
../payload/blacklist/domains/source/node.tag/prefix/node.def
../payload/blacklist/domains/source/node.tag/url
../payload/blacklist/domains/source/node.tag/url/node.def
../payload/blacklist/exclude
../payload/blacklist/exclude/node.def
../payload/blacklist/hosts
../payload/blacklist/hosts/dns-redirect-ip
../payload/blacklist/hosts/dns-redirect-ip/node.def
../payload/blacklist/hosts/exclude
../payload/blacklist/hosts/exclude/node.def
../payload/blacklist/hosts/include
../payload/blacklist/hosts/include/node.def
../payload/blacklist/hosts/node.def
../payload/blacklist/hosts/source
../payload/blacklist/hosts/source/node.def
../payload/blacklist/hosts/source/node.tag
../payload/blacklist/hosts/source/node.tag/description
../payload/blacklist/hosts/source/node.tag/description/node.def
../payload/blacklist/hosts/source/node.tag/prefix
../payload/blacklist/hosts/source/node.tag/prefix/node.def
../payload/blacklist/hosts/source/node.tag/url
../payload/blacklist/hosts/source/node.tag/url/node.def
../payload/blacklist/node.def
`
)

package tdata_test

import (
	"fmt"
	"testing"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/britannic/testutils"
)

func TestTdataCfg(t *testing.T) {
	want := map[string]string{
		"cfg":  cfg,
		"cfg2": cfg2,
		"cfg3": cfg3,
	}

	for k := range want {
		got, err := tdata.Get(k)
		OK(t, err)
		Equals(t, want[k], got)
	}
}

func TestTdataCfg2(t *testing.T) {
	r, err := tdata.Get("cfg2")
	OK(t, err)
	Equals(t, cfg2, r)
}

func TestTdataCfg3(t *testing.T) {
	r, err := tdata.Get("cfg3")
	OK(t, err)
	Equals(t, cfg3, r)
}

func TestTdataFileManifest(t *testing.T) {
	r, err := tdata.Get("fileManifest")
	OK(t, err)
	Equals(t, fileManifest, r)
}

func TestTdataDefault(t *testing.T) {
	_, err := tdata.Get("default")
	NotOK(t, err)
	Equals(t, fmt.Errorf("Get(default) is unknown!"), err)
}

var (
	// Cfg contains a valid full EdgeOS blacklist configuration
	cfg = `blacklist {
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
        source malwaredomainlist {
            description "127.0.0.1 based host and domain list"
            prefix "127.0.0.1 "
            url http://www.malwaredomainlist.com/hostslist/hosts.txt
        }
        source openphish {
            description "OpenPhish automatic phishing detection"
            prefix http
            url https://openphish.com/feed.txt
        }
        source someonewhocares {
            description "Zero based host and domain list"
            prefix 0.0.0.0
            url http://someonewhocares.org/hosts/zero/
        }
				source tasty {
						description "File source"
						dns-redirect-ip 0.0.0.0
						file /config/user-data/blist.hosts.src
				}
        source volkerschatz {
            description "Ad server blacklists"
            prefix http
            url http://www.volkerschatz.com/net/adpaths
        }
        source winhelp2002 {
            description "Zero based host and domain list"
						dns-redirect-ip 0.0.0.0
            prefix "0.0.0.0 "
            url http://winhelp2002.mvps.org/hosts.txt
        }
        source yoyo {
            description "Fully Qualified Domain Names only - no prefix to strip"
            prefix ""
            url http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext
        }
    }
	/* Warning: Do not remove the following line. */
	/* === vyatta-config-version: "config-management@1:conntrack@1:cron@1:dhcp-relay@1:dhcp-server@4:firewall@5:ipsec@5:nat@3:qos@1:quagga@2:system@4:ubnt-pptp@1:ubnt-util@1:vrrp@1:webgui@1:webproxy@1:zone-policy@1" === */
	/* Release version: v1.8.0.4853089.160219.1607 */
}`

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
	cfg3 = `blacklist {
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
        include beap.gemini.yahoo.com
    }
}`

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

package tdata

import (
	"fmt"
)

// Get returns r
func Get(s string) (r string, err error) {
	switch s {
	case "cfg":
		r = Cfg
	case "cfg2":
		r = CfgPartial
	case "cfg3":
		r = CfgMimimal
	case "fileManifest":
		r = FileManifest
	default:
		err = fmt.Errorf("Get(%v) is unknown!", s)
	}
	return r, err
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

	// CfgPartial contains a valid partial EdgeOS blacklist configuration
	CfgPartial = `blacklist {
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
        include beap.gemini.yahoo.com
    }
}`

	// DisabledCfg contains a disabled valid  EdgeOS blacklist configuration
	DisabledCfg = `blacklist {
    disabled true
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
        dns-redirect-ip 0.0.0.0
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
            prefix "0.0.0.0 "
            url http://winhelp2002.mvps.org/hosts.txt
        }
        source yoyo {
            description "Fully Qualified Domain Names only - no prefix to strip"
						file
						ip
            prefix ""
            url http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext
        }
    }
}`

	// ZeroHostSourcesCfg is a valid EdgeOS blacklist configuration with zero hosts sources
	ZeroHostSourcesCfg = `blacklist {
	disabled false
	dns-redirect-ip 0.0.0.0
	domains {
			dns-redirect-ip
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
			dns-redirect-ip
			include beap.gemini.yahoo.com
	}
}`

	// JSONcfg is JSON formatted blacklist configuration output
	JSONcfg = "{\n  \"nodes\": [{\n    \"blacklist\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [\n        \"122.2o7.net\",\n        \"1e100.net\",\n        \"adobedtm.com\",\n        \"akamai.net\",\n        \"amazon.com\",\n        \"amazonaws.com\",\n        \"apple.com\",\n        \"ask.com\",\n        \"avast.com\",\n        \"bitdefender.com\",\n        \"cdn.visiblemeasures.com\",\n        \"cloudfront.net\",\n        \"coremetrics.com\",\n        \"edgesuite.net\",\n        \"freedns.afraid.org\",\n        \"github.com\",\n        \"githubusercontent.com\",\n        \"google.com\",\n        \"googleadservices.com\",\n        \"googleapis.com\",\n        \"googleusercontent.com\",\n        \"gstatic.com\",\n        \"gvt1.com\",\n        \"gvt1.net\",\n        \"hb.disney.go.com\",\n        \"hp.com\",\n        \"hulu.com\",\n        \"images-amazon.com\",\n        \"msdn.com\",\n        \"paypal.com\",\n        \"rackcdn.com\",\n        \"schema.org\",\n        \"skype.com\",\n        \"smacargo.com\",\n        \"sourceforge.net\",\n        \"ssl-on9.com\",\n        \"ssl-on9.net\",\n        \"static.chartbeat.com\",\n        \"storage.googleapis.com\",\n        \"windows.net\",\n        \"yimg.com\",\n        \"ytimg.com\"\n        ]\n    },\n    \"domains\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [],\n      \"includes\": [\n        \"adsrvr.org\",\n        \"adtechus.net\",\n        \"advertising.com\",\n        \"centade.com\",\n        \"doubleclick.net\",\n        \"free-counter.co.uk\",\n        \"intellitxt.com\",\n        \"kiosked.com\"\n        ],\n      \"sources\": [{\n        \"malc0de\": {\n          \"disabled\": \"false\",\n          \"description\": \"List of zones serving malicious executables observed by malc0de.com/database/\",\n          \"prefix\": \"zone \",\n          \"file\": \"\",\n          \"url\": \"http://malc0de.com/bl/ZONES\"\n        }\n    }]\n    },\n    \"hosts\": {\n      \"disabled\": \"false\",\n      \"ip\": \"192.168.168.1\",\n      \"excludes\": [],\n      \"includes\": [\"beap.gemini.yahoo.com\"],\n      \"sources\": [{\n        \"adaway\": {\n          \"disabled\": \"false\",\n          \"description\": \"Blocking mobile ad providers and some analytics providers\",\n          \"prefix\": \"127.0.0.1 \",\n          \"file\": \"\",\n          \"url\": \"http://adaway.org/hosts.txt\"\n        },\n        \"tasty\": {\n          \"disabled\": \"false\",\n          \"description\": \"File source\",\n          \"prefix\": \"\",\n          \"file\": \"/config/user-data/blist.hosts.src\",\n          \"url\": \"\"\n        }\n    }]\n    }\n  }]\n}"

	// JSONcfgZeroHostSources is JSON formatted blacklist configuration output with zero sources for hosts
	JSONcfgZeroHostSources = "{\n  \"nodes\": [{\n    \"blacklist\": {\n      \"disabled\": \"false\",\n      \"ip\": \"0.0.0.0\",\n      \"excludes\": [\n        \"122.2o7.net\",\n        \"1e100.net\",\n        \"adobedtm.com\",\n        \"akamai.net\",\n        \"amazon.com\",\n        \"amazonaws.com\",\n        \"apple.com\",\n        \"ask.com\",\n        \"avast.com\",\n        \"bitdefender.com\",\n        \"cdn.visiblemeasures.com\",\n        \"cloudfront.net\",\n        \"coremetrics.com\",\n        \"edgesuite.net\",\n        \"freedns.afraid.org\",\n        \"github.com\",\n        \"githubusercontent.com\",\n        \"google.com\",\n        \"googleadservices.com\",\n        \"googleapis.com\",\n        \"googleusercontent.com\",\n        \"gstatic.com\",\n        \"gvt1.com\",\n        \"gvt1.net\",\n        \"hb.disney.go.com\",\n        \"hp.com\",\n        \"hulu.com\",\n        \"images-amazon.com\",\n        \"msdn.com\",\n        \"paypal.com\",\n        \"rackcdn.com\",\n        \"schema.org\",\n        \"skype.com\",\n        \"smacargo.com\",\n        \"sourceforge.net\",\n        \"ssl-on9.com\",\n        \"ssl-on9.net\",\n        \"static.chartbeat.com\",\n        \"storage.googleapis.com\",\n        \"windows.net\",\n        \"yimg.com\",\n        \"ytimg.com\"\n        ]\n    },\n    \"domains\": {\n      \"disabled\": \"false\",\n            \"excludes\": [],\n      \"includes\": [\n        \"adsrvr.org\",\n        \"adtechus.net\",\n        \"advertising.com\",\n        \"centade.com\",\n        \"doubleclick.net\",\n        \"free-counter.co.uk\",\n        \"intellitxt.com\",\n        \"kiosked.com\"\n        ],\n      \"sources\": [{\n        \"malc0de\": {\n          \"disabled\": \"false\",\n          \"description\": \"List of zones serving malicious executables observed by malc0de.com/database/\",\n          \"prefix\": \"zone \",\n          \"file\": \"\",\n          \"url\": \"http://malc0de.com/bl/ZONES\"\n        }\n    }]\n    },\n    \"hosts\": {\n      \"disabled\": \"false\",\n            \"excludes\": [],\n      \"includes\": [\"beap.gemini.yahoo.com\"],\n      \"sources\": [{}]\n    }\n  }]\n}"

	// JSONrawcfg is JSON unformatted blacklist configuration output
	JSONrawcfg = `{"blacklist":{"data":{},"disable":false,"excludes":["122.2o7.net","1e100.net","adobedtm.com","akamai.net","amazon.com","amazonaws.com","apple.com","ask.com","avast.com","bitdefender.com","cdn.visiblemeasures.com","cloudfront.net","coremetrics.com","edgesuite.net","freedns.afraid.org","github.com","githubusercontent.com","google.com","googleadservices.com","googleapis.com","googleusercontent.com","gstatic.com","gvt1.com","gvt1.net","hb.disney.go.com","hp.com","hulu.com","images-amazon.com","msdn.com","paypal.com","rackcdn.com","schema.org","skype.com","smacargo.com","sourceforge.net","ssl-on9.com","ssl-on9.net","static.chartbeat.com","storage.googleapis.com","windows.net","yimg.com","ytimg.com"],"includes":[],"ip":"0.0.0.0"},"domains":{"data":{"malc0de":{"desc":"List of zones serving malicious executables observed by malc0de.com/database/","disabled":false,"file":"","ip":"","name":"malc0de","prefix":"zone ","type":2,"url":"http://malc0de.com/bl/ZONES"}},"disable":false,"excludes":[],"includes":["adsrvr.org","adtechus.net","advertising.com","centade.com","doubleclick.net","free-counter.co.uk","intellitxt.com","kiosked.com"],"ip":"0.0.0.0"},"hosts":{"data":{"adaway":{"desc":"Blocking mobile ad providers and some analytics providers","disabled":false,"file":"","ip":"","name":"adaway","prefix":"127.0.0.1 ","type":3,"url":"http://adaway.org/hosts.txt"},"malwaredomainlist":{"desc":"127.0.0.1 based host and domain list","disabled":false,"file":"","ip":"","name":"malwaredomainlist","prefix":"127.0.0.1 ","type":3,"url":"http://www.malwaredomainlist.com/hostslist/hosts.txt"},"openphish":{"desc":"OpenPhish automatic phishing detection","disabled":false,"file":"","ip":"","name":"openphish","prefix":"http","type":3,"url":"https://openphish.com/feed.txt"},"someonewhocares":{"desc":"Zero based host and domain list","disabled":false,"file":"","ip":"","name":"someonewhocares","prefix":"0.0.0.0","type":3,"url":"http://someonewhocares.org/hosts/zero/"},"tasty":{"desc":"File source","disabled":false,"file":"/config/user-data/blist.hosts.src","ip":"0.0.0.0","name":"tasty","prefix":"","type":3,"url":""},"volkerschatz":{"desc":"Ad server blacklists","disabled":false,"file":"","ip":"","name":"volkerschatz","prefix":"http","type":3,"url":"http://www.volkerschatz.com/net/adpaths"},"winhelp2002":{"desc":"Zero based host and domain list","disabled":false,"file":"","ip":"0.0.0.0","name":"winhelp2002","prefix":"0.0.0.0 ","type":3,"url":"http://winhelp2002.mvps.org/hosts.txt"},"yoyo":{"desc":"Fully Qualified Domain Names only - no prefix to strip","disabled":false,"file":"","ip":"","name":"yoyo","prefix":"","type":3,"url":"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml\u0026showintro=1\u0026mimetype=plaintext"}},"disable":false,"excludes":[],"includes":["beap.gemini.yahoo.com"],"ip":"192.168.168.1"}}`

	// FileManifest is complete list of the blacklist config node templates
	FileManifest = `../payload/blacklist
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

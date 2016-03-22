package config

var (
	// Testdata2 contains a valid partial EdgeOS blacklist configuration
	Testdata2 = `blacklist {
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

	// Testdata contains a valid full EdgeOS blacklist configuration
	Testdata = `blacklist {
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
                    prefix ""
                    url http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext
                }
            }
        }`

	// FileManifest is complete list of the blacklist config node templates
	FileManifest = `blacklist
blacklist/disabled
blacklist/disabled/node.def
blacklist/dns-redirect-ip
blacklist/dns-redirect-ip/node.def
blacklist/domains
blacklist/domains/dns-redirect-ip
blacklist/domains/dns-redirect-ip/node.def
blacklist/domains/exclude
blacklist/domains/exclude/node.def
blacklist/domains/include
blacklist/domains/include/node.def
blacklist/domains/node.def
blacklist/domains/source
blacklist/domains/source/node.def
blacklist/domains/source/node.tag
blacklist/domains/source/node.tag/description
blacklist/domains/source/node.tag/description/node.def
blacklist/domains/source/node.tag/prefix
blacklist/domains/source/node.tag/prefix/node.def
blacklist/domains/source/node.tag/url
blacklist/domains/source/node.tag/url/node.def
blacklist/exclude
blacklist/exclude/node.def
blacklist/hosts
blacklist/hosts/dns-redirect-ip
blacklist/hosts/dns-redirect-ip/node.def
blacklist/hosts/exclude
blacklist/hosts/exclude/node.def
blacklist/hosts/include
blacklist/hosts/include/node.def
blacklist/hosts/node.def
blacklist/hosts/source
blacklist/hosts/source/node.def
blacklist/hosts/source/node.tag
blacklist/hosts/source/node.tag/description
blacklist/hosts/source/node.tag/description/node.def
blacklist/hosts/source/node.tag/prefix
blacklist/hosts/source/node.tag/prefix/node.def
blacklist/hosts/source/node.tag/url
blacklist/hosts/source/node.tag/url/node.def
blacklist/node.def`
)

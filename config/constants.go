// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

package config

// Testdata2 contains a valid partial EdgeOS blacklist configuration
var Testdata2 = `blacklist {
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
var Testdata = `blacklist {
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


# blacklist

© 2016 NJ Software. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE.txt file.

[UBNT EdgeMax](<a href="https://community.ubnt.com/t5/EdgeMAX/bd-p/EdgeMAX">https://community.ubnt.com/t5/EdgeMAX/bd-p/EdgeMAX</a>) dnsmasq Blacklist and Adware Blocking

NOTE: THIS IS NOT OFFICIAL UBIQUITI SOFTWARE AND THEREFORE NOT SUPPORTED OR ENDORSED BY Ubiquiti Networks®

[![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/britannic/blacklist/blob/master/LICENSE.txt) [![Alpha  Version](https://img.shields.io/badge/version-v0.05--alpha-red.svg)](https://github.com/britannic/blacklist) [![GoDoc](https://godoc.org/github.com/britannic/blacklist?status.svg)](https://godoc.org/github.com/britannic/blacklist) [![Build Status](https://travis-ci.org/britannic/blacklist.svg?branch=master)](https://travis-ci.org/britannic/blacklist) [![Coverage Status](https://coveralls.io/repos/github/britannic/blacklist/badge.svg?branch=master)](https://coveralls.io/github/britannic/blacklist?branch=master) [![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/github.com/britannic/blacklist)


### Overview
EdgeMax dnsmasq Blacklist and Adware Blocking is derived from the received wisdom found at [Ubiquiti Community](<a href="https://community.ubnt.com/t5/EdgeMAX/bd-p/EdgeMAX">https://community.ubnt.com/t5/EdgeMAX/bd-p/EdgeMAX</a>)

### Features
Generates configuration files used directly by dnsmasq to redirect DNS lookups
Integrated with the EdgeMax OS CLI

### Any FQDN in the blacklist will force dnsmasq to return the configured DNS redirect IP address
Compatibility

blacklist has been tested on the EdgeRouter Lite family of routers, versions v1.6.0-v1.8.0.

The script will also install a default blacklist setup, here is the stanza (show service dns forwarding):


	blacklist {
	    disabled false
	    dns-redirect-ip 0.0.0.0
	    domains {
	        exclude adobedtm.com
	        exclude apple.com
	        exclude coremetrics.com
	        exclude doubleclick.net
	        exclude google.com
	        exclude googleadservices.com
	        exclude googleapis.com
	        exclude hulu.com
	        exclude msdn.com
	        exclude paypal.com
	        exclude storage.googleapis.com
	        include adsrvr.org
	        include adtechus.net
	        include advertising.com
	        include centade.com
	        include doubleclick.net
	        include free-counter.co.uk
	        include kiosked.com
	        source malc0de.com {
	            description "List of zones serving malicious executables observed by malc0de.com/database/"
	            prefix "zone "
	            url <a href="http://malc0de.com/bl/ZONES">http://malc0de.com/bl/ZONES</a>
	        }
	    }
	    hosts {
	        exclude appleglobal.112.2o7.net
	        exclude autolinkmaker.itunes.apple.com
	        exclude cdn.visiblemeasures.com
	        exclude freedns.afraid.org
	        exclude hb.disney.go.com
	        exclude static.chartbeat.com
	        exclude survey.112.2o7.net
	        exclude ads.hulu.com
	        exclude ads-a-darwin.hulu.com
	        exclude ads-v-darwin.hulu.com
	        exclude track.hulu.com
	        include beap.gemini.yahoo.com
	        source openphish.com {
	            description "OpenPhish automatic phishing detection"
	            prefix http
	            url <a href="https://openphish.com/feed.txt">https://openphish.com/feed.txt</a>
	        }
	        source someonewhocares.org {
	            description "Zero based host and domain list"
	            prefix 0.0.0.0
	            url <a href="http://someonewhocares.org/hosts/zero/">http://someonewhocares.org/hosts/zero/</a>
	        }
	        source volkerschatz.com {
	            description "Ad server blacklists"
	            prefix http
	            url <a href="http://www.volkerschatz.com/net/adpaths">http://www.volkerschatz.com/net/adpaths</a>
	        }
	        source winhelp2002.mvps.org {
	            description "Zero based host and domain list"
	            prefix "0.0.0.0 "
	            url <a href="http://winhelp2002.mvps.org/hosts.txt">http://winhelp2002.mvps.org/hosts.txt</a>
	        }
	        source www.malwaredomainlist.com {
	            description "127.0.0.1 based host and domain list"
	            prefix "127.0.0.1 "
	            url <a href="http://www.malwaredomainlist.com/hostslist/hosts.txt">http://www.malwaredomainlist.com/hostslist/hosts.txt</a>
	        }
	        source yoyo.org {
	            description "Fully Qualified Domain Names only - no prefix to strip"
	            prefix ""
	            url <a href="http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&amp;showintro=1&amp;mimetype=plaintext">http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext</a>
	        }
	    }
	}

CLI commands to configure blacklist:


	delete service dns forwarding blacklist
	delete system task-scheduler task update_blacklists
	set service dns forwarding blacklist dns-redirect-ip 0.0.0.0
	set service dns forwarding blacklist disabled false
	# set service dns forwarding blacklist dns-redirect-ip 192.168.168.1
	set service dns forwarding blacklist domains include adsrvr.org
	set service dns forwarding blacklist domains include adtechus.net
	set service dns forwarding blacklist domains include advertising.com
	set service dns forwarding blacklist domains include centade.com
	set service dns forwarding blacklist domains include doubleclick.net
	set service dns forwarding blacklist domains include free-counter.co.uk
	set service dns forwarding blacklist domains include intellitxt.com
	set service dns forwarding blacklist domains include kiosked.com
	set service dns forwarding blacklist domains include patoghee.in
	set service dns forwarding blacklist domains source malc0de description 'List of zones serving malicious executables observed by malc0de.com/database/'
	set service dns forwarding blacklist domains source malc0de prefix 'zone '
	set service dns forwarding blacklist domains source malc0de url '<a href="http://malc0de.com/bl/ZONES">http://malc0de.com/bl/ZONES</a>'
	set service dns forwarding blacklist exclude 122.2o7.net
	set service dns forwarding blacklist exclude 1e100.net
	set service dns forwarding blacklist exclude adobedtm.com
	set service dns forwarding blacklist exclude akamai.net
	set service dns forwarding blacklist exclude akamaihd.net
	set service dns forwarding blacklist exclude amazon.com
	set service dns forwarding blacklist exclude amazonaws.com
	set service dns forwarding blacklist exclude apple.com
	set service dns forwarding blacklist exclude ask.com
	set service dns forwarding blacklist exclude avast.com
	set service dns forwarding blacklist exclude bitdefender.com
	set service dns forwarding blacklist exclude cdn.visiblemeasures.com
	set service dns forwarding blacklist exclude cloudfront.net
	set service dns forwarding blacklist exclude coremetrics.com
	set service dns forwarding blacklist exclude edgesuite.net
	set service dns forwarding blacklist exclude freedns.afraid.org
	set service dns forwarding blacklist exclude github.com
	set service dns forwarding blacklist exclude githubusercontent.com
	set service dns forwarding blacklist exclude google.com
	set service dns forwarding blacklist exclude googleadservices.com
	set service dns forwarding blacklist exclude googleapis.com
	set service dns forwarding blacklist exclude googleusercontent.com
	set service dns forwarding blacklist exclude gstatic.com
	set service dns forwarding blacklist exclude gvt1.com
	set service dns forwarding blacklist exclude gvt1.net
	set service dns forwarding blacklist exclude hb.disney.go.com
	set service dns forwarding blacklist exclude hp.com
	set service dns forwarding blacklist exclude hulu.com
	set service dns forwarding blacklist exclude images-amazon.com
	set service dns forwarding blacklist exclude live.com
	set service dns forwarding blacklist exclude microsoft.com
	set service dns forwarding blacklist exclude msdn.com
	set service dns forwarding blacklist exclude paypal.com
	set service dns forwarding blacklist exclude rackcdn.com
	set service dns forwarding blacklist exclude schema.org
	set service dns forwarding blacklist exclude shopify.com
	set service dns forwarding blacklist exclude skype.com
	set service dns forwarding blacklist exclude smacargo.com
	set service dns forwarding blacklist exclude sourceforge.net
	set service dns forwarding blacklist exclude ssl-on9.com
	set service dns forwarding blacklist exclude ssl-on9.net
	set service dns forwarding blacklist exclude sstatic.net
	set service dns forwarding blacklist exclude static.chartbeat.com
	set service dns forwarding blacklist exclude storage.googleapis.com
	set service dns forwarding blacklist exclude windows.net
	set service dns forwarding blacklist exclude yimg.com
	set service dns forwarding blacklist exclude ytimg.com
	set service dns forwarding blacklist hosts include beap.gemini.yahoo.com
	set service dns forwarding blacklist hosts source adaway description 'Blocking mobile ad providers and some analytics providers'
	set service dns forwarding blacklist hosts source adaway prefix '127.0.0.1 '
	set service dns forwarding blacklist hosts source adaway url '<a href="http://adaway.org/hosts.txt">http://adaway.org/hosts.txt</a>'
	# set service dns forwarding blacklist hosts source hpHosts description 'Ad and Tracking servers only'
	# set service dns forwarding blacklist hosts source hpHosts prefix 127.0.0.1
	# set service dns forwarding blacklist hosts source hpHosts url '<a href="http://hosts-file.net/ad_servers.txt">http://hosts-file.net/ad_servers.txt</a>'
	set service dns forwarding blacklist hosts source malwaredomainlist description '127.0.0.1 based host and domain list'
	set service dns forwarding blacklist hosts source malwaredomainlist prefix '127.0.0.1 '
	set service dns forwarding blacklist hosts source malwaredomainlist url '<a href="http://www.malwaredomainlist.com/hostslist/hosts.txt">http://www.malwaredomainlist.com/hostslist/hosts.txt</a>'
	set service dns forwarding blacklist hosts source openphish description 'OpenPhish automatic phishing detection'
	set service dns forwarding blacklist hosts source openphish prefix http
	set service dns forwarding blacklist hosts source openphish url '<a href="https://openphish.com/feed.txt">https://openphish.com/feed.txt</a>'
	set service dns forwarding blacklist hosts source someonewhocares description 'Zero based host and domain list'
	set service dns forwarding blacklist hosts source someonewhocares prefix 0.0.0.0
	set service dns forwarding blacklist hosts source someonewhocares url '<a href="http://someonewhocares.org/hosts/zero/">http://someonewhocares.org/hosts/zero/</a>'
	set service dns forwarding blacklist hosts source volkerschatz description 'Ad server blacklists'
	set service dns forwarding blacklist hosts source volkerschatz prefix http
	set service dns forwarding blacklist hosts source volkerschatz url '<a href="http://www.volkerschatz.com/net/adpaths">http://www.volkerschatz.com/net/adpaths</a>'
	set service dns forwarding blacklist hosts source winhelp2002 description 'Zero based host and domain list'
	set service dns forwarding blacklist hosts source winhelp2002 prefix '0.0.0.0 '
	set service dns forwarding blacklist hosts source winhelp2002 url '<a href="http://winhelp2002.mvps.org/hosts.txt">http://winhelp2002.mvps.org/hosts.txt</a>'
	set service dns forwarding blacklist hosts source yoyo description 'Fully Qualified Domain Names only - no prefix to strip'
	set service dns forwarding blacklist hosts source yoyo prefix ''
	set service dns forwarding blacklist hosts source yoyo url '<a href="http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&amp;showintro=1&amp;mimetype=plaintext">http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext</a>'
	set system task-scheduler task update_blacklists executable path /config/scripts/update-dnsmasq.pl
	set system task-scheduler task update_blacklists interval 1d

Notes:

In order to make this work properly, you will need to first ensure that your dnsmasq is correctly set up. An example configuration is posted below:


	show service dns forwarding
	 cache-size 2048
	 listen-on eth0
	 listen-on eth2
	 listen-on lo
	 name-server 208.67.220.220
	 name-server 208.67.222.222
	 name-server 2620:0:ccc::2
	 name-server 2620:0:ccd::2
	 options expand-hosts
	 options bogus-priv
	 options localise-queries
	 options domain=ubnt.home
	 options strict-order
	 options listen-address=127.0.0.1
	 system








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
blacklist {
        disabled false
        dns-redirect-ip 192.168.168.1
        domains {
            source NoBitCoin {
                description "Blocking Web Browser Bitcoin Mining"
                prefix 0.0.0.0
                url https://raw.githubusercontent.com/hoshsadiq/adblock-nocoin-list/master/hosts.txt
            }
            source malc0de {
                description "List of zones serving malicious executables observed by malc0de.com/database/"
                prefix zone
                url http://malc0de.com/bl/ZONES
            }
            source malwaredomains.com {
                description "Just Domains"
                url http://mirror1.malwaredomains.com/files/justdomains
            }
            source simple_tracking {
                description "Basic tracking list by Disconnect"
                url https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt
            }
            source zeus {
                description "abuse.ch ZeuS domain blocklist"
                url https://zeustracker.abuse.ch/blocklist.php?download=domainblocklist
            }
        source tasty {
            description "File source"
            dns-redirect-ip 10.10.10.10
            file ./internal/testdata/blist.hosts.src
            }
        }
        exclude 1e100.net
        exclude 2o7.net
        exclude adobedtm.com
        exclude akamai.net
        exclude akamaihd.net
        exclude amazon.com
        exclude amazonaws.com
        exclude apple.com
        exclude ask.com
        exclude avast.com
        exclude avira-update.com
        exclude bannerbank.com
        exclude bing.com
        exclude bit.ly
        exclude bitdefender.com
        exclude cdn.ravenjs.com
        exclude cdn.visiblemeasures.com
        exclude cloudfront.net
        exclude coremetrics.com
        exclude dropbox.com
        exclude ebay.com
        exclude edgesuite.net
        exclude evernote.com
        exclude express.co.uk
        exclude feedly.com
        exclude freedns.afraid.org
        exclude github.com
        exclude githubusercontent.com
        exclude global.ssl.fastly.net
        exclude google.com
        exclude googleads.g.doubleclick.net
        exclude googleadservices.com
        exclude googleapis.com
        exclude googletagmanager.com
        exclude googleusercontent.com
        exclude gstatic.com
        exclude gvt1.com
        exclude gvt1.net
        exclude hb.disney.go.com
        exclude herokuapp.com
        exclude hp.com
        exclude hulu.com
        exclude images-amazon.com
        exclude live.com
        exclude magnetmail1.net
        exclude microsoft.com
        exclude microsoftonline.com
        exclude msdn.com
        exclude msecnd.net
        exclude msftncsi.com
        exclude mywot.com
        exclude nsatc.net
        exclude paypal.com
        exclude pop.h-cdn.co
        exclude rackcdn.com
        exclude rarlab.com
        exclude schema.org
        exclude shopify.com
        exclude skype.com
        exclude smacargo.com
        exclude sourceforge.net
        exclude spotify.com
        exclude spotify.edgekey.net
        exclude spotilocal.com
        exclude ssl-on9.com
        exclude ssl-on9.net
        exclude sstatic.net
        exclude static.chartbeat.com
        exclude storage.googleapis.com
        exclude twimg.com
        exclude viewpoint.com
        exclude windows.net
        exclude xboxlive.com
        exclude yimg.com
        exclude ytimg.com
        include adk2x.com
        include adsrvr.org
        include adtechus.net
        include advertising.com
        include centade.com
        include doubleclick.net
        include fastplayz.com
        include free-counter.co.uk
        include hilltopads.net
        include intellitxt.com
        include kiosked.com
        include patoghee.in
        include themillionaireinpjs.com
        include traktrafficflow.com
        include wwwpromoter.com
        hosts {
            include ads.feedly.com
            include beap.gemini.yahoo.com
            source githubSteveBlack {
                description "Blacklists adware and malware websites"
                prefix 0.0.0.0
                url https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts
            }
            source hostsfile.org {
                description "hostsfile.org bad hosts blacklist"
                prefix 127.0.0.1
                url http://www.hostsfile.org/Downloads/hosts.txt
            }
            source openphish {
                description "OpenPhish automatic phishing detection"
                prefix http
                url https://openphish.com/feed.txt
            }
            source sysctl.org {
                description "This hosts file is a merged collection of hosts from Cameleon"
                prefix 127.0.0.1
                url http://sysctl.org/cameleon/hosts
            }
        }
    }
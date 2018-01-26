interfaces {
    ethernet eth0 {
        address dhcp
        address dhcpv6
        description "External WAN"
        dhcp-options {
            default-route update
            default-route-distance 210
            name-server no-update
        }
        duplex auto
        speed auto
    }
    ethernet eth1 {
        address 192.168.150.1/24
        description "Internal LAN"
        duplex auto
        mtu 1500
        speed auto
    }
    ethernet eth2 {
        address 192.168.200.1/24
        description DMZ
        duplex auto
        mtu 1500
        speed auto
    }
    ethernet eth3 {
        duplex auto
        speed auto
    }
    ethernet eth4 {
        duplex auto
        speed auto
    }
    loopback lo {
    }
    switch switch0 {
        mtu 1500
    }
}
service {
    dhcp-server {
        disabled false
        hostfile-update enable
        shared-network-name LAN0 {
            authoritative enable
            subnet 192.168.150.0/24 {
                default-router 192.168.150.1
                dns-server 192.168.150.1
                domain-name er-x.local
                lease 86400
                start 192.168.150.100 {
                    stop 192.168.150.200
                }
            }
        }
        static-arp disable
        use-dnsmasq disable
    }
    dns {
        forwarding {
            blacklist {
                disabled false
                dns-redirect-ip 0.0.0.0
                domains {
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
                exclude freedns.afraid.org
                exclude github.com
                exclude githubusercontent.com
                exclude global.ssl.fastly.net
                exclude google.com
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
                exclude microsoft.com
                exclude microsoftonline.com
                exclude msdn.com
                exclude msecnd.net
                exclude msftncsi.com
                exclude mywot.com
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
                exclude viewpoint.com
                exclude windows.net
                exclude xboxlive.com
                exclude yimg.com
                exclude ytimg.com
                hosts {
                    include beap.gemini.yahoo.com
                    source openphish {
                        description "OpenPhish automatic phishing detection"
                        prefix http
                        url https://openphish.com/feed.txt
                    }
                    source raw.github.com {
                        description "This hosts file is a merged collection of hosts from reputable sources"
                        prefix 0.0.0.0
                        url https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts
                    }
                    source sysctl.org {
                        description "This hosts file is a merged collection of hosts from Cameleon"
                        prefix 127.0.0.1
                        url http://sysctl.org/cameleon/hosts
                    }
                    source yoyo {
                        description "Fully Qualified Domain Names only - no prefix to strip"
                        prefix 127.0.0.1
                        url http://pgl.yoyo.org/as/serverlist.php
                    }
                }
            }
            cache-size 150
            listen-on eth1
            listen-on eth2
            name-server 208.67.220.220
            options bogus-priv
            options domain=er-x.local
            options expand-hosts
            options listen-address=::1
            options listen-address=127.0.0.1
            options localise-queries
            options strict-order
        }
    }
    gui {
        http-port 80
        https-port 443
        listen-address 192.168.150.1
        older-ciphers enable
    }
    ssh {
        disable-password-authentication
        port 22
        protocol-version v2
    }
    unms {
        connection wss://unifi.x.com:443+CWhNxTneTNDLolDE_XirvZwLwJDQwfnxXD1ZYS0SwYAAAAA+allowUntrustedCertificate
    }
}
system {
    conntrack {
        expect-table-size 4096
        hash-size 4096
        ignore {
            rule 10 {
                destination {
                    address 255.255.255.255
                }
            }
        }
        table-size 262144
    }
    domain-search {
        domain ashcreek.home
    }
    host-name er-x
    ip {
        override-hostname-ip 192.168.150.1
    }
    login {
        banner {
            post-login "\nWelcome to EdgeOS!\n"
            pre-login "\n\n\n\tWARNING *** WARNING *** WARNING *** WARNING *** WARNING\n\n\n\tWARNING: Criminal and civil penalties may be imposed for obtaining\n\tunauthorized access to this system or for causing intentional,\n\tunauthorized damage, deletion, alteration, or insertion of data.\n\tAny information stored, processed, or transmitted to this system\n\tmay be monitored, used, or disclosed by authorized personnel,\n\tincluding law enforcement. Email sysadmin@empirecreekcircle.com\n\tto gain access to this equipment if you need authorization.\n\n\n"
        }
        user nbnt {
            authentication {
                encrypted-password $6$zIyYjCe858Qmu1.mAfEHw4VSuSavlEIKbhQdvzz3qXJs/ygC8Jd0kaRaparu5eJI0T05iI2uvICN1xONowGDoxTwu/
                plaintext-password ""
                public-keys root@mainframe.local {
                    key AAAAB3NQABAAACAQDJOzuwOndcN1dDlgTPCkIVEJnDN07wmzUhxHEjWwBblHv/P3eSv9AQqBA46y1thlJyecBrEehgpXUECFUry2MEwyLcyUS7zPV/5zxRIAKCMRm+dhjfqwAFp7EiodYPO06dGWElQ0ND0+3MWJSd9TU+unnxipbGAySpobZtffNPdatJ0il5XxJULVgh3leyCqrLZQIYxK9X1wJWfD292mKJphoeYjlFdkBYq9oruG1cYyZrtt6x3Rf1Yp6/qQ2GHjAfumBT9AzgNeL6XbxiQGmOybdnNxfmqbKD5BIArXHaz4qysErhHwvqIbS/U5MHNPkEn/4OPbXEPkMqxWGPWuHhOokilIyFe8hdk0b5rJ9bYDn8uATtaXGL5nT1MjlLHTrWtkpH8eKfC/+hyXKZMNo2zpMr+b2gf769K7MgbIS6BxUsS0U43PneelO1G8RQZ1fjRy6WenA09oute2ARnvgMdcjqA+fzY8fgGpR837mXsOjW2zCP1pF/CYm56EyIG6cLp95dEt7VXclvIR/sTaOejE3jIXLW3mqsYkKxzCx6JCBDpi+Re6K2hEaHW80LvgMWs4wv3bjqsn5+RpctR7csp+tpbw6XYXpt0WsJDD0KGQ0t8KG9OVSQZ7xPM3T2SObAPQd/CRJlWK9ZE+Qk6xYWmWk65/1KTSKP+0qW7ZGETw==
                    type ssh-rsa
                }
            }
            full-name Admin
            level admin
        }
    }
    name-server 127.0.0.1
    ntp {
        server 0.ubnt.pool.ntp.org {
        }
        server 1.ubnt.pool.ntp.org {
        }
    }
    package {
        repository wheezy {
            components "main contrib non-free"
            distribution wheezy
            password ""
            url http://http.us.debian.org/debian/
            username ""
        }
    }
    syslog {
        global {
            archive {
                files 10
                size 250
            }
            facility all {
                level notice
            }
            facility cron {
                level err
            }
            facility protocols {
                level debug
            }
        }
    }
    task-scheduler {
        task update_blacklists {
            executable {
                path /config/scripts/update-dnsmasq
            }
            interval 1d
        }
    }
    time-zone America/Los_Angeles
}


/* Warning: Do not remove the following line. */
/* === vyatta-config-version: "config-management@1:conntrack@1:cron@1:dhcp-relay@1:dhcp-server@4:firewall@5:ipsec@5:nat@3:qos@1:quagga@2:system@4:ubnt-pptp@1:ubnt-udapi-server@1:ubnt-unms@1:ubnt-util@1:vrrp@1:webgui@1:webproxy@1:zone-policy@1" === */
/* Release version: v1.10.0-beta.3.5051713.180109.1605 */

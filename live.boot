firewall {
    all-ping enable
    broadcast-ping disable
    group {
        address-group media {
            address 192.168.4.255
            address 192.168.10.30
            address 192.168.10.35
            address 224.0.0.251
            address 239.255.255.250
            address 255.255.255.255
            address 192.168.50.1-192.168.50.255
            description "Media Address Group"
        }
        ipv6-network-group IPv6Bogons {
            description "IPv6 BOGON Addresses"
            ipv6-network ::/127
            ipv6-network ::ffff:0:0/96
            ipv6-network ::/96
            ipv6-network 100::/64
            ipv6-network 2001:10::/28
            ipv6-network 2001:db8::/32
            ipv6-network fc00::/7
            ipv6-network fec0::/10
            ipv6-network ff00::/8
            ipv6-network 2002::/24
            ipv6-network 2002:a00::/24
            ipv6-network 2002:7f00::/24
            ipv6-network 2002:a9fe::/32
            ipv6-network 2002:ac10::/28
            ipv6-network 2002:c000::/40
            ipv6-network 2002:c000:200::/40
            ipv6-network 2002:c0a8::/32
            ipv6-network 2002:c612::/31
            ipv6-network 2002:c633:6400::/40
            ipv6-network 2002:cb00:7100::/40
            ipv6-network 2002:e000::/20
            ipv6-network 2002:f000::/20
            ipv6-network 2002:ffff:ffff::/48
            ipv6-network 2001::/40
            ipv6-network 2001:0:a00::/40
            ipv6-network 2001:0:7f00::/40
            ipv6-network 2001:0:a9fe::/48
            ipv6-network 2001:0:ac10::/44
            ipv6-network 2001:0:c000::/56
            ipv6-network 2001:0:c000:200::/56
            ipv6-network 2001:0:c0a8::/48
            ipv6-network 2001:0:c612::/47
            ipv6-network 2001:0:c633:6400::/56
            ipv6-network 2001:0:cb00:7100::/56
            ipv6-network 2001:0:e000::/36
            ipv6-network 2001:0:f000::/36
            ipv6-network 2001:0:ffff:ffff::/64
        }
        ipv6-network-group IPv6media {
            description "IPv6 Media Network"
            ipv6-network 2601:1c0:7000:fcb7::/64
            ipv6-network fe80::2e9e:fcff:fe4f:0/127
        }
        network-group IPv4Bogons {
            description "IPv4 BOGON Addresses"
            network 10.0.0.0/8
            network 100.64.0.0/10
            network 127.0.0.0/8
            network 169.254.0.0/16
            network 172.16.0.0/12
            network 192.0.0.0/24
            network 192.0.2.0/24
            network 192.168.0.0/16
            network 198.18.0.0/15
            network 198.51.100.0/24
            network 203.0.113.0/24
            network 224.0.0.0/4
            network 240.0.0.0/4
        }
        port-group dhcpIPv4 {
            description "DHCP Port Group"
            port bootpc
            port bootps
        }
        port-group dhcpIPv6 {
            description "DHCP Port Group"
            port dhcpv6-client
            port dhcpv6-server
        }
        port-group email {
            description "Email Port Group"
            port imap2
            port imaps
            port smtp
            port ssmtp
        }
        port-group ftp {
            description "FTP Port Group"
            port ftp-data
            port ftp
            port ftps-data
            port ftps
            port sftp
        }
        port-group print {
            description "Print Port Group"
            port 1900
            port 3702
            port 5000
            port 5001
            port 5222
            port 5357
            port 8000
            port 8610
            port 8611
            port 8612
            port 8613
            port 9000
            port 9100
            port 9200
            port 9300
            port 9500
            port 9600
            port 9700
            port http
            port https
            port ipp
            port netbios-dgm
            port netbios-ns
            port netbios-ssn
            port printer
            port snmp-trap
            port snmp
        }
        port-group ssdp {
            description "SSDP Port Group"
            port 10102
            port 1900
            port 5354
            port afpovertcp
            port http
            port https
            port mdns
            port netbios-ns
        }
        port-group voip {
            description "VoIP Port Group"
            port 5060
            port 6060
            port 1026
        }
        port-group vpn {
            description "VPN Port Group"
            port isakmp
            port openvpn
            port l2tp
            port 4500
        }
    }
    ipv6-name ipv6-dmz-ext {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 1500 {
            action drop
            description "Block MDNS & SSDP access to Internet"
            destination {
                port mdns
            }
            protocol udp
        }
    }
    ipv6-name ipv6-dmz-gst {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
    }
    ipv6-name ipv6-dmz-int {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
    }
    ipv6-name ipv6-dmz-loc {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 400 {
            action drop
            description "Block mdx media address group"
            destination {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 1000 {
            action accept
            description "Permit access to local DNS"
            destination {
                port domain
            }
            protocol tcp_udp
        }
    }
    ipv6-name ipv6-dmz-mdx {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
    }
    ipv6-name ipv6-dmz-pbx {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
    }
    ipv6-name ipv6-ext-dmz {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action drop
            description "Drop IPv6 bogons"
            source {
                group {
                    ipv6-network-group IPv6Bogons
                }
            }
        }
        rule 3000 {
            action drop
            description "Drop brute force SSH from Internet"
            destination {
                port ssh
            }
            protocol tcp
            recent {
                count 3
                time 30
            }
        }
        rule 5000 {
            action accept
            description "Allow vpn traffic"
            destination {
                group {
                    port-group vpn
                }
            }
            protocol udp
        }
        rule 5500 {
            action accept
            description "Allow vpn PPTP"
            destination {
                port 1723
            }
            protocol tcp
        }
        rule 5600 {
            action accept
            description "Allow vpn ESP"
            protocol esp
        }
    }
    ipv6-name ipv6-ext-gst {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action drop
            description "Drop IPv6 bogons"
            source {
                group {
                    ipv6-network-group IPv6Bogons
                }
            }
        }
        rule 3000 {
            action drop
            description "Drop brute force SSH from Internet"
            destination {
                port ssh
            }
            protocol tcp
            recent {
                count 3
                time 30
            }
        }
    }
    ipv6-name ipv6-ext-int {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action drop
            description "Drop IPv6 bogons"
            source {
                group {
                    ipv6-network-group IPv6Bogons
                }
            }
        }
        rule 3000 {
            action drop
            description "Drop brute force SSH from Internet"
            destination {
                port ssh
            }
            protocol tcp
            recent {
                count 3
                time 30
            }
        }
    }
    ipv6-name ipv6-ext-loc {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action drop
            description "Drop IPv6 bogons"
            source {
                group {
                    ipv6-network-group IPv6Bogons
                }
            }
        }
        rule 6 {
            action accept
            description "Allow DHCPV6 responses from ISP"
            destination {
                port dhcpv6-client
            }
            protocol udp
            source {
                address fe80::/64
                port dhcpv6-server
            }
        }
        rule 500 {
            action drop
            description "Block IPv6-ICMP ping from the Internet"
            icmpv6 {
                type ping
            }
            protocol icmpv6
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 3000 {
            action drop
            description "Drop brute force SSH from Internet"
            destination {
                port ssh
            }
            protocol tcp
            recent {
                count 3
                time 30
            }
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
        rule 5000 {
            action accept
            description "Allow vpn traffic"
            destination {
                group {
                    port-group vpn
                }
            }
            protocol udp
        }
        rule 5500 {
            action accept
            description "Allow vpn PPTP"
            destination {
                port 1723
            }
            protocol tcp
        }
        rule 5600 {
            action accept
            description "Allow vpn ESP"
            protocol esp
        }
    }
    ipv6-name ipv6-ext-mdx {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action drop
            description "Drop IPv6 bogons"
            source {
                group {
                    ipv6-network-group IPv6Bogons
                }
            }
        }
        rule 3000 {
            action drop
            description "Drop brute force SSH from Internet"
            destination {
                port ssh
            }
            protocol tcp
            recent {
                count 3
                time 30
            }
        }
    }
    ipv6-name ipv6-ext-pbx {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action drop
            description "Drop IPv6 bogons"
            source {
                group {
                    ipv6-network-group IPv6Bogons
                }
            }
        }
        rule 3000 {
            action drop
            description "Drop brute force SSH from Internet"
            destination {
                port ssh
            }
            protocol tcp
            recent {
                count 3
                time 30
            }
        }
    }
    ipv6-name ipv6-gst-dmz {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
    }
    ipv6-name ipv6-gst-ext {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 1500 {
            action drop
            description "Block MDNS & SSDP access to Internet"
            destination {
                port mdns
            }
            protocol udp
        }
    }
    ipv6-name ipv6-gst-int {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    ipv6-network-group IPv6media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
        rule 2200 {
            action accept
            description "Permit Printer access"
            destination {
                group {
                    ipv6-network-group IPv6media
                    port-group print
                }
            }
            protocol tcp_udp
        }
    }
    ipv6-name ipv6-gst-loc {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 400 {
            action drop
            description "Block mdx media address group"
            destination {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 1000 {
            action accept
            description "Permit access to local DNS"
            destination {
                port domain
            }
            protocol tcp_udp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    ipv6-network-group IPv6media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
    }
    ipv6-name ipv6-gst-mdx {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    ipv6-network-group IPv6media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
    }
    ipv6-name ipv6-gst-pbx {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
    }
    ipv6-name ipv6-int-dmz {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
        rule 5000 {
            action accept
            description "Allow vpn traffic"
            destination {
                group {
                    port-group vpn
                }
            }
            protocol udp
        }
        rule 5500 {
            action accept
            description "Allow vpn PPTP"
            destination {
                port 1723
            }
            protocol tcp
        }
        rule 5600 {
            action accept
            description "Allow vpn ESP"
            protocol esp
        }
        rule 6000 {
            action accept
            description "Allow ADT Camera streams"
            destination {
                port 4301-4325
            }
            log enable
            protocol tcp_udp
        }
    }
    ipv6-name ipv6-int-ext {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 1500 {
            action drop
            description "Block MDNS & SSDP access to Internet"
            destination {
                port mdns
            }
            protocol udp
        }
    }
    ipv6-name ipv6-int-gst {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 400 {
            action accept
            description "Allow mdx to offer access to media address group"
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    ipv6-network-group IPv6media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
    }
    ipv6-name ipv6-int-loc {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow all connections"
        }
        rule 5 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 400 {
            action accept
            description "Allow mdx to offer access to media address group"
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 1000 {
            action accept
            description "Permit access to local DNS"
            destination {
                port domain
            }
            protocol tcp_udp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    ipv6-network-group IPv6media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
        rule 5000 {
            action accept
            description "Allow vpn traffic"
            destination {
                group {
                    port-group vpn
                }
            }
            protocol udp
        }
        rule 5500 {
            action accept
            description "Allow vpn PPTP"
            destination {
                port 1723
            }
            protocol tcp
        }
        rule 5600 {
            action accept
            description "Allow vpn ESP"
            protocol esp
        }
    }
    ipv6-name ipv6-int-mdx {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 3 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 5 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    ipv6-network-group IPv6media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
    }
    ipv6-name ipv6-int-pbx {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
    }
    ipv6-name ipv6-loc-dmz {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
    }
    ipv6-name ipv6-loc-ext {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 1500 {
            action drop
            description "Block MDNS & SSDP access to Internet"
            destination {
                port mdns
            }
            protocol udp
        }
    }
    ipv6-name ipv6-loc-gst {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 400 {
            action accept
            description "Allow mdx to offer access to media address group"
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    ipv6-network-group IPv6media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
    }
    ipv6-name ipv6-loc-int {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    ipv6-network-group IPv6media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
    }
    ipv6-name ipv6-loc-mdx {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    ipv6-network-group IPv6media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
    }
    ipv6-name ipv6-loc-pbx {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
    }
    ipv6-name ipv6-mdx-dmz {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
    }
    ipv6-name ipv6-mdx-ext {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 1500 {
            action drop
            description "Block MDNS & SSDP access to Internet"
            destination {
                port mdns
            }
            protocol udp
        }
    }
    ipv6-name ipv6-mdx-gst {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 400 {
            action accept
            description "Allow mdx to offer access to media address group"
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    ipv6-network-group IPv6media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
    }
    ipv6-name ipv6-mdx-int {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 3 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 5 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    ipv6-network-group IPv6media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
    }
    ipv6-name ipv6-mdx-loc {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 400 {
            action accept
            description "Allow mdx to offer access to media address group"
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 1000 {
            action accept
            description "Permit access to local DNS"
            destination {
                port domain
            }
            protocol tcp_udp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    ipv6-network-group IPv6media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    ipv6-network-group IPv6media
                }
            }
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
    }
    ipv6-name ipv6-mdx-pbx {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
    }
    ipv6-name ipv6-pbx-dmz {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
    }
    ipv6-name ipv6-pbx-ext {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 1500 {
            action drop
            description "Block MDNS & SSDP access to Internet"
            destination {
                port mdns
            }
            protocol udp
        }
    }
    ipv6-name ipv6-pbx-gst {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
    }
    ipv6-name ipv6-pbx-int {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
    }
    ipv6-name ipv6-pbx-loc {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
        rule 1000 {
            action accept
            description "Permit access to local DNS"
            destination {
                port domain
            }
            protocol tcp_udp
        }
    }
    ipv6-name ipv6-pbx-mdx {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV6 responses"
            destination {
                group {
                    port-group dhcpIPv6
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow IPv6-ICMP"
            protocol icmpv6
        }
    }
    ipv6-receive-redirects disable
    ipv6-src-route disable
    ip-src-route disable
    log-martians enable
    name dmz-ext {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 1500 {
            action drop
            description "Block MDNS & SSDP access to Internet"
            destination {
                port mdns
            }
            protocol udp
        }
    }
    name dmz-gst {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
    }
    name dmz-int {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
    }
    name dmz-loc {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 300 {
            action accept
            description "Permit access to pixel server"
            destination {
                address 192.168.168.1
            }
            protocol tcp
        }
        rule 400 {
            action drop
            description "Block mdx media address group"
            destination {
                group {
                    address-group media
                }
            }
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 1000 {
            action accept
            description "Permit access to local DNS"
            destination {
                port domain
            }
            protocol tcp_udp
        }
    }
    name dmz-mdx {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
    }
    name dmz-pbx {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 2 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
    }
    name ext-dmz {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action drop
            description "Drop IPv4 bogons"
            source {
                group {
                    network-group IPv4Bogons
                }
            }
        }
        rule 3000 {
            action drop
            description "Drop brute force SSH from Internet"
            destination {
                port ssh
            }
            protocol tcp
            recent {
                count 3
                time 30
            }
        }
        rule 5000 {
            action accept
            description "Allow vpn traffic"
            destination {
                group {
                    port-group vpn
                }
            }
            protocol udp
        }
        rule 5500 {
            action accept
            description "Allow vpn PPTP"
            destination {
                port 1723
            }
            protocol tcp
        }
        rule 5600 {
            action accept
            description "Allow vpn ESP"
            protocol esp
        }
    }
    name ext-gst {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action drop
            description "Drop IPv4 bogons"
            source {
                group {
                    network-group IPv4Bogons
                }
            }
        }
        rule 3000 {
            action drop
            description "Drop brute force SSH from Internet"
            destination {
                port ssh
            }
            protocol tcp
            recent {
                count 3
                time 30
            }
        }
    }
    name ext-int {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action drop
            description "Drop IPv4 bogons"
            source {
                group {
                    network-group IPv4Bogons
                }
            }
        }
        rule 3000 {
            action drop
            description "Drop brute force SSH from Internet"
            destination {
                port ssh
            }
            protocol tcp
            recent {
                count 3
                time 30
            }
        }
    }
    name ext-loc {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action drop
            description "Drop IPv4 bogons"
            source {
                group {
                    network-group IPv4Bogons
                }
            }
        }
        rule 5 {
            action accept
            description "Allow DHCPV4 responses from ISP"
            destination {
                port bootpc
            }
            protocol udp
            source {
                port bootps
            }
        }
        rule 500 {
            action drop
            description "Block ICMP ping from the Internet"
            icmp {
                type-name ping
            }
            protocol icmp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 3000 {
            action drop
            description "Drop brute force SSH from Internet"
            destination {
                port ssh
            }
            protocol tcp
            recent {
                count 3
                time 30
            }
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
        rule 5000 {
            action accept
            description "Allow vpn traffic"
            destination {
                group {
                    port-group vpn
                }
            }
            protocol udp
        }
        rule 5500 {
            action accept
            description "Allow vpn PPTP"
            destination {
                port 1723
            }
            protocol tcp
        }
        rule 5600 {
            action accept
            description "Allow vpn ESP"
            protocol esp
        }
    }
    name ext-mdx {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action drop
            description "Drop IPv4 bogons"
            source {
                group {
                    network-group IPv4Bogons
                }
            }
        }
        rule 3000 {
            action drop
            description "Drop brute force SSH from Internet"
            destination {
                port ssh
            }
            protocol tcp
            recent {
                count 3
                time 30
            }
        }
    }
    name ext-pbx {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action drop
            description "Drop IPv4 bogons"
            source {
                group {
                    network-group IPv4Bogons
                }
            }
        }
        rule 3000 {
            action drop
            description "Drop brute force SSH from Internet"
            destination {
                port ssh
            }
            protocol tcp
            recent {
                count 3
                time 30
            }
        }
    }
    name gst-dmz {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
    }
    name gst-ext {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 1500 {
            action drop
            description "Block MDNS & SSDP access to Internet"
            destination {
                port mdns
            }
            protocol udp
        }
    }
    name gst-int {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    address-group media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    address-group media
                }
            }
        }
        rule 2100 {
            action accept
            description "Permit IGMP"
            destination {
                group {
                    address-group media
                }
            }
            protocol igmp
        }
        rule 2200 {
            action accept
            description "Permit Printer access"
            destination {
                group {
                    address-group media
                    port-group print
                }
            }
            protocol tcp_udp
        }
    }
    name gst-loc {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 300 {
            action accept
            description "Permit access to pixel server"
            destination {
                address 192.168.168.1
            }
            protocol tcp
        }
        rule 400 {
            action drop
            description "Block mdx media address group"
            destination {
                group {
                    address-group media
                }
            }
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 1000 {
            action accept
            description "Permit access to local DNS"
            destination {
                port domain
            }
            protocol tcp_udp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    address-group media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    address-group media
                }
            }
        }
    }
    name gst-mdx {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    address-group media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    address-group media
                }
            }
        }
        rule 2100 {
            action accept
            description "Permit IGMP"
            destination {
                group {
                    address-group media
                }
            }
            protocol igmp
        }
    }
    name gst-pbx {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 2 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
    }
    name int-dmz {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
        rule 5000 {
            action accept
            description "Allow vpn traffic"
            destination {
                group {
                    port-group vpn
                }
            }
            protocol udp
        }
        rule 5500 {
            action accept
            description "Allow vpn PPTP"
            destination {
                port 1723
            }
            protocol tcp
        }
        rule 5600 {
            action accept
            description "Allow vpn ESP"
            protocol esp
        }
        rule 6000 {
            action accept
            description "Allow ADT Camera streams"
            destination {
                port 4301-4325
            }
            log enable
            protocol tcp_udp
        }
    }
    name int-ext {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 1500 {
            action drop
            description "Block MDNS & SSDP access to Internet"
            destination {
                port mdns
            }
            protocol udp
        }
    }
    name int-gst {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 400 {
            action accept
            description "Allow mdx to offer access to media address group"
            source {
                group {
                    address-group media
                }
            }
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    address-group media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    address-group media
                }
            }
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
    }
    name int-loc {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow all connections"
        }
        rule 4 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 300 {
            action accept
            description "Permit access to pixel server"
            destination {
                address 192.168.168.1
            }
            protocol tcp
        }
        rule 400 {
            action accept
            description "Allow mdx to offer access to media address group"
            source {
                group {
                    address-group media
                }
            }
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 1000 {
            action accept
            description "Permit access to local DNS"
            destination {
                port domain
            }
            protocol tcp_udp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    address-group media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    address-group media
                }
            }
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
        rule 5000 {
            action accept
            description "Allow vpn traffic"
            destination {
                group {
                    port-group vpn
                }
            }
            protocol udp
        }
        rule 5500 {
            action accept
            description "Allow vpn PPTP"
            destination {
                port 1723
            }
            protocol tcp
        }
        rule 5600 {
            action accept
            description "Allow vpn ESP"
            protocol esp
        }
    }
    name int-mdx {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 3 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    address-group media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    address-group media
                }
            }
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
    }
    name int-pbx {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 2 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
    }
    name loc-dmz {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
    }
    name loc-ext {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 1500 {
            action drop
            description "Block MDNS & SSDP access to Internet"
            destination {
                port mdns
            }
            protocol udp
        }
    }
    name loc-gst {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 400 {
            action accept
            description "Allow mdx to offer access to media address group"
            source {
                group {
                    address-group media
                }
            }
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    address-group media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    address-group media
                }
            }
        }
    }
    name loc-int {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    address-group media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    address-group media
                }
            }
        }
    }
    name loc-mdx {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    address-group media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    address-group media
                }
            }
        }
    }
    name loc-pbx {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
    }
    name mdx-dmz {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
    }
    name mdx-ext {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 1500 {
            action drop
            description "Block MDNS & SSDP access to Internet"
            destination {
                port mdns
            }
            protocol udp
        }
    }
    name mdx-gst {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 400 {
            action accept
            description "Allow mdx to offer access to media address group"
            source {
                group {
                    address-group media
                }
            }
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    address-group media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    address-group media
                }
            }
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
    }
    name mdx-int {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 3 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 4 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    address-group media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    address-group media
                }
            }
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
    }
    name mdx-loc {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow established connections"
            state {
                established enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 3 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 300 {
            action accept
            description "Permit access to pixel server"
            destination {
                address 192.168.168.1
            }
            protocol tcp
        }
        rule 400 {
            action accept
            description "Allow mdx to offer access to media address group"
            source {
                group {
                    address-group media
                }
            }
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 1000 {
            action accept
            description "Permit access to local DNS"
            destination {
                port domain
            }
            protocol tcp_udp
        }
        rule 2000 {
            action accept
            description "Permit MDNS & SSDP access"
            destination {
                group {
                    address-group media
                    port-group ssdp
                }
            }
            protocol tcp_udp
            source {
                group {
                    address-group media
                }
            }
        }
        rule 3100 {
            action accept
            description "Allow SSH"
            destination {
                port ssh
            }
            protocol tcp
        }
    }
    name mdx-pbx {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 2 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
    }
    name pbx-dmz {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 2 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
    }
    name pbx-ext {
        default-action drop
        enable-default-log
        rule 1 {
            action accept
            description "Allow all connections"
            state {
                established enable
                new enable
                related enable
            }
        }
        rule 2 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 1500 {
            action drop
            description "Block MDNS & SSDP access to Internet"
            destination {
                port mdns
            }
            protocol udp
        }
    }
    name pbx-gst {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 2 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
    }
    name pbx-int {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 2 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
    }
    name pbx-loc {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 2 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
        rule 1000 {
            action accept
            description "Permit access to local DNS"
            destination {
                port domain
            }
            protocol tcp_udp
        }
    }
    name pbx-mdx {
        default-action drop
        enable-default-log
        rule 1 {
            action drop
            description "Drop invalid packets"
            state {
                invalid enable
            }
        }
        rule 2 {
            action accept
            description "Allow DHCPV4 responses"
            destination {
                group {
                    port-group dhcpIPv4
                }
            }
            protocol udp
        }
        rule 510 {
            action accept
            description "Allow ICMP"
            protocol icmp
        }
    }
    options {
        mss-clamp {
            interface-type all
            mss 1452
        }
        mss-clamp6 {
            interface-type all
            mss 1452
        }
    }
    receive-redirects disable
    send-redirects enable
    source-validation disable
    syn-cookies enable
}
interfaces {
    ethernet eth0 {
        address 192.168.100.1/24
        description "Internal LAN"
        duplex auto
        ip {
            enable-proxy-arp
        }
        ipv6 {
            dup-addr-detect-transmits 1
            router-advert {
                cur-hop-limit 64
                link-mtu 0
                managed-flag true
                max-interval 600
                name-server fe80::de9f:dbff:fe28:8f6a
                other-config-flag true
                prefix ::/60 {
                    autonomous-flag true
                    on-link-flag true
                    valid-lifetime 2592000
                }
                reachable-time 0
                retrans-timer 0
                send-advert true
            }
        }
        mtu 1500
        speed auto
        vif 2 {
            address 10.1.1.1/24
            description "HVAC VLAN 2"
            ipv6 {
                dup-addr-detect-transmits 1
            }
            mtu 1500
        }
        vif 3 {
            address 172.16.0.1/24
            description "ADT VLAN 3"
            ipv6 {
                dup-addr-detect-transmits 1
            }
            mtu 1500
        }
        vif 4 {
            address 10.10.10.1/24
            description "VOIP VLAN 4"
            ipv6 {
                dup-addr-detect-transmits 1
            }
        }
        vif 5 {
            address 192.168.10.1/24
            description "Internal Users VLAN 5"
            ip {
                enable-proxy-arp
            }
            ipv6 {
                dup-addr-detect-transmits 1
            }
            mtu 1500
        }
        vif 6 {
            address 192.168.4.1/24
            description "Guest VLAN 6"
            ip {
                enable-proxy-arp
            }
            ipv6 {
                dup-addr-detect-transmits 1
            }
            mtu 1500
        }
        vif 7 {
            address 172.16.4.1/24
            description "Test VLAN 7"
            ip {
                enable-proxy-arp
            }
            ipv6 {
                dup-addr-detect-transmits 1
            }
            mtu 1500
        }
        vif 555 {
            address 192.168.50.1/24
            description Media
            ip {
                enable-proxy-arp
            }
            ipv6 {
                dup-addr-detect-transmits 1
            }
            mtu 1500
        }
    }
    ethernet eth1 {
        address 192.168.42.1/24
        description VoIP
        duplex auto
        ipv6 {
            dup-addr-detect-transmits 1
            router-advert {
                cur-hop-limit 64
                link-mtu 0
                managed-flag true
                max-interval 600
                name-server fe80::de9f:dbff:fe28:8f6a
                other-config-flag true
                prefix ::/60 {
                    autonomous-flag true
                    on-link-flag true
                    valid-lifetime 2592000
                }
                reachable-time 0
                retrans-timer 0
                send-advert true
            }
        }
        mtu 1500
        speed auto
    }
    ethernet eth2 {
        address dhcp
        description "External WAN"
        dhcp-options {
            default-route update
            default-route-distance 210
            name-server no-update
        }
        dhcpv6-pd {
            no-dns
            pd 0 {
                interface eth0 {
                    host-address ::1
                    no-dns
                    prefix-id :1
                    service slaac
                }
                interface eth0.2 {
                    host-address ::1
                    no-dns
                    prefix-id :2
                    service slaac
                }
                interface eth0.3 {
                    host-address ::1
                    no-dns
                    prefix-id :3
                    service slaac
                }
                interface eth0.4 {
                    host-address ::1
                    no-dns
                    prefix-id :4
                    service slaac
                }
                interface eth0.5 {
                    host-address ::1
                    no-dns
                    prefix-id :5
                    service slaac
                }
                interface eth0.6 {
                    host-address ::1
                    no-dns
                    prefix-id :6
                    service slaac
                }
                interface eth0.555 {
                    host-address ::1
                    no-dns
                    prefix-id :7
                    service slaac
                }
                interface eth1 {
                    host-address ::1
                    prefix-id :8
                    service slaac
                }
                prefix-length 60
            }
            rapid-commit enable
        }
        duplex auto
        mtu 1500
        speed auto
    }
    loopback lo {
    }
    pseudo-ethernet peth0 {
        address 192.168.168.1/24
        description "Pixel Server"
        link eth0
    }
}
port-forward {
    auto-firewall enable
    hairpin-nat enable
    lan-interface eth0
    lan-interface eth1
    rule 100 {
        description "Unifi Device Access for Cloud Key Controller inbound UDP traffic"
        forward-to {
            address 192.168.100.3
        }
        original-port http-alt
        protocol tcp
    }
    rule 200 {
        description "Unifi Device Access for Cloud Key Controller inbound TCP traffic"
        forward-to {
            address 192.168.100.3
        }
        original-port 3478
        protocol udp
    }
    rule 300 {
        description "UNMS Server"
        forward-to {
            address 192.168.10.98
            port 443
        }
        original-port 2443
    }
    wan-interface eth2
}
protocols {
    igmp-proxy {
        interface eth0 {
            role downstream
            threshold 1
        }
        interface eth0.5 {
            role downstream
            threshold 1
        }
        interface eth0.6 {
            role downstream
            threshold 1
        }
        interface eth0.7 {
            role downstream
            threshold 1
        }
        interface eth0.555 {
            alt-subnet 0.0.0.0/0
            role upstream
            threshold 1
        }
    }
}
service {
    bcast-relay {
        id 1 {
            description PlayFi
            interface eth0
            interface eth0.6
            interface eth0.555
            interface eth0.5
            port 10102
        }
        id 2 {
            description Roku
            interface eth0
            interface eth0.6
            interface eth0.555
            interface eth0.5
            port 1900
        }
    }
    dhcp-server {
        disabled false
        hostfile-update disable
        shared-network-name ADT {
            authoritative enable
            subnet 172.16.0.0/24 {
                default-router 172.16.0.1
                dns-server 192.168.100.1
                domain-name ashcreek.home
                lease 86400
                start 172.16.0.2 {
                    stop 172.16.0.2
                }
            }
        }
        shared-network-name HVAC {
            authoritative enable
            subnet 10.1.1.0/24 {
                default-router 10.1.1.1
                dns-server 192.168.100.1
                domain-name ashcreek.home
                lease 86400
                start 10.1.1.2 {
                    stop 10.1.1.4
                }
            }
        }
        shared-network-name LAN0 {
            authoritative enable
            subnet 192.168.100.0/24 {
                default-router 192.168.100.1
                dns-server 192.168.100.1
                domain-name ashcreek.home
                lease 86400
                start 192.168.100.5 {
                    stop 192.168.100.110
                }
                static-mapping CloudKey {
                    ip-address 192.168.100.3
                    mac-address 80:2a:a8:cd:f5:1c
                }
                static-mapping Neils-McAyre-LAN1-lan0 {
                    ip-address 192.168.100.70
                    mac-address ac:87:a3:21:8e:94
                }
                static-mapping Neils-McAyre-LAN2-lan0 {
                    ip-address 192.168.100.71
                    mac-address 00:e0:4d:68:00:53
                }
                static-mapping Neils-McAyre-WiFi-lan0 {
                    ip-address 192.168.100.72
                    mac-address 34:36:3b:83:a1:36
                }
                static-mapping UKAP {
                    ip-address 192.168.100.10
                    mac-address 04:18:d6:20:0a:18
                }
                static-mapping cisco-switch {
                    ip-address 192.168.100.15
                    mac-address e0:d1:73:bc:80:70
                }
                static-mapping tough-switch {
                    ip-address 192.168.100.20
                    mac-address dc:9f:db:80:ba:24
                }
                static-mapping unms-server {
                    ip-address 192.168.100.2
                    mac-address 04:18:d6:f1:57:46
                }
            }
        }
        shared-network-name Media {
            authoritative enable
            subnet 192.168.50.0/24 {
                bootfile-name /config/user-data/tftp/emrk-0.9c.bin
                bootfile-server 192.168.50.1
                default-router 192.168.50.1
                dns-server 208.67.220.220
                dns-server 208.67.222.222
                lease 86400
                start 192.168.50.30 {
                    stop 192.168.50.60
                }
                static-mapping FamilyRoomAppleTV {
                    ip-address 192.168.50.56
                    mac-address b8:78:2e:3b:55:7a
                }
                static-mapping FireTv {
                    ip-address 192.168.50.55
                    mac-address a0:02:dc:fe:eb:c4
                }
                static-mapping GreatRoomAppleTV {
                    ip-address 192.168.50.54
                    mac-address 1c:1a:c0:9b:cf:73
                }
                static-mapping OfficeAppleTV {
                    ip-address 192.168.50.58
                    mac-address c8:69:cd:3d:f6:c4
                }
                static-mapping PlayFi {
                    ip-address 192.168.50.50
                    mac-address 60:b6:06:00:ff:b6
                }
                static-mapping PlayFi_LAN {
                    ip-address 192.168.50.51
                    mac-address b0:1f:81:30:07:ab
                }
                static-mapping Roku {
                    ip-address 192.168.50.57
                    mac-address ac:3a:7a:ac:7a:d0
                }
                static-mapping Samsung-DVD {
                    ip-address 192.168.50.59
                    mac-address 00:23:99:1b:7f:26
                }
                static-mapping TiVo {
                    ip-address 192.168.50.53
                    mac-address 00:11:d9:3c:0d:bd
                }
                static-mapping Yamaha-RX-V473 {
                    ip-address 192.168.50.52
                    mac-address 00:a0:de:8f:03:04
                }
                tftp-server-name Media
            }
        }
        shared-network-name Test {
            authoritative enable
            subnet 10.10.10.0/24 {
                default-router 10.10.10.1
                dns-server 192.168.100.1
                domain-name ashcreek.home
                lease 86400
                start 10.10.10.2 {
                    stop 10.10.10.3
                }
            }
        }
        shared-network-name VoIP {
            authoritative enable
            subnet 192.168.42.0/24 {
                default-router 192.168.42.1
                dns-server 208.67.220.220
                dns-server 208.67.222.222
                domain-name ashcreek.home
                lease 86400
                start 192.168.42.5 {
                    stop 192.168.42.10
                }
                static-mapping ObiTalk {
                    ip-address 192.168.42.5
                    mac-address 9c:ad:ef:21:71:cb
                }
            }
        }
        shared-network-name WLAN0 {
            authoritative enable
            subnet 192.168.10.0/24 {
                default-router 192.168.10.1
                dns-server 192.168.100.1
                domain-name ashcreek.home
                lease 86400
                start 192.168.10.5 {
                    stop 192.168.10.240
                }
                static-mapping Beadle-ThinkPad {
                    ip-address 192.168.10.90
                    mac-address 24:77:03:7d:77:dc
                }
                static-mapping Canon-Printer {
                    ip-address 192.168.10.30
                    mac-address 2c:9e:fc:4f:a9:09
                }
                static-mapping MacBook-Air-LAN {
                    ip-address 192.168.10.71
                    mac-address ac:87:a3:21:8e:94
                }
                static-mapping MacBook-Air-WiFi {
                    ip-address 192.168.10.70
                    mac-address 34:36:3b:83:a1:36
                }
                static-mapping WDMyCloud {
                    ip-address 192.168.10.35
                    mac-address 00:90:a9:d6:dd:2e
                }
            }
        }
        shared-network-name WLAN6 {
            authoritative enable
            subnet 192.168.4.0/24 {
                default-router 192.168.4.1
                dns-server 192.168.100.1
                domain-name ashcreek.home
                lease 86400
                start 192.168.4.5 {
                    stop 192.168.4.120
                }
                static-mapping Neils-McAyre-LAN-wkvlan6 {
                    ip-address 192.168.4.70
                    mac-address ac:87:a3:21:8e:94
                }
                static-mapping Neils-McAyre-WiFi-wkvlan6 {
                    ip-address 192.168.4.71
                    mac-address 34:36:3b:83:a1:36
                }
                static-mapping RSA-Laptop {
                    ip-address 192.168.4.80
                    mac-address 3c:97:0e:6d:d3:e1
                }
                static-mapping RingDoorBell-WiFi-wkvlan6 {
                    ip-address 192.168.4.10
                    mac-address 00:1d:c9:2d:b0:70
                }
            }
        }
        use-dnsmasq enable
    }
    dns {
        dynamic {
            interface eth2 {
                service afraid {
                    host-name appledram.strangled.net
                    host-name fishbourne.chickenkiller.com
                    login britannic
                    password ****************
                    server freedns.afraid.org
                }
                service custom-cloudflare {
                    host-name appledram.empirecreekcircle.com
                    host-name testsvr.empirecreekcircle.com
                    host-name test.empirecreekcircle.com
                    host-name fishbourne.empirecreekcircle.com
                    login cloudflare@empirecreekcircle.com
                    options zone=empirecreekcircle.com
                    password ****************
                    protocol cloudflare
                    server www.cloudflare.com
                }
            }
        }
        forwarding {
            blacklist {
                disabled false
                dns-redirect-ip 192.168.168.1
                domains {
                    exclude bing.com
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
                        prefix "zone "
                        url http://malc0de.com/bl/ZONES
                    }
                    source malwaredomains.com {
                        description "Just domains"
                        prefix ""
                        url http://mirror1.malwaredomains.com/files/justdomains
                    }
                    source simple_tracking {
                        description "Basic tracking list by Disconnect"
                        prefix ""
                        url https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt
                    }
                    source zeus {
                        description "abuse.ch ZeuS domain blocklist"
                        prefix ""
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
                exclude help.evernote.com
                exclude herokuapp.com
                exclude hp.com
                exclude hulu.com
                exclude images-amazon.com
                exclude live.com
                exclude microsoft.com
                exclude msdn.com
                exclude msecnd.net
                exclude msftncsi.com
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
                        prefix "0.0.0.0 "
                        url https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts
                    }
                    source sysctl.org {
                        description "This hosts file is a merged collection of hosts from cameleon"
                        prefix "127.0.0.1	 "
                        url http://sysctl.org/cameleon/hosts
                    }
                    source yoyo {
                        description "Fully Qualified Domain Names only - no prefix to strip"
                        prefix ""
                        url http://pgl.yoyo.org/as/serverlist.phphostformat=nohtml&showintro=1&mimetype=plaintext
                    }
                    source tasty {
                        description "File source"
                        dns-redirect-ip 10.10.10.10
                        file ./internal/testdata/blist.hosts.src
                    }
                }
            }
            cache-size 150
            listen-on eth0
            listen-on eth0.2
            listen-on eth0.3
            listen-on eth0.4
            listen-on eth0.5
            listen-on eth0.555
            listen-on eth0.6
            listen-on eth1
            listen-on l2tp0
            listen-on l2tp1
            listen-on l2tp2
            listen-on l2tp3
            listen-on l2tp4
            listen-on l2tp5
            listen-on l2tp6
            listen-on l2tp9
            listen-on lo
            name-server 208.67.220.220
            name-server 208.67.222.222
            name-server 2620:0:ccc::2
            name-server 2620:0:ccd::2
            options bogus-priv
            options domain-needed
            options domain=ashcreek.home
            options enable-ra
            options except-interface=eth2
            options expand-hosts
            options listen-address=::1
            options listen-address=127.0.0.1
            options localise-queries
            options strict-order
            system
        }
    }
    gui {
        https-port 443
        listen-address 192.168.100.1
        older-ciphers enable
    }
    mdns {
        repeater {
            interface eth0.5
            interface eth0.6
            interface eth0
            interface eth0.555
        }
    }
    nat {
        rule 5010 {
            description "NAT Rule Set"
            log disable
            outbound-interface eth2
            protocol all
            type masquerade
        }
    }
    ssh {
        disable-password-authentication
        port 22
        protocol-version v2
    }
    unms {
        connection wss://192.168.10.98:443+zgC6g70kUTIAK87uctPcbEULlqlQO5gzaoSVeJYpAUG4NfKq+allowSelfSignedCertificate
    }
}
system {
    config-management {
        commit-archive {
            location ftp://edgeos:abochceHa5#!@wdmycloud.ashcreek.home/EdgeOs
        }
        commit-revisions 50
    }
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
        modules {
            ftp {
                disable
            }
            gre {
                disable
            }
            h323 {
                disable
            }
            pptp {
                disable
            }
            sip {
                disable
                enable-indirect-media
                enable-indirect-signalling
                port 5060
            }
            tftp {
                disable
            }
        }
        table-size 32768
        tcp {
            half-open-connections 512
            loose enable
            max-retrans 3
        }
        timeout {
            icmp 10
            other 600
            tcp {
                close 10
                close-wait 60
                established 86400
                fin-wait 10
                last-ack 10
                syn-recv 5
                syn-sent 5
                time-wait 10
            }
            udp {
                other 10
                stream 180
            }
        }
    }
    domain-name ashcreek.home
    host-name appledram
    ip {
        override-hostname-ip 192.168.100.1
    }
    login {
        banner {
            post-login "\nWelcome to EdgeOS!\n"
            pre-login "\n\n\n\tWARNING *** WARNING *** WARNING *** WARNING *** WARNING\n\n\n\tWARNING: Criminal and civil penalties may be imposed for obtaining\n\tunauthorized access to this system or for causing intentional,\n\tunauthorized damage, deletion, alteration, or insertion of data.\n\tAny information stored, processed, or transmitted to this system\n\tmay be monitored, used, or disclosed by authorized personnel,\n\tincluding law enforcement. Email sysadmin@empirecreekcircle.com\n\tto gain access to this equipment if you need authorization.\n\n\n"
        }
        user nbnt {
            authentication {
                encrypted-password ****************
                plaintext-password ****************
                public-keys root@MacBook-Air.home {
                    key ****************
                    type ssh-rsa
                }
            }
            full-name Admin
            level admin
        }
        user sysop {
            authentication {
                encrypted-password ****************
                plaintext-password ****************
                public-keys sysop@ubnt.ashcreek.home {
                    key ****************
                    type ssh-rsa
                }
            }
            full-name "System Operator"
            home-directory /home/sysop
            level operator
        }
    }
    name-server 127.0.0.1
    ntp {
        server 0.ubnt.pool.ntp.org {
        }
        server 1.ubnt.pool.ntp.org {
        }
        server 2.ubnt.pool.ntp.org {
        }
        server 3.ubnt.pool.ntp.org {
        }
    }
    offload {
        hwnat disable
        ipsec enable
        ipv4 {
            forwarding enable
            gre enable
            pppoe disable
            vlan enable
        }
        ipv6 {
            forwarding enable
            pppoe disable
            vlan enable
        }
    }
    package {
        repository wheezy {
            components "main contrib non-free"
            distribution wheezy
            password ****************
            url http://http.us.debian.org/debian/
            username ""
        }
        repository wheezy-backports {
            components main
            distribution wheezy-backports
            password ****************
            url http://http.us.debian.org/debian
            username ""
        }
    }
    static-host-mapping {
        host-name ADT-Pulse.ashcreek.home {
            alias adt
            alias pulse
            inet 172.16.0.2
        }
        host-name Canon-Printer.ashcreek.home {
            alias laser
            alias mfc
            inet 192.168.10.30
        }
        host-name CloudKey.ashcreek.home {
            inet 192.168.100.3
        }
        host-name FamilyRoomAppleTV.ashcreek.home {
            alias FamilyRoomAppleTV
            inet 192.168.50.56
        }
        host-name GreatRoomAppleTV.ashcreek.home {
            alias GreatRoomAppleTV
            inet 192.168.50.54
        }
        host-name HVAC-TL-WR741ND.ashcreek.home {
            inet 10.1.1.2
        }
        host-name OBi202-VOIP.ashcreek.home {
            alias obi
            alias voip
            alias obitalk
            alias phone
            inet 192.168.42.5
        }
        host-name OfficeAppleTV.ashcreek.home {
            alias OfficeAppleTV
            inet 192.168.50.58
        }
        host-name PlayFi.ashcreek.home {
            alias PlayFi
            inet 192.168.50.50
            inet 192.168.50.51
        }
        host-name RingDoorBell.ashcreek.home {
            alias RingDoorBell
            inet 192.168.4.10
        }
        host-name Roku.ashcreek.home {
            alias Roku
            inet 192.168.50.57
        }
        host-name Samsung-DVD.ashcreek.home {
            alias samsung
            inet 192.168.50.59
        }
        host-name ThinkPad.ashcreek.home {
            alias thinkpad
            inet 192.168.10.90
        }
        host-name TiVo.ashcreek.home {
            inet 192.168.50.53
        }
        host-name UKAP.ashcreek.home {
            alias ap
            alias uap
            alias WiFi
            alias unifi
            inet 192.168.100.10
        }
        host-name WDMyCloud.ashcreek.home {
            inet 192.168.10.35
        }
        host-name Yamaha-RX-V473.ashcreek.home {
            alias receiver
            alias yahama
            inet 192.168.50.52
        }
        host-name appledram.ashcreek.home {
            alias appledram
            alias router
            alias ubnt
            inet 192.168.100.1
        }
        host-name cisco-switch.ashcreek.home {
            alias cisco
            inet 192.168.100.15
        }
        host-name hvac.ashcreek.home {
            alias ac
            inet 10.1.1.2
        }
        host-name lenovo-lt.ashcreek.home {
            alias lenovo-lt
            alias rsa-laptop
            alias rsa-lt
            inet 192.168.4.80
        }
        host-name mac-ayre.ashcreek.home {
            alias macair
            alias macbook
            alias mcair
            alias neils-mac
            inet 192.168.10.70
            inet 192.168.10.71
            inet 192.168.100.70
            inet 192.168.100.71
            inet 192.168.100.72
            inet 192.168.4.70
            inet 192.168.4.71
        }
        host-name scratch.ashcreek.home {
            alias demo-router
            alias scratch
            inet 192.168.42.10
        }
        host-name tough-switch.ashcreek.home {
            alias switch
            inet 192.168.100.20
        }
        host-name wdmycloud.ashcreek.home {
            alias backup
            alias cloud
            inet 192.168.10.35
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
                path /config/scripts/update-dnsmasq.pl
            }
            interval 1d
        }
        task update_bogons {
            executable {
                path /config/scripts/bogon.py
            }
            interval 1d30m
        }
    }
    time-zone America/Los_Angeles
    traffic-analysis {
        dpi enable
        export enable
    }
}
vpn {
    ipsec {
        auto-update 180
        auto-firewall-nat-exclude enable
        esp-group dpd_client {
            compression disable
            lifetime 3600
            mode transport
            pfs dh-group26
            proposal 1 {
                encryption aes256
                hash sha1
            }
        }
        ike-group dpd_client {
            dead-peer-detection {
                action clear
                interval 30
                timeout 120
            }
            ikev2-reauth no
            key-exchange ikev1
            lifetime 28800
            proposal 1 {
                dh-group 26
                encryption aes256
                hash sha256
            }
        }
        ipsec-interfaces {
            interface eth0
            interface eth2
        }
        nat-networks {
            allowed-network 0.0.0.0/0 {
            }
        }
        nat-traversal enable
    }
    l2tp {
        remote-access {
            authentication {
                local-users {
                    username tsd {
                        password ****************
                    }
                }
                mode local
            }
            client-ip-pool {
                start 192.168.4.240
                stop 192.168.4.254
            }
            dhcp-interface eth2
            dns-servers {
                server-1 192.168.100.1
            }
            ipsec-settings {
                authentication {
                    mode pre-shared-secret
                    pre-shared-secret ****************
                }
                ike-lifetime 3600
            }
            mtu 1492
        }
    }
}
zone-policy {
    zone dmz {
        default-action drop
        description "DMZ Zone"
        from ext {
            firewall {
                ipv6-name ipv6-ext-dmz
                name ext-dmz
            }
        }
        from gst {
            firewall {
                ipv6-name ipv6-gst-dmz
                name gst-dmz
            }
        }
        from int {
            firewall {
                ipv6-name ipv6-int-dmz
                name int-dmz
            }
        }
        from loc {
            firewall {
                ipv6-name ipv6-loc-dmz
                name loc-dmz
            }
        }
        from mdx {
            firewall {
                ipv6-name ipv6-mdx-dmz
                name mdx-dmz
            }
        }
        from pbx {
            firewall {
                ipv6-name ipv6-pbx-dmz
                name pbx-dmz
            }
        }
        interface eth0.2
        interface eth0.3
        interface eth0.4
    }
    zone ext {
        default-action drop
        description "External Zone"
        from dmz {
            firewall {
                ipv6-name ipv6-dmz-ext
                name dmz-ext
            }
        }
        from gst {
            firewall {
                ipv6-name ipv6-gst-ext
                name gst-ext
            }
        }
        from int {
            firewall {
                ipv6-name ipv6-int-ext
                name int-ext
            }
        }
        from loc {
            firewall {
                ipv6-name ipv6-loc-ext
                name loc-ext
            }
        }
        from mdx {
            firewall {
                ipv6-name ipv6-mdx-ext
                name mdx-ext
            }
        }
        from pbx {
            firewall {
                ipv6-name ipv6-pbx-ext
                name pbx-ext
            }
        }
        interface eth2
    }
    zone gst {
        default-action drop
        description "Guest Zone"
        from dmz {
            firewall {
                ipv6-name ipv6-dmz-gst
                name dmz-gst
            }
        }
        from ext {
            firewall {
                ipv6-name ipv6-ext-gst
                name ext-gst
            }
        }
        from int {
            firewall {
                ipv6-name ipv6-int-gst
                name int-gst
            }
        }
        from loc {
            firewall {
                ipv6-name ipv6-loc-gst
                name loc-gst
            }
        }
        from mdx {
            firewall {
                ipv6-name ipv6-mdx-gst
                name mdx-gst
            }
        }
        from pbx {
            firewall {
                ipv6-name ipv6-pbx-gst
                name pbx-gst
            }
        }
        interface eth0.6
        interface eth0.7
        interface l2tp0
        interface l2tp1
        interface l2tp2
        interface l2tp3
        interface l2tp4
        interface l2tp5
        interface l2tp6
        interface l2tp7
        interface l2tp8
        interface l2tp9
    }
    zone int {
        default-action drop
        description "Internal Zone"
        from dmz {
            firewall {
                ipv6-name ipv6-dmz-int
                name dmz-int
            }
        }
        from ext {
            firewall {
                ipv6-name ipv6-ext-int
                name ext-int
            }
        }
        from gst {
            firewall {
                ipv6-name ipv6-gst-int
                name gst-int
            }
        }
        from loc {
            firewall {
                ipv6-name ipv6-loc-int
                name loc-int
            }
        }
        from mdx {
            firewall {
                ipv6-name ipv6-mdx-int
                name mdx-int
            }
        }
        from pbx {
            firewall {
                ipv6-name ipv6-pbx-int
                name pbx-int
            }
        }
        interface eth0
        interface eth0.5
        interface peth0
    }
    zone loc {
        default-action drop
        from dmz {
            firewall {
                ipv6-name ipv6-dmz-loc
                name dmz-loc
            }
        }
        from ext {
            firewall {
                ipv6-name ipv6-ext-loc
                name ext-loc
            }
        }
        from gst {
            firewall {
                ipv6-name ipv6-gst-loc
                name gst-loc
            }
        }
        from int {
            firewall {
                ipv6-name ipv6-int-loc
                name int-loc
            }
        }
        from mdx {
            firewall {
                ipv6-name ipv6-mdx-loc
                name mdx-loc
            }
        }
        from pbx {
            firewall {
                ipv6-name ipv6-pbx-loc
                name pbx-loc
            }
        }
        local-zone
    }
    zone mdx {
        default-action drop
        description "Media Zone"
        from dmz {
            firewall {
                ipv6-name ipv6-dmz-mdx
                name dmz-mdx
            }
        }
        from ext {
            firewall {
                ipv6-name ipv6-ext-mdx
                name ext-mdx
            }
        }
        from gst {
            firewall {
                ipv6-name ipv6-gst-mdx
                name gst-mdx
            }
        }
        from int {
            firewall {
                ipv6-name ipv6-int-mdx
                name int-mdx
            }
        }
        from loc {
            firewall {
                ipv6-name ipv6-loc-mdx
                name loc-mdx
            }
        }
        from pbx {
            firewall {
                ipv6-name ipv6-pbx-mdx
                name pbx-mdx
            }
        }
        interface eth0.555
    }
    zone pbx {
        default-action drop
        description "VoIP Zone"
        from dmz {
            firewall {
                ipv6-name ipv6-dmz-pbx
                name dmz-pbx
            }
        }
        from ext {
            firewall {
                ipv6-name ipv6-ext-pbx
                name ext-pbx
            }
        }
        from gst {
            firewall {
                ipv6-name ipv6-gst-pbx
                name gst-pbx
            }
        }
        from int {
            firewall {
                ipv6-name ipv6-int-pbx
                name int-pbx
            }
        }
        from loc {
            firewall {
                ipv6-name ipv6-loc-pbx
                name loc-pbx
            }
        }
        from mdx {
            firewall {
                ipv6-name ipv6-mdx-pbx
                name mdx-pbx
            }
        }
        interface eth1
    }
}


/* Warning: Do not remove the following line. */
/* === vyatta-config-version: "config-management@1:conntrack@1:cron@1:dhcp-relay@1:dhcp-server@4:firewall@5:ipsec@5:nat@3:qos@1:quagga@2:system@4:ubnt-pptp@1:ubnt-unms@1:ubnt-util@1:vrrp@1:webgui@1:webproxy@1:zone-policy@1" === */
/* Release version: v1.9.7.5001798.170720.0132 */

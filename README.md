# **UBNT edgeos-dnsmasq-blacklist dnsmasq DNS Blacklisting and Redirection**

[![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/britannic/blacklist/blob/master/LICENSE.txt "View the software license here")[![Version](https://img.shields.io/badge/version-v1.2.3-green.svg)](https://github.com/britannic/blacklist "Latest version")[![GoDoc](https://godoc.org/github.com/britannic/blacklist?status.svg)](https://godoc.org/github.com/britannic/blacklist  "Go documentation")[![Build Status](https://travis-ci.org/britannic/blacklist.svg?branch=master)](https://travis-ci.org/britannic/blacklist  "Build status for this version")[![Coverage Status](https://coveralls.io/repos/github/britannic/blacklist/badge.svg?branch=master "")](https://coveralls.io/github/britannic/blacklist?branch=master "Test coverage status for this version")[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/github.com/britannic/blacklist "Quality of Go code for this version")

Follow the conversation @ [community.ubnt.com](https://community.ubnt.com/t5/EdgeRouter/DNS-Adblocking-amp-Blacklisting-dnsmasq-Configuration/td-p/2215008/jump-to/first-unread-message "Follow the conversation about this software in the EdgeRouter forum (https://community.ubnt.com/t5/EdgeRouter/)")

## Donations and Sponsorship

Please show your thanks by donating to the project using [Square Cash](https://cash.me/$HelmRockSecurity/ "Securely send and receive cash without fees using Square Cash") or [PayPal](https://www.paypal.me/helmrocksecurity/)

[![Donate](https://img.shields.io/badge/Donate-%245-orange.svg?style=plastic)](https://cash.me/$HelmRockSecurity/5 "Give $5 using Square Cash (free money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2410-red.svg?style=plastic)](https://cash.me/$HelmRockSecurity/10 "Give $10 using Square Cash (free money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2415-yellow.svg?style=plastic)](https://cash.me/$HelmRockSecurity/15 "Give $15 using Square Cash (free money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2420-yellowgreen.svg?style=plastic)](https://cash.me/$HelmRockSecurity/20 "Give $20 using Square Cash (free money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2425-brightgreen.svg?style=plastic)](https://cash.me/$HelmRockSecurity/25 "Give $25 using Square Cash (free money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2450-ff69b4.svg?style=plastic)](https://cash.me/$HelmRockSecurity/50 "Give $50 using Square Cash (free money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%24100-blue.svg?style=plastic)](https://cash.me/$HelmRockSecurity/100 "Give $100 using Square Cash (free money transfer)")
[![Donate](https://img.shields.io/badge/Donate-Custom%20Amount-4B0082.svg?style=plastic)](https://cash.me/$HelmRockSecurity/ "Choose your own donation amount using Square Cash (free money transfer)")

[![Donate](https://img.shields.io/badge/Donate-%245-orange.svg?style=plastic)](https://paypal.me/helmrocksecurity/5 "Give $5 using PayPal (PayPal money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2410-red.svg?style=plastic)](https://paypal.me/helmrocksecurity/10 "Give $10 using PayPal (PayPal money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2415-yellow.svg?style=plastic)](https://paypal.me/helmrocksecurity/15 "Give $15 using PayPal (PayPal money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2420-yellowgreen.svg?style=plastic)](https://paypal.me/helmrocksecurity/20 "Give $20 using PayPal (PayPal money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2425-brightgreen.svg?style=plastic)](https://paypal.me/helmrocksecurity/25 "Give $25 using PayPal (PayPal money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%2450-ff69b4.svg?style=plastic)](https://paypal.me/helmrocksecurity/50 "Give $50 using PayPal (PayPal money transfer)")
[![Donate](https://img.shields.io/badge/Donate-%24100-blue.svg?style=plastic)](https://paypal.me/helmrocksecurity/100 "Give $100 using PayPal (PayPal money transfer)")
[![Donate](https://img.shields.io/badge/Donate-Custom%20Amount-4B0082.svg?style=plastic)](https://paypal.me/helmrocksecurity/ "Choose your own donation amount using PayPal (PayPal money transfer)")

We greatly appreciate any and all donations - Thank you! Funds go to maintaining development servers and networks.

## Note: This is 3rd party software and isn't supported or endorsed by Ubiquiti Networks®

## **Contents**

1. [Overview](#overview)
1. [Donate](#donations-and-sponsorship)
1. [Copyright](#copyright)
1. [Licenses](#licenses)
1. [Latest Version](#latest-version)
1. [Change Log](https://github.com/britannic/blacklist/blob/master/CHANGELOG.md)
1. [Features](#features)
1. [Compatibility](#compatibility)
1. [Installation](#installation)
    1. [Using apt-get](#apt-get-installation---erlite-3-erpoe-5-er-x-er-x-sfp-er4-unifi-gateway-3--unifi-gateway-4)
    1. [Using dpkg](#dpkg-installation---best-for-disk-space-constrained-routers)
1. [Upgrade](#upgrade)
1. [Reconfigure](#reconfigure)
1. [Removal](#removal)
1. [Frequently Asked Questions](#frequently-asked-questions)
   1. [Can I donate to project?](#donations-and-sponsorship)
   1. [Does the install backup my blacklist configuration before deleting it?](#does-the-install-backup-my-blacklist-configuration-before-deleting-it)
   1. [Does update-dnsmasq run automatically?](#does-update-dnsmasq-run-automatically)
   1. [How do I add or delete sources?](#how-do-i-add-or-delete-sources)
   1. [How do I back up my blacklist configuration and restore it later?](#how-do-i-back-up-my-blacklist-configuration-and-restore-it-later)
   1. [How do I configure dnsmasq?](#how-do-i-configure-dnsmasq)
   1. [How do I configure local file sources instead of internet based ones?](#how-do-i-configure-local-file-sources-instead-of-internet-based-ones)
   1. [How do I disable/enable dnsmasq blacklisting?](#how-do-i-disableenable-dnsmasq-blacklisting)
   1. [How do I exclude or include a host or a domain?](#how-do-i-exclude-or-include-a-host-or-a-domain)
   1. [How do I globally exclude or include hosts or a domains?](#how-do-i-globally-exclude-or-include-hosts-or-a-domains)
   1. [How do I use the command line switches?](#how-do-i-use-the-command-line-switches)
   1. [How do can keep my USG configuration after an upgrade, provision or reboot?](#how-do-i-keep-my-usg-configuration-after-an-upgrade-provision-or-reboot)
   1. [How does whitelisting work?](#how-does-whitelisting-work)
   1. [What is the difference between blocking domains and hosts?](#what-is-the-difference-between-blocking-domains-and-hosts)
   1. [Which blacklist sources are installed by default?](#which-blacklist-sources-are-installed-by-default)

## **Overview**

EdgeMax dnsmasq DNS blacklisting and redirection is inspired by the users at [EdgeMAX Community](https://community.ubnt.com/t5/EdgeMAX/bd-p/EdgeMAX/)

[[Top]](#contents)

## **Copyright**

* Copyright © 2020 [Helm Rock Consulting](https://www.helmrock.com/ "Visit Helm Rock Consulting at https://www.helmrock.com/")

[[Top]](#contents)

## **Licenses**

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
1. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

    THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
    ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
    WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
    DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
    ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
    (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
    LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
    ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
    (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
    SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

    The views and conclusions contained in the software and documentation are those
    of the authors and should not be interpreted as representing official policies,
    either expressed or implied, of the FreeBSD Project.

[[Top]](#contents)

## **Latest Version**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;[![Latest](https://img.shields.io/badge/Release-v1.2.3-green.svg)](https://github.com/britannic/blacklist/releases/latest  "Latest version")

### Release v1.2.2 (May 25, 2020)

* Add source OSID Mobile Blocklist
* Update documentation
* Check for internet connectivity before running an update
* Check and load /config/user-data/blacklist.cfg in preference to config.boot
* Fix bug in configure preventing use of the latest commit

[[Top]](#contents)

## **Change Log**

* See [changelog](https://github.com/britannic/blacklist/blob/master/CHANGELOG.md) for details.

[[Top]](#contents)

## **Features**

* Adds DNS blacklisting integration to the EdgeRouter configuration
* Generates configuration files used directly by dnsmasq to redirect dns lookups
* Integrated with the EdgeMax OS CLI
* Any FQDN in the blacklist will force dnsmasq to return the configured dns redirect IP address

[[Top]](#contents)

## **Compatibility**

* edgeos-dnsmasq-blacklist has been tested on the EdgeRouter ERLite-3, ERPoe-5, ER-X, ER4 UniFi Security Gateway USG3 and USG4 routers
  * EdgeMAX versions: v1.9.7+hotfix.4-v2.0.8, UniFi: v4.4.12-v4.4.44.5213871

[[Top]](#contents)

## **Installation**

* [Using apt-get](#apt-get-installation---erlite-3-erpoe-5-er-x-er-x-sfp-er4-unifi-gateway-3--unifi-gateway-4) - works for all routers
* [Using dpkg](#dpkg-installation---best-for-disk-space-constrained-routers) - best for disk space constrained routers

[[Top]](#contents)

### **apt-get Installation - ERLite-3, ERPoe-5, ER-X, ER-X-SFP, ER4, UniFi-Gateway-3 & UniFi-Gateway-4**

* Add the blacklist debian package repository using the router's CLI shell

```bash
configure
set system package repository blacklist components main
set system package repository blacklist description 'Britannic blacklist debian stretch repository'
set system package repository blacklist distribution stretch
set system package repository blacklist url 'https://raw.githubusercontent.com/britannic/debian-repo/master/blacklist/'
commit;save;exit
```

* Add the GPG signing key

```bash
sudo curl -L https://raw.githubusercontent.com/britannic/debian-repo/master/blacklist/public.key | sudo apt-key add -
```

* Update the system repositorities and install edgeos-dnsmasq-blacklist

```bash
sudo apt-get update && sudo apt-get install edgeos-dnsmasq-blacklist
```

[[Top]](#contents)

## **dpkg Installation - best for disk space constrained routers**

### **EdgeRouter ERLite-3, ERPoe-5, ER4, UniFi-Gateway-3 & UniFi-Gateway-4**

```bash
curl -L -O https://raw.githubusercontent.com/britannic/blacklist/master/edgeos-dnsmasq-blacklist_1.2.3_mips.deb
sudo dpkg -i edgeos-dnsmasq-blacklist_1.2.3_mips.deb
```

[[Top]](#contents)

### **EdgeRouter ER-X & ER-X-SFP**

* Ensure the router has enough space, by removing unnecessary files

```bash
sudo apt-get clean cache
delete system image
```

* Now download and install the edgeos-dnsmasq-blacklist package

```bash
curl -L -O https://raw.githubusercontent.com/britannic/blacklist/master/edgeos-dnsmasq-blacklist_1.2.3_mipsel.deb
sudo dpkg -i edgeos-dnsmasq-blacklist_1.2.3_mipsel.deb
```

[[Top]](#contents)

## **Upgrade**

* If the repository is set up and you are using apt-get:

```bash
sudo apt-get update && sudo apt-get install --only-upgrade edgeos-dnsmasq-blacklist
```

* Note, if you are using dpkg, it cannot upgrade packages, so follow these [instructions](#dpkg-installation---best-for-disk-space-constrained-routers) and the previous package version will be automatically removed before the new package version is installed

[[Top]](#contents)

## **Reconfigure**

* If the Unifi Security Gateway has been re-provisioned you might need to re-enable the blacklists, in order to do so run:

```bash
sudo dpkg-reconfigure edgeos-dnsmasq-blacklist
```

[[Top]](#contents)

## **Removal**

### **EdgeMAX - All Platforms**

```bash
sudo apt-get remove --purge edgeos-dnsmasq-blacklist
```

[[Top]](#contents)

## **Frequently Asked Questions**

### **How do I disable/enable dnsmasq blacklisting?**

* Use these CLI configure commands:
* Disable:

```bash
configure
set service dns forwarding blacklist disabled true
commit;save;exit
```

* Enable:

```bash
configure
set service dns forwarding blacklist disabled false
commit;save;exit
```

[[Top]](#contents)

### **Does the install backup my blacklist configuration before deleting it?**

* If a blacklist configuration already exists, the install routine will automatically back it up to /config/user-data/blacklist.$(date +'%FT%H%M%S').cmds

[[Top]](#contents)

### **How do I back up my blacklist configuration and restore it later?**

* use the following commands (make a note of the file name) in the shell (not in configure):

```bash
export DATE=$(date +'%FT%H%M%S'); echo "Backing up blacklist configuration to: /config/user-data/blacklist.${DATE}.cmds"; show configuration commands | grep blacklist > /config/user-data/blacklist.$(date +'%FT%H%M%S').cmds
```

* After installing the latest version, you can merge your backed up configuration:

```bash
configure
.  /config/user-data/blacklist.[date string].cmds
commit;save;exit
```

* If you prefer to delete the default configuration and restore your previous configuration, run these commands:

```bash
configure
delete service dns forwarding blacklist
.  /config/user-data/blacklist.[date string].cmds
commit;save;exit
```

[[Top]](#contents)

### **Which blacklist sources are installed by default?**

* Use these CLI shell commands to view the current sources or scan the log for previous downloads:

```bash
show configuration commands | match source
grep downloaded /var/log/update-dnsmasq.log
```

[[Top]](#contents)

### **How do I configure local file sources instead of internet based ones?**

* Use these commands to configure a local file source

```bash
set service dns forwarding blacklist hosts source myhosts description 'Blacklist file source'
set service dns forwarding blacklist hosts source myhosts dns-redirect-ip 0.0.0.0
set service dns forwarding blacklist hosts source myhosts file /config/user-data/blist.hosts.src
```

* File contents example for /config/user-data/blist.hosts.src:

```bash
gsmtop.net
click.buzzcity.net
ads.admoda.com
stats.pflexads.com
a.glcdn.co
wwww.adleads.com
ad.madvertise.de
apps.buzzcity.net
ads.mobgold.com
android.bcfads.com
req.appads.com
show.buzzcity.net
api.analytics.omgpop.com
r.edge.inmobicdn.net
www.mmnetwork.mobi
img.ads.huntmad.com
creative1cdn.mobfox.com
admicro2.vcmedia.vn
admicro1.vcmedia.vn
```

[[Top]](#contents)

### **How do I keep my USG configuration after an upgrade, provision or reboot?**

* Follow these [instructions](https://britannic.github.io/install-edgeos-packages/) on how to automatically install edgeos-dnsmasq-blacklist
* Generate and download a config.gateway.json file from your USG following these [instructions](https://help.ubnt.com/hc/en-us/articles/215458888-UniFi-How-to-further-customize-USG-configuration-with-config-gateway-json)
* Here's a sample [config.gateway.json](https://raw.githubusercontent.com/britannic/blacklist/master/config.gateway.json)
* Once the config.gateway.json has been generated, it will need to be uploaded to your **UniFi controller** per the [instructions](https://help.ubnt.com/hc/en-us/articles/215458888-UniFi-How-to-further-customize-USG-configuration-with-config-gateway-json)

[[Top]](#contents)

### **How do I add or delete sources?**

* Using the CLI configure command, to delete domains and hosts sources:

```bash
configure
delete service dns forwarding blacklist domains source malc0de
delete service dns forwarding blacklist hosts source yoyo.org
commit;save;exit
```

* To add a source, first check it can serve a text list and also note the prefix (if any) before the hosts or domains, e.g. [http://www.malwaredomainlist.com/](http://www.malwaredomainlist.com/) has this format:

```text
#               MalwareDomainList.com Hosts List           #
#   http://www.malwaredomainlist.com/hostslist/hosts.txt   #
#         Last updated: Mon, 04 Dec 17 19:18:42 +0000      #


127.0.0.1  localhost
127.0.0.1  0koryu0.easter.ne.jp
127.0.0.1  109-204-26-16.netconnexion.managedbroadband.co.uk
127.0.0.1  1866809.securefastserver.com
```

* So the prefix is "127.0.0.1  "
* Here's how to creating the source in the CLI:

```bash
configure
set service dns forwarding blacklist hosts source malwaredomainlist description '127.0.0.1 based host and domain list'
set service dns forwarding blacklist hosts source malwaredomainlist prefix '127.0.0.1  '
set service dns forwarding blacklist hosts source malwaredomainlist url 'http://www.malwaredomainlist.com/hostslist/hosts.txt'
commit;save;exit
```

[[Top]](#contents)

### **How do I globally exclude or include hosts or a domains?**

* Use these example commands to globally include or exclude blacklisted entries:

```bash
configure
set service dns forwarding blacklist exclude cdn.visiblemeasures.com
set service dns forwarding blacklist include www.nastywebsites.com
commit;save;exit
```

[[Top]](#contents)

### **How do I exclude or include a host or a domain?**

* Use these example commands to include or exclude blacklisted entries:

```bash
configure
set service dns forwarding blacklist domains exclude visiblemeasures.com
set service dns forwarding blacklist domains include domainsnastywebsites.com
set service dns forwarding blacklist hosts exclude cdn.visiblemeasures.com
set service dns forwarding blacklist hosts include www.nastywebsites.com
commit;save;exit
```

[[Top]](#contents)

### **How does whitelisting work?**

*dnsmasq will whitelist any entries in the configuration file domains and hosts (servers) with a hash in place of an IP address (the "#" force dnsmasq to forward the DNS request to the router's configured nameservers)

* i.e. servers (hosts)

```bash
server=/www.bing.com/#
```

* i.e. domains

```bash
address=/bing.com/#
```

[[Top]](#contents)

### **Does update-dnsmasq run automatically?**

* Yes, a scheduled task is created and run daily at midnight with a random start delay is used ensure other routers in the same time zone won't overload the source servers.
* The random start delay window is configured in seconds using this command - this example sets the start delay between 1-10800 seconds (0-3 hours):

```bash
set system task-scheduler task update_blacklists executable arguments 10800
```

* It can be reconfigured using these CLI configuration commands:

```bash
set system task-scheduler task update_blacklists executable path /config/scripts/update-dnsmasq-cronjob.sh
set system task-scheduler task update_blacklists executable arguments 10800
set system task-scheduler task update_blacklists interval 1d
```

* For example, to change the execution interval to every 6 hours, use this command:

```bash
set system task-scheduler task update_blacklists interval 6h
```

* In daily use, no additional interaction with update-dnsmasq is required. By default, cron will run update-dnsmasq at midnight each day to download the blacklist sources and update the dnsmasq configuration files in /etc/dnsmasq.d. dnsmasq will automatically be reloaded after the configuration file update is completed.

[[Top]](#contents)

### **How do I use the command line switches?**

* update-dnsmasq has the following commandline switches available:

```bash
/config/scripts/update-dnsmasq -h
    -dir string
            Override dnsmasq directory (default "/etc/dnsmasq.d")
    -f [full file path]
            [full file path] # Load a config.boot file
    -h   Display help
    -v   Verbose display
    -version
            Show version
```

[[Top]](#contents)

### **How do I configure dnsmasq?**

* dnsmasq may need to be configured to ensure blacklisting works correctly
  * Here is an example using the EdgeOS configuration shell

```bash
configure
set service dns forwarding cache-size 2048
set service dns forwarding except-interface [Your WAN i/f]
set service dns forwarding name-server [Your choice of IPv4 Internet Name-Server]
set service dns forwarding name-server [Your choice of IPv4 Internet Name-Server]
set service dns forwarding name-server [Your choice of IPv6 Internet Name-Server]
set service dns forwarding name-server [Your choice of IPv6 Internet Name-Server]
set service dns forwarding options bogus-priv
set service dns forwarding options domain-needed
set service dns forwarding options domain=mydomain.local
set service dns forwarding options enable-ra
set service dns forwarding options expand-hosts
set service dns forwarding options localise-queries
set service dns forwarding options strict-order
set service dns forwarding system
set system name-server 127.0.0.1
set system name-server '::1'
commit; save; exit
```

[[Top]](#contents)

### **What is the difference between blocking domains and hosts?**

* The difference lies in the order of update-dnsmasq's processing algorithm. Domains are processed first and take precedence over hosts, so that a blacklisted domain will force update-dnsmasq's source parser to exclude subsequent hosts from the same domain. This reduces dnsmasq's list of lookups, since it will automatically redirect hosts for a blacklisted domain.

[[Top]](#contents)

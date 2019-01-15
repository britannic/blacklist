// Package main is the UBNT edgeos-dnsmasq-blacklist dnsmasq DNS Blacklisting and Redirection.
//
// View the software license here (https://github.com/britannic/blacklist/blob/master/LICENSE.txt)Latest versionVersion (https://github.com/britannic/blacklist)Go documentationGoDoc (https://godoc.org/github.com/britannic/blacklist)Build status for this versionBuild Status (https://travis-ci.org/britannic/blacklist)Test coverage status for this versionCoverage Status (https://coveralls.io/github/britannic/blacklist?branch=master)Quality of Go code for this versionGo Report Card (https://goreportcard.com/report/github.com/britannic/blacklist)
//
// Follow the conversation @ community.ubnt.com (https://community.ubnt.com/t5/EdgeRouter/DNS-Adblocking-amp-Blacklisting-dnsmasq-Configuration/td-p/2215008/jump-to/first-unread-message "Follow the conversation about this software in the EdgeRouter forum (https://community.ubnt.com/t5/EdgeRouter/)")
//
// Donations and Sponsorship
//
// Please show your thanks by donating to the project using Securely send and receive cash without fees using Square CashSquare Cash (https://cash.me/$HelmRockSecurity/) or PayPal (https://www.paypal.me/helmrocksecurity/)
//
// Donate (https://cash.me/$HelmRockSecurity/5 "Give $5 using Square Cash (free money transfer)")
// Donate (https://cash.me/$HelmRockSecurity/10 "Give $10 using Square Cash (free money transfer)")
// Donate (https://cash.me/$HelmRockSecurity/15 "Give $15 using Square Cash (free money transfer)")
// Donate (https://cash.me/$HelmRockSecurity/20 "Give $20 using Square Cash (free money transfer)")
// Donate (https://cash.me/$HelmRockSecurity/25 "Give $25 using Square Cash (free money transfer)")
// Donate (https://cash.me/$HelmRockSecurity/50 "Give $50 using Square Cash (free money transfer)")
// Donate (https://cash.me/$HelmRockSecurity/100 "Give $100 using Square Cash (free money transfer)")
// Donate (https://cash.me/$HelmRockSecurity/ "Choose your own donation amount using Square Cash (free money transfer)")
//
// Donate (https://paypal.me/helmrocksecurity/5 "Give $5 using PayPal (PayPal money transfer)")
// Donate (https://paypal.me/helmrocksecurity/10 "Give $10 using PayPal (PayPal money transfer)")
// Donate (https://paypal.me/helmrocksecurity/15 "Give $15 using PayPal (PayPal money transfer)")
// Donate (https://paypal.me/helmrocksecurity/20 "Give $20 using PayPal (PayPal money transfer)")
// Donate (https://paypal.me/helmrocksecurity/25 "Give $25 using PayPal (PayPal money transfer)")
// Donate (https://paypal.me/helmrocksecurity/50 "Give $50 using PayPal (PayPal money transfer)")
// Donate (https://paypal.me/helmrocksecurity/100 "Give $100 using PayPal (PayPal money transfer)")
// Donate (https://paypal.me/helmrocksecurity/ "Choose your own donation amount using PayPal (PayPal money transfer)")
//
// We greatly appreciate any and all donations - Thank you! Funds go to maintaining development servers and networks.
//
// Note: This is 3rd party software and isn't supported or endorsed by Ubiquiti Networks®
//
// Contents
//
// • Overview (#overview)
//
// • Donate (#donations-and-sponsorship)
//
// • Copyright (#copyright)
//
// • Licenses (#licenses)
//
// • Latest Version (#latest-version)
//
// • Change Log (https://github.com/britannic/blacklist/blob/master/CHANGELOG.md)
//
// • Features (#features)
//
// • Compatibility (#compatibility)
//
// • Installation (#installation)
//
// • Using apt-get (#apt-get-installation---erlite-3-erpoe-5-er-x-er-x-sfp--unifi-gateway-3)
//
// • Using dpkg (#dpkg-installation---best-for-disk-space-constrained-routers)
//
// • Upgrade (#upgrade)
//
// • Removal (#removal)
//
// • Frequently Asked Questions (#frequently-asked-questions)
//
// • Can I donate to project? (#donations-and-sponsorship)
//
// • Does the install backup my blacklist configuration before deleting it? (#does-the-install-backup-my-blacklist-configuration-before-deleting-it)
//
// • Does update-dnsmasq run automatically? (#does-update-dnsmasq-run-automatically)
//
// • How do I add or delete sources? (#how-do-i-add-or-delete-sources)
//
// • How do I back up my blacklist configuration and restore it later? (#how-do-i-back-up-my-blacklist-configuration-and-restore-it-later)
//
// • How do I configure dnsmasq? (#how-do-i-configure-dnsmasq)
//
// • How do I configure local file sources instead of internet based ones? (#how-do-i-configure-local-file-sources-instead-of-internet-based-ones)
//
// • How do I disable/enable dnsmasq blacklisting? (#how-do-i-disableenable-dnsmasq-blacklisting)
//
// • How do I exclude or include a host or a domain? (#how-do-i-exclude-or-include-a-host-or-a-domain)
//
// • How do I globally exclude or include hosts or a domains? (#how-do-i-globally-exclude-or-include-hosts-or-a-domains)
//
// • How do I use the command line switches? (#how-do-i-use-the-command-line-switches)
//
// • How do can keep my USG configuration after an upgrade, provision or reboot? (#how-do-can-keep-my-usg-configuration-after-an-upgrade-provision-or-reboot)
//
// • How does whitelisting work? (#how-does-whitelisting-work)
//
// • What is the difference between blocking domains and hosts? (#what-is-the-difference-between-blocking-domains-and-hosts)
//
// • Which blacklist sources are installed by default? (#which-blacklist-sources-are-installed-by-default)
//
// Overview
//
// EdgeMax dnsmasq DNS blacklisting and redirection is inspired by the users at EdgeMAX Community (https://community.ubnt.com/t5/EdgeMAX/bd-p/EdgeMAX/)
//
// [Top] (#contents)
//
// Copyright
//
// • Copyright © Visit Helm Rock Consulting at https://www.helmrock.com/2019 Helm Rock Consulting (https://www.helmrock.com/)
//
// [Top] (#contents)
//
// Licenses
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//
// • Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.
//
//
// • Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
//
// The views and conclusions contained in the software and documentation are those
// of the authors and should not be interpreted as representing official policies,
// either expressed or implied, of the FreeBSD Project.
//
//
// [Top] (#contents)
//
// Latest Version
//
// &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Latest versionLatest (https://github.com/britannic/blacklist/releases/latest)
//
// Release v1.1.6.2 (April 24, 2018)
//
// • Code refactor
//
// • Global whitelist and blacklist configuration files now have their own prefix: "roots" i.e.
//
//   roots.global-blacklisted-domains.blacklist.conf
//
// [Top] (#contents)
//
// Change Log
//
// • See changelog (https://github.com/britannic/blacklist/blob/master/CHANGELOG.md) for details.
//
// [Top] (#contents)
//
// Features
//
// • Adds DNS blacklisting integration to the EdgeRouter configuration
//
// • Generates configuration files used directly by dnsmasq to redirect dns lookups
//
// • Integrated with the EdgeMax OS CLI
//
// • Any FQDN in the blacklist will force dnsmasq to return the configured dns redirect IP address
//
// [Top] (#contents)
//
// Compatibility
//
// • edgeos-dnsmasq-blacklist has been tested on the EdgeRouter ERLite-3, ERPoe-5, ER-X and UniFi Security Gateway USG-3 routers
//
// • EdgeMAX versions: v1.9.7+hotfix.4-v1.10.1, UniFi: v4.4.12-v4.4.18
//
// • integration could be adapted to work on VyOS and Vyatta derived ports, since  EdgeOS is a fork and port of Vyatta 6.3
//
// [Top] (#contents)
//
// Installation
//
// • Using apt-get (#apt-get-installation---erlite-3-erpoe-5-er-x-er-x-sfp--unifi-gateway-3) - works for all routers
//
// • Using dpkg (#dpkg-installation---best-for-disk-space-constrained-routers) - best for disk space constrained routers
//
// [Top] (#contents)
//
// apt-get Installation - ERLite-3, ERPoe-5, ER-X, ER-X-SFP & UniFi-Gateway-3
//
// • Add the blacklist debian package repository using the router's CLI shell
//
//   configure
//   set system package repository blacklist components main
//   set system package repository blacklist description 'Britannic blacklist debian wheezy repository'
//   set system package repository blacklist distribution wheezy
//   set system package repository blacklist url 'https://raw.githubusercontent.com/britannic/debian-repo/master/blacklist/'
//   commit;save;exit
//
// • Add the GPG signing key
//
//   sudo curl -L https://raw.githubusercontent.com/britannic/debian-repo/master/blacklist/public.key | sudo apt-key add -
//
// • Update the system repositorities and install edgeos-dnsmasq-blacklist
//
//   sudo apt-get update && sudo apt-get install edgeos-dnsmasq-blacklist
//
// [Top] (#contents)
//
// dpkg Installation - best for disk space constrained routers
//
// EdgeRouter ERLite-3, ERPoe-5 & UniFi-Gateway-3
//
//   curl -L -O https://raw.githubusercontent.com/britannic/blacklist/master/edgeos-dnsmasq-blacklist_1.1.6.2_mips.deb
//   sudo dpkg -i edgeos-dnsmasq-blacklist_1.1.6.2_mips.deb
//
// [Top] (#contents)
//
// EdgeRouter ER-X & ER-X-SFP
//
// • Ensure the router has enough space, by removing unnecessary files
//
//   sudo apt-get clean cache
//   delete system image
//
// • Now download and install the edgeos-dnsmasq-blacklist package
//
//   curl -L -O https://raw.githubusercontent.com/britannic/blacklist/master/edgeos-dnsmasq-blacklist_1.1.6.2_mipsel.deb
//   sudo dpkg -i edgeos-dnsmasq-blacklist_1.1.6.2_mipsel.deb
//
// [Top] (#contents)
//
// Upgrade
//
// • If the repository is set up and you are using apt-get:
//
//   sudo apt-get update && sudo apt-get upgrade edgeos-dnsmasq-blacklist
//
// • Note, if you are using dpkg, it cannot upgrade packages, so follow these instructions (#dpkg-installation---best-for-disk-space-constrained-routers) and the previous package version will be automatically removed before the new package version is installed
//
// [Top] (#contents)
//
// Removal
//
// EdgeMAX - All Platforms
//
//   sudo apt-get remove --purge edgeos-dnsmasq-blacklist
//
// [Top] (#contents)
//
// Frequently Asked Questions
//
// How do I disable/enable dnsmasq blacklisting?
//
// • Use these CLI configure commands:
//
// • Disable:
//
//   configure
//   set service dns forwarding blacklist disabled true
//   commit;save;exit
//
// • Enable:
//
//   configure
//   set service dns forwarding blacklist disabled false
//   commit;save;exit
//
// [Top] (#contents)
//
// Does the install backup my blacklist configuration before deleting it?
//
// • If a blacklist configuration already exists, the install routine will automatically back it up to /config/user-data/blacklist.$(date +'%FT%H%M%S').cmds
//
// [Top] (#contents)
//
// How do I back up my blacklist configuration and restore it later?
//
// • use the following commands (make a note of the file name):
//
//   export DATE=$(date +'%FT%H%M%S'); echo "Backing up blacklist configuration to: /config/user-data/blacklist.${DATE}.cmds"; show configuration commands | grep blacklist > /config/user-data/blacklist.$(date +'%FT%H%M%S').cmds
//
// • After installing the latest version, you can merge your backed up configuration:
//
//   configure
//   .  /config/user-data/blacklist.[date string].cmds
//   commit;save;exit
//
// • If you prefer to delete the default configuration and restore your previous configuration, run these commands:
//
//   configure
//   delete service dns forwarding blacklist
//   .  /config/user-data/blacklist.[date string].cmds
//   commit;save;exit
//
// [Top] (#contents)
//
// Which blacklist sources are installed by default?
//
// • You can use this command in the CLI shell to view the current sources after installation or view the log and see previous downloads:
//
//   show configuration commands | match blacklist | match source
//   more /var/log/update-dnsmasq.log
//
// [Top] (#contents)
//
// How do I configure local file sources instead of internet based ones?
//
// • Use these commands to configure a local file source
//
//   set service dns forwarding blacklist hosts source myhosts description 'Blacklist file source'
//   set service dns forwarding blacklist hosts source myhosts dns-redirect-ip 0.0.0.0
//   set service dns forwarding blacklist hosts source myhosts file /config/user-data/blist.hosts.src
//
// • File contents example for /config/user-data/blist.hosts.src:
//
//   gsmtop.net
//   click.buzzcity.net
//   ads.admoda.com
//   stats.pflexads.com
//   a.glcdn.co
//   wwww.adleads.com
//   ad.madvertise.de
//   apps.buzzcity.net
//   ads.mobgold.com
//   android.bcfads.com
//   req.appads.com
//   show.buzzcity.net
//   api.analytics.omgpop.com
//   r.edge.inmobicdn.net
//   www.mmnetwork.mobi
//   img.ads.huntmad.com
//   creative1cdn.mobfox.com
//   admicro2.vcmedia.vn
//   admicro1.vcmedia.vn
//
// [Top] (#contents)
//
// How do can keep my USG configuration after an upgrade, provision or reboot?
//
// • Follow these instructions (https://britannic.github.io/install-edgeos-packages/) on how to automatically install edgeos-dnsmasq-blacklist
//
// • Create a config.gateway.json file following these instructions (https://help.ubnt.com/hc/en-us/articles/215458888-UniFi-How-to-further-customize-USG-configuration-with-config-gateway-json)
//
// • Here's a sample config.gateway.json (https://raw.githubusercontent.com/britannic/blacklist/master/config.gateway.json)
//
// [Top] (#contents)
//
// How do I add or delete sources?
//
// • Using the CLI configure command, to delete domains and hosts sources:
//
//   configure
//   delete service dns forwarding blacklist domains source malc0de
//   delete service dns forwarding blacklist hosts source yoyo.org
//   commit;save;exit
//
// • To add a source, first check it can serve a text list and also note the prefix (if any) before the hosts or domains, e.g. http://www.malwaredomainlist.com/ (http://www.malwaredomainlist.com/) has this format:
//
//   #               MalwareDomainList.com Hosts List           #
//   #   http://www.malwaredomainlist.com/hostslist/hosts.txt   #
//   #         Last updated: Mon, 04 Dec 17 19:18:42 +0000      #
//
//
//   127.0.0.1  localhost
//   127.0.0.1  0koryu0.easter.ne.jp
//   127.0.0.1  109-204-26-16.netconnexion.managedbroadband.co.uk
//   127.0.0.1  1866809.securefastserver.com
//
// • So the prefix is "127.0.0.1  "
//
// • Here's how to creating the source in the CLI:
//
//   configure
//   set service dns forwarding blacklist hosts source malwaredomainlist description '127.0.0.1 based host and domain list'
//   set service dns forwarding blacklist hosts source malwaredomainlist prefix '127.0.0.1  '
//   set service dns forwarding blacklist hosts source malwaredomainlist url 'http://www.malwaredomainlist.com/hostslist/hosts.txt'
//   commit;save;exit
//
// [Top] (#contents)
//
// How do I globally exclude or include hosts or a domains?
//
// • Use these example commands to globally include or exclude blacklisted entries:
//
//   configure
//   set service dns forwarding blacklist exclude cdn.visiblemeasures.com
//   set service dns forwarding blacklist include www.nastywebsites.com
//   commit;save;exit
//
// [Top] (#contents)
//
// How do I exclude or include a host or a domain?
//
// • Use these example commands to include or exclude blacklisted entries:
//
//   configure
//   set service dns forwarding blacklist domains exclude visiblemeasures.com
//   set service dns forwarding blacklist domains include domainsnastywebsites.com
//   set service dns forwarding blacklist hosts exclude cdn.visiblemeasures.com
//   set service dns forwarding blacklist hosts include www.nastywebsites.com
//   commit;save;exit
//
// [Top] (#contents)
//
// How does whitelisting work?
//
// *dnsmasq will whitelist any entries in the configuration file domains and hosts (servers) with a hash in place of an IP address (the "#" force dnsmasq to forward the DNS request to the router's configured nameservers)
//
// • i.e. servers (hosts)
//
//   server=/www.bing.com/#
//
// • i.e. domains
//
//   address=/bing.com/#
//
// [Top] (#contents)
//
// Does update-dnsmasq run automatically?
//
// • Yes, a scheduled task is created and run daily at midnight with a random start delay is used ensure other routers in the same time zone won't overload the source servers.
//
// • The random start delay window is configured in seconds using this command - this example sets the start delay between 1-10800 seconds (0-3 hours):
//
//   set system task-scheduler task update_blacklists executable arguments 10800
//
// • It can be reconfigured using these CLI configuration commands:
//
//   set system task-scheduler task update_blacklists executable path /config/scripts/blacklist-cronjob.sh
//   set system task-scheduler task update_blacklists executable arguments 10800
//   set system task-scheduler task update_blacklists interval 1d
//
// • For example, to change the execution interval to every 6 hours, use this command:
//
//   set system task-scheduler task update_blacklists interval 6h
//
// • In daily use, no additional interaction with update-dnsmasq is required. By default, cron will run update-dnsmasq at midnight each day to download the blacklist sources and update the dnsmasq configuration files in /etc/dnsmasq.d. dnsmasq will automatically be reloaded after the configuration file update is completed.
//
// [Top] (#contents)
//
// How do I use the command line switches?
//
// • update-dnsmasq has the following commandline switches available:
//
//   /config/scripts/update-dnsmasq -h
//       -dir string
//               Override dnsmasq directory (default "/etc/dnsmasq.d")
//       -f [full file path]
//               [full file path] # Load a config.boot file
//       -h   Display help
//       -v   Verbose display
//       -version
//               Show version
//
// [Top] (#contents)
//
// How do I configure dnsmasq?
//
// • dnsmasq may need to be configured to ensure blacklisting works correctly
//
// • Here is an example using the EdgeOS configuration shell
//
//   configure
//   set service dns forwarding cache-size 2048
//   set service dns forwarding except-interface [Your WAN i/f]
//   set service dns forwarding name-server [Your choice of IPv4 Internet Name-Server]
//   set service dns forwarding name-server [Your choice of IPv4 Internet Name-Server]
//   set service dns forwarding name-server [Your choice of IPv6 Internet Name-Server]
//   set service dns forwarding name-server [Your choice of IPv6 Internet Name-Server]
//   set service dns forwarding options bogus-priv
//   set service dns forwarding options domain-needed
//   set service dns forwarding options domain=mydomain.local
//   set service dns forwarding options enable-ra
//   set service dns forwarding options expand-hosts
//   set service dns forwarding options localise-queries
//   set service dns forwarding options strict-order
//   set service dns forwarding system
//   set system name-server 127.0.0.1
//   set system name-server '::1'
//   commit; save; exit
//
// [Top] (#contents)
//
// What is the difference between blocking domains and hosts?
//
// • The difference lies in the order of update-dnsmasq's processing algorithm. Domains are processed first and take precedence over hosts, so that a blacklisted domain will force update-dnsmasq's source parser to exclude subsequent hosts from the same domain. This reduces dnsmasq's list of lookups, since it will automatically redirect hosts for a blacklisted domain.
//
// [Top] (#contents)
//
// blacklist
//
//
package main

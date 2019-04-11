# Changelog

## Releases

### Release v1.1.7.4 (April 11, 2019)

* Added support for EdgeOS 2.0.1
* Fixed config session detection bug
* Ensure all dnsmasq blacklist configuration files are removed for blacklist config delete and package removal
* Force dnsmasq restart after blacklist config delete

### Release v1.1.6.11 (March 9,2019)

* Fix for issue #8, /etc/dnsmasq.d/*blacklist.conf files aren't deleted after uninstallation

### Release v1.1.6.10 (February 3,2019)

* Removed "try set service dns forwarding blacklist disabled false" from the installation setup to prevent false positive issue in EdgeOS 2.0
* Forced update_dnsmasq to run after installation, to resolve non-detection of a new configuration in EdgeOS 2.0

### Release v1.1.6.9 (January 27, 2019)

* Added adjust.com to domain exclusions as this domain is being used by Microsoft for guiding users through a two-factor authentication setup process.

### Release v1.1.6.8 (January 20, 2019)

* Removed "set service dns forwarding blacklist disabled false" as it is broken in EdgeOS 2.0.0

### Release v1.1.6.7 (January 13, 2019)

* Removed hostfile.org as it is too agressive and causes a lot of false positive entries

### Release v1.1.6.3 (June 11, 2018)

* Removed domain source Malc0de as it is posting false positives with no means to notify the maintainer
* Added log message for sources that have no records extracted

```bash
[Source]: no records processed - check source and/or configuration
```

* Any source without records extracted, will still have a file written, but it will contain a message

```bash
# NO DATA WRITTEN - CHECK WHITELIST EXCLUSIONS
```

### Release v1.1.6.2 (April 24, 2018)

* Code refactor
* Global whitelist and blacklist configuration files now have their own prefix: "roots" i.e.

```bash
roots.global-blacklisted-domains.blacklist.conf
```

### Release v1.1.6.1 (April 13, 2018)

* Suppress log messages for predefined includes (blacklists) and excludes (whitelists)
* Changed Debian pre-remove and post-install script algorithm to detect an existing installation and reinstall a default configuration if none exists

### Release v1.1.5 (April 10, 2018)

* Add build architecture and OS information to "-version" argument
* Code refactor and parsing algorithm updates

### Release v1.1.4 (April 8, 2018)

* Performance enhancements for source entry processing
* Improved data counter metrics for found, extracted and dropped records

### Release v1.1.3 (April 5, 2018)

* Document updates

### Release v1.1.2 (April 5, 2018)

* Fixed bug that overwrote user configured blacklist settings during an upgrade

### Release v1.1.1 (April 4, 2018)

* Fixed bug that prevented pre-configured included hosts being correctly blacklisted

### Release v1.1.0 (April 3, 2018)

* Fixed minor log message bug to insert space between progname and "starting up"
* Renamed blacklist-cronjob.sh to update-dnsmasq-cronjob.sh to be consistent with update-dnsmasq
* Fixed a bug when scripted configure session isn't detected, resulting in wrong showconfig mode being used
* Fixed a bug in blacklist-cronjob.sh that inhibited the cron job delay
* Changed http error handling from fatal to error notification, so that update-dnsmasq can continue processing for sources that don't have problems and complete the update
* Added code to support dnsmasq configuration file whitelisting for domains and hosts (servers) using hash syntax (the "#" force dnsmasq to forward the DNS request to the configured nameservers)
* i.e. servers (hosts)

```bash
server=/www.bing.com/#
```

* i.e. domains

```bash
address=/bing.com/#
```

### Release v1.0.10 (February 27, 2018)

* Added functions to ensure all blacklist configuration files are removed from /etc/dnsmasq.d/ when uninstalling using

```bash
dpkg -P edgeos-dnsmasq-blacklist
```

* Or

```bash
apt-get remove --purge edgeos-dnsmasq-blacklist
```

* dnsmasq will be automatically restarted to remove stale redirects

### Release v1.0.9 (February 26, 2018)

* Added logic to not run the post installation script after an upgrade

### Release v1.0.8 (February 26, 2018)

* Algorithm to trap out of range cronjob arguments

### Release v1.0.7 (February 25, 2018)

* Adjust task-scheduler argument

### Release v1.0.6 (February 25, 2018)

* Fix bug in task-scheduler stanza to insert missing key work "system

### Release v1.0.5 (February 25, 2018)

* Nightly update-dnsmasq cron job now has a configurable argument to set how many seconds of random delay before starting

```bash
set system task-scheduler task update_blacklists executable arguments 60
set system task-scheduler task update_blacklists executable path /config/scripts/blacklist-cronjob.sh
set system task-scheduler task update_blacklists interval 1d
```

### Release v1.0.4 (February 24, 2018)

* Implemented starting nightly update-dnsmasq cron job at random times to prevent a datastorm if a lot of users are in the same time zone
  * Suggested by EdgeMax Community User [@sorvani](https://community.ubnt.com/t5/user/viewprofilepage/user-id/185589)

### Release v1.0.3 (February 23, 2018)

* Switching to debian respository installation using apt-get

### Release v1.0.2 (February 18, 2018)

* Updated sources
  * Changed hosts source [https://github.com/StevenBlack/hosts/](https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts) to version that only blacklists ads and malware

### Release v1.0.1 (February 18, 2018)

* Updated sources
  * Changed hosts source [https://github.com/StevenBlack/hosts/](https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews-gambling-porn/hosts) to version that doesn't blacklist social networks

### GA Release v1.0.0 (February 17, 2018)

* Updated sources
  * Added hosts source [https://github.com/StevenBlack/hosts/](https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews-gambling-porn/hosts)
  * Removed Yoyo and raw.github.com, since the new sources make them redundant
* Fixed bug where "domains" was reported as "hosts" source type

### Release Candidate v1.0.0.rc5 (February 4, 2018)

* Extended vyattacfg change group to include /opt/vyatta/config/tmp and /opt/vyatta/config/active

### Release Candidate v1.0.0.rc4 (February 1, 2018)

* Release candidate #4 v1.0.0
* Fixed bug to ensure pre-configured includes are processed first, so that pre-configured excludes won't drop them

### Release Candidate v1.0.0.rc3 (January 31, 2018)

* Release candidate #3 v1.0.0
* Added domains exclude nsatc.net as it was blocking MS Office logins

### Release Candidate v1.0.0.rc2 (January 31, 2018)

* Release candidate #2 v1.0.0
* New source added: [http://www.hostsfile.org/Downloads/hosts.txt](http://www.hostsfile.org/Downloads/hosts.txt)
* Added global exclude googleads.g.doubleclick.net to fix Google search results

### Release Candidate v1.0.0.rc1 (January 31, 2018)

* Release candidate #1 v1.0.0
* Includes pre-remove.sh back up routine

### Patch v0.0.12 (January 30, 2018)

* Added global exclude "evernote.com" as it is being false flagged by some sources

### Patch v0.0.11 (January 30, 2018)

* Reformatted update-dnsmasq.log output
* README Updates
  * Added FAQ
  * Refactored layout

### Patch v0.0.10 (January 28, 2018)

* Improved counters for statistics logging
* Increased test coverage of code
* Additional documentation
* Added config.boot file loader

### Patch v0.0.9 (January 24, 2018)

* Added logging for download errors and warnings for empty content
* Change HTTP user agent to emulate curl, to stop web servers from offering complex content
* Removed embedded tabs in source prefixes that were interpreted by the EdgeOS configure shell as a completion request,  preventing correct prefix matches

### Patch v0.0.8 (January 22, 2018)

* Removes redundant references to blacklist.t and perl modules
* Replace "▶" with ":" in log messages

### Release v0.0.7 (January 22, 2018)

* Debian package release for ease of installation, maintenance and updating. See README for instructions and general release notes.

### Pre-release v0.0.6 (January 20, 2018)

* Debian package script rough in

### Pre-release v0.0.5-alpha (June 5, 2016)

* Major code refactor

### Pre-release v0.0.4-alpha (May 11, 2016)

* Ground up rewrite to create self contained packages and simplify code base

### Pre-release v0.0.3-alpha (Jan 16, 2016)

* Alpha code release

### Pre-Alpha (Jan 15, 2016)

* Learning Go, rudimentary coding
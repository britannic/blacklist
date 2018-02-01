# Changelog

## Releases

### Release Candidate v1.0.0.rc4 (February 1, 2018)

* Fixed bug to ensure pre-configured includes are processed first, so that pre-configured excludes won't drop them

### Release Candidate v1.0.0.rc3 (January 31, 2018)

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
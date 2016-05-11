
# edgeos
    import "github.com/britannic/edgeos"

© 2016 NJ Software. All rights reserved. Use of this source code is governed by a BSD-style license that can be found in the LICENSE.txt file.

[UBNT EdgeMax](<a href="https://community.ubnt.com/t5/EdgeMAX/bd-p/EdgeMAX">https://community.ubnt.com/t5/EdgeMAX/bd-p/EdgeMAX</a>) dnsmasq Blacklist and Adware Blocking

NOTE: THIS IS NOT OFFICIAL UBIQUITI SOFTWARE AND THEREFORE NOT SUPPORTED OR ENDORSED BY Ubiquiti Networks®

[![License](<a href="https://img.shields.io/badge/license-BSD-blue.svg">https://img.shields.io/badge/license-BSD-blue.svg</a>)](<a href="https://github.com/britannic/blacklist/blob/master/LICENSE.txt">https://github.com/britannic/blacklist/blob/master/LICENSE.txt</a>) [![Alpha  Version](<a href="https://img.shields.io/badge/version-v0.03--alpha-red.svg">https://img.shields.io/badge/version-v0.03--alpha-red.svg</a>)](<a href="https://github.com/britannic/blacklist">https://github.com/britannic/blacklist</a>) [![GoDoc](<a href="https://godoc.org/github.com/britannic/blacklist?status.svg">https://godoc.org/github.com/britannic/blacklist?status.svg</a>)](<a href="https://godoc.org/github.com/britannic/blacklist">https://godoc.org/github.com/britannic/blacklist</a>) [![Build Status](<a href="https://travis-ci.org/britannic/blacklist.svg?branch=master">https://travis-ci.org/britannic/blacklist.svg?branch=master</a>)](<a href="https://travis-ci.org/britannic/blacklist">https://travis-ci.org/britannic/blacklist</a>) [![Build Status](<a href="https://drone.io/github.com/britannic/blacklist/status.png">https://drone.io/github.com/britannic/blacklist/status.png</a>)](<a href="https://drone.io/github.com/britannic/blacklist/latest">https://drone.io/github.com/britannic/blacklist/latest</a>) [![Coverage Status](<a href="https://coveralls.io/repos/github/britannic/blacklist/badge.svg?branch=master">https://coveralls.io/repos/github/britannic/blacklist/badge.svg?branch=master</a>)](<a href="https://coveralls.io/github/britannic/blacklist?branch=master">https://coveralls.io/github/britannic/blacklist?branch=master</a>) [![Go Report Card](<a href="https://goreportcard.com/badge/gojp/goreportcard">https://goreportcard.com/badge/gojp/goreportcard</a>)](<a href="https://goreportcard.com/report/github.com/britannic/blacklist">https://goreportcard.com/report/github.com/britannic/blacklist</a>)

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

Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.




## Constants
``` go
const (

    // Fext sets the dnsmasq configuration file extension
    Fext = "blacklist.conf"

    // Domains sets the domains string
    Domains = "domains"

    // Hosts sets the hosts string
    Hosts = "hosts"

    // PreCon sets the string for pre-configured
    PreCon = "pre-configured"

    // False is a string constant
    False = "false"

    // True is a string constant
    True = "true"
)
```

## Variables
``` go
var (
    // API sets the path and executable for the EdgeOS shell API
    API = "/bin/cli-shell-api"
)
```

## func APICmd
``` go
func APICmd() (r map[string]string)
```
APICmd returns a map of CLI commands


## func DeleteFile
``` go
func DeleteFile(f string) bool
```
DeleteFile removes a file if it exists


## func DiffArray
``` go
func DiffArray(a, b []string) (diff []string)
```
DiffArray returns the delta of two arrays


## func GetHTTP
``` go
func GetHTTP(method, URL string) (io.Reader, error)
```
GetHTTP creates http requests to download data


## func GetType
``` go
func GetType(in interface{}) (out interface{})
```
GetType returns the converted "in" type


## func Insession
``` go
func Insession() bool
```
Insession returns true if VyOS/EdgeOS configuration is in session


## func ListFiles
``` go
func ListFiles(dir string) (files []string, err error)
```
ListFiles returns a list of blacklist files


## func Load
``` go
func Load(action, level string) (reader io.Reader, err error)
```
Load reads the config using the EdgeOS/VyOS cli-shell-api


## func PurgeFiles
``` go
func PurgeFiles(files []string) (err error)
```
PurgeFiles removes any orphaned blacklist files that don't have sources


## func SHCmd
``` go
func SHCmd(a string) string
```
SHCmd returns the appropriate command for non-tty or tty configure context


## func ToBool
``` go
func ToBool(s string) bool
```
ToBool converts a string ("true" or "false") to it's boolean equivalent


## func WriteFile
``` go
func WriteFile(fname string, data io.Reader) (err error)
```
WriteFile writes blacklist data to storage



## type Config
``` go
type Config map[string]*EdgeOS
```
Config is a map of EdgeOS











### func (Config) Files
``` go
func (c Config) Files(dir string, nodes []string) (files []string)
```
Files returns a list of dnsmasq conf files from all srcs



### func (Config) FormatData
``` go
func (c Config) FormatData(fmttr string, data []string) (reader io.Reader, list List)
```
FormatData returns a io.Reader loaded with dnsmasq formatted data



### func (Config) Get
``` go
func (c Config) Get(path string) (e *EdgeOS)
```
Get returns a normalized EdgeOS data set



### func (Config) GetExcludes
``` go
func (c Config) GetExcludes(dex, ex List, nodes []string) (List, List)
```
GetExcludes collates the configured excludes and merges the ex/dex lists



### func (Config) WriteIncludes
``` go
func (c Config) WriteIncludes(dir string, nodes []string) (dex, ex List)
```
WriteIncludes writes pre-configure data to disk



## type Configure
``` go
type Configure interface {
    Get(path string) (e *EdgeOS)
    Files() []string
    FormatData(node string, data []string) (reader io.Reader, list List, err error)
}
```
Configure has methods for returning config data supersets











## type Content
``` go
type Content struct {
    // contains filtered or unexported fields
}
```
Content struct holds content data









### func NewContent
``` go
func NewContent(toRead interface{}) *Content
```
NewContent returns a new Content pointer




### func (\*Content) Read
``` go
func (c *Content) Read(p []byte) (n int, err error)
```


## type Data
``` go
type Data interface {
    Reader
    Writer
}
```
Data inteface implements the Reader and Writer interfaces











## type EdgeOS
``` go
type EdgeOS struct {
    Disabled bool
    Exc      []string
    Inc      []string
    IP       string
    Nodes    []Leaf
    Sources  []Srcs
}
```
EdgeOS struct for normalizing EdgeOS data.











## type Keys
``` go
type Keys []string
```
Keys is used for sorting operations on map keys











### func (Keys) Len
``` go
func (k Keys) Len() int
```
Len returns length of Keys



### func (Keys) Less
``` go
func (k Keys) Less(i, j int) bool
```
Less returns the smallest element



### func (Keys) Swap
``` go
func (k Keys) Swap(i, j int)
```
Swap swaps elements of a key array



## type Leaf
``` go
type Leaf struct {
    Data     map[string]*Srcs `json:"data, omitempty"`
    Disabled bool             `json:"disable"`
    Excludes []string         `json:"excludes, omitempty"`
    Includes []string         `json:"includes, omitempty"`
    IP       string           `json:"ip, omitempty"`
}
```
Leaf is a struct for EdgeOS configuration data











## type List
``` go
type List map[string]int
```
List is a map of int









### func GetSubdomains
``` go
func GetSubdomains(s string) (l List)
```
GetSubdomains returns a map of subdomains


### func MergeList
``` go
func MergeList(a, b List) List
```
MergeList combines two List maps




### func (List) KeyExists
``` go
func (l List) KeyExists(s string) bool
```
KeyExists returns true if the key exists



### func (List) String
``` go
func (l List) String() string
```


### func (List) SubKeyExists
``` go
func (l List) SubKeyExists(s string) bool
```
SubKeyExists returns true if part of all of the key matches



## type Lister
``` go
type Lister interface {
    KeyExists(s string) bool
    String() string
    SubKeyExists(s string) bool
}
```
Lister implements List methods











## type Nodes
``` go
type Nodes map[string]*Leaf
```
Nodes is a map of Leaf nodes









### func NewNodes
``` go
func NewNodes() Nodes
```
NewNodes implements a new Node map


### func ReadCfg
``` go
func ReadCfg(reader io.Reader) (Nodes, error)
```
ReadCfg extracts nodes from a EdgeOS/VyOS configuration structure




### func (Nodes) JSON
``` go
func (n Nodes) JSON() string
```
JSON returns raw print for the Blacklist struct



### func (Nodes) NewConfig
``` go
func (n Nodes) NewConfig() (c Config)
```
NewConfig returns an initialized Config map of struct EdgeOS



### func (Nodes) SortKeys
``` go
func (n Nodes) SortKeys() (pkeys Keys)
```
SortKeys returns an array of sorted strings



### func (Nodes) SortSKeys
``` go
func (n Nodes) SortSKeys(node string) (skeys Keys)
```
SortSKeys returns an array of sorted strings



### func (Nodes) String
``` go
func (n Nodes) String() string
```
String returns pretty print for the Blacklist struct



## type Reader
``` go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```
Reader implements the Read []byte method











## type Srcs
``` go
type Srcs struct {
    Desc     string `json:"desc, omitempty"`
    Disabled bool   `json:"disabled, omitempty"`
    File     string `json:"file, omitempty"`
    IP       string `json:"ip, omitempty"`
    List     List   `json:"-"`
    Name     string `json:"name"`
    No       int    `json:"-"`
    Prefix   string `json:"prefix"`
    Type     int    `json:"type, omitempty"`
    URL      string `json:"url, omitempty"`
}
```
Srcs holds download source information









### func Process
``` go
func Process(s *Srcs, dex, ex List, reader io.Reader) *Srcs
```
Process extracts hosts/domains from downloaded raw content




## type Writer
``` go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```
Writer implements the Write []byte method

















- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
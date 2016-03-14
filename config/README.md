
# config
    import "./config"

Package config provides methods and data structures for loading
an EdgeOS/VyOS configuration





## Variables
``` go
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
```
Testdata contains a valid full EdgeOS blacklist configuration

``` go
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
```
Testdata2 contains a valid partial EdgeOS blacklist configuration


## func APICmd
``` go
func APICmd() (r map[string]string)
```
APICmd returns a map of CLI commands


## func Insession
``` go
func Insession() bool
```
Insession returns true if VyOS/EdgeOS configuration is in session


## func Load
``` go
func Load(action string, level string) (r string, err error)
```
Load reads the config using the EdgeOS/VyOS cli-shell-api


## func ToBool
``` go
func ToBool(s string) (b bool)
```
ToBool converts a string ("true" or "false") to it's boolean equivalent



## type Blacklist
``` go
type Blacklist map[string]*Node
```
Blacklist type is a map of Nodes with string keys









### func Get
``` go
func Get(cfg string, root string) (b *Blacklist, err error)
```
Get extracts nodes from a EdgeOS/VyOS configuration structure




### func (Blacklist) String
``` go
func (b Blacklist) String() (result string)
```
String returns pretty print for the Blacklist struct



## type Dict
``` go
type Dict map[string]int
```
Dict is a common string key map of ints









### func GetSubdomains
``` go
func GetSubdomains(s string) (d Dict)
```
GetSubdomains returns a map of subdomains




### func (Dict) KeyExists
``` go
func (d Dict) KeyExists(s string) bool
```
KeyExists returns true if the key exists



### func (Dict) SubKeyExists
``` go
func (d Dict) SubKeyExists(s string) bool
```
SubKeyExists returns true if part of the key matches



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



## type Node
``` go
type Node struct {
    Disable          bool
    IP               string
    Exclude, Include []string
    Source           Source
}
```
Node configuration record











## type Source
``` go
type Source map[string]*Src
```
Source is a map of Srcs with string keys











## type Src
``` go
type Src struct {
    Desc    string
    Disable bool
    IP      string
    List    Dict
    Name    string
    No      int
    Prfx    string
    Type    string
    URL     string
}
```
Src record struct for Source map

















- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
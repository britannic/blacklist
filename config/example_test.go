package config_test

import (
	"encoding/json"
	"fmt"
	"os"

	c "github.com/britannic/blacklist/config"
)

func ExampleCfg() {
	b, e := c.Get(c.Testdata, "blacklist")
	if e != nil {
		os.Exit(1)
	}

	j, _ := json.MarshalIndent(b, "", "   ")
	fmt.Println(string(j))

	// Output: {
	//    "blacklist": {
	//       "Disable": false,
	//       "IP": "0.0.0.0",
	//       "Exclude": [
	//          "122.2o7.net",
	//          "1e100.net",
	//          "adobedtm.com",
	//          "akamai.net",
	//          "amazon.com",
	//          "amazonaws.com",
	//          "apple.com",
	//          "ask.com",
	//          "avast.com",
	//          "bitdefender.com",
	//          "cdn.visiblemeasures.com",
	//          "cloudfront.net",
	//          "coremetrics.com",
	//          "edgesuite.net",
	//          "freedns.afraid.org",
	//          "github.com",
	//          "githubusercontent.com",
	//          "google.com",
	//          "googleadservices.com",
	//          "googleapis.com",
	//          "googleusercontent.com",
	//          "gstatic.com",
	//          "gvt1.com",
	//          "gvt1.net",
	//          "hb.disney.go.com",
	//          "hp.com",
	//          "hulu.com",
	//          "images-amazon.com",
	//          "msdn.com",
	//          "paypal.com",
	//          "rackcdn.com",
	//          "schema.org",
	//          "skype.com",
	//          "smacargo.com",
	//          "sourceforge.net",
	//          "ssl-on9.com",
	//          "ssl-on9.net",
	//          "static.chartbeat.com",
	//          "storage.googleapis.com",
	//          "windows.net",
	//          "yimg.com",
	//          "ytimg.com"
	//       ],
	//       "Include": null,
	//       "Source": {}
	//    },
	//    "domains": {
	//       "Disable": false,
	//       "IP": "",
	//       "Exclude": null,
	//       "Include": [
	//          "adsrvr.org",
	//          "adtechus.net",
	//          "advertising.com",
	//          "centade.com",
	//          "doubleclick.net",
	//          "free-counter.co.uk",
	//          "intellitxt.com",
	//          "kiosked.com"
	//       ],
	//       "Source": {
	//          "malc0de": {
	//             "Desc": "",
	//             "Disable": false,
	//             "IP": "",
	//             "List": null,
	//             "Name": "malc0de",
	//             "Prfx": "",
	//             "Type": "domains",
	//             "URL": "http://malc0de.com/bl/ZONES"
	//          }
	//       }
	//    },
	//    "hosts": {
	//       "Disable": false,
	//       "IP": "",
	//       "Exclude": null,
	//       "Include": [
	//          "beap.gemini.yahoo.com"
	//       ],
	//       "Source": {
	//          "adaway": {
	//             "Desc": "",
	//             "Disable": false,
	//             "IP": "",
	//             "List": null,
	//             "Name": "adaway",
	//             "Prfx": "",
	//             "Type": "hosts",
	//             "URL": "http://adaway.org/hosts.txt"
	//          },
	//          "malwaredomainlist": {
	//             "Desc": "",
	//             "Disable": false,
	//             "IP": "",
	//             "List": null,
	//             "Name": "malwaredomainlist",
	//             "Prfx": "",
	//             "Type": "hosts",
	//             "URL": "http://www.malwaredomainlist.com/hostslist/hosts.txt"
	//          },
	//          "openphish": {
	//             "Desc": "",
	//             "Disable": false,
	//             "IP": "",
	//             "List": null,
	//             "Name": "openphish",
	//             "Prfx": "http",
	//             "Type": "hosts",
	//             "URL": "https://openphish.com/feed.txt"
	//          },
	//          "someonewhocares": {
	//             "Desc": "",
	//             "Disable": false,
	//             "IP": "",
	//             "List": null,
	//             "Name": "someonewhocares",
	//             "Prfx": "0.0.0.0",
	//             "Type": "hosts",
	//             "URL": "http://someonewhocares.org/hosts/zero/"
	//          },
	//          "volkerschatz": {
	//             "Desc": "",
	//             "Disable": false,
	//             "IP": "",
	//             "List": null,
	//             "Name": "volkerschatz",
	//             "Prfx": "http",
	//             "Type": "hosts",
	//             "URL": "http://www.volkerschatz.com/net/adpaths"
	//          },
	//          "winhelp2002": {
	//             "Desc": "",
	//             "Disable": false,
	//             "IP": "",
	//             "List": null,
	//             "Name": "winhelp2002",
	//             "Prfx": "",
	//             "Type": "hosts",
	//             "URL": "http://winhelp2002.mvps.org/hosts.txt"
	//          },
	//          "yoyo": {
	//             "Desc": "",
	//             "Disable": false,
	//             "IP": "",
	//             "List": null,
	//             "Name": "yoyo",
	//             "Prfx": "\"\"",
	//             "Type": "hosts",
	//             "URL": "http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml\u0026showintro=1\u0026mimetype=plaintext"
	//          }
	//       }
	//    }
	// }
}

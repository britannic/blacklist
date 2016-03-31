package config_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/global"
	. "github.com/britannic/testutils"
)

func check(t *testing.T) {
	if !true {
		t.Skip("Not implemented; skipping tests")
	}
}

func TestAPICmd(t *testing.T) {
	want := map[string]string{
		"cfExists":           fmt.Sprintf("%v cfExists", config.API),
		"cfReturnValue":      fmt.Sprintf("%v cfReturnValue", config.API),
		"cfReturnValues":     fmt.Sprintf("%v cfReturnValues", config.API),
		"exists":             fmt.Sprintf("%v exists", config.API),
		"existsActive":       fmt.Sprintf("%v existsActive", config.API),
		"getNodeType":        fmt.Sprintf("%v getNodeType", config.API),
		"inSession":          fmt.Sprintf("%v inSession", config.API),
		"isLeaf":             fmt.Sprintf("%v isLeaf", config.API),
		"isMulti":            fmt.Sprintf("%v isMulti", config.API),
		"isTag":              fmt.Sprintf("%v isTag", config.API),
		"listActiveNodes":    fmt.Sprintf("%v listActiveNodes", config.API),
		"listNodes":          fmt.Sprintf("%v listNodes", config.API),
		"returnActiveValue":  fmt.Sprintf("%v returnActiveValue", config.API),
		"returnActiveValues": fmt.Sprintf("%v returnActiveValues", config.API),
		"returnValue":        fmt.Sprintf("%v returnValue", config.API),
		"returnValues":       fmt.Sprintf("%v returnValues", config.API),
		"showCfg":            fmt.Sprintf("%v showCfg", config.API),
		"showConfig":         fmt.Sprintf("%v showConfig", config.API),
	}

	got := config.APICmd()

	Equals(t, want, got)

	for k := range got {
		v, ok := want[k]
		switch ok {
		case true:
			Equals(t, v, got[k])

		default:
			Equals(t, want[k], got[k])

		}
	}
}

func TestGet(t *testing.T) {
	b, err := config.Get(config.Testdata, global.Area.Root)
	OK(t, err)
	Equals(t, blist, fmt.Sprint(b.String()))

	want := errors.New("Configuration data is empty, cannot continue")

	b, got := config.Get("", global.Area.Root)

	NotEquals(t, nil, got)
	Equals(t, want.Error(), got.Error())
	Equals(t, &config.Blacklist{}, b)
}

func TestGetSubdomains(t *testing.T) {
	d := config.GetSubdomains("top.one.two.three.four.five.six.com")
	for key := range d {
		Assert(t, d.KeyExists(key), fmt.Sprintf("%v key doesn't exist", key), d)
	}
}

func TestKeyExists(t *testing.T) {
	full := "top.one.two.three.four.five.six.com"
	d := config.GetSubdomains(full)
	for key := range d {
		Assert(t, d.KeyExists(key), fmt.Sprintf("%v key doesn't exist", key), d)
	}
}

func TestSHcmd(t *testing.T) {
	type query struct {
		q, r string
	}
	testSrc := []*query{
		{
			q: "listNodes",
			r: "listActiveNodes",
		},
		{
			q: "returnValue",
			r: "returnActiveValue",
		},
		{
			q: "returnValues",
			r: "returnActiveValues",
		},
		{
			q: "exists",
			r: "existsActive",
		},
		{
			q: "showConfig",
			r: "showCfg",
		},
	}

	inSession := config.Insession()

	for _, rq := range testSrc {
		got := config.SHcmd(rq.q)
		switch inSession {
		case false:
			Equals(t, rq.r, got)

		default:
			Equals(t, rq.q, got)
		}
	}
}

func TestString(t *testing.T) {
	b, err := config.Get(config.Testdata, global.Area.Root)
	OK(t, err)

	Equals(t, blist, fmt.Sprint(b))
}

func TestSubKeyExists(t *testing.T) {
	full := "top.one.two.three.four.five.six.com"
	d := config.GetSubdomains(full)
	got := len(d)
	want := strings.Count(full, ".")

	Equals(t, want, got)

	for key := range d {
		Assert(t, d.SubKeyExists(key), fmt.Sprintf("%v sub key doesn't exist", key), d)
	}
}

func TestToBool(t *testing.T) {
	tests := map[string]bool{"false": false, "true": true, "fail": false}

	for k := range tests {
		Equals(t, tests[k], config.ToBool(k))
	}
}

var (
	keys = []string{
		"six.com",
		"five.six.com",
		"four.five.six.com",
		"three.four.five.six.com",
		"two.three.four.five.six.com",
		"one.two.three.four.five.six.com",
		"top.one.two.three.four.five.six.com",
	}

	blist = `Node: blacklist
	Disabled: false
	Redirect IP: 0.0.0.0
	Exclude(s):
		122.2o7.net
		1e100.net
		adobedtm.com
		akamai.net
		amazon.com
		amazonaws.com
		apple.com
		ask.com
		avast.com
		bitdefender.com
		cdn.visiblemeasures.com
		cloudfront.net
		coremetrics.com
		edgesuite.net
		freedns.afraid.org
		github.com
		githubusercontent.com
		google.com
		googleadservices.com
		googleapis.com
		googleusercontent.com
		gstatic.com
		gvt1.com
		gvt1.net
		hb.disney.go.com
		hp.com
		hulu.com
		images-amazon.com
		msdn.com
		paypal.com
		rackcdn.com
		schema.org
		skype.com
		smacargo.com
		sourceforge.net
		ssl-on9.com
		ssl-on9.net
		static.chartbeat.com
		storage.googleapis.com
		windows.net
		yimg.com
		ytimg.com

Node: domains
	Disabled: false
	Include(s):
		adsrvr.org
		adtechus.net
		advertising.com
		centade.com
		doubleclick.net
		free-counter.co.uk
		intellitxt.com
		kiosked.com
	Source: malc0de
		Disabled: false
		Description: List of zones serving malicious executables observed by malc0de.com/database/
		Prefix: "zone "
		URL: http://malc0de.com/bl/ZONES

Node: hosts
	Disabled: false
	Include(s):
		beap.gemini.yahoo.com
	Source: adaway
		Disabled: false
		Description: Blocking mobile ad providers and some analytics providers
		Prefix: "127.0.0.1 "
		URL: http://adaway.org/hosts.txt
	Source: malwaredomainlist
		Disabled: false
		Description: 127.0.0.1 based host and domain list
		Prefix: "127.0.0.1 "
		URL: http://www.malwaredomainlist.com/hostslist/hosts.txt
	Source: openphish
		Disabled: false
		Description: OpenPhish automatic phishing detection
		Prefix: "http"
		URL: https://openphish.com/feed.txt
	Source: someonewhocares
		Disabled: false
		Description: Zero based host and domain list
		Prefix: "0.0.0.0"
		URL: http://someonewhocares.org/hosts/zero/
	Source: volkerschatz
		Disabled: false
		Description: Ad server blacklists
		Prefix: "http"
		URL: http://www.volkerschatz.com/net/adpaths
	Source: winhelp2002
		Disabled: false
		Description: Zero based host and domain list
		Prefix: "0.0.0.0 "
		URL: http://winhelp2002.mvps.org/hosts.txt
	Source: yoyo
		Disabled: false
		Description: Fully Qualified Domain Names only - no prefix to strip
		Prefix: ""
		URL: http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext

`
)

package regx_test

import (
	"fmt"
	"testing"

	"github.com/britannic/blacklist/internal/regx"
	. "github.com/smartystreets/goconvey/convey"
)

type test struct {
	index  int
	input  string
	result string
}

type config map[string]test

func TestGet(t *testing.T) {
	Convey("Testing Get()", t, func() {
		for k := range c {
			act := regx.Get(k, c[k].input)
			So(len(act), ShouldBeGreaterThan, 0)
			So(act[c[k].index], ShouldEqual, c[k].result)
		}
	})
}

func TestRegex(t *testing.T) {
	Convey("Testing Regex()", t, func() {
		So(fmt.Sprint(regx.Obj), ShouldEqual, rxout)
	})
}

func TestStripPrefixAndSuffix(t *testing.T) {
	Convey("Testing StripPrefixAndSuffix()", t, func() {
		tests := []struct {
			exp    string
			line   string
			ok     bool
			prefix string
			rx     *regx.OBJ
		}{
			{
				exp:    "This is a complete sentence and should not be a comment.",
				line:   "/* This is a complete sentence and should not be a comment.",
				ok:     true,
				prefix: "/* ",
				rx:     regx.Obj,
			},
			{
				exp:    "verybad.phishing.sites.r.us.com",
				line:   "https://verybad.phishing.sites.r.us.com",
				ok:     true,
				prefix: "https://",
				rx:     regx.Obj,
			},
			{
				exp:    "verybad.phishing.sites.r.us.com",
				line:   "https://verybad.phishing.sites.r.us.com",
				ok:     true,
				prefix: "http",
				rx:     regx.Obj,
			},
			{
				exp:    "verybad.phishing.sites.r.us.com",
				line:   "verybad.phishing.sites.r.us.com",
				ok:     false,
				prefix: "http",
				rx:     regx.Obj,
			},
		}
		for _, tt := range tests {
			act, ok := tt.rx.StripPrefixAndSuffix(tt.line, tt.prefix)
			So(act, ShouldEqual, tt.exp)
			So(ok, ShouldEqual, tt.ok)
		}
	})
}

var (
	c = config{
		"cmnt": test{
			index:  1,
			input:  `/*Comment*/`,
			result: `Comment`,
		},
		"desc": test{
			index:  1,
			input:  `description "Descriptive text"`,
			result: `Descriptive text`,
		},
		"dsbl": test{
			index:  1,
			input:  `disabled false`,
			result: `false`,
		},
		"flip": test{
			index:  1,
			input:  `address=/.xunlei.com/0.0.0.0`,
			result: `0.0.0.0`,
		},
		"fqdn": test{
			index:  1,
			input:  `http:/123pagerank.com/*=UUID:272`,
			result: `123pagerank.com`,
		},
		"host": test{
			index:  1,
			input:  `address=/.xunlei.com/0.0.0.0`,
			result: `xunlei.com`,
		},
		"http": test{
			index:  1,
			input:  `https:/123pagerank.com/*=UUID:272`,
			result: `123pagerank.com/*=UUID:272`,
		},
		"ipbh": test{
			index:  1,
			input:  `dns-redirect-ip 0.0.0.0`,
			result: `0.0.0.0`,
		},
		"lbrc": test{
			index:  0,
			input:  `blacklist {`,
			result: `{`,
		},
		"leaf": test{
			index:  1,
			input:  `source volkerschatz {`,
			result: `source`,
		},
		"misc": test{
			index:  0,
			input:  `blacklist-bigot`,
			result: `blacklist-bigot`,
		},
		"mlti": test{
			index:  2,
			input:  `include adsrvr.org`,
			result: `adsrvr.org`,
		},
		"mpty": test{
			index:  0,
			input:  ``,
			result: ``,
		},
		"name": test{
			index:  1,
			input:  `Test "System"`,
			result: `Test`,
		},
		"node": test{
			index:  1,
			input:  `hosts {`,
			result: `hosts`,
		},
		"rbrc": test{
			index:  0,
			input:  `} blacklist`,
			result: `}`,
		},
		"sufx": test{
			index:  0,
			input:  `www.123pagerank.com/*=UUID`,
			result: `/*=UUID`,
		},
	}

	rxout = "CMNT: ^(?:[\\/*]+)(.*?)(?:[*\\/]+)$\nDESC: ^(?:description)+\\s\"?([^\"]+)?\"?$\nDSBL: ^(?:disabled)+\\s([\\S]+)$\nFLIP: ^(?:address=[/][.]{0,1}.*[/])(.*)$\nFQDN: \\b((?:(?:[^.-/]{0,1})[a-zA-Z0-9-_]{1,63}[-]{0,1}[.]{1})+(?:[a-zA-Z]{2,63}))\\b\nHOST: ^(?:address=[/][.]{0,1})(.*)(?:[/].*)$\nHTTP: (?:^(?:http|https){1}:)(?:\\/|%2f){1,2}(.*)\nIPBH: ^(?:dns-redirect-ip)+\\s([\\S]+)$\nLEAF: ^([\\S]+)+\\s([\\S]+)\\s[{]{1}$\nLBRC: [{]\nMISC: ^([\\w-]+)$\nMLTI: ^((?:include|exclude)+)\\s([\\S]+)$\nMPTY: ^$\nNAME: ^([\\w-]+)\\s[\"']{0,1}(.*?)[\"']{0,1}$\nNODE: ^([\\w-]+)\\s[{]{1}$\nRBRC: [}]\nSUFX: (?:#.*|\\{.*|[/[].*)\\z\n"
)

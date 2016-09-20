package regx_test

import (
	"fmt"
	"testing"

	"github.com/britannic/blacklist/internal/regx"
	. "github.com/smartystreets/goconvey/convey"
)

type test struct {
	index  int
	input  []byte
	result []byte
}

type config map[string]test

func TestGet(t *testing.T) {
	Convey("Testing Get()", t, func() {
		for k := range c {
			act := regx.Get([]byte(k), c[k].input)
			So(len(act), ShouldBeGreaterThan, 0)
			So(act[c[k].index], ShouldResemble, c[k].result)
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
			exp    []byte
			line   []byte
			ok     bool
			prefix string
			rx     *regx.OBJ
		}{
			{
				exp:    []byte("This is a complete sentence and should not be a comment."),
				line:   []byte("/* This is a complete sentence and should not be a comment."),
				ok:     true,
				prefix: "/* ",
				rx:     regx.Obj,
			},
			{
				exp:    []byte("verybad.phishing.sites.r.us.com"),
				line:   []byte("https://verybad.phishing.sites.r.us.com"),
				ok:     true,
				prefix: "https://",
				rx:     regx.Obj,
			},
			{
				exp:    []byte("verybad.phishing.sites.r.us.com"),
				line:   []byte("https://verybad.phishing.sites.r.us.com"),
				ok:     true,
				prefix: "http",
				rx:     regx.Obj,
			},
			{
				exp:    []byte("verybad.phishing.sites.r.us.com"),
				line:   []byte("verybad.phishing.sites.r.us.com"),
				ok:     false,
				prefix: "http",
				rx:     regx.Obj,
			},
		}
		for _, tt := range tests {
			act, ok := tt.rx.StripPrefixAndSuffix(tt.line, tt.prefix)
			So(act, ShouldResemble, tt.exp)
			So(ok, ShouldEqual, tt.ok)
		}
	})
}

var (
	c = config{
		"cmnt": test{
			index:  1,
			input:  []byte(`/*Comment*/`),
			result: []byte(`Comment`),
		},
		"desc": test{
			index:  1,
			input:  []byte(`description "Descriptive text"`),
			result: []byte(`Descriptive text`),
		},
		"dsbl": test{
			index:  1,
			input:  []byte(`disabled false`),
			result: []byte(`false`),
		},
		"flip": test{
			index:  1,
			input:  []byte(`address=/.xunlei.com/0.0.0.0`),
			result: []byte(`0.0.0.0`),
		},
		"fqdn": test{
			index:  1,
			input:  []byte(`http:/123pagerank.com/*=UUID:272`),
			result: []byte(`123pagerank.com`),
		},
		"host": test{
			index:  1,
			input:  []byte(`address=/.xunlei.com/0.0.0.0`),
			result: []byte(`xunlei.com`),
		},
		"http": test{
			index:  1,
			input:  []byte(`https:/123pagerank.com/*=UUID:272`),
			result: []byte(`123pagerank.com/*=UUID:272`),
		},
		"ipbh": test{
			index:  1,
			input:  []byte(`dns-redirect-ip 0.0.0.0`),
			result: []byte(`0.0.0.0`),
		},
		"lbrc": test{
			index:  0,
			input:  []byte(`blacklist {`),
			result: []byte(`{`),
		},
		"leaf": test{
			index:  1,
			input:  []byte(`source volkerschatz {`),
			result: []byte(`source`),
		},
		"misc": test{
			index:  0,
			input:  []byte(`blacklist-bigot`),
			result: []byte(`blacklist-bigot`),
		},
		"mlti": test{
			index:  2,
			input:  []byte(`include adsrvr.org`),
			result: []byte(`adsrvr.org`),
		},
		"mpty": test{
			index:  0,
			input:  []byte{},
			result: []byte{},
		},
		"name": test{
			index:  1,
			input:  []byte(`Test "System"`),
			result: []byte(`Test`),
		},
		"node": test{
			index:  1,
			input:  []byte(`hosts {`),
			result: []byte(`hosts`),
		},
		"rbrc": test{
			index:  0,
			input:  []byte(`} blacklist`),
			result: []byte(`}`),
		},
		"sufx": test{
			index:  0,
			input:  []byte(`www.123pagerank.com/*=UUID`),
			result: []byte(`/*=UUID`),
		},
	}

	rxout = "CMNT: ^(?:[\\/*]+)(.*?)(?:[*\\/]+)$\nDESC: ^(?:description)+\\s\"?([^\"]+)?\"?$\nDSBL: ^(?:disabled)+\\s([\\S]+)$\nFLIP: ^(?:address=[/][.]{0,1}.*[/])(.*)$\nFQDN: \\b((?:(?:[^.-/]{0,1})[a-zA-Z0-9-_]{1,63}[-]{0,1}[.]{1})+(?:[a-zA-Z]{2,63}))\\b\nHOST: ^(?:address=[/][.]{0,1})(.*)(?:[/].*)$\nHTTP: (?:^(?:http|https){1}:)(?:\\/|%2f){1,2}(.*)\nIPBH: ^(?:dns-redirect-ip)+\\s([\\S]+)$\nLEAF: ^([\\S]+)+\\s([\\S]+)\\s[{]{1}$\nLBRC: [{]\nMISC: ^([\\w-]+)$\nMLTI: ^((?:include|exclude)+)\\s([\\S]+)$\nMPTY: ^$\nNAME: ^([\\w-]+)\\s[\"']{0,1}(.*?)[\"']{0,1}$\nNODE: ^([\\w-]+)\\s[{]{1}$\nRBRC: [}]\nSUFX: (?:#.*|\\{.*|[/[].*)\\z\n"
)

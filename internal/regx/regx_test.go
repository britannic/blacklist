package regx_test

import (
	"testing"

	"github.com/britannic/blacklist/internal/regx"

	. "github.com/smartystreets/goconvey/convey"
)

type test struct {
	index  int
	input  []byte
	result []byte
}

type config map[regx.Leaf]test

func TestGet(t *testing.T) {
	o := regx.NewRegex()
	Convey("Testing Get()", t, func() {
		for k := range c {
			act := o.SubMatch(k, c[k].input)
			So(len(act), ShouldBeGreaterThan, 0)
			So(act[c[k].index], ShouldResemble, c[k].result)
		}
	})
}

func TestLeafString(t *testing.T) {
	Convey("Testing LeafString()", t, func() {
		var i regx.Leaf = 1000
		So(i.String(), ShouldEqual, "CMNT")
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
				rx:     regx.NewRegex(),
			},
			{
				exp:    []byte("verybad.phishing.sites.r.us.com"),
				line:   []byte("https://verybad.phishing.sites.r.us.com"),
				ok:     true,
				prefix: "https://",
				rx:     regx.NewRegex(),
			},
			{
				exp:    []byte("verybad.phishing.sites.r.us.com"),
				line:   []byte("https://verybad.phishing.sites.r.us.com"),
				ok:     true,
				prefix: "http",
				rx:     regx.NewRegex(),
			},
			{
				exp:    []byte("verybad.phishing.sites.r.us.com"),
				line:   []byte("verybad.phishing.sites.r.us.com"),
				ok:     false,
				prefix: "http",
				rx:     regx.NewRegex(),
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
		regx.CMNT: test{
			index:  1,
			input:  []byte(`/*Comment*/`),
			result: []byte(`Comment`),
		},
		regx.DESC: test{
			index:  1,
			input:  []byte(`description "Descriptive text"`),
			result: []byte(`Descriptive text`),
		},
		regx.DSBL: test{
			index:  1,
			input:  []byte(`disabled false`),
			result: []byte(`false`),
		},
		regx.FLIP: test{
			index:  1,
			input:  []byte(`address=/.xunlei.com/0.0.0.0`),
			result: []byte(`0.0.0.0`),
		},
		regx.FQDN: test{
			index:  1,
			input:  []byte(`http:/123pagerank.com/*=UUID:272`),
			result: []byte(`123pagerank.com`),
		},
		regx.HOST: test{
			index:  1,
			input:  []byte(`address=/.xunlei.com/0.0.0.0`),
			result: []byte(`xunlei.com`),
		},
		regx.HTTP: test{
			index:  1,
			input:  []byte(`https:/123pagerank.com/*=UUID:272`),
			result: []byte(`123pagerank.com/*=UUID:272`),
		},
		regx.IPBH: test{
			index:  1,
			input:  []byte(`dns-redirect-ip 0.0.0.0`),
			result: []byte(`0.0.0.0`),
		},
		regx.LBRC: test{
			index:  0,
			input:  []byte(`blacklist {`),
			result: []byte(`{`),
		},
		regx.LEAF: test{
			index:  1,
			input:  []byte(`source volkerschatz {`),
			result: []byte(`source`),
		},
		regx.MISC: test{
			index:  0,
			input:  []byte(`blacklist-bigot`),
			result: []byte(`blacklist-bigot`),
		},
		regx.MLTI: test{
			index:  2,
			input:  []byte(`include adsrvr.org`),
			result: []byte(`adsrvr.org`),
		},
		regx.MPTY: test{
			index:  0,
			input:  []byte{},
			result: []byte{},
		},
		regx.NAME: test{
			index:  1,
			input:  []byte(`Test "System"`),
			result: []byte(`Test`),
		},
		regx.NODE: test{
			index:  1,
			input:  []byte(`hosts {`),
			result: []byte(`hosts`),
		},
		regx.RBRC: test{
			index:  0,
			input:  []byte(`} blacklist`),
			result: []byte(`}`),
		},
		regx.SUFX: test{
			index:  0,
			input:  []byte(`www.123pagerank.com/*=UUID`),
			result: []byte(`/*=UUID`),
		},
	}

	// rxout = "CMNT: ^(?:[\\/*]+)(.*?)(?:[*\\/]+)$\nDESC: ^(?:description)+\\s\"?([^\"]+)?\"?$\nDSBL: ^(?:disabled)+\\s([\\S]+)$\nFLIP: ^(?:address=[/][.]{0,1}.*[/])(.*)$\nFQDN: \\b((?:(?:[^.-/]{0,1})[a-zA-Z0-9-_]{1,63}[-]{0,1}[.]{1})+(?:[a-zA-Z]{2,63}))\\b\nHOST: ^(?:address=[/][.]{0,1})(.*)(?:[/].*)$\nHTTP: (?:^(?:http|https){1}:)(?:\\/|%2f){1,2}(.*)\nIPBH: ^(?:dns-redirect-ip)+\\s([\\S]+)$\nLEAF: ^([\\S]+)+\\s([\\S]+)\\s[{]{1}$\nLBRC: [{]\nMISC: ^([\\w-]+)$\nMLTI: ^((?:include|exclude)+)\\s([\\S]+)$\nMPTY: ^$\nNAME: ^([\\w-]+)\\s[\"']{0,1}(.*?)[\"']{0,1}$\nNODE: ^([\\w-]+)\\s[{]{1}$\nRBRC: [}]\nSUFX: (?:#.*|\\{.*|[/[].*)\\z\n"
)

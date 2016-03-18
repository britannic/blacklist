package regx

import (
	"fmt"
	"testing"

	"github.com/britannic/blacklist/utils"
)

type test struct {
	index  int
	input  string
	result string
}

type config map[string]test

func TestGet(t *testing.T) {

	for k := range c {
		match := Get(k, c[k].input)
		if len(match) == 0 {
			t.Fatalf("%v results fail: %v", k, match)
		}
		if match[c[k].index] != c[k].result {
			t.Errorf("%v match fail: %v", k, match)
		}
	}
}

func TestRegex(t *testing.T) {
	rx := Regex
	rxtest := make(map[string][]byte)
	rxtest["want"] = append(rxtest["want"], rxout...)
	rxtest["got"] = append(rxtest["got"], fmt.Sprint(rx)...)

	if !utils.CmpHash(rxtest["got"], rxtest["want"]) {
		t.Errorf("Got: %v Want: %v", string(rxtest["got"]), string(rxtest["want"]))
	}
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

	rxout = `CMNT: ^(?:[\/*]+)(.*?)(?:[*\/]+)$
DESC: ^(?:description)+\s"?([^"]+)?"?$
DSBL: ^(?:disabled)+\s([\S]+)$
FLIP: ^(?:address=[/][.]{0,1}.*[/])(.*)$
FQDN: \b((?:(?:[^.-/]{0,1})[a-zA-Z0-9-_]{1,63}[-]{0,1}[.]{1})+(?:[a-zA-Z]{2,63}))\b
HOST: ^(?:address=[/][.]{0,1})(.*)(?:[/].*)$
HTTP: (?:^(?:http|https){1}:)(?:\/|%2f){1,2}(.*)
LEAF: ^(source)+\s([\S]+)\s[{]{1}$
LBRC: [{]
MISC: ^([\w-]+)$
MLTI: ^((?:include|exclude)+)\s([\S]+)$
MPTY: ^$
NAME: ^([\w-]+)\s["']{0,1}(.*?)["']{0,1}$
NODE: ^([\w-]+)\s[{]{1}$
RBRC: [}]
SUFX: (?:#.*|\{.*|[/[].*)\z
`
)

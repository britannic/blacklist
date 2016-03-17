package regx

import "testing"

type test struct {
	index  int
	input  string
	result string
}

type config map[string]test

func TestGet(t *testing.T) {
	c := config{
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

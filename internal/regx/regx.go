// Package regx provides regex objects for processing data in files and web content
package regx

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Leaf is a config label
type Leaf int

//go:generate stringer -type=Leaf
// Leaf label regx map keys
const (
	CMNT Leaf = iota + 1000
	DESC
	DSBL
	FLIP
	FQDN
	HOST
	HTTP
	IPBH
	LEAF
	LBRC
	MISC
	MLTI
	MPTY
	NAME
	NODE
	RBRC
	SUFX
)

// OBJ is a map of regex precompiled objects

type regexMap map[Leaf]*regexp.Regexp

// OBJ is a struct of regex precompiled objects
type OBJ struct {
	RX regexMap
}

// NewRegex returns a map of OBJ populated with a map of precompiled regex objects
func NewRegex() *OBJ {
	return &OBJ{
		RX: regexMap{
			CMNT: regexp.MustCompile(`^(?:[\/*]+)(.*?)(?:[*\/]+)$`),
			DESC: regexp.MustCompile(`^(?:description)+\s"?([^"]+)?"?$`),
			DSBL: regexp.MustCompile(`^(?:disabled)+\s([\S]+)$`),
			FLIP: regexp.MustCompile(`^(?:address=[/][.]{0,1}.*[/])(.*)$`),
			FQDN: regexp.MustCompile(`\b((?:(?:[^.-/]{0,1})[a-zA-Z0-9-_]{1,63}[-]{0,1}[.]{1})+(?:[a-zA-Z]{2,63}))\b`),
			HOST: regexp.MustCompile(`^(?:address=[/][.]{0,1})(.*)(?:[/].*)$`),
			HTTP: regexp.MustCompile(`(?:^(?:http|https){1}:)(?:\/|%2f){1,2}(.*)`),
			IPBH: regexp.MustCompile(`^(?:dns-redirect-ip)+\s([\S]+)$`),
			LBRC: regexp.MustCompile(`[{]`),
			LEAF: regexp.MustCompile(`^([\S]+)+\s([\S]+)\s[{]{1}$`),
			MISC: regexp.MustCompile(`^([\w-]+)$`),
			MLTI: regexp.MustCompile(`^((?:include|exclude)+)\s([\S]+)$`),
			MPTY: regexp.MustCompile(`^$`),
			NAME: regexp.MustCompile(`^([\w-]+)\s["']{0,1}(.*?)["']{0,1}$`),
			NODE: regexp.MustCompile(`^([\w-]+)\s[{]{1}$`),
			RBRC: regexp.MustCompile(`[}]`),
			SUFX: regexp.MustCompile(`(?:#.*|\{.*|[/[].*)\z`),
		},
	}
}

// SubMatch extracts the configuration value for a matched label
func (o *OBJ) SubMatch(t Leaf, b []byte) [][]byte {
	return o.RX[t].FindSubmatch(b)
}

func (o *OBJ) String() string {
	var a []string
	for k, v := range o.RX {
		a = append(a, fmt.Sprintf("%v: %v", k.String(), v))
	}
	sort.Strings(a)
	return strings.Join(a, "\n")
}

// StripPrefixAndSuffix strips the prefix and suffix
func (o *OBJ) StripPrefixAndSuffix(l []byte, p string) ([]byte, bool) {
	switch {
	case p == "http", p == "https":
		if !o.RX[HTTP].Match(l) {
			return l, false
		}
		l = o.RX[HTTP].FindSubmatch(l)[1]

	case bytes.HasPrefix(l, []byte(p)):
		l = bytes.TrimPrefix(l, []byte(p))
	}

	l = o.RX[SUFX].ReplaceAll(l, []byte{})
	l = bytes.Replace(l, []byte(`"`), []byte{}, -1)
	l = bytes.TrimSpace(l)
	return l, true
}

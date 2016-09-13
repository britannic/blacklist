// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package regx provides regex objects for processing data in files and web content
package regx

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// OBJ is a struct of regex precompiled objects
type OBJ struct {
	CMNT, DESC, DSBL, FLIP, FQDN, HOST, HTTP, IPBH, LEAF, LBRC, MISC, MLTI, MPTY, NAME, NODE, RBRC, SUFX *regexp.Regexp
}

// Obj is a struct of *OBJ populated with precompiled regex objects
var Obj = &OBJ{
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
}

// Get returns an array compiled regx OBJs
func Get(t, s string) (r []string) {
	rx := Obj
	switch t {
	case "cmnt":
		r = rx.CMNT.FindStringSubmatch(s)
	case "desc":
		r = rx.DESC.FindStringSubmatch(s)
	case "dsbl":
		r = rx.DSBL.FindStringSubmatch(s)
	case "flip":
		r = rx.FLIP.FindStringSubmatch(s)
	case "fqdn":
		r = rx.FQDN.FindStringSubmatch(s)
	case "host":
		r = rx.HOST.FindStringSubmatch(s)
	case "http":
		r = rx.HTTP.FindStringSubmatch(s)
	case "ipbh":
		r = rx.IPBH.FindStringSubmatch(s)
	case "lbrc":
		r = rx.LBRC.FindStringSubmatch(s)
	case "leaf":
		r = rx.LEAF.FindStringSubmatch(s)
	case "misc":
		r = rx.MISC.FindStringSubmatch(s)
	case "mlti":
		r = rx.MLTI.FindStringSubmatch(s)
	case "mpty":
		r = rx.MPTY.FindStringSubmatch(s)
	case "name":
		r = rx.NAME.FindStringSubmatch(s)
	case "node":
		r = rx.NODE.FindStringSubmatch(s)
	case "rbrc":
		r = rx.RBRC.FindStringSubmatch(s)
	case "sufx":
		r = rx.SUFX.FindStringSubmatch(s)
	}
	return r
}

// iter iterates over ints - use it in for loops
func iter(i int) []struct{} {
	return make([]struct{}, i)
}

func (rx *OBJ) String() (s string) {
	v := reflect.ValueOf(rx).Elem()
	for i := range iter(v.NumField()) {
		s += fmt.Sprintf("%v: %v\n", v.Type().Field(i).Name, v.Field(i).Interface())
	}
	return s
}

// StripPrefixAndSuffix strips the prefix and suffix
func (rx *OBJ) StripPrefixAndSuffix(line, prefix string) (string, bool) {
	switch {
	case prefix == "http", prefix == "https":
		if !rx.HTTP.MatchString(line) {
			return line, false
		}
		line = rx.HTTP.FindStringSubmatch(line)[1]

	case strings.HasPrefix(line, prefix):
		line = strings.TrimPrefix(line, prefix)
	}

	line = rx.SUFX.ReplaceAllString(line, "")
	line = strings.Replace(line, `"`, "", -1)
	line = strings.TrimSpace(line)
	return line, true
}

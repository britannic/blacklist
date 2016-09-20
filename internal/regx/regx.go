// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package regx provides regex objects for processing data in files and web content
package regx

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
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
func Get(t, b []byte) (r [][]byte) {
	rx := Obj
	switch string(t) {
	case "cmnt":
		r = rx.CMNT.FindSubmatch(b)
	case "desc":
		r = rx.DESC.FindSubmatch(b)
	case "dsbl":
		r = rx.DSBL.FindSubmatch(b)
	case "flip":
		r = rx.FLIP.FindSubmatch(b)
	case "fqdn":
		r = rx.FQDN.FindSubmatch(b)
	case "host":
		r = rx.HOST.FindSubmatch(b)
	case "http":
		r = rx.HTTP.FindSubmatch(b)
	case "ipbh":
		r = rx.IPBH.FindSubmatch(b)
	case "lbrc":
		r = rx.LBRC.FindSubmatch(b)
	case "leaf":
		r = rx.LEAF.FindSubmatch(b)
	case "misc":
		r = rx.MISC.FindSubmatch(b)
	case "mlti":
		r = rx.MLTI.FindSubmatch(b)
	case "mpty":
		r = rx.MPTY.FindSubmatch(b)
	case "name":
		r = rx.NAME.FindSubmatch(b)
	case "node":
		r = rx.NODE.FindSubmatch(b)
	case "rbrc":
		r = rx.RBRC.FindSubmatch(b)
	case "sufx":
		r = rx.SUFX.FindSubmatch(b)
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
func (rx *OBJ) StripPrefixAndSuffix(line []byte, prefix string) ([]byte, bool) {
	switch {
	case prefix == "http", prefix == "https":
		if !rx.HTTP.Match(line) {
			return line, false
		}
		line = rx.HTTP.FindSubmatch(line)[1]

	case bytes.HasPrefix(line, []byte(prefix)):
		line = bytes.TrimPrefix(line, []byte(prefix))
	}

	line = rx.SUFX.ReplaceAll(line, []byte{})
	line = bytes.Replace(line, []byte(`"`), []byte{}, -1)
	line = bytes.TrimSpace(line)
	return line, true
}

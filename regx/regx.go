// Copyright 2016 NJ Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.txt file.

// Package regx provides regex objects for processing data in files and web content
package regx

import (
	"fmt"
	"reflect"
	"regexp"
)

// RGX is a struct of regex precompiled objects
type RGX struct {
	CMNT, DESC, DSBL, FLIP, FQDN, HOST, HTTP, LEAF, LBRC, MISC, MLTI, MPTY, NAME, NODE, RBRC, SUFX *regexp.Regexp
}

// Get returns an array of the string and submatch
func Get(t, s string) (r []string) {
	rx := Regex
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
	return
}

func (rx *RGX) String() (result string) {
	v := reflect.ValueOf(rx).Elem()
	// values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		// values[i] = v.Field(i).Interface()

		result += fmt.Sprintf("%v: %v\n", v.Type().Field(i).Name, v.Field(i).Interface())
	}
	return result
}

// Regex is a struct of *re populated with precompiled regex objects
var Regex = &RGX{
	CMNT: regexp.MustCompile(`^(?:[\/*]+)(.*?)(?:[*\/]+)$`),
	DESC: regexp.MustCompile(`^(?:description)+\s"?([^"]+)?"?$`),
	DSBL: regexp.MustCompile(`^(?:disabled)+\s([\S]+)$`),
	FLIP: regexp.MustCompile(`^(?:address=[/][.]{0,1}.*[/])(.*)$`),
	FQDN: regexp.MustCompile(`\b((?:(?:[^.-/]{0,1})[a-zA-Z0-9-_]{1,63}[-]{0,1}[.]{1})+(?:[a-zA-Z]{2,63}))\b`),
	HOST: regexp.MustCompile(`^(?:address=[/][.]{0,1})(.*)(?:[/].*)$`),
	HTTP: regexp.MustCompile(`(?:^(?:http|https){1}:)(?:\/|%2f){1,2}(.*)`),
	LBRC: regexp.MustCompile(`[{]`),
	LEAF: regexp.MustCompile(`^(source)+\s([\S]+)\s[{]{1}$`),
	MISC: regexp.MustCompile(`^([\w-]+)$`),
	MLTI: regexp.MustCompile(`^((?:include|exclude)+)\s([\S]+)$`),
	MPTY: regexp.MustCompile(`^$`),
	NAME: regexp.MustCompile(`^([\w-]+)\s["']{0,1}(.*?)["']{0,1}$`),
	NODE: regexp.MustCompile(`^([\w-]+)\s[{]{1}$`),
	RBRC: regexp.MustCompile(`[}]`),
	SUFX: regexp.MustCompile(`(?:#.*|\{.*|[/[].*)\z`),
}

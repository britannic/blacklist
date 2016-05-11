package regx_test

import (
	"fmt"

	"github.com/britannic/blacklist/regx"
)

func ExampleOBJ() {
	rx := regx.Objects
	fmt.Println(rx)
	// Output: CMNT: ^(?:[\/*]+)(.*?)(?:[*\/]+)$
	// DESC: ^(?:description)+\s"?([^"]+)?"?$
	// DSBL: ^(?:disabled)+\s([\S]+)$
	// FLIP: ^(?:address=[/][.]{0,1}.*[/])(.*)$
	// FQDN: \b((?:(?:[^.-/]{0,1})[a-zA-Z0-9-_]{1,63}[-]{0,1}[.]{1})+(?:[a-zA-Z]{2,63}))\b
	// HOST: ^(?:address=[/][.]{0,1})(.*)(?:[/].*)$
	// HTTP: (?:^(?:http|https){1}:)(?:\/|%2f){1,2}(.*)
	// IPBH: ^(?:dns-redirect-ip)+\s([\S]+)$
	// LEAF: ^([\S]+)+\s([\S]+)\s[{]{1}$
	// LBRC: [{]
	// MISC: ^([\w-]+)$
	// MLTI: ^((?:include|exclude)+)\s([\S]+)$
	// MPTY: ^$
	// NAME: ^([\w-]+)\s["']{0,1}(.*?)["']{0,1}$
	// NODE: ^([\w-]+)\s[{]{1}$
	// RBRC: [}]
	// SUFX: (?:#.*|\{.*|[/[].*)\z
}

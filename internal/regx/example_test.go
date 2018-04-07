package regx_test

import (
	"fmt"
	"testing"

	"github.com/britannic/blacklist/internal/regx"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOBJ(t *testing.T) {
	rxmap := regx.NewRegex()
	Convey("Testing regx.Obj", t, func() {
		act := fmt.Sprint(rxmap)
		exp := `CMNT: ^(?:[\/*]+)(.*?)(?:[*\/]+)$
DESC: ^(?:description)+\s"?([^"]+)?"?$
DSBL: ^(?:disabled)+\s([\S]+)$
FLIP: ^(?:address=[/][.]{0,1}.*[/])(.*)$
FQDN: \b((?:(?:[^.-/]{0,1})[a-zA-Z0-9-_]{1,63}[-]{0,1}[.]{1})+(?:[a-zA-Z]{2,63}))\b
HOST: ^(?:address=[/][.]{0,1})(.*)(?:[/].*)$
HTTP: (?:^(?:http|https){1}:)(?:\/|%2f){1,2}(.*)
IPBH: ^(?:dns-redirect-ip)+\s([\S]+)$
LBRC: [{]
LEAF: ^([\S]+)+\s([\S]+)\s[{]{1}$
MISC: ^([\w-]+)$
MLTI: ^((?:include|exclude)+)\s([\S]+)$
MPTY: ^$
NAME: ^([\w-]+)\s["']{0,1}(.*?)["']{0,1}$
NODE: ^([\w-]+)\s[{]{1}$
RBRC: [}]
SUFX: (?:#.*|\{.*|[/[].*)\z`
		So(act, ShouldEqual, exp)
	})
}

package edgeos

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChkWeb(t *testing.T) {
	Convey("Testing TestChkWeb()", t, func() {
		tests := []struct {
			exp  bool
			port int
			site string
		}{
			{exp: true, site: "www.google.com", port: 443},
			{exp: true, site: "yahoo.com", port: 80},
			{exp: true, site: "bing.com", port: 443},
			{exp: false, site: "bigtop.@@@", port: 80},
		}
		for _, tt := range tests {
			Convey("Current test: "+tt.site, func() {
				got := ChkWeb(tt.site, tt.port)
				So(tt.exp, ShouldEqual, got)
			})
		}
	})

}

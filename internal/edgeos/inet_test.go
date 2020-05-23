package edgeos

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Chk_Web() returns true if DNS is working
func TestChk_Web(t *testing.T) {
	Convey("Testing TestChk_Web()", t, func() {
		tests := []struct {
			exp  bool
			port string
			site string
		}{
			{exp: true, site: "google.com", port: "53"},
			{exp: true, site: "google.com", port: "80"},
			{exp: true, site: "google.com", port: "443"},
			{exp: false, site: "bigtop.@@@", port: "80"},
		}
		// got := Chk_Web("google.com", "80")

		// So(got, ShouldBeTrue)
		for _, tt := range tests {
			got := Chk_Web(tt.site, tt.port)
			So(tt.exp, ShouldEqual, got)
		}
	})

}

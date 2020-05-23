package edgeos

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Chk_Web() returns true if DNS is working
func TestChk_Web(t *testing.T) {
	Convey("Testing Test_Chk_Web()", t, func() {
		tests := []struct {
			exp  bool
			port string
			site string
		}{
			{exp: true, site: "google.com", port: "53"},
			{exp: true, site: "google.com", port: "80"},
			{exp: true, site: "google.com", port: "443"},
			{exp: false, site: "bigtop.@@@", port: "443"},
		}

		for _, tt := range tests {
			So(tt.exp, ShouldEqual, Chk_Web(tt.site, tt.port))
		}
	})

}

package edgeos

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArea(t *testing.T) {
	Convey("Testing area()", t, func() {
		tests := []struct {
			exp  string
			name string
			s    *source
		}{
			{
				name: roots,
				s: &source{
					nType: root,
				},
				exp: roots,
			},
			{
				name: PreDomns,
				s: &source{
					nType: preDomn,
				},
				exp: domains,
			},
			{
				name: hosts,
				s: &source{
					nType: host,
				},
				exp: hosts,
			},
		}

		for _, tt := range tests {
			Convey("with "+tt.name, func() {
				So(tt.s.area(), ShouldEqual, tt.exp)
			})
		}
	})
}

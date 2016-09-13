package edgeos

import (
	"fmt"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMsgMethods(t *testing.T) {
	Convey("Testing Msg methods...", t, func() {
		tests := []struct {
			name  string
			dupes int
			gt    int
			new   int
			total int
			uniq  int
		}{
			{
				name:  "test-1",
				gt:    7,
				dupes: 0,
				new:   7,
				total: 0,
				uniq:  0,
			},
			{
				name:  "test-2",
				dupes: 0,
				gt:    34,
				new:   17,
				total: 17,
				uniq:  0,
			},
			{
				name:  "test-3",
				dupes: 0,
				gt:    1034,
				new:   1017,
				total: 17,
				uniq:  0,
			},
			{
				name:  "test-4",
				dupes: 51,
				gt:    103400,
				new:   101700,
				total: 1700,
				uniq:  103000,
			},
		}

		for _, tt := range tests {
			m := NewMsg(tt.name)
			m.New = tt.new
			m.Total = tt.total
			m.GetTotal()
			So(m.Total, ShouldEqual, tt.gt)

			m.Dupes = tt.dupes
			So(m.ReadDupes(), ShouldEqual, tt.dupes)

			m.New = tt.new
			So(m.ReadNew(), ShouldEqual, tt.new)

			m.Uniq = tt.uniq
			So(m.ReadUniq(), ShouldEqual, tt.uniq)
		}
	})
}

func TestMsgIncDupe(t *testing.T) {
	Convey("Testing incDupe()", t, func() {
		tests := []struct {
			name  string
			dupes int
			gt    int
		}{
			{
				name:  "test-1",
				dupes: 0,
				gt:    1,
			},
			{
				name:  "test-2",
				dupes: 5,
				gt:    6,
			},
			{
				name:  "test-3",
				dupes: 10000005,
				gt:    10000006,
			},
		}
		for _, tt := range tests {
			m := NewMsg(tt.name)
			m.Dupes = tt.dupes
			m.incDupe()
			So(m.Dupes, ShouldEqual, tt.gt)
		}
	})
}

func TestMsgIncNew(t *testing.T) {
	Convey("Testing incNew()", t, func() {
		tests := []struct {
			name string
			new  int
			gt   int
		}{
			{
				name: "test-1",
				new:  0,
				gt:   1,
			},
			{
				name: "test-2",
				new:  5,
				gt:   6,
			},
			{
				name: "test-3",
				new:  10000005,
				gt:   10000006,
			},
		}
		for _, tt := range tests {
			m := NewMsg(tt.name)
			m.New = tt.new
			m.incNew()
			So(m.New, ShouldEqual, tt.gt)
		}
	})
}

func TestMsgIncUniq(t *testing.T) {
	Convey("Testing IncUniq()", t, func() {
		tests := []struct {
			name string
			uniq int
			gt   int
		}{
			{
				name: "test-1",
				uniq: 0,
				gt:   1,
			},
			{
				name: "test-2",
				uniq: 5,
				gt:   6,
			},
			{
				name: "test-3",
				uniq: 10000005,
				gt:   10000006,
			},
		}
		for _, tt := range tests {
			m := NewMsg(tt.name)
			m.Uniq = tt.uniq
			m.incUniq()
			So(m.Uniq, ShouldEqual, tt.gt)
		}
	})
}

func TestNewMsg(t *testing.T) {
	Convey("Testing NewMsg('Vanilla')", t, func() {
		exp := &Msg{
			Name:    "Vanilla",
			RWMutex: &sync.RWMutex{},
		}
		So(NewMsg("Vanilla"), ShouldResemble, exp)
	})
}

func TestMsgString(t *testing.T) {
	Convey("Testing Msg.String()", t, func() {
		act := NewMsg("JSON-Print")
		exp := `{
	"Name": "JSON-Print",
	"dupes": 0,
	"new": 0,
	"total": 0,
	"uniq": 0
}`
		So(fmt.Sprint(act), ShouldEqual, exp)
	})
}

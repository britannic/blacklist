package edgeos

import (
	"fmt"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMsgGetTotal(t *testing.T) {
	Convey("Testing GetTotal()", t, func() {
		// sum := func(a int32, b int32) int32 { return a + b }
		tests := []struct {
			name  string
			dupes int32
			gt    int32
			new   int32
			total int32
			uniq  int32
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
		}

		for _, tt := range tests {
			m := NewMsg("test-1")
			m.New = tt.new
			m.Total = tt.total
			m.GetTotal()
			So(m.Total, ShouldEqual, tt.gt)
		}
	})
}

func TestMsgIncDupe(t *testing.T) {
	Convey("Testing IncDupe()", t, func() {
		tests := []struct {
			name  string
			dupes int32
			gt    int32
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
			m.IncDupe()
			So(m.Dupes, ShouldEqual, tt.gt)
		}
	})
}

func TestMsgIncNew(t *testing.T) {
	Convey("Testing IncNew()", t, func() {
		tests := []struct {
			name string
			new  int32
			gt   int32
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
			m.IncNew()
			So(m.New, ShouldEqual, tt.gt)
		}
	})
}

func TestMsgIncUniq(t *testing.T) {
	Convey("Testing IncUniq()", t, func() {
		tests := []struct {
			name string
			uniq int32
			gt   int32
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
			m.IncUniq()
			So(m.Uniq, ShouldEqual, tt.gt)
		}
	})
}

func TestNewMsg(t *testing.T) {
	Convey("Testing NewMsg('Vanilla')", t, func() {
		exp := &Msg{
			Name: "Vanilla",
			Rec: &Rec{
				RWMutex: &sync.RWMutex{},
			},
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

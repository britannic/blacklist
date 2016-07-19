package edgeos

import (
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoad(t *testing.T) {
	Convey("Testing Load()", t, func() {
		c := NewConfig(
			API("/bin/cli-shell-api"),
			Bash("/bin/bash"),
			InCLI("inSession"),
			Level("service dns forwarding"),
		)

		cfg, err := c.load("zBroken", "service dns forwarding")
		So(err, ShouldNotBeNil)

		cfg, err = c.load("showConfig", "")
		So(err, ShouldNotBeNil)

		r := CFGcli{Config: c}
		act, err := ioutil.ReadAll(r.read())
		So(err, ShouldBeNil)
		So(string(act), ShouldEqual, "")

		cfg, err = c.load("echo", "true")
		So(err, ShouldNotBeNil)
		So(cfg, ShouldResemble, []byte{})
	})
}

func TestPurgeFiles(t *testing.T) {
	Convey("Testing PurgeFiles()", t, func() {
		var (
			dir       = "/tmp"
			ext       = ".delete"
			purgeList []string
			exp       error
		)

		for i := 0; i < 10; i++ {
			fname := fmt.Sprintf("%v%v", i, ext)
			f, err := ioutil.TempFile(dir, fname)
			So(err, ShouldBeNil)
			purgeList = append(purgeList, f.Name())
		}

		err := purgeFiles(purgeList)
		So(err, ShouldBeNil)

		act := purgeFiles(purgeList)
		So(act, ShouldEqual, exp)

		act = purgeFiles([]string{"/dev/null"})
		exp = fmt.Errorf(`could not remove "/dev/null"`)
		So(act, ShouldResemble, exp)
	})
}

func TestAPICMD(t *testing.T) {
	Convey("Testing APICMD()", t, func() {
		tests := []struct {
			b    bool
			q, r string
		}{
			{
				b: false,
				q: "listNodes",
				r: "listNodes",
			},
			{
				b: true,
				q: "listNodes",
				r: "listActiveNodes",
			},
			{
				b: false,
				q: "listActiveNodes",
				r: "listNodes",
			},
			{
				b: false,
				q: "returnValue",
				r: "returnValue",
			},
			{
				b: true,
				q: "returnValue",
				r: "returnActiveValue",
			},
			{
				b: false,
				q: "returnActiveValue",
				r: "returnValue",
			},
			{
				b: false,
				q: "returnValues",
				r: "returnValues",
			},
			{
				b: true,
				q: "returnValues",
				r: "returnActiveValues",
			},
			{
				b: false,
				q: "returnActiveValues",
				r: "returnValues",
			},
			{
				b: false,
				q: "exists",
				r: "exists",
			},
			{
				b: true,
				q: "exists",
				r: "existsActive",
			},
			{
				b: false,
				q: "existsActive",
				r: "exists",
			},
			{
				b: false,
				q: "showCfg",
				r: "showCfg",
			},
			{
				b: true,
				q: "showCfg",
				r: "showConfig",
			},
			{
				b: false,
				q: "showConfig",
				r: "showCfg",
			},
		}

		for _, tt := range tests {
			So(apiCMD(tt.q, tt.b), ShouldEqual, tt.r)
		}

		c := NewConfig(
			API("/bin/cli-shell-api"),
			InCLI("inSession"),
			Level("service dns forwarding"),
		)
		act := fmt.Sprintf("%v %v", apiCMD("showConfig", c.InSession()), c.Level)
		exp := "showCfg service dns forwarding"
		So(act, ShouldEqual, exp)
	})
}

func TestDeleteFile(t *testing.T) {
	Convey("Testing DeleteFile()", t, func() {

		dir := "../testdata"
		ext := "delete.me"

		tests := []struct {
			name string
			f    string
			exp  bool
		}{
			{
				name: "exists",
				f:    fmt.Sprintf("%v%v", "goodFile", ext),
				exp:  true,
			},
			{
				name: "non-existent",
				f:    fmt.Sprintf("%v%v", "badFile", ext),
				exp:  true,
			},
		}

		for _, tt := range tests {
			switch tt.name {
			case "exists":
				f, err := ioutil.TempFile(dir, tt.f)
				So(err, ShouldBeNil)
				So(deleteFile(f.Name()), ShouldEqual, tt.exp)
			default:
				So(deleteFile(tt.f), ShouldEqual, tt.exp)
			}
		}
	})
}

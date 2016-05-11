package edgeos

import (
	"fmt"
	"os"
	"testing"

	. "github.com/britannic/testutils"
)

func TestAPICmd(t *testing.T) {
	want := map[string]string{
		"cfExists":           fmt.Sprintf("%v cfExists", API),
		"cfReturnValue":      fmt.Sprintf("%v cfReturnValue", API),
		"cfReturnValues":     fmt.Sprintf("%v cfReturnValues", API),
		"echo":               "echo",
		"exists":             fmt.Sprintf("%v exists", API),
		"existsActive":       fmt.Sprintf("%v existsActive", API),
		"getNodeType":        fmt.Sprintf("%v getNodeType", API),
		"inSession":          fmt.Sprintf("%v inSession", API),
		"isLeaf":             fmt.Sprintf("%v isLeaf", API),
		"isMulti":            fmt.Sprintf("%v isMulti", API),
		"isTag":              fmt.Sprintf("%v isTag", API),
		"listActiveNodes":    fmt.Sprintf("%v listActiveNodes", API),
		"listNodes":          fmt.Sprintf("%v listNodes", API),
		"pwd":                "pwd",
		"returnActiveValue":  fmt.Sprintf("%v returnActiveValue", API),
		"returnActiveValues": fmt.Sprintf("%v returnActiveValues", API),
		"returnValue":        fmt.Sprintf("%v returnValue", API),
		"returnValues":       fmt.Sprintf("%v returnValues", API),
		"showCfg":            fmt.Sprintf("%v showCfg", API),
		"showConfig":         fmt.Sprintf("%v showConfig", API),
	}

	got := APICmd()

	Equals(t, want, got)

	for k := range got {
		v, ok := want[k]
		switch ok {
		case true:
			Equals(t, v, got[k])

		default:
			Equals(t, want[k], got[k])
		}
	}
}

func TestSHCmd(t *testing.T) {
	type query struct {
		q, r string
	}
	testSrc := []*query{
		{
			q: "listNodes",
			r: "listActiveNodes",
		},
		{
			q: "returnValue",
			r: "returnActiveValue",
		},
		{
			q: "returnValues",
			r: "returnActiveValues",
		},
		{
			q: "exists",
			r: "existsActive",
		},
		{
			q: "showConfig",
			r: "showCfg",
		},
	}

	runTest := func(inSession bool) {
		for _, rq := range testSrc {
			got := SHCmd(rq.q)
			switch inSession {
			case false:
				Equals(t, rq.r, got)

			default:
				Equals(t, rq.q, got)
			}
		}
	}

	err := os.Setenv("_OFR_CONFIGURE", "ok")
	OK(t, err)
	Equals(t, "ok", os.ExpandEnv("$_OFR_CONFIGURE"))
	runTest(Insession())
	err = os.Setenv("_OFR_CONFIGURE", "")
	OK(t, err)
	Equals(t, "", os.ExpandEnv("$_OFR_CONFIGURE"))
	runTest(Insession())
}

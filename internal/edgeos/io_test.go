package edgeos

import (
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/britannic/testutils"
)

func TestAPICMD(t *testing.T) {
	want := map[string]string{
		"cfExists":           fmt.Sprintf("%v cfExists", api),
		"cfReturnValue":      fmt.Sprintf("%v cfReturnValue", api),
		"cfReturnValues":     fmt.Sprintf("%v cfReturnValues", api),
		"echo":               "echo",
		"exists":             fmt.Sprintf("%v exists", api),
		"existsActive":       fmt.Sprintf("%v existsActive", api),
		"getNodeType":        fmt.Sprintf("%v getNodeType", api),
		"inSession":          fmt.Sprintf("%v inSession", api),
		"isLeaf":             fmt.Sprintf("%v isLeaf", api),
		"isMulti":            fmt.Sprintf("%v isMulti", api),
		"isTag":              fmt.Sprintf("%v isTag", api),
		"listActiveNodes":    fmt.Sprintf("%v listActiveNodes", api),
		"listNodes":          fmt.Sprintf("%v listNodes", api),
		"returnActiveValue":  fmt.Sprintf("%v returnActiveValue", api),
		"returnActiveValues": fmt.Sprintf("%v returnActiveValues", api),
		"returnValue":        fmt.Sprintf("%v returnValue", api),
		"returnValues":       fmt.Sprintf("%v returnValues", api),
		"showCfg":            fmt.Sprintf("%v showCfg", api),
		"showConfig":         fmt.Sprintf("%v showConfig", api),
	}

	got := apiCMD()

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

func TestLoad(t *testing.T) {
	cfg, err := load("zBroken", "service dns forwarding")
	NotOK(t, err)

	cfg, err = load("showConfig", "")
	NotOK(t, err)

	cfg, err = load("echo", "Test")
	OK(t, err)
	Equals(t, "Test\n", cfg)

	r := CFGcli{}
	got, err := ioutil.ReadAll(r.Load())
	OK(t, err)
	Equals(t, "", string(got))
}

func TestPurgeFiles(t *testing.T) {
	var (
		dir       = "/tmp"
		ext       = ".delete"
		purgeList []string
		want      error
	)

	for i := 0; i < 10; i++ {
		fname := fmt.Sprintf("%v%v", i, ext)
		f, err := ioutil.TempFile(dir, fname)
		OK(t, err)
		purgeList = append(purgeList, f.Name())
	}

	err := purgeFiles(purgeList)
	OK(t, err)

	got := purgeFiles(purgeList)
	Equals(t, want, got)

	got = purgeFiles([]string{"/dev/null"})
	want = fmt.Errorf(`could not remove "/dev/null"`)
	Equals(t, want, got)
}

func TestSHCMD(t *testing.T) {
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

	inSession := insession()

	for _, rq := range testSrc {
		got := shCMD(rq.q)
		switch inSession {
		case false:
			Equals(t, rq.r, got)

		default:
			Equals(t, rq.q, got)
		}
	}
}

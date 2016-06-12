package edgeos

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/britannic/testutils"
)

func shuffleArray(slice []string) {
	rand.Seed(time.Now().UnixNano())
	n := len(slice)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func TestDiffArray(t *testing.T) {
	biggest := sort.StringSlice{"one", "two", "three", "four", "five", "six"}
	smallest := sort.StringSlice{"one", "two", "three"}
	want := sort.StringSlice{"five", "four", "six"}

	got := DiffArray(biggest, smallest)
	Equals(t, want, got)

	got = DiffArray(smallest, biggest)
	Equals(t, want, got)

	shuffleArray(biggest)
	got = DiffArray(smallest, biggest)
	Equals(t, want, got)

	shuffleArray(smallest)
	got = DiffArray(smallest, biggest)
	Equals(t, want, got)
}

func TestFormatData(t *testing.T) {
	reader := bytes.NewBufferString(tdata.Cfg)
	c, err := ReadCfg(reader)
	OK(t, err)
	NewParms(c).SetOpt(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Nodes([]string{"domains", "hosts"}),
	)

	for _, node := range c.Parms.Nodes {
		var (
			got       io.Reader
			gotList   = make(List)
			lines     []string
			wantBytes []byte
		)
		eq := getSeparator(node)

		getBytes := func() io.Reader {
			sort.Strings(c.bNodes[node].inc)
			return bytes.NewBuffer([]byte(strings.Join(c.bNodes[node].inc, "\n")))
		}

		b := bufio.NewScanner(getBytes())

		for b.Scan() {
			k := b.Text()
			lines = append(lines, fmt.Sprintf("address=%v%v/%v", eq, k, c.Get(node).ip)+"\n")
			gotList[k] = 0
		}

		sort.Strings(lines)
		wantBytes = []byte(strings.Join(lines, ""))

		fmttr := "address=" + eq + "%v/" + c.Get(node).ip
		got = formatData(fmttr, gotList)
		gotBytes, err := ioutil.ReadAll(got)
		OK(t, err)
		Equals(t, wantBytes[:], gotBytes[:])
		// fmt.Println(string(gotBytes[:]))
	}
}

func TestGetSubdomains(t *testing.T) {
	d := getSubdomains("top.one.two.three.four.five.six.intellitxt.com")

	for key := range d {
		Assert(t, d.keyExists(key), fmt.Sprintf("%v key doesn't exist", key), d)
	}
}

func TestGetType(t *testing.T) {
	tests := []struct {
		ntypestr string
		typeint  ntype
		typestr  string
	}{
		{typeint: domain, typestr: domains, ntypestr: "domain"},
		{typeint: host, typestr: hosts, ntypestr: "host"},
		{typeint: pre, typestr: preConf, ntypestr: "pre"},
		{typeint: root, typestr: blacklist, ntypestr: "root"},
		{typeint: unknown, typestr: notknown, ntypestr: "unknown"},
		{typeint: zone, typestr: zones, ntypestr: "zone"},
		{typeint: 100, typestr: notknown, ntypestr: "ntype(100)"},
	}

	for _, test := range tests {
		if test.typeint != 100 {
			Equals(t, test.typeint, typeStr(test.typestr))
			Equals(t, test.typestr, typeInt(test.typeint))
			Equals(t, test.typestr, getType(test.typeint))
			Equals(t, test.typeint, getType(test.typestr))
		}
		Equals(t, test.ntypestr, fmt.Sprint(test.typeint))
	}

}

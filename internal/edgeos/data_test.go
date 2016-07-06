package edgeos

import (
	"bufio"
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
	c := NewConfig(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Nodes([]string{domains, hosts}),
	)

	l := &CFGstatic{Cfg: tdata.Cfg}
	err := c.ReadCfg(l)
	OK(t, err)

	for _, node := range c.Parms.Nodes {
		var (
			got       io.Reader
			gotList   = make(list)
			lines     []string
			wantBytes []byte
		)
		eq := getSeparator(node)

		getBytes := func() io.Reader {
			sort.Strings(c.tree[node].inc)
			return strings.NewReader(strings.Join(c.tree[node].inc, "\n"))
		}

		b := bufio.NewScanner(getBytes())

		for b.Scan() {
			k := b.Text()
			lines = append(lines, fmt.Sprintf("address=%v%v/%v", eq, k, c.tree[node].ip)+"\n")
			gotList[k] = 0
		}

		sort.Strings(lines)
		wantBytes = []byte(strings.Join(lines, ""))

		fmttr := "address=" + eq + "%v/" + c.tree[node].ip
		got = formatData(fmttr, gotList)
		gotBytes, err := ioutil.ReadAll(got)
		OK(t, err)
		Equals(t, wantBytes, gotBytes)
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
		{typeint: 100, typestr: notknown, ntypestr: "ntype(100)"},
		{typeint: domn, typestr: domains, ntypestr: "domn"},
		{typeint: excDomn, typestr: ExcDomns, ntypestr: "excDomn"},
		{typeint: excHost, typestr: ExcHosts, ntypestr: "excHost"},
		{typeint: excRoot, typestr: ExcRoots, ntypestr: "excRoot"},
		{typeint: host, typestr: hosts, ntypestr: "host"},
		{typeint: preDomn, typestr: PreDomns, ntypestr: "preDomn"},
		{typeint: preHost, typestr: PreHosts, ntypestr: "preHost"},
		{typeint: root, typestr: blacklist, ntypestr: "root"},
		{typeint: unknown, typestr: notknown, ntypestr: "unknown"},
		{typeint: zone, typestr: zones, ntypestr: "zone"},
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

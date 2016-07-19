package edgeos

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/smartystreets/goconvey/convey"
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
	Convey("Testing DiffArray()", t, func() {
		biggest := sort.StringSlice{"one", "two", "three", "four", "five", "six"}
		smallest := sort.StringSlice{"one", "two", "three"}

		act := DiffArray(biggest, smallest)
		exp := sort.StringSlice{"five", "four", "six"}
		So(act, ShouldResemble, exp)

		act = DiffArray(smallest, biggest)
		So(act, ShouldResemble, exp)

		shuffleArray(biggest)
		act = DiffArray(smallest, biggest)
		So(act, ShouldResemble, exp)

		shuffleArray(smallest)
		act = DiffArray(smallest, biggest)
		So(act, ShouldResemble, exp)
	})
}

func TestFormatData(t *testing.T) {
	Convey("Testing FormatData()", t, func() {
		c := NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
			Nodes([]string{domains, hosts}),
		)

		So(c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)

		for _, node := range c.Parms.Nodes {
			var (
				act      io.Reader
				actList  = list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
				lines    []string
				expBytes []byte
			)

			eq := getSeparator(node)

			r := func() io.Reader {
				sort.Strings(c.tree[node].inc)
				return strings.NewReader(strings.Join(c.tree[node].inc, "\n"))
			}

			b := bufio.NewScanner(r())

			for b.Scan() {
				k := b.Text()
				lines = append(lines, fmt.Sprintf("address=%v%v/%v", eq, k, c.tree[node].ip)+"\n")
				actList.set(k, 0)
			}

			sort.Strings(lines)
			expBytes = []byte(strings.Join(lines, ""))

			fmttr := "address=" + eq + "%v/" + c.tree[node].ip
			act = formatData(fmttr, actList)
			actBytes, err := ioutil.ReadAll(act)

			So(err, ShouldBeNil)
			So(actBytes, ShouldResemble, expBytes)
		}
	})
}

func TestGetSubdomains(t *testing.T) {
	Convey("Testing GetSubdomains()", t, func() {
		d := getSubdomains("top.one.two.three.four.five.six.intellitxt.com")
		d.RWMutex = &sync.RWMutex{}

		for key := range d.entry {
			So(d.keyExists(key), ShouldBeTrue)
		}
	})
}

func TestGetType(t *testing.T) {
	Convey("Testing GetType()", t, func() {
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

		for _, tt := range tests {
			Convey("Testing GetType("+tt.ntypestr+")", func() {
				if tt.typeint != 100 {
					So(typeStr(tt.typestr), ShouldEqual, tt.typeint)
					So(typeInt(tt.typeint), ShouldEqual, tt.typestr)
					So(getType(tt.typeint), ShouldEqual, tt.typestr)
					So(getType(tt.typestr), ShouldEqual, tt.typeint)
				}
				So(fmt.Sprint(tt.typeint), ShouldEqual, tt.ntypestr)
			})
		}
	})
}

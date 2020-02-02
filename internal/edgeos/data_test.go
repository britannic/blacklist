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
	"sync"
	"testing"
	"time"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/smartystreets/goconvey/convey"
)

// logIt writes to io.Writer
func logIt(w io.Writer, s string) {
	io.Copy(w, strings.NewReader(s))
}

func shuffleArray(slice []string) {
	rand.Seed(time.Now().UnixNano())
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func TestDiffArray(t *testing.T) {
	Convey("Testing diffArray()", t, func() {
		biggest := sort.StringSlice{"one", "two", "three", "four", "five", "six"}
		smallest := sort.StringSlice{"one", "two", "three"}
		exp := sort.StringSlice{"five", "four", "six"}

		So(diffArray(biggest, smallest), ShouldResemble, exp)
		So(diffArray(smallest, biggest), ShouldResemble, exp)

		shuffleArray(biggest)
		So(diffArray(smallest, biggest), ShouldResemble, exp)

		shuffleArray(smallest)
		So(diffArray(smallest, biggest), ShouldResemble, exp)
	})
}

func TestFormatData(t *testing.T) {
	Convey("Testing FormatData()", t, func() {
		c := NewConfig(
			Dir("/tmp"),
			Ext("blacklist.conf"),
			Prefix("address=", "server="),
		)

		So(c.Blacklist(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)

		for _, node := range c.sortKeys() {
			var (
				actList = &list{RWMutex: &sync.RWMutex{}, entry: make(entry)}

				o = &source{
					ip: c.tree[node].ip,
					Env: &Env{
						Pfx: dnsPfx{
							domain: "address=",
							host:   "server=",
						},
					},
					nType: domn,
				}
				// pfx      = dnsPfx{domain: "address=", host: "server="}
				fmttr    = getDnsmasqPrefix(o)
				expBytes []byte
				lines    []string
			)

			r := func() io.Reader {
				sort.Strings(c.tree[node].inc)
				return strings.NewReader(strings.Join(c.tree[node].inc, "\n"))
			}

			b := bufio.NewScanner(r())
			for b.Scan() {
				k := b.Text()
				lines = append(lines, fmt.Sprintf(fmttr, k)+"\n")
				actList.set([]byte(k))
			}

			sort.Strings(lines)
			expBytes = []byte(strings.Join(lines, ""))
			actBytes, err := ioutil.ReadAll(formatData(fmttr, actList))

			So(err, ShouldBeNil)
			So(actBytes, ShouldResemble, expBytes)
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
			{typeint: root, typestr: rootNode, ntypestr: "root"},
			{typeint: unknown, typestr: notknown, ntypestr: "unknown"},
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

func TestNewWriter(t *testing.T) {
	Convey("Testing newWriter()", t, func() {

		tests := []struct {
			name   string
			exp    io.Writer
			expStr string
		}{
			{
				name: "vanilla",
				exp: func() io.Writer {
					var b bytes.Buffer
					return bufio.NewWriter(&b)
				}(),
				expStr: "Es ist Krieg!",
			},
		}
		for _, tt := range tests {
			act := NewWriter()
			Convey("Testing "+tt.name, func() {
				So(act, ShouldResemble, tt.exp)
				logIt(act, tt.expStr)
				var b bytes.Buffer
				want := bufio.NewWriter(&b)
				io.Copy(want, strings.NewReader(tt.expStr))
				So(act, ShouldResemble, want)
			})
		}
	})
}

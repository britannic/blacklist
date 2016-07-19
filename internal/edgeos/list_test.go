package edgeos

import (
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKeys(t *testing.T) {
	Convey("Testing Keys()", t, func() {
		var keys sort.StringSlice
		c := NewConfig(Nodes([]string{"domains", "hosts"}))
		So(c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)

		So(c.sortKeys(), ShouldResemble, sort.StringSlice{"blacklist", "domains", "hosts"})

		So(c.GetAll().Names(), ShouldResemble, sort.StringSlice{"includes.[1]", "includes.[9]", "malc0de", "malwaredomains.com", "openphish", "raw.github.com", "simple_tracking", "sysctl.org", "tasty", "volkerschatz", "yoyo", "zeus"})

		for _, k := range []string{"a", "b", "c", "z", "q", "s", "e", "i", "x", "m"} {
			keys = append(keys, k)
		}

		So(keys.Len(), ShouldEqual, 10)
	})
}

func TestKeyExists(t *testing.T) {
	Convey("Testing KeyExists()", t, func() {
		full := "top.one.two.three.four.five.six.intellitxt.com"
		d := getSubdomains(full)
		d.RWMutex = &sync.RWMutex{}

		for key := range d.entry {
			So(d.keyExists(key), ShouldBeTrue)
		}

		So(d.keyExists(`zKeyDoesn'tExist`), ShouldBeFalse)
	})
}

func TestMergeList(t *testing.T) {
	Convey("Testing MergeList()", t, func() {
		testList1 := list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
		testList2 := list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
		exp := list{RWMutex: &sync.RWMutex{}, entry: make(entry)}

		for i := 0; i < 20; i++ {
			exp.entry[string(i)] = 1
			switch {
			case i%2 == 0:
				testList1.entry[string(i)] = 1
			case i%2 != 0:
				testList2.entry[string(i)] = 1
			}
		}

		So(mergeList(testList1, testList2), ShouldResemble, exp)
	})
}

func TestString(t *testing.T) {
	Convey("Testing String()", t, func() {
		exp := `"a.applovin.com":0,
"a.glcdn.co":0,
"a.vserv.mobi":0,
"ad.leadboltapps.net":0,
"ad.madvertise.de":0,
"ad.where.com":0,
"ad1.adinfuse.com":0,
"ad2.adinfuse.com":0,
"adcontent.saymedia.com":0,
"adinfuse.com":0,
"admicro1.vcmedia.vn":0,
"admicro2.vcmedia.vn":0,
"admin.vserv.mobi":0,
"ads.adiquity.com":0,
"ads.admarvel.com":0,
"ads.admoda.com":0,
"ads.celtra.com":0,
"ads.flurry.com":0,
"ads.matomymobile.com":0,
"ads.mobgold.com":0,
"ads.mobilityware.com":0,
"ads.mopub.com":0,
`
		So(act.String(), ShouldEqual, exp)
	})
}

func TestSubKeyExists(t *testing.T) {
	Convey("Testing SubKeyExists()", t, func() {
		full := "top.one.two.three.four.five.six.com"
		d := getSubdomains(full)
		d.RWMutex = &sync.RWMutex{}

		key := `intellitxt.com`
		d.set(key, 0)
		So(d.subKeyExists(key), ShouldBeTrue)

		act := len(d.entry)
		exp := strings.Count(full, ".") + 1
		So(act, ShouldEqual, exp)

		for key = range d.entry {
			So(d.subKeyExists(key), ShouldBeTrue)
		}

		So(d.subKeyExists(`zKeyDoesn'tExist`), ShouldBeFalse)
	})
}

var (
	act = list{entry: entry{
		"a.applovin.com":         0,
		"a.glcdn.co":             0,
		"a.vserv.mobi":           0,
		"ad.leadboltapps.net":    0,
		"ad.madvertise.de":       0,
		"ad.where.com":           0,
		"ad1.adinfuse.com":       0,
		"ad2.adinfuse.com":       0,
		"adcontent.saymedia.com": 0,
		"adinfuse.com":           0,
		"admicro1.vcmedia.vn":    0,
		"admicro2.vcmedia.vn":    0,
		"admin.vserv.mobi":       0,
		"ads.adiquity.com":       0,
		"ads.admarvel.com":       0,
		"ads.admoda.com":         0,
		"ads.celtra.com":         0,
		"ads.flurry.com":         0,
		"ads.matomymobile.com":   0,
		"ads.mobgold.com":        0,
		"ads.mobilityware.com":   0,
		"ads.mopub.com":          0,
	},
	}
)

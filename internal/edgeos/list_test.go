package edgeos

import (
	"bytes"
	"sort"
	"sync"
	"testing"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKeys(t *testing.T) {
	Convey("Testing Keys()", t, func() {
		var keys sort.StringSlice
		c := NewConfig()
		So(c.ReadCfg(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)

		So(c.sortKeys(), ShouldResemble, sort.StringSlice{"blacklist", "domains", "hosts"})

		So(c.GetAll().Names(), ShouldResemble, sort.StringSlice{"blacklisted-servers", "blacklisted-subdomains", "malc0de", "malwaredomains.com", "openphish", "raw.github.com", "simple_tracking", "sysctl.org", "tasty", "volkerschatz", "yoyo", "zeus"})

		for _, k := range []string{"a", "b", "c", "z", "q", "s", "e", "i", "x", "m"} {
			keys = append(keys, k)
		}

		So(keys.Len(), ShouldEqual, 10)
	})
}

func TestKeyExists(t *testing.T) {
	Convey("Testing KeyExists()", t, func() {
		full := []byte("top.one.two.three.four.five.six.intellitxt.com")
		d := getSubdomains(full)
		d.RWMutex = &sync.RWMutex{}

		for k := range d.entry {
			So(d.keyExists([]byte(k)), ShouldBeTrue)
		}

		So(d.keyExists([]byte("zKeyDoesn'tExist")), ShouldBeFalse)
	})
}

func TestMergeList(t *testing.T) {
	Convey("Testing MergeList()", t, func() {
		testList1 := list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
		testList2 := list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
		exp := list{RWMutex: &sync.RWMutex{}, entry: make(entry)}

		for i := range Iter(20) {
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
		full := []byte("top.one.two.three.four.five.six.com")
		d := getSubdomains(full)
		d.RWMutex = &sync.RWMutex{}

		k := `intellitxt.com`
		d.set([]byte(k), 0)
		So(d.subKeyExists([]byte(k)), ShouldBeTrue)

		act := len(d.entry)
		exp := bytes.Count(full, []byte(".")) + 1
		So(act, ShouldEqual, exp)

		for k = range d.entry {
			So(d.subKeyExists([]byte(k)), ShouldBeTrue)
		}

		So(d.subKeyExists([]byte(`zKeyDoesn'tExist`)), ShouldBeFalse)
	})
}

var (
	act = list{
		entry: entry{
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

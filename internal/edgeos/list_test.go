package edgeos

import (
	"fmt"
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
		So(c.Blacklist(&CFGstatic{Cfg: tdata.Cfg}), ShouldBeNil)

		So(c.sortKeys(), ShouldResemble, sort.StringSlice{"blacklist", "domains", "hosts"})

		So(c.GetAll().Names(), ShouldResemble, sort.StringSlice{"blacklisted-servers", "blacklisted-subdomains", "global-blacklisted-domains", "malc0de", "malwaredomains.com", "openphish", "raw.github.com", "simple_tracking", "sysctl.org", "tasty", "volkerschatz", "yoyo", "zeus"})

		for _, k := range []string{"a", "b", "c", "z", "q", "s", "e", "i", "x", "m"} {
			keys = append(keys, k)
		}

		So(keys.Len(), ShouldEqual, 10)
	})
}

func TestKeyExists(t *testing.T) {
	exp := list{
		RWMutex: &sync.RWMutex{},
		entry: entry{"five.six.intellitxt.com": struct{}{},
			"four.five.six.intellitxt.com":                   struct{}{},
			"intellitxt.com":                                 struct{}{},
			"one.two.three.four.five.six.intellitxt.com":     struct{}{},
			"six.intellitxt.com":                             struct{}{},
			"three.four.five.six.intellitxt.com":             struct{}{},
			"top.one.two.three.four.five.six.intellitxt.com": struct{}{},
			"two.three.four.five.six.intellitxt.com":         struct{}{},
		},
	}
	Convey("Testing KeyExists()", t, func() {
		for _, k := range keyArray {
			So(exp.keyExists([]byte(k)), ShouldBeTrue)
		}
		So(exp.keyExists([]byte("zKeyDoesn'tExist")), ShouldBeFalse)
	})
}

func TestSubKeyExists(t *testing.T) {
	exp := list{
		RWMutex: &sync.RWMutex{},
		entry: entry{"five.six.intellitxt.com": struct{}{},
			"four.five.six.intellitxt.com":                   struct{}{},
			"intellitxt.com":                                 struct{}{},
			"one.two.three.four.five.six.intellitxt.com":     struct{}{},
			"six.intellitxt.com":                             struct{}{},
			"three.four.five.six.intellitxt.com":             struct{}{},
			"top.one.two.three.four.five.six.intellitxt.com": struct{}{},
			"two.three.four.five.six.intellitxt.com":         struct{}{},
		},
	}
	Convey("Testing KeyExists()", t, func() {
		for _, k := range keyArray {
			So(exp.subKeyExists([]byte(k)), ShouldBeTrue)
		}
		So(exp.subKeyExists([]byte("zKeyDoesn'tExist")), ShouldBeFalse)
		So(exp.subKeyExists([]byte("com")), ShouldBeFalse)
	})
}

func TestMerge(t *testing.T) {
	Convey("Testing merge()", t, func() {
		testList1 := list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
		testList2 := list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
		exp := list{RWMutex: &sync.RWMutex{}, entry: make(entry)}

		for i := range Iter(20) {
			exp.entry[fmt.Sprint(i)] = struct{}{}
			switch {
			case i%2 == 0:
				testList1.entry[fmt.Sprint(i)] = struct{}{}
			case i%2 != 0:
				testList2.entry[fmt.Sprint(i)] = struct{}{}
			}
		}
		testList1.merge(&testList2)

		So(testList1, ShouldResemble, exp)
	})
}

func TestString(t *testing.T) {
	Convey("Testing String()", t, func() {
		exp := `"a.applovin.com":{},
"a.glcdn.co":{},
"a.vserv.mobi":{},
"ad.leadboltapps.net":{},
"ad.madvertise.de":{},
"ad.where.com":{},
"ad1.adinfuse.com":{},
"ad2.adinfuse.com":{},
"adcontent.saymedia.com":{},
"adinfuse.com":{},
"admicro1.vcmedia.vn":{},
"admicro2.vcmedia.vn":{},
"admin.vserv.mobi":{},
"ads.adiquity.com":{},
"ads.admarvel.com":{},
"ads.admoda.com":{},
"ads.celtra.com":{},
"ads.flurry.com":{},
"ads.matomymobile.com":{},
"ads.mobgold.com":{},
"ads.mobilityware.com":{},
"ads.mopub.com":{},
`
		So(act.String(), ShouldEqual, exp)
	})
}

var (
	act = list{
		entry: entry{
			"a.applovin.com":         struct{}{},
			"a.glcdn.co":             struct{}{},
			"a.vserv.mobi":           struct{}{},
			"ad.leadboltapps.net":    struct{}{},
			"ad.madvertise.de":       struct{}{},
			"ad.where.com":           struct{}{},
			"ad1.adinfuse.com":       struct{}{},
			"ad2.adinfuse.com":       struct{}{},
			"adcontent.saymedia.com": struct{}{},
			"adinfuse.com":           struct{}{},
			"admicro1.vcmedia.vn":    struct{}{},
			"admicro2.vcmedia.vn":    struct{}{},
			"admin.vserv.mobi":       struct{}{},
			"ads.adiquity.com":       struct{}{},
			"ads.admarvel.com":       struct{}{},
			"ads.admoda.com":         struct{}{},
			"ads.celtra.com":         struct{}{},
			"ads.flurry.com":         struct{}{},
			"ads.matomymobile.com":   struct{}{},
			"ads.mobgold.com":        struct{}{},
			"ads.mobilityware.com":   struct{}{},
			"ads.mopub.com":          struct{}{},
		},
	}
	keyArray = [][]byte{
		[]byte("top.one.two.three.four.five.six.intellitxt.com"),
		[]byte("one.two.three.four.five.six.intellitxt.com"),
		[]byte("two.three.four.five.six.intellitxt.com"),
		[]byte("three.four.five.six.intellitxt.com"),
		[]byte("four.five.six.intellitxt.com"),
		[]byte("five.six.intellitxt.com"),
		[]byte("six.intellitxt.com"),
		[]byte("intellitxt.com"),
	}
)

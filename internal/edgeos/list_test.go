package edgeos

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/britannic/testutils"
)

func TestKeys(t *testing.T) {
	var keys sort.StringSlice
	l := &CFGstatic{Cfg: tdata.Cfg}
	c := NewConfig(Nodes([]string{"domains", "hosts"}))
	err := c.ReadCfg(l)
	OK(t, err)

	z := c.GetAll().Names()
	Equals(t, "[blacklist domains hosts]", fmt.Sprint(c.sortKeys()))
	Equals(t, "[includes.[1] includes.[9] malc0de malwaredomains.com openphish raw.github.com simple_tracking sysctl.org tasty volkerschatz yoyo zeus]", fmt.Sprint(z))

	for _, k := range []string{"a", "b", "c", "z", "q", "s", "e", "i", "x", "m"} {
		keys = append(keys, k)
	}

	Equals(t, 10, keys.Len())
}

func TestKeyExists(t *testing.T) {
	full := "top.one.two.three.four.five.six.intellitxt.com"
	d := getSubdomains(full)
	d.RWMutex = &sync.RWMutex{}

	for key := range d.entry {
		Assert(t, d.keyExists(key), fmt.Sprintf("%v key doesn't exist", key))
	}

	key := `zKeyDoesn'tExist`
	Assert(t, !d.keyExists(key), fmt.Sprintf("%v key shouldn't exist", key))
}

func TestMergeList(t *testing.T) {
	testList1 := list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
	testList2 := list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
	want := list{RWMutex: &sync.RWMutex{}, entry: make(entry)}

	for i := 0; i < 20; i++ {
		want.entry[string(i)] = 1
		switch {
		case i%2 == 0:
			testList1.entry[string(i)] = 1
		case i%2 != 0:
			testList2.entry[string(i)] = 1
		}
	}
	got := mergeList(testList1, testList2)
	Equals(t, want, got)
}

func TestString(t *testing.T) {
	want := `"a.applovin.com":0,
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

	Equals(t, want, gotKeys.String())
}

func TestSubKeyExists(t *testing.T) {
	full := "top.one.two.three.four.five.six.com"
	d := getSubdomains(full)
	d.RWMutex = &sync.RWMutex{}

	key := `intellitxt.com`
	d.set(key, 0)
	got := len(d.entry)
	want := strings.Count(full, ".") + 1

	Equals(t, want, got)

	for key = range d.entry {
		Assert(t, d.subKeyExists(key), fmt.Sprintf("%v sub key doesn't exist", key), d)
	}

	Assert(t, d.subKeyExists(key), fmt.Sprintf("%v key should exist", key))

	key = `zKeyDoesn'tExist`
	Assert(t, !d.subKeyExists(key), fmt.Sprintf("%v sub key shouldn't exist", key))

}

var (
	gotKeys = list{entry: entry{
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

package edgeos

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/britannic/testutils"
)

func TestKeys(t *testing.T) {
	var keys sort.StringSlice
	l := &CFGstatic{Cfg: tdata.Cfg}
	b, err := ReadCfg(l)
	OK(t, err)

	Equals(t, "[blacklist domains hosts]", fmt.Sprint(b.sortKeys()))
	Equals(t, "[adaway malwaredomainlist openphish someonewhocares tasty volkerschatz winhelp2002 yoyo]", fmt.Sprint(b.sortSKeys("hosts")))

	for _, k := range []string{"a", "b", "c", "z", "q", "s", "e", "i", "x", "m"} {
		keys = append(keys, k)
	}

	Equals(t, 10, keys.Len())
}

func TestKeyExists(t *testing.T) {
	full := "top.one.two.three.four.five.six.intellitxt.com"
	d := getSubdomains(full)
	for key := range d {
		Assert(t, d.keyExists(key), fmt.Sprintf("%v key doesn't exist", key))
	}

	key := `zKeyDoesn'tExist`
	Assert(t, !d.keyExists(key), fmt.Sprintf("%v key shouldn't exist", key))
}

func TestMergeList(t *testing.T) {
	testList1 := make(List)
	testList2 := make(List)
	want := make(List)

	for i := 0; i < 20; i++ {
		want[string(i)] = 1
		switch {
		case i%2 == 0:
			testList1[string(i)] = 1
		case i%2 != 0:
			testList2[string(i)] = 1
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
	key := `intellitxt.com`
	d[key] = 0
	got := len(d)
	want := strings.Count(full, ".") + 1

	Equals(t, want, got)

	for key = range d {
		Assert(t, d.subKeyExists(key), fmt.Sprintf("%v sub key doesn't exist", key), d)
	}

	Assert(t, d.subKeyExists(key), fmt.Sprintf("%v key should exist", key))

	key = `zKeyDoesn'tExist`
	Assert(t, !d.subKeyExists(key), fmt.Sprintf("%v sub key shouldn't exist", key))

}

var (
	gotKeys = List{
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
	}
)

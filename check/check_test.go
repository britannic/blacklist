package check_test

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/britannic/blacklist/check"
	"github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/data"
	"github.com/britannic/blacklist/global"
	. "github.com/britannic/testutils"
)

var (
	blacklist        *config.Blacklist
	live             = &check.Cfg{Blacklist: blacklist}
	dmsqdir, logfile string
)

func init() {
	switch global.WhatArch {
	case global.TargetArch:
		dmsqdir = global.DmsqDir

	default:
		dmsqdir = "../testdata"
	}

	var err error
	live.Blacklist, err = config.Get(config.Testdata, global.Area.Root)

	if err != nil {
		log.Fatal("Couldn't load config.Testdata")
	}
}

func TestBlacklistings(t *testing.T) {
	a := &check.Args{
		Fname: dmsqdir + "/%v" + ".pre-configured" + global.Fext,
	}

	Assert(t, live.Blacklistings(a), "Blacklistings failed.", a)

	badData := `blacklist {
        disabled false
        dns-redirect-ip 0.0.0.0
        domains {
            include broken.adsrvr.org
            include broken.adtechus.net
            include broken.advertising.com
            include broken.centade.com
            include broken.doubleclick.net
            include broken.free-counter.co.uk
            include broken.intellitxt.com
            include broken.kiosked.com
        }
        hosts {
            include broken.beap.gemini.yahoo.com
						include broken.beap.gemini.msn.com
        }
    }`

	var err error
	failed := &check.Cfg{Blacklist: blacklist}
	failed.Blacklist, err = config.Get(badData, global.Area.Root)
	OK(t, err)

	Assert(t, !failed.Blacklistings(a), "Blacklistings should have failed.", failed)

	a = &check.Args{
		Fname: dmsqdir + "/%v" + ".--BROKEN--" + global.Fext,
	}

	Assert(t, !live.Blacklistings(a), "Blacklistings should have failed.", a)
}

func TestExclusions(t *testing.T) {

	a := &check.Args{
		Ex:  data.GetExcludes(*live.Blacklist),
		Dir: dmsqdir,
	}

	Assert(t, live.Exclusions(a), "Exclusions failure.", a)

	a.Dir = "broken directory"
	Assert(t, !live.Exclusions(a), "Exclusions should have failed.", a)

	a.Ex = make(config.Dict)
	b := *live.Blacklist

	badexcludes := b[global.Area.Domains].Include
	for _, k := range badexcludes {
		a.Ex[k] = 0
	}

	Assert(t, !live.Exclusions(a), "Exclusions should have failed.", a)
}

func TestExcludedDomains(t *testing.T) {
	a := &check.Args{
		Ex:  data.GetExcludes(*live.Blacklist),
		Dex: make(config.Dict),
		Dir: dmsqdir,
	}

	Assert(t, live.ExcludedDomains(a), "Excluded domains failure.", a)
}

func TestConfFiles(t *testing.T) {
	a := &check.Args{
		Dir:   dmsqdir,
		Fname: dmsqdir + `/*` + global.Fext,
	}

	Assert(t, live.ConfFiles(a), "Problems with dnsmasq configuration files.", a)
}

func TestConfFilesContent(t *testing.T) {
	a := &check.Args{
		Dir:   dmsqdir,
		Fname: dmsqdir + `/*` + global.Fext,
	}
	Assert(t, live.ConfFilesContent(a), "Problems with dnsmasq configuration files.", a)
}

func TestConfIP(t *testing.T) {
	a := &check.Args{
		Dir: dmsqdir,
	}

	Assert(t, live.ConfIP(a), "DNS redirect IP configuration failed", a)
}

func TestConfTemplates(t *testing.T) {
	a := &check.Args{
		Data: fileManifest,
		Dir:  `../payload/blacklist`,
	}

	Assert(t, check.ConfTemplates(a), "Configuration template nodes do not match", a)
}

// func TestIsDisabled(t *testing.T) {
// 	a := &Args{
// 		Dir:   dmsqdir,
// 		Fname: dmsqdir + `/*` + global.Fext,
// 	}
//
// }

// func TestIPRedirection(t *testing.T) {
// 	a := &check.Args{
// 		Dir: dmsqdir,
// 	}
// 	if live.IPRedirection(a) != nil {
// 		t.Errorf("Problems with IP redirection: %v", err)
// 	}
// }

var fileManifest = `../payload/blacklist
../payload/blacklist/disabled
../payload/blacklist/disabled/node.def
../payload/blacklist/dns-redirect-ip
../payload/blacklist/dns-redirect-ip/node.def
../payload/blacklist/domains
../payload/blacklist/domains/dns-redirect-ip
../payload/blacklist/domains/dns-redirect-ip/node.def
../payload/blacklist/domains/exclude
../payload/blacklist/domains/exclude/node.def
../payload/blacklist/domains/include
../payload/blacklist/domains/include/node.def
../payload/blacklist/domains/node.def
../payload/blacklist/domains/source
../payload/blacklist/domains/source/node.def
../payload/blacklist/domains/source/node.tag
../payload/blacklist/domains/source/node.tag/description
../payload/blacklist/domains/source/node.tag/description/node.def
../payload/blacklist/domains/source/node.tag/prefix
../payload/blacklist/domains/source/node.tag/prefix/node.def
../payload/blacklist/domains/source/node.tag/url
../payload/blacklist/domains/source/node.tag/url/node.def
../payload/blacklist/exclude
../payload/blacklist/exclude/node.def
../payload/blacklist/hosts
../payload/blacklist/hosts/dns-redirect-ip
../payload/blacklist/hosts/dns-redirect-ip/node.def
../payload/blacklist/hosts/exclude
../payload/blacklist/hosts/exclude/node.def
../payload/blacklist/hosts/include
../payload/blacklist/hosts/include/node.def
../payload/blacklist/hosts/node.def
../payload/blacklist/hosts/source
../payload/blacklist/hosts/source/node.def
../payload/blacklist/hosts/source/node.tag
../payload/blacklist/hosts/source/node.tag/description
../payload/blacklist/hosts/source/node.tag/description/node.def
../payload/blacklist/hosts/source/node.tag/prefix
../payload/blacklist/hosts/source/node.tag/prefix/node.def
../payload/blacklist/hosts/source/node.tag/url
../payload/blacklist/hosts/source/node.tag/url/node.def
../payload/blacklist/node.def
`

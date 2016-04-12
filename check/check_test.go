package check_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	. "github.com/britannic/blacklist/check"
	"github.com/britannic/blacklist/config"
	"github.com/britannic/blacklist/data"
	"github.com/britannic/blacklist/global"
	"github.com/britannic/blacklist/tdata"
	. "github.com/britannic/testutils"
)

var (
	blacklist *config.Blacklist
	live      = &Cfg{Blacklist: blacklist}
	dmsqdir   string
	log       *logrus.Logger
	hook      *test.Hook
)

func init() {
	global.SetVars(global.WhatArch)
	switch global.WhatArch {
	case global.TargetArch:
		dmsqdir = global.DmsqDir

	default:
		dmsqdir = "../tdata"
	}

	log, hook = test.NewNullLogger()
	var err error

	live.Blacklist, err = config.Get(tdata.Cfg, global.Area.Root)
	if err != nil {
		fmt.Print(fmt.Errorf("Couldn't load tdata.Cfg"))
		os.Exit(1)
	}
}

func TestBlacklistings(t *testing.T) {
	a := &Args{
		Fname: dmsqdir + "/%v" + ".pre-configured" + global.Fext,
		Log:   log,
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
	failed := &Cfg{Blacklist: blacklist}
	failed.Blacklist, err = config.Get(badData, global.Area.Root)
	OK(t, err)

	Assert(t, !failed.Blacklistings(a), "Blacklistings should have failed.", failed)

	a.Fname = dmsqdir + "/%v" + ".--BROKEN--" + global.Fext

	badData = `blacklist {
        disabled false
        dns-redirect-ip 0.0.0.0
        hosts {
            include broken.beap.gemini.yahoo.com
						include broken.beap.gemini.msn.com
        }
    }`

	failed = &Cfg{Blacklist: blacklist}
	failed.Blacklist, err = config.Get(badData, global.Area.Root)
	OK(t, err)

	Assert(t, !failed.Blacklistings(a), "Blacklistings should have failed.", a)

	Equals(t, "Includes not correct in ../tdata/hosts.--BROKEN--.blacklist.conf\n\tGot: []\n\tWant: [broken.beap.gemini.yahoo.com broken.beap.gemini.msn.com]", hook.LastEntry().Message)
	Equals(t, logrus.ErrorLevel.String(), hook.LastEntry().Level.String())
}

func TestExclusions(t *testing.T) {

	a := &Args{
		Ex:  data.GetExcludes(*live.Blacklist),
		Dir: dmsqdir,
		Log: log,
	}

	Assert(t, live.Exclusions(a), fmt.Sprintf("Exclusions failure - last log entry: %v", hook.LastEntry().Message))

	a.Dir = "broken directory"
	Assert(t, !live.Exclusions(a), fmt.Sprintf("Exclusions should have failed - last log entry: %v", hook.LastEntry().Message))

	a.Ex = make(config.Dict)
	b := *live.Blacklist

	badexcludes := b[global.Area.Domains].Include
	for _, k := range badexcludes {
		a.Ex[k] = 0
	}

	Assert(t, !live.Exclusions(a), "Exclusions should have failed.", a)
}

func TestExcludedDomains(t *testing.T) {
	a := &Args{
		Ex:  data.GetExcludes(*live.Blacklist),
		Dex: make(config.Dict),
		Dir: dmsqdir,
		Log: log,
	}

	Assert(t, live.ExcludedDomains(a), fmt.Sprintf("Excluded domains failure - last log entry: %v", hook.LastEntry().Message))

	a.Dir = "--BROKEN--"
	a.Ex = make(config.Dict)
	a.Dex = make(config.Dict)

	Assert(t, !live.ExcludedDomains(a), fmt.Sprintf("Excluded domains failure - last log entry: %v", hook.LastEntry().Message))

	Equals(t, "Error getting file: --BROKEN--/domains.pre-configured.blacklist.conf, error: open --BROKEN--/domains.pre-configured.blacklist.conf: no such file or directory\n", hook.LastEntry().Message)
	Equals(t, logrus.ErrorLevel.String(), hook.LastEntry().Level.String())

	var (
		err     error
		badData = `blacklist {
				disabled false
				dns-redirect-ip 0.0.0.0
				hosts {
						include broken.beap.gemini.yahoo.com
						include broken.beap.gemini.msn.com
				}
		}`
	)

	failed := &Cfg{Blacklist: blacklist}
	failed.Blacklist, err = config.Get(badData, global.Area.Root)
	OK(t, err)

	Assert(t, !failed.ExcludedDomains(a), fmt.Sprintf("Excluded domains failure - last log entry: %v", hook.LastEntry().Message))

	Equals(t, "Error getting file: --BROKEN--/hosts.pre-configured.blacklist.conf, error: open --BROKEN--/hosts.pre-configured.blacklist.conf: no such file or directory\n", hook.LastEntry().Message)
	Equals(t, logrus.ErrorLevel.String(), hook.LastEntry().Level.String())
}

func TestConfFiles(t *testing.T) {
	a := &Args{
		Dir:   dmsqdir,
		Fname: dmsqdir + `/*` + global.Fext,
		Log:   log,
	}

	Assert(t, live.ConfFiles(a), "Problems with dnsmasq configuration files.", a)
}

func TestConfFilesContent(t *testing.T) {
	a := &Args{
		Dir:   dmsqdir,
		Fname: dmsqdir + `/*` + global.Fext,
		Log:   log,
	}
	Assert(t, live.ConfFilesContent(a), "Problems with dnsmasq configuration files.", a)
}

func TestConfIP(t *testing.T) {
	a := &Args{
		Dir: dmsqdir,
		Log: log,
	}

	Assert(t, live.ConfIP(a), "DNS redirect IP configuration failed", a)
}

func TestConfTemplates(t *testing.T) {
	a := &Args{
		Data: tdata.FileManifest,
		Dir:  `../payload/blacklist`,
		Log:  log,
	}

	Assert(t, ConfTemplates(a), "Configuration template nodes do not match", a)
}

// func TestIsDisabled(t *testing.T) {
// 	a := &Args{
// 		Dir:   dmsqdir,
// 		Fname: dmsqdir + `/*` + global.Fext,
// 	}
//
// }

func TestIPRedirection(t *testing.T) {
	if global.WhatArch != global.TargetArch {
		t.SkipNow()
	}

	a := &Args{
		Dir: dmsqdir,
		Log: log,
	}

	if !live.IPRedirection(a) {
		t.Errorf("Problems with IP redirection!")
	}
}

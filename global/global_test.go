package global_test

import (
	"testing"

	"github.com/britannic/blacklist/global"
)

func TestGlobalVars(t *testing.T) {
	type globals struct {
		debug                             bool
		dmsqdir, fsfx, fstr, root, whatOS string
	}

	glob := &globals{
		debug:   global.Dbg,
		dmsqdir: global.DmsqDir,
		fsfx:    global.Fext,
		fstr:    global.FStr,
		root:    global.Root,
		whatOS:  global.WhatOS,
	}

	switch {
	case glob.debug:
		t.Errorf("%+v shouldn't be %v", global.Dbg, glob.debug)
	case glob.dmsqdir != "/etc/dnsmasq.d":
		if global.WhatOS != "darwin" {
			t.Errorf(`%+v should be = "/etc/dnsmasq.d"  not %v`, global.DmsqDir, glob.dmsqdir)
		}
	case glob.fsfx != ".blacklist.conf":
		t.Errorf(`%+v should be = ".blacklist.conf"  not %v`, global.FStr, glob.fsfx)
	case glob.fstr != `%v/%v.%v.blacklist.conf`:
		t.Errorf(`%+v should be = %q not %v`, global.FStr, `%v/%v.%v.blacklist.conf`, glob.fstr)
	case glob.root != "blacklist":
		t.Errorf(`%+v should be = "blacklist"  not %v`, global.Root, glob.root)
	}
}

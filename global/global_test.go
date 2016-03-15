package global_test

import (
	"testing"

	"github.com/britannic/blacklist/global"
)

func TestGlobalVars(t *testing.T) {
	type globals struct {
		debug                     bool
		dmsqdir, fsfx, fstr, root string
	}

	glob := &globals{
		debug:   global.Dbg,
		dmsqdir: global.DmsqDir,
		fsfx:    global.Fext,
		fstr:    global.FStr,
		root:    global.Root,
	}

	switch {
	case glob.debug:
		t.Errorf("%+v shouldn't be %v", glob.debug, glob.debug)
	case glob.dmsqdir != "/etc/dnsmasq.d":
		t.Errorf(`%+v should be = "/etc/dnsmasq.d"  not %v`, glob.dmsqdir, glob.dmsqdir)
	case glob.fsfx != ".blacklist.conf":
		t.Errorf(`%+v should be = ".blacklist.conf"  not %v`, glob.fsfx, glob.fsfx)
	case glob.fstr != `%v/%v.%v.blacklist.conf`:
		t.Errorf(`%+v should be = %q not %v`, glob.fstr, `%v/%v.%v.blacklist.conf`, glob.fstr)
	case glob.root != "blacklist":
		t.Errorf(`%+v should be = "blacklist"  not %v`, glob.root, glob.root)
	}
}

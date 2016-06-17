package edgeos

import (
	"testing"

	"github.com/britannic/blacklist/internal/tdata"
	. "github.com/britannic/testutils"
)

func TestObjectString(t *testing.T) {
	r := &CFGstatic{Cfg: tdata.Cfg}
	exp := "[\nDesc:\t \"Blocking mobile ad providers and some analytics providers\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"adaway\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://adaway.org/hosts.txt\"\n \nDesc:\t \"List of zones serving malicious executables observed by malc0de.com/database/\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"urls\"\nName:\t \"malc0de\"\nnType:\t \"domain\"\nPrefix:\t \"zone \"\nType:\t \"domains\"\nURL:\t \"http://malc0de.com/bl/ZONES\"\n \nDesc:\t \"127.0.0.1 based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"malwaredomainlist\"\nnType:\t \"host\"\nPrefix:\t \"127.0.0.1 \"\nType:\t \"hosts\"\nURL:\t \"http://www.malwaredomainlist.com/hostslist/hosts.txt\"\n \nDesc:\t \"OpenPhish automatic phishing detection\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"openphish\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"https://openphish.com/feed.txt\"\n \nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"domain\"\nPrefix:\t \"\"\nType:\t \"domains\"\nURL:\t \"\"\n \nDesc:\t \"pre-configured\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"pre-configured\"\nName:\t \"pre-configured\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"Zero based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"someonewhocares\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0\"\nType:\t \"hosts\"\nURL:\t \"http://someonewhocares.org/hosts/zero/\"\n \nDesc:\t \"File source\"\nDisabled: false\nFile:\t \"/config/user-data/blist.hosts.src\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"files\"\nName:\t \"tasty\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"\"\n \nDesc:\t \"Ad server blacklists\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"volkerschatz\"\nnType:\t \"host\"\nPrefix:\t \"http\"\nType:\t \"hosts\"\nURL:\t \"http://www.volkerschatz.com/net/adpaths\"\n \nDesc:\t \"Zero based host and domain list\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"0.0.0.0\"\nLtype:\t \"urls\"\nName:\t \"winhelp2002\"\nnType:\t \"host\"\nPrefix:\t \"0.0.0.0 \"\nType:\t \"hosts\"\nURL:\t \"http://winhelp2002.mvps.org/hosts.txt\"\n \nDesc:\t \"Fully Qualified Domain Names only - no prefix to strip\"\nDisabled: false\nFile:\t \"\"\nIP:\t \"192.168.168.1\"\nLtype:\t \"urls\"\nName:\t \"yoyo\"\nnType:\t \"host\"\nPrefix:\t \"\"\nType:\t \"hosts\"\nURL:\t \"http://pgl.yoyo.org/as/serverlist.php?hostformat=nohtml&showintro=1&mimetype=plaintext\"\n]"
	c, err := ReadCfg(r)
	OK(t, err)
	c.Parms = NewParms()
	c.SetOpt(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Nodes([]string{"domains", "hosts"}),
		STypes([]string{"files", "pre-configured", "urls"}),
	)
	Equals(t, exp, c.GetAll().String())
}

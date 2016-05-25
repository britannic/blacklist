package edgeos

import (
	"testing"

	. "github.com/britannic/testutils"
)

func TestOption(t *testing.T) {
	p := &Parms{}
	Equals(t, Parms{debug: false, poll: 0, test: false, verbosity: 0}, *p)
	prev := p.SetOpt(Debug(true), Poll(10), Test(true), Verbosity(2))
	Equals(t, Parms{debug: true, poll: 10, test: true, verbosity: 2}, *p)
	p.SetOpt(prev)
	Equals(t, Parms{debug: false, poll: 0, test: false, verbosity: 0}, *p)
}

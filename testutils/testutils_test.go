package testutils_test

import (
	"fmt"
	"testing"

	. "github.com/britannic/blacklist/testutils"
)

func TestAssert(t *testing.T) {
	True := func() bool {
		return true
	}
	Assert(t, True(), "We have a problem Houston!", true)
}

func TestEquals(t *testing.T) {
	Equals(t, true, true)
}

func TestNotEquals(t *testing.T) {
	NotEquals(t, 10, 0)
}

func TestOK(t *testing.T) {
	OK(t, nil)
}

func TestNotOK(t *testing.T) {
	NotOK(t, fmt.Errorf("We have a problem Houston!"))
}

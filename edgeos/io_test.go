package edgeos

import (
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/britannic/testutils"
)

func TestPurgeFiles(t *testing.T) {
	var (
		dir       = "/tmp"
		ext       = ".delete"
		purgeList []string
		want      error
	)

	for i := 0; i < 10; i++ {
		fname := fmt.Sprintf("%v%v", i, ext)
		f, err := ioutil.TempFile(dir, fname)
		OK(t, err)
		purgeList = append(purgeList, f.Name())
	}

	err := purgeFiles(purgeList)
	OK(t, err)

	got := purgeFiles(purgeList)
	Equals(t, want, got)

	got = purgeFiles([]string{"/dev/null"})
	want = fmt.Errorf(`could not remove "/dev/null"`)
	Equals(t, want, got)
}

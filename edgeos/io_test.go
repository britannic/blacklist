package edgeos

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "github.com/britannic/testutils"
)

func TestListFiles(t *testing.T) {
	want := []string{
		"./test/domains.malc0de.blacklist.conf", "./test/hosts.adaway.blacklist.conf",
	}
	got, err := ListFiles("./test")
	OK(t, err)
	Equals(t, want, got)

	got, err = ListFiles("./ztestz@z")
	NotOK(t, err)
	Equals(t, []string(nil), got)
}

func TestDeleteFile(t *testing.T) {
	f, err := ioutil.TempFile("./test", "test.delete.this.file.txt")
	OK(t, err)
	Equals(t, true, DeleteFile(f.Name()))
	OK(t, err)

	Equals(t, true, DeleteFile(f.Name()))
	OK(t, err)

	Equals(t, false, DeleteFile("/dev/null"))
	OK(t, err)
}

func TestPurgeFiles(t *testing.T) {
	var (
		errArray  []string
		dir       = "/tmp"
		ext       = ".delete"
		purgeList []string
	)

	for i := 0; i < 10; i++ {
		fname := fmt.Sprintf("%v%v", i, ext)
		f, err := ioutil.TempFile(dir, fname)
		OK(t, err)
		purgeList = append(purgeList, f.Name())
	}

	err := PurgeFiles(purgeList)
	OK(t, err)

	got := PurgeFiles(purgeList)
	for _, fname := range purgeList {
		errArray = append(errArray, fmt.Sprintf("%q: stat %v: no such file or directory", fname, fname))
	}
	want := fmt.Errorf("%v", strings.Join(errArray, "\n"))
	Equals(t, want, got)

	got = PurgeFiles([]string{"/dev/null"})
	want = fmt.Errorf(`could not remove "/dev/null"`)
	Equals(t, want, got)
}

func TestWriteFile(t *testing.T) {
	writeFileTests := []struct {
		data  io.Reader
		dir   string
		fname string
		ok    bool
		want  string
	}{
		{
			data:  NewContent("The rest is history!"),
			dir:   "/tmp",
			fname: "Test.util.WriteFile",
			ok:    true,
			want:  "",
		},
		{
			data:  NewContent([]byte{84, 104, 101, 32, 114, 101, 115, 116, 32, 105, 115, 32, 104, 105, 115, 116, 111, 114, 121, 33}),
			dir:   "/tmp",
			fname: "Test.util.WriteFile",
			ok:    true,
			want:  "",
		},
		{
			data:  NewContent("This shouldn't be written!"),
			dir:   "",
			fname: "/tmp/z/d/c/r/c:reallybadfile.zfts",
			ok:    false,
			want:  `unable to open file: /tmp/z/d/c/r/c:reallybadfile.zfts for writing, error: open /tmp/z/d/c/r/c:reallybadfile.zfts: no such file or directory`,
		},
	}

	for _, test := range writeFileTests {
		switch test.ok {
		case true:
			f, err := ioutil.TempFile(test.dir, test.fname)
			OK(t, err)
			err = WriteFile(f.Name(), test.data)
			OK(t, err)
			os.Remove(f.Name())

		default:
			err := WriteFile(test.fname, test.data)
			NotOK(t, err)
			Equals(t, `unable to open file: /tmp/z/d/c/r/c:reallybadfile.zfts for writing, error: open /tmp/z/d/c/r/c:reallybadfile.zfts: no such file or directory`, err.Error())
		}
	}
}

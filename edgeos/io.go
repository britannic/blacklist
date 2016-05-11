package edgeos

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Data inteface implements the Reader and Writer interfaces
type Data interface {
	Reader
	Writer
}

// Reader implements the Read []byte method
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer implements the Write []byte method
type Writer interface {
	Write(p []byte) (n int, err error)
}

// Content struct holds content data
type Content struct {
	read interface{}
	eof  bool
}

// NewContent returns a new Content pointer
func NewContent(toRead interface{}) *Content {
	return &Content{read: toRead, eof: false}
}

func (c *Content) Read(p []byte) (n int, err error) {
	if c.eof {
		return 0, io.EOF
	}

	var data []byte
	switch c.read.(type) {
	case string:
		data = []byte(c.read.(string))
	default:
		data = c.read.([]byte)
	}

	for i, b := range data {
		p[i] = b
	}

	c.eof = true
	return len(data), nil
}

// DeleteFile removes a file if it exists
func DeleteFile(f string) bool {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return true
	}

	if err := os.Remove(f); err != nil {
		return false
	}

	return true
}

// ListFiles returns a list of blacklist files
func ListFiles(dir string) (files []string, err error) {
	dlist, err := ioutil.ReadDir(dir)
	if err != nil {
		return files, err
	}

	for _, f := range dlist {
		if strings.Contains(f.Name(), Fext) {
			files = append(files, dir+"/"+f.Name())
		}
	}
	return files, err
}

// PurgeFiles removes any orphaned blacklist files that don't have sources
func PurgeFiles(files []string) (err error) {
	var errArray []string

NEXT:
	for _, file := range files {
		if _, err = os.Stat(file); os.IsNotExist(err) {
			errArray = append(errArray, fmt.Sprintf("%q: %v", file, err))
			continue NEXT
		}
		if !DeleteFile(file) {
			errArray = append(errArray, fmt.Sprintf("could not remove %q", file))
		}
	}
	switch len(errArray) > 0 {
	case true:
		err = fmt.Errorf("%v", strings.Join(errArray, "\n"))
		return err
	}

	return err
}

// WriteFile writes blacklist data to storage
func WriteFile(fname string, data io.Reader) (err error) {
	var (
		f *os.File
		n int
	)

	f, err = os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("unable to open file: %v for writing, error: %v", fname, err)
	}

	defer func() error {
		if err = f.Close(); err != nil {
			return err
		}
		return nil
	}()

	w := bufio.NewWriter(f)
	buf := make([]byte, 1024)

NEXT:
	for {

		n, err = data.Read(buf)
		switch {
		case err != nil && err != io.EOF:
			return err
		case n == 0:
			break NEXT
		}

		if _, err = w.Write(buf[:n]); err != nil {
			return err
		}

		if err = w.Flush(); err != nil {
			return err
		}
	}

	return err
}

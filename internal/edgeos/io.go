package edgeos

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// CFGcli is for configurations loaded via the EdgeOS CFGcli
type CFGcli struct {
	*Config
	Cfg string
}

// CFGstatic is for configurations loaded via the EdgeOS CFGstatic
type CFGstatic struct {
	*Config
	Cfg string
}

var apiCMDs = map[string]string{
	"cfExists":           "cfExists",
	"cfReturnValue":      "cfReturnValue",
	"cfReturnValues":     "cfReturnValues",
	"echo":               "true",
	"exists":             "exists",
	"existsActive":       "existsActive",
	"getNodeType":        "getNodeType",
	"inSession":          "inSession",
	"isLeaf":             "isLeaf",
	"isMulti":            "isMulti",
	"isTag":              "isTag",
	"listActiveNodes":    "listActiveNodes",
	"listNodes":          "listNodes",
	"returnActiveValue":  "returnActiveValue",
	"returnActiveValues": "returnActiveValues",
	"returnValue":        "returnValue",
	"returnValues":       "returnValues",
	"showCfg":            "showCfg",
	"showConfig":         "showConfig",
}

func active(a string, inS bool) string {
	switch inS {
	case true:
		switch a {
		case "exists":
			a = "existsActive"
		case "listNodes":
			a = "listActiveNodes"
		case "returnValue":
			a = "returnActiveValue"
		case "returnValues":
			a = "returnActiveValues"
		case "showCfg":
			a = "showConfig"

		}
	default:
		switch a {
		case "existsActive":
			a = "exists"
		case "listActiveNodes":
			a = "listNodes"
		case "returnActiveValue":
			a = "returnValue"
		case "returnActiveValues":
			a = "returnValues"
		case "showConfig":
			a = "showCfg"
		}
	}
	return a
}

// apiCMD returns a map of CLI commands
func apiCMD(action string, inCLI bool) string {
	return apiCMDs[active(action, inCLI)]
}

// deleteFile removes a file if it exists
func deleteFile(f string) bool {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return true
	}
	if err := os.Remove(f); err != nil {
		return false
	}
	return true
}

// getFile reads a file and returns a *bufio.Scanner instance
func getFile(fname string) (io.Reader, error) {
	return os.Open(fname)
}

// read returns an EdgeOS config file loaded in to a string
func (c *CFGcli) read() io.Reader {
	b, err := c.load("showConfig", c.Level)
	if err != nil {
		log.Print(err)
	}
	return bytes.NewReader(b)
}

// read returns an EdgeOS config file loaded in to a string
func (c *CFGstatic) read() io.Reader {
	return strings.NewReader(c.Cfg)
}

// purgeFiles removes any orphaned blacklist files that don't have sources
func purgeFiles(files []string) error {
	var errArray []string
NEXT:
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue NEXT
		}
		if !deleteFile(file) {
			errArray = append(errArray, fmt.Sprintf("could not remove %q", file))
		}
	}
	switch len(errArray) > 0 {
	case true:
		return fmt.Errorf("%v", strings.Join(errArray, "\n"))
	}
	return nil
}

// func (c *CFGcli) String() string {
// 	c.ReadCfg(c)
// 	return c.String()
// }
//
// func (c *CFGstatic) String() string {
// 	c.ReadCfg(c)
// 	return c.String()
// }

package edgeos

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// api sets the path and executable for the EdgeOS shell api
const (
	api     = "/bin/cli-shell-api"
	service = "service dns forwarding"
)

// apiCMD returns a map of CLI commands
func apiCMD() (r map[string]string) {
	r = make(map[string]string)
	c := []string{
		"cfExists",
		"cfReturnValue",
		"cfReturnValues",
		"echo",
		"exists",
		"existsActive",
		"getNodeType",
		"inSession",
		"isLeaf",
		"isMulti",
		"isTag",
		"listActiveNodes",
		"listNodes",
		"returnActiveValue",
		"returnActiveValues",
		"returnValue",
		"returnValues",
		"showCfg",
		"showConfig",
	}
	for _, k := range c {
		switch k {
		case "echo":
			r[k] = k
		default:
			r[k] = fmt.Sprintf("%v %v", api, k)
		}
	}
	return r
}

// getFile reads a file and returns a *bufio.Scanner instance
func getFile(fname string) (io.Reader, error) {
	return os.Open(fname)
}

// insession returns true if VyOS/EdgeOS configuration is in session
func insession() bool {
	cmd := exec.Command(api, "inSession")
	var out bytes.Buffer
	cmd.Stdout = &out
	if ok := cmd.Run(); ok == nil {
		if out.String() == "0" {
			return true
		}
	}
	return false
}

// load reads the config using the EdgeOS/VyOS cli-shell-api
func load(action string, level string) (r string, err error) {
	action = shCMD(action)
	x := apiCMD()
	if _, ok := x[action]; !ok {
		return r, fmt.Errorf("API command %v is invalid", level)
	}

	cmd := exec.Command("/bin/bash")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%v %v", x[action], level))

	stdout, err := cmd.Output()
	if err != nil {
		return r, err
	}

	r = string(stdout)
	return r, err
}

// LoadCfg returns an EdgeOS config file string and error
func LoadCfg() (string, error) {
	return load("showConfig", service)
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

// shCMD returns the appropriate command for non-tty or tty context
func shCMD(a string) (action string) {
	if !insession() {
		action = a
		switch a {
		case "listNodes":
			action = "listActiveNodes"
		case "returnValue":
			action = "returnActiveValue"
		case "returnValues":
			action = "returnActiveValues"
		case "exists":
			action = "existsActive"
		case "showConfig":
			action = "showCfg"
		}
	}
	return action
}

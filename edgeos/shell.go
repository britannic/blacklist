package edgeos

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var (
	// API sets the path and executable for the EdgeOS shell API
	API       = "/bin/cli-shell-api"
	inSession = "inSession"
)

// CMDline is an interface for exec.Command, allows for test mock ups
// type CMDline interface {
// 	Command(string, ...string) *exec.Cmd
// }

// APICmd returns a map of CLI commands
func APICmd() (r map[string]string) {
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
		"pwd",
		"returnActiveValue",
		"returnActiveValues",
		"returnValue",
		"returnValues",
		"showCfg",
		"showConfig",
	}

	for _, k := range c {
		switch k {
		case "echo", "pwd":
			r[k] = k
		default:
			r[k] = fmt.Sprintf("%v %v", API, k)
		}
	}
	return r
}

// RunCMDline struct
// type RunCMDline struct{}

// Command is a wrapper for exec.Command
// func (r RunCMDline) Command(command string, args ...string) *exec.Cmd {
// 	out := exec.Command(command, args...)
// 	return out
// }

// Insession returns true if VyOS/EdgeOS configuration is in session
func Insession() bool {
	if ok := os.ExpandEnv("$_OFR_CONFIGURE"); ok == "ok" {
		return true
	}
	return false
}

// Load reads the config using the EdgeOS/VyOS cli-shell-api
func Load(action, level string) (reader io.Reader, err error) {
	action = SHCmd(action)
	x := APICmd()
	if _, ok := x[action]; !ok {
		return reader, fmt.Errorf("API command %v is invalid", level)
	}

	cmd := exec.Command("/bin/bash")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%v %v", x[action], level))

	stdout, err := cmd.Output()
	if err != nil {
		return reader, err
	}

	reader = bytes.NewBuffer(stdout)
	return reader, err
}

// SHCmd returns the appropriate command for non-tty or tty configure context
func SHCmd(a string) string {
	// execute := RunCMDline{}
	if !Insession() {
		switch a {
		case "listNodes":
			a = "listActiveNodes"
		case "returnValue":
			a = "returnActiveValue"
		case "returnValues":
			a = "returnActiveValues"
		case "exists":
			a = "existsActive"
		case "showConfig":
			a = "showCfg"
		}
	}
	return a
}

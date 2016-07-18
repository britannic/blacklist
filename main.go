package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	e "github.com/britannic/blacklist/internal/edgeos"
)

const (
	all   = "all"
	files = "file"
	pre   = "pre-configured"
	urls  = "url"
)

var (
	// Version vars updated by go build -ldflags
	build   = "UNKNOWN"
	githash = "UNKNOWN"
	version = "UNKNOWN"

	exitCmd    = os.Exit
	logFatalln = log.Fatalln
	logPrintf  = log.Printf
	objex      = []e.IFace{
		e.ExRtObj,
		e.ExDmObj,
		e.ExHtObj,
		e.PreDObj,
		e.PreHObj,
		e.FileObj,
		e.URLdObj,
		e.URLhObj,
	}
)

func main() {
	c := setUpEnv()
	logFatalln(removeStaleFiles(c))
	logFatalln(processObjects(c, objex))

	reloadDNS(c)
}

// basename removes directory components and file extensions.
func basename(s string) string {
	// Discard last '/' and everything before.
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break
		}
	}

	// Preserve everything before last '.'.
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}
	return s
}

func (o *opts) initEdgeOS() *e.Config {
	return e.NewConfig(
		e.API("/bin/cli-shell-api"),
		e.Arch(runtime.GOARCH),
		e.Bash("/bin/bash"),
		e.Cores(2),
		e.Debug(*o.Debug),
		e.Dir(o.setDir(*o.ARCH)),
		e.DNSsvc("service dnsmasq restart"),
		e.Ext("blacklist.conf"),
		e.File(*o.File),
		e.FileNameFmt("%v/%v.%v.%v"),
		e.InCLI("inSession"),
		e.Level("service dns forwarding"),
		e.Method("GET"),
		e.Nodes([]string{"domains", "hosts"}),
		e.Poll(*o.Poll),
		e.Prefix("address="),
		e.LTypes([]string{files, e.PreDomns, e.PreHosts, urls}),
		e.Timeout(30*time.Second),
		e.WCard(e.Wildcard{Node: "*s", Name: "*"}),
	)
}

func processObjects(c *e.Config, objects []e.IFace) error {
	for _, o := range objects {
		ct, err := c.NewContent(o)
		if err != nil {
			return err
		}

		err = c.ProcessContent(ct)
		if err != nil {
			return err
		}
	}
	return nil
}

func reloadDNS(c *e.Config) {
	b, err := c.ReloadDNS()
	if err != nil {
		logPrintf("ReloadDNS(): %v\n error: %v\n", string(b), err)
		exitCmd(1)
	}
	logPrintf("ReloadDNS(): %v\n", string(b))
	return
}

func removeStaleFiles(c *e.Config) error {
	fmt.Println(c.Wildcard)
	if err := c.GetAll().Files().Remove(); err != nil {
		return fmt.Errorf("c.GetAll().Files().Remove() error: %v\n", err)
	}
	return nil
}

func setUpEnv() *e.Config {
	o := getOpts()
	o.Init("blacklist", flag.ExitOnError)
	o.setArgs()

	c := o.initEdgeOS()
	c.ReadCfg(o.getCFG(c))

	return c
}

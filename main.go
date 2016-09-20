package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	e "github.com/britannic/blacklist/internal/edgeos"
	logging "github.com/op/go-logging"
)

const (
	all   = "all"
	files = "file"
	pre   = "pre-configured"
	urls  = "url"
)

var (
	// updated by go build -ldflags
	build   = "UNKNOWN"
	githash = "UNKNOWN"
	version = "UNKNOWN"
	// ---

	exitCmd  = os.Exit
	log, err = newLog()
	logError = func(args ...interface{}) {
		log.Error(args)
	}

	logErrorf = func(s string, args ...interface{}) {
		log.Errorf(s, args)
	}

	logCrit    = log.Critical
	logFatalln = func(args ...interface{}) {
		logCrit(args)
		exitCmd(1)
	}

	logFile    = "blacklist.log"
	logInfo    = log.Info
	logInfof   = log.Infof
	logPrintf  = logInfof
	logPrintln = logInfo

	objex = []e.IFace{
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

func newLog() (*logging.Logger, error) {
	fdFmt := logging.MustStringFormatter(
		`%{level:.4s}[%{id:03x}]%{time:2006-01-02 15:04:05.000} ▶ %{message}`,
	)

	scrFmt := logging.MustStringFormatter(
		`%{color:bold}%{level:.4s}%{color:reset}[%{id:03x}]%{time:15:04:05.000} ▶ %{message}`,
	)

	fd, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	fdlog := logging.NewLogBackend(fd, "", 0)
	fdFmttr := logging.NewBackendFormatter(fdlog, fdFmt)

	scr := logging.NewLogBackend(os.Stderr, "", 0)
	scrFmttr := logging.NewBackendFormatter(scr, scrFmt)

	logging.SetBackend(fdFmttr, scrFmttr)

	return logging.MustGetLogger(basename(os.Args[0])), err
}

func main() {

	c := setUpEnv()
	logInfo("Starting up...")
	if err := removeStaleFiles(c); err != nil {
		logFatalln(err)
	}

	// if err := processObjects(c, objex); err != nil {
	// 	logFatalln(err)
	// }

	logInfo("Shutting down...")
	// reloadDNS(c)
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

	// Preserve everything before last '.'
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
		e.Dbug(*o.Dbug),
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
		e.Logger(log),
		e.LTypes([]string{files, e.PreDomns, e.PreHosts, urls}),
		e.Timeout(30*time.Second),
		e.Verb(*o.Verb),
		e.WCard(e.Wildcard{Node: "*s", Name: "*"}),
		e.Writer(ioutil.Discard),
	)
}

func processObjects(c *e.Config, objects []e.IFace) error {
	for _, o := range objects {
		ct, err := c.NewContent(o)
		if err != nil {
			return err
		}

		if err = c.ProcessContent(ct); err != nil {
			return err
		}
	}
	return nil
}

func reloadDNS(c *e.Config) {
	b, err := c.ReloadDNS()
	if err != nil {
		logErrorf("ReloadDNS(): \n error: %v\n", string(b), err)
		exitCmd(1)
	}
	logPrintf("ReloadDNS(): %v\n", string(b))
}

func removeStaleFiles(c *e.Config) error {
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

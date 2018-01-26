package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	e "github.com/britannic/blacklist/internal/edgeos"
	"github.com/britannic/mflag"
	logging "github.com/op/go-logging"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	files = "file"
	urls  = "url"
)

var (
	// updated by go build -ldflags
	build   = "UNKNOWN"
	githash = "UNKNOWN"
	version = "UNKNOWN"
	// ---

	progname = basename(os.Args[0])
	exitCmd  = os.Exit
	fdFmttr  logging.Backend
	log      = newLog()

	logErrorf = func(s string, args ...interface{}) {
		log.Errorf(s, args)
	}

	logCrit = log.Critical

	logFatalln = func(args ...interface{}) {
		err := fmt.Sprintf("%v", args)
		logCrit(err)
		exitCmd(1)
	}

	logFile   = fmt.Sprintf("/var/log/%s.log", progname)
	logInfo   = log.Info
	logInfof  = log.Infof
	logNotice = log.Notice
	logPrintf = logInfof
	// logPrintln = logInfo

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

func main() {
	var (
		env *e.Config
		err error
	)

	if env, err = initEnv(); err != nil {
		logErrorf("%s shutting down.", err)
		exitCmd(0)
	}

	logNotice(fmt.Sprintf("Starting blacklist update..."))

	logInfo("Removing stale blacklists...")
	if err = removeStaleFiles(env); err != nil {
		logFatalln(err)
	}

	// _, _ = context.WithTimeout(context.Background(), c.Timeout)

	if !env.Disabled {
		if err := processObjects(env, objex); err != nil {
			logFatalln(err)
		}
	}

	reloadDNS(env)

	logNotice(fmt.Sprintf("Blacklist update completed......"))
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
		e.Disabled(false),
		e.Dbug(*o.Dbug),
		e.Dir(o.setDir(*o.ARCH)),
		e.DNSsvc("/etc/init.d/dnsmasq restart"),
		e.Ext("blacklist.conf"),
		e.File(*o.File),
		e.FileNameFmt("%v/%v.%v.%v"),
		e.InCLI("inSession"),
		e.Level("service dns forwarding"),
		e.Method("GET"),
		e.Prefix("address="),
		e.Logger(log),
		e.LTypes([]string{files, e.PreDomns, e.PreHosts, urls}),
		e.Timeout(30*time.Second),
		e.Verb(*o.Verb),
		e.WCard(e.Wildcard{Node: "*s", Name: "*"}),
		e.Writer(ioutil.Discard),
	)
}

func initEnv() (env *e.Config, err error) {
	if env, err = setUpEnv(); err != nil {
		d := killFiles(env)

		logInfo(progname + ": commencing dnsmasq blacklist update...")
		logInfo("Removing stale blacklists...")

		if err = d.Remove(); err != nil {
			logFatalln(err)
		}
		exitCmd(0)
	}
	return env, err
}

// killFiles returns an empty *e.CFile string array
func killFiles(env *e.Config) *e.CFile {
	return &e.CFile{Names: []string{}, Parms: env.Parms}
}

// newLog returns a logging.Logger pointer
func newLog() *logging.Logger {
	if runtime.GOOS == "darwin" {
		logFile = fmt.Sprintf("/tmp/%s.log", progname)
	}
	fdFmt := logging.MustStringFormatter(
		`%{level:.4s}[%{id:03x}]%{time:2006-01-02 15:04:05.000}: %{message}`,
	)

	fd, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}

	fdlog := logging.NewLogBackend(fd, progname+": ", 0)
	fdFmttr = logging.NewBackendFormatter(fdlog, fdFmt)

	sysFmttr, err := logging.NewSyslogBackend(progname + ": ")
	if err != nil {
		fmt.Println(err)
	}

	logging.SetBackend(fdFmttr, sysFmttr)

	return logging.MustGetLogger(progname)
}

// processObjects processes local sources, downloads Internet sources and creates
// dnsmasq configuration files
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

// reloadDNS reloads the latest processed dnsmasq configuration files
func reloadDNS(c *e.Config) {
	b, err := c.ReloadDNS()
	if err != nil {
		logErrorf("ReloadDNS(): \n error: %v\n", string(b), err)
		exitCmd(1)
	}
	logPrintf("ReloadDNS(): %v\n", string(b))
}

// removeStaleFiles deletes redundant files
func removeStaleFiles(c *e.Config) error {
	if err := c.GetAll().Files().Remove(); err != nil {
		return fmt.Errorf("c.GetAll().Files().Remove() error: %v", err)
	}
	return nil
}

// screenLog adds stderr logging output to the screen
func screenLog() {
	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		if runtime.GOOS == "darwin" {
			logFile = fmt.Sprintf("/tmp/%s.log", progname)
		}
		scrFmt := logging.MustStringFormatter(
			`%{color:bold}%{level:.4s}%{color:reset}[%{id:03x}]%{time:15:04:05.000}: %{message}`,
		)

		scr := logging.NewLogBackend(os.Stderr, progname+": ", 0)
		colors := []string{
			logging.CRITICAL: logging.ColorSeqBold(logging.ColorMagenta),
			logging.DEBUG:    logging.ColorSeqBold(logging.ColorCyan),
			logging.ERROR:    logging.ColorSeqBold(logging.ColorRed),
			logging.INFO:     logging.ColorSeqBold(logging.ColorGreen),
			logging.NOTICE:   logging.ColorSeqBold(logging.ColorBlue),
			logging.WARNING:  logging.ColorSeqBold(logging.ColorYellow),
		}
		scr.Color = true
		scr.ColorConfig = colors

		scrFmttr := logging.NewBackendFormatter(scr, scrFmt)
		sysFmttr, err := logging.NewSyslogBackend(progname + ": ")
		if err != nil {
			fmt.Println(err)
		}
		if runtime.GOOS != "amd64" {
			logging.SetBackend(fdFmttr, scrFmttr, sysFmttr)
		}
	}
}

func setUpEnv() (*e.Config, error) {
	o := getOpts()
	o.Init("blacklist", mflag.ExitOnError)
	o.setArgs()
	c := o.initEdgeOS()
	err := c.ReadCfg(o.getCFG(c))

	return c, err
}

package main

import (
	"fmt"
	"os"
	"runtime"

	e "github.com/britannic/blacklist/internal/edgeos"
	logging "github.com/britannic/go-logging"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	// updated by go build -ldflags
	architecture = "UNKNOWN"
	build        = "UNKNOWN"
	githash      = "UNKNOWN"
	hostOS       = "UNKNOWN"
	version      = "UNKNOWN"
	// ----------------------------

	boldcolors = []string{
		logging.CRITICAL: logging.ColorSeqBold(logging.ColorMagenta),
		logging.ERROR:    logging.ColorSeqBold(logging.ColorRed),
		logging.INFO:     logging.ColorSeqBold(logging.ColorGreen),
		logging.WARNING:  logging.ColorSeqBold(logging.ColorYellow),
		logging.NOTICE:   logging.ColorSeqBold(logging.ColorCyan),
		logging.DEBUG:    logging.ColorSeqBold(logging.ColorBlue),
	}
	exitCmd      = os.Exit
	fdFmttr      logging.Backend
	haveTerm     = inTerminal
	initEnvirons = initEnv
	log          = newLog(prefix)
	logCritf     = log.Criticalf
	logErrorf    = func(f string, args ...interface{}) { log.Errorf(f, args...) }
	logFatalf    = func(f string, args ...interface{}) { logCritf(f, args...); exitCmd(1) }
	logFile      = fmt.Sprintf("/var/log/%s.log", progname)
	logInfo      = log.Info
	logInfof     = log.Infof
	logNoticef   = log.Noticef
	logPrintf    = logInfof
	progname     = basename(os.Args[0])
	prefix       = fmt.Sprintf("%s: ", progname)
	objex        = []e.IFace{
		e.PreRObj,
		e.PreDObj,
		e.PreHObj,
		e.ExRtObj,
		e.ExDmObj,
		e.ExHtObj,
		e.FileObj,
		e.URLdObj,
		e.URLhObj,
	}
)

func main() {
	c, err := initEnvirons()
	if err != nil {
		logErrorf("%s shutting down.", err.Error())
		exitCmd(0)
	}

	c.Debug(fmt.Sprintf("Dumping commandline args: %v", os.Args[1:]))
	c.Debug(fmt.Sprintf("Dumping env variables: %v", c))
	logNoticef("%v", "Starting blacklist update...")

	logInfo("Removing stale blacklists...")
	if err = removeStaleFiles(c); err != nil {
		logFatalf("%v", err.Error())
	}

	// _, _ = context.WithTimeout(context.Background(), c.Timeout)

	if !c.Disabled {
		if err := processObjects(c, objex); err != nil {
			logErrorf("%v", err.Error())
		}
	}

	c.GetTotalStats()
	reloadDNS(c)
	logNoticef("%v", "Blacklist update completed......")
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

// files returns an empty *e.CFile string array
func files(c *e.Config) *e.CFile {
	return &e.CFile{Names: []string{}, Env: c.Env}
}

func initEnv() (c *e.Config, err error) {
	if c, err = setUpEnv(); err != nil {
		fmt.Fprintf(os.Stderr, "Removing stale dnsmasq blaclist files, because %v\n", err.Error())
		if err = files(c).Remove(); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err.Error())
		}
		exitCmd(0)
	}
	return c, err
}

func inTerminal() bool {
	return terminal.IsTerminal(int(os.Stdin.Fd()))
}

// newLog returns a logging.Logger pointer
func newLog(prefix string) *logging.Logger {
	if runtime.GOOS == "darwin" {
		logFile = fmt.Sprintf("/tmp/%s.log", progname)
	}
	fdFmt := logging.MustStringFormatter(
		`%{level:.4s}[%{id:03x}]%{time:2006-01-02 15:04:05.000}: %{message}`,
	)

	fd, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	fdlog := logging.NewLogBackend(fd, "", 0)
	fdFmttr = logging.NewBackendFormatter(fdlog, fdFmt)

	sysFmttr, err := logging.NewSyslogBackend(progname + ": ")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
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
	if b, err := c.ReloadDNS(); err != nil {
		logErrorf("ReloadDNS(): %v\n error: %v\n", string(b), err.Error())
		exitCmd(1)
	}
	logPrintf("%s", "Successfully restarted dnsmasq")
}

// removeStaleFiles deletes redundant files
func removeStaleFiles(c *e.Config) error {
	if err := c.GetAll().Files().Remove(); err != nil {
		return fmt.Errorf("problem removing stale files: %v", err.Error())
	}
	return nil
}

// screenLog adds stderr logging output to the screen
func screenLog(prefix string) logging.LeveledBackend {
	if haveTerm() {
		var (
			err      error
			scrFmt   = `%{color:bold}%{level:.4s}%{color:reset}[%{id:03x}]%{time:15:04:05.000}: %{message}`
			sysFmttr *logging.SyslogBackend
		)

		if runtime.GOOS == "darwin" {
			logFile = fmt.Sprintf("/tmp/%s.log", progname)
		}

		if sysFmttr, err = logging.NewSyslogBackend(prefix); err != nil {
			fmt.Println(err.Error())
		}

		return logging.SetBackend(
			logging.NewBackendFormatter(
				newScreenLogBackend(boldcolors, prefix),
				logging.MustStringFormatter(scrFmt),
			),
			fdFmttr,
			sysFmttr,
		)
	}
	return nil
}

func newScreenLogBackend(colors []string, prefix string) *logging.LogBackend {
	scr := logging.NewLogBackend(os.Stderr, prefix, 0)
	if len(colors) > 0 {
		scr.ColorConfig = boldcolors
		scr.Color = true
	}
	return scr
}

func setUpEnv() (c *e.Config, err error) {
	o := getOpts()
	o.setArgs()
	c = o.initEdgeOS()
	err = c.ReadCfg(o.getCFG(c))
	return c, err
}

package main

import (
	"fmt"
	"os"
	"runtime"

	logging "github.com/britannic/go-logging"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	boldcolors = []string{
		logging.CRITICAL: logging.ColorSeqBold(logging.ColorMagenta),
		logging.ERROR:    logging.ColorSeqBold(logging.ColorRed),
		logging.INFO:     logging.ColorSeqBold(logging.ColorGreen),
		logging.WARNING:  logging.ColorSeqBold(logging.ColorYellow),
		logging.NOTICE:   logging.ColorSeqBold(logging.ColorCyan),
		logging.DEBUG:    logging.ColorSeqBold(logging.ColorBlue),
	}
	fdFmttr    logging.Backend
	haveTerm   = inTerminal
	log        = newLog(prefix)
	logCritf   = log.Criticalf
	logErrorf  = func(f string, args ...interface{}) { log.Errorf(f, args...) }
	logFatalf  = func(f string, args ...interface{}) { logCritf(f, args...); exitCmd(1) }
	logFile    = setLogFile(runtime.GOOS)
	logInfo    = log.Info
	logInfof   = log.Infof
	logNoticef = log.Noticef
	logPrintf  = logInfof
)

// inTerminal returns true if the current terminal is interactive
func inTerminal() bool {
	return terminal.IsTerminal(int(os.Stdin.Fd()))
}

// setLogFile returns a log directory and file name dependent on the current OS
func setLogFile(os string) string {
	if os == "darwin" {
		return fmt.Sprintf("/tmp/%s.log", prog)
	}
	return fmt.Sprintf("/var/log/%s.log", prog)
}

// newLog returns a logging.Logger pointer
func newLog(prefix string) *logging.Logger {
	fdFmt := logging.MustStringFormatter(
		`%{level:.4s}[%{id:03x}]%{time:2006-01-02 15:04:05.000}: %{message}`,
	)
	// nolint
	fd, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	fdlog := logging.NewLogBackend(fd, "", 0)
	fdFmttr = logging.NewBackendFormatter(fdlog, fdFmt)

	sysFmttr, err := logging.NewSyslogBackend(prog + ": ")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	logging.SetBackend(fdFmttr, sysFmttr)

	return logging.MustGetLogger(prog)
}

func newScreenLogBackend(colors []string, prefix string) *logging.LogBackend {
	scr := logging.NewLogBackend(os.Stderr, prefix, 0)
	if len(colors) > 0 {
		scr.ColorConfig = boldcolors
		scr.Color = true
	}
	return scr
}

// screenLog adds stderr logging output to the screen
func screenLog(prefix string) logging.LeveledBackend {
	if haveTerm() {
		var (
			err      error
			scrFmt   = `%{color:bold}%{level:.4s}%{color:reset}[%{id:03x}]%{time:15:04:05.000}: %{message}`
			sysFmttr *logging.SyslogBackend
		)

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

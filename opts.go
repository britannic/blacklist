package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	e "github.com/britannic/blacklist/internal/edgeos"
	"github.com/britannic/blacklist/internal/tdata"
	"github.com/britannic/mflag"
)

// opts struct for command line options and setting initial variables
type opts struct {
	*mflag.FlagSet
	ARCH    *string
	Dbug    *bool
	DNSdir  *string
	DNStmp  *string
	File    *string
	Help    *bool
	MIPSLE  *string
	MIPS64  *string
	OS      *string
	Test    *bool
	Verb    *bool
	Version *bool
}

// cleanArgs removes flags when code is being tested
func cleanArgs(args []string) (r []string) {
	for _, a := range args {
		switch {
		case strings.HasPrefix(a, "-test"), strings.HasPrefix(a, "-convey"):
			continue
		default:
			r = append(r, a)
		}
	}
	return r
}

// getCFG returns a e.ConfLoader
func (o *opts) getCFG(c *e.Config) e.ConfLoader {
	if *o.File != "" {
		var (
			err error
			f   []byte
			r   io.Reader
		)

		if r, err = e.GetFile(*o.File); err != nil {
			logFatalf("cannot open configuration file %s!", *o.File)
		}

		if f, err = ioutil.ReadAll(r); err != nil {
			logFatalf("cannot read configuration file %s!", *o.File)
		}
		return &e.CFGstatic{Config: c, Cfg: string(f)}
	}
	switch *o.ARCH {
	case *o.MIPSLE, *o.MIPS64:
		return &e.CFGcli{Config: c}
	}
	return &e.CFGstatic{Config: c, Cfg: tdata.Live}
	// return &e.CFGstatic{Config: c, Cfg: tdata.Get("none")}

}

// getOpts returns command line flags and values or displays help
func getOpts() *opts {
	var (
		flags mflag.FlagSet
		o     = &opts{
			FlagSet: &flags,
			ARCH:    flags.String("arch", runtime.GOARCH, "Set EdgeOS CPU architecture", false),
			DNSdir:  flags.String("dir", "/etc/dnsmasq.d", "Override dnsmasq directory", true),
			DNStmp:  flags.String("tmp", "/tmp", "Override dnsmasq temporary directory", false),
			Dbug:    flags.Bool("debug", false, "Enable Debug mode", false),
			File:    flags.String("f", "", "`<file>` # Load a config.boot file", true),
			Help:    flags.Bool("h", false, "Display help", true),
			MIPS64:  flags.String("mips64", "mips64", "Override target EdgeOS CPU architecture", false),
			MIPSLE:  flags.String("mipsle", "mipsle", "Override target EdgeOS CPU architecture", false),
			OS:      flags.String("os", runtime.GOOS, "Override native EdgeOS OS", false),
			Test:    flags.Bool("dryrun", false, "Run config and data validation tests", false),
			Verb:    flags.Bool("v", false, "Verbose display", true),
			Version: flags.Bool("version", false, "Show version", true),
		}
	)
	flags.Init(prog, mflag.ExitOnError)
	flags.Usage = o.PrintDefaults

	return o
}

func (o *opts) initEdgeOS() *e.Config {
	dnsmasq := "/etc/init.d/dnsmasq"
	if _, err := os.Stat(dnsmasq); os.IsNotExist(err) {
		dnsmasq = "/bin/systemctl restart dnsmasq"
	}
	return e.NewConfig(
		e.API("/bin/cli-shell-api"),
		e.Arch(runtime.GOARCH),
		e.Bash("/bin/bash"),
		e.Cores(2),
		e.Disabled(false),
		e.Dbug(*o.Dbug),
		e.Dir(o.setDir(*o.ARCH)),
		e.DNSsvc(dnsmasq),
		e.Ext("blacklist.conf"),
		e.File(*o.File),
		e.FileNameFmt("%v/%v.%v.%v"),
		e.InCLI("inSession"),
		e.Level("service dns forwarding"),
		e.Method("GET"),
		e.Prefix("address=", "server="),
		e.Logger(log),
		e.Timeout(30*time.Second),
		e.Verb(*o.Verb),
		e.WCard(e.Wildcard{Node: "*s", Name: "*"}),
		// e.Writer(ioutil.Discard),
	)
}

// setArgs retrieves arguments entered on the command line
func (o *opts) setArgs() {
	if o.Parse(cleanArgs((os.Args[1:]))) != nil {
		// o.Usage()
		exitCmd(0)
	}

	if *o.Dbug {
		screenLog("")
		e.Dbug(*o.Dbug)
	}

	if *o.Help {
		o.Usage()
		exitCmd(0)
	}

	if *o.Test {
		fmt.Println("Test activated!")
		exitCmd(0)
	}

	if *o.Verb {
		screenLog("")
	}

	if *o.Version {
		fmt.Printf(
			" Build Information:\n"+
				"   Version:\t\t\t%s\n"+
				"   Date:\t\t\t%s\n"+
				"   CPU:\t\t\t\t%v\n"+
				"   OS:\t\t\t\t%v\n"+
				"   Git hash:\t\t\t%v\n\n"+
				" This software comes with ABSOLUTELY NO WARRANTY.\n"+
				" %s is free software, and you are\n"+
				" welcome to redistribute it under the terms of\n"+
				" the Simplified BSD License.\n",
			version,
			build,
			architecture,
			hostOS,
			githash,
			prog,
		)
		exitCmd(0)
	}
}

// setDir sets the directory according to the host CPU arch
func (o *opts) setDir(arch string) string {
	switch arch {
	case *o.MIPSLE, *o.MIPS64:
		return *o.DNSdir
	}
	return *o.DNStmp
}

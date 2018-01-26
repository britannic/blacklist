package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/britannic/blacklist/internal/edgeos"
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

// setDir sets the directory according to the host CPU arch
func (o *opts) setDir(arch string) (dir string) {
	switch arch {
	case *o.MIPSLE, *o.MIPS64:
		dir = *o.DNSdir
	default:
		dir = *o.DNStmp
	}
	return dir
}

// getCFG returns a edgeos.ConfLoader
func (o *opts) getCFG(c *edgeos.Config) (r edgeos.ConfLoader) {
	if *o.File != "" {
		var (
			f      []byte
			err    error
			reader io.Reader
		)

		if reader, err = edgeos.GetFile(*o.File); err != nil {
			logFatalln(fmt.Sprintf("Cannot open configuration file %s!", *o.File))
		}

		f, _ = ioutil.ReadAll(reader)
		r = &edgeos.CFGstatic{Config: c, Cfg: string(f)}
		return r
	}
	switch *o.ARCH {
	case *o.MIPSLE, *o.MIPS64:
		r = &edgeos.CFGcli{Config: c}
	default:
		r = &edgeos.CFGstatic{Config: c, Cfg: tdata.Live}
	}
	return r
}

// getOpts returns command line flags and values or displays help
func getOpts() *opts {
	var (
		flags mflag.FlagSet
		o     = &opts{
			ARCH:    flags.String("arch", runtime.GOARCH, "Set EdgeOS CPU architecture", false),
			Dbug:    flags.Bool("debug", false, "Enable debug mode", false),
			DNSdir:  flags.String("dir", "/etc/dnsmasq.d", "Override dnsmasq directory", true),
			DNStmp:  flags.String("tmp", "/tmp", "Override dnsmasq temporary directory", false),
			Help:    flags.Bool("h", false, "Display help", true),
			File:    flags.String("f", "", "`<file>` # Load a configuration file", true),
			FlagSet: &flags,
			MIPSLE:  flags.String("mipsle", "mipsle", "Override target EdgeOS CPU architecture", false),
			MIPS64:  flags.String("mips64", "mips64", "Override target EdgeOS CPU architecture", false),
			OS:      flags.String("os", runtime.GOOS, "Override native EdgeOS OS", false),
			Test:    flags.Bool("t", false, "Run config and data validation tests", false),
			Verb:    flags.Bool("v", false, "Verbose display", true),
			Version: flags.Bool("version", false, "Show version", true),
		}
	)
	flags.Init("blacklist", mflag.ExitOnError)
	flags.Usage = o.PrintDefaults

	return o
}

// cleanArgs removes flags when code is being tested
func cleanArgs(args []string) (r []string) {
NEXT:
	for _, a := range args {
		switch {
		case strings.HasPrefix(a, "-test"):
			continue NEXT
		case strings.HasPrefix(a, "-convey"):
			continue NEXT
		default:
			r = append(r, a)
		}
	}
	return r
}

// setArgs retrieves arguments entered on the command line
func (o *opts) setArgs() {
	if err := o.Parse(cleanArgs((os.Args[1:]))); err != nil {
		o.Usage()
		exitCmd(0)
	}

	switch {
	case *o.Help:
		o.Usage()
		exitCmd(0)

	case *o.Test:
		fmt.Println("Test activated!")
		exitCmd(0)

	case *o.Verb:
		screenLog()

	case *o.Version:
		fmt.Printf(" Version:\t\t%s\n Build date:\t\t%s\n Git short hash:\t%v\n\n This software comes with ABSOLUTELY NO WARRANTY.\n %s is free software, and you are\n welcome to redistribute it under the terms of\n the Simplified BSD License.\n", version, build, githash, progname)
		exitCmd(0)
	}
}

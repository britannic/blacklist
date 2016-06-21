package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	e "github.com/britannic/blacklist/internal/edgeos"
	"github.com/britannic/blacklist/internal/tdata"
	"github.com/fatih/structs"
)

// Opts struct for command line options and setting initial variables
type Opts struct {
	*flag.FlagSet
	ARCH    *string
	Debug   *bool
	DNSdir  *string
	DNStmp  *string
	File    *string
	MIPS64  *string
	OS      *string
	Poll    *int
	Test    *bool
	Verb    *bool
	Version *bool
}

// SetDir sets the directory according to the host CPU arch
func (o *Opts) SetDir(arch string) (dir string) {
	switch arch {
	case *o.MIPS64:
		dir = *o.DNSdir
	default:
		dir = *o.DNStmp
	}
	return dir
}

// GetCFG returns a loaded *Config
func (o *Opts) getCFG(arch string) (c *e.Config, err error) {
	c = &e.Config{Parms: &e.Parms{}}
	switch arch {
	case *o.MIPS64:
		r := &e.CFGcli{}
		c, err = e.ReadCfg(r)
	default:
		r := &e.CFGstatic{Cfg: tdata.Cfg}
		c, err = e.ReadCfg(r)
	}
	return c, err
}

func (o *Opts) String() (result string) {
	for _, name := range structs.Names(&Opts{}) {
		result += name + "\n"
	}

	return result
}

// getOpts returns command line flags and values or displays help
func getOpts() *Opts {
	var flags flag.FlagSet
	flags.Init("blacklist", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options]\n\n", basename(os.Args[0]))
		flags.PrintDefaults()
	}

	return &Opts{
		ARCH:    flags.String("arch", runtime.GOARCH, "Set EdgeOS CPU architecture"),
		Debug:   flags.Bool("debug", false, "Enable debug mode"),
		DNSdir:  flags.String("dir", "/etc/dnsmasq.d", "Override dnsmasq directory"),
		DNStmp:  flags.String("tmp", "/tmp", "Override dnsmasq temporary directory"),
		File:    flags.String("f", "", "<file> # Load a configuration file"),
		FlagSet: &flags,
		MIPS64:  flags.String("mips64", "mips64", "Override target EdgeOS CPU architecture"),
		OS:      flags.String("os", runtime.GOOS, "Override native EdgeOS OS"),
		Poll:    flags.Int("i", 5, "Polling interval"),
		Test:    flags.Bool("test", false, "Run config and data validation tests"),
		Verb:    flags.Bool("v", false, "Verbose display"),
		Version: flags.Bool("version", false, "Show program version number"),
	}
}

func (o *Opts) setArgs(fn func(int)) {
	if os.Args[1:] != nil {
		if err := o.Parse(os.Args[1:]); err != nil {
			o.Usage()
		}
	}

	switch {
	case *o.Test:
		fmt.Println("Test activated!")
		fn(0)

	case *o.Version:
		fmt.Printf(" Version:\t\t%s\n Build date:\t\t%s\n Git short hash:\t%v\n", version, build, githash)
		fn(0)
	}
}

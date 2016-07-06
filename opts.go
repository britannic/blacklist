package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/britannic/blacklist/internal/edgeos"
	"github.com/britannic/blacklist/internal/tdata"
)

// opts struct for command line options and setting initial variables
type opts struct {
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

// setDir sets the directory according to the host CPU arch
func (o *opts) setDir(arch string) (dir string) {
	switch arch {
	case *o.MIPS64:
		dir = *o.DNSdir
	default:
		dir = *o.DNStmp
	}
	return dir
}

// getCFG returns a e.ConfLoader
func (o *opts) getCFG(c *edgeos.Config) (r edgeos.ConfLoader) {
	switch *o.ARCH {
	case *o.MIPS64:
		r = &edgeos.CFGcli{Config: c}
	default:
		r = &edgeos.CFGstatic{Config: c, Cfg: tdata.Cfg}
	}
	return r
}

// getOpts returns command line flags and values or displays help
func getOpts() *opts {
	var flags flag.FlagSet
	flags.Init("blacklist", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options]\n\n", basename(os.Args[0]))
		flags.PrintDefaults()
	}

	return &opts{
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

func (o *opts) setArgs(fn func(int)) {
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

func (o *opts) String() string {
	type pArray struct {
		n string
		i int
		v string
	}

	var fields []pArray

	maxLen := func(pA []pArray) int {
		smallest := len(pA[0].n)
		largest := len(pA[0].n)
		for i := range pA {
			if len(pA[i].n) > largest {
				largest = len(pA[i].n)
			} else if len(pA[i].n) < smallest {
				smallest = len(pA[i].n)
			}
		}
		return largest
	}

	visitor := func(a *flag.Flag) {
		field := pArray{n: fmt.Sprint(a.Name), v: fmt.Sprint(a.Value)}
		fields = append(fields, field)
	}

	o.VisitAll(visitor)

	max := maxLen(fields)
	pad := func(s string) string {
		i := len(s)
		repeat := max - i + 1
		return strings.Repeat(" ", repeat)
	}

	s := "FlagSet\n"
	for _, field := range fields {
		if field.v == "" {
			field.v = "**not initialized**"
		}
		s += fmt.Sprintf("%v:%v%q\n", strings.ToUpper(field.n), pad(field.n), field.v)
	}

	return s
}

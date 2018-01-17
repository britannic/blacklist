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

type omitFlags map[string]bool

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
			ARCH:    flags.String("arch", runtime.GOARCH, "Set EdgeOS CPU architecture"),
			Dbug:    flags.Bool("debug", false, "Enable debug mode"),
			DNSdir:  flags.String("dir", "/etc/dnsmasq.d", "Override dnsmasq directory"),
			DNStmp:  flags.String("tmp", "/tmp", "Override dnsmasq temporary directory"),
			Help:    flags.Bool("h", false, "Display help"),
			File:    flags.String("f", "", "`<file>` # Load a configuration file"),
			FlagSet: &flags,
			MIPSLE:  flags.String("mipsle", "mipsle", "Override target EdgeOS CPU architecture"),
			MIPS64:  flags.String("mips64", "mips64", "Override target EdgeOS CPU architecture"),
			OS:      flags.String("os", runtime.GOOS, "Override native EdgeOS OS"),
			Test:    flags.Bool("t", false, "Run config and data validation tests"),
			Verb:    flags.Bool("v", false, "Verbose display"),
			Version: flags.Bool("version", false, "Show version"),
		}
	)
	flags.Init("blacklist", mflag.ExitOnError)
	// flags.Usage = o.PrintDefaults
	flags.Usage = func() {
		o.printDefaults(
			omitFlags{
				"arch":   true,
				"debug":  true,
				"f":      true,
				"mipsle": true,
				"mips64": true,
				"os":     true,
				"t":      true,
				"tmp":    true,
			})
	}

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
		fmt.Printf(" Version:\t\t%s\n Build date:\t\t%s\n Git short hash:\t%v\n", version, build, githash)
		exitCmd(0)
	}
}

// printDefaults prints to standard error the default values of all
// defined command-line flags in the set. See the documentation for
// the global function PrintDefaults for more information.
func (o *opts) printDefaults(omit omitFlags) {
	o.VisitAll(func(f *mflag.Flag) {
		if !omit[f.Name] {
			s := fmt.Sprintf("  -%s", f.Name) // Two spaces before -; see next two comments.

			name, usage := mflag.UnquoteUsage(f)
			if len(name) > 0 {
				s += " " + name
			}
			// Boolean flags of one ASCII letter are so common we
			// treat them specially, putting their usage on the same line.
			if len(s) <= 4 { // space, space, '-', 'x'.
				s += "\t"
			} else {
				// Four spaces before the tab triggers good alignment
				// for both 4- and 8-space tab stops.
				s += "\n    \t"
			}
			s += usage
			if !mflag.IsZeroValue(f, f.DefValue) {
				if _, ok := f.Value.(*mflag.StringValue); ok {
					// put quotes on the value
					s += fmt.Sprintf(" (default %q)", f.DefValue)
				} else {
					s += fmt.Sprintf(" (default %v)", f.DefValue)
				}
			}
			fmt.Fprint(o.FlagSet.Out(), s, "\n")
		}
	})
}

func (o *opts) String() string {
	var s string
	o.VisitAll(func(f *mflag.Flag) {
		s += fmt.Sprintf("  -%s", f.Name) // Two spaces before -; see next two comments.

		name, usage := mflag.UnquoteUsage(f)
		if len(name) > 0 {
			s += " " + name
		}
		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if len(s) <= 4 { // space, space, '-', 'x'.
			s += "\t"
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			s += "\n    \t"
		}
		s += usage
		if !mflag.IsZeroValue(f, f.DefValue) {
			if _, ok := f.Value.(*mflag.StringValue); ok {
				// put quotes on the value
				s += fmt.Sprintf(" (default %q)", f.DefValue)
			} else {
				s += fmt.Sprintf(" (default %v)", f.DefValue)
			}
		}
		s = fmt.Sprint(s, "\n")

	})

	return s
}

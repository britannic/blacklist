package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	e "github.com/britannic/blacklist/internal/edgeos"
	"github.com/britannic/blacklist/internal/tdata"
	"github.com/fatih/structs"
)

var (
	// Versioning vars updated by go build -ldflags
	build   = "UNKNOWN"
	githash = "UNKNOWN"
	version = "UNKNOWN"

	cores    = runtime.NumCPU()
	dnsDir   = "/etc/dnsmasq.d"
	dnsTmp   = "/tmp"
	mips64   = "mips64"
	whatOS   = runtime.GOOS
	whatArch = runtime.GOARCH
)

func main() {
	o := getOpts()
	o.Init("blacklist", flag.ExitOnError)

	if os.Args[1:] != nil {
		if err := o.Parse(os.Args[1:]); err != nil {
			o.Usage()
		}
	}

	switch {
	case *o.Test:
		fmt.Println("Test activated!")
		os.Exit(0)

	case *o.Version:
		fmt.Printf(" Version:\t\t%s\n Build date:\t\t%s\n Git short hash:\t%v\n", version, build, githash)
		os.Exit(0)
	}

	c, err := getCFG(whatArch)
	if err != nil {
		log.Fatalf("Couldn't load configuration: %v", err)
	}
	p := e.NewParms(c)
	_ = p.SetOpt(
		e.Cores(runtime.NumCPU()),
		e.Debug(*o.Debug),
		e.Dir(setDir(whatArch)),
		e.Ext("blacklist.conf"),
		e.File(*o.File),
		e.Method("GET"),
		e.Poll(*o.Poll),
		e.STypes([]string{"files", "pre-configured", "urls"}),
	)
	fmt.Println(p)
}

// Opts struct for command line options
type Opts struct {
	*flag.FlagSet
	Debug   *bool
	File    *string
	Poll    *int
	Test    *bool
	Verb    *bool
	Version *bool
}

func (o *Opts) String() (result string) {
	for _, name := range structs.Names(&Opts{}) {
		result += name + "\n"
	}

	return result
}

// getOpts returns command line flags and values or displays help
func getOpts() Opts {
	flags := flag.NewFlagSet("blacklist", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options]\n\n", basename(os.Args[0]))
		flags.PrintDefaults()
	}

	return Opts{
		FlagSet: flags,
		File:    flags.String("f", "", "<file> # Load a configuration file"),
		Debug:   flags.Bool("debug", false, "Enable debug mode"),
		Poll:    flags.Int("i", 5, "Polling interval"),
		Test:    flags.Bool("test", false, "Run config and data validation tests"),
		Verb:    flags.Bool("v", false, "Verbose display"),
		Version: flags.Bool("version", false, "# show program version number"),
	}
}

func setDir(arch string) (dir string) {
	switch arch {
	case mips64:
		dir = dnsDir
	default:
		dir = dnsTmp
	}
	return dir
}

func getCFG(arch string) (c *e.Config, err error) {
	var cfg string
	c = &e.Config{}
	switch arch {
	case mips64:
		if cfg, err = e.LoadCfg(); err != nil {
			return c, err
		}
		c, err = e.ReadCfg(bytes.NewBufferString(cfg))
	default:
		c, err = e.ReadCfg(bytes.NewBufferString(tdata.Cfg))
	}
	return c, err
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

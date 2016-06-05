package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	e "github.com/britannic/blacklist/edgeos"
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
	flag.CommandLine.Parse(os.Args[1:])

	switch {
	case *o.Test:
		os.Exit(0)

	case *o.Version:
		fmt.Printf(" Version:\t\t%s\n Build date:\t\t%s\n Git short hash:\t%v\n", version, build, githash)
		os.Exit(0)
	}

	c := e.Config{}
	p := e.NewParms(&c)
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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] param>\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	return Opts{
		File:    flag.String("f", "", "<file> # Load a configuration file"),
		Debug:   flag.Bool("debug", false, "Enable debug mode"),
		Poll:    flag.Int("i", 5, "Polling interval"),
		Test:    flag.Bool("test", false, "Run config and data validation tests"),
		Verb:    flag.Bool("v", false, "Verbose display"),
		Version: flag.Bool("version", false, "# show program version number"),
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

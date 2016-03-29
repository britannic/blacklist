package main

import (
	"flag"

	"github.com/fatih/structs"
)

// Opts struct for command line options
type Opts struct {
	Debug   *bool
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

// getopts returns legal command lines flags and values or displays help
func getopts() (options Opts) {

	// options.file = flag.String("f", "/config/config.boot", "<file> # Load a configuration file")
	options.Poll = flag.Int("i", 5, "Polling interval")
	options.Debug = flag.Bool("debug", false, "Enable debug mode")
	options.Test = flag.Bool("test", false, "Run config and data validation tests")
	options.Verb = flag.Bool("v", false, "Verbose display")
	options.Version = flag.Bool("version", false, "# show program version number")

	// flag.Parse()

	return options
}

package main

import (
	"flag"

	"github.com/fatih/structs"
)

// parms struct for command line options
type parms struct {
	Debug   *bool
	File    *string
	Poll    *int
	Test    *bool
	Verb    *bool
	Version *bool
}

func (o *parms) String() (opts string) {
	for _, name := range structs.Names(&parms{}) {
		opts += name + "\n"
	}

	return opts
}

// getparms returns legal command lines flags and values or displays help
func getOpts() parms {
	return parms{
		Debug:   flag.Bool("debug", false, "Enable debug mode"),
		File:    flag.String("f", "", "<file> # Load a configuration file"),
		Poll:    flag.Int("i", 5, "Polling interval"),
		Test:    flag.Bool("test", false, "Run config and data validation tests"),
		Verb:    flag.Bool("v", false, "Verbose display"),
		Version: flag.Bool("version", false, "# show program version number"),
	}
}

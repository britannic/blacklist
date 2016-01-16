package main

import "flag"
import "fmt"
import "os"

var Build string = "UNKNOWN"
var Version string = "UNKNOWN"
var Githash string = "UNKNOWN,

type Opts struct {
	file    *string
	v       *bool
	version *bool
}

func getOpts() Opts {
	var options Opts
	options.file = flag.String("f", "/config/config.boot", "<file> # Load a configuration file")
	options.v = flag.Bool("v", false, "Verbose display")
	options.version = flag.Bool("version", false, "# show program version number")
	flag.Parse()

	return options
}

func main() {
	var options Opts = getOpts()
	fmt.Println(*options.v)
	fmt.Println(*options.file)

	if *options.version {
		fmt.Printf("%s version: %s, build date: %s\n", os.Args[0], Version, Build)
	}
}

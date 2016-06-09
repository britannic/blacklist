package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	e "github.com/britannic/blacklist/internal/edgeos"
)

var (
	// Version vars updated by go build -ldflags
	build   = "UNKNOWN"
	githash = "UNKNOWN"
	version = "UNKNOWN"
)

func main() {
	// m := &e.Mvars{
	// 	DNSdir:   "/etc/dnsmasq.d",
	// 	DNStmp:   "/tmp",
	// 	MIPS64:   "mips64",
	// 	WhatOS:   runtime.GOOS,
	// 	WhatArch: runtime.GOARCH,
	// }
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

	c, err := o.GetCFG(*o.ARCH)
	if err != nil {
		log.Fatalf("Couldn't load configuration: %v", err)
	}
	p := e.NewParms(c)
	_ = p.SetOpt(
		e.Cores(runtime.NumCPU()),
		e.Debug(*o.Debug),
		e.Dir(o.SetDir(*o.ARCH)),
		e.Ext(".blacklist.conf"),
		e.File(*o.File),
		e.Method("GET"),
		e.Poll(*o.Poll),
		e.STypes([]string{"files", "pre-configured", "urls"}),
	)
	fmt.Println(p)
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

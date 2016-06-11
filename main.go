package main

import (
	"flag"
	"log"
	"runtime"
	"syscall"

	e "github.com/britannic/blacklist/internal/edgeos"
)

var (
	// Version vars updated by go build -ldflags
	build   = "UNKNOWN"
	githash = "UNKNOWN"
	version = "UNKNOWN"
)

func main() {
	var (
		c   *e.Config
		err error
		o   = getOpts()
	)

	o.Init("blacklist", flag.ExitOnError)
	o.setArgs(func(code int) {
		syscall.Exit(code)
	})

	c, err = o.getCFG(*o.ARCH)
	if err != nil {
		log.Fatalf("Couldn't load configuration: %v", err)
	}

	p := e.NewParms(c)
	_ = p.SetOpt(
		e.Cores(runtime.NumCPU()),
		e.Debug(*o.Debug),
		e.Dir(o.SetDir(*o.ARCH)),
		e.Ext(".blacklist.conf"),
		e.Excludes(c.Get("all").Excludes()),
		// e.FileNameFmt(src.Parms.dir + "/%v.%v." + src.Parms.ext),
		e.File(*o.File),
		e.Method("GET"),
		e.Poll(*o.Poll),
		e.STypes([]string{"files", "pre-configured", "urls"}),
	)
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

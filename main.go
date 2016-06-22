package main

import (
	"flag"
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
	const (
		all = "all"
		pre = "pre-configured"
	)
	var o = getOpts()

	o.Init("blacklist", flag.ExitOnError)
	o.setArgs(func(code int) {
		syscall.Exit(code)
	})

	c := e.NewConfig(
		e.API("/bin/cli-shell-api"),
		e.Cores(runtime.NumCPU()),
		e.Debug(*o.Debug),
		e.Dir(o.SetDir(*o.ARCH)),
		e.Ext("blacklist.conf"),
		e.File(*o.File),
		e.FileNameFmt("%v/%v.%v.%v"),
		e.Level("service dns forwarding"),
		e.Method("GET"),
		e.Nodes([]string{"domains", "hosts"}),
		e.Poll(*o.Poll),
		e.Prefix("address="),
		e.STypes([]string{"file", pre, "url"}),
	)

	c.ReadCfg(o.getCFG())

	c.SetOpt(
		e.Excludes(c.Excludes("all")),
	)

	c.GetAll().Files().Remove()
	c.GetAll(pre).GetContent().ProcessContent()
	c.GetAll("url").GetContent().ProcessContent()
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

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"

	e "github.com/britannic/blacklist/edgeos"
	"github.com/britannic/blacklist/tdata"
)

var (
	// Versioning vars updated by go build -ldflags
	build   = "UNKNOWN"
	githash = "UNKNOWN"
	version = "UNKNOWN"

	cores    = runtime.NumCPU()
	dir      = setDir()
	dnsDir   = "/etc/dnsmasq.d"
	dnsTmp   = "/tmp"
	nodes    = []string{"domains", "hosts"}
	whatOS   = runtime.GOOS
	whatArch = runtime.GOARCH
)

func main() {
	var (
		reader io.Reader
		err    error
	)

	runtime.GOMAXPROCS(cores)

	o := getOpts()
	p := &e.Parms{}
	flag.CommandLine.Parse(os.Args[1:])

	switch {
	case *o.Debug:
		p.SetOpt(e.Debug(*o.Debug))

	case *o.File != "":
		reader, err = os.Open(*o.File)
		if err != nil {
			log.Fatalf("%v", err)
		}

	case *o.File == "" && whatArch != "mips64":
		reader = bytes.NewBufferString(tdata.Cfg)

	case *o.File == "":
		reader, err = e.Load("showCfg", "service dns forwarding")

	case *o.Poll != 5:
		p.SetOpt(e.Poll(time.Duration(*o.Poll) * time.Second))

	case *o.Test:
		code := 0
		os.Exit(code)

	case *o.Version:
		fmt.Printf(" Version:\t\t%s\n Build date:\t\t%s\n Git short hash:\t%v\n", version, build, githash)
		os.Exit(0)

		// case *o.Verb:
		// 	g.LogOutput = "screen"
		// 	logger := &g.Set{
		// 		Level:  logrus.DebugLevel,
		// 		Output: g.LogOutput,
		// 	}
		//
		// 	log = g.LogInit(logger)
	}

	b, err := e.ReadCfg(reader)
	if err != nil {
		log.Fatalf("%v", err)
	}

	c := b.NewConfig()
	removeFiles(c)
	allNodes := []string{"blacklist"}
	allNodes = append(allNodes, nodes...)

	c.GetExcludes(allNodes)

	for _, node := range nodes {
		for _, src := range c.Source(node, e.PreConf) {
			data := c.Process(&src, c.GetIncludes(node))
			fname := fmt.Sprintf("%v/%v.%v", dir, node, e.Fext)
			e.WriteFile(fname, data)
		}
	}
}

func removeFiles(c *e.Config) {
	files, err := e.ListFiles(dir)
	if err != nil {
		log.Printf("%v", err)
	}

	purgefiles := e.DiffArray(c.Files(dir, nodes), files)
	if err = e.PurgeFiles(purgefiles); err != nil {
		log.Printf("%v", err)
	}
}

func setDir() string {
	switch whatArch {
	case "mips64":
		return dnsDir
	default:
		return dnsTmp
	}
}

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"

	e "github.com/britannic/blacklist/edgeos"
)

var (
	// Versioning vars updated by go build -ldflags
	build   = "UNKNOWN"
	githash = "UNKNOWN"
	version = "UNKNOWN"

	cores = runtime.NumCPU()
	dbg   = false
	dir   = "/tmp"
	nodes = []string{"blacklist", "domains", "hosts"}
	poll  = time.Second * 2 // poll
)

func main() {
	var (
		reader io.Reader
		err    error
	)

	runtime.GOMAXPROCS(cores)

	o := getOpts()
	a := os.Args[1:]
	flag.CommandLine.Parse(a)

	switch {
	case *o.Debug:
		dbg = true

	case *o.File != "":
		reader, err = os.Open(*o.File)
		if err != nil {
			log.Fatalf("%v", err)
		}

	case *o.Poll != 5:
		poll = time.Duration(*o.Poll) * time.Second
		log.Printf("Poll duration %v", poll)

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

	if *o.File == "" {
		reader, err = e.Load("showCfg", "service dns forwarding")
	}

	b, err := e.ReadCfg(reader)
	if err != nil {
		log.Fatalf("%v", err)
	}

	c := b.NewConfig()

	files, err := e.ListFiles(dir)
	if err != nil {
		log.Printf("%v", err)
	}

	purgefiles := e.DiffArray(c.Files(dir, nodes), files)
	if err = e.PurgeFiles(purgefiles); err != nil {
		log.Printf("%v", err)
	}

	dex, ex := c.WriteIncludes(dir, nodes)
	dex, ex = c.GetExcludes(dex, ex, nodes)

	fmt.Println(dex)
	fmt.Println(ex)
}

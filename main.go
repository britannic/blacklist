package main

import (
	"fmt"
	"os"

	e "github.com/britannic/blacklist/internal/edgeos"
)

var (
	// updated by go build -ldflags
	architecture = "UNKNOWN"
	build        = "UNKNOWN"
	githash      = "UNKNOWN"
	hostOS       = "UNKNOWN"
	version      = "UNKNOWN"
	// ----------------------------

	exitCmd      = os.Exit
	initEnvirons = initEnv
	prog         = basename(os.Args[0])
	prefix       = fmt.Sprintf("%s: ", prog)
	defCfgFile   = "/config/user-data/blacklist.cfg"
)

func main() {
	objex := []e.IFace{
		e.PreRObj,
		e.PreDObj,
		e.PreHObj,
		e.ExRtObj,
		e.ExDmObj,
		e.ExHtObj,
		e.FileObj,
		e.URLdObj,
		e.URLhObj,
	}

	c, err := initEnvirons()
	if err != nil {
		logErrorf("%s shutting down.", err.Error())
		exitCmd(0)
	}

	c.Debug(fmt.Sprintf("Dumping commandline args: %v", os.Args[1:]))
	c.Debug(fmt.Sprintf("Dumping env variables: %v", c))
	logNoticef("%v", "Starting blacklist update...")

	if !e.ChkWeb("www.google.com", 443) {
		logFatalf("%s", "No internet access, aborting blacklist update!")
	}

	logInfo("Checking for stale blacklists...")
	if err = removeStaleFiles(c); err != nil {
		logFatalf("%v", err.Error())
	}

	// _, _ = context.WithTimeout(context.Background(), c.Timeout)

	if !c.Disabled {
		if err := processObjects(c, objex); err != nil {
			logErrorf("%v", err.Error())
		}
	}

	c.GetTotalStats()
	reloadDNS(c)
	logNoticef("%v", "Blacklist update completed......")
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

	// Preserve everything before last '.'
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}
	return s
}

// files returns an empty *e.CFile string array
func files(c *e.Config) *e.CFile {
	return &e.CFile{Names: []string{}, Env: c.Env}
}

func initEnv() (c *e.Config, err error) {
	o := getOpts()
	o.setArgs()
	c = o.initEdgeOS()
	if err = c.Blacklist(o.getCFG(c)); err != nil {
		fmt.Fprintf(os.Stderr, "Removing stale dnsmasq blacklist files, because %v\n", err.Error())
		if err = files(c).Remove(); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err.Error())
		}
		reloadDNS(c)
		exitCmd(0)
	}
	return c, err
}

// processObjects processes local sources, downloads Internet sources and creates
// dnsmasq configuration files
func processObjects(c *e.Config, objects []e.IFace) error {
	for _, o := range objects {
		ct, err := c.NewContent(o)
		if err != nil {
			return err
		}
		if err = c.ProcessContent(ct); err != nil {
			return err
		}
	}
	return nil
}

// reloadDNS reloads the latest processed dnsmasq configuration files
func reloadDNS(c *e.Config) {
	if b, err := c.ReloadDNS(); err != nil {
		logErrorf("ReloadDNS(): %v\n error: %v\n", string(b), err.Error())
		exitCmd(1)
	}
	logPrintf("%s", "Successfully restarted dnsmasq")
}

// removeStaleFiles deletes redundant files
func removeStaleFiles(c *e.Config) error {
	if err := c.GetAll().Files().Remove(); err != nil {
		return fmt.Errorf("problem removing stale files: %v", err.Error())
	}
	return nil
}

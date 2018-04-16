package edgeos

import (
	"encoding/json"
	"io"
	"runtime"
	"sync"
	"time"

	logging "github.com/britannic/go-logging"
)

// Env is struct of parameters
type Env struct {
	ctr
	ioWriter io.Writer
	Log      *logging.Logger
	API      string        `json:"API,omitempty"`
	Arch     string        `json:"Arch,omitempty"`
	Bash     string        `json:"Bash,omitempty"`
	Cores    int           `json:"Cores,omitempty"`
	Disabled bool          `json:"Disabled"`
	Dbug     bool          `json:"Dbug,omitempty"`
	Dex      list          `json:"Dex,omitempty"`
	Dir      string        `json:"Dir,omitempty"`
	DNSsvc   string        `json:"dnsmasq service,omitempty"`
	Exc      list          `json:"Exc,omitempty"`
	Ext      string        `json:"dnsmasq fileExt.,omitempty"`
	File     string        `json:"File,omitempty"`
	FnFmt    string        `json:"File name fmt,omitempty"`
	InCLI    string        `json:"-"`
	Level    string        `json:"CLI Path,omitempty"`
	Method   string        `json:"HTTP method,omitempty"`
	Pfx      dnsPfx        `json:"Prefix,omitempty"`
	Test     bool          `json:"Test,omitempty"`
	Timeout  time.Duration `json:"Timeout,omitempty"`
	Verb     bool          `json:"Verbosity,omitempty"`
	Wildcard/*..........*/ `json:"Wildcard,omitempty"`
}

// dnsPfx defines the prefix entries in the dnsmasq configuration file
type dnsPfx struct {
	domain string
	host   string
}

// Wildcard struct sets globbing wildcards for filename searches
type Wildcard struct {
	Node string `json:"Node,omitempty"`
	Name string `json:"Name,omitempty"`
}

// Option is a recursive function
type Option func(c *Config) Option

// Debug logs debug messages when the Dbug flag is true
func (p *Env) Debug(s ...interface{}) {
	if p.Dbug {
		p.Log.Debug(s...)
	}
}

// SetOpt sets the specified options passed as Env and returns an option to restore the last set of arg's previous values
func (c *Config) SetOpt(opts ...Option) Option {
	// apply all the options, and replace each with its inverse
	for i, opt := range opts {
		opts[i] = opt(c)
	}

	for i, j := 0, len(opts)-1; i <= j; i, j = i+1, j-1 {
		opts[i], opts[j] = opts[j], opts[i]
	}

	return func(c *Config) Option {
		return c.SetOpt(opts...)
	}
}

// Arch sets target CPU architecture
func Arch(arch string) Option {
	return func(c *Config) Option {
		previous := c.Arch
		c.Arch = arch
		return Arch(previous)
	}
}

// API sets the EdgeOS CLI API command
func API(s string) Option {
	return func(c *Config) Option {
		previous := c.API
		c.API = s
		return API(previous)
	}
}

// Bash sets the shell processor
func Bash(cmd string) Option {
	return func(c *Config) Option {
		previous := c.Bash
		c.Bash = cmd
		return Bash(previous)
	}
}

// Cores sets max CPU cores
func Cores(i int) Option {
	return func(c *Config) Option {
		previous := c.Cores
		runtime.GOMAXPROCS(i)
		c.Cores = i
		return Cores(previous)
	}
}

// Disabled toggles Disabled
func Disabled(b bool) Option {
	return func(c *Config) Option {
		previous := c.Disabled
		c.Disabled = b
		return Disabled(previous)
	}
}

// Dbug toggles Debug level on or off
func Dbug(b bool) Option {
	return func(c *Config) Option {
		previous := c.Dbug
		c.Dbug = b
		return Dbug(previous)
	}
}

// Dir sets directory location
func Dir(d string) Option {
	return func(c *Config) Option {
		previous := c.Dir
		c.Dir = d
		return Dir(previous)
	}
}

// DNSsvc sets dnsmasq restart command
func DNSsvc(d string) Option {
	return func(c *Config) Option {
		previous := c.DNSsvc
		c.DNSsvc = d
		return DNSsvc(previous)
	}
}

// Ext sets the blacklist file n extension
func Ext(e string) Option {
	return func(c *Config) Option {
		previous := c.Ext
		c.Ext = e
		return Ext(previous)
	}
}

// File sets the EdgeOS configuration file
func File(f string) Option {
	return func(c *Config) Option {
		previous := c.File
		c.File = f
		return File(previous)
	}
}

// FileNameFmt sets the EdgeOS configuration file name format
func FileNameFmt(f string) Option {
	return func(c *Config) Option {
		previous := c.FnFmt
		c.FnFmt = f
		return FileNameFmt(previous)
	}
}

// InCLI sets the CLI inSession command
func InCLI(in string) Option {
	return func(c *Config) Option {
		previous := c.InCLI
		c.InCLI = in
		return InCLI(previous)
	}
}

// Level sets the EdgeOS API CLI level
func Level(s string) Option {
	return func(c *Config) Option {
		previous := c.Level
		c.Level = s
		return Level(previous)
	}
}

// Logger sets a pointer to the logger
func Logger(l *logging.Logger) Option {
	return func(c *Config) Option {
		previous := c.Log
		c.Log = l
		return Logger(previous)
	}
}

// Method sets the HTTP method
func Method(method string) Option {
	return func(c *Config) Option {
		previous := c.Method
		c.Method = method
		return Method(previous)
	}
}

// NewConfig returns a new *Config initialized with the parameter options passed to it
func NewConfig(opts ...Option) *Config {
	c := Config{
		tree: make(tree),
		Env: &Env{
			ctr: make(ctr),
			Dex: list{RWMutex: &sync.RWMutex{}, entry: make(entry)},
			Exc: list{RWMutex: &sync.RWMutex{}, entry: make(entry)},
		},
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

// Prefix sets the dnsmasq configuration address line prefix
func Prefix(d string, h string) Option {
	return func(c *Config) Option {
		pd := c.Pfx.domain
		ph := c.Pfx.host
		c.Pfx = dnsPfx{domain: d, host: h}
		return Prefix(pd, ph)
	}
}

// String method to implement fmt.Print interface
func (p *Env) String() string {
	out, _ := json.MarshalIndent(p, "", "\t")
	return string(out)
}

// Test toggles testing mode on or off
func Test(b bool) Option {
	return func(c *Config) Option {
		previous := c.Test
		c.Test = b
		return Test(previous)
	}
}

// Timeout sets how long before an unresponsive goroutine is aborted
func Timeout(t time.Duration) Option {
	return func(c *Config) Option {
		previous := c.Timeout
		c.Timeout = t
		return Timeout(previous)
	}
}

// Verb sets the verbosity level to v
func Verb(b bool) Option {
	return func(c *Config) Option {
		previous := c.Verb
		c.Verb = b
		return Verb(previous)
	}
}

// WCard sets file globbing wildcard values
func WCard(w Wildcard) Option {
	return func(c *Config) Option {
		previous := c.Wildcard
		c.Wildcard = w
		return WCard(previous)
	}
}

// Writer provides an address for anything that can use io.Writer
func Writer(w io.Writer) Option {
	return func(c *Config) Option {
		previous := c.ioWriter
		c.ioWriter = w
		return Writer(previous)
	}
}

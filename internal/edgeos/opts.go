package edgeos

import (
	"encoding/json"
	"io"
	"runtime"
	"sync"
	"time"

	logging "github.com/op/go-logging"
)

// Parms is struct of parameters
type Parms struct {
	ioWriter io.Writer
	*logging.Logger
	API     string        `json:"API, omitempty"`
	Arch    string        `json:"Arch, omitempty"`
	Bash    string        `json:"Bash, omitempty"`
	Cores   int           `json:"Cores, omitempty"`
	Debug   bool          `json:"Debug, omitempty"`
	Dex     list          `json:"Dex, omitempty"`
	Dir     string        `json:"Dir, omitempty"`
	DNSsvc  string        `json:"dnsmasq service, omitempty"`
	Exc     list          `json:"Exc, omitempty"`
	Ext     string        `json:"dnsmasq fileExt., omitempty"`
	File    string        `json:"File, omitempty"`
	FnFmt   string        `json:"File name fmt, omitempty"`
	InCLI   string        `json:"-"`
	Level   string        `json:"CLI Path, omitempty"`
	Ltypes  []string      `json:"Leaf nodes, omitempty"`
	Method  string        `json:"HTTP method, omitempty"`
	Nodes   []string      `json:"Nodes, omitempty"`
	Pfx     string        `json:"Prefix, omitempty"`
	Poll    int           `json:"Poll, omitempty"`
	Test    bool          `json:"Test, omitempty"`
	Timeout time.Duration `json:"Timeout, omitempty"`
	Verb    bool          `json:"Verbosity, omitempty"`
	Wildcard/*..........*/ `json:"Wildcard, omitempty"`
}

// Wildcard struct sets globbing wildcards for filename searches
type Wildcard struct {
	Node string `json:"omitempty"`
	Name string `json:"omitempty"`
}

// Option is a recursive function
type Option func(c *Config) Option

// NewConfig returns a new *Config initialized with the parameter options passed to it
func NewConfig(opts ...Option) *Config {
	c := Config{
		tree: make(tree),
		Parms: &Parms{
			Dex: list{RWMutex: &sync.RWMutex{}, entry: make(entry)},
			Exc: list{RWMutex: &sync.RWMutex{}, entry: make(entry)},
		},
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

// SetOpt sets the specified options passed as Parms and returns an option to restore the last set of arg's previous values
func (c *Config) SetOpt(opts ...Option) (previous Option) {
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

// Debug toggles debug level on or off
func Debug(b bool) Option {
	return func(c *Config) Option {
		previous := c.Debug
		c.Debug = b
		return Debug(previous)
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

// File sets the EdgeOS configuration file n
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

// Logger sets the EdgeOS API CLI level
func Logger(l *logging.Logger) Option {
	return func(c *Config) Option {
		previous := c.Logger
		c.Logger = l
		return Logger(previous)
	}
}

// LTypes sets an array of legal types used by Source
func LTypes(s []string) Option {
	return func(c *Config) Option {
		previous := c.Ltypes
		c.Ltypes = s
		return LTypes(previous)
	}
}

// Prefix sets the dnsmasq configuration address line prefix
func Prefix(l string) Option {
	return func(c *Config) Option {
		previous := c.Pfx
		c.Pfx = l
		return Prefix(previous)
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

// Nodes sets the node ns array
func Nodes(nodes []string) Option {
	return func(c *Config) Option {
		previous := c.Parms.Nodes
		c.Parms.Nodes = nodes
		return Nodes(previous)
	}
}

// Poll sets the polling interval in seconds
func Poll(t int) Option {
	return func(c *Config) Option {
		previous := c.Poll
		c.Poll = t
		return Poll(previous)
	}
}

// String method to implement fmt.Print interface
func (p *Parms) String() string {
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

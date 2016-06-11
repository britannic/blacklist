package edgeos

import (
	"fmt"
	"runtime"
	"strings"
)

// Parms is struct of parameters
type Parms struct {
	arch      string
	cores     int
	dir       string
	dex       List
	debug     bool
	exc       List
	ext       string
	fnfmt     string
	file      string
	method    string
	nodes     []string
	poll      int
	stypes    []string
	test      bool
	verbosity int
}

// Option sets is a recursive function
type Option func(p *Parms) Option

// SetOpt sets the specified options passed as Parms and returns an option to restore the last arg's previous value
func (p *Parms) SetOpt(opts ...Option) (previous Option) {
	// apply all the options, and replace each with its inverse
	for i, opt := range opts {
		opts[i] = opt(p)
	}

	// Reverse the list of inverses, since we want them to be applied in reverse order
	for i, j := 0, len(opts)-1; i <= j; i, j = i+1, j-1 {
		opts[i], opts[j] = opts[j], opts[i]
	}

	return func(p *Parms) Option {
		return p.SetOpt(opts...)
	}
}

// GetOpt retrieves the value of an Option
// func (p *Parms) GetOpt(opt string) string {
// 	return
// }

// Arch sets target CPU architecture
func Arch(arch string) Option {
	return func(p *Parms) Option {
		previous := p.arch
		p.arch = arch
		return Arch(previous)
	}
}

// Cores sets max CPU cores
func Cores(i int) Option {
	return func(p *Parms) Option {
		previous := p.cores
		runtime.GOMAXPROCS(i)
		p.cores = i
		return Cores(previous)
	}
}

// Debug toggles debug level on or off
func Debug(b bool) Option {
	return func(p *Parms) Option {
		previous := p.debug
		p.debug = b
		return Debug(previous)
	}
}

// Dir sets directory location
func Dir(d string) Option {
	return func(p *Parms) Option {
		previous := p.dir
		p.dir = d
		return Dir(previous)
	}
}

// Excludes sets nodes exclusions
func Excludes(l List) Option {
	return func(p *Parms) Option {
		previous := p.exc
		p.exc = l
		return Excludes(previous)
	}
}

// Ext sets the blacklist file n extension
func Ext(e string) Option {
	return func(p *Parms) Option {
		previous := p.ext
		p.ext = e
		return Ext(previous)
	}
}

// File sets the EdgeOS configuration file n
func File(f string) Option {
	return func(p *Parms) Option {
		previous := p.file
		p.file = f
		return File(previous)
	}
}

// FileNameFmt sets the EdgeOS configuration file n
func FileNameFmt(f string) Option {
	return func(p *Parms) Option {
		previous := p.fnfmt
		p.fnfmt = f
		return File(previous)
	}
}

// Method sets the HTTP method
func Method(method string) Option {
	return func(p *Parms) Option {
		previous := p.method
		p.method = method
		return Method(previous)
	}
}

// NewParms sets a new *Parms instance
func NewParms(c *Config) *Parms {
	c.Parms = &Parms{
		dex: make(List),
		exc: make(List),
	}
	return c.Parms
}

// Nodes sets the node ns array
func Nodes(nodes []string) Option {
	return func(p *Parms) Option {
		previous := p.nodes
		p.nodes = nodes
		return Nodes(previous)
	}
}

// Poll sets the polling interval in seconds
func Poll(t int) Option {
	return func(p *Parms) Option {
		previous := p.poll
		p.poll = t
		return Poll(previous)
	}
}

// String method to implement fmt.Print interface
func (p Parms) String() string {
	max := 9
	pad := func(i int) string {
		repeat := max - i + 1
		return strings.Repeat(" ", repeat)
	}

	getVal := func(v interface{}) string {
		return strings.Replace(fmt.Sprint(v), "\n", "", -1)
	}

	fields := []struct {
		n string
		i int
		v string
	}{
		{n: "arch", i: 4, v: getVal(p.arch)},
		{n: "cores", i: 5, v: getVal(p.cores)},
		{n: "dir", i: 3, v: getVal(p.dir)},
		{n: "debug", i: 5, v: getVal(p.debug)},
		{n: "exc", i: 3, v: getVal(p.exc)},
		{n: "ext", i: 3, v: getVal(p.ext)},
		{n: "file", i: 4, v: getVal(p.file)},
		{n: "method", i: 6, v: getVal(p.method)},
		{n: "nodes", i: 5, v: getVal(p.nodes)},
		{n: "poll", i: 4, v: getVal(p.poll)},
		{n: "stypes", i: 6, v: getVal(p.stypes)},
		{n: "test", i: 4, v: getVal(p.test)},
		{n: "verbosity", i: 9, v: getVal(p.verbosity)},
	}

	r := fmt.Sprintln("edgeos.Parms{")
	for _, field := range fields {
		r += fmt.Sprintf("%v:%v%v\n", field.n, pad(field.i), field.v)
	}

	r += fmt.Sprintln("}")

	return r
}

// STypes sets an array of legal types used by Source
func STypes(s []string) Option {
	return func(p *Parms) Option {
		previous := p.stypes
		p.stypes = s
		return STypes(previous)
	}
}

// Test toggles testing mode on or off
func Test(b bool) Option {
	return func(p *Parms) Option {
		previous := p.test
		p.test = b
		return Test(previous)
	}
}

// Verbosity sets the verbosity level to v
func Verbosity(i int) Option {
	return func(p *Parms) Option {
		previous := p.verbosity
		p.verbosity = i
		return Verbosity(previous)
	}
}

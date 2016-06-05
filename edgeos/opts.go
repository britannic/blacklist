package edgeos

import (
	"fmt"
	"runtime"
)

// parms is struct of parameters
type parms struct {
	cores     int
	dir       string
	debug     bool
	exc       List
	ext       string
	file      string
	method    string
	nodes     []string
	poll      int
	stypes    []string
	test      bool
	verbosity int
}

// Option sets is a recursive function
type Option func(p *parms) Option

// SetOpt sets the specified options passed as parms and returns an option to restore the last arg's previous value
func (p *parms) SetOpt(opts ...Option) (previous Option) {
	// apply all the options, and replace each with its inverse
	for i, opt := range opts {
		opts[i] = opt(p)
	}

	// reverse the list of inverses, since we want them to be applied in reverse order
	for i, j := 0, len(opts)-1; i <= j; i, j = i+1, j-1 {
		opts[i], opts[j] = opts[j], opts[i]
	}

	return func(p *parms) Option {
		return p.SetOpt(opts...)
	}
}

// Cores sets max CPU cores
func Cores(i int) Option {
	return func(p *parms) Option {
		previous := p.cores
		runtime.GOMAXPROCS(i)
		p.cores = i
		return Cores(previous)
	}
}

// Debug toggles debug level on or off
func Debug(b bool) Option {
	return func(p *parms) Option {
		previous := p.debug
		p.debug = b
		return Debug(previous)
	}
}

// Dir sets directory location
func Dir(d string) Option {
	return func(p *parms) Option {
		previous := p.dir
		p.dir = d
		return Dir(previous)
	}
}

// Excludes sets nodes exclusions
func Excludes(l List) Option {
	return func(p *parms) Option {
		previous := p.exc
		p.exc = l
		return Excludes(previous)
	}
}

// Ext sets the blacklist file name extension
func Ext(e string) Option {
	return func(p *parms) Option {
		previous := p.ext
		p.ext = e
		return Ext(previous)
	}
}

// File sets the EdgeOS configuration file name
func File(f string) Option {
	return func(p *parms) Option {
		previous := p.file
		p.file = f
		return File(previous)
	}
}

// Method sets the HTTP method
func Method(method string) Option {
	return func(p *parms) Option {
		previous := p.method
		p.method = method
		return Method(previous)
	}
}

// NewParms sets a new *parms instance
func NewParms(c *Config) *parms {
	c.parms = &parms{}
	return c.parms
}

// Nodes sets the node names array
func Nodes(nodes []string) Option {
	return func(p *parms) Option {
		previous := p.nodes
		p.nodes = nodes
		return Nodes(previous)
	}
}

// Poll sets the polling interval in seconds
func Poll(t int) Option {
	return func(p *parms) Option {
		previous := p.poll
		p.poll = t
		return Poll(previous)
	}
}

// String method to implement fmt.Print interface
func (p *parms) String() (s string) {
	s += fmt.Sprintln("edgeos.parms{")
	s += fmt.Sprintf("cores:\t\t%v\n", p.cores)
	s += fmt.Sprintf("dir:\t\t%q\n", p.dir)
	s += fmt.Sprintf("debug:\t\t%t\n", p.debug)

	for k := range p.exc {
		s += fmt.Sprintf("exc:\t\t%q: %v\n", k, p.exc[k])
	}

	s += fmt.Sprintf("ext:\t\t%q\n", p.ext)
	s += fmt.Sprintf("file:\t\t%q\n", p.file)
	s += fmt.Sprintf("method:\t\t%q\n", p.method)

	for _, node := range p.nodes {
		s += fmt.Sprintf("node:\t\t%q\n", node)
	}

	s += fmt.Sprintf("poll:\t\t%v\n", p.poll)
	s += fmt.Sprintf("stypes:\t\t%v\n", p.stypes)
	s += fmt.Sprintf("test:\t\t%t\n", p.test)
	s += fmt.Sprintf("verbosity:\t%q\n", p.verbosity)
	s += fmt.Sprintln("}")

	return s
}

// STypes sets an array of legal types used by Source
func STypes(s []string) Option {
	return func(p *parms) Option {
		previous := p.stypes
		p.stypes = s
		return STypes(previous)
	}
}

// Test toggles testing mode on or off
func Test(b bool) Option {
	return func(p *parms) Option {
		previous := p.test
		p.test = b
		return Test(previous)
	}
}

// Verbosity sets the verbosity level to v
func Verbosity(i int) Option {
	return func(p *parms) Option {
		previous := p.verbosity
		p.verbosity = i
		return Verbosity(previous)
	}
}

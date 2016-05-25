package edgeos

import "time"

// Parms is struct of parameters
type Parms struct {
	debug     bool
	poll      time.Duration
	test      bool
	verbosity int
}

// Option sets is a recursive function
type Option func(p *Parms) Option

// SetOpt sets the specified options passed as parms and returns an option to restore the last arg's previous value
func (p *Parms) SetOpt(opts ...Option) (previous Option) {
	// apply all the options, and replace each with its inverse
	for i, opt := range opts {
		opts[i] = opt(p)
	}

	// reverse the list of inverses, since we want them to be applied in reverse order
	for i, j := 0, len(opts)-1; i <= j; i, j = i+1, j-1 {
		opts[i], opts[j] = opts[j], opts[i]
	}

	return func(p *Parms) Option {
		return p.SetOpt(opts...)
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

// Poll sets the polling interval in seconds
func Poll(t time.Duration) Option {
	return func(p *Parms) Option {
		previous := p.poll
		p.poll = t
		return Poll(previous)
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

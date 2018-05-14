package edgeos

import (
	"bufio"
	"bytes"
	"io"
	"sync"

	"github.com/britannic/blacklist/internal/regx"
)

// IFace type for labeling interface types
type IFace int

// IFace types for labeling Content interfaces
const (
	notfound       = -1
	Invalid  IFace = iota + 100
	ExDmObj
	ExHtObj
	ExRtObj
	FileObj
	FylDObj
	FylHObj
	PreDObj
	PreHObj
	PreRObj
	URLdObj
	URLhObj
)

type bList struct {
	file string
	r    io.Reader
	size int
}

// Contenter is an interface for handling the different file/http data sources
type Contenter interface {
	Find(string) int
	GetList() *Objects
	Len() int
	SetURL(string, string)
	String() string
}

// ExcDomnObjects struct of *Objects for domain exclusions
type ExcDomnObjects struct {
	*Objects
}

// ExcHostObjects struct of *Objects for host exclusions
type ExcHostObjects struct {
	*Objects
}

// ExcRootObjects struct of *Objects for global domain exclusions
type ExcRootObjects struct {
	*Objects
}

// FIODataObjects struct of *Objects for files
type FIODataObjects struct {
	*Objects
}

// FIODomnObjects struct of *Objects for files
type FIODomnObjects struct {
	*Objects
}

// FIOHostObjects struct of *Objects for files
type FIOHostObjects struct {
	*Objects
}

// PreDomnObjects struct of *Objects for pre-configured domains content
type PreDomnObjects struct {
	*Objects
}

// PreHostObjects struct of *Objects for pre-configured hosts content
type PreHostObjects struct {
	*Objects
}

// PreRootObjects struct of *Objects for pre-configured hosts content
type PreRootObjects struct {
	*Objects
}

// URLDomnObjects struct of *Objects for domain URLs
type URLDomnObjects struct {
	*Objects
}

// URLHostObjects struct of *Objects for host URLs
type URLHostObjects struct {
	*Objects
}

// Find returns the int position of an Objects' element
func (e *ExcDomnObjects) Find(s string) int {
	for i, o := range e.src {
		if o.name == s {
			return i
		}
	}
	return notfound
}

// Find returns the int position of an Objects' element
func (e *ExcHostObjects) Find(s string) int {
	for i, o := range e.src {
		if o.name == s {
			return i
		}
	}
	return notfound
}

// Find returns the int position of an Objects' element
func (e *ExcRootObjects) Find(s string) int {
	for i, o := range e.src {
		if o.name == s {
			return i
		}
	}
	return notfound
}

// Find returns the int position of an Objects' element
func (f *FIODataObjects) Find(s string) int {
	for i, o := range f.src {
		if o.name == s {
			return i
		}
	}
	return notfound
}

// Find returns the int position of an Objects' element
// func (f *FIODomnObjects) Find(s string) int {
// 	for i, o := range f.src {
// 		if o.name == s {
// 			return i
// 		}
// 	}
// 	return notfound
// }

// Find returns the int position of an Objects' element
// func (f *FIOHostObjects) Find(s string) int {
// 	for i, o := range f.src {
// 		if o.name == s {
// 			return i
// 		}
// 	}
// 	return notfound
// }

// Find returns the int position of an Objects' element
func (p *PreDomnObjects) Find(s string) int {
	for i, o := range p.src {
		if o.name == s {
			return i
		}
	}
	return notfound
}

// Find returns the int position of an Objects' element
func (p *PreHostObjects) Find(s string) int {
	for i, o := range p.src {
		if o.name == s {
			return i
		}
	}
	return notfound
}

// Find returns the int position of an Objects' element
func (p *PreRootObjects) Find(s string) int {
	for i, o := range p.src {
		if o.name == s {
			return i
		}
	}
	return notfound
}

// Find returns the int position of an Objects' element
func (u *URLHostObjects) Find(s string) int {
	for i, o := range u.src {
		if o.name == s {
			return i
		}
	}
	return notfound
}

// Find returns the int position of an Objects' element
func (u *URLDomnObjects) Find(s string) int {
	for i, o := range u.src {
		if o.name == s {
			return i
		}
	}
	return notfound
}

// GetList implements the Contenter interface for ExcDomnObjects
func (e *ExcDomnObjects) GetList() *Objects {
	for _, o := range e.src {
		if o.nType == excDomn {
			if o.exc != nil {
				o.r = o.excludes()
				o.Env = e.Env
			}
		}
	}
	return e.Objects
}

// GetList implements the Contenter interface for ExcHostObjects
func (e *ExcHostObjects) GetList() *Objects {
	for _, o := range e.src {
		if o.nType == excHost {
			if o.exc != nil {
				o.r = o.excludes()
				o.Env = e.Env
			}
		}
	}
	return e.Objects
}

// GetList implements the Contenter interface for ExcRootObjects
func (e *ExcRootObjects) GetList() *Objects {
	for _, o := range e.src {
		if o.nType == excRoot {
			if o.exc != nil {
				o.r = o.excludes()
				o.Env = e.Env
			}
		}
	}
	return e.Objects
}

// GetList implements the Contenter interface for FIODataObjects
func (f *FIODataObjects) GetList() *Objects {
	var responses = make(chan *source, len(f.src))

	for _, s := range f.src {
		s.Env = f.Env
		go func(s *source) {
			s.r, s.err = GetFile(s.file)
			responses <- s
		}(s)
	}

	for range f.src {
		response := <-responses
		f.src[f.Find(response.name)] = response
	}
	close(responses)
	return f.Objects
}

// GetList implements the Contenter interface for FIODataObjects
// func (f *FIODomnObjects) GetList() *Objects {
// 	var responses = make(chan *source, len(f.src))

// 	for _, s := range f.src {
// 		s.Env = f.Env
// 		go func(s *source) {
// 			s.r, s.err = GetFile(s.file)
// 			responses <- s
// 		}(s)
// 	}

// 	for range f.src {
// 		response := <-responses
// 		f.src[f.Find(response.name)] = response
// 	}
// 	close(responses)
// 	return f.Objects
// }

// GetList implements the Contenter interface for FIODataObjects
// func (f *FIOHostObjects) GetList() *Objects {
// 	var responses = make(chan *source, len(f.src))

// 	for _, s := range f.src {
// 		s.Env = f.Env
// 		go func(s *source) {
// 			s.r, s.err = GetFile(s.file)
// 			responses <- s
// 		}(s)
// 	}

// 	for range f.src {
// 		response := <-responses
// 		f.src[f.Find(response.name)] = response
// 	}
// 	close(responses)
// 	return f.Objects
// }

// GetList implements the Contenter interface for PreDomnObjects
func (p *PreDomnObjects) GetList() *Objects {
	for _, o := range p.src {
		if o.ltype == PreDomns && o.inc != nil {
			o.r = o.includes()
			o.Env = p.Env
		}
	}
	return p.Objects
}

// GetList implements the Contenter interface for PreHostObjects
func (p *PreHostObjects) GetList() *Objects {
	for _, o := range p.src {
		if o.ltype == PreHosts && o.inc != nil {
			o.r = o.includes()
			o.Env = p.Env
		}
	}
	return p.Objects
}

// GetList implements the Contenter interface for PreRootObjects
func (p *PreRootObjects) GetList() *Objects {
	for _, o := range p.src {
		if o.ltype == PreRoots && o.inc != nil {
			o.r = o.includes()
			o.Env = p.Env
		}
	}
	return p.Objects
}

// GetList implements the Contenter interface for URLHostObjects
func (u *URLDomnObjects) GetList() *Objects {
	var responses = make(chan *source, len(u.src))

	for _, s := range u.src {
		s.Env = u.Env
		go func(s *source) {
			responses <- download(s)
		}(s)
	}

	for range u.src {
		response := <-responses
		u.src[u.Find(response.name)] = response
	}
	close(responses)
	return u.Objects
}

// GetList implements the Contenter interface for URLHostObjects
func (u *URLHostObjects) GetList() *Objects {
	var responses = make(chan *source, len(u.src))

	for _, s := range u.src {
		s.Env = u.Env
		go func(s *source) {
			responses <- download(s)
		}(s)
	}

	for range u.src {
		response := <-responses
		u.src[u.Find(response.name)] = response

	}
	close(responses)
	return u.Objects
}

// Len returns how many sources there are
func (e *ExcDomnObjects) Len() int { return len(e.src) }

// Len returns how many sources there are
func (e *ExcHostObjects) Len() int { return len(e.src) }

// Len returns how many sources there are
func (e *ExcRootObjects) Len() int { return len(e.src) }

// Len returns how many sources there are
func (f *FIODataObjects) Len() int { return len(f.src) }

// Len returns how many sources there are
// func (f *FIODomnObjects) Len() int { return len(f.src) }

// Len returns how many sources there are
// func (f *FIOHostObjects) Len() int { return len(f.src) }

// Len returns how many sources there are
func (p *PreDomnObjects) Len() int { return len(p.src) }

// Len returns how many sources there are
func (p *PreHostObjects) Len() int { return len(p.src) }

// Len returns how many sources there are
func (p *PreRootObjects) Len() int { return len(p.src) }

// Len returns how many sources there are
func (u *URLDomnObjects) Len() int { return len(u.src) }

// Len returns how many sources there are
func (u *URLHostObjects) Len() int { return len(u.src) }

// Process extracts hosts/domains from downloaded raw content
func (s *source) process() *bList {
	var (
		l                        = list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
		area                     = typeInt(s.nType)
		b                        = bufio.NewScanner(s.r)
		dropped, extracted, kept int
		find                     = regx.NewRegex()
		ok                       bool
	)

	for b.Scan() {
		line := bytes.ToLower(bytes.TrimSpace(b.Bytes()))

		switch {
		case bytes.HasPrefix(line, []byte("#")), bytes.HasPrefix(line, []byte("//")), bytes.HasPrefix(line, []byte("<")):
			continue
		case bytes.HasPrefix(line, []byte(s.prefix)):
			if line, ok = find.StripPrefixAndSuffix(line, s.prefix); ok {
				for _, fqdn := range find.RX[regx.FQDN].FindAll(line, -1) {
					extracted++
					if s.Dex.subKeyExists(fqdn) {
						dropped++
						continue
					}
					if !s.Exc.keyExists(fqdn) {
						kept++
						s.Exc.set(fqdn)
						l.set(fqdn)
						continue
					}
					dropped++
				}
			}
		}
	}

	switch s.nType {
	case domn, excDomn, excRoot:
		s.Dex.merge(l)
	}

	s.sum(area, dropped, extracted, kept)

	return &bList{
		file: s.filename(area),
		r:    formatData(getDnsmasqPrefix(s), l),
		size: kept,
	}
}

// SetURL sets the Object's url field value
func (e *ExcDomnObjects) SetURL(name, url string) {
	for _, o := range e.src {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (e *ExcHostObjects) SetURL(name, url string) {
	for _, o := range e.src {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (e *ExcRootObjects) SetURL(name, url string) {
	for _, o := range e.src {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (f *FIODataObjects) SetURL(name, url string) {
	for _, o := range f.src {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
// func (f *FIODomnObjects) SetURL(name, url string) {
// 	for _, o := range f.src {
// 		if o.name == name {
// 			o.url = url
// 		}
// 	}
// }

// SetURL sets the Object's url field value
// func (f *FIOHostObjects) SetURL(name, url string) {
// 	for _, o := range f.src {
// 		if o.name == name {
// 			o.url = url
// 		}
// 	}
// }

// SetURL sets the Object's url field value
func (p *PreDomnObjects) SetURL(name, url string) {
	for _, o := range p.src {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (p *PreHostObjects) SetURL(name, url string) {
	for _, o := range p.src {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (p *PreRootObjects) SetURL(name, url string) {
	for _, o := range p.src {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (u *URLDomnObjects) SetURL(name, url string) {
	for _, o := range u.src {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (u *URLHostObjects) SetURL(name, url string) {
	for _, o := range u.src {
		if o.name == name {
			o.url = url
		}
	}
}

func (e *ExcDomnObjects) String() string { return e.Objects.String() }
func (e *ExcHostObjects) String() string { return e.Objects.String() }
func (e *ExcRootObjects) String() string { return e.Objects.String() }
func (f *FIODataObjects) String() string { return f.Objects.String() }

// func (f *FIODomnObjects) String() string { return f.Objects.String() }
// func (f *FIOHostObjects) String() string { return f.Objects.String() }
func (p *PreDomnObjects) String() string { return p.Objects.String() }
func (p *PreHostObjects) String() string { return p.Objects.String() }
func (p *PreRootObjects) String() string { return p.Objects.String() }
func (u *URLDomnObjects) String() string { return u.Objects.String() }
func (u *URLHostObjects) String() string { return u.Objects.String() }

func (i IFace) String() string {
	switch i {
	case ExDmObj:
		return ExcDomns
	case ExHtObj:
		return ExcHosts
	case ExRtObj:
		return ExcRoots
	case FileObj:
		return files
	case PreDObj:
		return PreDomns
	case PreHObj:
		return PreHosts
	case PreRObj:
		return PreRoots
	case URLhObj, URLdObj:
		return urls
	}
	return notknown
}

// func sourceIFace(n string) IFace {
//  if n == domains {
// 	 return PreDObj
//  }
// return PreHObj
// }

package edgeos

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/britannic/blacklist/internal/regx"
)

// IFace type for labeling interface types
type IFace int

// IFace types for labeling interface types
const (
	Invalid IFace = iota + 100
	ExDmObj
	ExHtObj
	ExRtObj
	FileObj
	PreDObj
	PreHObj
	URLdObj
	URLhObj
)

type bList struct {
	file string
	r    io.Reader
	size int
}

// Contenter is a Content interface
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

// PreDomnObjects struct of *Objects for pre-configured domains content
type PreDomnObjects struct {
	*Objects
}

// PreHostObjects struct of *Objects for pre-configured hosts content
type PreHostObjects struct {
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
func (e *ExcDomnObjects) Find(elem string) int {
	for i, o := range e.src {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (e *ExcHostObjects) Find(elem string) int {
	for i, o := range e.src {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (e *ExcRootObjects) Find(elem string) int {
	for i, o := range e.src {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (f *FIODataObjects) Find(elem string) int {
	for i, o := range f.src {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (p *PreDomnObjects) Find(elem string) int {
	for i, o := range p.src {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (p *PreHostObjects) Find(elem string) int {
	for i, o := range p.src {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (u *URLHostObjects) Find(elem string) int {
	for i, o := range u.src {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (u *URLDomnObjects) Find(elem string) int {
	for i, o := range u.src {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// GetList implements the Contenter interface for ExcDomnObjects
func (e *ExcDomnObjects) GetList() *Objects {
	for _, o := range e.src {
		if o.nType == excDomn {
			if o.exc != nil {
				o.r = o.excludes()
				o.Env = e.Objects.Env
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
				o.Env = e.Objects.Env
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
				o.Env = e.Objects.Env
			}
		}
	}
	return e.Objects
}

// GetList implements the Contenter interface for FIODataObjects
func (f *FIODataObjects) GetList() *Objects {
	var responses = make(chan *source, len(f.src))
	defer close(responses)

	for _, o := range f.src {
		o.Env = f.Objects.Env
		go func(o *source) {
			o.r, o.err = GetFile(o.file)
			responses <- o
		}(o)
	}

	for range f.src {
		response := <-responses
		f.src[f.Find(response.name)] = response
	}

	return f.Objects
}

// GetList implements the Contenter interface for PreDomnObjects
func (p *PreDomnObjects) GetList() *Objects {
	for _, o := range p.src {
		if o.ltype == PreDomns && o.inc != nil {
			o.r = o.includes()
			o.Env = p.Objects.Env
		}
	}
	return p.Objects
}

// GetList implements the Contenter interface for PreHostObjects
func (p *PreHostObjects) GetList() *Objects {
	for _, o := range p.src {
		if o.ltype == PreHosts && o.inc != nil {
			o.r = o.includes()
			o.Env = p.Objects.Env
		}
	}
	return p.Objects
}

// GetList implements the Contenter interface for URLHostObjects
func (u *URLDomnObjects) GetList() *Objects {
	var responses = make(chan *source, len(u.src))

	defer close(responses)

	for _, o := range u.src {
		o.Env = u.Objects.Env
		go func(o *source) {
			responses <- getHTTP(o)
		}(o)
	}

	for range u.src {
		response := <-responses
		u.src[u.Find(response.name)] = response
	}

	return u.Objects
}

// GetList implements the Contenter interface for URLHostObjects
func (u *URLHostObjects) GetList() *Objects {
	var responses = make(chan *source, len(u.src))

	defer close(responses)

	for _, o := range u.src {
		o.Env = u.Objects.Env
		go func(o *source) {
			responses <- getHTTP(o)
		}(o)
	}

	for range u.src {
		response := <-responses
		u.src[u.Find(response.name)] = response

	}

	return u.Objects
}

// Len returns how many sources there are
func (e *ExcDomnObjects) Len() int { return len(e.Objects.src) }

// Len returns how many sources there are
func (e *ExcHostObjects) Len() int { return len(e.Objects.src) }

// Len returns how many sources there are
func (e *ExcRootObjects) Len() int { return len(e.Objects.src) }

// Len returns how many sources there are
func (f *FIODataObjects) Len() int { return len(f.Objects.src) }

// Len returns how many sources there are
func (p *PreDomnObjects) Len() int { return len(p.Objects.src) }

// Len returns how many sources there are
func (p *PreHostObjects) Len() int { return len(p.Objects.src) }

// Len returns how many sources there are
func (u *URLDomnObjects) Len() int { return len(u.Objects.src) }

// Len returns how many sources there are
func (u *URLHostObjects) Len() int { return len(u.Objects.src) }

// Process extracts hosts/domains from downloaded raw content
func (o *source) process() *bList {
	var (
		add                   = list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
		area                  = typeInt(o.nType)
		b                     = bufio.NewScanner(o.r)
		drop, extracted, kept int
		find                  = regx.NewRegex()
		ok                    bool
	)

	for b.Scan() {
		line := bytes.TrimSpace(bytes.ToLower(b.Bytes()))

		switch {
		case bytes.HasPrefix(line, []byte("#")), bytes.HasPrefix(line, []byte("//")), bytes.HasPrefix(line, []byte("<")):
			continue
		case bytes.HasPrefix(line, []byte(o.prefix)):
			if line, ok = find.StripPrefixAndSuffix(line, o.prefix); ok {
				fqdns := find.RX[regx.FQDN].FindAll(line, -1)

				for _, fqdn := range fqdns {
					extracted++
					if o.Dex.subKeyExists(fqdn) {
						drop++
						continue
					}
					if !o.Exc.keyExists(fqdn) {
						kept++
						o.Exc.set(fqdn, 0)
						add.set(fqdn, 0)
						continue
					}
					drop++
				}
			}
		}
	}

	switch o.nType {
	case domn, excDomn, excRoot:
		o.Dex = mergeList(o.Dex, add)
	}

	// Let's do some accounting
	atomic.AddInt32(&o.ctr[area].dropped, int32(drop))
	atomic.AddInt32(&o.ctr[area].extracted, int32(extracted))
	atomic.AddInt32(&o.ctr[area].kept, int32(kept))

	if kept != 0 && area != rootNode {
		o.Log.Infof("%s: downloaded: %d", o.name, extracted)
		o.Log.Infof("%s: extracted: %d", o.name, kept)
		o.Log.Infof("%s: dropped: %d", o.name, drop)
	}

	return &bList{
		file: o.filename(area),
		r:    formatData(getDnsmasqPrefix(o), add),
		size: kept,
	}
}

func (o *source) filename(area string) string {
	switch o.nType {
	case excDomn, excRoot, preDomn, preRoot:
		return fmt.Sprintf(o.FnFmt, o.Dir, domains, o.name, o.Ext)
	case excHost, preHost:
		return fmt.Sprintf(o.FnFmt, o.Dir, hosts, o.name, o.Ext)
	}
	return fmt.Sprintf(o.FnFmt, o.Dir, area, o.name, o.Ext)
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
func (p *PreDomnObjects) String() string { return p.Objects.String() }
func (p *PreHostObjects) String() string { return p.Objects.String() }
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
	case URLhObj, URLdObj:
		return urls
	}
	return notknown
}

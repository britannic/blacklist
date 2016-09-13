package edgeos

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/britannic/blacklist/internal/regx"
)

// IFace type for labeling interface types
type IFace int

// IFace types for labeling interface types
const (
	Invalid IFace = iota + 100
	ExRtObj
	ExDmObj
	ExHtObj
	FileObj
	PreDObj
	PreHObj
	URLdObj
	URLhObj
)

type bList struct {
	file string
	r    io.Reader
}

// Contenter is a Content interface
type Contenter interface {
	Find(string) int
	GetList() *Objects
	Len() int
	SetURL(string, string)
	String() string
}

// ExcDomnObjects implements GetList for domain exclusions
type ExcDomnObjects struct {
	*Objects
}

// ExcHostObjects implements GetList for host exclusions
type ExcHostObjects struct {
	*Objects
}

// ExcRootObjects implements GetList for global domain exclusions
type ExcRootObjects struct {
	*Objects
}

// FIODataObjects implements GetList for files
type FIODataObjects struct {
	*Objects
}

// PreDomnObjects implements GetList for pre-configured domains content
type PreDomnObjects struct {
	*Objects
}

// PreHostObjects implements GetList for pre-configured hosts content
type PreHostObjects struct {
	*Objects
}

// URLDomnObjects implements GetList for URLs
type URLDomnObjects struct {
	*Objects
}

// URLHostObjects implements GetList for URLs
type URLHostObjects struct {
	*Objects
}

// Find returns the int position of an Objects' element
func (e *ExcDomnObjects) Find(elem string) int {
	for i, o := range e.x {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (e *ExcHostObjects) Find(elem string) int {
	for i, o := range e.x {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (e *ExcRootObjects) Find(elem string) int {
	for i, o := range e.x {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (f *FIODataObjects) Find(elem string) int {
	for i, o := range f.x {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (p *PreDomnObjects) Find(elem string) int {
	for i, o := range p.x {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (p *PreHostObjects) Find(elem string) int {
	for i, o := range p.x {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (u *URLHostObjects) Find(elem string) int {
	for i, o := range u.x {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (u *URLDomnObjects) Find(elem string) int {
	for i, o := range u.x {
		if o.name == elem {
			return i
		}
	}
	return -1
}

// GetList implements the Contenter interface for ExcDomnObjects
func (e *ExcDomnObjects) GetList() *Objects {
	for _, o := range e.x {
		switch o.nType {
		case excDomn:
			if o.exc != nil {
				o.r = o.excludes()
				o.Parms = e.Objects.Parms
			}
		}
	}
	return e.Objects
}

// GetList implements the Contenter interface for ExcHostObjects
func (e *ExcHostObjects) GetList() *Objects {
	for _, o := range e.x {
		switch o.nType {
		case excHost:
			if o.exc != nil {
				o.r = o.excludes()
				o.Parms = e.Objects.Parms
			}
		}
	}
	return e.Objects
}

// GetList implements the Contenter interface for ExcRootObjects
func (e *ExcRootObjects) GetList() *Objects {
	for _, o := range e.x {
		switch o.nType {
		case excRoot:
			if o.exc != nil {
				o.r = o.excludes()
				o.Parms = e.Objects.Parms
			}
		}
	}
	return e.Objects
}

// GetList implements the Contenter interface for FIODataObjects
func (f *FIODataObjects) GetList() *Objects {
	var responses = make(chan *object, len(f.x))

	defer close(responses)

	for _, o := range f.x {
		o.Parms = f.Objects.Parms
		go func(o *object) {
			o.r, o.err = getFile(o.file)
			responses <- o
		}(o)
	}

	for _ = range Iter(len(f.x)) {
		select {
		case response := <-responses:
			f.x[f.Find(response.name)] = response
		}
	}

	return f.Objects
}

// GetList implements the Contenter interface for PreDomnObjects
func (p *PreDomnObjects) GetList() *Objects {
	for _, o := range p.x {
		if o.ltype == PreDomns && o.inc != nil {
			o.r = o.includes()
			o.Parms = p.Objects.Parms
		}
	}
	return p.Objects
}

// GetList implements the Contenter interface for PreHostObjects
func (p *PreHostObjects) GetList() *Objects {
	for _, o := range p.x {
		if o.ltype == PreHosts && o.inc != nil {
			o.r = o.includes()
			o.Parms = p.Objects.Parms
		}
	}
	return p.Objects
}

// GetList implements the Contenter interface for URLHostObjects
func (u *URLDomnObjects) GetList() *Objects {
	var responses = make(chan *object, len(u.x))

	defer close(responses)

	for _, o := range u.x {
		o.Parms = u.Objects.Parms
		go func(o *object) {
			responses <- getHTTP(o)
		}(o)
	}

	for i := 0; i < len(u.x); i++ {
		select {
		case response := <-responses:
			u.x[u.Find(response.name)] = response
		}
	}

	return u.Objects
}

// GetList implements the Contenter interface for URLHostObjects
func (u *URLHostObjects) GetList() *Objects {
	var responses = make(chan *object, len(u.x))

	defer close(responses)

	for _, o := range u.x {
		o.Parms = u.Objects.Parms
		go func(o *object) {
			responses <- getHTTP(o)
		}(o)
	}

	for i := 0; i < len(u.x); i++ {
		select {
		case response := <-responses:
			u.x[u.Find(response.name)] = response
		}
	}

	return u.Objects
}

// Len returns how many objects there are
func (e *ExcDomnObjects) Len() int { return len(e.Objects.x) }

// Len returns how many objects there are
func (e *ExcHostObjects) Len() int { return len(e.Objects.x) }

// Len returns how many objects there are
func (e *ExcRootObjects) Len() int { return len(e.Objects.x) }

// Len returns how many objects there are
func (f *FIODataObjects) Len() int { return len(f.Objects.x) }

// Len returns how many objects there are
func (p *PreDomnObjects) Len() int { return len(p.Objects.x) }

// Len returns how many objects there are
func (p *PreHostObjects) Len() int { return len(p.Objects.x) }

// Len returns how many objects there are
func (u *URLDomnObjects) Len() int { return len(u.Objects.x) }

// Len returns how many objects there are
func (u *URLHostObjects) Len() int { return len(u.Objects.x) }

// Process extracts hosts/domains from downloaded raw content
func (o *object) process(m chan *Msg) *bList {
	var (
		add = list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
		b   = bufio.NewScanner(o.r)
		d   = NewMsg(o.Name)
		rx  = regx.Obj
	)

NEXT:
	for b.Scan() {
		line := strings.ToLower(b.Text())
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "#"), strings.HasPrefix(line, "//"):
			continue NEXT

		case strings.HasPrefix(line, o.prefix):
			var ok bool

			if line, ok = rx.StripPrefixAndSuffix(line, o.prefix); ok {
				fqdns := rx.FQDN.FindAllString(line, -1)

			FQDN:
				for _, fqdn := range fqdns {
					isDEX := o.Dex.subKeyExists(fqdn)
					isEXC := o.Exc.keyExists(fqdn)

					switch {
					case isDEX:
						// o.Dex.inc(fqdn)
						d.incDupe()
						m <- d
						continue FQDN

					case isEXC:
						if add.keyExists(fqdn) {
							// add.inc(fqdn)
							d.incDupe()
							m <- d
						}
						// o.Exc.inc(fqdn)

					case !isEXC:
						o.Exc.set(fqdn, 0)
						add.set(fqdn, 0)
						d.incNew()
						m <- d
					}
				}
			}
		default:
			continue NEXT
		}
	}

	switch o.nType {
	case domn, excDomn, excRoot:
		o.Dex = mergeList(o.Dex, add)
	}

	fmttr := o.Pfx + getSeparator(getType(o.nType).(string)) + "%v/" + o.ip

	d.Done = true
	m <- d

	return &bList{
		file: fmt.Sprintf(o.FnFmt, o.Dir, getType(o.nType).(string), o.name, o.Ext),
		r:    formatData(fmttr, add),
	}
}

// ProcessContent processes the Contents array
func (c *Config) ProcessContent(m chan *Msg, cts ...Contenter) error {
	var (
		errs      []string
		getErrors chan error
	)

	if len(cts) < 1 {
		return errors.New("Empty Contenter interface{} passed to ProcessContent()")
	}

	for _, ct := range cts {
		for _, o := range ct.GetList().x {
			getErrors = make(chan error)
			_ = NewMsg(o.Name) //TODO

			if o.err != nil {
				errs = append(errs, o.err.Error())
			}

			go func(o *object) {
				switch o.nType {
				case excDomn, excHost, excRoot:
					o.process(m)
					getErrors <- nil
				default:
					getErrors <- o.process(m).writeFile()
				}
			}(o)

			for i := 0; i < len(cts); i++ {
				select {
				case err := <-getErrors:
					if err != nil {
						errs = append(errs, err.Error())
					}
				}
			}
		}
		close(getErrors)
	}

	if errs != nil {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}

	return nil
}

// SetURL sets the Object's url field value
func (e *ExcDomnObjects) SetURL(name, url string) {
	for _, o := range e.x {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (e *ExcHostObjects) SetURL(name, url string) {
	for _, o := range e.x {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (e *ExcRootObjects) SetURL(name, url string) {
	for _, o := range e.x {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (f *FIODataObjects) SetURL(name, url string) {
	for _, o := range f.x {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (p *PreDomnObjects) SetURL(name, url string) {
	for _, o := range p.x {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (p *PreHostObjects) SetURL(name, url string) {
	for _, o := range p.x {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (u *URLDomnObjects) SetURL(name, url string) {
	for _, o := range u.x {
		if o.name == name {
			o.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (u *URLHostObjects) SetURL(name, url string) {
	for _, o := range u.x {
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

func (i IFace) String() (s string) {
	switch i {
	case ExDmObj:
		s = ExcDomns
	case ExHtObj:
		s = ExcHosts
	case ExRtObj:
		s = ExcRoots
	case FileObj:
		s = files
	case PreDObj:
		s = PreDomns
	case PreHObj:
		s = PreHosts
	case URLhObj, URLdObj:
		s = urls
	default:
		s = notknown
	}
	return s
}

// writeFile saves hosts/domains data to disk
func (b *bList) writeFile() error {
	w, err := os.Create(b.file)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, b.r)
	return err
}

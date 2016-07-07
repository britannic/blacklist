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

type blist struct {
	file string
	r    io.Reader
}

// Contenter is a Content interface
type Contenter interface {
	Find(elem string) int
	GetList() *Objects
	SetURL(name string, url string)
	String() string
}

// FIODataObjects implements GetList for files
type FIODataObjects struct {
	*Objects
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

// PreDomnObjects implements GetList for pre-configured domains content
type PreDomnObjects struct {
	*Objects
}

// PreHostObjects implements GetList for pre-configured hosts content
type PreHostObjects struct {
	*Objects
}

// URLHostObjects implements GetList for URLs
type URLHostObjects struct {
	*Objects
}

// URLDomnObjects implements GetList for URLs
type URLDomnObjects struct {
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
				o.err = nil
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
				o.err = nil
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
	for _, o := range f.x {
		if o.ltype == files && o.file != "" {
			o.r, o.err = getFile(o.file)
			o.Parms = f.Objects.Parms
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
	var wg sync.WaitGroup
	wg.Add(len(u.x))
	responses := make(chan *object, len(u.x))

	for _, o := range u.x {
		o.Parms = u.Objects.Parms
		go func(o *object) {
			defer wg.Done()
			responses <- getHTTP(o)
		}(o)
	}

	go func() {
		for response := range responses {
			u.x[u.Find(response.name)] = response
		}
	}()

	wg.Wait()
	return u.Objects
}

// GetList implements the Contenter interface for URLHostObjects
func (u *URLHostObjects) GetList() *Objects {
	var wg sync.WaitGroup
	wg.Add(len(u.x))
	responses := make(chan *object, len(u.x))

	for _, o := range u.x {
		o.Parms = u.Objects.Parms
		go func(o *object) {
			defer wg.Done()
			responses <- getHTTP(o)
		}(o)
	}

	go func() {
		for response := range responses {
			u.x[u.Find(response.name)] = response
		}
	}()

	wg.Wait()
	return u.Objects
}

// Process extracts hosts/domains from downloaded raw content
func (o *object) process() *blist {
	var (
		b   = bufio.NewScanner(o.r)
		rx  = regx.Obj
		add = list{RWMutex: &sync.RWMutex{}, entry: make(entry)}
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
					isEX := o.Exc.keyExists(fqdn)

					switch {
					case isDEX:
						o.Dex.inc(fqdn)
						continue FQDN

					case isEX:
						if add.keyExists(fqdn) {
							add.inc(fqdn)
						}
						o.Exc.inc(fqdn)

					case !isEX:
						o.Exc.set(fqdn, 0)
						add.set(fqdn, 0)
					}
				}
			}
		default:
			continue NEXT
		}
	}

	switch o.nType {
	case domn, excDomn, excRoot:
		o.Parms.Dex = mergeList(o.Parms.Dex, add)
	}

	fmttr := o.Parms.Pfx + getSeparator(getType(o.nType).(string)) + "%v/" + o.ip

	return &blist{
		file: fmt.Sprintf(o.Parms.FnFmt, o.Parms.Dir, getType(o.nType).(string), o.name, o.Parms.Ext),
		r:    formatData(fmttr, add),
	}
}

// ProcessContent processes the Contents array
func (c *Config) ProcessContent(cts ...Contenter) error {
	var (
		errs      []string
		getErrors = make(chan error)
		wg        sync.WaitGroup
	)

	for _, ct := range cts {
		wg.Add(len(ct.GetList().x))
		for _, o := range ct.GetList().x {
			go func(o *object) {
				defer wg.Done()
				switch o.nType {
				case excDomn, excHost, excRoot:
					o.process()
				default:
					getErrors <- o.process().writeFile()
				}
			}(o)
		}
	}

	go func(errs []string) {
		for err := range getErrors {
			errs = append(errs, fmt.Sprint(err))
		}
	}(errs)

	wg.Wait()

	if errs != nil {
		return errors.New(strings.Join(errs, "\n"))
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
func (b *blist) writeFile() error {
	w, err := os.Create(b.file)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, b.r)
	return err
}

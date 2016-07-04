package edgeos

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

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
	URLsObj
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

// URLDataObjects implements GetList for URLs
type URLDataObjects struct {
	*Objects
}

// Find returns the int position of an Objects' element
func (e *ExcDomnObjects) Find(elem string) int {
	for i, obj := range e.x {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (e *ExcHostObjects) Find(elem string) int {
	for i, obj := range e.x {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (e *ExcRootObjects) Find(elem string) int {
	for i, obj := range e.x {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (f *FIODataObjects) Find(elem string) int {
	for i, obj := range f.x {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (p *PreDomnObjects) Find(elem string) int {
	for i, obj := range p.x {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (p *PreHostObjects) Find(elem string) int {
	for i, obj := range p.x {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element
func (u *URLDataObjects) Find(elem string) int {
	for i, obj := range u.x {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

// GetList implements the Contenter interface for ExcDomnObjects
func (e *ExcDomnObjects) GetList() *Objects {
	// var c Contents
	for _, obj := range e.x {
		switch obj.nType {
		case excDomn:
			if obj.exc != nil {
				obj.r = obj.excludes()
				obj.err = nil
				obj.Parms = e.Objects.Parms
			}
		}
	}
	return e.Objects
}

// GetList implements the Contenter interface for ExcHostObjects
func (e *ExcHostObjects) GetList() *Objects {
	for _, obj := range e.x {
		switch obj.nType {
		case excHost:
			if obj.exc != nil {
				obj.r = obj.excludes()
				obj.err = nil
				obj.Parms = e.Objects.Parms
			}
		}
	}
	return e.Objects
}

// GetList implements the Contenter interface for ExcRootObjects
func (e *ExcRootObjects) GetList() *Objects {
	for _, obj := range e.x {
		switch obj.nType {
		case excRoot:
			if obj.exc != nil {
				obj.r = obj.excludes()
				obj.Parms = e.Objects.Parms
			}
		}
	}
	return e.Objects
}

// GetList implements the Contenter interface for FIODataObjects
func (f *FIODataObjects) GetList() *Objects {
	for _, obj := range f.x {
		if obj.ltype == files && obj.file != "" {
			obj.r, obj.err = getFile(obj.file)
			obj.Parms = f.Objects.Parms
		}
	}
	return f.Objects
}

// GetList implements the Contenter interface for PreDomnObjects
func (p *PreDomnObjects) GetList() *Objects {
	for _, obj := range p.x {
		if obj.ltype == PreDomns && obj.inc != nil {
			obj.r = obj.includes()
			obj.Parms = p.Objects.Parms
		}
	}
	return p.Objects
}

// GetList implements the Contenter interface for PreHostObjects
func (p *PreHostObjects) GetList() *Objects {
	for _, obj := range p.x {
		if obj.ltype == PreHosts && obj.inc != nil {
			obj.r = obj.includes()
			obj.Parms = p.Objects.Parms
		}
	}
	return p.Objects
}

// GetList implements the Contenter interface for URLDataObjects
func (u *URLDataObjects) GetList() *Objects {
	for _, obj := range u.x {
		if obj.ltype == urls && obj.url != "" {
			obj.r, obj.err = getHTTP(u.Parms.Method, obj.url)
			obj.Parms = u.Objects.Parms
		}
	}
	return u.Objects
}

// Process extracts hosts/domains from downloaded raw content
func (o *object) process() *blist {
	var (
		b     = bufio.NewScanner(o.r)
		rx    = regx.Obj
		sList = make(list)
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
					isDEX := o.Parms.Dex.subKeyExists(fqdn)
					isEX := o.Parms.Exc.keyExists(fqdn)

					switch {
					case isDEX:
						o.Parms.Dex[fqdn]++
						continue FQDN

					case isEX:
						if sList.keyExists(fqdn) {
							sList[fqdn]++
						}
						o.Parms.Exc[fqdn]++

					case !isEX:
						o.Parms.Exc[fqdn] = 0
						sList[fqdn] = 0
					}
				}
			}
		default:
			continue NEXT
		}
	}

	switch o.nType {
	case domn, excDomn, excRoot:
		o.Parms.Dex = mergeList(o.Parms.Dex, sList)
	}

	fmttr := o.Parms.Pfx + getSeparator(getType(o.nType).(string)) + "%v/" + o.ip

	return &blist{
		file: fmt.Sprintf(o.Parms.FnFmt, o.Parms.Dir, getType(o.nType).(string), o.name, o.Parms.Ext),
		r:    formatData(fmttr, sList),
	}
}

// ProcessContent processes the Contents array
func (c *Config) ProcessContent(cts ...Contenter) error {
	var errs []string
	for _, ct := range cts {
		for _, src := range ct.GetList().x {
			switch src.nType {
			case excDomn, excHost, excRoot:
				src.process()

			default:
				if err := src.process().writeFile(); err != nil {
					errs = append(errs, fmt.Sprint(err))
				}
			}
		}
	}

	if errs != nil {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

// SetURL sets the Object's url field value
func (e *ExcDomnObjects) SetURL(name, url string) {
	for _, obj := range e.x {
		if obj.name == name {
			obj.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (e *ExcHostObjects) SetURL(name, url string) {
	for _, obj := range e.x {
		if obj.name == name {
			obj.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (e *ExcRootObjects) SetURL(name, url string) {
	for _, obj := range e.x {
		if obj.name == name {
			obj.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (f *FIODataObjects) SetURL(name, url string) {
	for _, obj := range f.x {
		if obj.name == name {
			obj.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (p *PreDomnObjects) SetURL(name, url string) {
	for _, obj := range p.x {
		if obj.name == name {
			obj.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (p *PreHostObjects) SetURL(name, url string) {
	for _, obj := range p.x {
		if obj.name == name {
			obj.url = url
		}
	}
}

// SetURL sets the Object's url field value
func (u *URLDataObjects) SetURL(name, url string) {
	for _, obj := range u.x {
		if obj.name == name {
			obj.url = url
		}
	}
}

func (e *ExcDomnObjects) String() string { return e.Objects.String() }
func (e *ExcHostObjects) String() string { return e.Objects.String() }
func (e *ExcRootObjects) String() string { return e.Objects.String() }
func (f *FIODataObjects) String() string { return f.Objects.String() }
func (p *PreDomnObjects) String() string { return p.Objects.String() }
func (p *PreHostObjects) String() string { return p.Objects.String() }
func (u *URLDataObjects) String() string { return u.Objects.String() }

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
	case URLsObj:
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

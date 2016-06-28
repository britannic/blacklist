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

// iFace type for labeling interface types
type iFace int

// iFace types for labeling interface types
const (
	Invalid iFace = iota + 100
	ContObj
	FileObj
	PreDomnObj
	PreHostObj
	URLObj
)

const cntnt = "contents"

type blist struct {
	file string
	r    io.Reader
}

// Contenter is a Content interface
type Contenter interface {
	Find(elem string) int
	GetBlacklist() *Contents
	SetURL(name string, url string)
}

// Content is a struct of blacklist content
type Content struct {
	*Object
	*Parms
	err error
	r   io.Reader
}

// FileObjects implements GetBlacklist for files
type FileObjects struct {
	*Objects
}

// PreDomainObjects implements GetBlacklist for pre-configured domains content
type PreDomainObjects struct {
	*Objects
}

// PreHostObjects implements GetBlacklist for pre-configured hosts content
type PreHostObjects struct {
	*Objects
}

// URLObjects implements GetBlacklist for URLs
type URLObjects struct {
	*Objects
}

// Contents is an array of *content
type Contents []*Content

// Find returns the int position of an Objects' element in the StringSlice
func (c *Contents) Find(elem string) int {
	if elem == cntnt {
		return 0
	}
	return -1
}

// Find returns the int position of an Objects' element in the StringSlice
func (f *FileObjects) Find(elem string) int {
	for i, obj := range f.S {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element in the StringSlice
func (p *PreDomainObjects) Find(elem string) int {
	for i, obj := range p.S {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element in the StringSlice
func (p *PreHostObjects) Find(elem string) int {
	for i, obj := range p.S {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

// Find returns the int position of an Objects' element in the StringSlice
func (u *URLObjects) Find(elem string) int {
	for i, obj := range u.S {
		if obj.name == elem {
			return i
		}
	}
	return -1
}

// GetBlacklist implements a dummy Contenter interface for Contents
func (c *Contents) GetBlacklist() *Contents {
	return c
}

// GetBlacklist implements the Contenter interface for FileObjects
func (f *FileObjects) GetBlacklist() *Contents {
	var c Contents
	for _, obj := range f.S {
		if obj.ltype == files && obj.file != "" {
			b, err := getFile(obj.file)
			c = append(c, &Content{
				err:    err,
				Object: obj,
				Parms:  f.Parms,
				r:      b,
			})
		}
	}
	return &c
}

// GetBlacklist implements the Contenter interface for PreDomainObjects
func (p *PreDomainObjects) GetBlacklist() *Contents {
	var c Contents
	for _, obj := range p.S {
		if obj.ltype == PreDomns && obj.inc != nil {
			c = append(c, &Content{
				err:    nil,
				Object: obj,
				Parms:  p.Parms,
				r:      obj.Includes(),
			})
		}
	}
	return &c
}

// GetBlacklist implements the Contenter interface for PreHostObjects
func (p *PreHostObjects) GetBlacklist() *Contents {
	var c Contents
	for _, obj := range p.S {
		if obj.ltype == PreHosts && obj.inc != nil {
			c = append(c, &Content{
				err:    nil,
				Object: obj,
				Parms:  p.Parms,
				r:      obj.Includes(),
			})
		}
	}
	return &c
}

// GetBlacklist implements the Contenter interface for URLObjects
func (u *URLObjects) GetBlacklist() *Contents {
	var c Contents
	for _, obj := range u.S {
		if obj.ltype == urls && obj.url != "" {
			reader, err := GetHTTP(u.Parms.Method, obj.url)
			c = append(c, &Content{
				err:    err,
				Object: obj,
				Parms:  u.Parms,
				r:      reader,
			})
		}
	}
	return &c
}

// Process extracts hosts/domains from downloaded raw content
func (c *Content) process() *blist {
	var (
		b     = bufio.NewScanner(c.r)
		rx    = regx.Objects
		sList = make(List)
	)

NEXT:
	for b.Scan() {
		line := strings.ToLower(b.Text())
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "#"), strings.HasPrefix(line, "//"):
			continue NEXT

		case strings.HasPrefix(line, c.prefix):
			var ok bool

			if line, ok = rx.StripPrefixAndSuffix(line, c.prefix); ok {
				fqdns := rx.FQDN.FindAllString(line, -1)

			FQDN:
				for _, fqdn := range fqdns {
					isDEX := c.Parms.Dex.subKeyExists(fqdn)
					isEX := c.Parms.Exc.keyExists(fqdn)

					switch {
					case isDEX:
						c.Parms.Dex[fqdn]++
						continue FQDN

					case isEX:
						if sList.keyExists(fqdn) {
							sList[fqdn]++
						}
						c.Parms.Exc[fqdn]++

					case !isEX:
						c.Parms.Exc[fqdn] = 0
						sList[fqdn] = 0
					}
				}
			}
		default:
			continue NEXT
		}
	}

	if c.nType == domain {
		c.Parms.Dex = mergeList(c.Parms.Dex, sList)
	}

	fmttr := c.Parms.Pfx + getSeparator(getType(c.nType).(string)) + "%v/" + c.ip

	return &blist{
		file: fmt.Sprintf(c.Parms.FnFmt, c.Parms.Dir, getType(c.nType).(string), c.name, c.Parms.Ext),
		r:    formatData(fmttr, sList),
	}
}

// ProcessContent processes the Contents array
func (c *Config) ProcessContent(ct Contenter) error {
	var errs []string
	for _, src := range *ct.GetBlacklist() {
		if err := src.process().WriteFile(); err != nil {
			errs = append(errs, fmt.Sprint(err))
		}
	}

	if errs != nil {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

// SetURL sets a field value
func (f *FileObjects) SetURL(name, url string) {
	for _, obj := range f.S {
		if obj.name == name {
			obj.url = url
		}
	}
}

// SetURL sets a field value
func (p *PreDomainObjects) SetURL(name, url string) {
	for _, obj := range p.S {
		if obj.name == name {
			obj.url = url
		}
	}
}

// SetURL sets a field value
func (p *PreHostObjects) SetURL(name, url string) {
	for _, obj := range p.S {
		if obj.name == name {
			obj.url = url
		}
	}
}

// SetURL sets a field value
func (u *URLObjects) SetURL(name, url string) {
	for _, obj := range u.S {
		if obj.name == name {
			obj.url = url
		}
	}
}

// SetURL sets a field value
func (c *Contents) SetURL(name, url string) {
	if name == cntnt && url == cntnt {
		_ = true
	}
}

func (f *FileObjects) String() string      { return f.Objects.String() }
func (p *PreDomainObjects) String() string { return p.Objects.String() }
func (p *PreHostObjects) String() string   { return p.Objects.String() }
func (u *URLObjects) String() string       { return u.Objects.String() }

func (i iFace) String() (s string) {
	switch i {
	case ContObj:
		s = cntnt
	case FileObj:
		s = files
	case PreDomnObj:
		s = PreDomns
	case PreHostObj:
		s = PreHosts
	case URLObj:
		s = urls
	default:
		s = notknown
	}
	return s
}

// WriteFile saves hosts/domains data to disk
func (b *blist) WriteFile() error {
	w, err := os.Create(b.file)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, b.r)
	return err
}

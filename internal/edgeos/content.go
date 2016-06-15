package edgeos

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/britannic/blacklist/internal/regx"
)

type blist struct {
	file string
	r    io.Reader
}

// Contenter is a Content interface
type Contenter interface {
	Process() io.Reader
}

// Content is a struct of blacklist content
type Content struct {
	*Object
	Contenter
	err error
	*Parms
	r io.Reader
}

// Contents is an array of *content
type Contents []*Content

// GetContent returns a Content struct
func (objs *Objects) GetContent() *Contents {
	var c Contents
	for _, o := range objs.S {
		switch o.ltype {
		case preConf:
			if o.inc != nil {
				c = append(c, &Content{
					err:    nil,
					Object: o,
					Parms:  o.Parms,
					r:      o.Includes(),
				})
			}

		case "files":
			if o.file != "" {
				b, err := getFile(o.file)
				c = append(c, &Content{
					err:    err,
					Object: o,
					Parms:  o.Parms,
					r:      b,
				})
			}

		case "urls":
			if o.url != "" {
				reader, err := GetHTTP(o.Parms.Method, o.url)
				c = append(c, &Content{
					err:    err,
					Object: o,
					Parms:  o.Parms,
					r:      reader,
				})
			}
		}
	}
	return &c
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

// Process extracts hosts/domains from downloaded raw content
func (c *Content) Process() *blist {
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
						// isList := sList.keyExists(fqdn)
						if sList.keyExists(fqdn) {
							sList[fqdn]++
						}
						c.Parms.Exc[fqdn]++

					// case isList:
					// 	sList[fqdn]++

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

// ProcessContent iterates through the Contents array and processes each
func (c *Contents) ProcessContent() {
	for _, src := range *c {
		if err := src.Process().WriteFile(); err != nil {
			log.Println(err)
		}
	}
}

func (c *Contents) String() (result string) {
	for _, src := range *c {
		b := src.Process()
		got, _ := ioutil.ReadAll(b.r)
		result += fmt.Sprintf("File: %#v\nData:\n%v\n", b.file, string(got))
	}
	return result
}

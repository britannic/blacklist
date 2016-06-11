package edgeos

import (
	"bufio"
	"io"
	"io/ioutil"
	"strings"

	"github.com/britannic/blacklist/internal/regx"
)

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
					r:      o.Includes(),
				})
			}

		case "files":
			if o.file != "" {
				b, err := getFile(o.file)
				c = append(c, &Content{
					err:    err,
					Object: o,
					r:      b,
				})
			}

		case "urls":
			if o.url != "" {
				reader, err := GetHTTP(o.Parms.method, o.url)
				c = append(c, &Content{
					err:    err,
					Object: o,
					r:      reader,
				})
			}
		}
	}
	return &c
}

// WriteFile saves hosts/domains data to disk
func (c *Contents) WriteFile() (err error) {
	for _, content := range *c {
		var b []byte
		if b, err = ioutil.ReadAll(content.process()); err != nil {
			return err
		}
		fname := content.String()
		return ioutil.WriteFile(fname, b, 0644)
	}
	return err
}

// process extracts hosts/domains from downloaded raw content
func (c *Content) process() io.Reader {
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

			if line, ok = stripPrefixAndSuffix(line, c.prefix, rx); ok {
				fqdns := rx.FQDN.FindAllString(line, -1)

			FQDN:
				for _, fqdn := range fqdns {
					isDEX := c.dex.subKeyExists(fqdn)
					isEX := c.Parms.exc.keyExists(fqdn)
					isList := sList.keyExists(fqdn)

					switch {
					case isDEX:
						continue FQDN

					case isEX:
						if isList {
							sList[fqdn]++
						}
						c.Parms.exc[fqdn]++

					case isList:
						sList[fqdn]++

					case !isEX:
						c.Parms.exc[fqdn] = 0
						sList[fqdn] = 0
					}
				}
			}
		default:
			continue NEXT
		}
	}

	if c.nType == domain {
		c.dex = mergeList(c.dex, sList)
	}

	fmttr := "address=" + getSeparator(getType(c.nType).(string)) + "%v/" + c.ip
	return formatData(fmttr, sList)
}

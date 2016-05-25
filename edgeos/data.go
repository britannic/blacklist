package edgeos

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/britannic/blacklist/regx"
)

// Bnodes determines which nodes are used to look up configuration data
var Bnodes = []string{Root, Domains, Hosts}

// Configure has methods for returning config data supersets
// type Configure interface {
// 	Disabled(node string) bool
// 	Excludes(node string) []string
// 	Files() []string
// 	FormatData(node string, data []string) (reader io.Reader, list List, err error)
// 	Get(node string) (e *EdgeOS)
// 	Includes(node string) []string
// 	IP(node string) string
// 	SourceFiles(node string) (fSrcs []string)
// 	Sources(node string) []Srcs
// 	SourceURLs(node string) (fSrcs []string)
// }

// Object is a map of *EdgeOS
type Object map[string]*EdgeOS

// Config is a struct for configuring the EdgeOS methods
type Config struct {
	ex, dex List
	o       Object
}

// EdgeOS struct for normalizing EdgeOS data.
type EdgeOS struct {
	Disabled bool
	Exc      []string
	Inc      []string
	IP       string
	Nodes    []Leaf
	Sources  []Srcs
}

// DiffArray returns the delta of two arrays
func DiffArray(a, b []string) (diff []string) {
	var biggest, smallest []string

	switch {
	case len(a) > len(b), len(a) == len(b):
		biggest = a
		smallest = b

	case len(a) < len(b):
		biggest = b
		smallest = a
	}

	dmap := make(List)
	for _, k := range smallest {
		dmap[k] = 0
	}

	for _, k := range biggest {
		if !dmap.KeyExists(k) {
			diff = append(diff, k)
		}
	}

	sort.Strings(diff)
	return diff
}

// Disabled returns the node is true or false
func (c *Config) Disabled(node string) bool {
	return c.o[node].Disabled
}

// Excludes returns an array of excluded blacklist domains/hosts
func (c *Config) Excludes(node string) []string {
	return c.o[node].Exc
}

// Files returns a list of dnsmasq conf files from all srcs
func (c *Config) Files(dir string, nodes []string) (files []string) {
	format := dir + "/%v.%v." + Fext
	for _, node := range nodes {
		for _, src := range c.Get(node).Sources {
			files = append(files, fmt.Sprintf(format, getType(src.Type), src.Name))
		}
	}
	sort.Strings(files)
	return files
}

// FormatData returns an io.Reader loaded with dnsmasq formatted data
func FormatData(fmttr string, data List) io.Reader {
	var (
		b     []byte
		lines []string
	)
	for k := range data {
		lines = append(lines, fmt.Sprintf(fmttr+"\n", k))
	}

	sort.Strings(lines)
	for _, line := range lines {
		b = append(b, line...)
	}

	return bytes.NewBuffer(b)
}

// Get returns a normalized EdgeOS data set
func (c *Config) Get(node string) (e *EdgeOS) {
	return c.o[node]
}

// GetExcludes collates the configured excludes and merges the ex/dex lists
func (c *Config) GetExcludes(nodes []string) {
	for _, node := range nodes {
		for _, k := range c.Excludes(node) {
			switch node {
			case Domains:
				c.dex[k] = 0
			case Root, Hosts:
				c.ex[k] = 0
			}
		}
	}
}

// GetIncludes processes the configured excludes and returns an io.Reader of exclusions
func (c *Config) GetIncludes(node string) io.Reader {
	inc := []byte{}
	for _, k := range c.Includes(node) {
		inc = append(inc, []byte(k+"\n")...)
	}
	return bytes.NewBuffer(inc)
}

func (n Nodes) getLeaves(node string) (leaves []Leaf) {
	leaves = []Leaf{}
	if node == "all" {
		for _, k := range n.SortKeys() {
			leaves = append(leaves, (*n[k]))
		}
		return leaves
	}
	leaves = append(leaves, (*n[node]))
	return leaves
}

// GetSeparator returns the dnsmasq conf file delimiter
func GetSeparator(node string) (sep string) {
	sep = "/"
	if node == Domains {
		sep = "/."
	}
	return sep
}

func (n Nodes) getSrcs(node string) (srcs []Srcs) {
	if node != blacklist {

		if len(n[node].Includes) > 0 {
			n[node].Data[PreConf] = &Srcs{
				List:   make(List),
				Name:   PreConf,
				Prefix: "",
				Type:   getType(node).(int),
			}
		}

		for _, src := range n.SortSKeys(node) {
			srcs = append(srcs, (*n[node].Data[src]))
		}
	}
	return srcs
}

// Includes returns an array of included blacklist domains/hosts
func (c *Config) Includes(node string) []string {
	return c.o[node].Inc
}

// IP returns the configured node IP, or the root node's IP if ""
func (c *Config) IP(node string) string {
	if c.o[node].IP == "" && node != Root {
		c.o[node].IP = c.o[Root].IP
	}
	return c.o[node].IP
}

// NewConfig returns an initialized Config map of struct EdgeOS
func (n Nodes) NewConfig() (c *Config) {
	c = &Config{
		ex:  make(List),
		dex: make(List),
		o:   make(Object),
	}

	for _, node := range Bnodes {
		ip := n[node].IP
		if ip == "" {
			ip = n[blacklist].IP
		}
		c.o[node] = &EdgeOS{
			Disabled: n[node].Disabled,
			Exc:      n[node].Excludes,
			Inc:      n[node].Includes,
			IP:       ip,
			Nodes:    n.getLeaves(node),
			Sources:  n.getSrcs(node),
		}
	}

	return c
}

// Process extracts hosts/domains from downloaded raw content
func (c *Config) Process(s *Srcs, reader io.Reader) io.Reader {
	rx := regx.Objects
	b := bufio.NewScanner(reader)

NEXT:
	for b.Scan() {
		line := strings.ToLower(b.Text())
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "#"), strings.HasPrefix(line, "//"):
			continue NEXT

		case strings.HasPrefix(line, s.Prefix):
			var ok bool

			if line, ok = stripPrefixAndSuffix(line, s.Prefix, rx); ok {
				fqdns := rx.FQDN.FindAllString(line, -1)

			FQDN:
				for _, fqdn := range fqdns {
					isDEX := c.dex.SubKeyExists(fqdn)
					isEX := c.ex.KeyExists(fqdn)
					isList := s.List.KeyExists(fqdn)

					switch {
					case isDEX:
						continue FQDN

					case isEX:
						if isList {
							s.List[fqdn]++
						}
						c.ex[fqdn]++

					case isList:
						s.List[fqdn]++

					case !isEX:
						c.ex[fqdn] = 0
						s.List[fqdn] = 0
					}
				}
			}
		default:
			continue NEXT
		}
	}

	fmttr := "address=" + GetSeparator(getType(s.Type).(string)) + "%v/" + s.IP
	return FormatData(fmttr, s.List)
}

func shuffleArray(slice []string) {
	rand.Seed(time.Now().UnixNano())
	n := len(slice)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Sources returns an array of all the node's Srcs
func (c *Config) Sources(node string) []Srcs {
	return c.o[node].Sources
}

// Source returns an array subset of the node's Srcs specified by stype
func (c *Config) Source(node, stype string) []Srcs {
	s := []Srcs{}
	for _, src := range c.o[node].Sources {
		switch stype {
		case "files":
			if src.File != "" {
				s = append(s, src)
			}
		case PreConf:
			if src.Name == PreConf {
				if src.IP == "" {
					src.IP = c.o[node].IP
				}
				s = append(s, src)
			}
		case "urls":
			s = append(s, src)
		}
	}

	return s
}

// stripPrefixAndSuffix strips the prefix and suffix
func stripPrefixAndSuffix(line, prefix string, rx *regx.OBJ) (string, bool) {
	switch {
	case prefix == "http":
		if !rx.HTTP.MatchString(line) {
			return line, false
		}
		line = rx.HTTP.FindStringSubmatch(line)[1]

	case strings.HasPrefix(line, prefix):
		line = strings.TrimPrefix(line, prefix)
	}

	line = rx.SUFX.ReplaceAllString(line, "")
	line = strings.Replace(line, `"`, "", -1)
	line = strings.TrimSpace(line)
	return line, true
}

// UpdateList converts []string to map of List
func UpdateList(data []string, list List) {
	for _, k := range data {
		list[k] = 0
	}
}

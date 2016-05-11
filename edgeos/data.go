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

// Configure has methods for returning config data supersets
type Configure interface {
	Get(path string) (e *EdgeOS)
	Files() []string
	FormatData(node string, data []string) (reader io.Reader, list List, err error)
}

// Config is a map of EdgeOS
type Config map[string]*EdgeOS

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

// Files returns a list of dnsmasq conf files from all srcs
func (c Config) Files(dir string, nodes []string) (files []string) {
	format := dir + "/%v.%v." + Fext
	for _, node := range nodes {
		if len(c.Get(node).Inc) > 0 {
			files = append(files, fmt.Sprintf(format, node, PreCon))
		}

		for _, src := range c.Get(node).Sources {
			files = append(files, fmt.Sprintf(format, GetType(src.Type), src.Name))
		}
	}
	sort.Strings(files)
	return files
}

// Get returns a normalized EdgeOS data set
func (c Config) Get(path string) (e *EdgeOS) {
	return c[path]
}

// GetExcludes collates the configured excludes and merges the ex/dex lists
func (c Config) GetExcludes(dex, ex List, nodes []string) (List, List) {
	for _, node := range nodes {
		for _, k := range c.Get(node).Exc {
			switch node {
			case Domains:
				dex[k] = 0
			case blacklist, Hosts:
				ex[k] = 0
			}
		}
	}
	return dex, ex
}

func (n Nodes) getLeaves(node string) (leaves []Leaf) {
	if node == "all" {
		for _, k := range n.SortKeys() {
			leaves = append(leaves, (*n[k]))
		}
		return leaves
	}
	leaves = append(leaves, (*n[node]))
	return leaves
}

// func (n Nodes) getLists(node string) (lists []Lists) {
// 	for node := range n {
// 		lists = append(lists, (*n[node].Data[snode].List))
// 	}
// 	return lists
// }

// getSeparator returns the dnsmasq conf file delimiter
func getSeparator(node string) (sep string) {
	sep = "/"
	if node == Domains {
		sep = "/."
	}
	return sep
}

func (n Nodes) getSrcs(node string) (srcs []Srcs) {
	// for _, node := range n.SortKeys() {
	if node != blacklist {
		for _, src := range n.SortSKeys(node) {
			srcs = append(srcs, (*n[node].Data[src]))
		}
	}
	// }
	return srcs
}

// FormatData returns a io.Reader loaded with dnsmasq formatted data
func (c Config) FormatData(fmttr string, data []string) (reader io.Reader, list List) {
	var (
		b     []byte
		lines []string
	)

	list = make(List)
	for _, k := range data {
		list[k] = 0
		lines = append(lines, fmt.Sprintf(fmttr, k))
	}

	sort.Strings(lines)
	for _, line := range lines {
		b = append(b, line...)
	}

	reader = bytes.NewBuffer(b)
	return reader, list
}

// NewConfig returns an initialized Config map of struct EdgeOS
func (n Nodes) NewConfig() (c Config) {
	c = make(Config)

	for node := range n {
		ip := n[node].IP
		if ip == "" {
			ip = n[blacklist].IP
		}
		c[node] = &EdgeOS{
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
func Process(s *Srcs, dex, ex List, reader io.Reader) *Srcs {
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
			var ok bool // We have to declare ok here, to fix var line shadow bug
			line, ok = stripPrefixAndSuffix(line, s.Prefix, rx)
			if ok {
				fqdns := rx.FQDN.FindAllString(line, -1)

			FQDN:
				for _, fqdn := range fqdns {
					isDEX := dex.SubKeyExists(fqdn)
					isEX := ex.KeyExists(fqdn)
					isList := s.List.KeyExists(fqdn)

					switch {
					case isDEX:
						continue FQDN

					case isEX:
						if isList {
							s.List[fqdn]++
						}
						ex[fqdn]++

					case isList:
						s.List[fqdn]++

					case !isEX:
						ex[fqdn] = 0
						s.List[fqdn] = 0

					}
				}
			}
		default:
			continue NEXT
		}
	}
	return s
}

func shuffleArray(slice []string) {
	rand.Seed(time.Now().UnixNano())
	n := len(slice)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
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

// WriteIncludes writes pre-configure data to disk
func (c Config) WriteIncludes(dir string, nodes []string) (dex, ex List) {
	for _, node := range nodes {
		var (
			reader   io.Reader
			includes = c.Get(node).Inc
		)

		fmttr := "address=" + getSeparator(node) + "%v/" + c.Get(node).IP + "\n"

		switch node {
		case blacklist:
			reader, dex = c.FormatData(fmttr, includes)

		case Domains:
			reader, dex = c.FormatData(fmttr, includes)

		case Hosts:
			reader, ex = c.FormatData(fmttr, includes)
		}

		if len(includes) > 0 {
			WriteFile(fmt.Sprintf("%v/%v.%v.%v", dir, node, PreCon, Fext), reader)
		}

	}
	return dex, ex
}

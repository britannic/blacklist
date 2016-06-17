package edgeos

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/britannic/blacklist/internal/regx"
)

// BooltoStr converts a boolean ("true" or "false") to a string equivalent
func BooltoStr(b bool) string {
	if b {
		return True
	}
	return False
}

// DiffArray returns the delta of two arrays
func DiffArray(a, b []string) (diff sort.StringSlice) {
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
		if !dmap.keyExists(k) {
			diff = append(diff, k)
		}
	}
	diff.Sort()
	return diff
}

// formatData returns an io.Reader loaded with dnsmasq formatted data
func formatData(fmttr string, data List) io.Reader {
	var lines sort.StringSlice
	for k := range data {
		lines = append(lines, fmt.Sprintf(fmttr+"\n", k))
	}
	lines.Sort()
	return bytes.NewBuffer([]byte(strings.Join(lines, "")))
}

// getSeparator returns the dnsmasq conf file delimiter
func getSeparator(node string) string {
	if node == domains {
		return "/."
	}
	return "/"
}

// getSubdomains returns a map of subdomains
func getSubdomains(s string) (l List) {
	l = make(List)
	keys := strings.Split(s, ".")
	for i := 0; i < len(keys)-1; i++ {
		key := strings.Join(keys[i:], ".")
		l[key] = 0
	}
	return l
}

// getType returns the converted "in" type
func getType(in interface{}) (out interface{}) {
	switch in.(type) {
	case ntype:
		out = typeInt(in.(ntype))
	case string:
		out = typeStr(in.(string))
	}
	return out
}

// ReadCfg extracts nodes from a EdgeOS/VyOS configuration structure
func ReadCfg(r ConfLoader) (*Config, error) {
	var (
		tnode  string
		b      = bufio.NewScanner(r.Load())
		branch string
		nodes  = make([]string, 2)
		rx     = regx.Objects
		s      *Object
		sCfg   = Config{Parms: &Parms{}, bNodes: make(bNodes)}
	)

LINE:
	for b.Scan() {
		line := strings.TrimSpace(b.Text())

		switch {
		case rx.MLTI.MatchString(line):
			incExc := regx.Get("mlti", line)
			switch incExc[1] {
			case "exclude":
				sCfg.bNodes[tnode].exc = append(sCfg.bNodes[tnode].exc, incExc[2])

			case "include":
				sCfg.bNodes[tnode].inc = append(sCfg.bNodes[tnode].inc, incExc[2])
			}

		case rx.NODE.MatchString(line):
			node := regx.Get("node", line)
			tnode = node[1]
			nodes = append(nodes, tnode)
			s = newObject()
			sCfg.bNodes[tnode] = s

		case rx.LEAF.MatchString(line):
			srcName := regx.Get("leaf", line)
			branch = srcName[2]
			nodes = append(nodes, srcName[1])

			if srcName[1] == src {
				s.name = branch
				s.nType = getType(tnode).(ntype)
			}

		case rx.DSBL.MatchString(line):
			sCfg.bNodes[tnode].disabled = StrToBool(regx.Get("dsbl", line)[1])

		case rx.IPBH.MatchString(line) && nodes[len(nodes)-1] != src:
			sCfg.bNodes[tnode].ip = regx.Get("ipbh", line)[1]

		case rx.NAME.MatchString(line):
			name := regx.Get("name", line)

			switch name[1] {
			case "description":
				s.desc = name[2]

			case blackhole:
				s.ip = name[2]

			case "file":
				s.file = name[2]
				s.ltype = name[1]
				sCfg.bNodes[tnode].data[branch] = s
				s = newObject() // reset s for the next loop

			case "prefix":
				s.prefix = name[2]

			case "url":
				s.ltype = name[1]
				s.url = name[2]
				sCfg.bNodes[tnode].data[branch] = s
				s = newObject() // reset s for the next loop

			}

		case rx.DESC.MatchString(line) || rx.CMNT.MatchString(line) || rx.MISC.MatchString(line):
			continue LINE

		case rx.RBRC.MatchString(line):
			nodes = nodes[:len(nodes)-1] // pop last node
			tnode = nodes[len(nodes)-1]
		}
	}

	if len(sCfg.bNodes) < 1 {
		return &sCfg, errors.New("Configuration data is empty, cannot continue")
	}
	return &sCfg, nil
}

// Source returns a map of sources
func (d data) Source(ltype string) *Objects {
	b := false
	var p *Parms
	objs := []*Object{}
	for _, k := range d.sortSKeys() {
		if !b {
			if p = d[k].Parms; p.Dir != "" {
				b = true
			}
		}
		switch {
		case ltype == d[k].ltype:
			objs = append(objs, d[k])
		case ltype == "all":
			objs = append(objs, d[k])
		}
	}
	return &Objects{Parms: p, S: objs}
}

// StrToBool converts a string ("true" or "false") to it's boolean equivalent
func StrToBool(s string) bool {
	if strings.ToLower(s) == True {
		return true
	}
	return false
}

func typeInt(i ntype) (s string) {
	switch i {
	case domain:
		s = domains
	case host:
		s = hosts
	case pre:
		s = preConf
	case root:
		s = blacklist
	case unknown:
		s = notknown
	case zone:
		s = zones
	}
	return s
}

func typeStr(s string) (i ntype) {
	switch s {
	case domains:
		i = domain
	case hosts:
		i = host
	case preConf:
		i = pre
	case blacklist:
		i = root
	case notknown:
		i = unknown
	case zones:
		i = zone
	}
	return i
}

// Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.
package edgeos

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/britannic/blacklist/internal/regx"
)

// bNodes is a map of leaf nodes
type bNodes map[string]*Object

// ConfLoader interface defines configuration load method
type ConfLoader interface {
	Load() io.Reader
}

// CFile holds an array of file names
type CFile struct {
	*Parms
	names []string
	nType ntype
}

// Config is a struct of configuration fields
type Config struct {
	*Parms
	bNodes
}

const (
	agent     = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/601.4.4 (KHTML, like Gecko) Version/9.0.3 Safari/601.4.4`
	all       = "all"
	blackhole = "dns-redirect-ip"
	blacklist = "blacklist"
	dbg       = false
	disabled  = "disabled"
	domains   = "domains"
	files     = "file"
	hosts     = "hosts"
	notknown  = "unknown"
	preNoun   = "pre-configured"
	rootNode  = blacklist
	src       = "source"
	urls      = "url"
	zones     = "zones"

	// ExcDomns labels domain exclusions
	ExcDomns = "domn-excludes"
	// ExcHosts labels host exclusions
	ExcHosts = "host-excludes"
	// ExcRoots labels global domain exclusions
	ExcRoots = "root-excludes"
	// False is a string constant
	False = "false"
	// PreDomns designates string label for preconfigured blacklisted domains
	PreDomns = preNoun + "-domain"
	// PreHosts designates string label for preconfigured blacklisted hosts
	PreHosts = preNoun + "-host"
	// True is a string constant
	True = "true"
)

func (c *Config) addExc(node string) *Objects {
	var (
		ltype string
		o     = &Objects{}
	)

	switch node {
	case domains:
		ltype = ExcDomns
	case hosts:
		ltype = ExcHosts
	case rootNode:
		ltype = ExcRoots
	}

	o.Parms = c.Parms
	o.S = append(o.S, &Object{
		desc:  ltype + " exclusions",
		exc:   c.bNodes[node].exc,
		ip:    c.bNodes.getIP(node),
		ltype: ltype,
		name:  ltype,
		nType: getType(ltype).(ntype),
		Parms: c.Parms,
	})
	return o
}

func (c *Config) addInc(node string) *Object {
	var (
		inc   = c.bNodes[node].inc
		ltype string
		n     ntype
	)
	if len(inc) > 0 {
		switch node {
		case domains:
			ltype = getType(preDomn).(string)
			n = getType(ltype).(ntype)
		case hosts:
			ltype = getType(preHost).(string)
			n = getType(ltype).(ntype)
		}
		return &Object{
			desc:  ltype + " blacklist content",
			inc:   inc,
			ip:    c.bNodes.getIP(node),
			ltype: ltype,
			name:  fmt.Sprintf("includes.[%v]", len(inc)),
			nType: n,
			Parms: c.Parms,
		}
	}
	return nil
}

// CreateObject returns an interface of the requested iFace type
func (c *Config) CreateObject(i iFace) (Contenter, error) {
	var (
		ltype = i.String()
		o     *Objects
	)

	switch ltype {
	case ExcDomns:
		o = c.addExc(domains)
	case ExcHosts:
		o = c.addExc(hosts)
	case ExcRoots:
		o = c.addExc(rootNode)
	default:
		o = c.GetAll(ltype)
	}

	switch i {
	case ExDmObj:
		return &ExcDomnObjects{Objects: o}, nil
	case ExHtObj:
		return &ExcHostObjects{Objects: o}, nil
	case ExRtObj:
		return &ExcRootObjects{Objects: o}, nil
	case FileObj:
		return &FIODataObjects{Objects: o}, nil
	case PreDObj:
		return &PreDomnObjects{Objects: o}, nil
	case PreHObj:
		return &PreHostObjects{Objects: o}, nil
	case URLsObj:
		return &URLDataObjects{Objects: o}, nil
	default:
		return nil, errors.New("Invalid interface requested")
	}
}

// excludes returns a string array of excludes
func (c *Config) excludes(nodes ...string) List {
	var exc []string
	switch nodes {
	case nil:
		for _, k := range c.Nodes() {
			if len(c.bNodes[k].exc) != 0 {
				exc = append(exc, c.bNodes[k].exc...)
			}
		}
	default:
		for _, node := range nodes {
			exc = append(exc, c.bNodes[node].exc...)
		}
	}
	return updateList(exc)
}

// Get returns an *Object for a given node
func (c *Config) Get(node string) *Objects {
	o := &Objects{Parms: c.Parms, S: []*Object{}}

	switch node {
	case all:
		for _, node := range c.Parms.Nodes {
			o.addObj(c, node)
		}
	default:
		o.addObj(c, node)
	}
	return o
}

// GetAll returns an array of Objects
func (c *Config) GetAll(ltypes ...string) *Objects {
	var (
		newDomns = true
		newHosts = true
		o        = &Objects{Parms: c.Parms}
	)

	for _, node := range c.Parms.Nodes {
		switch ltypes {
		case nil:
			o.addObj(c, node)
		default:
			for _, ltype := range ltypes {
				switch ltype {
				case PreDomns:
					if newDomns && node == domains {
						o.S = append(o.S, c.addInc(node))
						newDomns = false
					}
				case PreHosts:
					if newHosts && node == hosts {
						o.S = append(o.S, c.addInc(node))
						newHosts = false
					}
				default:
					obj := c.validate(node).S
					for i := range obj {
						if obj[i].ltype == ltype {
							o.S = append(o.S, obj[i])
						}
					}
				}
			}
		}
	}
	return o
}

// InSession returns true if VyOS/EdgeOS configuration is in session
func (c *Config) InSession() bool {
	return os.ExpandEnv("$_OFR_CONFIGURE") == "ok"
}

// load reads the config using the EdgeOS/VyOS cli-shell-api
func (c *Config) load(action string, level string) ([]byte, error) {
	cmd := exec.Command(c.Bash)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%v %v %v", c.API, apiCMD(action, c.InSession()), level))

	return cmd.Output()
}

// Nodes returns an array of configured nodes
func (c *Config) Nodes() (nodes []string) {
	for k := range c.bNodes {
		nodes = append(nodes, k)
	}
	sort.Strings(nodes)

	return nodes
}

// ReadCfg extracts nodes from a EdgeOS/VyOS configuration structure
func (c *Config) ReadCfg(r ConfLoader) error {
	var (
		tnode  string
		b      = bufio.NewScanner(r.Load())
		branch string
		nodes  []string
		rx     = regx.Objects
		s      *Object
	)

LINE:
	for b.Scan() {
		line := strings.TrimSpace(b.Text())

		switch {
		case rx.MLTI.MatchString(line):
			incExc := regx.Get("mlti", line)
			switch incExc[1] {
			case "exclude":
				c.bNodes[tnode].exc = append(c.bNodes[tnode].exc, incExc[2])

			case "include":
				c.bNodes[tnode].inc = append(c.bNodes[tnode].inc, incExc[2])
			}

		case rx.NODE.MatchString(line):
			node := regx.Get("node", line)
			tnode = node[1]
			nodes = append(nodes, tnode)
			s = newObject()
			c.bNodes[tnode] = s

		case rx.LEAF.MatchString(line):
			srcName := regx.Get("leaf", line)
			branch = srcName[2]
			nodes = append(nodes, srcName[1])

			if srcName[1] == src {
				s = newObject()
				s.name = branch
				s.nType = getType(tnode).(ntype)
			}

		case rx.DSBL.MatchString(line):
			c.bNodes[tnode].disabled = StrToBool(regx.Get("dsbl", line)[1])

		case rx.IPBH.MatchString(line) && nodes[len(nodes)-1] != src:
			c.bNodes[tnode].ip = regx.Get("ipbh", line)[1]

		case rx.NAME.MatchString(line):
			name := regx.Get("name", line)

			switch name[1] {
			case "description":
				s.desc = name[2]

			case blackhole:
				s.ip = name[2]

			case files:
				s.file = name[2]
				s.ltype = name[1]
				c.bNodes[tnode].Objects.S = append(c.bNodes[tnode].Objects.S, s)
				s = newObject() // reset s for the next loop

			case "prefix":
				s.prefix = name[2]

			case urls:
				s.ltype = name[1]
				s.url = name[2]
				c.bNodes[tnode].Objects.S = append(c.bNodes[tnode].Objects.S, s)
			}

		case rx.DESC.MatchString(line) || rx.CMNT.MatchString(line) || rx.MISC.MatchString(line):
			continue LINE

		case rx.RBRC.MatchString(line):
			if len(nodes) > 1 {
				nodes = nodes[:len(nodes)-1] // pop last node
				tnode = nodes[len(nodes)-1]
			}
		}
	}

	if len(c.bNodes) < 1 {
		return errors.New("Configuration data is empty, cannot continue")
	}

	return nil
}

// readDir returns a listing of dnsmasq formatted blacklist configuration files
func (c *CFile) readDir(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

// ReloadDNS reloads the dnsmasq configuration
func (c *Config) ReloadDNS() ([]byte, error) {
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = strings.NewReader(c.DNSsvc)

	return cmd.CombinedOutput()
}

// Remove deletes a CFile array of file names
func (c *CFile) Remove() error {

	pattern := fmt.Sprintf(c.FnFmt, c.Dir, c.Wildcard.Node, c.Wildcard.Name, c.Parms.Ext)
	dlist, err := c.readDir(pattern)
	if err != nil {
		return err
	}

	return purgeFiles(DiffArray(c.names, dlist))
}

// String returns pretty print for the Blacklist struct
func (c *Config) String() (result string) {
	indent := 1
	cmma := comma
	cnt := len(c.sortKeys())
	result += fmt.Sprintf("{\n%v%q: [{\n", tabs(indent), "nodes")

	for i, pkey := range c.sortKeys() {
		if i == cnt-1 {
			cmma = null
		}

		indent++
		result += fmt.Sprintf("%v%q: {\n", tabs(indent), pkey)
		indent++

		result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), disabled,
			BooltoStr(c.bNodes[pkey].disabled))
		result = is(indent, result, "ip", c.bNodes[pkey].ip)
		result += getJSONArray(&cfgJSON{array: c.bNodes[pkey].exc, pk: pkey, leaf: "excludes", indent: indent})

		if pkey != rootNode {
			result += getJSONArray(&cfgJSON{array: c.bNodes[pkey].inc, pk: pkey, leaf: "includes", indent: indent})
			result += getJSONsrcArray(&cfgJSON{Config: c, pk: pkey, indent: indent})
		}

		indent--
		result += fmt.Sprintf("%v}%v\n", tabs(indent), cmma)
		indent--
	}

	result += tabs(indent) + "}]\n}"

	return result
}

// String implements string method
func (c *CFile) String() string {
	return strings.Join(c.names, "\n")
}

// Strings returns a sorted array of strings.
func (c *CFile) Strings() []string {
	sort.Strings(c.names)

	return c.names
}

// LTypes returns an array of configured nodes
func (c *Config) LTypes() []string {
	return c.Parms.Ltypes
}

func (b bNodes) getIP(node string) (ip string) {
	switch b[node].ip {
	case "":
		ip = b[rootNode].ip
	default:
		ip = b[node].ip
	}
	return ip
}

func (b bNodes) validate(node string) *Objects {
	for _, obj := range b[node].Objects.S {
		if obj.ip == "" {
			obj.ip = b.getIP(node)
		}
	}
	return &b[node].Objects
}

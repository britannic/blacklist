//Package dnsmasq parses dnsmasq.conf address and server name IP mapping files
package dnsmasq

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"strings"
)

// Host is a container for IP addresses
type Host struct {
	IP     string `json:"IP,omitempty"`
	Server bool   `json:"Server,omitempty"`
}

// Conf is map of Hosts
type Conf map[string]Host

type confLoader interface {
	read() io.Reader
}

const (
	address = "address="
	server  = "server="
)

func ipv4(ip string) bool {
	return net.IP([]byte(ip)).To4() == nil
}

func matchIP(ip string, ips []string) bool {
	var b bool
	for _, dns := range ips {
		b = ip == dns && ipv4(dns)
	}
	return b
}

func fetchHost(k, ip string) bool {
	ips, err := net.LookupHost(k)
	if err != nil {
		return false
	}
	return matchIP(ip, ips)
}

// Redirect return true if the resolved IP address matches the correct IP (redirected or normal)
func (c Conf) Redirect(k, ip string) bool {
	if _, ok := c[k]; ok {
		if c[k].Server && c[k].IP == "#" {
			return !fetchHost(c[k].IP, ip)
		}
		return fetchHost(c[k].IP, ip)
	}
	return false
}

// Parse extracts host to IP mappings from a dnsmasq configuration file
func (c Conf) Parse(r confLoader) error {
	b := bufio.NewScanner(r.read())

	for b.Scan() {
		line := bytes.TrimSpace(b.Bytes())
		b := bytes.Split(line, []byte("/"))

		switch {
		case len(b) < 3:
			return errors.New("no dnsmasq configuration mapping entries found")
		case bytes.HasPrefix(b[0], []byte(address)):
			k := string(b[1])
			c[k] = Host{
				IP:     string(b[2]),
				Server: false,
			}
		case bytes.HasPrefix(b[0], []byte(server)):
			k := string(b[1])
			c[k] = Host{
				IP:     string(b[2]),
				Server: true,
			}
		}
	}
	return nil
}

// Mapping holds a dnsmasq configuration file contents
type Mapping struct {
	Contents string
}

func (m *Mapping) read() io.Reader {
	return strings.NewReader(m.Contents)
}

func (c Conf) String() string {
	j, _ := json.Marshal(c)
	return string(j)
}

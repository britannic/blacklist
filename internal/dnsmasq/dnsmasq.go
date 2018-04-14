//Package dnsmasq parses dnsmasq.conf address and server name IP mapping files
package dnsmasq

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"os"
)

const (
	address = "address="
	server  = "server="
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

// Mapping holds a dnsmasq configuration file contents
type Mapping struct {
	Contents []byte
}

// ConfigFile reads a file and returns an io.Reader
func ConfigFile(f string) (io.Reader, error) {
	return os.Open(f)
}

func fetchHost(k, ip string) bool {
	ips, err := net.LookupHost(k)
	if err != nil {
		return false
	}
	return matchIP(ip, ips)
}

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

// Parse extracts host to IP mappings from a dnsmasq configuration file
func (c Conf) Parse(r confLoader) error {
	b := bufio.NewScanner(r.read())

	for b.Scan() {
		l := bytes.TrimSpace(b.Bytes())
		d := bytes.Split(l, []byte("/"))

		switch {
		case len(d) < 3:
			return errors.New("no dnsmasq configuration mapping entries found")
		case bytes.HasPrefix(d[0], []byte(address)):
			c[string(d[1])] = Host{IP: string(d[2]), Server: false}
		case bytes.HasPrefix(d[0], []byte(server)):
			c[string(d[1])] = Host{IP: string(d[2]), Server: true}
		}
	}
	return nil
}

func (m *Mapping) read() io.Reader {
	return bytes.NewReader(m.Contents)
}

// Redirect returns true if the resolved IP address matches the correct IP (redirected or normal)
func (c Conf) Redirect(k, ip string) bool {
	if _, ok := c[k]; ok {
		if c[k].Server && c[k].IP == "#" {
			return !fetchHost(c[k].IP, ip)
		}
		return fetchHost(c[k].IP, ip)
	}
	return false
}

func (c Conf) String() string {
	j, _ := json.Marshal(c)
	return string(j)
}

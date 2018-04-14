package dnsmasq

import (
	"encoding/json"
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFetchHost(t *testing.T) {
	tests := []struct {
		conf Conf
		exp  bool
		ip   string
		key  string
		name string
	}{
		{
			ip:   "0.0.0.0",
			key:  "badguy_s.com",
			conf: Conf{"badguys.com": Host{IP: "0.0.0.0", Server: false}},
			exp:  false,
			name: "badguys.com",
		},
		{
			ip:   "127.0.0.1",
			key:  "localhoster",
			conf: Conf{"localhost": Host{IP: "127.0.0.1", Server: false}},
			exp:  false,
			name: "localhoster",
		},
		{
			ip:   "127.0.0.1",
			key:  "localhost",
			conf: Conf{"localhost": Host{IP: "#", Server: true}},
			exp:  true,
			name: "localServer",
		},
		{
			ip:   "127.0.0.1",
			key:  "localhost",
			conf: Conf{"localhost": Host{IP: "127.0.0.1", Server: false}},
			exp:  true,
			name: "localhost",
		},
		{
			ip:   "127.0.0.1",
			exp:  false,
			name: "no name",
		},
		{
			ip:   "::1",
			key:  "localhost",
			conf: Conf{"localhost": Host{IP: "127.0.0.1", Server: false}},
			exp:  false,
			name: "localhost IPv6",
		},
	}

	Convey("Testing String()", t, func() {
		for _, tt := range tests {
			Convey("current test "+tt.name, func() {
				So(fetchHost(tt.key, tt.ip), ShouldEqual, tt.exp)
				So(tt.conf.Redirect(tt.key, tt.ip), ShouldEqual, tt.exp)
			})
		}
	})
}

func TestParse(t *testing.T) {
	tests := []struct {
		act string
		Host
		err    error
		exp    string
		name   string
		reader Mapping
	}{
		{
			Host: Host{
				IP:     "127.0.0.1",
				Server: false,
			},
			act:    `{"badguys.com":{"IP":"0.0.0.0"}}`,
			err:    nil,
			exp:    "127.0.0.1",
			name:   "badguys.com",
			reader: Mapping{Contents: `address=/badguys.com/0.0.0.0`},
		},
		{
			Host: Host{
				IP:     "127.0.0.1",
				Server: true,
			},
			act:    `{"xrated.com":{"IP":"0.0.0.0","Server":true}}`,
			err:    nil,
			exp:    "127.0.0.1",
			name:   "xrated.com",
			reader: Mapping{Contents: `server=/xrated.com/0.0.0.0`},
		},
		{
			act:  `{}`,
			err:  errors.New("no dnsmasq configuration mapping entries found"),
			exp:  "127.0.0.1",
			name: "No dnsmasq entry",
			reader: Mapping{Contents: `# All files in this directory will be read by dnsmasq as 
# configuration files, except if their names end in 
# ".dpkg-dist",".dpkg-old" or ".dpkg-new"
#
# This can be changed by editing /etc/default/dnsmasq`},
		},
	}
	Convey("Conf map should show each map entry", t, func() {
		c := make(Conf)
		for _, tt := range tests {
			Convey("current test: "+tt.name, func() {
				if err := c.Parse(&tt.reader); err != nil {
					So(err.Error(), ShouldEqual, tt.err.Error())
				}
				j, err := json.Marshal(c)
				So(err, ShouldBeNil)
				So(string(j), ShouldEqual, tt.act)
			})
		}
	})
}

func TestString(t *testing.T) {
	tests := []struct {
		conf Conf
		exp  string
	}{
		{
			conf: Conf{"badguys.com": Host{IP: "0.0.0.0", Server: false}},
			exp:  `{"badguys.com":{"IP":"0.0.0.0"}}`,
		},
		{
			exp: `null`,
		},
	}

	Convey("Testing String()", t, func() {
		for _, tt := range tests {
			So(tt.conf.String(), ShouldEqual, tt.exp)
		}
	})
}

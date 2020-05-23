// Package edgeos provides methods and structures to retrieve, parse and render EdgeOS configuration data and files.
package edgeos

import (
	"net"
)

// Chk_Web() returns true if DNS is working
func Chk_Web(site string) bool {
	if _, err := net.LookupIP(site); err != nil {
		return false
	}
	return true
}

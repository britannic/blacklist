package edgeos

import (
	"net"
)

// Chk_Web() returns true if DNS is working
func Chk_Web(site, port string) bool {
	conn, err := net.Dial("tcp", net.JoinHostPort(site, port))
	conn.Close()
	return err == nil
}

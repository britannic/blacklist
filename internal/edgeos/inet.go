package edgeos

import (
	"fmt"
	"net"
)

// Chk_Web() returns true if DNS is working
func Chk_Web(site, port string) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", site, port))
	conn.Close()
	return err == nil
}

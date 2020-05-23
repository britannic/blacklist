package edgeos

import (
	"net"
	"time"
)

// ChkWeb() returns true if DNS is working
func ChkWeb(site, port string) bool {
	timeOut := 3 * time.Second
	conn, err := net.DialTimeout("tcp4", net.JoinHostPort(site, port), timeOut)
	conn.Close()
	return err == nil
}

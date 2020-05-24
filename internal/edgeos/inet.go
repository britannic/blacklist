package edgeos

import (
	"fmt"
	"net"
	"time"
)

// ChkWeb() returns true if DNS is working
func ChkWeb(site string, port int) bool {
	target := fmt.Sprintf("%s:%d", site, port)
	timeOut := 3 * time.Second
	conn, err := net.DialTimeout("tcp4", target, timeOut)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

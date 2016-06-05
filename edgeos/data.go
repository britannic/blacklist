package edgeos

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// getSeparator returns the dnsmasq conf file delimiter
func getSeparator(node string) string {
	if node == domains {
		return "/."
	}
	return "/"
}

// formatData returns an io.Reader loaded with dnsmasq formatted data
func formatData(fmttr string, data List) io.Reader {
	var lines []string
	for k := range data {
		lines = append(lines, fmt.Sprintf(fmttr+"\n", k))
	}
	sort.Strings(lines)
	return bytes.NewBuffer([]byte(strings.Join(lines, "")))
}

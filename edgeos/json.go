package edgeos

import "encoding/json"

// JSON returns raw print for the Blacklist struct
func (n Nodes) JSON() string {
	j, _ := json.Marshal(n)
	return string(j)
}

// String returns pretty print for the Blacklist struct
func (n Nodes) String() string {
	j, _ := json.MarshalIndent(n, "", "  ")
	return string(j)
}

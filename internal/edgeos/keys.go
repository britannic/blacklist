package edgeos

import "sort"

// sortFlags returns the flags as a slice in lexicographical sorted order.
func (c *Config) sortKeys() (pkeys sort.StringSlice) {
	pkeys = make(sort.StringSlice, len(c.bNodes))
	i := 0
	for k := range c.bNodes {
		pkeys[i] = k
		i++
	}
	pkeys.Sort()

	return pkeys
}

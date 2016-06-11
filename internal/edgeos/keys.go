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

// sortSKeys returns an array of sorted Keys
func (c *Config) sortSKeys(node string) (skeys sort.StringSlice) {
	skeys = make(sort.StringSlice, len(c.bNodes[node].data))
	i := 0
	for k := range c.bNodes[node].data {
		skeys[i] = k
		i++
	}
	skeys.Sort()

	return skeys
}

// sortSKeys returns an array of custom sorted strings
func (d data) sortSKeys() (skeys sort.StringSlice) {
	for k := range d {
		if d[k].ltype != preConf {
			skeys = append(skeys, k)
			// skeys[i] = k
			// i++
		}
	}
	skeys.Sort()

	if _, ok := d[preConf]; ok {
		skeys = append(sort.StringSlice{preConf}, skeys...)
	}
	return skeys
}

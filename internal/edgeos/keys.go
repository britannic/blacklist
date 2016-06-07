package edgeos

import "sort"

// Keys is used for sorting operations on map Keys
type Keys []string

// len returns length of Keys
func (k Keys) Len() int { return len(k) }

// less returns the smallest element
func (k Keys) Less(i, j int) bool { return k[i] < k[j] }

// Swap swaps elements of a key array
func (k Keys) Swap(i, j int) { k[i], k[j] = k[j], k[i] }

// sortKeys returns an array of sorted Keys
func (c *Config) sortKeys() (pkeys Keys) {
	for pkey := range c.bNodes {
		pkeys = append(pkeys, pkey)
	}
	sort.Sort(Keys(pkeys))
	return pkeys
}

// sortSKeys returns an array of sorted Keys
func (c *Config) sortSKeys(node string) (skeys Keys) {
	for skey := range c.bNodes[node].data {
		skeys = append(skeys, skey)
	}
	sort.Sort(Keys(skeys))
	return skeys
}

// sortSKeys returns an array of custom sorted strings
func (d data) sortSKeys() (skeys Keys) {
	for skey := range d {
		if d[skey].ltype != preConf {
			skeys = append(skeys, skey)
		}
	}
	sort.Sort(Keys(skeys))
	if _, ok := d[preConf]; ok {
		skeys = append(Keys{preConf}, skeys...)
	}
	return skeys
}

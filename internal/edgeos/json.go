package edgeos

import "fmt"

const (
	comma = ","
	enter = "\n"
	null  = ""
	tab   = "  "
)

func tabs(t int) (r string) {
	if t <= 0 {
		return r
	}
	for i := 0; i < t; i++ {
		r += tab
	}
	return r
}

type cfgJSON struct {
	array []string
	*Config
	indent       int
	leaf, pk, sk string
}

func getJSONdisabled(c *cfgJSON) (d string) {
	d = False
	switch c.sk {
	case null:
		d = BooltoStr(c.bNodes[c.pk].disabled)
	default:
		d = BooltoStr(c.bNodes[c.pk].data[c.sk].disabled)
	}
	return d
}

func getJSONsrcIP(c *Config, pkey string) (result string) {
	b := c.bNodes
	if len(b[pkey].ip) > 0 {
		result += fmt.Sprintf("%q: %q,\n", "ip", b[pkey].ip)
	}
	return result
}

func getJSONArray(c *cfgJSON) (result string) {
	indent := c.indent
	cmma := comma
	ret := enter
	result += fmt.Sprintf("%v%q: [", tabs(indent), c.leaf)
	cnt := len(c.array)

	switch {
	case c.pk != rootNode && cnt == 0:
		result += "],\n"
		return result

	case cnt == 1:
		ret = null
		indent = 0

	case cnt > 1:
		result += enter
		indent++
	}

	if cnt > 0 {
		for i, s := range c.array {
			if i == cnt-1 {
				cmma = null
			}
			result += fmt.Sprintf("%v%q%v%v", tabs(indent), s, cmma, ret)
		}

		cmma = comma

		if c.pk == rootNode {
			cmma = null
		}

		result += fmt.Sprintf("%v]%v\n", tabs(indent), cmma)
	}

	return result
}

func getJSONsrcArray(c *cfgJSON) (result string) {
	var (
		i      int
		indent = c.indent
		skeys  = c.Config.sortSKeys(c.pk)
		cnt    = len(skeys)
	)

	if cnt == 0 {
		result += fmt.Sprintf("%v%q: [{}]\n", tabs(c.indent), "sources")
		return result
	}

	result += fmt.Sprintf("%v%q: [{%v", tabs(c.indent), "sources", enter)

	for i, c.sk = range skeys {
		if _, ok := c.Config.bNodes[c.pk].data[c.sk]; ok {
			cmma := comma
			indent = c.indent + 1

			if i == cnt-1 {
				cmma = null
			}

			d := getJSONdisabled(&cfgJSON{Config: c.Config, pk: c.pk, sk: c.sk})
			s := c.bNodes[c.pk].data[c.sk]
			result += fmt.Sprintf("%v%q: {\n", tabs(indent), c.sk)
			indent++
			result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), disabled, d)
			result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), "description", s.desc)
			result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), "prefix", s.prefix)
			result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), "file", s.file)
			result += fmt.Sprintf("%v%q: %q\n", tabs(indent), "url", s.url)
			indent--
			result += fmt.Sprintf("%v}%v%v", tabs(indent), cmma, enter)
		}
	}

	indent -= 2
	result += fmt.Sprintf("%v}]%v", tabs(indent), enter)
	return result
}

// String returns pretty print for the Blacklist struct
func (c *Config) String() (result string) {
	indent := 1
	cmma := comma
	cnt := len(c.sortKeys())
	result += fmt.Sprintf("{\n%v%q: [{\n", tabs(indent), "nodes")

	for i, pkey := range c.sortKeys() {

		if i == cnt-1 {
			cmma = null
		}

		indent++
		result += fmt.Sprintf("%v%q: {\n", tabs(indent), pkey)

		indent++
		result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), disabled, getJSONdisabled(&cfgJSON{Config: c, pk: pkey}))

		result += tabs(indent) + getJSONsrcIP(c, pkey)

		result += getJSONArray(&cfgJSON{array: c.bNodes[pkey].exc, pk: pkey, leaf: "excludes", indent: indent})

		if pkey != rootNode {
			result += getJSONArray(&cfgJSON{array: c.bNodes[pkey].inc, pk: pkey, leaf: "includes", indent: indent})
		}

		if pkey != rootNode {
			result += getJSONsrcArray(&cfgJSON{Config: c, pk: pkey, indent: indent})
		}

		indent--
		result += fmt.Sprintf("%v}%v\n", tabs(indent), cmma)
		indent--
	}

	result += tabs(indent) + "}]\n}"
	return result
}

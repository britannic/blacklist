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

// func getJSONdisabled(c *cfgJSON) (d string) {
// 	d = False
// 	switch c.sk {
// 	case null:
// 		d = BooltoStr(c.bNodes[c.pk].disabled)
// 	default:
// 		d = BooltoStr(c.bNodes[c.pk].data[c.sk].disabled)
// 	}
// 	return d
// }

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
		// skeys  = c.Config.sortSKeys(c.pk)
		cnt    = len(c.bNodes[c.pk].Objects.S)
		i      int
		indent = c.indent
		s      *Object
	)

	if cnt == 0 {
		result += fmt.Sprintf("%v%q: [{}]\n", tabs(c.indent), "sources")
		return result
	}

	result += fmt.Sprintf("%v%q: [{%v", tabs(c.indent), "sources", enter)

	for i, s = range c.bNodes[c.pk].Objects.S {
		// if _, ok := c.bNodes[c.pk].data[c.sk]; ok {
		cmma := comma
		indent = c.indent + 1

		if i == cnt-1 {
			cmma = null
		}

		// d := getJSONdisabled(&cfgJSON{Config: c.Config, pk: c.pk, sk: c.sk})
		// s := c.bNodes[c.pk].data[c.sk]
		result += fmt.Sprintf("%v%q: {\n", tabs(indent), s.name)
		indent++
		result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), disabled, BooltoStr(s.disabled))
		result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), "description", s.desc)
		result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), "prefix", s.prefix)
		result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), "file", s.file)
		result += fmt.Sprintf("%v%q: %q\n", tabs(indent), urls, s.url)
		indent--
		result += fmt.Sprintf("%v}%v%v", tabs(indent), cmma, enter)
		// }
	}

	indent -= 2
	result += fmt.Sprintf("%v}]%v", tabs(indent), enter)
	return result
}

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

func is(indent int, result, title, s string) string {
	if s != "" {
		result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), title, s)
		return result
	}
	return result
}

func getJSONsrcArray(c *cfgJSON) (result string) {
	var (
		cnt    = len(c.bNodes[c.pk].objects.obs)
		i      int
		indent = c.indent
		s      *object
	)

	if cnt == 0 {
		result += fmt.Sprintf("%v%q: [{}]\n", tabs(c.indent), "sources")
		return result
	}

	result += fmt.Sprintf("%v%q: [{%v", tabs(c.indent), "sources", enter)

	for i, s = range c.bNodes[c.pk].objects.obs {
		cmma := comma
		indent = c.indent + 1

		if i == cnt-1 {
			cmma = null
		}

		result += fmt.Sprintf("%v%q: {\n", tabs(indent), s.name)
		indent++
		result += fmt.Sprintf("%v%q: %q,\n", tabs(indent), disabled, BooltoStr(s.disabled))
		result = is(indent, result, "description", s.desc)
		result = is(indent, result, "ip", s.ip)
		result = is(indent, result, "prefix", s.prefix)
		result = is(indent, result, files, s.file)
		result = is(indent, result, urls, s.url)
		indent--
		result += fmt.Sprintf("%v}%v%v", tabs(indent), cmma, enter)
	}

	indent -= 2
	result += fmt.Sprintf("%v}]%v", tabs(indent), enter)
	return result
}

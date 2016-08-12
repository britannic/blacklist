package edgeos

import "fmt"

const (
	comma = ","
	enter = "\n"
	null  = ""
	tab   = "  "
)

type cfgJSON struct {
	*Config
	array        []string
	indent       int
	leaf, pk, sk string
}

func tabs(t int) (s string) {
	if t <= 0 {
		return s
	}
	for i := 0; i < t; i++ {
		s += tab
	}
	return s
}

func getJSONArray(c *cfgJSON) (js string) {
	cma := comma
	cnt := len(c.array)
	ind := c.indent
	js += fmt.Sprintf("%v%q: [", tabs(ind), c.leaf)
	ret := enter

	switch {
	case c.pk != rootNode && cnt == 0:
		js += "],\n"
		return js

	case cnt == 1:
		ret = null
		ind = 0

	case cnt > 1:
		js += enter
		ind++
	}

	if cnt > 0 {
		for i, s := range c.array {
			if i == cnt-1 {
				cma = null
			}
			js += fmt.Sprintf("%v%q%v%v", tabs(ind), s, cma, ret)
		}

		cma = comma
		if c.pk == rootNode {
			cma = null
		}

		js += fmt.Sprintf("%v]%v\n", tabs(ind), cma)
	}

	return js
}

func is(ind int, js, title, s string) string {
	if s != "" {
		js += fmt.Sprintf("%v%q: %q,\n", tabs(ind), title, s)
		return js
	}
	return js
}

func getJSONsrcArray(c *cfgJSON) (js string) {
	var (
		cnt = len(c.tree[c.pk].Objects.x)
		i   int
		ind = c.indent
		o   *object
	)

	if cnt == 0 {
		js += fmt.Sprintf("%v%q: [{}]\n", tabs(c.indent), "sources")
		return js
	}

	js += fmt.Sprintf("%v%q: [{%v", tabs(c.indent), "sources", enter)

	for i, o = range c.tree[c.pk].Objects.x {
		cma := comma
		ind = c.indent + 1

		if i == cnt-1 {
			cma = null
		}

		js += fmt.Sprintf("%v%q: {\n", tabs(ind), o.name)
		ind++
		js += fmt.Sprintf("%v%q: %q,\n", tabs(ind), disabled, booltoStr(o.disabled))
		js = is(ind, js, "description", o.desc)
		js = is(ind, js, "ip", o.ip)
		js = is(ind, js, "prefix", o.prefix)
		js = is(ind, js, files, o.file)
		js = is(ind, js, urls, o.url)
		ind--
		js += fmt.Sprintf("%v}%v%v", tabs(ind), cma, enter)
	}

	ind -= 2
	js += fmt.Sprintf("%v}]%v", tabs(ind), enter)
	return js
}

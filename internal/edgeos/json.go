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
	array    []string
	indent   int
	leaf, pk string
}

func tabs(t int) (s string) {
	if t <= 0 {
		return s
	}
	for range Iter(t) {
		s += tab
	}
	return s
}

func getJSONArray(c *cfgJSON) (js string) {
	ø := comma
	cnt := len(c.array)
	ȹ := c.indent
	js += fmt.Sprintf("%v%q: [", tabs(ȹ), c.leaf)
	eof := enter

	switch {
	case cnt == 0:
		js += "]," + enter
		return js

	case cnt == 1:
		eof = null
		ȹ = 0

	case cnt > 1:
		js += enter
		ȹ++
	}

	if cnt > 0 {
		for i, s := range c.array {
			if i == cnt-1 {
				ø = null
			}
			js += fmt.Sprintf("%v%q%v%v", tabs(ȹ), s, ø, eof)
		}

		ø = comma
		js += fmt.Sprintf("%v]%v\n", tabs(ȹ), ø)
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

func getJSONsrcArray(c *cfgJSON) string {
	var (
		cnt = len(c.tree[c.pk].src)
		i   int
		js  string
		o   *source
		ø   = comma
		ȹ   = c.indent
	)

	if cnt == 0 {
		js += fmt.Sprintf("%v%q: [{}]\n", tabs(c.indent), "sources")
		return js
	}

	js += fmt.Sprintf("%v%q: [{%v", tabs(c.indent), "sources", enter)

	for i, o = range c.tree[c.pk].src {
		ȹ = c.indent + 1

		if i == cnt-1 {
			ø = null
		}

		js += fmt.Sprintf("%v%q: {\n", tabs(ȹ), o.name)
		ȹ++
		js += fmt.Sprintf("%v%q: %q,\n", tabs(ȹ), disabled, booltoStr(o.disabled))
		js = is(ȹ, js, "description", o.desc)
		js = is(ȹ, js, "ip", o.ip)
		js = is(ȹ, js, "prefix", o.prefix)
		js = is(ȹ, js, files, o.file)
		js = is(ȹ, js, urls, o.url)
		ȹ--
		js += fmt.Sprintf("%v}%v%v", tabs(ȹ), ø, enter)
	}

	ȹ -= 2
	js += fmt.Sprintf("%v}]%v", tabs(ȹ), enter)
	return js
}

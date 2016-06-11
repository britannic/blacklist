package edgeos

import (
	"fmt"
	"strings"

	"github.com/britannic/oldBlist2/regx"
)

// Excludes returns a List map of blacklist exclusions
func (o *Object) Excludes() List {
	return UpdateList(o.exc)
}

// String pretty prints Object
func (o *Object) String() (r string) {
	r += fmt.Sprintf("\nDesc:\t %q\n", o.desc)
	r += fmt.Sprintf("Disabled: %v\n", o.disabled)
	r += fmt.Sprintf("File:\t %q\n", o.file)
	r += fmt.Sprintf("IP:\t %q\n", o.ip)
	r += fmt.Sprintf("Ltype:\t %q\n", o.ltype)
	r += fmt.Sprintf("Name:\t %q\n", o.name)
	r += fmt.Sprintf("nType:\t %q\n", o.nType)
	r += fmt.Sprintf("Prefix:\t %q\n", o.prefix)
	r += fmt.Sprintf("Type:\t %q\n", getType(o.nType))
	r += fmt.Sprintf("URL:\t %q\n", o.url)
	return r
}

// stripPrefixAndSuffix strips the prefix and suffix
func stripPrefixAndSuffix(line, prefix string, rx *regx.OBJ) (string, bool) {
	switch {
	case prefix == "http":
		if !rx.HTTP.MatchString(line) {
			return line, false
		}
		line = rx.HTTP.FindStringSubmatch(line)[1]

	case strings.HasPrefix(line, prefix):
		line = strings.TrimPrefix(line, prefix)
	}

	line = rx.SUFX.ReplaceAllString(line, "")
	line = strings.Replace(line, `"`, "", -1)
	line = strings.TrimSpace(line)
	return line, true
}

package edgeos

import "fmt"

// Excludes returns a List map of blacklist exclusions
func (o *Object) Excludes() List {
	return UpdateList(o.exc)
}

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

// Swap swaps elements of a slice of Objects
// func (o Objects) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

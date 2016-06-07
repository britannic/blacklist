package edgeos

// GetContent returns a Content struct
func (objs *Objects) GetContent() *Contents {
	var c Contents
	for _, o := range objs.S {
		switch o.ltype {
		case preConf:
			if o.inc != nil {
				c = append(c, &Content{
					err:    nil,
					Object: o,
					r:      o.Includes(),
				})
			}

		case "files":
			if o.file != "" {
				b, err := getFile(o.file)
				c = append(c, &Content{
					err:    err,
					Object: o,
					r:      b,
				})
			}

		case "urls":
			if o.url != "" {
				reader, err := GetHTTP(o.Parms.method, o.url)
				c = append(c, &Content{
					err:    err,
					Object: o,
					r:      reader,
				})
			}
		}
	}
	return &c
}

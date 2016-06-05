package edgeos

// GetContent returns a Content struct
func (objs *Objects) GetContent() *Contents {
	var c Contents
	for _, o := range *objs {
		switch o.ltype {
		case preConf:
			if o.inc != nil {
				c = append(c, &Content{
					err:    nil,
					object: o,
					r:      o.Includes(),
				})
			}

		case "files":
			if o.file != "" {
				c = append(c, &Content{
					err:    nil,
					object: o,
					r:      nil,
				})
			}

		case "urls":
			if o.url != "" {
				reader, err := GetHTTP(o.parms.method, o.url)
				c = append(c, &Content{
					err:    err,
					object: o,
					r:      reader,
				})
			}
		}
	}
	return &c
}

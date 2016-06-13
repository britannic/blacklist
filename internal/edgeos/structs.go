package edgeos

import "io"

type blist struct {
	file   string
	reader io.Reader
}

// bNodes is a map of leaf nodes
type bNodes map[string]*Object

// CFile holds an array of file names
type CFile struct {
	*Parms
	names []string
	nType ntype
}

// Config is a struct of configuration fields
type Config struct {
	bNodes
	*Parms
}

// data is a map[string] of *Object
type data map[string]*Object

// Object struct for normalizing EdgeOS data.
type Object struct {
	*Parms
	data
	desc     string
	disabled bool
	exc      []string
	file     string
	inc      []string
	ip       string
	ltype    string
	name     string
	nType    ntype
	prefix   string
	url      string
}

// Objects is a struct of []*Objects
type Objects struct {
	S []*Object
	*Parms
}

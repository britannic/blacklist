package edgeos

import "io"

type blist struct {
	r io.Reader
}

// bNodes is a map of leaf nodes
type bNodes map[string]*object

// CFile holds an array of file names
type CFile struct {
	*parms
	names []string
}

// Config is a struct of configuration fields
type Config struct {
	bNodes
	*parms
}

// Content is a struct of blacklist content
type Content struct {
	err error
	*object
	r io.Reader
}

// Contents is an array of *content
type Contents []*Content

type data map[string]*object

// object struct for normalizing EdgeOS data.
type object struct {
	*parms
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

// Objects is a map[string] of *objects
type Objects []*object

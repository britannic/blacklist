package edgeos

import (
	"io"
	"os"
)

// ConfLoader interface defines load method
type ConfLoader interface {
	Load() io.Reader
}

// OSinformer implements os.FileInfo methods
type OSinformer interface {
	ReadDir(string) ([]os.FileInfo, error)
	Remove() error
}

// CFGstatic is for configurations loaded via the EdgeOS CFGstatic
type CFGstatic struct {
	Cfg string
}

// CFGcli is for configurations loaded via the EdgeOS CFGcli
type CFGcli struct {
	Cfg string
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

// Objects is a struct of []*Object
type Objects struct {
	S []*Object
	*Parms
}

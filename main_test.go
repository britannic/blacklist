package main

import (
	"testing"

	. "github.com/britannic/testutils"
)

func TestBuild(t *testing.T) {
	want := map[string]string{
		"build":   build,
		"githash": githash,
		"version": version,
	}

	for k := range want {
		Equals(t, "UNKNOWN", want[k])
	}
}

func TestGetOpts(t *testing.T) {
	// global.Args = []string{"-debug", "-i", "8", "-test", "-v", "-version", "-help"}
	o := getOpts()
	want := `Debug
File
Poll
Test
Verb
Version
`
	Equals(t, want, o.String())

	tests := []struct {
		name string
		test interface{}
		exp  interface{}
	}{
		{
			name: "o.Debug",
			test: o.Debug,
			exp:  true,
		},
		{
			name: "o.File",
			test: o.File,
			exp:  "",
		},
		{
			name: "o.Poll",
			test: o.Poll,
			exp:  8,
		},
		{
			name: "o.Test",
			test: o.Test,
			exp:  true,
		},
		{
			name: "o.Verb",
			test: o.Verb,
			exp:  true,
		},
		{
			name: "o.Version",
			test: o.Version,
			exp:  true,
		},
	}

	for _, run := range tests {
		switch run.test.(type) {
		case bool:
			Equals(t, run.exp.(bool), run.test.(bool))

		case string:
			Equals(t, run.exp.(string), run.test.(string))

		case int:
			Equals(t, run.exp.(int), run.test.(int))
		}
	}
}

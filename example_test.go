package main

import "github.com/britannic/blacklist/global"

func Examplemain() {
	global.Args = []string{"-version"}
	main()
	// Output:
	// Version:		UNKNOWN
	// Build date:		UNKNOWN
	// Git short hash:	UNKNOWN
}

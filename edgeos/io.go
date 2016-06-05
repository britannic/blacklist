package edgeos

import (
	"fmt"
	"os"
	"strings"
)

// purgeFiles removes any orphaned blacklist files that don't have sources
func purgeFiles(files []string) error {
	var errArray []string

NEXT:
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			// errArray = append(errArray, fmt.Sprintf("%q: %v", file, err))
			continue NEXT
		}
		if !deleteFile(file) {
			errArray = append(errArray, fmt.Sprintf("could not remove %q", file))
		}
	}
	switch len(errArray) > 0 {
	case true:
		err := fmt.Errorf("%v", strings.Join(errArray, "\n"))
		return err
	}

	return nil
}

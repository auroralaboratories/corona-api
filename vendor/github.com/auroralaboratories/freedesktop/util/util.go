package util

import (
	"os"
)

// Retrive a named environment variable or return a fallback value in one call
func Getenv(name string, fallback string) string {
	if rv := os.Getenv(name); rv == `` {
		return fallback
	} else {
		return rv
	}
}

// Attempt to open a given file read-only and close it to verify existence and readability
func FileExistsAndIsReadable(name string) bool {
	file, err := os.Open(name)
	defer file.Close()
	return (err == nil)
}

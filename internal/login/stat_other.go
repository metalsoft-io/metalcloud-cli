//go:build !unix

package login

import "os"

// unixStat returns false: Windows has no Unix modes and owners.
func unixStat(path string) (os.FileMode, int, bool) {
	return 0, 0, false
}

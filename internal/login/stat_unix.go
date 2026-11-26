//go:build unix

package login

import (
	"os"
	"syscall"
)

// unixStat returns the permission bits of path and its owner's uid. It returns false when path can't
// be looked at.
func unixStat(path string) (os.FileMode, int, bool) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, 0, false
	}
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, false
	}
	return info.Mode().Perm(), int(st.Uid), true
}

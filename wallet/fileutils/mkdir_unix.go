//go:build !windows

package fileutils

import (
	"golang.org/x/sys/unix"
	"os"
)

func MkdirAll(p string) error {
	mask := unix.Umask(0)  // umask 0000
	defer unix.Umask(mask) // recover umask
	return os.MkdirAll(p, 0777)
}

func WriteFile(p string, data []byte) error {
	ask := unix.Umask(0)   // umask 0000
	defer unix.Umask(mask) // recover umask
	return os.WriteFile(p, data, 0777)
}

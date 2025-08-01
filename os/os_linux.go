//go:build linux

package os

import (
	"github.com/anlakii/wallify/os/linux"
)

func SetWallpaper(filePath string) error {
	return linux.SetWallpaper(filePath)
}

func Resolution() (string, error) {
	return linux.Resolution()
}

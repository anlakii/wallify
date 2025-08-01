package os

import (
	"runtime"
	"github.com/anlakii/wallify/darwin"
	"github.com/anlakii/wallify/linux"
	"github.com/anlakii/wallify/types"
)

type WallpaperManager interface {
	SetWallpaper(filePath string) error
	Resolution() (types.Resolution, error)
}

func NewWallpaperManager() WallpaperManager {
	switch runtime.GOOS {
	case "darwin":
		return &darwin.DarwinManager{}
	case "linux":
		return &linux.LinuxManager{}
	}
	return nil
}

package os

import "github.com/anlakii/wallify/types"

type WallpaperManager struct{}

func (wm *WallpaperManager) SetWallpaper(filePath string) error {
	return SetWallpaper(filePath)
}

func (wm *WallpaperManager) Resolution() (types.Resolution, error) {
	res := GetResolution()
	return types.Resolution{Width: res.Width, Height: res.Height}, nil
}

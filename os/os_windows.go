package os

import (
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/anlakii/wallify/types"
)

const (
	setDesktopWallpaper  = 0x0014
	updateINIFile        = 0x01
	sendWinINIChange     = 0x02
	smCxPrimaryMonitor = 76
	smCyPrimaryMonitor = 77
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	systemParametersInfo = user32.NewProc("SystemParametersInfoW")
	getSystemMetrics     = user32.NewProc("GetSystemMetrics")
)

type WallpaperManager struct{}

func (wm *WallpaperManager) SetWallpaper(filePath string) error {
	path, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	_, _, err = systemParametersInfo.Call(
		setDesktopWallpaper,
		0,
		uintptr(unsafe.Pointer(pathPtr)),
		updateINIFile|sendWinINIChange,
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (wm *WallpaperManager) Resolution() (types.Resolution, error) {
	width, _, err := getSystemMetrics.Call(smCxPrimaryMonitor)
	if width == 0 {
		return types.Resolution{}, err
	}
	height, _, err := getSystemMetrics.Call(smCyPrimaryMonitor)
	if height == 0 {
		return types.Resolution{}, err
	}
	return types.Resolution{
		Width:  uint(width),
		Height: uint(height),
	}, nil
}

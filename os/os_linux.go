package os

import (
	"os/exec"
	"regexp"
	"strconv"
	"github.com/anlakii/wallify/types"
)

type WallpaperManager struct{}

func (wm *WallpaperManager) SetWallpaper(filePath string) error {
	cmd := exec.Command("feh", "--bg-fill", filePath)
	return cmd.Run()
}

func (wm *WallpaperManager) Resolution() (types.Resolution, error) {
	cmd := exec.Command("xrandr")
	output, err := cmd.Output()
	if err != nil {
		return types.Resolution{}, err
	}

	re := regexp.MustCompile(`(\d+)x(\d+)\s+\d+\.\d+\*\+`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) == 3 {
		width, _ := strconv.Atoi(matches[1])
		height, _ := strconv.Atoi(matches[2])
		return types.Resolution{Width: uint(width), Height: uint(height)}, nil
	}

	return types.Resolution{}, nil
}

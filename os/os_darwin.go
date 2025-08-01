//go:build darwin

package os

import (
	"fmt"
	"github.com/anlakii/wallify/os/darwin"
)

func SetWallpaper(filePath string) error {
	return darwin.SetWallpaper(filePath)
}

func Resolution() (string, error) {
	res := darwin.GetResolution()
	return fmt.Sprintf("%dx%d", res.Width, res.Height), nil
}

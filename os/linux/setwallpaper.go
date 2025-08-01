//go:build linux

package linux

import (
	"os/exec"
)

func SetWallpaper(filePath string) error {
	cmd := exec.Command("feh", "--bg-fill", filePath)
	return cmd.Run()
}

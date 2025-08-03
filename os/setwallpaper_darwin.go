package os

/*
#cgo LDFLAGS: -framework AppKit -framework Foundation
#include <stdlib.h>

int SetWallpaperForAllSpaces(const char *filePath);
*/
import "C"
import (
    "fmt"
    "unsafe"
)

func SetWallpaper(filePath string) error {
    cFilePath := C.CString(filePath)
    defer C.free(unsafe.Pointer(cFilePath))
    
    result := C.SetWallpaperForAllSpaces(cFilePath)
    
    if result != 0 {
        return fmt.Errorf("failed to set wallpaper for all spaces")
    }
    return nil
}

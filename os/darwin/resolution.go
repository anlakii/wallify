package darwin

/*
#cgo LDFLAGS: -framework AppKit -framework CoreGraphics
#include <CoreGraphics/CoreGraphics.h>
#include <stdlib.h>

typedef struct {
    int width;
    int height;
} Resolution;

Resolution GetNativeResolution() {
    Resolution res = {0, 0};
    CGDirectDisplayID displayID = CGMainDisplayID();
    CGDisplayModeRef currentMode = CGDisplayCopyDisplayMode(displayID);

    if (currentMode != NULL) {
        res.width = (int)CGDisplayModeGetPixelWidth(currentMode);
        res.height = (int)CGDisplayModeGetPixelHeight(currentMode);
        CGDisplayModeRelease(currentMode);
    }
    return res;
}
*/
import "C"

type Resolution struct {
    Width  uint
    Height uint
}

func GetResolution() Resolution {
    nativeRes := C.GetNativeResolution()
    
    res := Resolution{
        Width:  uint(nativeRes.width),
        Height: uint(nativeRes.height),
    }
    
	return res
}

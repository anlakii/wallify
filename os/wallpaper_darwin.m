#import <AppKit/AppKit.h>
#import <Foundation/Foundation.h>

int SetWallpaperForAllSpaces(const char *filePath) {
    @autoreleasepool {
        NSString *path = [NSString stringWithUTF8String:filePath];
        NSURL *url = [NSURL fileURLWithPath:path];
        NSError *error = nil;

        NSWorkspace *workspace = [NSWorkspace sharedWorkspace];
        
        NSArray *screens = [NSScreen screens];
        for (NSScreen *screen in screens) {
            NSDictionary *options = [NSDictionary dictionary];
            BOOL success = [workspace setDesktopImageURL:url forScreen:screen options:options error:&error];
            if (!success) {
                return 1;
            }
        }

        system("killall WallpaperAgent");
        
        return 0;
    }
}

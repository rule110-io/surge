#import <Cocoa/Cocoa.h>

extern void HandleFile(char *);

@interface AppDelegate : NSObject <NSApplicationDelegate, NSUserNotificationCenterDelegate>

@end

@interface Document : NSDocument

@end
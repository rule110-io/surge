#import <Cocoa/Cocoa.h>

extern void HandleFile(char *);
extern void ShowNotification(char *, char *);

@interface AppDelegate : NSObject <NSApplicationDelegate, NSUserNotificationCenterDelegate>

@end

@interface Document : NSDocument

@end
#import <Cocoa/Cocoa.h>

extern void HandleURL(char *);
extern void VisualModeSwitched();
extern char *GetOsxMode();

@interface GoPasser : NSObject
+ (void)handleGetURLEvent:(NSAppleEventDescriptor *)event;
+ (void)visualModeChanged:(NSNotification *)notif;
@end

void StartURLHandler(void);
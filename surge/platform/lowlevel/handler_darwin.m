#include "handler_darwin.h"

@implementation GoPasser
+ (void)handleGetURLEvent:(NSAppleEventDescriptor *)event
{
	HandleURL([[[event paramDescriptorForKeyword:keyDirectObject] stringValue] UTF8String]);
}
+(void)visualModeChanged:(NSNotification *)notif
{
  VisualModeSwitched();
}
@end

void StartURLHandler(void) {
	NSAppleEventManager *appleEventManager = [NSAppleEventManager sharedAppleEventManager];
    [appleEventManager setEventHandler:[GoPasser class]
                           andSelector:@selector(handleGetURLEvent:)
                         forEventClass:kInternetEventClass andEventID:kAEGetURL];
  [[NSDistributedNotificationCenter defaultCenter] addObserver:[GoPasser class] selector:@selector(visualModeChanged:) name:@"AppleInterfaceThemeChangedNotification" object:nil];

}

char* GetOsxMode(void){
  NSString *osxMode = [[NSUserDefaults standardUserDefaults] stringForKey:@"AppleInterfaceStyle"];
  return [osxMode UTF8String];
}
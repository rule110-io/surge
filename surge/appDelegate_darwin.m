#include "appDelegate_darwin.h"

@interface AppDelegate ()

@end

@implementation AppDelegate

-(BOOL)application:(NSApplication *)sender openFile:(NSString *)filename
{
   NSLog(@"%@", filename);
   YES;
}
 
-(void)application:(NSApplication *)sender openFiles:(NSArray *)filenames
{
   NSLog(@"%@", filenames);
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    // Insert code here to initialize your application
}


- (void)applicationWillTerminate:(NSNotification *)aNotification {
    // Insert code here to tear down your application
}


@end

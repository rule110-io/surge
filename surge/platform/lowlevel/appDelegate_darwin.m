#include "appDelegate_darwin.h"

@interface AppDelegate ()

@end

@implementation AppDelegate

-(BOOL)application:(NSApplication *)sender openFile:(NSString *)filename
{
   YES;
}
 
-(void)application:(NSApplication *)sender openFiles:(NSArray *)filenames
{
   NSLog(@"%@", filenames);
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification
{
    [[NSUserNotificationCenter defaultUserNotificationCenter] setDelegate:self];
}


- (void)applicationWillTerminate:(NSNotification *)aNotification {
    // Insert code here to tear down your application
}

- (BOOL)userNotificationCenter:(NSUserNotificationCenter *)center shouldPresentNotification:(NSUserNotification *)notification{
    return YES;
}




@end

@interface Document ()

@end

@implementation Document

- (instancetype)init {
    self = [super init];
    if (self) {
        // Add your subclass-specific initialization here.
    }
    return self;
}

+ (BOOL)autosavesInPlace {
    return YES;
}


- (NSString *)windowNibName {
    // Override returning the nib file name of the document
    // If you need to use a subclass of NSWindowController or if your document supports multiple NSWindowControllers, you should remove this method and override -makeWindowControllers instead.
    return @"Document";
}


- (NSData *)dataOfType:(NSString *)typeName error:(NSError **)outError {
    // Insert code here to write your document to data of the specified type. If outError != NULL, ensure that you create and set an appropriate error if you return nil.
    // Alternatively, you could remove this method and override -fileWrapperOfType:error:, -writeToURL:ofType:error:, or -writeToURL:ofType:forSaveOperation:originalContentsURL:error: instead.
    [NSException raise:@"UnimplementedMethod" format:@"%@ is unimplemented", NSStringFromSelector(_cmd)];
    return nil;
}


- (BOOL)readFromData:(NSData *)data ofType:(NSString *)typeName error:(NSError **)outError {
    NSData *dataFromFile = [data retain];
    NSString *myString = [[NSString alloc] initWithData:dataFromFile encoding:NSUTF8StringEncoding];

    NSLog(@"Data received: %@", myString);
    HandleFile([myString UTF8String]);
    // Do something about the image
    //if (outError != NULL) {
    //    NSLog(@"Error when loading data ...");
    //    *outError = [NSError errorWithDomain:NSOSStatusErrorDomain code:unimpErr userInfo:NULL];
    //    return NO;
    //}
    return YES;
}


@end

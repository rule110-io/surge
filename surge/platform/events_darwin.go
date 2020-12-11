package platform

//#cgo CFLAGS: -x objective-c
//#cgo LDFLAGS: -framework Cocoa
//#include "handler_darwin.h"
//#include "appDelegate_darwin.h"
import "C"
import (
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

//export HandleURL
func HandleURL(u *C.char) {
	labelText <- C.GoString(u)
}

//export VisualModeSwitched
func VisualModeSwitched() {
	SetVisualModeLikeOS()
}

//export HandleFile
func HandleFile(u *C.char) {
	filestring = C.GoString(u)
}

//WatchOSXHandler .
func WatchOSXHandler() {
	for true {
		if len(magnetstring) > 0 {
			//do act
			decodedMagetstring, err := url.QueryUnescape(magnetstring)
			if err != nil {
				log.Fatal(err)
				return
			}
			//go ParsePayloadString(decodedMagetstring)
			askUser("startDownloadMagnetLinks", "{files : ["+decodedMagetstring+"]}")

			//reregister URLHandler
			C.StartURLHandler()
			//stream chan string into string
			magnetstring = <-labelText
		}
		if len(filestring) > 0 {
			//decode file contents
			//push in array
			askUser("startDownloadMagnetLinks", "{files : ["+filestring+"]}")
			//reregister URLHandler
			filestring = ""
		}
		time.Sleep(time.Second)
	}
}

//SetVisualModeLikeOS .
func SetVisualModeLikeOS() int {
	mode := C.GoString(C.GetOsxMode())
	if mode == "" {
		//light mode
		DbWriteSetting("DarkMode", "false")
		wailsRuntime.Events.Emit("darkThemeEvent", "false")
		return 0
	} else {
		//dark mode
		DbWriteSetting("DarkMode", "true")
		wailsRuntime.Events.Emit("darkThemeEvent", "true")
		return 1
	}
}

// InitOSHandler .
func InitOSHandler() {
	// the event handler blocks!, so buffer the channel at least once to get the first message
	labelText = make(chan string, 1)

	//initially register OSX event handler
	C.StartURLHandler()

	//stream chan string into string
	magnetstring = <-labelText
	mode = <-appearance
}

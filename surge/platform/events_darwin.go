package platform

//#cgo CFLAGS: -x objective-c
//#cgo LDFLAGS: -framework Cocoa
//#include "handler_darwin.h"
//#include "appDelegate_darwin.h"
import "C"
import (
	"net/url"
	"time"

	"log"
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
				log.Panic(err)
				return
			}
			//go ParsePayloadString(decodedMagetstring)
			AskUser("startDownloadMagnetLinks", decodedMagetstring)

			//reregister URLHandler
			C.StartURLHandler()
			//stream chan string into string
			magnetstring = <-labelText
		}
		if len(filestring) > 0 {
			//decode file contents
			//push in array
			AskUser("startDownloadMagnetLinks", filestring)
			//reregister URLHandler
			filestring = ""
		}
		time.Sleep(time.Second)
	}
}

//SetVisualModeLikeOS .
func SetVisualModeLikeOS() {
	mode := C.GoString(C.GetOsxMode())
	if mode == "" {
		//light mode
		setVisualModeRef(0)
	} else {
		//dark mode
		setVisualModeRef(1)
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

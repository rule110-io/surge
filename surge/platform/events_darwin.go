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
	"github.com/wailsapp/wails"
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
			AskUser("startDownloadMagnetLinks", "{files : ["+decodedMagetstring+"]}")

			//reregister URLHandler
			C.StartURLHandler()
			//stream chan string into string
			magnetstring = <-labelText
		}
		if len(filestring) > 0 {
			//decode file contents
			//push in array
			AskUser("startDownloadMagnetLinks", "{files : ["+filestring+"]}")
			//reregister URLHandler
			filestring = ""
		}
		time.Sleep(time.Second)
	}
}

//SetVisualModeLikeOS .
func SetVisualModeLikeOS(wailsRuntime *wails.Runtime) int {
	mode := C.GoString(C.GetOsxMode())
	if mode == "" {
		//light mode
		return 0
	} else {
		//dark mode
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

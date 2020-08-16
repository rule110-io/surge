package surge

//#cgo CFLAGS: -x objective-c
//#cgo LDFLAGS: -framework Cocoa
//#include "handler_darwin.h"
//#include "appDelegate_darwin.h"
import "C"
import (
	"log"
	"net/url"
	"time"
)

//export HandleURL
func HandleURL(u *C.char) {
	labelText <- C.GoString(u)
}

//export VisualModeSwitched
func VisualModeSwitched() {
	setVisualModeLikeOSX()
}

//export HandleFile
func HandleFile(u *C.char) {
	filestring = C.GoString(u)
}

func watchOSXHandler() {
	for true {
		if len(magnetstring) > 0 {
			//do act
			decodedMagetstring, err := url.QueryUnescape(magnetstring)
			if err != nil {
				log.Fatal(err)
				return
			}
			go ParsePayloadString(decodedMagetstring)
			pushNotification("Received magnet link:", decodedMagetstring)
			//reregister URLHandler
			C.StartURLHandler()
			//stream chan string into string
			magnetstring = <-labelText
		}
		if len(filestring) > 0 {
			//go ParsePayloadString(filestring)
			pushNotification("File opened with content:", filestring)
			//reregister URLHandler
			filestring = ""
		}
		time.Sleep(time.Second)
	}
}

func setVisualModeLikeOSX() {
	mode := C.GoString(C.GetOsxMode())
	if mode == "" {
		//light mode
		DbWriteSetting("DarkMode", "false")
		wailsRuntime.Events.Emit("darkThemeEvent", "false")
	} else {
		//dark mode
		DbWriteSetting("DarkMode", "true")
		wailsRuntime.Events.Emit("darkThemeEvent", "true")
	}
}

func initOSXHandler() {
	// the event handler blocks!, so buffer the channel at least once to get the first message
	labelText = make(chan string, 1)

	//initially register OSX event handler
	C.StartURLHandler()

	//stream chan string into string
	magnetstring = <-labelText
	mode = <-appearance
}

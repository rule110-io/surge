package surge

//#cgo CFLAGS: -x objective-c
//#cgo LDFLAGS: -framework Foundation
//#include "handler_darwin.h"
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
		time.Sleep(time.Second)
	}
}

func initOSXHandler() {
	// the event handler blocks!, so buffer the channel at least once to get the first message
	labelText = make(chan string, 1)

	//initially register OSX event handler
	C.StartURLHandler()

	//stream chan string into string
	magnetstring = <-labelText
}

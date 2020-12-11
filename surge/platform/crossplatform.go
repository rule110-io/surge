package platform

import "github.com/wailsapp/wails"

var labelText chan string
var appearance chan string
var filestring = ""
var magnetstring = ""
var mode = ""

var wailsRuntime *wails.Runtime

type setVisualMode func(int)

var setVisualModeRef setVisualMode

// SetWailsRuntime binds the runtime
func SetWailsRuntime(runtime *wails.Runtime, setVisualModeFunc setVisualMode) {
	wailsRuntime = runtime
	setVisualModeRef = setVisualModeFunc
}

//AskUser emit ask user event
func AskUser(context string, payload string) {
	//log.Println("Emitting Event: ", "notificationEvent", title, text)
	wailsRuntime.Events.Emit("userEvent", context, payload)
}

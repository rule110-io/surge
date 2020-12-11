package platform

import "github.com/wailsapp/wails"

var labelText chan string
var appearance chan string
var filestring = ""
var magnetstring = ""
var mode = ""

var wailsRuntime *wails.Runtime

// SetWailsRuntime binds the runtime
func SetWailsRuntime(runtime *wails.Runtime) {
	wailsRuntime = runtime
}

//AskUser emit ask user event
func AskUser(context string, payload string) {
	//log.Println("Emitting Event: ", "notificationEvent", title, text)
	wailsRuntime.Events.Emit("userEvent", context, payload)
}

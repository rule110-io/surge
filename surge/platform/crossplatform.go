package platform

import "github.com/wailsapp/wails"

var labelText chan string
var appearance chan string
var filestring = ""
var magnetstring = ""
var mode = ""

//AskUser emit ask user event
func AskUser(wailsRuntime *wails.Runtime, context string, payload string) {
	//log.Println("Emitting Event: ", "notificationEvent", title, text)
	wailsRuntime.Events.Emit("userEvent", context, payload)
}

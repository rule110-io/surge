package platform

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var labelText chan string
var appearance chan string
var filestring = ""
var magnetstring = ""
var mode = ""

var wailsRuntime *context.Context

type setVisualMode func(int)

var setVisualModeRef setVisualMode

// SetWailsRuntime binds the runtime
func SetWailsRuntime(ctx *context.Context, setVisualModeFunc setVisualMode) {
	wailsRuntime = ctx
	setVisualModeRef = setVisualModeFunc
}

//AskUser emit ask user event
func AskUser(context string, payload string) {
	runtime.EventsEmit(*wailsRuntime, "userEvent", context, payload)
}

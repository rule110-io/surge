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

var wailsContext *context.Context

type setVisualMode func(int)

var setVisualModeRef setVisualMode

// SetWailsContext binds the runtime
func SetWailsContext(ctx *context.Context, setVisualModeFunc setVisualMode) {
	wailsContext = ctx
	setVisualModeRef = setVisualModeFunc
}

//AskUser emit ask user event
func AskUser(context string, payload string) {
	runtime.EventsEmit(*wailsContext, "userEvent", context, payload)
}

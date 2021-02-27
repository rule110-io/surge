package main

import (
	"os"

	"log"

	"github.com/leaanthony/mewn"
	surge "github.com/rule110-io/surge/backend"
	"github.com/rule110-io/surge/backend/platform"
	"github.com/wailsapp/wails"
)

var wailsRuntime *wails.Runtime
var arguments []string

// Stats
type Stats struct {
	log *wails.CustomLogger
}

//WailsRuntime
type WailsRuntime struct {
	runtime *wails.Runtime
}

// WailsInit
func (s *Stats) WailsInit(runtime *wails.Runtime) error {
	s.log = runtime.Log.New("Stats")
	go surge.WailsBind(runtime)

	return nil
}

//WailsShutdown (does not trigger in debug environment, found end of main to be more reliable)
/*func (s *WailsRuntime) WailsShutdown() {
	surge.StopClient()
}*/

func main() {
	defer surge.RecoverAndLog()

	keepRunning := platform.ProcessStartupArgs(os.Args, &surge.FrontendReady)
	if !keepRunning {
		return
	}

	//surge.HashFile("C:\\Users\\mitch\\Downloads\\surge_remote\\surge-0.2.0-beta.windows.zip")

	stats := &Stats{}

	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]
	log.Println(argsWithProg)
	log.Println(argsWithoutProg)

	//invoked with a download
	if len(argsWithoutProg) > 0 {
		arguments = os.Args[1:]
	}

	//Initialize folder structures on os filesystem
	newlyCreated, err := platform.InitializeFolders()
	if err != nil {
		log.Panic("Error on startup", err.Error())
	}
	surge.InitializeDb()
	surge.InitializeLog()
	defer surge.CloseDb()
	if newlyCreated {
		// seems like this is the first time starting the app
		//set tour to active and default mode to light
		surge.DbWriteSetting("Tour", "true")
		surge.DbWriteSetting("DarkMode", "false")
	}

	surge.StartClient(arguments)

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:     1144,
		Height:    790,
		Resizable: true,
		Title:     "Surge",
		JS:        js,
		CSS:       css,
		Colour:    "#131313",
	})

	app.Bind(stats)
	app.Bind(&surge.MiddlewareFunctions{})

	app.Run()

	surge.StopClient()
}

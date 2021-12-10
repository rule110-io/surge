package main

import (
	"embed"
	"os"

	"log"

	surge "github.com/rule110-io/surge/backend"
	"github.com/rule110-io/surge/backend/platform"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/dist
var assets embed.FS

var arguments []string

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

	//stats := &Stats{}

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

	app := NewApp()

	wails.Run(&options.App{
		Title:             "Surge",
		Width:             1320,
		Height:            570,
		MinWidth:          1320,
		MinHeight:         570,
		MaxWidth:          1920,
		MaxHeight:         740,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		RGBA:              &options.RGBA{R: 33, G: 37, B: 43, A: 255},
		Assets:            assets,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnShutdown:        app.shutdown,
		Bind: []interface{}{
			&surge.MiddlewareFunctions{},
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "Surge v0.5 - P2P on steroids",
				Message: "© 2020-2022 rule110",
			},
		},
	})

}

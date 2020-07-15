package main

import (
	"github.com/leaanthony/mewn"
	"github.com/rule110-io/surge-ui/surge"
	"github.com/wailsapp/wails"
)

var wailsRuntime *wails.Runtime

func getLocalFiles() []surge.File {
	return surge.LocalFiles
}

func getRemoteFiles() []surge.File {
	return surge.ListedFiles
}

func fetchRemoteFiles() {
	topicEncoded := surge.TopicEncode(surge.TestTopic)
	go surge.GetSubscriptions(topicEncoded)
}

func scanLocalFiles() {
	go surge.ScanLocal()
}

func downloadFile(Hash string) {
	go surge.DownloadFile(Hash)
}

// Stats .
type Stats struct {
	log *wails.CustomLogger
}

// WailsInit .
func (s *Stats) WailsInit(runtime *wails.Runtime) error {
	s.log = runtime.Log.New("Stats")
	go surge.Start(runtime)

	return nil
}

func main() {
	stats := &Stats{}
	surge.InitializeDb()
	defer surge.CloseDb()

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "surge-ui",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	app.Bind(getLocalFiles)
	app.Bind(getRemoteFiles)
	app.Bind(downloadFile)
	app.Bind(fetchRemoteFiles)
	app.Bind(scanLocalFiles)
	app.Bind(stats)

	app.Run()

}

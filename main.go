package main

import (
	"github.com/leaanthony/mewn"
	"github.com/rule110-io/surge-ui/surge"
	"github.com/wailsapp/wails"
)

var wailsRuntime *wails.Runtime

func getLocalFiles(Skip int, Take int) surge.SearchQueryResult {
	left := Skip
	right := Skip + Take

	if left > len(surge.LocalFiles) {
		left = len(surge.LocalFiles)
	}

	if right > len(surge.LocalFiles) {
		right = len(surge.LocalFiles)
	}

	return surge.SearchQueryResult{
		Result: surge.LocalFiles[left:right],
		Count:  len(surge.LocalFiles),
	}
}

func getRemoteFiles(Query string, Skip int, Take int) surge.SearchQueryResult {
	return surge.SearchFile(Query, Skip, Take)
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
	app.Bind(stats)
	app.Bind(getLocalFiles)
	app.Bind(getRemoteFiles)
	app.Bind(downloadFile)
	app.Bind(fetchRemoteFiles)
	app.Bind(scanLocalFiles)

	app.Run()

}

package main

import (
	"log"
	"github.com/rule110-io/surge-ui/surge"
  "github.com/leaanthony/mewn"
  "github.com/wailsapp/wails"
)

var wailsRuntime *wails.Runtime

func getLocalFiles() []surge.File {
  return surge.LocalFiles
}

func getRemoteFiles() []surge.File {
  return surge.ListedFiles
}

func getSessions() []surge.Session {
  return surge.Sessions
}

func fetchRemoteFiles() {
  topicEncoded := surge.TopicEncode(surge.TestTopic)
  go surge.GetSubscriptions(topicEncoded)
}

func scanLocalFiles() {
  go surge.ScanLocal()
}

func downloadFile(Addr string, Size int64, FileID string) {
  go surge.DownloadFile(Addr, Size, FileID)
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
  log.Println(stats)


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
  app.Bind(getSessions)
  app.Bind(downloadFile)
  app.Bind(fetchRemoteFiles)
  app.Bind(scanLocalFiles)
  app.Bind(stats)
  
  app.Run()
}


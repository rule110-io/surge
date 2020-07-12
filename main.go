package main

import (
  "github.com/leaanthony/mewn"
  "github.com/wailsapp/wails"
)

func getLocalFiles() []SurgeFile {
  return localFiles
}

func getRemoteFiles() []SurgeFile {
  return listedFiles
}

func getSessions() []SurgeSession {
  return sessions
}

func fetchRemoteFiles() {
  topicEncoded := surgeTopicEncode(testTopic)
  go getSubscriptions(topicEncoded)
}

func scanLocalFiles() {
  go surgeScanLocal()
}

func main() {
  go SurgeStart()

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
  
  app.Run()
}


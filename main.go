package main

import (
  "github.com/leaanthony/mewn"
  "github.com/wailsapp/wails"
  "time"
  "fmt"
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

func main() {
  go SurgeStart()


  time.Sleep(time.Second * 10)
  fmt.Println(localFiles)
  fmt.Println(listedFiles)
  
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
  
  app.Run()
}


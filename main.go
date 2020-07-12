package main

import (
  "github.com/leaanthony/mewn"
  "github.com/wailsapp/wails"
)

func getLocalFilesTestString() string {
  return queryPayload
}

func getSessions() []SurgeSession {
  return sessions
}

func main() {
  SurgeStart()

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
  app.Bind(getLocalFilesTestString)
  app.Bind(getSessions)
  
  app.Run()
}

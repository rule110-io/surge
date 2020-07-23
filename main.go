package main

import (
	"github.com/leaanthony/mewn"
	"github.com/rule110-io/surge-ui/surge"
	"github.com/wailsapp/wails"
)

var wailsRuntime *wails.Runtime

func getLocalFiles(Skip int, Take int) surge.LocalFilePageResult {

	trackedFiles := surge.GetTrackedFiles()

	for i := 0; i < len(trackedFiles); i++ {
		trackedFiles[i].ChunkMap = nil
	}

	left := Skip
	right := Skip + Take

	if left > len(trackedFiles) {
		left = len(trackedFiles)
	}

	if right > len(trackedFiles) {
		right = len(trackedFiles)
	}

	return surge.LocalFilePageResult{
		Result: trackedFiles[left:right],
		Count:  len(trackedFiles),
	}
}

func getRemoteFiles(Query string, Skip int, Take int) surge.SearchQueryResult {
	return surge.SearchFile(Query, Skip, Take)
}

func getFileChunkMap(Hash string, Size int) string {
	if Size == 0 {
		Size = 400
	}
	return surge.GetFileChunkMapString(Hash, Size)
}

func downloadFile(Hash string) bool {
	return surge.DownloadFile(Hash)
}

func setDownloadPause(Hash string, State bool) {
	surge.SetFilePause(Hash, State)
}

func openFile(Hash string) {
	surge.OpenFileByHash(Hash)
}

func openFolder(Hash string) {
	surge.OpenFolderByHash(Hash)
}

func seedFile() bool {
	path, err := surge.OpenFileDialog()
	if err != nil {
		return false
	}
	return surge.SeedFile(path)
}

func removeFile(Hash string, FromDisk bool) bool {
	return surge.RemoveFile(Hash, FromDisk)
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
	app.Bind(setDownloadPause)
	app.Bind(openFile)
	app.Bind(openFolder)
	app.Bind(getFileChunkMap)
	app.Bind(seedFile)
	app.Bind(removeFile)

	app.Run()

}

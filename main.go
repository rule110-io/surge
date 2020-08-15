package main

import (
	"log"
	"os"

	"github.com/leaanthony/mewn"
	"github.com/rule110-io/surge-ui/surge"
	"github.com/wailsapp/wails"
)

var wailsRuntime *wails.Runtime
var arguments []string

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
	return surge.GetFileChunkMapStringByHash(Hash, Size)
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

//RemoteClientOnlineModel holds info of remote clients
type RemoteClientOnlineModel struct {
	NumKnown  int
	NumOnline int
}

func getNumberOfRemoteClient() RemoteClientOnlineModel {
	total, online := surge.GetNumberOfRemoteClient()

	return RemoteClientOnlineModel{
		NumKnown:  total,
		NumOnline: online,
	}
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
	go surge.Start(runtime, arguments)

	return nil
}

func main() {
	stats := &Stats{}
	surge.InitializeDb()
	defer surge.CloseDb()

	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]
	log.Println(argsWithProg)
	log.Println(argsWithoutProg)
	//test string
	//arguments = []string{"surge://|file|Futurama.S07E01.BRRip.x264-ION10.mp4|219091405|cd0731496277102a869dacb0e99b7708c2b708824b647ffeb267de4743b7856e|e5579685272ad6d162d263a498da6eda0f35db97626dc2ecff788e9675298b67|/"}

	//invoked with a download
	if len(argsWithoutProg) > 0 {
		arguments = os.Args[1:]
	}

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:     1144,
		Height:    768,
		Resizable: false,
		Title:     "surge-ui",
		JS:        js,
		CSS:       css,
		Colour:    "#131313",
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
	app.Bind(getNumberOfRemoteClient)

	app.Run()

}

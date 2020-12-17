package main

import (
	"bufio"
	"fmt"

	"github.com/leaanthony/mewn"
	"github.com/rule110-io/surge-ui/mailslot"
	"github.com/rule110-io/surge-ui/surge"
	"github.com/rule110-io/surge-ui/surge/platform"
	log "github.com/sirupsen/logrus"
	"github.com/sqweek/dialog"
	"github.com/wailsapp/wails"
)

var wailsRuntime *wails.Runtime
var arguments []string

const instanceid = "0a7332a0-8ddf-450e-9ff4-306cdca3278e"
const mailSlotName = `\\.\mailslot\surgeclient`

func getLocalFiles(Skip int, Take int) surge.LocalFilePageResult {

	trackedFiles := surge.GetTrackedFiles()

	totalNum := len(trackedFiles)

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

	//Subset
	trackedFiles = trackedFiles[left:right]

	for i := 0; i < len(trackedFiles); i++ {
		surge.ListedFilesLock.Lock()

		for _, file := range surge.ListedFiles {
			trackedFiles[i].Seeders = []string{surge.GetMyAddress()}
			if file.FileHash == trackedFiles[i].FileHash {
				trackedFiles[i].Seeders = file.Seeders
				trackedFiles[i].Seeders = append(trackedFiles[i].Seeders, surge.GetMyAddress())
				trackedFiles[i].SeederCount = len(trackedFiles[i].Seeders)
				break
			}
		}

		if len(trackedFiles[i].Seeders) == 0 && (trackedFiles[i].IsUploading || trackedFiles[i].IsHashing) {
			trackedFiles[i].Seeders = []string{surge.GetMyAddress()}
			trackedFiles[i].SeederCount = len(trackedFiles[i].Seeders)
		}

		surge.ListedFilesLock.Unlock()
	}

	return surge.LocalFilePageResult{
		Result: trackedFiles,
		Count:  totalNum,
	}
}

func getRemoteFiles(Query string, OrderBy string, IsDesc bool, Skip int, Take int) surge.SearchQueryResult {
	return surge.SearchFile(Query, OrderBy, IsDesc, Skip, Take)
}

func getPublicKey() string {
	return surge.GetMyAddress()
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

func openLink(Link string) {
	surge.OpenOSPath(Link)
}

func openLog() {
	surge.OpenLogFile()
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
	path := surge.OpenFileDialog()
	if path == "" {
		return false
	}
	return surge.SeedFile(path)
}

func removeFile(Hash string, FromDisk bool) bool {
	return surge.RemoveFile(Hash, FromDisk)
}

func writeSetting(Key string, Value string) bool {
	err := surge.DbWriteSetting(Key, Value)
	return err != nil
}

func readSetting(Key string) string {
	val, _ := surge.DbReadSetting(Key)
	return val
}

func startDownloadMagnetLinks(Magnetlinks string) bool {
	//need to parse Magnetlinks array and download all of them
	go surge.ParsePayloadString(Magnetlinks)
	return true
}

// Stats .
type Stats struct {
	log *wails.CustomLogger
}

// WailsInit .
func (s *Stats) WailsInit(runtime *wails.Runtime) error {
	s.log = runtime.Log.New("Stats")
	go surge.WailsBind(runtime)

	return nil
}

//WailsRuntime .
type WailsRuntime struct {
	runtime *wails.Runtime
}

//WailsShutdown .
func (s *WailsRuntime) WailsShutdown() {
	surge.Stop()
}

func processStartupArgs(args []string) bool {
	argsWithProg := args
	argsWithoutProg := args[1:]
	log.Println(argsWithProg)
	log.Println(argsWithoutProg)

	singleInstanceErr := mailslot.SingleInstance(instanceid)
	isSingleInstance := singleInstanceErr == nil

	if isSingleInstance {
		go mailSlotListen()
	}

	//In the case of no arguments we either exit when surge is running or startup main client
	if len(argsWithoutProg) == 0 {
		if !isSingleInstance {
			dialog.Message("%s", "Another instance of surge is already running").Title("Surge Startup Failed").Error()
			return false
		}
		return true
	}

	//In case of arguments we will either notify our main client, or if that is not running become the mainclient.
	if !isSingleInstance {
		mailSlotSend(argsWithoutProg[0])
		return false
	}

	//If single instance and params, we start main and askUser
	//platform.AskUser("startDownloadMagnetLinks", "{files : ["+argsWithoutProg[0]+"]}")

	return true
}

func mailSlotListen() {
	ms, err := mailslot.New(mailSlotName, 0, mailslot.MAILSLOT_WAIT_FOREVER)
	defer ms.Close()

	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(ms)

	for true {
		testData, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Println(testData)

		//surge.ParsePayloadString(testData)

		//Insert new file from arguments and start download
		platform.AskUser("startDownloadMagnetLinks", "{files : ["+testData+"]}")
	}
}

func mailSlotSend(message string) {
	ms, err := mailslot.Open(mailSlotName)
	defer ms.Close()

	if err != nil {
		panic(err)
	}

	ms.Write([]byte(message + "\n"))
}

func main() {
	keepRunning := processStartupArgs([]string{"sdfsdf", "surge://|file|Project Hospital.iso|204963840|8fa0743b91b218b21fd21f65cfc18e862229d32f881c9832b6ca3d46cb1f5416|9f5471e6de22e1f1c8d3b8c4216167fade2cabd17c6fe53d950eef3147c148b7|/"})

	if !keepRunning {
		return
	}

	stats := &Stats{}
	surge.InitializeDb()
	surge.InitializeLog()
	defer surge.CloseDb()

	surge.Start()

	js := mewn.String("./frontend/dist/app.js")
	css := mewn.String("./frontend/dist/app.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:     1144,
		Height:    768,
		Resizable: false,
		Title:     "Surge",
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
	app.Bind(openLink)
	app.Bind(openLog)
	app.Bind(openFolder)
	app.Bind(getFileChunkMap)
	app.Bind(seedFile)
	app.Bind(removeFile)
	app.Bind(getNumberOfRemoteClient)
	app.Bind(writeSetting)
	app.Bind(readSetting)
	app.Bind(startDownloadMagnetLinks)
	app.Bind(getPublicKey)

	app.Run()

}

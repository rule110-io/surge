package surge

import (
	"fmt"
	"os"
	"sync"
	"time"

	"log"

	"github.com/rule110-io/surge-ui/surge/platform"
	"github.com/rule110-io/surge-ui/surge/sessionmanager"

	bitmap "github.com/boljen/go-bitmap"
	movavg "github.com/mxmCherry/movavg"
	"github.com/wailsapp/wails"
)

//FrontendReady is a flag to check if frontend is ready
var FrontendReady = false
var subscribers []string

var localFileName string
var sendSize int64
var receivedSize int64

var startTime = time.Now()

var clientOnlineMap map[string]bool

var downloadBandwidthAccumulator map[string]int
var uploadBandwidthAccumulator map[string]int

var fileBandwidthMap map[string]BandwidthMA

var zeroBandwidthMap map[string]bool

var chunksInTransit map[string]bool

var clientOnlineMapLock = &sync.Mutex{}
var clientOnlineRefreshingLock = &sync.Mutex{}
var bandwidthAccumulatorMapLock = &sync.Mutex{}
var chunkInTransitLock = &sync.Mutex{}

//Sessions .
//var Sessions []*sessionmanager.Session

//var testReader *bufio.Reader

var workerCount = 0

//Sessions collection lock
var sessionsWriteLock = &sync.Mutex{}

//ListedFilesLock lock this whenever you're reading or mutating the ListedFiles collection
var ListedFilesLock = &sync.Mutex{}

// File holds all info of a tracked file in surge
type File struct {
	FileName      string
	FileSize      int64
	FileHash      string
	Path          string
	NumChunks     int
	IsDownloading bool
	IsUploading   bool
	IsPaused      bool
	IsMissing     bool
	IsHashing     bool
	ChunkMap      []byte
	ChunksShared  int
	seeders       []string
	seederCount   int
}

//NumClientsStruct .
type NumClientsStruct struct {
	Online int
}

// FileListing struct for all frontend file listing props
type FileListing struct {
	FileName     string
	FileSize     int64
	FileHash     string
	Seeders      []string
	NumChunks    int
	IsTracked    bool
	IsAvailable  bool
	SeederCount  int
	ChunksShared int
}

// LocalFileListing is a wrapper for a local db file for the frontend
type LocalFileListing struct {
	FileName      string
	FileSize      int64
	FileHash      string
	Path          string
	NumChunks     int
	IsDownloading bool
	IsUploading   bool
	IsPaused      bool
	IsMissing     bool
	IsHashing     bool
	ChunksShared  int
	Seeders       []string
	SeederCount   int
	Progress      float32
}

// FileStatusEvent holds update info on download progress
type FileStatusEvent struct {
	FileHash          string
	Progress          float32
	Status            string
	DownloadBandwidth int
	UploadBandwidth   int
	NumChunks         int
	ChunkMap          string
	ChunksShared      int
}

//BandwidthMA tracks moving average for download and upload bandwidth
type BandwidthMA struct {
	Download movavg.MA
	Upload   movavg.MA
}

//ListedFiles are remote files that can be downloaded
var ListedFiles []File

var wailsRuntime *wails.Runtime

var numClientsStore *wails.Store

// WailsBind is a binding function at startup
func WailsBind(runtime *wails.Runtime) {
	wailsRuntime = runtime
	platform.SetWailsRuntime(wailsRuntime, SetVisualMode)

	//Mac specific functions
	go platform.InitOSHandler()
	platform.SetVisualModeLikeOS()

	numClients := NumClientsStruct{
		Online: 0,
	}

	numClientsStore = wailsRuntime.Store.New("numClients", numClients)

	updateNumClientStore()

	//Wait for our client to initialize, perhaps there is no internet connectivity
	tryCount := 1
	for !clientInitialized {
		time.Sleep(time.Second)
		if tryCount%10 == 0 {
			pushError("Connection to NKN not yet established 0", "do you have an active internet connection?")
		}
		tryCount++
	}
	updateNumClientStore()

	//Get subs first synced then grab file queries for those subs
	GetSubscriptions()

	//Startup async processes to continue processing subs/files and updating gui
	go updateGUI()
	go rescanPeers()

	FrontendReady = true
}

//SetVisualMode Sets the visualmode
func SetVisualMode(visualMode int) {
	if visualMode == 0 {
		//light mode
		DbWriteSetting("DarkMode", "false")
		wailsRuntime.Events.Emit("darkThemeEvent", "false")
	} else if visualMode == 1 {
		//dark mode
		DbWriteSetting("DarkMode", "true")
		wailsRuntime.Events.Emit("darkThemeEvent", "true")
	}
}

// Start initializes surge
func Start(args []string) {

	//Initialize all our global data maps
	clientOnlineMap = make(map[string]bool)
	downloadBandwidthAccumulator = make(map[string]int)
	uploadBandwidthAccumulator = make(map[string]int)
	zeroBandwidthMap = make(map[string]bool)
	fileBandwidthMap = make(map[string]BandwidthMA)
	chunksInTransit = make(map[string]bool)

	//Initialize our surge nkn client
	go InitializeClient(args)
}

func chunkMapFull(s []byte, num int) bool {
	//No chunkmap means no download was initiated, all chunks are local
	if s == nil {
		return true
	}

	for i := 0; i < num; i++ {
		if bitmap.Get(s, i) == false {
			return false
		}
	}
	return true
}

func chunksDownloaded(s []byte, num int) int {
	//No chunkmap means no download was initiated, all chunks are local
	if s == nil {
		return num
	}

	chunksLocalNum := 0
	for i := 0; i < num; i++ {
		if bitmap.Get(s, i) == true {
			chunksLocalNum++
		}
	}
	return chunksLocalNum
}

func updateGUI() {

	for true {
		time.Sleep(time.Second)

		log.Println("Active Workers:", workerCount)
		fmt.Println("Active Workers:", workerCount)

		log.Println("Active Sessions:", sessionmanager.GetSessionLength())
		fmt.Println("Active Sessions:", sessionmanager.GetSessionLength())

		//Create session aggregate maps for file
		fileProgressMap := make(map[string]float32)

		totalDown := 0
		totalUp := 0

		//Insert uploads
		allFiles := dbGetAllFiles()
		for _, file := range allFiles {
			if file.IsUploading {
				fileProgressMap[file.FileHash] = 1
			}
			key := file.FileHash

			if file.IsPaused {
				continue
			}

			if file.IsDownloading {
				numChunksLocal := chunksDownloaded(file.ChunkMap, file.NumChunks)
				progress := float32(float64(numChunksLocal) / float64(file.NumChunks))
				fileProgressMap[file.FileHash] = progress

				if progress >= 1.0 {
					platform.ShowNotification("Download Finished", "Download for "+file.FileName+" finished!")
					pushNotification("Download Finished", file.FileName)
					file.IsDownloading = false
					file.IsUploading = true
					dbInsertFile(file)
					go AddToSeedString(file)
				}
			}

			down, up := fileBandwidth(key)
			totalDown += down
			totalUp += up

			if zeroBandwidthMap[key] == false || down+up != 0 {
				statusEvent := FileStatusEvent{
					FileHash:          key,
					Progress:          fileProgressMap[key],
					DownloadBandwidth: down,
					UploadBandwidth:   up,
					NumChunks:         file.NumChunks,
					ChunkMap:          GetFileChunkMapString(&file, 156),
					ChunksShared:      file.ChunksShared,
				}
				wailsRuntime.Events.Emit("fileStatusEvent", statusEvent)
			}

			zeroBandwidthMap[key] = down+up == 0
		}

		//log.Println("Emitting globalBandwidthUpdate: ", totalDown, totalUp)
		if zeroBandwidthMap["total"] == false || totalDown+totalUp != 0 {
			wailsRuntime.Events.Emit("globalBandwidthUpdate", totalDown, totalUp)
		}

		zeroBandwidthMap["total"] = totalDown+totalUp == 0
	}
}

func fileBandwidth(FileID string) (Download int, Upload int) {

	//Get accumulator
	bandwidthAccumulatorMapLock.Lock()
	downAccu := downloadBandwidthAccumulator[FileID]
	downloadBandwidthAccumulator[FileID] = 0

	upAccu := uploadBandwidthAccumulator[FileID]
	uploadBandwidthAccumulator[FileID] = 0
	bandwidthAccumulatorMapLock.Unlock()

	if fileBandwidthMap[FileID].Download == nil {
		fileBandwidthMap[FileID] = BandwidthMA{
			Download: movavg.ThreadSafe(movavg.NewSMA(10)),
			Upload:   movavg.ThreadSafe(movavg.NewSMA(10)),
		}
	}

	fileBandwidthMap[FileID].Download.Add(float64(downAccu))
	fileBandwidthMap[FileID].Upload.Add(float64(upAccu))

	return int(fileBandwidthMap[FileID].Download.Avg()), int(fileBandwidthMap[FileID].Upload.Avg())
}

//ByteCountSI converts filesize in bytes to human readable text
func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func getFileSize(path string) (size int64) {
	fi, err := os.Stat(path)
	if err != nil {
		return -1
	}
	// get the size
	return fi.Size()
}

// Stats .
type Stats struct {
	log *wails.CustomLogger
}

func getListedFileByHash(Hash string) *File {

	var selectedFile *File = nil

	ListedFilesLock.Lock()
	for _, file := range ListedFiles {
		if file.FileHash == Hash {
			selectedFile = &file
			break
		}
	}
	ListedFilesLock.Unlock()

	return selectedFile
}

//GetFileChunkMapString returns the chunkmap in hex for a file given by hash
func GetFileChunkMapString(file *File, Size int) string {

	outputSize := Size
	inputSize := file.NumChunks

	stepSize := float64(inputSize) / float64(outputSize)
	stepSizeInt := int(stepSize)

	var boolBuffer = ""
	if inputSize >= outputSize {

		for i := 0; i < outputSize; i++ {
			localCount := 0
			for j := 0; j < stepSizeInt; j++ {
				local := bitmap.Get(file.ChunkMap, int(float64(i)*stepSize)+j)
				if local {
					localCount++
				} else {
					if localCount == 0 {
						boolBuffer += "0"
					} else {
						boolBuffer += "1"
					}
					break
				}
			}
			if localCount == stepSizeInt {
				boolBuffer += "2"
			}
		}
	} else {
		iter := float64(0)
		for i := 0; i < outputSize; i++ {
			local := bitmap.Get(file.ChunkMap, int(iter))
			if local {
				boolBuffer += "2"
			} else {
				boolBuffer += "0"
			}
			iter += stepSize
		}
	}
	return boolBuffer
}

//GetFileChunkMapStringByHash returns the chunkmap in hex for a file given by hash
func GetFileChunkMapStringByHash(Hash string, Size int) string {

	file, err := dbGetFile(Hash)
	if err != nil {
		return ""
	}
	return GetFileChunkMapString(file, Size)
}

//SetFilePause sets a file IsPaused state for by file hash
func SetFilePause(Hash string, State bool) {

	fileWriteLock.Lock()
	file, err := dbGetFile(Hash)
	if err != nil {
		pushNotification("Failed To Pause", "Could not find the file to pause.")
		return
	}
	file.IsPaused = State
	dbInsertFile(*file)
	fileWriteLock.Unlock()

	msg := "Paused"
	if State == false {
		msg = "Resumed"
	}
	pushNotification("Download "+msg, file.FileName)
}

//RemoveFile removes file from surge db and optionally from disk
func RemoveFile(Hash string, FromDisk bool) bool {

	fileWriteLock.Lock()

	if FromDisk {
		file, err := dbGetFile(Hash)
		if err != nil {
			log.Println("Error on remove file (read db)", err.Error())
			pushError("Error on remove file (read db)", err.Error())
			return false
		}
		err = os.Remove(file.Path)
		if err != nil {
			log.Println("Error on remove file from disk", err.Error())
			pushError("Error on remove file from disk", err.Error())
		}
	}

	err := dbDeleteFile(Hash)
	if err != nil {
		log.Println("Error on remove file (read db)", err.Error())
		pushError("Error on remove file (read db)", err.Error())
		return false
	}
	fileWriteLock.Unlock()

	//Rebuild entirely
	dbFiles := dbGetAllFiles()
	go BuildSeedString(dbFiles)

	return true
}

//GetMyAddress returns current client address
func GetMyAddress() string {
	for !clientInitialized {
		time.Sleep(time.Millisecond * 50)
	}
	return client.Addr().String()
}

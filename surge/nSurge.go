package surge

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"log"

	"github.com/rule110-io/surge-ui/surge/constants"
	"github.com/rule110-io/surge-ui/surge/platform"
	"github.com/rule110-io/surge-ui/surge/sessionmanager"

	bitmap "github.com/boljen/go-bitmap"
	movavg "github.com/mxmCherry/movavg"
	"github.com/wailsapp/wails"
)

// SurgeActive is true when client is operational
var SurgeActive bool = false
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
	ChunkMap      []byte
	ChunksShared  int
	Seeders       []string
	SeederCount   int
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
	tryCount := 0
	for !clientInitialized {
		if tryCount%10 == 0 {
			pushError("Connection to NKN not yet established", "do you have an active internet connection?")
		}
		time.Sleep(time.Second)
		tryCount++
	}
	updateNumClientStore()

	//Get subs first synced then grab file queries for those subs
	GetSubscriptions()

	//Startup async processes to continue processing subs/files and updating gui
	go updateGUI()
	go rescanPeers()
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

//DownloadFile downloads the file
func DownloadFile(Hash string) bool {

	//Addr string, Size int64, FileID string
	file := getListedFileByHash(Hash)
	if file == nil {
		pushError("Error on download file", "No listed file with hash: "+Hash)
	}

	fileSeeders := file.seeders

	pushNotification("Download Started", file.FileName)

	remoteFolder, err := platform.GetRemoteFolder()
	if err != nil {
		log.Println("Remote folder does not exist")
	}

	// If the file doesn't exist allocate it
	var path = remoteFolder + string(os.PathSeparator) + file.FileName
	AllocateFile(path, file.FileSize)
	numChunks := int((file.FileSize-1)/int64(ChunkSize)) + 1

	//When downloading from remote enter file into db
	dbFile, err := dbGetFile(Hash)
	log.Println(dbFile)
	if err != nil {
		file.Path = path
		file.NumChunks = numChunks
		file.ChunkMap = bitmap.NewSlice(numChunks)
		file.IsDownloading = true
		dbInsertFile(*file)
	}

	//Create a random fetch sequence
	randomChunks := make([]int, numChunks)
	for i := 0; i < numChunks; i++ {
		randomChunks[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(randomChunks), func(i, j int) { randomChunks[i], randomChunks[j] = randomChunks[j], randomChunks[i] })

	seederAlternator := 0
	mutateSeederLock := sync.Mutex{}
	appendChunkLock := sync.Mutex{}

	//Give the seeder a fair start with timers when a download is initiated
	//Potentionally this seeder was last queried 60 seconds ago for files and otherwise idle but online
	for _, seeder := range fileSeeders {
		sessionmanager.UpdateActivity(seeder)
	}

	downloadJob := func(terminateFlag *bool) {

		//Used to terminate the rescanning of peers
		terminate := func(flag *bool) {
			*flag = true
		}
		defer terminate(terminateFlag)

		for i := 0; i < numChunks; i++ {

			newFileData := getListedFileByHash(Hash)
			if newFileData != nil {
				file = newFileData
			}

			for len(fileSeeders) == 0 {
				time.Sleep(time.Second * 5)
			}

			dbFile, err := dbGetFile(file.FileHash)

			//Check if file is still tracked in surge
			if err != nil {
				log.Println("Download Job Treminated", "File no longer in DB")
				return
			}

			//Pause if file is paused
			for err == nil && dbFile.IsPaused {
				time.Sleep(time.Second * 5)
				dbFile, err = dbGetFile(file.FileHash)
				if err != nil {
					break
				}
			}

			for workerCount >= NumWorkers {
				time.Sleep(time.Millisecond)
			}
			workerCount++

			//Create a async job to download a chunk
			requestChunkJob := func(chunkID int) {

				success := false
				downloadSeederAddr := ""

				if len(fileSeeders) > seederAlternator {
					//Get seeder
					downloadSeederAddr = fileSeeders[seederAlternator]
					session, existing := sessionmanager.GetExistingSession(downloadSeederAddr, constants.WorkerGetSessionTimeout)

					if existing {
						success = RequestChunk(session, file.FileHash, int32(chunkID))
					} else {
						success = false
					}
				}

				//if download fails append the chunk to remaining to retry later
				if !success {
					appendChunkLock.Lock()
					randomChunks = append(randomChunks, chunkID)
					numChunks++
					appendChunkLock.Unlock()

					workerCount--

					//TODO: Think about alternatives from straight dropping the seeder.
					mutateSeederLock.Lock()
					fileSeeders = removeStringFromSlice(fileSeeders, downloadSeederAddr)
					mutateSeederLock.Unlock()

					//return out of job
					return
				}

				//If chunk is requested add to transit map
				chunkKey := file.FileHash + "_" + strconv.Itoa(chunkID)

				chunkInTransitLock.Lock()
				chunksInTransit[chunkKey] = true
				chunkInTransitLock.Unlock()

				//Sleep and check if entry still exists in transit map.
				sleepWorker := true
				inTransit := true
				receiveTimeoutCounter := 0

				for sleepWorker {
					time.Sleep(time.Second)
					//fmt.Println(string("\033[36m"), "Worker Sleeping", string("\033[0m"))

					//Check if connection is lost
					_, sessionExists := sessionmanager.GetExistingSession(downloadSeederAddr, constants.WorkerGetSessionTimeout)
					if !sessionExists {
						//if session no longer exists
						fmt.Println(string("\033[36m"), "session no longer exists", string("\033[0m"))
						inTransit = true
						sleepWorker = false
						break
					}

					//Check if received
					isInTransit := chunksInTransit[chunkKey]
					if !isInTransit {
						//if no longer in transit, continue workers
						fmt.Println(string("\033[36m"), "no longer in transit, continue workers", string("\033[0m"))
						inTransit = false
						sleepWorker = false
						break
					} else if receiveTimeoutCounter >= constants.WorkerChunkReceiveTimeout {
						//if timeout is triggered, leave in transit.
						fmt.Println(string("\033[36m"), "timeout is triggered, leave in transit.", string("\033[0m"))
						inTransit = true
						sleepWorker = false
						break
					}
					receiveTimeoutCounter++
				}

				//If its still in transit abort
				if inTransit {
					appendChunkLock.Lock()
					randomChunks = append(randomChunks, chunkID)
					numChunks++
					appendChunkLock.Unlock()

					workerCount--

					//TODO: Think about alternatives from straight dropping the seeder.
					mutateSeederLock.Lock()
					fileSeeders = removeStringFromSlice(fileSeeders, downloadSeederAddr)
					mutateSeederLock.Unlock()

					//return out of job
					return
				}
			}

			//get chunk id
			appendChunkLock.Lock()
			chunkid := randomChunks[i]
			appendChunkLock.Unlock()

			go requestChunkJob(chunkid)

			mutateSeederLock.Lock()
			seederAlternator++
			if seederAlternator > len(fileSeeders)-1 {
				seederAlternator = 0
			}
			mutateSeederLock.Unlock()
		}
	}

	scanForSeeders := func(terminateFlag *bool) {

		//While we are not terminated scan for new peers
		for *terminateFlag == false {
			time.Sleep(time.Second * 5)

			newFile := getListedFileByHash(Hash)
			if newFile != nil {
				//Check for new sessions
				mutateSeederLock.Lock()
				fileSeeders = []string{}
				for i := 0; i < len(newFile.seeders); i++ {
					_, existing := sessionmanager.GetExistingSession(newFile.seeders[i], 60)
					if existing {
						fileSeeders = append(fileSeeders, newFile.seeders[i])
					}
				}
				mutateSeederLock.Unlock()
			}
		}
	}

	terminateFlag := false
	go downloadJob(&terminateFlag)
	go scanForSeeders(&terminateFlag)

	return true
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

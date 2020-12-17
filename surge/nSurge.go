package surge

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rule110-io/surge-ui/surge/platform"
	log "github.com/sirupsen/logrus"

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
var Sessions []*Session

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
	Seeders       []string
	Path          string
	NumChunks     int
	IsDownloading bool
	IsUploading   bool
	IsPaused      bool
	IsMissing     bool
	IsHashing     bool
	ChunkMap      []byte
	SeederCount   int
	ChunksShared  int
}

//NumClientsStruct .
type NumClientsStruct struct {
	Subscribed int
	Online     int
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

// Session is a wrapper for everything needed to maintain a surge session
type Session struct {
	FileHash   string
	FileSize   int64
	Downloaded int64
	Uploaded   int64
	session    net.Conn
	reader     *bufio.Reader
	file       *os.File
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

var numClientsSubscribed int = 0
var numClientsOnline int = 0

var numClientsStore *wails.Store

// WailsBind is a binding function at startup
func WailsBind(runtime *wails.Runtime) {
	wailsRuntime = runtime
	platform.SetWailsRuntime(wailsRuntime, SetVisualMode)

	//Mac specific functions
	go platform.InitOSHandler()
	platform.SetVisualModeLikeOS()

	numClients := NumClientsStruct{
		Subscribed: 0,
		Online:     0,
	}

	numClientsStore = wailsRuntime.Store.New("numClients", numClients)

	//Get subs first synced then grab file queries for those subs
	GetSubscriptions()
	go queryRemoteForFiles()

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
func Start() {

	//Initialize all our global data maps
	clientOnlineMap = make(map[string]bool)
	downloadBandwidthAccumulator = make(map[string]int)
	uploadBandwidthAccumulator = make(map[string]int)
	zeroBandwidthMap = make(map[string]bool)
	fileBandwidthMap = make(map[string]BandwidthMA)
	chunksInTransit = make(map[string]bool)

	//Initialize folder structures on os filesystem
	newlyCreated, err := platform.InitializeFolders()
	if err != nil {
		pushError("Error on startup", err.Error())
	}
	if newlyCreated {
		// seems like this is the first time starting the app
		//set tour to active
		DbWriteSetting("Tour", "true")
		//set default mode to light
		DbWriteSetting("DarkMode", "false")
	}

	//Initialize our surge nkn client
	initialSuccess := InitializeClient(false)
	if !initialSuccess {
		go InitializeClient(true)
	}
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

func updateGUI() {
	defer RecoverAndLog()
	for true {
		time.Sleep(time.Second)

		//Create session aggregate maps for file
		fileProgressMap := make(map[string]float32)

		sessionsWriteLock.Lock()
		for _, session := range Sessions {
			//log.Println("Active session:", session.session.RemoteAddr().String())
			if session.FileSize == 0 {
				continue
			}

			fileProgressMap[session.FileHash] = float32(float64(session.Downloaded) / float64(session.FileSize))

			if session.Downloaded >= session.FileSize {
				fileEntry, err := dbGetFile(session.FileHash)
				if err != nil {
					pushError("Error on download complete", err.Error())
					continue
				}

				if chunkMapFull(fileEntry.ChunkMap, fileEntry.NumChunks) {
					platform.ShowNotification("Download Finished", "Download for "+fileEntry.FileName+" finished!")
					pushNotification("Download Finished", fileEntry.FileName)
					session.session.Close()

					fileEntry.IsDownloading = false
					fileEntry.IsUploading = true
					dbInsertFile(*fileEntry)
					go AddToSeedString(*fileEntry)
				}
			}
		}
		sessionsWriteLock.Unlock()

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

		log.Println("Active Workers:", workerCount)
		fmt.Println("Active Workers:", workerCount)
	}
}

func fileBandwidth(FileID string) (Download int, Upload int) {
	defer RecoverAndLog()
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

// WailsInit .
func (s *Stats) WailsInit(runtime *wails.Runtime) error {
	s.log = runtime.Log.New("Stats")
	runtime.Events.Emit("notificationEvent", "Backend Init", "just a test")
	return nil
}

func getListedFileByHash(Hash string) *File {

	defer RecoverAndLog()
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
	defer RecoverAndLog()

	file := getListedFileByHash(Hash)
	if file == nil {
		pushError("Error on download file", "No listed file with hash: "+Hash)
	}

	downloadSessions := []*Session{}

	// Create  sessions
	for i := 0; i < len(file.Seeders); i++ {
		surgeSession, err := createSession(file, file.Seeders[i])
		if err != nil {
			log.Println("Could not create session for download", Hash, file.Seeders[i])
			continue
		}
		go initiateSession(surgeSession)

		downloadSessions = append(downloadSessions, surgeSession)
	}

	if len(downloadSessions) == 0 {
		pushNotification("Download Session Failed, failed to connect to all seeders.", file.FileName)
		return false
	}

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

			for len(downloadSessions) == 0 {
				time.Sleep(time.Second * 5)
			}

			//Pause if file is paused
			dbFile, err := dbGetFile(file.FileHash)
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
				defer RecoverAndLog()

				success := false
				downloadSeeder := &Session{}
				downloadSeederAddr := ""

				if len(downloadSessions) > seederAlternator {
					//Get seeder
					downloadSeeder = downloadSessions[seederAlternator]
					if downloadSeeder != nil && downloadSeeder.session != nil {
						downloadSeederAddr = downloadSeeder.session.RemoteAddr().String()
						success = RequestChunk(downloadSeeder, file.FileHash, int32(chunkID))
					}
				}

				//if download fails append the chunk to remaining to retry later
				if !success {
					appendChunkLock.Lock()
					randomChunks = append(randomChunks, chunkID)
					numChunks++
					appendChunkLock.Unlock()

					workerCount--

					//Drop the seeder
					mutateSeederLock.Lock()
					downloadSessions = removeAndCloseSessionOrdered(downloadSessions, downloadSeederAddr)
					log.Println("Lost connection", "Dropping 1 Session for Download "+file.FileName)
					mutateSeederLock.Unlock()

					if len(downloadSessions) == 0 {
						return
					}
				}

				//If chunk is requested add to transit map
				chunkKey := file.FileHash + "_" + strconv.Itoa(chunkID)

				chunkInTransitLock.Lock()
				chunksInTransit[chunkKey] = true
				chunkInTransitLock.Unlock()

				//Sleep for 30 seconds, check if entry still exists in transit map.
				time.Sleep(time.Second * 60)
				inTransit := chunksInTransit[chunkKey]

				//If its still in transit abort
				if inTransit {
					appendChunkLock.Lock()
					randomChunks = append(randomChunks, chunkID)
					numChunks++
					appendChunkLock.Unlock()

					workerCount--

					//Drop the seeder
					mutateSeederLock.Lock()
					downloadSessions = removeAndCloseSessionOrdered(downloadSessions, downloadSeederAddr)
					log.Println("Lost connection", "Dropping 1 Session for Download "+file.FileName)
					mutateSeederLock.Unlock()

					if len(downloadSessions) == 0 {
						return
					}
				}
			}

			//get chunk id
			appendChunkLock.Lock()
			chunkid := randomChunks[i]
			appendChunkLock.Unlock()

			go requestChunkJob(chunkid)

			mutateSeederLock.Lock()
			seederAlternator++
			if seederAlternator > len(downloadSessions)-1 {
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
				for i := 0; i < len(newFile.Seeders); i++ {
					//Check if the newFile seeder is not already part of the downloadSessions
					alreadyAdded := false
					for j := 0; j < len(downloadSessions); j++ {
						if downloadSessions[j].session == nil {
							continue
						}
						if downloadSessions[j].session.RemoteAddr().String() == newFile.Seeders[i] {
							alreadyAdded = true
							break
						}
					}

					//Skip this entry
					if alreadyAdded {
						continue
					}

					surgeSession, err := createSession(newFile, newFile.Seeders[i])
					if err != nil {
						log.Println("Could not create session for download", Hash, newFile.Seeders[i])
						continue
					}

					dbFile, err := dbGetFile(Hash)
					if err == nil && dbFile != nil {
						//Prime the session with known bytes downloaded
						surgeSession.Downloaded = int64(dbFile.NumChunks-len(randomChunks)) * ChunkSize
						//If the last chunk is set, we want to deduct the missing bytes because its not a complete chunk
						lastChunkSet := bitmap.Get(dbFile.ChunkMap, dbFile.NumChunks-1)
						if lastChunkSet {
							overshotBytes := int64(dbFile.NumChunks)*int64(ChunkSize) - dbFile.FileSize
							surgeSession.Downloaded -= overshotBytes
						}

						go initiateSession(surgeSession)

						mutateSeederLock.Lock()
						downloadSessions = append(downloadSessions, surgeSession)
						mutateSeederLock.Unlock()
					}
				}
			}
		}
	}

	terminateFlag := false
	go downloadJob(&terminateFlag)
	go scanForSeeders(&terminateFlag)

	return true
}

//SearchQueryResult is a paging query result for file searches
type SearchQueryResult struct {
	Result []FileListing
	Count  int
}

//LocalFilePageResult is a paging query result for tracked files
type LocalFilePageResult struct {
	Result []File
	Count  int
}

//SearchFile runs a paged query
func SearchFile(Query string, OrderBy string, IsDesc bool, Skip int, Take int) SearchQueryResult {
	defer RecoverAndLog()
	var results []FileListing

	ListedFilesLock.Lock()
	for _, file := range ListedFiles {
		if strings.Contains(strings.ToLower(file.FileName), strings.ToLower(Query)) || strings.Contains(strings.ToLower(file.FileHash), strings.ToLower(Query)) {

			result := FileListing{
				FileName:    file.FileName,
				FileHash:    file.FileHash,
				FileSize:    file.FileSize,
				Seeders:     file.Seeders,
				NumChunks:   file.NumChunks,
				SeederCount: len(file.Seeders),
			}

			tracked, err := dbGetFile(result.FileHash)

			//only add non-local files to the result
			if err != nil && tracked == nil {
				results = append(results, result)
			}

		}
	}
	ListedFilesLock.Unlock()

	switch OrderBy {
	case "FileName":
		if !IsDesc {
			sort.Sort(sortByFileNameAsc(results))
		} else {
			sort.Sort(sortByFileNameDesc(results))
		}
	case "FileSize":
		if !IsDesc {
			sort.Sort(sortByFileSizeAsc(results))
		} else {
			sort.Sort(sortByFileSizeDesc(results))
		}
	default:
		if !IsDesc {
			sort.Sort(sortBySeederCountAsc(results))
		} else {
			sort.Sort(sortBySeederCountDesc(results))
		}
	}

	left := Skip
	right := Skip + Take

	if left > len(results) {
		left = len(results)
	}

	if right > len(results) {
		right = len(results)
	}

	return SearchQueryResult{
		Result: results[left:right],
		Count:  len(results),
	}
}

//GetTrackedFiles returns all files tracked in surge client
func GetTrackedFiles() []File {
	defer RecoverAndLog()
	return dbGetAllFiles()
}

//GetFileChunkMapString returns the chunkmap in hex for a file given by hash
func GetFileChunkMapString(file *File, Size int) string {
	defer RecoverAndLog()
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
	defer RecoverAndLog()
	file, err := dbGetFile(Hash)
	if err != nil {
		return ""
	}
	return GetFileChunkMapString(file, Size)
}

//SetFilePause sets a file IsPaused state for by file hash
func SetFilePause(Hash string, State bool) {
	defer RecoverAndLog()
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

//OpenFileDialog uses platform agnostic package for a file dialog
func OpenFileDialog() string {
	defer RecoverAndLog()
	selectedFile := wailsRuntime.Dialog.SelectFile()
	return selectedFile
}

//RemoveFile removes file from surge db and optionally from disk
func RemoveFile(Hash string, FromDisk bool) bool {
	defer RecoverAndLog()

	//Close sessions for this file
	for _, session := range Sessions {
		if session.FileHash == Hash {
			closeSession(session)
			break
		}
	}

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
			log.Println("Error on remove file from disk - removing from surge", err.Error())
			pushError("Error on remove file from disk - removing from surge", err.Error())
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

package surge

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	bitmap "github.com/boljen/go-bitmap"
	nkn "github.com/nknorg/nkn-sdk-go"
	"github.com/wailsapp/wails"
)

// SurgeActive is true when client is operational
var SurgeActive bool = false

//ChunkSize is size of chunk in bytes (256 kB)
const ChunkSize = 1024 * 256

//NumClients is the number of NKN clients
const NumClients = 8

//NumWorkers is the total number of concurrent chunk fetches allowed
const NumWorkers = 16

const localPath = "local"
const remotePath = "remote"

//OS folder permission bitflags
const (
	osRead       = 04
	osWrite      = 02
	osEx         = 01
	osUserShift  = 6
	osGroupShift = 3
	osOthShift   = 0

	osUserR   = osRead << osUserShift
	osUserW   = osWrite << osUserShift
	osUserX   = osEx << osUserShift
	osUserRw  = osUserR | osUserW
	osUserRwx = osUserRw | osUserX

	osGroupR   = osRead << osGroupShift
	osGroupW   = osWrite << osGroupShift
	osGroupX   = osEx << osGroupShift
	osGroupRw  = osGroupR | osGroupW
	osGroupRwx = osGroupRw | osGroupX

	osOthR   = osRead << osOthShift
	osOthW   = osWrite << osOthShift
	osOthX   = osEx << osOthShift
	osOthRw  = osOthR | osOthW
	osOthRwx = osOthRw | osOthX

	osAllR   = osUserR | osGroupR | osOthR
	osAllW   = osUserW | osGroupW | osOthW
	osAllX   = osUserX | osGroupX | osOthX
	osAllRw  = osAllR | osAllW
	osAllRwx = osAllRw | osGroupX
)

var localFileName string
var sendSize int64
var receivedSize int64

var startTime = time.Now()

var client *nkn.MultiClient

//Sessions .
var Sessions []*Session

//var testReader *bufio.Reader

var workerCount = 0

// File holds all file listing info of a seeded file
type File struct {
	FileName      string
	FileSize      int64
	FileHash      string
	Seeder        string
	Path          string
	NumChunks     int
	IsDownloading bool
	IsUploading   bool
	IsPaused      bool
	ChunkMap      []byte
}

// Session is a wrapper for everything needed to maintain a surge session
type Session struct {
	FileHash        string
	FileSize        int64
	Downloaded      int64
	Uploaded        int64
	deltaDownloaded int64
	bandwidthQueue  []int
	session         net.Conn
	reader          *bufio.Reader
}

// DownloadStatusEvent holds update info on download progress
type DownloadStatusEvent struct {
	FileHash  string
	Progress  float32
	Status    string
	Bandwidth int
	NumChunks int
}

//ListedFiles are remote files that can be downloaded
var ListedFiles []File

//LocalFiles are files that can be seeded
var LocalFiles []File

var wailsRuntime *wails.Runtime

// Start initializes surge
func Start(runtime *wails.Runtime) {
	wailsRuntime = runtime
	var dirFileMode os.FileMode
	dirFileMode = os.ModeDir | (osUserRwx | osAllR)

	//Ensure local and remote folders exist
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		os.Mkdir(localPath, dirFileMode)
	}
	if _, err := os.Stat(remotePath); os.IsNotExist(err) {
		os.Mkdir(remotePath, dirFileMode)
	}

	account := InitializeAccount()

	var err error

	client, err = nkn.NewMultiClient(account, "", NumClients, false, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		<-client.OnConnect.C

		pushNotification("Client Connected", "Successfully connected to the NKN network")

		client.Listen(nil)
		SurgeActive = true
		go Listen()
	}

	go ScanLocal()
	topicEncoded := TopicEncode(TestTopic)
	GetSubscriptions(topicEncoded)

	tracked := GetTrackedFiles()
	for i := 0; i < len(tracked); i++ {
		go restartDownload(tracked[i].FileHash)
	}

	go updateGUI()
}

func updateGUI() {
	for true {
		time.Sleep(time.Second)

		for _, session := range Sessions {
			if session.FileSize == 0 {
				continue
			}

			//Take bandwith delta
			deltaBandwidth := int(session.Downloaded - session.deltaDownloaded)
			session.deltaDownloaded = session.Downloaded

			//Append to queue
			session.bandwidthQueue = append(session.bandwidthQueue, deltaBandwidth)

			//Dequeue if queue > 10
			if len(session.bandwidthQueue) > 10 {
				session.bandwidthQueue = session.bandwidthQueue[1:]
			}

			var bandwidthMA10 = 0
			for i := 0; i < len(session.bandwidthQueue); i++ {
				bandwidthMA10 += session.bandwidthQueue[i]
			}
			bandwidthMA10 /= len(session.bandwidthQueue)

			fileInfo, err := dbGetFile(session.FileHash)
			if err != nil {
				log.Panicln(err)
			}

			statusEvent := DownloadStatusEvent{
				FileHash:  session.FileHash,
				Progress:  float32(float64(session.Downloaded) / float64(session.FileSize)),
				Status:    "Downloading",
				Bandwidth: bandwidthMA10,
				NumChunks: fileInfo.NumChunks,
			}
			log.Println("Emitting downloadStatusEvent: ", statusEvent)
			wailsRuntime.Events.Emit("downloadStatusEvent", statusEvent)

			//Download completed
			/*var completed = true
			for i := 0; i < fileInfo.NumChunks; i++ {
				if bitmap.Get(fileInfo.ChunkMap, i) == false {
					completed = false
					break
				}
			}*/

			if session.Downloaded == session.FileSize {
				pushNotification("Download Finished", getListedFileByHash(session.FileHash).FileName)
				session.session.Close()

				fileEntry, err := dbGetFile(session.FileHash)
				if err != nil {
					log.Panicln(err)
				}
				fileEntry.IsDownloading = false
				dbInsertFile(*fileEntry)
			}
		}
	}
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

func sendSeedSubscription(Topic string, Payload string) {
	txnHash, err := client.Subscribe("", Topic, 4320, Payload, nil)
	if err != nil {
		log.Println("Probably already subscribed", err)
	} else {
		log.Println("Subscribed: ", txnHash)
	}
}

//GetSubscriptions .
func GetSubscriptions(Topic string) {
	//Empty file cache
	ListedFiles = []File{}
	//fileBox.Children = []fyne.CanvasObject{}

	subscribers, err := client.GetSubscribers(Topic, 0, 100, true, true)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range subscribers.SubscribersInTxPool.Map {
		subscribers.Subscribers.Map[k] = v
	}

	for k, v := range subscribers.Subscribers.Map {
		if len(v) > 0 {
			SendQueryRequest(k, "Testing query functionality.")
		}
	}

	fmt.Println(ListedFiles)

}

// Stats .
type Stats struct {
	log *wails.CustomLogger
}

// WailsInit .
func (s *Stats) WailsInit(runtime *wails.Runtime) error {
	s.log = runtime.Log.New("Stats")
	runtime.Events.Emit("notificationEvent", "Backend Init", "just a test")
	log.Println("TESTING TESTING TESTING")
	return nil
}

func getListedFileByHash(Hash string) *File {
	for _, file := range ListedFiles {
		if file.FileHash == Hash {
			return &file
		}
	}
	return nil
}

//DownloadFile downloads the file
func DownloadFile(Hash string) {
	//Addr string, Size int64, FileID string

	file := getListedFileByHash(Hash)
	if file == nil {
		log.Panic("No listed file with hash", Hash)
	}

	// Create a sessions
	surgeSession, err := createSession(file)
	if err != nil {
		log.Println("Could not create session for download", Hash)
		pushNotification("Download Session Failed", file.FileName)
		return
	}
	go initiateSession(surgeSession)

	pushNotification("Download Started", file.FileName)

	// If the file doesn't exist allocate it
	var path = remotePath + "/" + file.FileName
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

	downloadJob := func() {
		for i := 0; i < numChunks; i++ {
			workerCount++
			go RequestChunk(surgeSession, file.FileHash, int32(i))

			for workerCount >= NumWorkers {
				time.Sleep(time.Millisecond)
			}
		}
	}
	go downloadJob()
}

func pushNotification(title string, text string) {
	log.Println("Emitting Event: ", "notificationEvent", title, text)
	wailsRuntime.Events.Emit("notificationEvent", title, text)
}

//SearchQueryResult is a paging query result for file searches
type SearchQueryResult struct {
	Result []File
	Count  int
}

//SearchFile runs a paged query
func SearchFile(Query string, Skip int, Take int) SearchQueryResult {
	var results []File

	for _, file := range ListedFiles {
		if strings.Contains(strings.ToLower(file.FileName), strings.ToLower(Query)) {
			results = append(results, file)
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
	return dbGetAllFiles()
}

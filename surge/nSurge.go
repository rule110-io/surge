package surge

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
	nkn "github.com/nknorg/nkn-sdk-go"
	"github.com/wailsapp/wails"
)

// SurgeActive is true when client is operational
var SurgeActive bool = false

//ChunkSize is size of chunk in bytes (256 kB)
const ChunkSize = 1024 * 256

//NumClients is the number of NKN clients
const NumClients = 16

const localPath = "local"
const remotePath = "remote"

var localFileName string
var sendSize int64
var receivedSize int64

var startTime = time.Now()

var client *nkn.MultiClient

//Sessions .
var Sessions []Session

//var testReader *bufio.Reader

var workerCount = 0

var chunksTotal int
var chunksRequested int
var chunksReceived int

// File holds all file listing info of a seeded file
type File struct {
	FileName string
	FileSize int64
	MD5Hash  string
	Seeder   string
}

// Session is a wrapper for everything needed to maintain a surge session
type Session struct {
	session net.Conn
	reader  *bufio.Reader
}

//ListedFiles are remote files that can be downloaded
var ListedFiles []File
//LocalFiles are files that can be seeded
var LocalFiles []File

var wailsRuntime *wails.Runtime

// Start initializes surge
func Start(runtime *wails.Runtime) {
	wailsRuntime = runtime

	//Ensure local and remote folders exist
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		os.Mkdir(localPath, os.ModeDir)
	}
	if _, err := os.Stat(remotePath); os.IsNotExist(err) {
		os.Mkdir(remotePath, os.ModeDir)
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

	/*app := app.New()
	window = app.NewWindow("nSurge")
	window.Resize(fyne.NewSize(800, -1))

	accountLabel = widget.NewEntry()
	accountLabel.SetText(client.Address())
	accountLabel.SetReadOnly(true)

	testLabel = widget.NewLabel("- idle - ")
	progressBar = widget.NewProgressBar()

	fileBox = widget.NewVBox()
	fileScroller := widget.NewScrollContainer(fileBox)
	fileScroller.SetMinSize(fyne.NewSize(-1, 300))

	contentBox = widget.NewVBox(
		widget.NewLabel("Your address"),
		accountLabel,
		testLabel,
		progressBar,
		widget.NewButton("Fetch Remote Files", func() {
			topicEncoded := surgeTopicEncode(testTopic)
			getSubscriptions(topicEncoded)
		}),
		widget.NewLabel("- Remote Files -"),
		fileScroller,
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	)

	window.SetContent(contentBox)
	*/
	go ScanLocal()
	topicEncoded := TopicEncode(TestTopic)
	GetSubscriptions(topicEncoded)

	go updateGUI()

	//window.ShowAndRun()
}

func updateGUI() {
	/*for true {
		time.Sleep(time.Millisecond * 16)

		if chunksTotal > 0 && chunksReceived < chunksTotal {
			remainingChunks := chunksTotal - chunksReceived
			activeWorkers := chunksRequested - chunksReceived
			testLabel.SetText("Remaining chunks: " + strconv.Itoa(remainingChunks) + " Active Workers: " + strconv.Itoa(activeWorkers))
			progressBar.SetValue(float64(chunksReceived) / float64(chunksTotal))
		} else {
			testLabel.SetText("- idle -")
			progressBar.SetValue(0)
		}
	}*/
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
	log.Println("TESTING TESTING TESTING");
	return nil
}
//DownloadFile downloads the file
func DownloadFile(Addr string, Size int64, FileID string) {

	// Create a sessions
	var err error

	sessionConfing := nkn.GetDefaultSessionConfig()
	sessionConfing.MTU = 16384
	dialConfig := &nkn.DialConfig{
		SessionConfig: sessionConfing,
	}

	downloadSession, err := client.DialWithConfig(Addr, dialConfig)
	if err != nil {
		log.Fatal(err)
	}
	downloadReader := bufio.NewReader(downloadSession)

	surgeSession := Session{
		reader:  downloadReader,
		session: downloadSession,
	}

	pushNotification("Download Started", "downloading file: " + "'"+FileID+"'")

	go initiateSession(surgeSession)

	// If the file doesn't exist allocate it
	var path = remotePath + "/" + FileID
	AllocateFile(path, Size)

	chunksRequested = 0
	chunksReceived = 0
	//Try send request to self
	var numChunks = uint32((Size-1)/int64(ChunkSize)) + 1
	chunksTotal = int(numChunks)

	downloadJob := func() {
		for i := uint32(0); i < numChunks; i++ {
			workerCount++
			chunksRequested++
			go RequestChunk(surgeSession, FileID, int32(i))

			for workerCount >= 256 {
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
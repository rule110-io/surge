package surge

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"log"

	bitmap "github.com/boljen/go-bitmap"
	nkn "github.com/nknorg/nkn-sdk-go"
	"github.com/rule110-io/surge/backend/constants"
	"github.com/rule110-io/surge/backend/messaging"
	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/mutexes"
	pb "github.com/rule110-io/surge/backend/payloads"
	"github.com/rule110-io/surge/backend/platform"
	"github.com/rule110-io/surge/backend/sessionmanager"
	"google.golang.org/protobuf/proto"
)

//NumClientsStruct struct to hold number of online clients
type NumClientsStruct struct {
	Online int
}

//FrontendReady flags whether frontend is ready to receive events etc
var FrontendReady = false

var workerMap map[string]int

//ListedFiles are remote files that can be downloaded
var ListedFiles []models.File

var wailsContext *context.Context

//whether the nkn client is initialized
var clientInitialized = false

//The nkn client
var client *nkn.MultiClient

//NumClientsStruct .

//var numClientsStore *wails.Store

// WailsBind is a binding function at startup
func WailsBind(ctx *context.Context) {

	wailsContext = ctx

	platform.SetWailsContext(ctx, SetVisualMode)

	//Mac specific functions
	//go platform.InitOSHandler()
	platform.SetVisualModeLikeOS()

	// numClients := NumClientsStruct{
	// 	Online: 0,
	// }

	//numClientsStore = runtime.Store.New("numClients", numClients)

	updateNumClientStore()

	//Wait for our client to initialize, perhaps there is no internet connectivity
	tryCount := 1
	for !clientInitialized {
		time.Sleep(time.Second)
		if tryCount%10 == 0 {
			pushError("Connection to NKN not yet established", "do you have an active internet connection?")
		}
		tryCount++
	}
	updateNumClientStore()

	//Startup async processes to continue processing subs/files and updating gui
	go updateFileDataWorker()

	FrontendReady = true
	log.Println("Frontend connected")
}

//InitializeClient Initiates the surge client and instantiates connection with the NKN network
func InitializeClient(args []string) bool {
	var err error

	account := InitializeAccount()
	client, err = nkn.NewMultiClient(account, "", getNumberClients(), false, &nkn.ClientConfig{
		ConnectRetries:    1000,
		SeedRPCServerAddr: GetBootstrapRPC(),
	})
	if err != nil {
		pushError(err.Error(), "do you have an active internet connection?")
	}

	<-client.OnConnect.C
	clientInitialized = true
	sessionmanager.Initialize(client, onClientConnected, onClientDisconnected)

	pushNotification("Client Connected", "Successfully connected to the NKN network")

	client.Listen(nil)
	go Listen()

	//get all files in the DB
	dbFiles := dbGetAllFiles()
	var filesOnDisk []models.File

	//for each file in DB
	for i := 0; i < len(dbFiles); i++ {
		//if local path of file is still valid
		if FileExists(dbFiles[i].Path) {
			filesOnDisk = append(filesOnDisk, dbFiles[i])
		} else {
			dbFiles[i].IsMissing = true
			dbFiles[i].IsDownloading = false
			dbFiles[i].IsUploading = false
			dbInsertFile(dbFiles[i])
		}
	}

	for i := 0; i < len(filesOnDisk); i++ {
		if filesOnDisk[i].IsDownloading {
			go restartDownload(filesOnDisk[i].FileHash)
		}
	}

	messaging.Initialize(client, client.Account(), MessageReceived)

	//Get the transaction fee setting
	TransactionFee, err = DbReadSetting("defaultTxFee")
	if err != nil {
		DbWriteSetting("defaultTxFee", "0")
		TransactionFee = "0"
	}

	go autoSubscribeWorker()

	go platform.WatchOSXHandler()

	//Insert new file from arguments and start download
	if len(args) > 0 && len(args[0]) > 0 {
		platform.AskUser("startDownloadMagnetLinks", "{files : ["+args[0]+"]}")
	}

	return true
}

//StartClient Starts the surge client
func StartClient(args []string) {

	//Initialize all our global data maps
	workerMap = make(map[string]int)
	downloadBandwidthAccumulator = make(map[string]int)
	uploadBandwidthAccumulator = make(map[string]int)
	zeroBandwidthMap = make(map[string]bool)
	fileBandwidthMap = make(map[string]models.BandwidthMA)
	chunksInTransit = make(map[string]bool)

	//Initialize our surge nkn client
	InitializeFileSeedTracker()
	InitializeTopicsManager()
	InitializeClient(args)

	//If we have no subs, subscribe to official
	if len(topicsMap) == 0 {
		subscribeToSurgeTopic(constants.SurgeOfficialTopic, true)
	}
}

//StopClient Stops the surge client and cleans up
func StopClient() {

	//Persist our connections for future bootstraps
	PersistRPC(client)

	client.Close()
}

//DownloadFileByHash Downloads a file by providing a hash
func DownloadFileByHash(Hash string) bool {

	//Addr string, Size int64, FileID string
	file := getListedFileByHash(Hash)
	if file == nil {
		pushError("Error on download file", "No listed file with hash: "+Hash)
		return false
	}

	pushNotification("Download Started", file.FileName)

	remoteFolder, err := GetDownloadFolderPath()
	if err != nil {
		pushError("Error on download file", "Could not access download folder at path: "+remoteFolder)
	}

	// If the file doesn't exist allocate it
	var path = remoteFolder + string(os.PathSeparator) + file.FileName

	isAllocated := AllocateFile(path, file.FileSize)
	if !isAllocated {
		return false
	}

	numChunks := int((file.FileSize-1)/int64(constants.ChunkSize)) + 1

	//When downloading from remote enter file into db
	_, err = dbGetFile(Hash)
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
	//rand.Shuffle(len(randomChunks), func(i, j int) { randomChunks[i], randomChunks[j] = randomChunks[j], randomChunks[i] })

	downloadChunks(file, randomChunks)

	return true
}

// Restarts a file download by providing a hash
func restartDownload(Hash string) {
	file, err := dbGetFile(Hash)
	if err != nil {
		pushError("Error on restart download", err.Error())
		return
	}

	//Get missing chunk indices
	var missingChunks []int
	for i := 0; i < file.NumChunks; i++ {
		if !bitmap.Get(file.ChunkMap, i) {
			missingChunks = append(missingChunks, i)
		}
	}

	numChunks := len(missingChunks)

	//Nothing more to download
	if numChunks == 0 {
		return
	}

	rand.Seed(time.Now().UnixNano())
	//rand.Shuffle(numChunks, func(i, j int) { missingChunks[i], missingChunks[j] = missingChunks[j], missingChunks[i] })

	log.Println("Restarting Download for", file.FileName)

	downloadChunks(file, missingChunks)
}

// fetches the number of clients connected and stores it
func updateNumClientStore() {
	numConnections := sessionmanager.GetSessionLength()
	if clientInitialized {
		numConnections++
	}
}

// listens for incoming sessions
func listenForIncomingSessions() {

	for !client.IsClosed() {
		listenSession, err := client.Accept()
		if err != nil {
			pushError("Error on client accept", err.Error())
			continue
		}

		sessionmanager.AcceptSession(listenSession)
	}
}

// Listen will listen to incoming requests for chunks
func Listen() {
	go listenForIncomingSessions()
}

func onClientConnected(session *sessionmanager.Session, isDialIn bool) {
	go updateNumClientStore()
	addr := session.Session.RemoteAddr().String()

	log.Println("Client Connected", addr)

	go listenToSession(session)
}

func onClientDisconnected(addr string) {
	go updateNumClientStore()

	mutexes.ListedFilesLock.Lock()
	defer mutexes.ListedFilesLock.Unlock()

	RemoveSeeder(addr)
	log.Println("Client Disconnected", addr)

	//Remove empty seeders listings
	for i := 0; i < len(ListedFiles); i++ {
		if !AnySeeders(ListedFiles[i].FileHash) {
			// Remove the element at index i from a.
			ListedFiles[i] = ListedFiles[len(ListedFiles)-1] // Copy last element to index i.
			ListedFiles[len(ListedFiles)-1] = models.File{}  // Erase last element (write zero value).
			ListedFiles = ListedFiles[:len(ListedFiles)-1]   // Truncate slice.
			i--
		}
	}
}

func listenToSession(Session *sessionmanager.Session) {
	defer RecoverAndLog()

	addr := Session.Session.RemoteAddr().String()

	log.Println("Initiate Session", addr)

	for Session.Session != nil {
		data, chunkType, err := SessionRead(Session)
		if err != nil {
			log.Println("Session read failed, closing session error:", err)
			break
		}

		sessionmanager.UpdateActivity(Session.Session.RemoteAddr().String())

		switch chunkType {
		case constants.SurgeChunkID:
			//Write add to download internally after parsing data
			processChunk(Session, data)
		}
	}
}

func processChunk(Session *sessionmanager.Session, Data []byte) {

	//Try to parse SurgeMessage
	surgeMessage := &pb.SurgeMessage{}
	if err := proto.Unmarshal(Data, surgeMessage); err != nil {
		log.Panic("Failed to parse surge message:", err)
	}

	//Write add to download
	mutexes.BandwidthAccumulatorMapLock.Lock()
	downloadBandwidthAccumulator[surgeMessage.FileID] += len(Data)
	mutexes.BandwidthAccumulatorMapLock.Unlock()

	//Data nill means its a request for data
	if surgeMessage.Data == nil {
		go TransmitChunk(Session, surgeMessage.FileID, surgeMessage.ChunkID)
	} else { //If data is not nill we are receiving data

		//When we receive a chunk mark it as no longer in transit
		chunkKey := surgeMessage.FileID + "_" + strconv.Itoa(int(surgeMessage.ChunkID))

		mutexes.ChunkInTransitLock.Lock()
		chunksInTransit[chunkKey] = false
		mutexes.ChunkInTransitLock.Unlock()

		mutexes.WorkerMapLock.Lock()
		workerMap[Session.Session.RemoteAddr().String()]--
		if workerMap[Session.Session.RemoteAddr().String()] < 0 {
			workerMap[Session.Session.RemoteAddr().String()] = 0
		}
		mutexes.WorkerMapLock.Unlock()

		go WriteChunk(surgeMessage.FileID, surgeMessage.ChunkID, surgeMessage.Data)
	}
}

//SeedFilepath generates everything needed to seed a file
func SeedFilepath(Path string, Topic string) bool {

	permissions := GetTopicPermissions(Topic, GetAccountAddress())

	if !permissions.CanWrite {
		pushError("Seed File Error", "no write permission for this topic.")
	}

	log.Println("Seeding file", Path)

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Panic(err)
	}
	randomHash := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	fileName := filepath.Base(Path)
	fileSize := surgeGetFileSize(Path)
	numChunks := int((fileSize-1)/int64(constants.ChunkSize)) + 1
	chunkMap := bitmap.NewSlice(numChunks)

	//Local files are always fully available, set all chunks to 1
	for i := 0; i < numChunks; i++ {
		bitmap.Set(chunkMap, i, true)
	}

	//Append to local files
	localFile := models.File{
		FileName:      fileName,
		FileSize:      fileSize,
		FileHash:      randomHash,
		Path:          Path,
		NumChunks:     numChunks,
		ChunkMap:      chunkMap,
		IsUploading:   false,
		IsDownloading: false,
		IsHashing:     true,
		Topic:         Topic,
	}

	//Check if file is already seeded
	_, err = dbGetFile(localFile.FileHash)
	if err == nil {
		//File already seeding
		pushError("Seed File Error", fileName+" already seeding.")
		return false
	}

	//When seeding a new file enter file into db
	dbInsertFile(localFile)

	go hashFile(randomHash)

	return true
}

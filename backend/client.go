package surge

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"log"

	bitmap "github.com/boljen/go-bitmap"
	nkn "github.com/nknorg/nkn-sdk-go"
	"github.com/rule110-io/surge/backend/constants"
	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/mutexes"
	pb "github.com/rule110-io/surge/backend/payloads"
	"github.com/rule110-io/surge/backend/platform"
	"github.com/rule110-io/surge/backend/sessionmanager"
	"github.com/wailsapp/wails"
	"google.golang.org/protobuf/proto"
)

type NumClientsStruct struct {
	Online int
}

var FrontendReady = false
var workerCount = 0

//ListedFiles are remote files that can be downloaded
var ListedFiles []models.File

var wailsRuntime *wails.Runtime

//whether the nkn client is initialized
var clientInitialized = false

//The nkn client
var client *nkn.MultiClient

var queryPayload = ""

var clientOnlineMap map[string]bool

//NumClientsStruct .

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
	go updateFileDataWorker()
	go rescanPeers()

	FrontendReady = true
}

// Initiates the surge client and instantiates connection with the NKN network
func InitializeClient(args []string) bool {
	var err error

	account := InitializeAccount()
	client, err = nkn.NewMultiClient(account, "", constants.NumClients, false, &nkn.ClientConfig{
		ConnectRetries: 1000,
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

	dbFiles := dbGetAllFiles()
	var filesOnDisk []models.File

	for i := 0; i < len(dbFiles); i++ {
		if FileExists(dbFiles[i].Path) {
			filesOnDisk = append(filesOnDisk, dbFiles[i])
		} else {
			dbFiles[i].IsMissing = true
			dbFiles[i].IsDownloading = false
			dbFiles[i].IsUploading = false
			dbInsertFile(dbFiles[i])
		}
	}

	go BuildSeedString(filesOnDisk)
	for i := 0; i < len(filesOnDisk); i++ {
		if filesOnDisk[i].IsDownloading {
			go restartDownload(filesOnDisk[i].FileHash)
		}
	}

	go autoSubscribeWorker()

	go platform.WatchOSXHandler()

	//Insert new file from arguments and start download
	if args != nil && len(args) > 0 && len(args[0]) > 0 {
		platform.AskUser("startDownloadMagnetLinks", "{files : ["+args[0]+"]}")
	}

	return true
}

// Starts the surge client
func StartClient(args []string) {

	//Initialize all our global data maps
	clientOnlineMap = make(map[string]bool)
	downloadBandwidthAccumulator = make(map[string]int)
	uploadBandwidthAccumulator = make(map[string]int)
	zeroBandwidthMap = make(map[string]bool)
	fileBandwidthMap = make(map[string]models.BandwidthMA)
	chunksInTransit = make(map[string]bool)

	//Initialize our surge nkn client
	go InitializeClient(args)
}

// Stops the surge client and cleans up
func StopClient() {
	client.Close()
	client = nil
}

// Downloads a file by providing a hash
func DownloadFileByHash(Hash string) bool {

	//Addr string, Size int64, FileID string
	file := getListedFileByHash(Hash)
	if file == nil {
		pushError("Error on download file", "No listed file with hash: "+Hash)
	}

	pushNotification("Download Started", file.FileName)

	remoteFolder, err := platform.GetRemoteFolder()
	if err != nil {
		log.Println("Remote folder does not exist")
	}

	// If the file doesn't exist allocate it
	var path = remoteFolder + string(os.PathSeparator) + file.FileName
	AllocateFile(path, file.FileSize)
	numChunks := int((file.FileSize-1)/int64(constants.ChunkSize)) + 1

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
		if bitmap.Get(file.ChunkMap, i) == false {
			missingChunks = append(missingChunks, i)
		}
	}

	numChunks := len(missingChunks)

	//Nothing more to download
	if numChunks == 0 {
		platform.ShowNotification("Download Finished", "Download for "+file.FileName+" finished!")
		pushNotification("Download Finished", file.FileName)
		file.IsDownloading = false
		file.IsUploading = true
		dbInsertFile(*file)
		go AddToSeedString(*file)
		return
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(numChunks, func(i, j int) { missingChunks[i], missingChunks[j] = missingChunks[j], missingChunks[i] })

	log.Println("Restarting Download Creation Session for", file.FileName)

	downloadChunks(file, missingChunks)
}

// fetches the number of clients connected and stores it
func updateNumClientStore() {
	numConnections := sessionmanager.GetSessionLength()
	if clientInitialized {
		numConnections++
	}
	numClientsStore.Update(func(data NumClientsStruct) NumClientsStruct {
		return NumClientsStruct{
			Online: numConnections,
		}
	})
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

	fmt.Println(string("\033[36m"), "Client Connected", addr, string("\033[0m"))

	go listenToSession(session)

	if isDialIn {
		fmt.Println(string("\033[36m"), "Sending file query to accepted client", addr, string("\033[0m"))
		go SendQueryRequest(addr, "Testing query functionality.")
	}
}

func onClientDisconnected(addr string) {
	go updateNumClientStore()

	//Remove this address from remote file seeders
	mutexes.ListedFilesLock.Lock()
	for i := 0; i < len(ListedFiles); i++ {
		ListedFiles[i].Seeders = removeStringFromSlice(ListedFiles[i].Seeders, addr)
		ListedFiles[i].SeederCount = len(ListedFiles[i].Seeders)
		fmt.Println(string("\033[31m"), "onClientDisconnected", ListedFiles[i].FileName, "seeders remaining:", ListedFiles[i].SeederCount, string("\033[0m"))
	}

	//Remove empty seeders listings
	for i := 0; i < len(ListedFiles); i++ {
		if len(ListedFiles[i].Seeders) == 0 {
			// Remove the element at index i from a.
			ListedFiles[i] = ListedFiles[len(ListedFiles)-1] // Copy last element to index i.
			ListedFiles[len(ListedFiles)-1] = models.File{}  // Erase last element (write zero value).
			ListedFiles = ListedFiles[:len(ListedFiles)-1]   // Truncate slice.
			i--
		}
	}

	mutexes.ListedFilesLock.Unlock()
}

func listenToSession(Session *sessionmanager.Session) {
	defer RecoverAndLog()

	addr := Session.Session.RemoteAddr().String()

	fmt.Println(string("\033[31m"), "Initiate Session", addr, string("\033[0m"))

	for Session.Session != nil {
		data, chunkType, err := SessionRead(Session)
		fmt.Println(string("\033[31m"), "Read data from session", addr, string("\033[0m"))

		if err != nil {
			log.Println("Session read failed, closing session error:", err)
			break
		}

		sessionmanager.UpdateActivity(Session.Session.RemoteAddr().String())

		switch chunkType {
		case constants.SurgeChunkID:
			//Write add to download internally after parsing data
			processChunk(Session, data)
			break
		case constants.SurgeQueryRequestID:
			processQueryRequest(Session, data)
			//Write add to download
			mutexes.BandwidthAccumulatorMapLock.Lock()
			downloadBandwidthAccumulator["DISCOVERY"] += len(data)
			mutexes.BandwidthAccumulatorMapLock.Unlock()
			break
		case constants.SurgeQueryResponseID:
			processQueryResponse(Session, data)
			//Write add to download
			mutexes.BandwidthAccumulatorMapLock.Lock()
			downloadBandwidthAccumulator["DISCOVERY"] += len(data)
			mutexes.BandwidthAccumulatorMapLock.Unlock()
			break
		}
	}
}

func downloadChunks(file *models.File, randomChunks []int) {
	fileID := file.FileHash
	file = getListedFileByHash(fileID)

	for file == nil {
		time.Sleep(time.Second)
		file = getListedFileByHash(fileID)
	}

	numChunks := len(randomChunks)

	seederAlternator := 0
	mutateSeederLock := sync.Mutex{}
	appendChunkLock := sync.Mutex{}

	//Give the seeder a fair start with timers when a download is initiated
	//Potentionally this seeder was last queried 60 seconds ago for files and otherwise idle but online
	for _, seeder := range file.Seeders {
		sessionmanager.UpdateActivity(seeder)
	}

	downloadJob := func(terminateFlag *bool) {

		//Used to terminate the rescanning of peers
		terminate := func(flag *bool) {
			*flag = true
		}
		defer terminate(terminateFlag)

		for i := 0; i < numChunks; i++ {
			fmt.Println(string("\033[36m"), "Preparing Chunk Fetch", string("\033[0m"))
			file = getListedFileByHash(fileID)

			for file == nil || len(file.Seeders) == 0 {
				time.Sleep(time.Second * 5)
				fmt.Println(string("\033[36m"), "SLEEPING NO SEEDERS FOR FILE", string("\033[0m"))
				file = getListedFileByHash(fileID)
			}

			dbFile, err := dbGetFile(file.FileHash)

			//Check if file is still tracked in surge
			if err != nil {
				log.Println("Download Job Terminated", "File no longer in DB")
				return
			}

			//Pause if file is paused
			for err == nil && dbFile.IsPaused {
				time.Sleep(time.Second * 5)
				dbFile, err = dbGetFile(file.FileHash)
				if err != nil {
					log.Println("Download Job Terminated", "File no longer in DB")
					return
				}

				//Coming out of a pause situation we reset our received timer
				if !dbFile.IsPaused {
					//Give the seeder a fair start with timers when a download is initiated
					//Potentionally this seeder was last queried 60 seconds ago for files and otherwise idle but online
					for _, seeder := range file.Seeders {
						sessionmanager.UpdateActivity(seeder)
					}
				}
			}

			for workerCount >= constants.NumWorkers {
				time.Sleep(time.Millisecond)
			}
			workerCount++

			//Create a async job to download a chunk
			requestChunkJob := func(chunkID int) {

				success := false
				downloadSeederAddr := ""

				mutateSeederLock.Lock()
				if len(file.Seeders) > seederAlternator {
					//Get seeder
					downloadSeederAddr = file.Seeders[seederAlternator]
					session, existing := sessionmanager.GetExistingSessionWithoutClosing(downloadSeederAddr, constants.WorkerGetSessionTimeout)

					if existing {
						success = RequestChunk(session, file.FileHash, int32(chunkID))
					} else {
						success = false
					}
				}
				mutateSeederLock.Unlock()

				//if download fails append the chunk to remaining to retry later
				if !success {
					appendChunkLock.Lock()
					randomChunks = append(randomChunks, chunkID)
					numChunks++
					appendChunkLock.Unlock()

					workerCount--
					//TODO: Remove this clamp, dont double count timeouted arrivals
					if workerCount < 0 {
						workerCount = 0
					}

					//This file was not available at this time from this seeder, drop seeder for file.
					mutateSeederLock.Lock()
					for i := 0; i < len(ListedFiles); i++ {
						if ListedFiles[i].FileHash == fileID {
							ListedFiles[i].Seeders = removeStringFromSlice(ListedFiles[i].Seeders, downloadSeederAddr)
							ListedFiles[i].SeederCount = len(ListedFiles[i].Seeders)
							file = &ListedFiles[i]
							break
						}
					}
					mutateSeederLock.Unlock()

					//return out of job
					return
				}

				//If chunk is requested add to transit map
				chunkKey := file.FileHash + "_" + strconv.Itoa(chunkID)

				mutexes.ChunkInTransitLock.Lock()
				chunksInTransit[chunkKey] = true
				mutexes.ChunkInTransitLock.Unlock()

				//Sleep and check if entry still exists in transit map.
				sleepWorker := true
				inTransit := true
				receiveTimeoutCounter := 0

				for sleepWorker {
					time.Sleep(time.Second)
					//fmt.Println(string("\033[36m"), "Worker Sleeping", string("\033[0m"))

					//Check if connection is lost
					_, sessionExists := sessionmanager.GetExistingSessionWithoutClosing(downloadSeederAddr, constants.WorkerGetSessionTimeout)
					if !sessionExists {
						//if session no longer exists
						fmt.Println(string("\033[36m"), "session no longer exists while waiting for chunk to arrive for", downloadSeederAddr, string("\033[0m"))

						inTransit = true
						sleepWorker = false
						break
					}

					//Check if received
					mutexes.ChunkInTransitLock.Lock()
					isInTransit := chunksInTransit[chunkKey]
					mutexes.ChunkInTransitLock.Unlock()

					if !isInTransit {
						//if no longer in transit, continue workers
						inTransit = false
						sleepWorker = false
						break
					}
					if receiveTimeoutCounter >= constants.WorkerChunkReceiveTimeout {
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
					//TODO: Remove this clamp, dont double count timeouted arrivals
					if workerCount < 0 {
						workerCount = 0
					}

					//This file was not available at this time from this seeder, drop seeder for file.
					mutateSeederLock.Lock()
					for i := 0; i < len(ListedFiles); i++ {
						if ListedFiles[i].FileHash == fileID {
							ListedFiles[i].Seeders = removeStringFromSlice(ListedFiles[i].Seeders, downloadSeederAddr)
							ListedFiles[i].SeederCount = len(ListedFiles[i].Seeders)
							file = &ListedFiles[i]
							break
						}
					}
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

			seederAlternator++
			if seederAlternator > len(file.Seeders)-1 {
				seederAlternator = 0
			}
		}
	}

	terminateFlag := false
	go downloadJob(&terminateFlag)
}

// SendQueryRequest sends a query to a client on session
func SendQueryRequest(Addr string, Query string) bool {

	surgeSession, exists := sessionmanager.GetExistingSession(Addr, constants.SendQueryRequestSessionTimeout)

	if !exists {
		return false
	}

	msg := &pb.SurgeQuery{
		Query: Query,
	}
	msgSerialized, err := proto.Marshal(msg)
	if err != nil {
		log.Panic("Failed to encode surge message:", err)
		return false
	}

	fmt.Println(string("\033[31m"), "Send Query Request", Addr, string("\033[0m"))
	written, err := SessionWrite(surgeSession, msgSerialized, constants.SurgeQueryRequestID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
	if err != nil {
		log.Println("Failed to send Surge Request:", err)
		return false
	}

	//Write add to upload
	mutexes.BandwidthAccumulatorMapLock.Lock()
	uploadBandwidthAccumulator["DISCOVERY"] += written
	mutexes.BandwidthAccumulatorMapLock.Unlock()

	return true
}

// SendQueryResponse sends a query to a client on session
func SendQueryResponse(Session *sessionmanager.Session, Query string) {

	b := []byte(queryPayload)
	fmt.Println(string("\033[31m"), "Send Query Response", Session.Session.RemoteAddr().String(), string("\033[0m"))
	written, err := SessionWrite(Session, b, constants.SurgeQueryResponseID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
	if err != nil {
		log.Println("Failed to send Surge Ruquest:", err)
	}
	//Write add to upload
	mutexes.BandwidthAccumulatorMapLock.Lock()
	uploadBandwidthAccumulator["DISCOVERY"] += written
	mutexes.BandwidthAccumulatorMapLock.Unlock()
}

func processQueryRequest(Session *sessionmanager.Session, Data []byte) {

	//Try to parse SurgeMessage
	surgeQuery := &pb.SurgeQuery{}
	if err := proto.Unmarshal(Data, surgeQuery); err != nil {
		log.Panic("Failed to parse surge message:", err)
	}
	log.Println("Query received", surgeQuery.Query)

	SendQueryResponse(Session, surgeQuery.Query)
}

func processQueryResponse(Session *sessionmanager.Session, Data []byte) {

	//Try to parse SurgeMessage
	s := string(Data)
	seeder := Session.Session.RemoteAddr().String()

	fmt.Println(string("\033[36m"), "file query response received", seeder, string("\033[0m"))

	mutexes.ListedFilesLock.Lock()

	//Parse the response
	payloadSplit := strings.Split(s, "surge://")
	for j := 0; j < len(payloadSplit); j++ {
		data := strings.Split(payloadSplit[j], "|")

		if len(data) < 3 {
			continue
		}

		fileSize, _ := strconv.ParseInt(data[3], 10, 64)
		numChunks := int((fileSize-1)/int64(constants.ChunkSize)) + 1

		newListing := models.File{
			FileLocation: "remote",
			FileName:     data[2],
			FileSize:     fileSize,
			FileHash:     data[4],
			Seeders:      []string{seeder},
			Path:         "",
			NumChunks:    numChunks,
			ChunkMap:     nil,
			SeederCount:  1,
		}

		//Replace existing, or remove.
		var replace = false
		for l := 0; l < len(ListedFiles); l++ {
			if ListedFiles[l].FileHash == newListing.FileHash {

				//if the seeder is unique add it as an additional seeder for the file
				ListedFiles[l].Seeders = append(ListedFiles[l].Seeders, seeder)
				ListedFiles[l].Seeders = distinctStringSlice(ListedFiles[l].Seeders)
				ListedFiles[l].SeederCount = len(ListedFiles[l].Seeders)

				replace = true
				break
			}
		}
		//Unique listing so we add
		if replace == false {
			ListedFiles = append(ListedFiles, newListing)
		}

		fmt.Println(string("\033[33m"), "Filename", newListing.FileName, "FileHash", newListing.FileHash, string("\033[0m"))

		log.Println("Query response new file: ", newListing.FileName, " seeder: ", seeder)
	}
	mutexes.ListedFilesLock.Unlock()
}

// RequestChunk sends a request to an address for a specific chunk of a specific file
func RequestChunk(Session *sessionmanager.Session, FileID string, ChunkID int32) bool {
	if Session == nil || Session.Session == nil {
		return false
	}

	msg := &pb.SurgeMessage{
		FileID:  FileID,
		ChunkID: ChunkID,
	}
	msgSerialized, err := proto.Marshal(msg)
	if err != nil {
		log.Panic("Failed to encode surge message:", err)
	} else {
		fmt.Println(string("\033[31m"), "Request Chunk", FileID, ChunkID, string("\033[0m"))

		written, err := SessionWrite(Session, msgSerialized, constants.SurgeChunkID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
		if err != nil {
			log.Println("Failed to request chunk", err)
			return false
		}

		//Write add to upload
		mutexes.BandwidthAccumulatorMapLock.Lock()
		uploadBandwidthAccumulator[FileID] += written
		mutexes.BandwidthAccumulatorMapLock.Unlock()
	}

	return true
}

// TransmitChunk tranmits target file chunk to address
func TransmitChunk(Session *sessionmanager.Session, FileID string, ChunkID int32) {
	defer RecoverAndLog()

	//Open file

	mutexes.FileWriteLock.Lock()
	fileInfo, err := dbGetFile(FileID)
	if err != nil {
		log.Println("Error on transmit chunk - file not in db", err.Error())
		return
	}
	fileInfo.ChunksShared++
	dbInsertFile(*fileInfo)
	mutexes.FileWriteLock.Unlock()

	file, err := os.Open(fileInfo.Path)

	//When we have an OS read error on the file mark the file as missing, stop down and uploads on it
	if err != nil {
		log.Println("Error on transmit chunk - file read failure", err.Error())

		mutexes.FileWriteLock.Lock()
		fileInfo.IsMissing = true
		fileInfo.IsDownloading = false
		fileInfo.IsUploading = false
		dbInsertFile(*fileInfo)
		mutexes.FileWriteLock.Unlock()

		return
	}
	defer file.Close()

	//Read the requested chunk
	chunkOffset := int64(ChunkID) * constants.ChunkSize
	buffer := make([]byte, constants.ChunkSize)
	bytesread, err := file.ReadAt(buffer, chunkOffset)

	if err != nil {
		if err != io.EOF {
			log.Println("Error on transmit chunk - read chunk failed: ", ChunkID, err.Error())
			return
		}
	}

	//Create the proto data
	dataReply := &pb.SurgeMessage{
		FileID:  FileID,
		ChunkID: ChunkID,
		Data:    buffer[:bytesread],
	}
	dateReplySerialized, err := proto.Marshal(dataReply)
	if err != nil {
		log.Panic("Error on transmit chunk - chunk serialization error", err.Error())
		return
	}

	//Transmit the chunk
	fmt.Println(string("\033[31m"), "Transmit Chunk", FileID, ChunkID, string("\033[0m"))
	written, err := SessionWrite(Session, dateReplySerialized, constants.SurgeChunkID) //Client.Send(nkn.NewStringArray(Addr), dateReplySerialized, nil)
	if err != nil {
		log.Println("Error on transmit chunk - failed to write to session", err.Error())
		return
	}
	log.Println("Chunk transmitted: ", bytesread, " bytes")

	//Write add to upload
	mutexes.BandwidthAccumulatorMapLock.Lock()
	uploadBandwidthAccumulator[FileID] += written
	mutexes.BandwidthAccumulatorMapLock.Unlock()
}

func processChunk(Session *sessionmanager.Session, Data []byte) {

	//Try to parse SurgeMessage
	surgeMessage := &pb.SurgeMessage{}
	if err := proto.Unmarshal(Data, surgeMessage); err != nil {
		log.Panic("Failed to parse surge message:", err)
	}
	fmt.Println(string("\033[31m"), "PROCESSING CHUNK", string("\033[0m"))

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

		go WriteChunk(surgeMessage.FileID, surgeMessage.ChunkID, surgeMessage.Data)
	}
}

// SessionWrite writes to session
func SessionWrite(Session *sessionmanager.Session, Data []byte, ID byte) (written int, err error) {

	if Session == nil || Session.Session == nil {
		return 0, errors.New("write to session error, session nil")
	}
	//Package identifier to know what we are sending
	packID := make([]byte, 1)
	packID[0] = ID

	//Create buffer of 4 bytes to put the size of the package
	buff := make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, uint32(len(Data)))

	//append pack and buff
	buff = append(packID, buff...)

	//Write data
	buff = append(buff, Data...)

	//Session.session.SetWriteDeadline(time.Now().Add(60 * time.Second))
	_, err = Session.Session.Write(buff)
	if err != nil {
		return 0, err
	}

	return len(buff), err
}

//SessionRead reads from session
func SessionRead(Session *sessionmanager.Session) (data []byte, ID byte, err error) {

	headerBuffer := make([]byte, 5) //int32 size of header + 1 for packid

	// the header of 4 bytes + 1 for packid
	_, err = io.ReadFull(Session.Reader, headerBuffer)
	if err != nil {
		if err.Error() == "session closed" {
			log.Println(err)
			return nil, 0x0, err
		}
		log.Println(err)
		return nil, 0x0, err
	}

	//Get the packid
	packID := headerBuffer[0]
	log.Println(packID)

	//Get the size from the bytes
	sizeBytes := append(headerBuffer[:0], headerBuffer[1:]...)

	size := binary.LittleEndian.Uint32(sizeBytes)

	data = make([]byte, size)

	// read the full message, or return an error
	_, err = io.ReadFull(Session.Reader, data[:int(size)])
	if err != nil {
		log.Println(err)
		return nil, 0x0, err
	}

	return data, packID, err
}

//SeedFile generates everything needed to seed a file
func SeedFilepath(Path string) bool {

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
		FileLocation:  "local",
		FileName:      fileName,
		FileSize:      fileSize,
		FileHash:      randomHash,
		Path:          Path,
		NumChunks:     numChunks,
		ChunkMap:      chunkMap,
		IsUploading:   false,
		IsDownloading: false,
		IsHashing:     true,
	}

	//Check if file is already seeded
	_, err = dbGetFile(localFile.FileHash)
	if err == nil {
		//File already seeding
		pushError("Seed failed", fileName+" already seeding.")
		return false
	}

	//When seeding a new file enter file into db
	dbInsertFile(localFile)

	go hashFile(randomHash)

	return true
}

//BuildSeedString builds a string of seeded files to share with clients
func BuildSeedString(dbFiles []models.File) {

	newQueryPayload := ""
	for _, dbFile := range dbFiles {
		magnet := surgeGenerateMagnetLink(dbFile.FileName, dbFile.FileSize, dbFile.FileHash, client.Addr().String())
		log.Println("Magnet:", magnet)

		if dbFile.IsUploading {
			//Add to payload
			payload := surgeGenerateTopicPayload(dbFile.FileName, dbFile.FileSize, dbFile.FileHash)
			//log.Println(payload)
			newQueryPayload += payload
		}
	}
	queryPayload = newQueryPayload
}

//AddToSeedString adds to existing seed string
func AddToSeedString(dbFile models.File) {

	//Add to payload
	payload := surgeGenerateTopicPayload(dbFile.FileName, dbFile.FileSize, dbFile.FileHash)
	//log.Println(payload)
	queryPayload += payload

	//Make sure you're subscribed when seeding a file
	go subscribeToSurgeTopic()
}

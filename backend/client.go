package surge

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"log"

	bitmap "github.com/boljen/go-bitmap"
	nkn "github.com/nknorg/nkn-sdk-go"
	"github.com/rule110-io/surge/backend/constants"
	pb "github.com/rule110-io/surge/backend/payloads"
	"github.com/rule110-io/surge/backend/platform"
	"github.com/rule110-io/surge/backend/sessionmanager"
	"google.golang.org/protobuf/proto"
)

//whether the nkn client is initialized
var clientInitialized = false

//The nkn client
var client *nkn.MultiClient

var queryPayload = ""

//InitializeClient initializes connection with nkn
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
	var filesOnDisk []File

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

//Stop cleanup for surge
func Stop() {
	client.Close()
	client = nil
}

//Function that automatically grabs subscriptions for nkn topic
func rescanPeers() {

	for true {
		time.Sleep(constants.RescanPeerInterval)
		GetSubscriptions()
	}
}

func autoSubscribeWorker() {

	//As long as the client is running subscribe
	for true {
		//Only subscribe when this client is hosting anything
		hosting := false

		files := dbGetAllFiles()
		for i := 0; i < len(files); i++ {
			if files[i].IsUploading {
				hosting = true
				break
			}
		}

		if hosting {
			subscribeToSurgeTopic()
		}

		time.Sleep(time.Second * 20 * constants.SubscriptionDuration)
	}
}

func subscribeToSurgeTopic() {
	Topic := TopicEncode(constants.PublicTopic)
	txnHash, err := client.Subscribe("", Topic, constants.SubscriptionDuration, "Surge Beta Client", nil)
	if err != nil {
		log.Println("Probably already subscribed", err)
	} else {
		log.Println("Subscribed: ", txnHash)
	}
}

//GetSubscriptions .
func GetSubscriptions() {

	Topic := TopicEncode(constants.PublicTopic)

	subResponse, err := client.GetSubscribers(Topic, 0, 100, true, true)
	if err != nil {
		pushError(err.Error(), "do you have an active internet connection?")
		return
	}

	for k, v := range subResponse.SubscribersInTxPool.Map {
		subResponse.Subscribers.Map[k] = v
	}

	subscribers = []string{}
	for k, v := range subResponse.Subscribers.Map {
		if len(v) > 0 {
			if k != client.Addr().String() {
				subscribers = append(subscribers, k)
			}
		}
	}

	fmt.Println(string("\033[36m"), "Get Subscriptions", len(subscribers), string("\033[0m"))

	for _, sub := range subscribers {
		connectAndQueryJob := func(addr string) {
			_, err := sessionmanager.GetSession(addr, constants.GetSessionDialTimeout)
			if err == nil {
				fmt.Println(string("\033[36m"), "Sending file query to subscriber", addr, string("\033[0m"))
				go SendQueryRequest(addr, "Testing query functionality.")
			}
		}
		go connectAndQueryJob(sub)
	}
}

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
	ListedFilesLock.Lock()
	for i := 0; i < len(ListedFiles); i++ {
		ListedFiles[i].seeders = removeStringFromSlice(ListedFiles[i].seeders, addr)
		ListedFiles[i].seederCount = len(ListedFiles[i].seeders)
		fmt.Println(string("\033[31m"), "onClientDisconnected", ListedFiles[i].FileName, "seeders remaining:", ListedFiles[i].seederCount, string("\033[0m"))
	}

	//Remove empty seeders listings
	for i := 0; i < len(ListedFiles); i++ {
		if len(ListedFiles[i].seeders) == 0 {
			// Remove the element at index i from a.
			ListedFiles[i] = ListedFiles[len(ListedFiles)-1] // Copy last element to index i.
			ListedFiles[len(ListedFiles)-1] = File{}         // Erase last element (write zero value).
			ListedFiles = ListedFiles[:len(ListedFiles)-1]   // Truncate slice.
			i--
		}
	}

	ListedFilesLock.Unlock()
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
			bandwidthAccumulatorMapLock.Lock()
			downloadBandwidthAccumulator["DISCOVERY"] += len(data)
			bandwidthAccumulatorMapLock.Unlock()
			break
		case constants.SurgeQueryResponseID:
			processQueryResponse(Session, data)
			//Write add to download
			bandwidthAccumulatorMapLock.Lock()
			downloadBandwidthAccumulator["DISCOVERY"] += len(data)
			bandwidthAccumulatorMapLock.Unlock()
			break
		}
	}
}

//DownloadFile downloads the file
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

func downloadChunks(file *File, randomChunks []int) {
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
	for _, seeder := range file.seeders {
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

			for file == nil || len(file.seeders) == 0 {
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
					for _, seeder := range file.seeders {
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
				if len(file.seeders) > seederAlternator {
					//Get seeder
					downloadSeederAddr = file.seeders[seederAlternator]
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
							ListedFiles[i].seeders = removeStringFromSlice(ListedFiles[i].seeders, downloadSeederAddr)
							ListedFiles[i].seederCount = len(ListedFiles[i].seeders)
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
					_, sessionExists := sessionmanager.GetExistingSessionWithoutClosing(downloadSeederAddr, constants.WorkerGetSessionTimeout)
					if !sessionExists {
						//if session no longer exists
						fmt.Println(string("\033[36m"), "session no longer exists while waiting for chunk to arrive for", downloadSeederAddr, string("\033[0m"))

						inTransit = true
						sleepWorker = false
						break
					}

					//Check if received
					chunkInTransitLock.Lock()
					isInTransit := chunksInTransit[chunkKey]
					chunkInTransitLock.Unlock()

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
							ListedFiles[i].seeders = removeStringFromSlice(ListedFiles[i].seeders, downloadSeederAddr)
							ListedFiles[i].seederCount = len(ListedFiles[i].seeders)
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
			if seederAlternator > len(file.seeders)-1 {
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
	bandwidthAccumulatorMapLock.Lock()
	uploadBandwidthAccumulator["DISCOVERY"] += written
	bandwidthAccumulatorMapLock.Unlock()

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
	bandwidthAccumulatorMapLock.Lock()
	uploadBandwidthAccumulator["DISCOVERY"] += written
	bandwidthAccumulatorMapLock.Unlock()
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

	ListedFilesLock.Lock()

	//Parse the response
	payloadSplit := strings.Split(s, "surge://")
	for j := 0; j < len(payloadSplit); j++ {
		data := strings.Split(payloadSplit[j], "|")

		if len(data) < 3 {
			continue
		}

		fileSize, _ := strconv.ParseInt(data[3], 10, 64)
		numChunks := int((fileSize-1)/int64(constants.ChunkSize)) + 1

		newListing := File{
			FileName:    data[2],
			FileSize:    fileSize,
			FileHash:    data[4],
			seeders:     []string{seeder},
			Path:        "",
			NumChunks:   numChunks,
			ChunkMap:    nil,
			seederCount: 1,
		}

		//Replace existing, or remove.
		var replace = false
		for l := 0; l < len(ListedFiles); l++ {
			if ListedFiles[l].FileHash == newListing.FileHash {

				//if the seeder is unique add it as an additional seeder for the file
				ListedFiles[l].seeders = append(ListedFiles[l].seeders, seeder)
				ListedFiles[l].seeders = distinctStringSlice(ListedFiles[l].seeders)
				ListedFiles[l].seederCount = len(ListedFiles[l].seeders)

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
	ListedFilesLock.Unlock()
}

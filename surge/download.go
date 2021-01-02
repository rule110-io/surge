package surge

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	bitmap "github.com/boljen/go-bitmap"
	"github.com/rule110-io/surge-ui/surge/constants"
	"github.com/rule110-io/surge-ui/surge/platform"
	"github.com/rule110-io/surge-ui/surge/sessionmanager"
)

//DownloadFile downloads the file
func DownloadFile(Hash string) bool {

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
	fileName := file.FileName

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
			file = getListedFileByHash(fileID)
			for file == nil || len(file.seeders) == 0 {
				time.Sleep(time.Second * 5)
				fmt.Println(string("\033[36m"), "SLEEPING NO SEEDERS FOR FILE", fileName, string("\033[0m"))
				file = getListedFileByHash(fileID)
			}

			dbFile, err := dbGetFile(file.FileHash)

			//Check if file is still tracked in surge
			if err != nil {
				log.Println("Download Job Terminated", "File no longer in DB")
				pushError("Download Job Terminated", "File no longer in DB")
				return
			}

			//Pause if file is paused
			for err == nil && dbFile.IsPaused {
				time.Sleep(time.Second * 5)
				dbFile, err = dbGetFile(file.FileHash)
				if err != nil {
					log.Println("Download Job Terminated", "File no longer in DB")
					pushError("Download Job Terminated", "File no longer in DB")
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

			for workerCount >= NumWorkers {
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
					session, existing := sessionmanager.GetExistingSessionWithoutClosing(downloadSeederAddr, constants.WorkerGetSessionTimeout, "Get Download Session for Worker - WorkerGetSessionTimeout")

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

					//TODO: Think about alternatives from straight dropping the seeder.
					mutateSeederLock.Lock()
					file.seeders = removeStringFromSlice(file.seeders, downloadSeederAddr)
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
					_, sessionExists := sessionmanager.GetExistingSessionWithoutClosing(downloadSeederAddr, constants.WorkerGetSessionTimeout, "Worker waiting for potential timeout get session - WorkerGetSessionTimeout")
					if !sessionExists {
						//if session no longer exists
						fmt.Println(string("\033[36m"), "session no longer exists", string("\033[0m"))
						fmt.Println(string("\033[36m"), downloadSeederAddr, sessionmanager.GetSessionsString(), string("\033[0m"))

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
					file.seeders = removeStringFromSlice(file.seeders, downloadSeederAddr)
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
			if seederAlternator > len(file.seeders)-1 {
				seederAlternator = 0
			}
			mutateSeederLock.Unlock()
		}
	}

	/*scanForSeeders := func(terminateFlag *bool) {

		//While we are not terminated scan for new peers
		for *terminateFlag == false {
			time.Sleep(time.Second * 5)

			newFile := getListedFileByHash(fileID)
			if newFile != nil {
				//Check for new sessions
				mutateSeederLock.Lock()
				fileSeeders = []string{}
				for i := 0; i < len(newFile.seeders); i++ {
					_, existing := sessionmanager.GetExistingSession(newFile.seeders[i], 60, "Scan for seeders in download session timeout")
					if existing {
						fileSeeders = append(fileSeeders, newFile.seeders[i])
					}
				}
				mutateSeederLock.Unlock()
			}
		}
	}*/

	terminateFlag := false
	go downloadJob(&terminateFlag)
	//go scanForSeeders(&terminateFlag)
}

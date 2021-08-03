package surge

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	bitmap "github.com/boljen/go-bitmap"
	movavg "github.com/mxmCherry/movavg"
	"github.com/rule110-io/surge/backend/constants"
	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/mutexes"
	pb "github.com/rule110-io/surge/backend/payloads"
	"github.com/rule110-io/surge/backend/sessionmanager"
	"google.golang.org/protobuf/proto"
)

var downloadBandwidthAccumulator map[string]int
var uploadBandwidthAccumulator map[string]int

var fileBandwidthMap map[string]models.BandwidthMA

var zeroBandwidthMap map[string]bool

var chunksInTransit map[string]bool

//sets the current bandwith of a file
func fileBandwidth(FileID string) (Download int, Upload int) {

	//Get accumulator
	mutexes.BandwidthAccumulatorMapLock.Lock()
	downAccu := downloadBandwidthAccumulator[FileID]
	downloadBandwidthAccumulator[FileID] = 0

	upAccu := uploadBandwidthAccumulator[FileID]
	uploadBandwidthAccumulator[FileID] = 0
	mutexes.BandwidthAccumulatorMapLock.Unlock()

	if fileBandwidthMap[FileID].Download == nil {
		fileBandwidthMap[FileID] = models.BandwidthMA{
			Download: movavg.ThreadSafe(movavg.NewSMA(10)),
			Upload:   movavg.ThreadSafe(movavg.NewSMA(10)),
		}
	}

	fileBandwidthMap[FileID].Download.Add(float64(downAccu))
	fileBandwidthMap[FileID].Upload.Add(float64(upAccu))

	return int(fileBandwidthMap[FileID].Download.Avg()), int(fileBandwidthMap[FileID].Upload.Avg())
}

func downloadChunks(file *models.File, randomChunks []int, mutateSeederLock *sync.Mutex, activeSeeders *[]string) {
	fileID := file.FileHash

	//todo: lock seeders
	for file == nil || len(*activeSeeders) == 0 {
		time.Sleep(time.Second)
		file = getListedFileByHash(fileID)
	}

	numChunks := len(randomChunks)

	seederAlternator := 0
	appendChunkLock := sync.Mutex{}

	//Give the seeder a fair start with timers when a download is initiated
	//Potentionally this seeder was last queried 60 seconds ago for files and otherwise idle but online
	//todo: Lock seeders
	for _, seeder := range *activeSeeders {
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

			dbFile, err := dbGetFile(fileID)

			//Check if file is still tracked in surge
			if err != nil {
				log.Println("Download Job Terminated", "File no longer in DB")
				return
			}

			//Pause if file is paused
			for err == nil && dbFile.IsPaused {
				time.Sleep(time.Second * 5)
				dbFile, err = dbGetFile(fileID)
				if err != nil {
					log.Println("Download Job Terminated", "File no longer in DB")
					return
				}

				//Coming out of a pause situation we reset our received timer
				if !dbFile.IsPaused {
					//Give the seeder a fair start with timers when a download is initiated
					//Potentionally this seeder was last queried 60 seconds ago for files and otherwise idle but online
					//todo: lockseeders?
					for _, seeder := range *activeSeeders {
						sessionmanager.UpdateActivity(seeder)
					}
				}
			}

			//Create a async job to download a chunk
			requestChunkJob := func(chunkID int, downloadSeederAddr string) {

				success := false

				session, err := sessionmanager.GetSession(downloadSeederAddr, constants.WorkerGetSessionTimeout)

				if err == nil {
					success = RequestChunk(session, fileID, int32(chunkID))
				} else {
					success = false
				}

				//if download fails append the chunk to remaining to retry later
				if !success {
					appendChunkLock.Lock()
					randomChunks = append(randomChunks, chunkID)
					numChunks++
					appendChunkLock.Unlock()

					//TODO: Remove this clamp, dont double count timeouted arrivals
					mutexes.WorkerMapLock.Lock()
					workerMap[downloadSeederAddr]--
					if workerMap[downloadSeederAddr] < 0 {
						workerMap[downloadSeederAddr] = 0
					}
					mutexes.WorkerMapLock.Unlock()

					//return out of job
					return
				}

				//If chunk is requested add to transit map
				chunkKey := fileID + "_" + strconv.Itoa(chunkID)

				mutexes.ChunkInTransitLock.Lock()
				chunksInTransit[chunkKey] = true
				mutexes.ChunkInTransitLock.Unlock()

				//Sleep and check if entry still exists in transit map.
				sleepWorker := true
				inTransit := true
				receiveTimeoutCounter := 0

				for sleepWorker {
					time.Sleep(time.Second)

					//Check if connection is lost
					_, sessionExists := sessionmanager.GetExistingSession(downloadSeederAddr, constants.WorkerGetSessionTimeout)
					if !sessionExists {
						//if session no longer exists
						fmt.Println(string("\033[36m"), "session no longer exists while waiting for chunk to arrive for", downloadSeederAddr, string("\033[0m"))

						fmt.Println(string("\033[36m"), "Attempting to reconnect with a new session", downloadSeederAddr, string("\033[0m"))
						_, err := sessionmanager.GetSession(downloadSeederAddr, 10)
						if err == nil {
							sleepWorker = false
						}
						inTransit = true
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

					//TODO: Remove this clamp, dont double count timeouted arrivals
					mutexes.WorkerMapLock.Lock()
					workerMap[downloadSeederAddr]--
					if workerMap[downloadSeederAddr] < 0 {
						workerMap[downloadSeederAddr] = 0
					}
					mutexes.WorkerMapLock.Unlock()

					//This file was not available at this time from this seeder, drop seeder for file.
					mutateSeederLock.Lock()
					for i := 0; i < len(ListedFiles); i++ {
						if ListedFiles[i].FileHash == fileID {
							removeStringFromSlicePtr(activeSeeders, downloadSeederAddr)
							break
						}
					}
					mutateSeederLock.Unlock()
					//return out of job
					return
				}
			}

			downloadSeederAddr := ""

			//spin to seeder with workers available
			spinForSeeder := true
			for spinForSeeder {

				mutateSeederLock.Lock()
				numSeeds := len(*activeSeeders)
				mutateSeederLock.Unlock()

				for numSeeds == 0 {
					time.Sleep(time.Second)

					listingRefreshed := getListedFileByHash(fileID)

					if listingRefreshed != nil {
						newSeeders := getListedFileByHash(fileID).Seeders
						for _, v := range newSeeders {
							_, err := sessionmanager.GetSession(v, constants.GetSessionDialTimeout)
							if err != nil {
								log.Println(err)
							} else {
								mutateSeederLock.Lock()
								*activeSeeders = append(*activeSeeders, v)
								mutateSeederLock.Unlock()
								sessionmanager.UpdateActivity(v)
							}
						}
						mutateSeederLock.Lock()
						numSeeds = len(*activeSeeders)
						mutateSeederLock.Unlock()
					}
				}

				downloadSeederAddr = (*activeSeeders)[seederAlternator]

				//If seeder selected exceeds worker limit skip
				mutexes.WorkerMapLock.Lock()
				seedWorkerNum := workerMap[downloadSeederAddr]
				mutexes.WorkerMapLock.Unlock()

				if seedWorkerNum >= constants.NumWorkers {
					seederAlternator++
					if seederAlternator > len(*activeSeeders)-1 {
						seederAlternator = 0

						//When weve spun a complete seeder cycle we sleep
						time.Sleep(time.Millisecond)
					}
				} else {
					//If the seeder has room for more workers we accept
					spinForSeeder = false
					mutexes.WorkerMapLock.Lock()
					workerMap[downloadSeederAddr]++
					mutexes.WorkerMapLock.Unlock()
				}
			}

			//get chunk id
			appendChunkLock.Lock()
			chunkid := randomChunks[i]
			appendChunkLock.Unlock()

			go requestChunkJob(chunkid, downloadSeederAddr)
		}
	}

	terminateFlag := false
	go downloadJob(&terminateFlag)
}

func chunksDownloaded(s []byte, num int) int {
	//No chunkmap means no download was initiated, all chunks are local
	if s == nil {
		return num
	}

	chunksLocalNum := 0
	for i := 0; i < num; i++ {
		if bitmap.Get(s, i) {
			chunksLocalNum++
		}
	}
	return chunksLocalNum
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

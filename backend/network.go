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

func downloadChunks(file *models.File, randomChunks []int) {
	log.Println("Starting download for file:", file.FileName, file.FileHash, "size:", file.FileSize)

	fileID := file.FileHash

	//todo: lock seeders
	for file == nil || !AnySeeders(fileID) {
		time.Sleep(time.Second)

		//TODO: this might not be needed anymore now we have the file seeder tracker.
		file = getListedFileByHash(fileID)
	}

	numChunks := len(randomChunks)

	seederAlternator := 0
	appendChunkLock := sync.Mutex{}
	recreateSessionLock := sync.Mutex{}
	lastRecreateTime := int64(0)

	//Give the seeder a fair start with timers when a download is initiated
	//Potentionally this seeder was last queried 60 seconds ago for files and otherwise idle but online
	//todo: Lock seeders
	for _, seeder := range GetSeeders(fileID) {
		sessionmanager.UpdateActivity(seeder)
	}

	downloadJob := func(terminateFlag *bool) {

		//Used to terminate the rescanning of peers
		terminate := func(flag *bool) {
			*flag = true
		}
		defer terminate(terminateFlag)

		for i := 0; i < numChunks; i++ {
			dbFile, err := dbGetFile(fileID)

			//Check if file is still tracked in surge
			if err != nil {
				log.Println("Download Job Terminated", "File no longer in DB")
				return
			}

			//Pause if file is paused
			for dbFile.IsPaused {
				return
			}

			//Create a async job to download a chunk
			requestChunkJob := func(chunkID int, downloadSeederAddr string) {
				requeue := func() {
					fmt.Println("Chunk ID", chunkID, " failed, and is being listed to be fetched again.")
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
				}

				recreateSessionLock.Lock()
				session, err := sessionmanager.GetSession(downloadSeederAddr)
				recreateSessionLock.Unlock()

				successRequest := false
				if err == nil {
					successRequest = RequestChunk(session, fileID, int32(chunkID))
				} else {
					//No session could be made, drop the file seeder.
					RemoveFileSeeder(fileID, downloadSeederAddr)
				}

				//if download fails return
				if !successRequest {
					requeue()
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
					/*existingSession, any := sessionmanager.GetExistingSession(downloadSeederAddr, constants.WorkerGetSessionTimeout)
					if existingSession == nil || !any {
						fmt.Println(string("\033[36m"), "No session with remote client available while waiting for requested chunk to arrive.", downloadSeederAddr, string("\033[0m"))
						inTransit = true
						sleepWorker = false
						break
					}*/

					//Check if received
					mutexes.ChunkInTransitLock.Lock()
					isInTransit := chunksInTransit[chunkKey]
					mutexes.ChunkInTransitLock.Unlock()

					if !isInTransit {
						//chunk received! if no longer in transit, continue workers
						inTransit = false
						sleepWorker = false
						break
					}

					_, currentSessionExists := sessionmanager.GetExistingSessionWithoutClosing(downloadSeederAddr, constants.WorkerGetSessionTimeout)

					if receiveTimeoutCounter >= constants.WorkerChunkReceiveTimeout && !currentSessionExists {
						//if timeout is triggered, leave in transit.
						log.Println(string("\033[36m"), "timeout is triggered, leave in transit.", string("\033[0m"))
						inTransit = true
						sleepWorker = false

						//Try replacing the session with a new one.
						lockTime := time.Now().Unix()

						recreateSessionLock.Lock()
						if lastRecreateTime > lockTime {
							//a new session was created.
							requeue()
							log.Println(string("\033[36m"), "using newly replaced session.", string("\033[0m"))
							recreateSessionLock.Unlock()
							return
						}

						_, err := sessionmanager.ReplaceSession(downloadSeederAddr)
						if err == nil {
							lastRecreateTime = time.Now().Unix()
						}
						recreateSessionLock.Unlock()

						if err == nil {
							log.Println(string("\033[36m"), "session was replaced, continue downloading.", string("\033[0m"))
						} else {
							log.Println(string("\033[36m"), "a new session could not be created, stop downloading.", string("\033[0m"))
							//RemoveSeeder(downloadSeederAddr)
							sessionmanager.CloseSession(downloadSeederAddr)
							requeue()
							return
						}

						break
					}
					receiveTimeoutCounter++
				}

				//If its still in transit abort
				if inTransit {
					requeue()
					return
				}
			}

			downloadSeederAddr := ""

			//spin to seeder with workers available
			spinForSeeder := true
			for spinForSeeder {
				for !AnySeeders(fileID) {
					fmt.Println(string("\033[36m"), "sleeping for seeders.", string("\033[0m"))
					time.Sleep(time.Second)
				}

				downloadSeederAddr = GetSeeders(fileID)[seederAlternator]

				//If seeder selected exceeds worker limit skip
				mutexes.WorkerMapLock.Lock()
				seedWorkerNum := workerMap[downloadSeederAddr]
				mutexes.WorkerMapLock.Unlock()

				if seedWorkerNum >= getNumberWorkers() {
					seederAlternator++
					if seederAlternator > len(GetSeeders(fileID))-1 {
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

	_, err = os.Stat(fileInfo.Path)
	//When we have an OS read error on the file mark the file as missing, stop down and uploads on it
	if err != nil {
		log.Println("Error on transmit chunk - file no longer at path, stopping upload.", err.Error())

		mutexes.FileWriteLock.Lock()
		fileInfo.IsMissing = true
		fileInfo.IsDownloading = false
		fileInfo.IsUploading = false
		fileInfo.IsAvailable = false
		dbInsertFile(*fileInfo)
		mutexes.FileWriteLock.Unlock()

		return
	}

	file, err := os.Open(fileInfo.Path)
	if err != nil {
		log.Println("Error on transmit chunk - file read failure", err.Error())
		return
	}

	//Read the requested chunk
	chunkOffset := int64(ChunkID) * constants.ChunkSize
	buffer := make([]byte, constants.ChunkSize)
	bytesread, err := file.ReadAt(buffer, chunkOffset)
	file.Close()

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
			log.Println("Session read", err)
			return nil, 0x0, err
		}
		log.Println("Session read", err)
		return nil, 0x0, err
	}

	//Get the packid
	packID := headerBuffer[0]

	//Get the size from the bytes
	sizeBytes := append(headerBuffer[:0], headerBuffer[1:]...)

	size := binary.LittleEndian.Uint32(sizeBytes)

	data = make([]byte, size)

	// read the full message, or return an error
	_, err = io.ReadFull(Session.Reader, data[:int(size)])
	if err != nil {
		log.Println("Session read", err)
		return nil, 0x0, err
	}

	return data, packID, err
}

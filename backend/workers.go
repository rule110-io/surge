// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This file contains all surge workers
	Workers run in a separate thread as long as the app gets terminated
*/

package surge

import (
	"fmt"
	"log"
	"time"

	"github.com/rule110-io/surge/backend/constants"
	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/mutexes"
	"github.com/rule110-io/surge/backend/platform"
	"github.com/rule110-io/surge/backend/sessionmanager"
)

//Function that automatically grabs subscriptions for nkn topic
func rescanPeersWorker() {
	for true {
		time.Sleep(constants.RescanPeerInterval)
		GetSubscriptions()
	}
}

// makes sure the client is regularly subscribed to the surge topic
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

// takes care that file data is regularly updated and stored in the database
func updateFileDataWorker() {

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

		statusBundle := []models.FileStatusEvent{}

		//Insert uploads
		allFiles := dbGetAllFiles()
		for _, file := range allFiles {
			if file.IsUploading {
				fileProgressMap[file.FileHash] = 1
			}
			key := file.FileHash

			//if file.IsPaused {
			//	continue
			//}

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
				statusEvent := models.FileStatusEvent{
					FileHash:          key,
					Progress:          fileProgressMap[key],
					DownloadBandwidth: down,
					UploadBandwidth:   up,
					NumChunks:         file.NumChunks,
					ChunkMap:          GetFileChunkMapString(&file, 156),
					ChunksShared:      file.ChunksShared,
				}
				statusBundle = append(statusBundle, statusEvent)
			}

			zeroBandwidthMap[key] = down+up == 0
		}

		//Add peer discovery global bandwidth
		mutexes.BandwidthAccumulatorMapLock.Lock()
		totalDown += downloadBandwidthAccumulator["DISCOVERY"]
		totalUp += uploadBandwidthAccumulator["DISCOVERY"]
		downloadBandwidthAccumulator["DISCOVERY"] = 0
		uploadBandwidthAccumulator["DISCOVERY"] = 0
		mutexes.BandwidthAccumulatorMapLock.Unlock()

		if zeroBandwidthMap["total"] == false || totalDown+totalUp != 0 {
			wailsRuntime.Events.Emit("globalBandwidthUpdate", statusBundle, totalDown, totalUp)
		}

		zeroBandwidthMap["total"] = totalDown+totalUp == 0
	}
}
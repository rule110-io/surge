package surge

import (
	"strconv"
	"strings"
	"sync"

	"log"

	"github.com/rule110-io/surge/backend/constants"
)

var fileWriteLock = &sync.Mutex{}

// needs to be rewritten! Should only parse the payload string and return the file objects. No download start!
// ParsePayloadString parses payload of files
func ParsePayloadString(s string) {

	payloadSplit := strings.Split(s, "surge://")
	for j := 0; j < len(payloadSplit); j++ {
		data := strings.Split(payloadSplit[j], "|")

		if len(data) < 3 {
			continue
		}

		fileSize, _ := strconv.ParseInt(data[3], 10, 64)
		numChunks := int((fileSize-1)/int64(constants.ChunkSize)) + 1

		seeder := strings.Split(data[5], ",")

		newListing := File{
			FileName:  data[2],
			FileSize:  fileSize,
			FileHash:  data[4],
			seeders:   seeder,
			Path:      "",
			NumChunks: numChunks,
			ChunkMap:  nil,
		}

		ListedFilesLock.Lock()
		//Replace existing, or remove.
		var replace = false
		for l := 0; l < len(ListedFiles); l++ {
			if ListedFiles[l].FileHash == newListing.FileHash {
				//if the seeder is unique add it as an additional seeder for the file
				ListedFiles[l].seeders = append(ListedFiles[l].seeders, seeder...)
				replace = true
				break
			}
		}
		//Unique listing so we add
		if replace == false {
			ListedFiles = append(ListedFiles, newListing)
		}
		ListedFilesLock.Unlock()

		log.Println("Program paramater new file: ", newListing.FileName, " seeder: ", newListing.seeders)

		go DownloadFileByHash(newListing.FileHash)
	}
}

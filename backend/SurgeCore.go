package surge

import (
	"strconv"
	"strings"

	"log"

	"github.com/rule110-io/surge/backend/constants"
	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/mutexes"
)

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

		newListing := models.GeneralFile{
			FileLocation: "remote",
			FileName:     data[2],
			FileSize:     fileSize,
			FileHash:     data[4],
			Seeders:      seeder,
			Path:         "",
			NumChunks:    numChunks,
			ChunkMap:     nil,
		}

		mutexes.ListedFilesLock.Lock()
		//Replace existing, or remove.
		var replace = false
		for l := 0; l < len(ListedFiles); l++ {
			if ListedFiles[l].FileHash == newListing.FileHash {
				//if the seeder is unique add it as an additional seeder for the file
				ListedFiles[l].Seeders = append(ListedFiles[l].Seeders, seeder...)
				replace = true
				break
			}
		}
		//Unique listing so we add
		if replace == false {
			ListedFiles = append(ListedFiles, newListing)
		}
		mutexes.ListedFilesLock.Unlock()

		log.Println("Program paramater new file: ", newListing.FileName, " seeder: ", newListing.Seeders)

		go DownloadFileByHash(newListing.FileHash)
	}
}

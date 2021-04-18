package surge

import (
	b64 "encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/rule110-io/surge/backend/constants"
	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/mutexes"
)

func removeStringFromSlice(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func distinctStringSlice(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

type sortBySeederCountAsc []models.File

func (a sortBySeederCountAsc) Len() int           { return len(a) }
func (a sortBySeederCountAsc) Less(i, j int) bool { return a[i].SeederCount < a[j].SeederCount }
func (a sortBySeederCountAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type sortBySeederCountDesc []models.File

func (a sortBySeederCountDesc) Len() int           { return len(a) }
func (a sortBySeederCountDesc) Less(i, j int) bool { return a[i].SeederCount > a[j].SeederCount }
func (a sortBySeederCountDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type sortByFileNameAsc []models.File

func (a sortByFileNameAsc) Len() int { return len(a) }
func (a sortByFileNameAsc) Less(i, j int) bool {
	return strings.ToLower(a[i].FileName) < strings.ToLower(a[j].FileName)
}
func (a sortByFileNameAsc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortByFileNameDesc []models.File

func (a sortByFileNameDesc) Len() int { return len(a) }
func (a sortByFileNameDesc) Less(i, j int) bool {
	return strings.ToLower(a[i].FileName) > strings.ToLower(a[j].FileName)
}
func (a sortByFileNameDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortByFileSizeAsc []models.File

func (a sortByFileSizeAsc) Len() int           { return len(a) }
func (a sortByFileSizeAsc) Less(i, j int) bool { return a[i].FileSize < a[j].FileSize }
func (a sortByFileSizeAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type sortByFileSizeDesc []models.File

func (a sortByFileSizeDesc) Len() int           { return len(a) }
func (a sortByFileSizeDesc) Less(i, j int) bool { return a[i].FileSize > a[j].FileSize }
func (a sortByFileSizeDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

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

//TopicEncode .
func TopicEncode(topic string) string {
	return "SRG_" + strings.ReplaceAll(b64.StdEncoding.EncodeToString([]byte(topic)), "=", "-")
}

func surgeGenerateTopicPayload(fileName string, sizeInBytes int64, hash string, topic string) string {
	//Example payload
	//surge://|file|The_Two_Towers-The_Purist_Edit-Trailer.avi|14997504|965c013e991ee246d63d45ea71954c4d|/

	return "surge://|file|" + fileName + "|" + strconv.FormatInt(sizeInBytes, 10) + "|" + hash + "|" + topic + "|/"
}

func surgeGenerateMagnetLink(fileName string, sizeInBytes int64, hash string, seeder string, topic string) string {
	//Example payload
	//surge://|file|The_Two_Towers-The_Purist_Edit-Trailer.avi|14997504|965c013e991ee246d63d45ea71954c4d|/
	if seeder == "" {
		seeder = client.Addr().String()
	}

	return "surge://|file|" + fileName + "|" + strconv.FormatInt(sizeInBytes, 10) + "|" + hash + "|" + seeder + "|" + topic + "|/"
}

func hashFile(randomHash string) {
	//Clean up after were done here, even when we fail we dont want these randomhash files in db
	defer dbDeleteFile(randomHash)

	dbFile, err := dbGetFile(randomHash)
	if err != nil {
		pushError("File Hash Failed", "Could find dbEntry for hash "+randomHash)
	}

	hashString, err := HashFile(dbFile.Path)
	if err != nil {
		log.Println(err)
		pushError("File Hash Failed", "Could not hash file at "+dbFile.Path)
	}

	dbFile.IsUploading = true
	dbFile.IsHashing = false
	dbFile.FileHash = hashString
	dbInsertFile(*dbFile)

	//Add to payload
	AddToSeedString(*dbFile)
	pushNotification("Now seeding", dbFile.FileName)
}

// ParsePayloadString parses payload of files
func ParsePayloadString(s string) []models.File {

	files := []models.File{}
	payloadSplit := strings.Split(s, "surge://")
	for j := 0; j < len(payloadSplit); j++ {
		data := strings.Split(payloadSplit[j], "|")

		if len(data) < 3 {
			continue
		}

		fileSize, _ := strconv.ParseInt(data[3], 10, 64)
		numChunks := int((fileSize-1)/int64(constants.ChunkSize)) + 1

		seeder := strings.Split(data[5], ",")

		newListing := models.File{
			FileLocation: "remote",
			FileName:     data[2],
			FileSize:     fileSize,
			FileHash:     data[4],
			Seeders:      seeder,
			Path:         "",
			NumChunks:    numChunks,
			ChunkMap:     nil,
			Topic:        data[5],
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
		files = append(files, newListing)
	}
	return files
}

func filterFile(ss []models.File, test func(models.File) bool) (ret []models.File) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

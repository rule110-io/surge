package surge

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"log"

	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/platform"

	"github.com/xujiajun/nutsdb"
)

const fileBucketName = "fileBucket"
const settingBucketName = "settingsBucket"

var db *nutsdb.DB

//InitializeDb initializes db
func InitializeDb() {
	var err error
	opt := nutsdb.DefaultOptions

	opt.Dir = platform.GetSurgeDir() + string(os.PathSeparator) + "db"
	db, err = nutsdb.Open(opt)
	if err != nil {
		log.Panic(err)
	}

	dbGetAllFiles()

}

//CloseDb .
func CloseDb() {
	db.Close()
}

func dbInsertFile(File models.GeneralFile) {
	if err := db.Update(
		func(tx *nutsdb.Tx) error {

			fileKey := []byte(File.FileHash)
			fileBytes, _ := json.Marshal(File)

			if err := tx.Put(fileBucketName, fileKey, fileBytes, 0); err != nil {
				return err
			}
			return nil
		}); err != nil {
		log.Panic(err)
	}
}

func dbGetFile(Key string) (*models.GeneralFile, error) {
	result := &models.GeneralFile{}

	if err := db.View(
		func(tx *nutsdb.Tx) error {
			fileKey := []byte(Key)
			e, err := tx.Get(fileBucketName, fileKey)
			if err != nil {
				return err
			}

			err = json.Unmarshal(e.Value, result)
			return err
		}); err != nil {
		return nil, err
	}

	return result, nil
}

func dbGetAllFiles() []models.GeneralFile {
	files := []models.GeneralFile{}

	if err := db.View(
		func(tx *nutsdb.Tx) error {
			entries, err := tx.GetAll(fileBucketName)
			if err != nil {
				return err
			}

			for _, entry := range entries {

				newFile := &models.GeneralFile{}
				json.Unmarshal(entry.Value, newFile)
				files = append(files, *newFile)
			}

			return nil
		}); err != nil {
		log.Println(err)
	} else {
		return files
	}
	return files
}

func dbDeleteFile(Hash string) error {
	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte(Hash)
			if err := tx.Delete(fileBucketName, key); err != nil {
				return err
			}
			return nil
		}); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Stores or updates a key with a given value
func DbWriteSetting(Name string, value string) error {
	err := db.Update(
		func(tx *nutsdb.Tx) error {

			keyBytes := []byte(Name)
			valueBytes := []byte(value)

			if err := tx.Put(settingBucketName, keyBytes, valueBytes, 0); err != nil {
				return err
			}
			return nil
		})
	return err
}

// Reads a key and returns value
func DbReadSetting(Name string) (string, error) {
	result := ""
	key := []byte(Name)

	if err := db.View(
		func(tx *nutsdb.Tx) error {
			bytes, err := tx.Get(settingBucketName, key)
			if err != nil {
				return err
			}

			result = string(bytes.Value)

			return err
		}); err != nil {
		return "", err
	}

	return result, nil
}

//SearchQueryResult is a paging query result for file searches
type SearchQueryResult struct {
	Result []models.GeneralFile
	Count  int
}

//LocalFilePageResult is a paging query result for tracked files
type LocalFilePageResult struct {
	Result []models.GeneralFile
	Count  int
}

//SearchRemoteFile runs a paged query
func SearchRemoteFile(Query string, OrderBy string, IsDesc bool, Skip int, Take int) SearchQueryResult {

	var results []models.GeneralFile

	ListedFilesLock.Lock()
	for _, file := range ListedFiles {
		if strings.Contains(strings.ToLower(file.FileName), strings.ToLower(Query)) || strings.Contains(strings.ToLower(file.FileHash), strings.ToLower(Query)) && file.FileLocation == "remote" {

			result := models.GeneralFile{
				FileName:    file.FileName,
				FileHash:    file.FileHash,
				FileSize:    file.FileSize,
				Seeders:     file.Seeders,
				NumChunks:   file.NumChunks,
				SeederCount: len(file.Seeders),
			}

			tracked, err := dbGetFile(result.FileHash)

			//only add non-local files to the result
			if err != nil && tracked == nil {
				results = append(results, result)
			}

		}
	}
	ListedFilesLock.Unlock()

	switch OrderBy {
	case "FileName":
		if !IsDesc {
			sort.Sort(sortByFileNameAsc(results))
		} else {
			sort.Sort(sortByFileNameDesc(results))
		}
	case "FileSize":
		if !IsDesc {
			sort.Sort(sortByFileSizeAsc(results))
		} else {
			sort.Sort(sortByFileSizeDesc(results))
		}
	default:
		if !IsDesc {
			sort.Sort(sortBySeederCountAsc(results))
		} else {
			sort.Sort(sortBySeederCountDesc(results))
		}
	}

	left := Skip
	right := Skip + Take

	if left > len(results) {
		left = len(results)
	}

	if right > len(results) {
		right = len(results)
	}

	return SearchQueryResult{
		Result: results[left:right],
		Count:  len(results),
	}
}

//SearchLocalFile runs a paged query
func SearchLocalFile(Query string, OrderBy string, IsDesc bool, Skip int, Take int) LocalFilePageResult {

	resultFiles := []models.GeneralFile{}

	allFiles := dbGetAllFiles()
	for _, file := range allFiles {
		if strings.Contains(strings.ToLower(file.FileName), strings.ToLower(Query)) || strings.Contains(strings.ToLower(file.FileHash), strings.ToLower(Query)) && file.FileLocation == "local" {
			resultFiles = append(resultFiles, file)
		}
	}

	totalNum := len(resultFiles)

	switch OrderBy {
	default:
		if !IsDesc {
			sort.Sort(sortByFileNameAsc(resultFiles))
		} else {
			sort.Sort(sortByFileNameDesc(resultFiles))
		}
	}

	left := Skip
	right := Skip + Take

	if left > len(resultFiles) {
		left = len(resultFiles)
	}

	if right > len(resultFiles) {
		right = len(resultFiles)
	}

	//Subset
	resultFiles = resultFiles[left:right]
	resultListings := []models.GeneralFile{}

	ListedFilesLock.Lock()
	for i := 0; i < len(resultFiles); i++ {
		listing := models.GeneralFile{
			ChunksShared:  resultFiles[i].ChunksShared,
			FileHash:      resultFiles[i].FileHash,
			FileName:      resultFiles[i].FileName,
			FileSize:      resultFiles[i].FileSize,
			IsDownloading: resultFiles[i].IsDownloading,
			IsHashing:     resultFiles[i].IsHashing,
			IsMissing:     resultFiles[i].IsMissing,
			IsPaused:      resultFiles[i].IsPaused,
			IsUploading:   resultFiles[i].IsUploading,
			NumChunks:     resultFiles[i].NumChunks,
			Path:          resultFiles[i].Path,
		}

		if listing.IsUploading {
			listing.Seeders = []string{GetAccountAddress()}
		} else {
			listing.Seeders = []string{}
		}

		for _, file := range ListedFiles {
			if file.FileHash == listing.FileHash {
				listing.Seeders = append(listing.Seeders, file.Seeders...)
				break
			}
		}
		listing.SeederCount = len(listing.Seeders)
		//If file is downloading set progress
		if listing.IsDownloading || listing.IsPaused {
			numChunksLocal := chunksDownloaded(resultFiles[i].ChunkMap, listing.NumChunks)
			listing.Progress = float32(float64(numChunksLocal) / float64(listing.NumChunks))
		}

		resultListings = append(resultListings, listing)

	}
	ListedFilesLock.Unlock()

	return LocalFilePageResult{
		Result: resultListings,
		Count:  totalNum,
	}
}

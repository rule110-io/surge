// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This file contains DB related functions
	It takes care of initializing the db as well as querying and processing DB entries
*/

package surge

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"log"

	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/mutexes"
	"github.com/rule110-io/surge/backend/platform"

	"github.com/xujiajun/nutsdb"
)

const fileBucketName = "fileBucket"
const settingBucketName = "settingsBucket"

var db *nutsdb.DB

type FileFilterState int

const (
	All = iota
	Downloading
	Seeding
	Completed
	Paused
)

//InitializeDb initializes db
func InitializeDb() {
	var err error
	opt := nutsdb.DefaultOptions

	opt.Dir = platform.GetSurgeDir() + string(os.PathSeparator) + "db"
	db, err = nutsdb.Open(opt)
	if err != nil {
		log.Panic(err)
	}
}

//CloseDb .
func CloseDb() {
	db.Close()
}

// Gets all Files in the DB
func dbGetAllFiles() []models.File {
	files := []models.File{}

	if err := db.View(
		func(tx *nutsdb.Tx) error {
			entries, err := tx.GetAll(fileBucketName)
			if err != nil {
				return err
			}

			for _, entry := range entries {

				newFile := &models.File{}
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

// Gets a File by providing the fileHash
func dbGetFile(Hash string) (*models.File, error) {
	result := &models.File{}

	if err := db.View(
		func(tx *nutsdb.Tx) error {
			fileKey := []byte(Hash)
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

// Inserts a File to the DB
func dbInsertFile(File models.File) {
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

// Deletes a File by providing the fileHash
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

//DbWriteSetting Stores or updates a key with a given value
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

//DbReadSetting Reads a key and returns value
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

//PagedQueryResult is a paging query result for file searches
type PagedQueryResult struct {
	Result []models.File
	Count  int
}

//PagedQueryResult is a paging query result for file searches
type PagedQueryRemoteResult struct {
	Result []models.FileListing
	Count  int
}

//SearchRemoteFile runs a paged query
func SearchRemoteFile(Topic string, Query string, OrderBy string, IsDesc bool, Skip int, Take int) PagedQueryRemoteResult {

	var results []models.FileListing

	mutexes.ListedFilesLock.Lock()
	for _, file := range ListedFiles {

		if file.Topic == Topic {

			if strings.Contains(strings.ToLower(file.FileName), strings.ToLower(Query)) || strings.Contains(strings.ToLower(file.FileHash), strings.ToLower(Query)) {

				localFile, _ := dbGetFile(file.FileHash)

				if localFile != nil {
					result := models.FileListing{
						FileName:      file.FileName,
						FileHash:      file.FileHash,
						FileSize:      file.FileSize,
						NumChunks:     file.NumChunks,
						Topic:         file.Topic,
						NumSeeders:    len(GetSeeders(file.FileHash)),
						IsTracked:     true,
						IsDownloading: file.IsDownloading,
						IsUploading:   file.IsUploading,
					}
					results = append(results, result)
				} else {
					result := models.FileListing{
						FileName:      file.FileName,
						FileHash:      file.FileHash,
						FileSize:      file.FileSize,
						NumChunks:     file.NumChunks,
						Topic:         file.Topic,
						NumSeeders:    len(GetSeeders(file.FileHash)),
						IsTracked:     false,
						IsDownloading: file.IsDownloading,
						IsUploading:   file.IsUploading,
					}
					results = append(results, result)
				}

			}
		}
	}
	mutexes.ListedFilesLock.Unlock()

	switch OrderBy {
	case "FileName":
		if !IsDesc {
			sort.Sort(sortByListingFileNameAsc(results))
		} else {
			sort.Sort(sortByListingFileNameDesc(results))
		}
	case "FileSize":
		if !IsDesc {
			sort.Sort(sortByListingFileSizeAsc(results))
		} else {
			sort.Sort(sortByListingFileSizeDesc(results))
		}
	default:
		if !IsDesc {
			sort.Sort(sortByListingSeederCountAsc(results))
		} else {
			sort.Sort(sortByListingSeederCountDesc(results))
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

	return PagedQueryRemoteResult{
		Result: results[left:right],
		Count:  len(results),
	}
}

//SearchLocalFile runs a paged query
func SearchLocalFile(Query string, filterState FileFilterState, OrderBy string, IsDesc bool, Skip int, Take int) PagedQueryResult {

	resultFiles := []models.File{}

	allFiles := dbGetAllFiles()
	for _, file := range allFiles {
		if strings.Contains(strings.ToLower(file.FileName), strings.ToLower(Query)) || strings.Contains(strings.ToLower(file.FileHash), strings.ToLower(Query)) {
			resultFiles = append(resultFiles, file)
		}
	}

	fileFilterFunc := func(f models.File) bool { return true }

	//Filter files on filter state
	switch filterState {
	/*case All: //added for clarity
	fileFilterFunc = func(f models.File) bool { return true }
	break*/
	case Downloading:
		fileFilterFunc = func(f models.File) bool { return f.IsDownloading }
	case Seeding:
		fileFilterFunc = func(f models.File) bool { return f.IsUploading }
	case Completed:
		fileFilterFunc = func(f models.File) bool { return f.IsAvailable }
	case Paused:
		fileFilterFunc = func(f models.File) bool { return f.IsPaused }
	}

	resultFiles = filterFile(resultFiles, fileFilterFunc)

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
	resultListings := []models.File{}

	mutexes.ListedFilesLock.Lock()
	for i := 0; i < len(resultFiles); i++ {
		listing := models.File{
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
			Topic:         resultFiles[i].Topic,
		}

		//If file is downloading set progress
		if listing.IsDownloading || listing.IsPaused {
			numChunksLocal := chunksDownloaded(resultFiles[i].ChunkMap, listing.NumChunks)
			listing.Progress = float32(float64(numChunksLocal) / float64(listing.NumChunks))
		}

		resultListings = append(resultListings, listing)

	}
	mutexes.ListedFilesLock.Unlock()

	return PagedQueryResult{
		Result: resultListings,
		Count:  totalNum,
	}
}

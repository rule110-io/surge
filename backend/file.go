// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This file contains functions related to files in surge
*/

package surge

import (
	"log"
	"os"

	bitmap "github.com/boljen/go-bitmap"
	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/mutexes"
)

func getListedFileByHash(Hash string) *models.File {

	var selectedFile *models.File = nil

	mutexes.ListedFilesLock.Lock()
	for _, file := range ListedFiles {
		if file.FileHash == Hash {
			selectedFile = &file
			break
		}
	}
	mutexes.ListedFilesLock.Unlock()

	return selectedFile
}

//GetFileChunkMapString returns the chunkmap in hex for a file given by hash
func GetFileChunkMapString(file *models.File, Size int) string {

	outputSize := Size
	inputSize := file.NumChunks

	stepSize := float64(inputSize) / float64(outputSize)
	stepSizeInt := int(stepSize)

	var boolBuffer = ""
	if inputSize >= outputSize {

		for i := 0; i < outputSize; i++ {
			localCount := 0
			for j := 0; j < stepSizeInt; j++ {
				local := bitmap.Get(file.ChunkMap, int(float64(i)*stepSize)+j)
				if local {
					localCount++
				} else {
					if localCount == 0 {
						boolBuffer += "0"
					} else {
						boolBuffer += "1"
					}
					break
				}
			}
			if localCount == stepSizeInt {
				boolBuffer += "2"
			}
		}
	} else {
		iter := float64(0)
		for i := 0; i < outputSize; i++ {
			local := bitmap.Get(file.ChunkMap, int(iter))
			if local {
				boolBuffer += "2"
			} else {
				boolBuffer += "0"
			}
			iter += stepSize
		}
	}
	return boolBuffer
}

//SetFilePause sets a file IsPaused state for by file hash
func SetFilePause(Hashes []string, State bool) {

	mutexes.FileWriteLock.Lock()
	defer mutexes.FileWriteLock.Unlock()

	for _, hash := range Hashes {
		file, err := dbGetFile(hash)
		if err != nil {
			pushNotification("Failed To Pause", "Could not find the file to pause.")
		}

		if file.IsPaused != State {
			file.IsPaused = State
			dbInsertFile(*file)

			if !file.IsPaused && file.IsDownloading {
				go restartDownload(file.FileHash)
			}
		}
	}
}

//RemoveFileByHash removes file from surge db and optionally from disk
func RemoveFileByHash(Hash string, FromDisk bool) bool {

	mutexes.FileWriteLock.Lock()

	file, err := dbGetFile(Hash)
	if !FromDisk && err != nil {
		pushError("Error on remove file (read db)", err.Error())
		return false
	}

	if FromDisk {
		err = os.Remove(file.Path)
		if err != nil {
			pushError("Error on remove file from disk", err.Error())
		}
	}

	err = dbDeleteFile(Hash)
	if err != nil {
		pushError("Error on remove file (read db)", err.Error())
		return false
	}
	mutexes.FileWriteLock.Unlock()

	log.Println("Removing file:", file.FileName, file.FileHash, "from disk:", FromDisk)
	AnnounceRemoveFile(file.Topic, file.FileHash)
	return true
}

//GetFileChunkMapStringByHash returns the chunkmap in hex for a file given by hash
func GetFileChunkMapStringByHash(Hash string, Size int) string {

	file, err := dbGetFile(Hash)
	if err != nil {
		return ""
	}
	return GetFileChunkMapString(file, Size)
}

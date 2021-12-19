// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This file contains all local input and output functions
	This can be notifications sent to the frontend or storing procedures
*/

package surge

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	bitmap "github.com/boljen/go-bitmap"
	"github.com/rule110-io/surge/backend/constants"
	"github.com/rule110-io/surge/backend/mutexes"
	"github.com/rule110-io/surge/backend/platform"
	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func emitNotificationEvent(event string, title string, text string) {
	runtime.EventsEmit(*wailsContext, "notificationEvent", title, text, time.Now().Unix())
}

func pushNotification(title string, text string) {
	//If wails frontend is not yet bound, we wait in a task to not block main thread
	if !FrontendReady {
		waitAndPush := func() {
			for !FrontendReady {
				time.Sleep(time.Millisecond * 50)
			}
			emitNotificationEvent("notificationEvent", title, text)
		}
		go waitAndPush()
	} else {
		emitNotificationEvent("notificationEvent", title, text)
	}
}

func pushError(title string, text string) {
	//If wails frontend is not yet bound, we wait in a task to not block main thread
	if !FrontendReady {
		waitAndPush := func() {
			for !FrontendReady {
				time.Sleep(time.Millisecond * 50)
			}
			emitNotificationEvent("errorEvent", title, text)
		}
		go waitAndPush()
	} else {
		emitNotificationEvent("errorEvent", title, text)
	}
}

//SetVisualMode Sets the visualmode
func SetVisualMode(visualMode int) {
	if visualMode == 0 {
		//light mode
		DbWriteSetting("DarkMode", "false")
		runtime.EventsEmit(*wailsContext, "darkThemeEvent", "false")
	} else if visualMode == 1 {
		//dark mode
		DbWriteSetting("DarkMode", "true")
		runtime.EventsEmit(*wailsContext, "darkThemeEvent", "true")
	}
}

//WriteChunk writes a chunk to disk
func WriteChunk(FileID string, ChunkID int32, Chunk []byte) {
	defer RecoverAndLog()

	fileInfo, err := dbGetFile(FileID)
	if err != nil {
		log.Println("Error on write chunk (db get)", err.Error())
		return
	}
	remoteFolder, err := platform.GetRemoteFolder()
	if err != nil {
		pushError("Error on write chunk (GetRemoteFolder)", err.Error())
		return
	}
	path := remoteFolder + string(os.PathSeparator) + fileInfo.FileName
	osFile, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		pushError("Error on write chunk (OpenFile)", err.Error())
		return
	}
	defer osFile.Close()

	if err != nil {
		pushError("Error on write chunk (sessionmanager.GetFile)", err.Error())
		return
	}

	chunkOffset := int64(ChunkID) * constants.ChunkSize
	bytesWritten, err := osFile.WriteAt(Chunk, chunkOffset)
	if err != nil {
		pushError("Error on write chunk (file write)", err.Error())
		return
	}
	//Success
	log.Println("Chunk written to disk: ", bytesWritten, " bytes")

	//Update bitmap async as this has a lock in it but does not have to be waited for
	setBitMap := func() {
		mutexes.FileWriteLock.Lock()

		//Set chunk to available in the map
		fileInfo, err := dbGetFile(FileID)
		if err != nil {
			pushError("Error on chunk write (db get)", err.Error())
			return
		}
		bitmap.Set(fileInfo.ChunkMap, int(ChunkID), true)
		dbInsertFile(*fileInfo)

		mutexes.FileWriteLock.Unlock()
	}
	go setBitMap()
}

//OpenOSPath Open a file, directory, or URI using the OS's default application for that object type. Don't wait for the open command to complete.
func OpenOSPath(Path string) {
	open.Start(Path)
}

//OpenFileByHash opens a file with OS default application for object type
func OpenFileByHash(Hash string) {

	fileInfo, err := dbGetFile(Hash)
	if err != nil {
		pushError("Error on open file", err.Error())
		return
	}
	OpenOSPath(fileInfo.Path)
}

//OpenFolderByHash opens the folder containing the file by hash in os
func OpenFolderByHash(Hash string) {

	fileInfo, err := dbGetFile(Hash)
	if err != nil {
		pushError("Error on open folder", err.Error())
		return
	}
	OpenOSPath(filepath.Dir(fileInfo.Path))
}

// AllocateFile Allocates a file on disk at path with size in bytes
func AllocateFile(path string, size int64) bool {

	fd, err := os.Create(path)
	if err != nil {
		pushError("File disk allocation error", "file could not be created at "+path)
		return false
	}
	_, err = fd.Seek(size-1, 0)
	if err != nil {
		pushError("File disk allocation error", "file was created but could not be read at "+path)
		return false
	}
	_, err = fd.Write([]byte{0})
	if err != nil {
		pushError("File disk allocation error", "file was read but could not be written at "+path)
		return false
	}
	err = fd.Close()
	if err != nil {
		pushError("File disk allocation error", "file was written but could not be released at "+path)
		return false
	}

	return true
}

// HashFile generates hash for file given filepath
func HashFile(filePath string) (string, error) {

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := sha256.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)

	//Convert the bytes to a string
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil

}

func surgeGetFileSize(path string) int64 {

	fi, err := os.Stat(path)
	if err != nil {
		log.Panic("Error on get filesize", err)
	}
	// get the size
	return fi.Size()
}

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

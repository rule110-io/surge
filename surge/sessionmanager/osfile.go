package sessionmanager

import (
	"log"
	"os"
	"sync"
)

var fileManagerLock = sync.Mutex{}

//A map to hold nkn files
var fileMap map[string]*os.File

//GetFile returns a file for given address
func GetFile(FileID string, Path string) (*os.File, error) {
	fileManagerLock.Lock()
	defer fileManagerLock.Unlock()

	//Check for an existing file
	file, exists := fileMap[Path]

	//create if it doesnt exist
	if !exists {
		var err error = nil
		file, err = openFile(Path)
		if err == nil {
			fileMap[FileID] = file
		}
	}

	return file, nil
}

//CloseFile closes the file
func CloseFile(FileHash string) {
	fileManagerLock.Lock()
	defer fileManagerLock.Unlock()

	file, exists := fileMap[FileHash]

	//Close nkn file, nill out the pointers
	if exists && file != nil {
		file.Sync()
		file.Close()
		file = nil
	}

	//Delete from the map
	delete(fileMap, FileHash)

	log.Println("file closed for: ", FileHash)
}

func openFile(Path string) (file *os.File, err error) {
	file, err = os.OpenFile(Path, os.O_RDWR, 0644)
	return file, err
}

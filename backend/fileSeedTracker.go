package surge

import "sync"

var fileSeedMap map[string][]string
var fileSeedLock = sync.Mutex{}

func InitializeFileSeedTracker() {
	fileSeedMap = make(map[string][]string)
	fileSeedLock = sync.Mutex{}
}

func AddFileSeeder(fileHash string, addr string) {
	fileSeedLock.Lock()
	defer fileSeedLock.Unlock()

	//check if slice exists, otherwise we must create it
	_, exists := fileSeedMap[fileHash]
	if !exists {
		fileSeedMap[fileHash] = []string{}
	}

	//Append and distinct so we dont double up
	fileSeedMap[fileHash] = append(fileSeedMap[fileHash], addr)
	fileSeedMap[fileHash] = distinctStringSlice(fileSeedMap[fileHash])
}

func RemoveFileSeeder(fileHash string, addr string) {
	fileSeedLock.Lock()
	defer fileSeedLock.Unlock()

	removeFileSeeder(fileHash, addr)
}

func removeFileSeeder(fileHash string, addr string) {
	fileSeedMap[fileHash] = removeStringFromSlice(fileSeedMap[fileHash], addr)

	//TODO: should we set to nil the filemap key when len = 0?
}

func RemoveSeeder(addr string) {
	fileSeedLock.Lock()
	defer fileSeedLock.Unlock()

	for k := range fileSeedMap {
		removeFileSeeder(k, addr)
	}
}

func AnySeeders(fileHash string) bool {
	fileSeedLock.Lock()
	defer fileSeedLock.Unlock()

	seeders, exists := fileSeedMap[fileHash]
	if !exists {
		return false
	}
	return len(seeders) > 0
}

func GetSeeders(fileHash string) []string {
	fileSeedLock.Lock()
	defer fileSeedLock.Unlock()

	return fileSeedMap[fileHash]
}

package surge

import (
	"sort"

	"github.com/rule110-io/surge/backend/constants"
	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/mutexes"
	"github.com/rule110-io/surge/backend/platform"
)

//MiddlewareFunctions struct to hold wails runtime for all middleware implementations
type MiddlewareFunctions struct {
}

func (s *MiddlewareFunctions) GetLocalFiles(Query string, filterState FileFilterState, OrderBy string, IsDesc bool, Skip int, Take int) PagedQueryResult {
	return SearchLocalFile(Query, filterState, OrderBy, IsDesc, Skip, Take)
}

//GetRemoteFiles gets remote files
func (s *MiddlewareFunctions) GetRemoteFiles(Topic string, Query string, OrderBy string, IsDesc bool, Skip int, Take int) PagedQueryRemoteResult {
	return SearchRemoteFile(Topic, Query, OrderBy, IsDesc, Skip, Take)
}

//DownloadFile download file by hash
func (s *MiddlewareFunctions) DownloadFile(Hash string) bool {
	return DownloadFileByHash(Hash)
}

//SetDownloadPause set pause state by hash (bool)
func (s *MiddlewareFunctions) SetDownloadPause(Hashes []string, State bool) {
	SetFilePause(Hashes, State)
}

//GetPublicKey retrieves account pubkey
func (s *MiddlewareFunctions) GetPublicKey() string {
	return GetAccountAddress()
}

//GetFileChunkMap returns chunkmap for file by hash with desired length param
func (s *MiddlewareFunctions) GetFileChunkMap(Hash string, Size int) string {
	if Size == 0 {
		Size = 400
	}
	return GetFileChunkMapStringByHash(Hash, Size)
}

//OpenFile open file by given hash (os)
func (s *MiddlewareFunctions) OpenFile(Hash string) {
	OpenFileByHash(Hash)
}

//OpenLink open url given by param
func (s *MiddlewareFunctions) OpenLink(Link string) {
	OpenOSPath(Link)
}

//OpenLog opens log file in os
func (s *MiddlewareFunctions) OpenLog() {
	OpenLogFile()
}

//OpenFolder opens folder for file given by hash
func (s *MiddlewareFunctions) OpenFolder(Hash string) {
	OpenFolderByHash(Hash)
}

func (s *MiddlewareFunctions) SeedFile(Topic string) bool {
	path := platform.OpenFileDialog()
	if path == "" {
		return false
	}
	return SeedFilepath(path, Topic)
}

//RemoveFile remove file from surge (and os) by hash
func (s *MiddlewareFunctions) RemoveFile(Hash string, FromDisk bool) bool {
	return RemoveFileByHash(Hash, FromDisk)
}

//WriteSetting generic kvs setting store
func (s *MiddlewareFunctions) WriteSetting(Key string, Value string) bool {
	err := DbWriteSetting(Key, Value)
	return err != nil
}

//ReadSetting generic kvs setting store
func (s *MiddlewareFunctions) ReadSetting(Key string) string {
	val, _ := DbReadSetting(Key)
	return val
}

//StartDownloadMagnetLinks initiate a download by magnet
func (s *MiddlewareFunctions) StartDownloadMagnetLinks(Magnetlinks string) bool {
	//need to parse Magnetlinks array and download all of them
	files := ParsePayloadString(Magnetlinks)
	for i := 0; i < len(files); i++ {
		go DownloadFileByHash(files[i].FileHash)
	}
	return true
}

//SubscribeToTopic subscribes to given topic
func (s *MiddlewareFunctions) SubscribeToTopic(Topic string) {
	if len(Topic) == 0 {
		pushError("Error on Subscribe", "topic name of length zero.")
	} else {
		subscribeToSurgeTopic(Topic, true)
	}
}

func (s *MiddlewareFunctions) UnsubscribeFromTopic(Topic string) {
	unsubscribeFromSurgeTopic(Topic)
}

func (s *MiddlewareFunctions) GetTopicSubscriptions() []string {
	topicNames := []string{}

	mutexes.TopicsMapLock.Lock()
	for _, v := range topicsMap {
		topicNames = append(topicNames, v.Name)
	}
	mutexes.TopicsMapLock.Unlock()
	sort.Strings(topicNames)

	return topicNames
}

type FileDetails struct {
	FileID  string
	Seeders []string
}

func (s *MiddlewareFunctions) GetFileDetails(FileHash string) FileDetails {
	return FileDetails{
		FileID:  FileHash,
		Seeders: GetSeeders(FileHash),
	}
}
func (s *MiddlewareFunctions) GetTopicDetails(Topic string) models.TopicInfo {

	return GetTopicInfo(Topic)
}

func (s *MiddlewareFunctions) GetOfficialTopicName() string {
	return constants.SurgeOfficialTopic
}

//

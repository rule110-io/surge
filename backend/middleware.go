package surge

import (
	"github.com/rule110-io/surge/backend/platform"
	"github.com/wailsapp/wails"
)

//MiddlewareFunctions struct to hold wails runtime for all middleware implementations
type MiddlewareFunctions struct {
	r *wails.Runtime
}

//GetLocalFiles gets local files
func (s *MiddlewareFunctions) GetLocalFiles(Query string, OrderBy string, IsDesc bool, Skip int, Take int) PagedQueryResult {
	return SearchLocalFile(Query, OrderBy, IsDesc, Skip, Take)
}

//GetRemoteFiles gets remote files
func (s *MiddlewareFunctions) GetRemoteFiles(Query string, OrderBy string, IsDesc bool, Skip int, Take int) PagedQueryResult {
	return SearchRemoteFile(Query, OrderBy, IsDesc, Skip, Take)
}

//DownloadFile download file by hash
func (s *MiddlewareFunctions) DownloadFile(Hash string) bool {
	return DownloadFileByHash(Hash)
}

//SetDownloadPause set pause state by hash (bool)
func (s *MiddlewareFunctions) SetDownloadPause(Hash string, State bool) {
	SetFilePause(Hash, State)
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

//SeedFile initiate os dialog to seed a file
func (s *MiddlewareFunctions) SeedFile() bool {
	path := platform.OpenFileDialog()
	if path == "" {
		return false
	}
	return SeedFilepath(path)
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

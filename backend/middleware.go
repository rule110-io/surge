package surge

import (
	"github.com/rule110-io/surge/backend/platform"
	"github.com/wailsapp/wails"
)

type MiddlewareFunctions struct {
	r *wails.Runtime
}

func (s *MiddlewareFunctions) GetLocalFiles(Query string, filterState FileFilterState, OrderBy string, IsDesc bool, Skip int, Take int) PagedQueryResult {
	return SearchLocalFile(Query, filterState, OrderBy, IsDesc, Skip, Take)
}

func (s *MiddlewareFunctions) GetRemoteFiles(Query string, OrderBy string, IsDesc bool, Skip int, Take int) PagedQueryResult {
	return SearchRemoteFile(Query, OrderBy, IsDesc, Skip, Take)
}

func (s *MiddlewareFunctions) DownloadFile(Hash string) bool {
	return DownloadFileByHash(Hash)
}

func (s *MiddlewareFunctions) SetDownloadPause(Hash string, State bool) {
	SetFilePause(Hash, State)
}

func (s *MiddlewareFunctions) GetPublicKey() string {
	return GetAccountAddress()
}

func (s *MiddlewareFunctions) GetFileChunkMap(Hash string, Size int) string {
	if Size == 0 {
		Size = 400
	}
	return GetFileChunkMapStringByHash(Hash, Size)
}

func (s *MiddlewareFunctions) OpenFile(Hash string) {
	OpenFileByHash(Hash)
}

func (s *MiddlewareFunctions) OpenLink(Link string) {
	OpenOSPath(Link)
}

func (s *MiddlewareFunctions) OpenLog() {
	OpenLogFile()
}

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

func (s *MiddlewareFunctions) RemoveFile(Hash string, FromDisk bool) bool {
	return RemoveFileByHash(Hash, FromDisk)
}

func (s *MiddlewareFunctions) WriteSetting(Key string, Value string) bool {
	err := DbWriteSetting(Key, Value)
	return err != nil
}

func (s *MiddlewareFunctions) ReadSetting(Key string) string {
	val, _ := DbReadSetting(Key)
	return val
}

func (s *MiddlewareFunctions) StartDownloadMagnetLinks(Magnetlinks string) bool {
	//need to parse Magnetlinks array and download all of them
	files := ParsePayloadString(Magnetlinks)
	for i := 0; i < len(files); i++ {
		go DownloadFileByHash(files[i].FileHash)
	}
	return true
}

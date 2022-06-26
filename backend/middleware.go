package surge

import (
	"fmt"
	"strconv"

	"github.com/nknorg/nkn-sdk-go"
	"github.com/rule110-io/surge/backend/constants"
	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/mutexes"
	"github.com/rule110-io/surge/backend/sessionmanager"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	path, _ := runtime.OpenFileDialog(*wailsContext, runtime.OpenDialogOptions{
		Title: "Select File",
	})
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

func (s *MiddlewareFunctions) GetTopicSubscriptions() []models.TopicInfo {
	return GetTopicsWithPermissions()
}

type FileDetails struct {
	FileID           string
	Seeders          []SeederDetails
	NumChunks        int
	ChunksDownloaded int
	ChunksShared     int
	BytesDownloaded  int64
	BytesUploaded    int64
	DateTimeAdded    int64
}

type SeederDetails struct {
	PublicKey     string
	Workers       int
	ActiveSession bool
	LastActivity  int64
}

func (s *MiddlewareFunctions) GetFileDetails(FileHash string) FileDetails {

	file, err := dbGetFile(FileHash)
	if err != nil {
		pushError("Error on getting file details", err.Error())
		return FileDetails{}
	}

	chunksDownloaded := chunksDownloaded(file.ChunkMap, file.NumChunks)
	byteDown := int64(chunksDownloaded) * int64(constants.ChunkSize)
	byteUp := int64(file.ChunksShared) * int64(constants.ChunkSize)

	seederDetails := []SeederDetails{}
	seeders := GetSeeders(FileHash)

	for _, v := range seeders {

		mutexes.WorkerMapLock.Lock()
		workerCount := workerMap[v]
		mutexes.WorkerMapLock.Unlock()

		session, exists := sessionmanager.GetExistingSessionWithoutClosing(v, constants.GetSessionDialTimeout)

		sessionActive := false
		lastActivity := int64(-1)

		if exists {
			sessionActive = true
			lastActivity = session.LastActivityUnix
		}

		seederDetails = append(seederDetails, SeederDetails{
			PublicKey:     v,
			Workers:       workerCount,
			ActiveSession: sessionActive,
			LastActivity:  lastActivity,
		})
	}

	return FileDetails{
		FileID:           file.FileHash,
		Seeders:          seederDetails,
		NumChunks:        file.NumChunks,
		ChunksDownloaded: chunksDownloaded,
		ChunksShared:     file.ChunksShared,
		BytesDownloaded:  byteDown,
		BytesUploaded:    byteUp,
		DateTimeAdded:    file.DateTimeAdded,
	}
}
func (s *MiddlewareFunctions) GetTopicDetails(Topic string) models.TopicInfo {
	return GetTopicInfo(Topic)
}

func (s *MiddlewareFunctions) GetOfficialTopicName() string {
	return constants.SurgeOfficialTopic
}

func (s *MiddlewareFunctions) SetDownloadFolder() bool {
	path, _ := runtime.OpenDirectoryDialog(*wailsContext, runtime.OpenDialogOptions{
		Title: "Select Download Folder",
	})
	if path == "" {
		return false
	}
	DbWriteSetting("downloadFolder", path)
	return true
}

func (s *MiddlewareFunctions) GetWalletAddress() string {
	return WalletAddress()
}

func (s *MiddlewareFunctions) GetWalletBalance() string {
	return WalletBalance()
}
func TransferToPk(PubKey string, Amount string, Fee string) string {
	walletAddr, _ := nkn.ClientAddrToWalletAddr(PubKey)
	_, hash := WalletTransfer(walletAddr, Amount, Fee)
	return hash
}

func (s *MiddlewareFunctions) GetTxFee() string {
	return TransactionFee
}

func (s *MiddlewareFunctions) SetTxFee(Fee string) {
	fmt.Println("tx fee set", Fee)
	TransactionFee = Fee
	DbWriteSetting("defaultTxFee", Fee)
}

func (s *MiddlewareFunctions) Tip(FileHash string, Amount string, Fee string) {
	fmt.Println(FileHash, Amount, Fee)
	amountFloat, err := strconv.ParseFloat(Amount, 64)

	if err != nil {
		pushError("Error on tip", "Invalid amount.")
		return
	}

	if amountFloat < 0.00000001 {
		pushError("Error on tip", "Minimum tip amount is 0.00000001")
		return
	}

	seeders := GetSeeders(FileHash)
	if len(seeders) == 0 {
		pushError("Error on tip", "No seeders found to tip.")
		return
	}

	share := amountFloat / float64(len(GetSeeders(FileHash)))
	calculatedFee := CalculateFee(Fee)

	feePerTxFloat, _ := strconv.ParseFloat(calculatedFee, 64)
	totalFeeFloat := feePerTxFloat * float64(len(seeders))

	isEnough, balanceError := ValidateBalanceForTransaction(amountFloat, totalFeeFloat)
	if !isEnough {
		pushError("Error on tip", balanceError.Error())
		return
	}

	for _, v := range seeders {
		walletAddr, _ := nkn.ClientAddrToWalletAddr(v)
		success, hash := WalletTransfer(walletAddr, fmt.Sprintf("%f", share), calculatedFee)
		fmt.Println(success, hash, walletAddr, share)

		if !success {
			break
		}
	}
}

package surge

import (
	"time"

	nkn "github.com/nknorg/nkn-sdk-go"
	"github.com/rule110-io/surge-ui/surge/platform"
	log "github.com/sirupsen/logrus"
)

//ChunkSize is size of chunk in bytes (256 kB)
const ChunkSize = 1024 * 256

//NumClients is the number of NKN clients
const NumClients = 8

//NumWorkers is the total number of concurrent chunk fetches allowed
const NumWorkers = 32

//whether the nkn client is initialized
var clientInitialized = false

//The nkn client
var client *nkn.MultiClient

//duration of a subscription blocktime is ~20sec
const subscriptionDuration = 180 // 180 is approximately one hour

//InitializeClient initializes connection with nkn
func InitializeClient() {
	var err error

	account := InitializeAccount()
	client, err = nkn.NewMultiClient(account, "", NumClients, false, nil)
	clientInitialized = true

	if err != nil {
		pushError("Error on startup", err.Error())
	} else {
		<-client.OnConnect.C

		pushNotification("Client Connected", "Successfully connected to the NKN network")

		client.Listen(nil)
		SurgeActive = true
		go Listen()

		dbFiles := dbGetAllFiles()
		var filesOnDisk []File

		for i := 0; i < len(dbFiles); i++ {
			if FileExists(dbFiles[i].Path) {
				filesOnDisk = append(filesOnDisk, dbFiles[i])
			} else {
				dbFiles[i].IsMissing = true
				dbFiles[i].IsDownloading = false
				dbFiles[i].IsUploading = false
				dbInsertFile(dbFiles[i])
			}
		}

		go BuildSeedString(filesOnDisk)
		for i := 0; i < len(filesOnDisk); i++ {
			go restartDownload(filesOnDisk[i].FileHash)
		}

		go autoSubscribeWorker()
		go GetSubscriptions()

		go queryRemoteForFiles()

		go platform.WatchOSXHandler()

		go rescanPeers()
	}
}

//Stop cleanup for surge
func Stop() {
	client.Close()
	client = nil
}

//Function that automatically grabs subscriptions for nkn topic
func rescanPeers() {
	defer RecoverAndLog()
	for true {
		time.Sleep(time.Minute)
		go GetSubscriptions()
	}
}

//GetNumberOfRemoteClient returns number of clients and online clients
func GetNumberOfRemoteClient() (int, int) {
	return numClientsSubscribed, numClientsOnline
}

func autoSubscribeWorker() {
	defer RecoverAndLog()

	//As long as the client is running subscribe
	for true {
		//Only subscribe when this client is hosting anything
		hosting := false

		files := dbGetAllFiles()
		for i := 0; i < len(files); i++ {
			if files[i].IsUploading {
				hosting = true
				break
			}
		}

		if hosting {
			subscribeToSurgeTopic()
		}

		time.Sleep(time.Second * 20 * subscriptionDuration)
	}
}

func subscribeToSurgeTopic() {
	Topic := TopicEncode(TestTopic)
	txnHash, err := client.Subscribe("", Topic, subscriptionDuration, "Surge Beta Client", nil)
	if err != nil {
		log.Println("Probably already subscribed", err)
	} else {
		log.Println("Subscribed: ", txnHash)
	}
}

//GetSubscriptions .
func GetSubscriptions() {
	defer RecoverAndLog()

	Topic := TopicEncode(TestTopic)

	subResponse, err := client.GetSubscribers(Topic, 0, 100, true, true)
	if err != nil {
		pushError("Error on get subscriptions", err.Error())
		return
	}

	for k, v := range subResponse.SubscribersInTxPool.Map {
		subResponse.Subscribers.Map[k] = v
	}

	subscribers = []string{}
	for k, v := range subResponse.Subscribers.Map {
		if len(v) > 0 {
			if k != client.Addr().String() {
				subscribers = append(subscribers, k)
			}
		}
	}
}

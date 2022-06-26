package surge

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/rule110-io/surge/backend/constants"
	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/mutexes"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var topicsMap map[string]models.Topic
var topicEncodedSubcribeStateMap map[string]int
var startupSubscribe = true

const topicsMapBucketKey = "topicBucket"

func InitializeTopicsManager() {
	topicsMap = make(map[string]models.Topic)
	topicEncodedSubcribeStateMap = make(map[string]int)

	//Load from db
	mapString, err := DbReadSetting(topicsMapBucketKey)
	if err == nil {
		mapBytes := []byte(mapString)
		err := json.Unmarshal(mapBytes, &topicsMap)
		if err != nil {
			log.Println("Failed to unmarshal setting for topics", err)
		}
	}
}

func subscribeToSurgeTopic(topicName string, applySafeLock bool) {

	if applySafeLock {
		mutexes.TopicsMapLock.Lock()
		defer mutexes.TopicsMapLock.Unlock()
	}

	topicEncoded := TopicEncode(topicName)
	subscriptionActive := IsSubscriptionActive(topicEncoded)

	if _, ok := topicsMap[topicName]; !ok {
		topicModel := models.Topic{
			Name:        topicName,
			NameEncoded: topicEncoded,
		}

		topicsMap[topicName] = topicModel

		//Save to our bucket
		mapBytes, err := json.Marshal(topicsMap)
		if err == nil {
			mapString := string(mapBytes)
			DbWriteSetting(topicsMapBucketKey, mapString)
		}
	}

	//If we dont have an active sub resubscribe.
	if !subscriptionActive {
		subscribeToPubSub(topicEncoded)
	}

	//Only announce files if the client is first starting up, or when we are newly subscribed.
	if startupSubscribe || !subscriptionActive {
		AnnounceFiles(topicEncoded)
	}
}

func unsubscribeFromSurgeTopic(topicName string) {
	mutexes.TopicsMapLock.Lock()
	defer mutexes.TopicsMapLock.Unlock()

	if topic, ok := topicsMap[topicName]; ok {
		unsubscribeToPubSub(topic.NameEncoded)
	}

	//Delete from map
	delete(topicsMap, topicName)

	//Save to our bucket
	mapBytes, err := json.Marshal(topicsMap)
	if err == nil {
		mapString := string(mapBytes)
		DbWriteSetting(topicsMapBucketKey, mapString)
	}
}

func resubscribeToTopics() {
	mutexes.TopicsMapLock.Lock()
	defer mutexes.TopicsMapLock.Unlock()
	for _, topic := range topicsMap {

		subscribeToSurgeTopic(topic.Name, false)
	}
	startupSubscribe = false
}

//GetTopicInfo returns info about the topic given
func GetTopicInfo(topicName string) models.TopicInfo {

	topicEncoded := TopicEncode(topicName)
	subCount, _ := client.GetSubscribersCount(topicEncoded)

	//count files with topic
	fileCount := 0
	for _, v := range ListedFiles {
		bitSetVar := 0
		if v.Topic == topicName {
			bitSetVar = 1
		}
		fileCount += bitSetVar
	}

	state := 0

	//get topic state
	knownState, any := topicEncodedSubcribeStateMap[topicEncoded]
	if any {
		state = knownState
	}

	return models.TopicInfo{
		Name:              topicName,
		Subscribers:       subCount,
		FileCount:         fileCount,
		Permissions:       GetTopicPermissions(topicName, GetAccountAddress()),
		SubscriptionState: state,
	}
}
func GetTopicPermissions(topicName string, clientAddr string) models.TopicPermissions {
	if topicName != constants.SurgeOfficialTopic {
		return models.TopicPermissions{
			CanRead:  true,
			CanWrite: true,
		}
	}

	//Check if user is from team
	if clientAddr == constants.TeamAddressA ||
		clientAddr == constants.TeamAddressB ||
		clientAddr == constants.TeamAddressC {
		return models.TopicPermissions{
			CanRead:  true,
			CanWrite: true,
		}
	} else {
		return models.TopicPermissions{
			CanRead:  true,
			CanWrite: false,
		}
	}
}

func GetTopicsWithPermissions() []models.TopicInfo {
	topicNames := []string{}

	mutexes.TopicsMapLock.Lock()
	for _, v := range topicsMap {
		topicNames = append(topicNames, v.Name)
	}
	mutexes.TopicsMapLock.Unlock()
	sort.Strings(topicNames)

	//Get objects from names
	modelData := []models.TopicInfo{}

	for _, v := range topicNames {
		entry := models.TopicInfo{
			Name:        v,
			Permissions: GetTopicPermissions(v, GetAccountAddress()),
		}
		modelData = append(modelData, entry)
	}

	return modelData
}

func updateTopicSubscriptionState(TopicEncoded string, NewState int) {
	topicEncodedSubcribeStateMap[TopicEncoded] = NewState
	runtime.EventsEmit(*wailsContext, "topicsUpdated")

	fmt.Println("Subscription state updated", TopicEncoded, NewState)
}

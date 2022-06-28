package surge

import (
	"encoding/json"
	"log"
	"sort"

	"github.com/rule110-io/surge/backend/constants"
	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/mutexes"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var topicsMap map[string]models.Topic
var topicEncodedSubcribeStateMap map[string]int

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

func subscribeToSurgeTopic(topicName string, applySafeLock bool) (bool, error) {

	if applySafeLock {
		mutexes.TopicsMapLock.Lock()
		defer mutexes.TopicsMapLock.Unlock()
	}

	topicEncoded := TopicEncode(topicName)
	subscriptionActive, err := IsSubscriptionActive(topicEncoded)
	if err != nil {
		return false, err
	}

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

	previousState := topicEncodedSubcribeStateMap[topicEncoded]

	subscribeSuccess := true
	//If we dont have an active sub resubscribe.
	if !subscriptionActive {
		if !subscribeToPubSub(topicEncoded) {
			subscribeSuccess = false
		}
	} else {
		updateTopicSubscriptionState(topicEncoded, 2)
	}

	//Only announce files if the client is first starting up, or when we are newly subscribed.
	if subscribeSuccess {

		if previousState == 0 || previousState == 1 {
			AnnounceFiles(topicEncoded)
		}
		//first startup, were already subscribed, set the state.
		updateTopicSubscriptionState(topicEncoded, 2)
	}

	return subscribeSuccess, nil
}

func unsubscribeFromSurgeTopic(topicName string) bool {
	mutexes.TopicsMapLock.Lock()
	defer mutexes.TopicsMapLock.Unlock()

	if topic, ok := topicsMap[topicName]; ok {
		unsubSuccess := unsubscribeToPubSub(topic.NameEncoded)
		if !unsubSuccess {
			return false
		}
	}

	//Delete from map
	delete(topicsMap, topicName)

	//Save to our bucket
	mapBytes, err := json.Marshal(topicsMap)
	if err == nil {
		mapString := string(mapBytes)
		DbWriteSetting(topicsMapBucketKey, mapString)
	}

	return true
}

func resubscribeToTopics() {
	mutexes.TopicsMapLock.Lock()
	defer mutexes.TopicsMapLock.Unlock()

	for _, topic := range topicsMap {
		_, err := subscribeToSurgeTopic(topic.Name, false)
		if err != nil {

			//set all topics to pending state
			pushError("Topic connection error", err.Error())
			for _, disableTopic := range topicsMap {
				updateTopicSubscriptionState(disableTopic.NameEncoded, 1)
			}

			break
		}
	}
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

		state := 0
		//get topic state
		knownState, any := topicEncodedSubcribeStateMap[TopicEncode(v)]
		if any {
			state = knownState
		}

		entry := models.TopicInfo{
			Name:              v,
			Permissions:       GetTopicPermissions(v, GetAccountAddress()),
			SubscriptionState: state,
		}
		modelData = append(modelData, entry)
	}

	return modelData
}

func updateTopicSubscriptionState(TopicEncoded string, NewState int) {

	previousValue, exists := topicEncodedSubcribeStateMap[TopicEncoded]
	isChanged := !exists || previousValue != NewState

	if isChanged {
		topicEncodedSubcribeStateMap[TopicEncoded] = NewState
		if FrontendReady {
			runtime.EventsEmit(*wailsContext, "topicsUpdated")
		}
	}
}

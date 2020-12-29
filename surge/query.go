package surge

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	pb "github.com/rule110-io/surge-ui/payloads"
	"github.com/rule110-io/surge-ui/surge/constants"
	"github.com/rule110-io/surge-ui/surge/sessionmanager"
	"google.golang.org/protobuf/proto"
)

var queryPayload = ""

// SendQueryRequest sends a query to a client on session
func SendQueryRequest(Addr string, Query string) bool {

	surgeSession, exists := sessionmanager.GetExistingSession(Addr, constants.SendQueryRequestSessionTimeout, "Send query request timeout - SendQueryRequestSessionTimeout")

	if !exists {
		return false
	}

	msg := &pb.SurgeQuery{
		Query: Query,
	}
	msgSerialized, err := proto.Marshal(msg)
	if err != nil {
		log.Panic("Failed to encode surge message:", err)
		return false
	}

	fmt.Println(string("\033[31m"), "Send Query Request", Addr, string("\033[0m"))
	written, err := SessionWrite(surgeSession, msgSerialized, constants.SurgeQueryRequestID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
	if err != nil {
		log.Println("Failed to send Surge Request:", err)
		return false
	}

	//Write add to upload
	bandwidthAccumulatorMapLock.Lock()
	uploadBandwidthAccumulator["DISCOVERY"] += written
	bandwidthAccumulatorMapLock.Unlock()

	return true
}

// SendQueryResponse sends a query to a client on session
func SendQueryResponse(Session *sessionmanager.Session, Query string) {

	b := []byte(queryPayload)
	fmt.Println(string("\033[31m"), "Send Query Response", Session.Session.RemoteAddr().String(), string("\033[0m"))
	written, err := SessionWrite(Session, b, constants.SurgeQueryResponseID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
	if err != nil {
		log.Println("Failed to send Surge Ruquest:", err)
	}
	//Write add to upload
	bandwidthAccumulatorMapLock.Lock()
	uploadBandwidthAccumulator["DISCOVERY"] += written
	bandwidthAccumulatorMapLock.Unlock()
}

func processQueryRequest(Session *sessionmanager.Session, Data []byte) {

	//Try to parse SurgeMessage
	surgeQuery := &pb.SurgeQuery{}
	if err := proto.Unmarshal(Data, surgeQuery); err != nil {
		log.Panic("Failed to parse surge message:", err)
	}
	log.Println("Query received", surgeQuery.Query)

	SendQueryResponse(Session, surgeQuery.Query)
}

func processQueryResponse(Session *sessionmanager.Session, Data []byte) {

	//Try to parse SurgeMessage
	s := string(Data)
	seeder := Session.Session.RemoteAddr().String()

	fmt.Println(string("\033[36m"), "file query response received", seeder, string("\033[0m"))

	ListedFilesLock.Lock()

	//Parse the response
	payloadSplit := strings.Split(s, "surge://")
	for j := 0; j < len(payloadSplit); j++ {
		data := strings.Split(payloadSplit[j], "|")

		if len(data) < 3 {
			continue
		}

		fileSize, _ := strconv.ParseInt(data[3], 10, 64)
		numChunks := int((fileSize-1)/int64(ChunkSize)) + 1

		newListing := File{
			FileName:    data[2],
			FileSize:    fileSize,
			FileHash:    data[4],
			seeders:     []string{seeder},
			Path:        "",
			NumChunks:   numChunks,
			ChunkMap:    nil,
			seederCount: 1,
		}

		//Replace existing, or remove.
		var replace = false
		for l := 0; l < len(ListedFiles); l++ {
			if ListedFiles[l].FileHash == newListing.FileHash {

				//if the seeder is unique add it as an additional seeder for the file
				ListedFiles[l].seeders = append(ListedFiles[l].seeders, seeder)
				ListedFiles[l].seeders = distinctStringSlice(ListedFiles[l].seeders)
				ListedFiles[l].seederCount = len(ListedFiles[l].seeders)

				replace = true
				break
			}
		}
		//Unique listing so we add
		if replace == false {
			ListedFiles = append(ListedFiles, newListing)
		}

		fmt.Println(string("\033[33m"), "Filename", newListing.FileName, "FileHash", newListing.FileHash, string("\033[0m"))

		log.Println("Query response new file: ", newListing.FileName, " seeder: ", seeder)
	}
	ListedFilesLock.Unlock()
}

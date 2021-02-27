package surge

import (
	"math/rand"
	"strings"
	"time"

	nkn "github.com/nknorg/nkn-sdk-go"
	"github.com/rule110-io/surge/backend/constants"
)

//PersistRPC will persist all current nkn RPC connections for future bootstrapping the client
func PersistRPC(client *nkn.MultiClient) {
	connectedClients := client.GetClients()

	//We dont persist our connection if count is low, if our connections are scarce better to fall back to default connections
	if len(connectedClients) > 3 {
		rpcAddresses := []string{}
		for _, conClient := range connectedClients {
			rAddr := "http://" + conClient.GetNode().RPCAddr
			rpcAddresses = append(rpcAddresses, rAddr)
		}

		//Persist the addresses in our db
		rpcAddrString := strings.Join(rpcAddresses[:], ",")
		DbWriteSetting("rpcCache", rpcAddrString)
	}
}

//GetBootstrapRPC returns the rpc collection to connect to nkn
func GetBootstrapRPC() *nkn.StringArray {
	rpcAddrString, err := DbReadSetting("rpcCache")

	rpcAddresses := []string{}
	if err == nil {
		rpcAddresses = strings.Split(rpcAddrString, ",")

		//Shuffle the addresses
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(rpcAddresses), func(i, j int) { rpcAddresses[i], rpcAddresses[j] = rpcAddresses[j], rpcAddresses[i] })

		//Append fallback default rpc
		rpcAddresses = append(rpcAddresses, constants.DefaultRPCAddress)
		return nkn.NewStringArray(rpcAddresses...)
	}
	return nil
}

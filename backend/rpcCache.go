package surge

import (
	"context"
	"strings"
	"sync"

	nkn "github.com/nknorg/nkn-sdk-go"
	nknPb "github.com/nknorg/nkn/v2/pb"
	"github.com/rule110-io/surge/backend/constants"
)

const (
	rpcTimeout = 3000 // in millisecond
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
	var rpcAddresses []string

	rpcAddrString, err := DbReadSetting("rpcCache")
	if err == nil {
		rpcAddresses = strings.Split(rpcAddrString, ",")

		// Find the available ones and sort by latency
		filteredRPCAddresses, err := FilterSeedRPCServer(context.Background(), rpcAddresses, rpcTimeout)
		if err == nil {
			rpcAddresses = filteredRPCAddresses
		}
	}

	//Append fallback default rpc
	rpcAddresses = append(rpcAddresses, constants.DefaultRPCAddress)

	return nkn.NewStringArray(rpcAddresses...)
}

// FilterSeedRPCServer gets the state of rpc node list, remove unavailable nodes
// (network failure or not in persist finish state), and sort available ones by
// latency.
func FilterSeedRPCServer(ctx context.Context, maybeRPCAddrs []string, timeout int32) ([]string, error) {
	var wg sync.WaitGroup
	var lock sync.Mutex
	rpcAddrs := make([]string, 0, len(maybeRPCAddrs))

	for _, addr := range maybeRPCAddrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			nodeState, err := nkn.GetNodeStateContext(ctx, &nkn.RPCConfig{
				SeedRPCServerAddr: nkn.NewStringArray(addr),
				RPCTimeout:        timeout,
			})
			if err != nil {
				return
			}
			if nodeState.SyncState != nknPb.SyncState_name[int32(nknPb.SyncState_PERSIST_FINISHED)] {
				return
			}
			lock.Lock()
			rpcAddrs = append(rpcAddrs, addr)
			lock.Unlock()
		}(addr)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
	}

	return rpcAddrs, nil
}

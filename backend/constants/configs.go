package constants

const (
	//ChunkSize is size of chunk in bytes (256 kB)
	ChunkSize = 1024 * 256

	//NumClients is the number of NKN clients
	NumClients = 8

	//NumWorkers is the total number of concurrent chunk fetches allowed
	NumWorkers = 32

	//duration of a subscription blocktime is ~20sec
	SubscriptionDuration = 180
)

package surge

import (
	bitmap "github.com/boljen/go-bitmap"
	movavg "github.com/mxmCherry/movavg"
	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/mutexes"
)

var downloadBandwidthAccumulator map[string]int
var uploadBandwidthAccumulator map[string]int

var fileBandwidthMap map[string]models.BandwidthMA

var zeroBandwidthMap map[string]bool

var chunksInTransit map[string]bool

//sets the current bandwith of a file
func fileBandwidth(FileID string) (Download int, Upload int) {

	//Get accumulator
	mutexes.BandwidthAccumulatorMapLock.Lock()
	downAccu := downloadBandwidthAccumulator[FileID]
	downloadBandwidthAccumulator[FileID] = 0

	upAccu := uploadBandwidthAccumulator[FileID]
	uploadBandwidthAccumulator[FileID] = 0
	mutexes.BandwidthAccumulatorMapLock.Unlock()

	if fileBandwidthMap[FileID].Download == nil {
		fileBandwidthMap[FileID] = models.BandwidthMA{
			Download: movavg.ThreadSafe(movavg.NewSMA(10)),
			Upload:   movavg.ThreadSafe(movavg.NewSMA(10)),
		}
	}

	fileBandwidthMap[FileID].Download.Add(float64(downAccu))
	fileBandwidthMap[FileID].Upload.Add(float64(upAccu))

	return int(fileBandwidthMap[FileID].Download.Avg()), int(fileBandwidthMap[FileID].Upload.Avg())
}

func chunksDownloaded(s []byte, num int) int {
	//No chunkmap means no download was initiated, all chunks are local
	if s == nil {
		return num
	}

	chunksLocalNum := 0
	for i := 0; i < num; i++ {
		if bitmap.Get(s, i) == true {
			chunksLocalNum++
		}
	}
	return chunksLocalNum
}

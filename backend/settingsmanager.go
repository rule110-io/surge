package surge

import (
	"strconv"

	"github.com/rule110-io/surge/backend/constants"
	"github.com/rule110-io/surge/backend/platform"
)

func getDownloadFolderPath() (string, error) {
	folder, err := DbReadSetting("downloadFolder")
	if err == nil && len(folder) > 0 {
		return folder, nil
	} else {
		folder, err = platform.GetRemoteFolder()
		if err == nil {
			return folder, nil
		}
	}
	return "", err
}

func clamp(val int, min int, max int) int {
	if val > max {
		return max
	} else if val < min {
		return min
	}
	return val
}

func getNumberClients() int {
	num, err := DbReadSetting("numClients")
	if err == nil && len(num) > 0 {
		val, _ := strconv.Atoi(num)
		return clamp(val, constants.NumClientsMin, constants.NumClientsMax)
	} else {
		return constants.NumClients
	}
}
func getNumberWorkers() int {
	num, err := DbReadSetting("numWorkers")
	if err == nil && len(num) > 0 {
		val, _ := strconv.Atoi(num)
		return clamp(val, constants.NumWorkersMin, constants.NumWorkersMax)
	} else {
		return constants.NumWorkers
	}
}

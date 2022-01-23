package platform

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"log"

	"github.com/rule110-io/surge/mailslot"
	"github.com/sqweek/dialog"
)

const instanceid = "0a7332a0-8ddf-450e-9ff4-306cdca3278e"
const mailSlotName = `\\.\mailslot\surgeclient`

//ProcessStartupArgs handles startup args logic
func ProcessStartupArgs(args []string, frontendReady *bool) bool {
	lastArg := args[len(args)-1]

	singleInstanceErr := mailslot.SingleInstance(instanceid)
	isSingleInstance := singleInstanceErr == nil

	if isSingleInstance {
		go mailSlotListen()
	}

	magnetString := ""

	//Check if param is magnet
	magnetFound := strings.Contains(lastArg, "surge://|file|")
	if magnetFound {
		magnetString = lastArg
	}

	//Check if param is filepath to src
	surgeFileFound := strings.Contains(lastArg, ".slc")

	if surgeFileFound {
		_, err := os.Stat(lastArg)

		// If the file doesn't exist, create it
		if os.IsNotExist(err) {
			dialog.Message("%s", "Failed to read surge file from \r\n"+lastArg).Title("Surge Startup Failed").Error()
		} else { //else read seed from file
			file, err := os.Open(lastArg)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			magnetStringBytes, err := ioutil.ReadAll(file)
			magnetString = string(magnetStringBytes)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	//In the case of no arguments we either exit when surge is running or startup main client
	if magnetString == "" {
		if !isSingleInstance {
			dialog.Message("%s", "Another instance of surge is already running\r\n"+strings.Join(args, ",")).Title("Surge Startup Failed").Error()
			return false
		}
		return true
	}

	//In case of arguments we will either notify our main client, or if that is not running become the mainclient.
	if !isSingleInstance {
		mailSlotSend(magnetString)
		return false
	}

	fireStartDownload := func() {
		//If wails frontend is not yet binded, we wait in a task to not block main thread
		if !*frontendReady {
			waitAndPush := func() {
				for !*frontendReady {
					time.Sleep(time.Millisecond * 50)
				}
				AskUser("startDownloadMagnetLinks", magnetString)
			}
			go waitAndPush()
		} else {
			AskUser("startDownloadMagnetLinks", magnetString)
		}
	}
	go fireStartDownload()

	return true
}

func mailSlotListen() {
	ms, err := mailslot.New(mailSlotName, 0, mailslot.MAILSLOT_WAIT_FOREVER)
	if err != nil {
		panic(err)
	}

	defer ms.Close()

	reader := bufio.NewReader(ms)

	for {
		testData, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Println(testData)

		//Insert new file from arguments and start download
		AskUser("startDownloadMagnetLinks", testData)
	}
}

func mailSlotSend(message string) {
	ms, err := mailslot.Open(mailSlotName)
	if err != nil {
		panic(err)
	}

	defer ms.Close()

	ms.Write([]byte(message + "\n"))
}

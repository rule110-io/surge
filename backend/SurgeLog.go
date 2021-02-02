package surge

import (
	"log"
	"os"

	"github.com/rule110-io/surge/backend/platform"
	"gopkg.in/natefinch/lumberjack.v2"
)

const logPath = "surge.log"

//InitializeLog init for log file
func InitializeLog() {

	//var err error

	var dir = platform.GetSurgeDir()

	var logPathOS = dir + string(os.PathSeparator) + logPath

	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	//file, err := os.OpenFile(logPathOS, os.O_WRONLY|os.O_CREATE, 0644)
	//if err != nil {
	//	log.Panic(err)
	//}

	log.SetOutput(&lumberjack.Logger{
		Filename:   logPathOS,
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     10,   //days
		Compress:   true, // disabled by default
	})

	// Only log the warning severity or above.
	//log.SetLevel(log.WarnLevel)
}

//OpenLogFile opens a log file with OS default application for object type
func OpenLogFile() {
	var err error

	var dir = platform.GetSurgeDir()

	var logPathOS = dir + string(os.PathSeparator) + logPath

	if err != nil {
		pushError("Error on open log", err.Error())
		return
	}

	OpenOSPath(logPathOS)
}

package surge

import (
	"os"
	"os/user"
	"runtime"

	"golang.org/x/sys/windows"
)

const remotePath = "downloads"

var remoteFolder = ""

//InitializeFolders initializes folder structures needed, returns if folders are created (not yet existing)
func InitializeFolders() bool {
	newCreated := false

	var dirFileMode os.FileMode
	var dir = GetSurgeDir()
	dirFileMode = os.ModeDir | (osUserRwx | osAllR)

	myself, err := user.Current()
	if err != nil {
		pushError("Error on startup", err.Error())
	}

	if runtime.GOOS == "windows" {
		homedir, _ := windows.KnownFolderPath(windows.FOLDERID_Downloads, 0)
		remoteFolder = homedir + string(os.PathSeparator) + "surge_" + remotePath
	} else {
		homedir := myself.HomeDir
		remoteFolder = homedir + string(os.PathSeparator) + "Downloads" + string(os.PathSeparator) + "surge_" + remotePath
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, dirFileMode)
		newCreated = true
	}

	//Ensure remote folders exist
	if _, err := os.Stat(remoteFolder); os.IsNotExist(err) {
		os.Mkdir(remoteFolder, dirFileMode)
	}

	return newCreated
}

package platform

import (
	"os"

	"golang.org/x/sys/windows"
)

//InitializeFolders initializes folder structures needed, returns if folders are created (not yet existing)
func InitializeFolders() (bool, error) {
	newCreated := false

	var dir = GetSurgeDir()
	var dirFileMode os.FileMode
	dirFileMode = os.ModeDir | (osUserRwx | osAllR)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, dirFileMode)
		newCreated = true
	}

	remoteFolder, err := GetRemoteFolder()
	if err != nil {
		return false, err
	}

	//Ensure remote folders exist
	if _, err := os.Stat(remoteFolder); os.IsNotExist(err) {
		os.Mkdir(remoteFolder, dirFileMode)
	}

	return newCreated, nil
}

//GetSurgeDir returns the surge dir
func GetSurgeDir() string {
	return os.Getenv("APPDATA") + string(os.PathSeparator) + "Surge"
}

//GetRemoteFolder returns the download dir
func GetRemoteFolder() (string, error) {
	homedir, _ := windows.KnownFolderPath(windows.FOLDERID_Downloads, 0)
	return homedir + string(os.PathSeparator) + "surge_downloads", nil
}

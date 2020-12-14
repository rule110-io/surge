package platform

import (
	"os"

	"golang.org/x/sys/windows"
)

//GetSurgeDir returns the surge dir
func GetSurgeDir() string {
	return os.Getenv("APPDATA") + string(os.PathSeparator) + "Surge"
}

//GetRemoteFolder returns the download dir
func GetRemoteFolder() (string, error) {
	homedir, _ := windows.KnownFolderPath(windows.FOLDERID_Downloads, 0)
	return homedir + string(os.PathSeparator) + "surge_downloads", nil
}

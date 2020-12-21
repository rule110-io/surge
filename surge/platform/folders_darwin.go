package platform

import (
	"os"
	"os/user"
)

//GetSurgeDir returns the surge dir
func GetSurgeDir() string {
	return os.Getenv("HOME") + string(os.PathSeparator) + ".surge"
}

//GetRemoteFolder returns the download dir
func GetRemoteFolder() (string, error) {
	myself, err := user.Current()
	if err != nil {
		return "", err
	}
	homedir := myself.HomeDir
	return homedir + string(os.PathSeparator) + "Downloads" + string(os.PathSeparator) + "surge_downloads", nil

}

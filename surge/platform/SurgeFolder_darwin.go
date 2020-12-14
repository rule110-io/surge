package platform

import (
	"os"
	"os/user"
)

//OS folder permission bitflags
const (
	osRead       = 04
	osWrite      = 02
	osEx         = 01
	osUserShift  = 6
	osGroupShift = 3
	osOthShift   = 0

	osUserR   = osRead << osUserShift
	osUserW   = osWrite << osUserShift
	osUserX   = osEx << osUserShift
	osUserRw  = osUserR | osUserW
	osUserRwx = osUserRw | osUserX

	osGroupR   = osRead << osGroupShift
	osGroupW   = osWrite << osGroupShift
	osGroupX   = osEx << osGroupShift
	osGroupRw  = osGroupR | osGroupW
	osGroupRwx = osGroupRw | osGroupX

	osOthR   = osRead << osOthShift
	osOthW   = osWrite << osOthShift
	osOthX   = osEx << osOthShift
	osOthRw  = osOthR | osOthW
	osOthRwx = osOthRw | osOthX

	osAllR   = osUserR | osGroupR | osOthR
	osAllW   = osUserW | osGroupW | osOthW
	osAllX   = osUserX | osGroupX | osOthX
	osAllRw  = osAllR | osAllW
	osAllRwx = osAllRw | osGroupX
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

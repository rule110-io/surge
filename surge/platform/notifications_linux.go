package platform

import (
	"os"
	"os/exec"
	"path/filepath"

	"log"
)

// ShowNotification .
func ShowNotification(title string, text string) {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panic(err)
	}

	Icon := filepath.Join(dir, "appicon.png")

	exec.Command("notify-send", "-i", Icon, title, text).Run()
}

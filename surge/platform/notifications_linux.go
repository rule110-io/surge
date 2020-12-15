package platform

import (
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// ShowNotification .
func ShowNotification(title string, text string) error {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	Icon := filepath.Join(dir, "appicon.png")

	exec.Command("notify-send", "-i", Icon, title, text).Run()
}

package surge

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/toast.v1"
)

func showNotification(title string, text string) {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	Icon := filepath.Join(dir, "appicon.png")

	notification := toast.Notification{
		AppID:   "Surge",
		Title:   title,
		Message: text,
		Icon:    Icon, // This file must exist (remove this line if it doesn't)
	}
	err = notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}

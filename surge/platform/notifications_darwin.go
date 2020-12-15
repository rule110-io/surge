package platform

import (
	"os"
	"path/filepath"

	notifier "github.com/deckarep/gosx-notifier"
	log "github.com/sirupsen/logrus"
)

// ShowNotification .
func ShowNotification(title string, text string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	Icon := filepath.Join(dir, "appicon.png")

	notification := notifier.Notification{
		Group:   "com.rule110.surge",
		Title:   title,
		Message: text,
		Sound:   notifier.Glass,
		AppIcon: Icon,
	}

	notification.Push()
}

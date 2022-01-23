package platform

import (
	"os"
	"path/filepath"

	"log"

	notifier "github.com/deckarep/gosx-notifier"
)

// ShowNotification .
func ShowNotification(title string, text string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panic(err)
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

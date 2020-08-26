package surge

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/toast.v1"
)

func watchOSXHandler() {
}

func initOSHandler() {

}

func setVisualModeLikeOS() {

}

func showNotification(title string, text string) {
	notification := toast.Notification{
		AppID:   "Surge",
		Title:   title,
		Message: text,
	}
	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}

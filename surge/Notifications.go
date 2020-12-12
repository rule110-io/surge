package surge

import "time"

func pushNotification(title string, text string) {
	//If wails frontend is not yet binded, we wait in a task to not block main thread
	if wailsRuntime == nil {

		waitAndPush := func() {
			for wailsRuntime == nil {
				time.Sleep(50)
			}
			wailsRuntime.Events.Emit("notificationEvent", title, text)
		}
		go waitAndPush()
	} else {
		wailsRuntime.Events.Emit("notificationEvent", title, text)
	}
}

func pushError(title string, text string) {
	//If wails frontend is not yet binded, we wait in a task to not block main thread
	if wailsRuntime == nil {
		waitAndPush := func() {
			for wailsRuntime == nil {
				time.Sleep(50)
			}
			wailsRuntime.Events.Emit("errorEvent", title, text)
		}
		go waitAndPush()
	} else {
		wailsRuntime.Events.Emit("errorEvent", title, text)
	}
}

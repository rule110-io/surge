package surge

import "time"

func pushNotification(title string, text string) {
	//If wails frontend is not yet binded, we wait in a task to not block main thread
	if !FrontendReady {
		waitAndPush := func() {
			for !FrontendReady {
				time.Sleep(time.Millisecond * 50)
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
	if !FrontendReady {
		waitAndPush := func() {
			for !FrontendReady {
				time.Sleep(time.Millisecond * 50)
			}
			wailsRuntime.Events.Emit("errorEvent", title, text)
		}
		go waitAndPush()
	} else {
		wailsRuntime.Events.Emit("errorEvent", title, text)
	}
}

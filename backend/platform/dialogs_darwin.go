package platform

import (
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//OpenFileDialog uses platform agnostic package for a file dialog
func OpenFileDialog() string {
	selectedFile, err := runtime.OpenFileDialog(*wailsContext, runtime.OpenDialogOptions{
		Title: "Select File",
	})
	if err != nil {
		log.Panic("Error on file opening", err.Error())
	}
	return selectedFile
}

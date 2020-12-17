package platform

import "github.com/sqweek/dialog"

//OpenFileDialog uses platform agnostic package for a file dialog
func OpenFileDialog() string {
	selectedFile := dialog.File().Load()
	return selectedFile
}

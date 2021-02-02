package platform

//OpenFileDialog uses platform agnostic package for a file dialog
func OpenFileDialog() string {
	selectedFile := wailsRuntime.Dialog.SelectFile()
	return selectedFile
}

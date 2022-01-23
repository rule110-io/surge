package models

type FileListing struct {
	FileName      string
	FileHash      string
	FileSize      int64
	NumChunks     int
	NumSeeders    int
	Topic         string
	IsTracked     bool
	IsDownloading bool
	IsUploading   bool
}

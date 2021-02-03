package models

// declaring a student struct
type GeneralFile struct {
	FileLocation  string
	FileName      string
	FileSize      int64
	FileHash      string
	Seeders       []string
	Path          string
	NumChunks     int
	IsDownloading bool
	IsUploading   bool
	IsPaused      bool
	IsMissing     bool
	IsHashing     bool
	IsTracked     bool
	IsAvailable   bool
	ChunkMap      []byte
	ChunksShared  int
	SeederCount   int
	Progress      float32
}

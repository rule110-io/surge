// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Model for File
	A File model describes any type of file handled by surge - regardless of being remote or local
*/

package models

type File struct {
	FileName      string
	FileSize      int64
	FileHash      string
	Path          string //only for local
	NumChunks     int
	IsDownloading bool
	IsUploading   bool
	IsPaused      bool
	IsMissing     bool
	IsHashing     bool //only for local
	IsTracked     bool //only for local
	IsAvailable   bool //only for local
	ChunkMap      []byte
	ChunksShared  int
	Progress      float32 //only for remote
	Topic         string
	DateTimeAdded int64
}

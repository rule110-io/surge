// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Model for File
	A File model describes any type of file handled by surge - regardless of being remote or local
*/

package models

type FileTransfer struct {
	FileName      string
	FileSize      int64
	FileHash      string
	IsDownloading bool
	IsUploading   bool
	IsPaused      bool
	IsMissing     bool
	IsHashing     bool
	IsTracked     bool
	IsAvailable   bool
	Progress      float32
	Topic         string
	NumSeeders    int
}

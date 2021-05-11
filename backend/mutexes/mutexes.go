// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This file contains mutexes
	They help to circumvent concurrent map rewrites
*/

package mutexes

import "sync"

var BandwidthAccumulatorMapLock = &sync.Mutex{}

// Mutex for reading or mutating the File model
var ChunkInTransitLock = &sync.Mutex{}

// Mutex for reading or mutating the File model
var FileWriteLock = &sync.Mutex{}

// Mutex for reading or mutating the ListedFiles collection
var ListedFilesLock = &sync.Mutex{}

var WorkerMapLock = &sync.Mutex{}
// Mutex for reading or mutating the TopicsMap collection
var TopicsMapLock = &sync.Mutex{}

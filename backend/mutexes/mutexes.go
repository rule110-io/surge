package mutexes

import "sync"

var BandwidthAccumulatorMapLock = &sync.Mutex{}
var ChunkInTransitLock = &sync.Mutex{}
var FileWriteLock = &sync.Mutex{}

//ListedFilesLock lock this whenever you're reading or mutating the ListedFiles collection
var ListedFilesLock = &sync.Mutex{}

// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This are the constants for surge configurations
	Changing them can lead to unforeseen consequences. So make sure you understand what you're doing.
*/

package constants

import "time"

const (
	//ChunkSize is size of chunk in bytes (1024 kB)
	ChunkSize = 1024 * 1024

	//NumClients is the number of NKN clients
	NumClients = 8

	//NumWorkers is the total number of concurrent chunk fetches allowed
	NumWorkers = 8

	//duration of a subscription blocktime is ~20sec
	SubscriptionDuration = 180

	//RescanPeerInterval the frequency of which subscriptions and file queries are polled
	RescanPeerInterval = time.Minute
)
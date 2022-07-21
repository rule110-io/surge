// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This are the constants for surge configurations
	Changing them can lead to unforeseen consequences. So make sure you understand what you're doing.
*/

package constants

const (
	//ChunkSize is size of chunk in bytes (1024 kB)
	ChunkSize = 1024 * 1024

	//NumClients is the number of NKN clients
	NumClients    = 4
	NumClientsMin = 1
	NumClientsMax = 8

	//NumWorkers is the total number of concurrent chunk fetches allowed
	NumWorkers    = 8
	NumWorkersMin = 1
	NumWorkersMax = 12

	//duration of a subscription blocktime is ~20sec
	SubscriptionDuration = 4000

	//official surge wallets
	TeamAddressA = "7a48870a43d1512e467e8df103b1dee8d908f297ffe1fb45e81317965597bc7c"
	TeamAddressB = "44734f736b31e522e9be64a812cf42d0822c765f4bc13404d3169ff8e3d54c9e"
	TeamAddressC = "68a10e26288b9e97fc97362eb935574dd3db74004f0081918ee121d15ed1d29b"

	//Meta payload for onchain transactions
	TransactionMeta = "Surge 2.0 Client"
)

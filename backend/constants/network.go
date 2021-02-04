// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This are the constants for network configurations
	Changing them can lead to unforeseen consequences. So make sure you understand what you're doing.
*/

package constants

const (
	//PublicTopic the public topic for subscriptions
	PublicTopic = "privateTest"

	//SurgeChunkID .
	SurgeChunkID byte = 0x001

	//SurgeQueryRequestID .
	SurgeQueryRequestID byte = 0x002

	//SurgeQueryResponseID .
	SurgeQueryResponseID byte = 0x003

	//NknClientDialTimeout time before timeout error on dial with nkn client
	NknClientDialTimeout = 10000

	//WorkerChunkReceiveTimeout is the time till a chunk request is considered a timeout and the chunk is requeued
	WorkerChunkReceiveTimeout = 60 //seconds

	//WorkerGetSessionTimeout when the session activity is older than this value the worker considers the session lost and moves on
	WorkerGetSessionTimeout = 10 //seconds

	//SendQueryRequestSessionTimeout when the session activity is older than this value the send query request is not sent
	SendQueryRequestSessionTimeout = 60 //seconds

	//GetSessionDialTimeout time till dial timeout
	GetSessionDialTimeout = 60
)

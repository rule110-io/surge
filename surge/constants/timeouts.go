package constants

import "time"

const (
	//NknClientDialTimeout time before timeout error on dial with nkn client
	NknClientDialTimeout = 10000

	//RescanPeerInterval the frequency of which subscriptions and file queries are polled
	RescanPeerInterval = time.Minute

	//WorkerChunkReceiveTimeout is the time till a chunk request is considered a timeout and the chunk is requeued
	WorkerChunkReceiveTimeout = 60 //seconds

	//WorkerGetSessionTimeout when the session activity is older than this value the worker considers the session lost and moves on
	WorkerGetSessionTimeout = 10 //seconds

	//SendQueryRequestSessionTimeout when the session activity is older than this value the send query request is not sent
	SendQueryRequestSessionTimeout = 60 //seconds

	//GetSessionDialTimeout time till dial timeout
	GetSessionDialTimeout = 60
)

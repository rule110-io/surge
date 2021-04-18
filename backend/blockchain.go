// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This file contains all blockchain related functions
*/

package surge

import (
	"fmt"
	"log"

	"github.com/rule110-io/surge/backend/constants"
	"github.com/rule110-io/surge/backend/mutexes"
	"github.com/rule110-io/surge/backend/sessionmanager"
)

var subscribers []string

//GetSubscriptions .
func GetSubscriptions() {

	mutexes.TopicsMapLock.Lock()

	for _, topic := range topicsMap {

		subResponse, err := client.GetSubscribers(topic.NameEncoded, 0, 100, true, true)
		if err != nil {
			pushError(err.Error(), "do you have an active internet connection?")
			return
		}

		for k, v := range subResponse.SubscribersInTxPool.Map {
			subResponse.Subscribers.Map[k] = v
		}

		subscribers = []string{}
		for k, v := range subResponse.Subscribers.Map {
			if len(v) > 0 {
				if k != client.Addr().String() {
					subscribers = append(subscribers, k)
				}
			}
		}

		fmt.Println(string("\033[36m"), "Get Subscriptions", len(subscribers), string("\033[0m"))

		for _, sub := range subscribers {
			connectAndQueryJob := func(addr string) {
				_, err := sessionmanager.GetSession(addr, constants.GetSessionDialTimeout)
				if err == nil {
					fmt.Println(string("\033[36m"), "Sending file query to subscriber", addr, string("\033[0m"))
					go SendQueryRequest(addr, "Testing query functionality.")
				}
			}
			go connectAndQueryJob(sub)
		}
	}

	mutexes.TopicsMapLock.Unlock()
}

func subscribeToPubSub(topic string) {
	txnHash, err := client.Subscribe("", topic, constants.SubscriptionDuration, "Surge Beta Client", nil)
	if err != nil {
		log.Println("Probably already subscribed", err)
	} else {
		log.Println("Subscribed: ", txnHash)
	}
}

func unsubscribeToPubSub(topic string) {
	txnHash, err := client.Unsubscribe("", topic, nil)
	if err != nil {
		log.Println("Probably already subscribed", err)
	} else {
		log.Println("Subscribed: ", txnHash)
	}
}

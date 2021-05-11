// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This file contains all blockchain related functions
*/

package surge

import (
	"log"

	"github.com/rule110-io/surge/backend/constants"
)

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

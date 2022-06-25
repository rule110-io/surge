// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This file contains all blockchain related functions
*/

package surge

import (
	"log"

	"github.com/nknorg/nkn-sdk-go"
	"github.com/rule110-io/surge/backend/constants"
)

var TransactionFee string

func subscribeToPubSub(topic string) {
	config := &nkn.DefaultTransactionConfig
	config.Fee = CalculateFee(TransactionFee)

	txnHash, err := client.Subscribe("", topic, constants.SubscriptionDuration, "Surge Beta Client", config)
	if err != nil {
		log.Println("Probably already subscribed to:", topic, "error:", err)
	} else {
		log.Println("Subscribed: ", topic, txnHash, "fee paid:", config.Fee)
	}
}

func unsubscribeToPubSub(topic string) {
	config := &nkn.DefaultTransactionConfig
	config.Fee = CalculateFee(TransactionFee)

	txnHash, err := client.Unsubscribe("", topic, config)
	if err != nil {
		log.Println("Probably not subscribed to:", topic, "error:", err)
	} else {
		log.Println("Unsubscribed: ", topic, txnHash, "fee paid:", config.Fee)
	}
}

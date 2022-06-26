// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This file contains all blockchain related functions
*/

package surge

import (
	"log"
	"strconv"

	"github.com/nknorg/nkn-sdk-go"
	"github.com/rule110-io/surge/backend/constants"
)

var TransactionFee string

func subscribeToPubSub(topic string) {
	config := &nkn.DefaultTransactionConfig
	config.Fee = CalculateFee(TransactionFee)

	feeFloat, _ := strconv.ParseFloat(config.Fee, 64)
	hasBalance, _ := ValidateBalanceForTransaction(0, feeFloat)
	if !hasBalance {
		pushError("Error on subscribe to topic", "Not enough fee in wallet, consider depositing NKN or if possible lower transaction fees in the wallet settings.")
		updateTopicSubscriptionState(topic, 0)
		return
	}

	updateTopicSubscriptionState(topic, 1)
	txnHash, err := client.Subscribe("", topic, constants.SubscriptionDuration, "Surge Beta Client", config)
	if err != nil {
		log.Println("Probably already subscribed to:", topic, "error:", err)
	} else {
		log.Println("Subscribed: ", topic, txnHash, "fee paid:", config.Fee)
	}
	updateTopicSubscriptionState(topic, 2)
}

func unsubscribeToPubSub(topic string) {
	config := &nkn.DefaultTransactionConfig
	config.Fee = CalculateFee(TransactionFee)

	feeFloat, _ := strconv.ParseFloat(config.Fee, 64)
	hasBalance, _ := ValidateBalanceForTransaction(0, feeFloat)
	if !hasBalance {
		pushError("Error on unsubscribe to topic", "Not enough fee in wallet, consider depositing NKN or if possible lower transaction fees in the wallet settings.")
		updateTopicSubscriptionState(topic, 2)
	}
	updateTopicSubscriptionState(topic, 1)

	txnHash, err := client.Unsubscribe("", topic, config)
	if err != nil {
		log.Println("Probably not subscribed to:", topic, "error:", err)
	} else {
		log.Println("Unsubscribed: ", topic, txnHash, "fee paid:", config.Fee)
	}
	updateTopicSubscriptionState(topic, 0)
}

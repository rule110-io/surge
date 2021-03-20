// Copyright 2021 rule101. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	This is the surge account management code
	It takes care of initializing an account as well as getting specific information about a users account
*/

package surge

import (
	"io/ioutil"
	"os"
	"time"

	"log"

	nkn "github.com/nknorg/nkn-sdk-go"
	"github.com/rule110-io/surge/backend/platform"
)

const accountPath = "account.surge"

// InitializeAccount will create an account file and return it. If there is already an account in place it will just return the existing account.
func InitializeAccount() *nkn.Account {
	var seed []byte

	var err error
	var dir = platform.GetSurgeDir()

	var accountPathOS = dir + string(os.PathSeparator) + accountPath

	_, err = os.Stat(accountPathOS)

	// If the file doesn't exist, create it
	if os.IsNotExist(err) {
		account, err := nkn.NewAccount(nil)
		seed = account.Seed()

		f, err := os.OpenFile(accountPathOS, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Panic(err)
		}
		if _, err := f.Write(seed); err != nil {
			log.Panic(err)
		}
		if err := f.Close(); err != nil {
			log.Panic(err)
		}
	} else { //else read seed from file
		file, err := os.Open(accountPathOS)
		if err != nil {
			log.Panic(err)
		}
		defer file.Close()

		seed, err = ioutil.ReadAll(file)
	}

	account, err := nkn.NewAccount(seed)
	if err != nil {
		log.Panic(err)
	}

	return account
}

//GetAccountAddress returns current client address
func GetAccountAddress() string {
	for !clientInitialized {
		time.Sleep(time.Millisecond * 50)
	}
	return client.Addr().String()
}

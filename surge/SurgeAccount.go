package surge

import (
	"io/ioutil"
	"log"
	"os"

	nkn "github.com/nknorg/nkn-sdk-go"
)

const accountPath = "account.surge"

// InitializeAccount will either create or fetch existing account
func InitializeAccount() *nkn.Account {
	var seed []byte

	_, err := os.Stat(accountPath)

	// If the file doesn't exist, create it
	if os.IsNotExist(err) {
		account, err := nkn.NewAccount(nil)
		seed = account.Seed()

		f, err := os.OpenFile(accountPath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write(seed); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	} else { //else read seed from file
		file, err := os.Open(accountPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		seed, err = ioutil.ReadAll(file)
	}

	account, err := nkn.NewAccount(seed)
	if err != nil {
		log.Fatal(err)
	}

	return account
}

package surge

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"

	nkn "github.com/nknorg/nkn-sdk-go"
)

const accountPath = "account.surge"

// InitializeAccount will either create or fetch existing account
func InitializeAccount() *nkn.Account {
	var seed []byte

	var err error
	var dir = ""

	if runtime.GOOS == "darwin" {
		dir, _ = os.UserHomeDir()
		dir = dir + string(os.PathSeparator) + ".surge"
	} else {
		dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
	}

	var accountPathOS = dir + string(os.PathSeparator) + accountPath

	_, err = os.Stat(accountPathOS)

	// If the file doesn't exist, create it
	if os.IsNotExist(err) {
		account, err := nkn.NewAccount(nil)
		seed = account.Seed()

		f, err := os.OpenFile(accountPathOS, os.O_WRONLY|os.O_CREATE, 0644)
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
		file, err := os.Open(accountPathOS)
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

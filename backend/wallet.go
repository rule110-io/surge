package surge

import (
	"log"

	nkn "github.com/nknorg/nkn-sdk-go"
)

func WalletAddress() string {
	wallet, _ := nkn.ClientAddrToWalletAddr(GetAccountAddress())
	return wallet
}

func WalletTransfer(address string, amount string, fee string) (bool, string) {
	config := &nkn.DefaultTransactionConfig
	config.Fee = fee

	result, err := client.Transfer(address, amount, config)
	if err != nil {
		pushError("Transfer failed", err.Error())
		return false, ""
	}
	log.Println("Transfered " + amount + " nkn to " + address + " txHash: " + result)
	return true, result
}

func WalletBalance() string {
	amount, err := client.Balance()
	if err != nil {
		pushError("Transfer failed", err.Error())
		return "-1"
	}
	return amount.String()
}

/*Wallet Features
- Import wallet from private key
- Export wallet private key
- Send transaction with (amount, fee, toAddress)
- WalletInfo (retrieve personal wallet address + wallet balance)
- Set/Get transaction fee default


- Optional:
seed + password wallet files instead of private keys?
get network average fee for last x amount of blocks
*/

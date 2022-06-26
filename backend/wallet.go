package surge

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	nkn "github.com/nknorg/nkn-sdk-go"
	"github.com/rule110-io/surge/backend/openapi"
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

func CalculateFee(Fee string) string {
	avgFee := openapi.GetAvgFee()
	avgFeeFloat, _ := strconv.ParseFloat(avgFee, 64)

	feePercent := 0.2
	lowFee := avgFeeFloat - avgFeeFloat*feePercent
	highFee := avgFeeFloat + avgFeeFloat*feePercent

	switch Fee {
	case "0":
		return "0"
	case "33":
		return fmt.Sprintf("%f", lowFee)
	case "66":
		return avgFee
	case "100":
		return fmt.Sprintf("%f", highFee)
	}

	return "0"
}

//ValidateBalanceForTransaction returns a boolean for whether there is enough balance to make a transation
func ValidateBalanceForTransaction(Amount float64, Fee float64) (bool, error) {
	if Amount < 0.00000001 {
		return false, errors.New("Minimum tip amount is 0.00000001")
	}

	balance := WalletBalance()
	balanceFloat, _ := strconv.ParseFloat(balance, 64)

	if Amount+Fee >= balanceFloat {
		return false, errors.New("Not enough nkn available required: " + fmt.Sprintf("%f", Amount+Fee) + " available: " + balance)
	}

	return true, nil
}

/*Wallet Features
- Import wallet from private key
- Export wallet private key
- ✔ Send transaction with (amount, fee, toAddress)
- ✔ WalletInfo (retrieve personal wallet address + wallet balance)
- ✔ Set/Get transaction fee default


- Optional:
seed + password wallet files instead of private keys?
get network average fee for last x amount of blocks
*/

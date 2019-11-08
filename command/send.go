package main

import (
	"github.com/spf13/cobra"
	"github.com/MatrixAINetwork/go-AIMan/AIMan"
	"github.com/MatrixAINetwork/go-AIMan/Accounts"
	"github.com/MatrixAINetwork/go-AIMan/providers"
	"github.com/MatrixAINetwork/go-AIMan/transactions"
	"math/big"
	"fmt"
	"github.com/pkg/errors"
)

func newSendCommand() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send a Man Transaction",
		// PreRunE block for this command will check to make sure enrollment
		// information exists before running the command
		RunE: normalSend,
	}
	return sendCmd
}
func ParseBig256(s string) (*big.Int, bool) {
	if s == "" {
		return new(big.Int), true
	}
	var bigint *big.Int
	var ok bool
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		bigint, ok = new(big.Int).SetString(s[2:], 16)
	} else {
		bigint, ok = new(big.Int).SetString(s, 10)
	}
	if ok && bigint.BitLen() > 256 {
		bigint, ok = nil, false
	}
	return bigint, ok
}

func normalSend(cmd *cobra.Command, args []string) error {
	from := args[0]
	to := args[1]
	value := args[2]
	amount,ok := ParseBig256(value)
	amount.Mul(amount,big.NewInt(1e18))
	if !ok {
		return errors.New("amount value set error")
	}
	aiMan := AIMan.NewAIMan(providers.NewHTTPProvider(walletConfig.RPC, 100, true))
	nonce,err := aiMan.Man.GetTransactionCount(from,"latest")
	if err != nil {
		return err
	}
	keystorePath := "./keystore"
	manager := Accounts.NewKeystoreManager(keystorePath, walletConfig.ChainID)
	manager.Unlock(from,"R7c5Rsrj1Q7r4d5fp")
	trans := transactions.NewTransaction(nonce.Uint64(),to,amount,200000,big.NewInt(18e9),
		nil,0,0,0)
	raw,err := manager.SignTx(trans,from)
	if err != nil {
		return err
	}
	txID, err := aiMan.Man.SendRawTransaction(raw)
	if err != nil {
		return err
	}
	fmt.Println(txID)
	return nil
}
package main

import (
	"fmt"
	"github.com/MatrixAINetwork/go-wallet/config"
	"path/filepath"
	"github.com/spf13/cobra"
	"github.com/MatrixAINetwork/go-AIMan/Accounts"
	"github.com/MatrixAINetwork/go-AIMan/AIMan"
	"github.com/MatrixAINetwork/go-AIMan/manager"
	"math/big"
)

func newBalanceCommand() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "balance",
		Short: "balance a Man Private key",
		// PreRunE block for this command will check to make sure enrollment
		// information exists before running the command
		RunE: Balance,
	}
	return sendCmd
}

func Balance(cmd *cobra.Command, args []string) error {
	allAccounts := make([]accountInfo,0)
	config.LoadJSON(filepath.Join(walletConfig.KeystorePath,"MidAccount.json"),&allAccounts)
	manager := &manager.Manager{AIMan.NewAIMan(newHTTPProvider()),
		Accounts.NewKeystoreManager(walletConfig.KeystorePath, walletConfig.ChainID)}
	man := new(big.Int).SetUint64(1e18)
	for _,info := range allAccounts {
		manager.Man.GetBalance(info.Address,"lastest")
		balance,_ := manager.Man.GetBalance(info.Address,"latest")
		bal := balance[0].Balance.ToInt()
		bal.Div(bal,man)
		fmt.Println(info.Address,bal)
	}
	return nil
}

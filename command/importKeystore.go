package main

import (
	"github.com/spf13/cobra"
	"github.com/MatrixAINetwork/go-AIMan/Accounts"
	"fmt"
	"github.com/MatrixAINetwork/go-wallet/config"
	"github.com/MatrixAINetwork/go-matrix/crypto"
	"path/filepath"
	"github.com/MatrixAINetwork/go-matrix/base58"
)

func newImportCommand() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "import",
		Short: "import a Man Private key",
		// PreRunE block for this command will check to make sure enrollment
		// information exists before running the command
		RunE: Import,
	}
	return sendCmd
}
type PrivateInfo struct {
	Address string `json:"address"`
	PrivateKey string `json:"privateKey"`
}
func Import(cmd *cobra.Command, args []string) error {
	file := args[0]
	privateInfos := make([]PrivateInfo,0)
	config.LoadJSON(file,&privateInfos)
	manager := Accounts.NewKeystoreManager(walletConfig.KeystorePath, walletConfig.ChainID)
	allAccounts := make([]accountInfo,0)
	for _,info := range privateInfos  {
		priKey,err := crypto.HexToECDSA(info.PrivateKey)
		if err!=nil {
			fmt.Println("private key is error", info.PrivateKey)
		}
		password := string(randPassPhrase())
		address := crypto.PubkeyToAddress(priKey.PublicKey)
		manAddr := base58.Base58EncodeToString("MAN",address)
		if manAddr != info.Address {
			fmt.Println("Address is Error","OldAddr=",info.Address,"newAddr=",manAddr)
		}
		manager.Keystore.ImportECDSA(priKey,password)
		allAccounts = append(allAccounts,accountInfo{manAddr,password})
	}
	err := config.SaveJSON(filepath.Join(walletConfig.KeystorePath,"import.json"),&allAccounts)
	if err != nil {
		return err
	}

	return nil
}
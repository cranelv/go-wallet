package main

import (
	"github.com/spf13/cobra"
	"math/big"
	"github.com/MatrixAINetwork/go-AIMan/AIMan"
	"github.com/MatrixAINetwork/go-AIMan/providers"
	"github.com/MatrixAINetwork/go-AIMan/Accounts"
	"github.com/MatrixAINetwork/go-AIMan/transactions"
	"strconv"
	"time"
	"fmt"
)

func newRushCommand() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "rush",
		Short: "rush to Send large number Man Transactions",
		// PreRunE block for this command will check to make sure enrollment
		// information exists before running the command
		RunE: normalRush,
	}
	return sendCmd
}
func newRushIPCCommand() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "rushIPC",
		Short: "rush to Send large number Man Transactions use IPC Message",
		// PreRunE block for this command will check to make sure enrollment
		// information exists before running the command
		RunE: normalRushIPC,
	}
	return sendCmd
}
func newRushIPCCommand1() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "rushIPC1",
		Short: "rush to Send large number Man Transactions use IPC Message",
		// PreRunE block for this command will check to make sure enrollment
		// information exists before running the command
		RunE: normalRushIPC1,
	}
	return sendCmd
}
func normalRush(cmd *cobra.Command, args []string) error {
	aiMan := AIMan.NewAIMan(providers.NewHTTPProvider(args[4],100,false))
	return rushCommand(args,aiMan,270)
}
func normalRushIPC(cmd *cobra.Command, args []string) error {
	aiMan := AIMan.NewAIMan(providers.NewIPCProvider(args[4]))
	return rushCommand(args,aiMan,10000)
}
func normalRushIPC1(cmd *cobra.Command, args []string) error {
	aiMan := AIMan.NewAIMan(providers.NewIPCProvider(args[4]))
	return rushCommand1(args,aiMan,10000)
}
func rushCommand(args []string,aiMan *AIMan.AIMan,cap int) error {
	from := args[0]
	nonce,err := aiMan.Man.GetTransactionCount(from,"latest")
	if err != nil {
		return err
	}
	password := args[1]
	keystorePath := args[3]
	chainID,err := strconv.Atoi(args[5])
	manager := Accounts.NewKeystoreManager(keystorePath, int64(chainID))
	manager.Unlock(from,password)
	nonce1 := nonce.Uint64()
	size,err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	for i:=0;i<size;i++{
		err := TxPackage(nonce1,from,manager,aiMan,cap)
		if err != nil {
			return err
		}
		nonce1 += uint64(cap)
	}
	return nil
}
func rushCommand1(args []string,aiMan *AIMan.AIMan,cap int) error {
	from := args[0]
	nonce,err := aiMan.Man.GetTransactionCount(from,"latest")
	if err != nil {
		return err
	}
	password := args[1]
	keystorePath := args[3]
	chainID,err := strconv.Atoi(args[5])
	manager := Accounts.NewKeystoreManager(keystorePath, int64(chainID))
	manager.Unlock(from,password)
	nonce1 := nonce.Uint64()
	size,err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	raws := make([]interface{},size*cap)
	to := "MAN.4BRmmxsC9iPPDyr8CRpRKUcp7GAww"
	amount := big.NewInt(10000)
	for i:=0;i<size*cap;i++{
		trans := transactions.NewTransaction(nonce1,to,amount,200000,big.NewInt(18e9),
			nil,0,0,0)
		raw,err := manager.SignTx(trans,from)
		if err != nil{
			return err
		}
		raws[i] = raw
		nonce1++
	}
	fmt.Println(len(raws))
	for i:=0;i<size;i++{
		aiMan.Man.BatchSendRawTransaction(raws[i*cap:(i+1)*cap])
	fmt.Println(i*cap,(i+1)*cap)
		time.Sleep(time.Millisecond)
	}
	return nil
}
func TxPackage(nonce uint64,from string,manager *Accounts.KeystoreManager,aiMan *AIMan.AIMan,cap int) error {
	size := cap
	raws := make([]interface{},size)
	to := "MAN.4BRmmxsC9iPPDyr8CRpRKUcp7GAww"
	amount := big.NewInt(10000)
	for i:=0;i<size;i++{
		trans := transactions.NewTransaction(nonce,to,amount,200000,big.NewInt(18e9),
			nil,0,0,0)
		raw,err := manager.SignTx(trans,from)
		if err != nil{
			return err
		}
		raws[i] = raw
		nonce++
	}
	aiMan.Man.BatchSendRawTransaction(raws)
	return nil
}
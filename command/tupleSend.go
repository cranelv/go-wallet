package main

import (
	"github.com/spf13/cobra"
	"math/big"
	"github.com/MatrixAINetwork/go-AIMan/AIMan"
	"github.com/MatrixAINetwork/go-AIMan/providers"
	"github.com/MatrixAINetwork/go-AIMan/Accounts"
	"github.com/MatrixAINetwork/go-AIMan/transactions"
	"fmt"
	"errors"
	"github.com/MatrixAINetwork/go-AIMan/waiting"
	"time"
	"github.com/MatrixAINetwork/go-AIMan/manager"
	"crypto/rand"
	"github.com/MatrixAINetwork/go-matrix/base58"
	"github.com/MatrixAINetwork/go-wallet/config"
	"path/filepath"
	"strings"
)
var walletConfig = config.NewWalletConfig()
func newHTTPProvider()*providers.HTTPProvider{
	secure := false
	rpc := walletConfig.RPC
	if strings.Index(walletConfig.RPC,"http://") == 0 {
		rpc = rpc[len("http://"):len(walletConfig.RPC)]
	}else if strings.Index(walletConfig.RPC,"https://") == 0 {
		rpc = rpc[len("https://"):len(walletConfig.RPC)]
		secure = true
	}
	return providers.NewHTTPProvider(rpc, 100, secure)
}
func newTupleCommand() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "tuple",
		Short: "tuple Send a Man Transaction",
		// PreRunE block for this command will check to make sure enrollment
		// information exists before running the command
		RunE: tupleSend,
	}
	return sendCmd
}
func towei(args string)(*big.Int,error){
	amount,ok := new(big.Float).SetString(args)
	if !ok {
		return nil,errors.New("amount value set error")
	}
	amount.Mul(amount,big.NewFloat(1.0e18))
	value1 := new(big.Int)
	value1,_ = amount.Int(value1)
	return value1,nil
}
func tupleSend(cmd *cobra.Command, args []string) error {
	from := args[0]
	to := args[1]
	value := args[2]
	amount,err := towei(value)
	password := args[3]
	if err != nil {
		return err
	}
//	amount.Mul(amount,big.NewInt(1e18))
	gas := big.NewInt(18e9)
	gas.Mul(gas,big.NewInt(21000))
	manager := &manager.Manager{AIMan.NewAIMan(newHTTPProvider()),
		Accounts.NewKeystoreManager(walletConfig.KeystorePath, walletConfig.ChainID)}
	if amount.Sign() == 0 {
		balance,err:= manager.Man.GetBalance(from,"latest")
		if err != nil{
			return err
		}
		amount = balance[0].Balance.ToInt()
		amount.Sub(amount,gas)
		if amount.Sign() <= 0 {
			return nil
		}
	}
	accounts,err := newAccounts(manager,walletConfig.SendNum)
	if err != nil {
		return err
	}
	allAccounts := make([]accountInfo,0)
	config.LoadJSON(filepath.Join(walletConfig.KeystorePath,"MidAccount.json"),&allAccounts)
	allAccounts = append(allAccounts,accounts...)
	err = config.SaveJSON(filepath.Join(walletConfig.KeystorePath,"MidAccount.json"),&allAccounts)
	if err != nil {
		return err
	}
	froms := []accountInfo{accountInfo{Address:from,PassPhrase:password}}
	froms = append(froms,accounts...)
	for i:=0;i<len(accounts);i++  {
		err := stepSend(froms[i].Address,froms[i+1].Address,froms[i].PassPhrase,amount,manager)
		if err != nil{
			return err
		}
		amount.Sub(amount,gas)
	}
	err = stepSend(froms[len(accounts)].Address,to,froms[len(accounts)].PassPhrase,amount,manager)
	if err != nil{
		return err
	}
	return nil
}
func randPassPhrase() []byte {
	phraseText := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	phrase := []byte(phraseText)
	length := len(phrase)
	buff := make([]byte,16)
	rand.Reader.Read(buff)
	var retText [16]byte
	for i,tx := range buff {
		retText[i] = phrase[int(tx)%length]
	}
	return retText[:]
}
type accountInfo struct {
	Address string `json:"address"`
	PassPhrase string `json:"password"`
}
func newAccounts(manager *manager.Manager,num int)([]accountInfo,error){
	accInfo := make([]accountInfo,num)
	for i:=0;i<num;i++ {
		phrase := string(randPassPhrase())
		acc,err := manager.Keystore.NewAccount(phrase)
		if err != nil {
			return nil,err
		}
		accInfo[i].Address = base58.Base58EncodeToString("MAN",acc.Address)
		accInfo[i].PassPhrase = phrase
	}
	return accInfo,nil
}

func stepSend(from,to,password string, amount *big.Int,manager *manager.Manager) error{
	manager.Unlock(from,password)
	blockNumber, err := manager.Man.GetBlockNumber()
	if err != nil {
		return err
	}
	nonce,err := manager.Man.GetTransactionCount(from,"latest")
	if err != nil {
		return err
	}
	trans := transactions.NewTransaction(nonce.Uint64(),to,amount,21000,big.NewInt(18e9),
		nil,0,0,0)
	raw,err := manager.SignTx(trans,from)
	if err != nil {
		return err
	}
	txID, err := manager.Man.SendRawTransaction(raw)
	if err != nil {
		return err
	}
	fmt.Println(txID)
	wait3 := waiting.NewMultiWaiting(waiting.NewWaitBlockHeight(manager,blockNumber.Uint64()+10),
		waiting.NewWaitTime(200*time.Second),
		waiting.NewWaitTxReceipt(manager,txID))
	index := wait3.Waiting()
	if index != 2{
		return errors.New("Time Out")
	}
	return nil
}
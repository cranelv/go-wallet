package config

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"path/filepath"
	"os"
)

func findLine(data []byte, offset int64) (line int) {
	line = 1
	for i, r := range string(data) {
		if int64(i) >= offset {
			return
		}
		if r == '\n' {
			line++
		}
	}
	return
}

func LoadJSON(file string, val interface{}) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, val); err != nil {
		if syntaxerr, ok := err.(*json.SyntaxError); ok {
			line := findLine(content, syntaxerr.Offset)
			return fmt.Errorf("JSON syntax error at %v:%v: %v", file, line, err)
		}
		return fmt.Errorf("JSON unmarshal error in %v: %v", file, err)
	}
	return nil
}
type WalletConfig struct {
	RPC string `json:"rpc"`
	SendNum int `json:"sendnum"`
	KeystorePath string `json:"keystorepath"`
	ChainID int64 `json:"chainid"`
}
func getCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}
func NewWalletConfig() *WalletConfig{
	config := &WalletConfig{}
	path := getCurrentDirectory()
	LoadJSON(filepath.Join(path,"config.json"),&config)
	return config
}
func SaveJSON(file string,val interface{}) error  {
	out, _ := json.MarshalIndent(val, "", "  ")
	if err := ioutil.WriteFile(file, out, 0644); err != nil {
		return err
	}
	return nil

}

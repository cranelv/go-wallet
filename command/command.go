package main

import (
	"os"
	"github.com/spf13/viper"
	"github.com/spf13/cobra"
	"strings"
)

func RunMain(args []string) error {
	// Save the os.Args
	saveOsArgs := os.Args
	os.Args = args

	// Execute the command
	cmdName := ""
	if len(args) > 1 {
		cmdName = args[1]
	}
	ccmd := NewCommand(cmdName)
	err := ccmd.Execute()

	// Restore original os.Args
	os.Args = saveOsArgs

	return err
}
type ClientCmd struct {
	// name of the sub command
	name string
	// rootCmd is the base command for this Wallet
	rootCmd *cobra.Command
	// My viper instance
	myViper *viper.Viper
	// cfgFileName is the name of the configuration file
	cfgFileName string
	// homeDirectory is the location of the client's home directory
	homeDirectory string
	// Set to log level
	logLevel string
}

// NewCommand returns new ClientCmd ready for running
func NewCommand(name string) *ClientCmd {
	c := &ClientCmd{
		myViper: viper.New(),
	}
	c.name = strings.ToLower(name)
	c.init()
	return c
}

// Execute runs this ClientCmd
func (c *ClientCmd) Execute() error {
	return c.rootCmd.Execute()
}

// init initializes the ClientCmd instance
// It intializes the cobra root and sub commands and
// registers command flgs with viper
func (c *ClientCmd) init() {
	c.rootCmd = &cobra.Command{
		Use:   "MANWallet",
		Short: "Matrix AI Network Blockchain wallet",
	}
	c.rootCmd.AddCommand(newSendCommand())
	c.rootCmd.AddCommand(newRushCommand())
	c.rootCmd.AddCommand(newRushIPCCommand())
	c.rootCmd.AddCommand(newRushIPCCommand1())
	//c.rootCmd.AddCommand
}


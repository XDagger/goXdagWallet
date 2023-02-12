package cli

import (
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"goXdagWallet/components"
	"goXdagWallet/config"
	"goXdagWallet/xdago/base58"
	"goXdagWallet/xdago/cryptography"
	"goXdagWallet/xdago/secp256k1"
	"goXdagWallet/xlog"
	"os"
	"strings"
	"time"
)

var WalletAccount components.WalletState
var spin = spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner

func NewCli(accountStatus int) {
	fmt.Println("Starting xdag wallet, version ", config.GetConfig().Version)
	xlog.Info("Starting xdag wallet, version ", config.GetConfig().Version)
	WalletAccount.HasAccount = accountStatus >= 0
	WalletAccount.WalletType = accountStatus
	if WalletAccount.HasAccount {
		connectWallet()
	} else {
		registerWallet()
	}
	OpenAndRunWallet()
}

func connectWallet() {
	fmt.Println("Connecting wallet....")
	if WalletAccount.WalletType == components.HAS_BOTH {
		selectWallet()
	} else {
		WalletAccount.Password = showPassword()
	}
}
func selectWallet() {
	prompt := promptui.Select{
		Label: "Select wallet type",
		Items: []string{"Mnemonic", "Non Mnemonic"},
	}

	_, result, err := prompt.Run()

	for err != nil {
		fmt.Printf("Input selection failed %v\n", err)
		_, result, err = prompt.Run()
	}

	if result == "Mnemonic" {
		WalletAccount.WalletType = components.HAS_ONLY_BIP
	} else {
		WalletAccount.WalletType = components.HAS_ONLY_XDAG
	}
	WalletAccount.Password = showPassword()
}
func registerWallet() {
	fmt.Println("Registering wallet....")
	prompt := promptui.Select{
		Label: "Select Wallet Type",
		Items: []string{"Create Mnemonic", "Import Mnemonic", "Import Non Mnemonic"},
	}

	_, result, err := prompt.Run()

	for err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		_, result, err = prompt.Run()
	}

	if result == "Import Non Mnemonic" {
		WalletAccount.WalletType = components.HAS_ONLY_XDAG
		inputFolderPath()
		WalletAccount.Password = showPassword()
	} else {
		WalletAccount.WalletType = components.HAS_ONLY_BIP
		if result == "Create Mnemonic" {
			pwd := showPassword()
			reshowPassword(pwd)
			w, res := components.NewBipWallet(pwd, 128)
			if res {
				components.BipWallet = w
				b := cryptography.ToBytesAddress(components.BipWallet.GetDefKey())
				components.BipAddress = base58.ChkEnc(b[:])
				OpenAndRunWallet()
			}

		} else { // Import Mnemonic
			inputFilePath()
			WalletAccount.Password = showPassword()
		}
	}
}

func OpenAndRunWallet() {
	fmt.Println("Initializing cryptography...")
	fmt.Println("Reading wallet...")
	spin.Start()
	if WalletAccount.WalletType == components.HAS_ONLY_XDAG {
		res := components.ConnectXdagWallet()
		spin.Stop()
		if res == 0 {
			RunWallet(WalletAccount.WalletType)
		} else {
			fmt.Println("Password incorrect")
		}
	} else if WalletAccount.WalletType == components.HAS_ONLY_BIP {
		res := components.ConnectBipWallet(WalletAccount.Password)
		spin.Stop()
		if res {
			b := cryptography.ToBytesAddress(components.BipWallet.GetDefKey())
			components.BipAddress = base58.ChkEnc(b[:])
			RunWallet(WalletAccount.WalletType)
		} else {
			fmt.Println("Password incorrect")
		}
	}
}

var validateCmd = func(input string) error {
	if input != "help" && input != "exit" && input != "account" &&
		input != "balance" && !strings.HasPrefix(input, "xfer ") {
		return errors.New("unknown command, input 'help' to list available commands")
	}
	if strings.HasSuffix(input, "xfer ") {
		items := strings.Fields(input)
		if len(items) != 4 {
			return errors.New("transfer command parameters error")
		}

	}

	return nil
}

func RunWallet(walletExists int) {
	prompt := promptui.Prompt{
		Label:    "Command > ",
		Validate: validateCmd,
	}
	for {
		result, err := prompt.Run()

		for err != nil {
			fmt.Printf("Input command failed %v\n", err)
			result, err = prompt.Run()
		}
		switch result {
		case "help":
			fmt.Println("---------------------------------------------------------")
			fmt.Println("      help -- display commands list")
			fmt.Println("      exit -- exit cli wallet")
			fmt.Println("   account -- display address of wallet account")
			fmt.Println("   balance -- display balance of wallet account")
			fmt.Println("xfer V A R -- transfer V coins to address A with remark R")
			fmt.Println("---------------------------------------------------------")
			break
		case "exit":
			os.Exit(0)
		case "account":
			if walletExists == components.HAS_ONLY_BIP {
				fmt.Println(components.BipAddress)
			} else if walletExists == components.HAS_ONLY_XDAG {
				fmt.Println(components.XdagAddress)
			}
			break
		case "balance":
			var balance string
			var err error
			spin.Start()
			if walletExists == components.HAS_ONLY_BIP {
				balance, err = components.BalanceRpc(components.BipAddress)
			} else if walletExists == components.HAS_ONLY_XDAG {
				balance, err = components.BalanceRpc(components.XdagAddress)
			}
			spin.Stop()
			if err != nil {
				fmt.Println("Get balance failed", err)
			} else {
				fmt.Println(balance)
			}
			break
		}

		if strings.HasPrefix(result, "xfer ") {
			var fromAddress string
			var fromKey *secp256k1.PrivateKey
			items := strings.Fields(result)
			spin.Start()
			if walletExists == components.HAS_ONLY_BIP {
				fromAddress = components.BipAddress
				fromKey = components.BipWallet.GetDefKey()
			} else if walletExists == components.HAS_ONLY_XDAG {
				fromAddress = components.XdagAddress
				fromKey = components.XdagKey
			}
			err := components.TransferRpc(fromAddress, items[2], items[1], items[3], fromKey)
			spin.Stop()
			if err != nil {
				fmt.Println("Transfer failed", err)
			} else {
				fmt.Println("Transfer has been committed. Please wait for a while to get it done")
			}
		}
	}
}

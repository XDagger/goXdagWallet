package cli

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"goXdagWallet/components"
	"goXdagWallet/config"
	"time"
)

var WalletAccount components.WalletState

func NewCli(accountStatus int) {
	fmt.Println("Starting xdag wallet, version ", config.GetConfig().Version)
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
	fmt.Println("Connecting wallet ....")
	if WalletAccount.WalletType == components.HAS_BOTH {
		selectWallet()
	} else {
		showPassword()
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
	fmt.Println("Registering wallet ....")
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
		} else { // Import Mnemonic
			inputFilePath()
			WalletAccount.Password = showPassword()
		}
	}
}

func OpenAndRunWallet() {
	fmt.Println("Initializing cryptography...")
	fmt.Println("Reading wallet...")
	s := spinner.New(spinner.CharSets[33], 100*time.Millisecond) // Build our new spinner
	s.Start()
	if WalletAccount.WalletType == components.HAS_ONLY_XDAG {
		res := components.ConnectXdagWallet()
		if res == 0 {
			RunWallet(WalletAccount.WalletType)
		} else {
			fmt.Println("Password incorrect")
			return
		}
	} else if WalletAccount.WalletType == components.HAS_ONLY_BIP {
		res := components.ConnectBipWallet()
		if res {
			RunWallet(WalletAccount.WalletType)
		} else {
			fmt.Println("Password incorrect")
		}
	}
}

func RunWallet(walletExists int) {

}

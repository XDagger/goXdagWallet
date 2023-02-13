package cli

import (
	"errors"
	"fmt"
	"goXdagWallet/components"
	"goXdagWallet/config"
	"goXdagWallet/xdago/base58"
	"goXdagWallet/xdago/common"
	"goXdagWallet/xdago/cryptography"
	"goXdagWallet/xdago/secp256k1"
	bip "goXdagWallet/xdago/wallet"
	"goXdagWallet/xlog"
	"os"
	"path"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
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
		copy(components.Password[:], WalletAccount.Password)
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
	copy(components.Password[:], WalletAccount.Password)
}
func registerWallet() {
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
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
			var psDone bool
			var passwd string
			for !psDone {
				passwd = showPassword()
				psDone = reshowPassword(passwd)
			}

			w, res := components.NewBipWallet(passwd, 128)
			if res {
				components.BipWallet = w
				b := cryptography.ToBytesAddress(components.BipWallet.GetDefKey())
				components.BipAddress = base58.ChkEnc(b[:])
				components.PwdStr = passwd
				WalletAccount.Password = passwd
				OpenAndRunWallet()
			}

		} else { // Import Mnemonic
			filePath := inputFilePath()
			WalletAccount.Password = showPassword()

			pathDest := path.Join(pwd, common.BIP32_WALLET_FOLDER)
			if err := os.RemoveAll(pathDest); err != nil {
				fmt.Println("Clear dir failed", err)
			}
			if err := os.MkdirAll(pathDest, 0666); err != nil {
				fmt.Println("Make dir failed", err)
			}
			dirDest := path.Join(pathDest, common.BIP32_WALLET_FILE_NAME)
			components.PwdStr = WalletAccount.Password
			components.BipWallet, err = bip.ImportWalletFromMnemonicFile(filePath, dirDest, WalletAccount.Password)
			if err != nil {
				fmt.Println("Import mnemonic failed", err)
			} else {
				b := cryptography.ToBytesAddress(components.BipWallet.GetDefKey())
				components.BipAddress = base58.ChkEnc(b[:])
				OpenAndRunWallet()
			}

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
		input != "balance" && !strings.HasPrefix(input, "xfer ") &&
		!strings.HasPrefix(input, "export ") {
		return errors.New("unknown command, input 'help' to list available commands")
	}
	if strings.HasSuffix(input, "xfer ") {
		items := strings.Fields(input)
		if len(items) != 4 {
			return errors.New("transfer command parameters error")
		}

	}

	if strings.HasSuffix(input, "export ") {
		items := strings.Fields(input)
		if len(items) != 2 {
			return errors.New("export command parameters error")
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
			if walletExists == components.HAS_ONLY_BIP {
				fmt.Println("  export P -- export mnemonic to file P")
			}
			fmt.Println("---------------------------------------------------------")
		case "exit":
			os.Exit(0)
		case "account":
			if walletExists == components.HAS_ONLY_BIP {
				fmt.Println(components.BipAddress)
			} else if walletExists == components.HAS_ONLY_XDAG {
				fmt.Println(components.XdagAddress)
			}
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

		if strings.HasPrefix(result, "export ") {
			items := strings.Fields(result)
			err := components.BipWallet.ExportMnemonic(items[1])
			if err != nil {
				fmt.Println("Export mnemonic failed", err)
			} else {
				fmt.Println("Export mnemonic success")
			}
		}
	}
}

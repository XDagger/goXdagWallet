package components

import (
	"fmt"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"goXdagWallet/xdago/base58"
	"goXdagWallet/xdago/common"
	"goXdagWallet/xdago/cryptography"
	xdagoUtils "goXdagWallet/xdago/utils"
	bip "goXdagWallet/xdago/wallet"
	"goXdagWallet/xlog"
	"os"
	"path"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"golang.org/x/exp/utf8string"
)

var chanBalance = make(chan int, 1)

func Xdag_Wallet_fount() int {
	hasXdagWallet := 0
	hasBip32Wallet := 0
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	pathName := path.Join(pwd, "xdagj_dat", "dnet_key.dat")
	// change current working directory
	os.Chdir(pwd)

	fi, err := os.Stat(pathName)
	if err != nil {
		hasXdagWallet = -1
	} else if fi.Size() != 2048 {
		hasXdagWallet = -1
	}
	pathName = path.Join(pwd, common.BIP32_WALLET_FOLDER, common.BIP32_WALLET_FILE_NAME)
	//os.Chdir(pwd)

	fi, err = os.Stat(pathName)
	if err != nil {
		hasBip32Wallet = -1
	} else if fi.Size() < 125 {
		hasBip32Wallet = -1
	}
	if hasXdagWallet == -1 && hasBip32Wallet == -1 {
		return -1 // no wallet
	}
	if hasXdagWallet == 0 && hasBip32Wallet == 0 { // has both wallets
		return 0
	} else if hasXdagWallet == 0 { // only has xdag wallet
		return 1
	} else {
		return 2 // only has bip32_bip44 wallet
	}
}

func ConnectBipWallet() bool {
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	wallet := bip.NewWallet(path.Join(pwd, common.BIP32_WALLET_FOLDER, common.BIP32_WALLET_FILE_NAME))
	res := wallet.UnlockWallet(PwdStr)
	if res {
		BipWallet = &wallet
	}
	return res
}
func NewBipWallet(password string) (*bip.Wallet, bool) {
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	wallet := bip.NewWallet(path.Join(pwd, common.BIP32_WALLET_FOLDER, common.BIP32_WALLET_FILE_NAME))
	wallet.UnlockWallet(password)
	wallet.InitializeHdWallet(bip.NewMnemonic())
	wallet.AddAccountWithNextHdKey()
	res := wallet.Flush()

	return &wallet, res
}

func ValidateXdagAddress(address string) bool {
	_, err := xdagoUtils.Address2Hash(address)
	return err == nil
}

func ValidateBipAddress(address string) bool {
	_, err := checkBase58Address(address)
	return err == nil
}

func ValidateRemark(remark string) bool {
	return utf8string.NewString(remark).IsASCII()
}

func NewWalletWindow(walletExists int) {
	if WalletWindow != nil {
		return
	}

	if walletExists != HAS_ONLY_XDAG {
		b := cryptography.ToBytesAddress(BipWallet.GetDefKey())
		BipAddress = base58.ChkEnc(b[:])
	}

	getBalance()
	LogonWindow.Win.Hide()
	w := WalletApp.NewWindow(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version) +
		getTestTitle())
	WalletWindow = w
	w.SetMaster()
	LogonWindow.Win.Content().Resize(fyne.NewSize(0, 0))
	tabs := container.NewAppTabs()
	if walletExists != HAS_ONLY_BIP {
		tabs.Append(container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabAccount")+"-0",
			theme.HomeIcon(), AccountPage(XdagAddress, XdagBalance, WalletWindow)))
	}
	tabs.Append(container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabAccount")+"-1",
		theme.FolderOpenIcon(), BipPage(BipAddress, BipBalance, WalletWindow)))
	tabs.Append(container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabTransfer"),
		theme.MailSendIcon(), TransferPage(WalletWindow)))
	tabs.Append(container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabHistory"),
		theme.ContentPasteIcon(), HistoryPage(WalletWindow)))
	tabs.Append(container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabDonate"),
		theme.NewThemedResource(donateIconRes), DonatePage(WalletWindow)))
	tabs.Append(container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabAbout"),
		theme.InfoIcon(), AboutPage(WalletWindow)))

	if fyne.CurrentDevice().IsMobile() {
		tabs.SetTabLocation(container.TabLocationBottom)
	} else {
		tabs.SetTabLocation(container.TabLocationLeading)
	}

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(640, 480))
	w.CenterOnScreen()
	go checkBalance()
	w.SetOnClosed(func() {
		xlog.CleanXdagLog(xlog.StdXdagLog)
		chanBalance <- 1
		WalletApp.Quit()
		os.Exit(0)
	})
	w.Show()
}

func checkBalance() {
	for {
		select {
		case <-chanBalance:
			return
		case <-time.After(time.Second * 130):
			getBalance()
		}
	}
}

func getBalance() {
	if LogonWindow.WalletExists == HAS_ONLY_BIP || LogonWindow.WalletExists == HAS_BOTH {
		balance, err := BalanceRpc(BipAddress)
		if err == nil {
			if balance != "" {
				BipBalance = balance
				BipAccountBalance.Set(balance)
			} else {
				xlog.Error("get bip32 account balance error.")
			}
		} else {
			xlog.Error(err)
		}
	}
	if LogonWindow.WalletExists == HAS_ONLY_XDAG || LogonWindow.WalletExists == HAS_BOTH {
		balance, err := BalanceRpc(XdagAddress)
		if err == nil {
			if balance != "" {
				XdagBalance = balance
				AccountBalance.Set(balance)
			} else {
				xlog.Error("get xdag account balance error.")
			}
		} else {
			xlog.Error(err)
		}
	}
}

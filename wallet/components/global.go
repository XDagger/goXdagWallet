package components

import (
	"goXdagWallet/xdago/secp256k1"
	xdagoUtils "goXdagWallet/xdago/utils"
	bip "goXdagWallet/xdago/wallet"

	"fyne.io/fyne/v2"
)

const (
	WALLET_NOT_FOUND = -1
	HAS_BOTH         = 0
	HAS_ONLY_XDAG    = 1
	HAS_ONLY_BIP     = 2
)

const XDAG_FEE float64 = 0.1

var LogonWindow LogonWin
var WalletWindow fyne.Window
var WalletApp fyne.App
var Password [256]byte
var PwdStr string
var XdagAddress string
var OldAddresses []string
var XdagBalance string
var XdagKey *secp256k1.PrivateKey
var BipWallet *bip.Wallet
var BipAddress string
var BipBalance string
var AddressVerify = make(map[string]xdagoUtils.VerifyData)
var OldKeys []*secp256k1.PrivateKey

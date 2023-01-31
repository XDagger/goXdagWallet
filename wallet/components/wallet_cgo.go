package components

//#cgo darwin LDFLAGS: -L${SRCDIR}/../../clib -lxdag_runtime_Darwin -L/usr/lib -lsecp256k1 -lm -Llocal/opt/openssl/lib -lssl -lcrypto
//#cgo linux LDFLAGS: -L${SRCDIR}/../../clib -lxdag_runtime_Linux -L/usr/lib -lsecp256k1 -lssl -lcrypto -lm
//#cgo windows LDFLAGS: -L${SRCDIR}/../../clib -lxdag_runtime_Windows -L/usr/lib -L/usr/local/lib -lsecp256k1 -lssl -lcrypto -lm -lws2_32
//#include "../../clib/xdag_runtime.h"
//#include "callback.h"
//#include <stdlib.h>
//#include <string.h>
/*
 typedef const char cchar_t;
*/
import "C"
import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"goXdagWallet/xdago/base58"
	"goXdagWallet/xdago/common"
	"goXdagWallet/xdago/cryptography"
	"goXdagWallet/xdago/secp256k1"
	xdagoUtils "goXdagWallet/xdago/utils"
	bip "goXdagWallet/xdago/wallet"
	"goXdagWallet/xlog"
	"os"
	"path"
	"strings"
	"time"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"golang.org/x/exp/utf8string"
)

var chanBalance = make(chan int, 1)

// var regDone = make(chan int, 1)

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
	}
	if fi.Size() != 2048 {
		hasXdagWallet = -1
	}
	pathName = path.Join(pwd, common.BIP32_WALLET_FOLDER, common.BIP32_WALLET_FILE_NAME)
	//os.Chdir(pwd)

	fi, err = os.Stat(pathName)
	if err != nil {
		hasBip32Wallet = -1
	}
	if fi.Size() < 125 {
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
		return 2 // only has bip32 or bip44 wallet
	}
}
func ConnectXdagWallet() int32 {
	var testnet int
	if config.GetConfig().Option.IsTestNet {
		testnet = 1
	}
	res := C.init_password_callback(C.int(testnet))
	result := int32(res)
	fmt.Println(result)
	if result == 0 {
		k := getDefaultKey()
		if k == nil {
			xlog.Error("get default key failed.")
			fmt.Println("get default key failed.")
			return -4
		} else {
			XdagKey = secp256k1.PrivKeyFromBytes(k)
			addr, err := xdagoUtils.AddressFromStorage()
			if err != nil {
				xlog.Error(err)
				return -5
			} else {
				XdagAddress = addr
				xlog.Info(addr)
				//block := transactionBlock(addr, "4smXToYpMy1648T3PXpBRZ8zSey5c6Sy7", "test", 100.5, XdagKey)
				//xlog.Info(block)
			}
		}
	}
	return result
}

func ConnectBipWallet() bool {
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	wallet := bip.NewWallet(path.Join(pwd, common.BIP32_WALLET_FOLDER, common.BIP32_WALLET_FILE_NAME))
	res := wallet.UnlockWallet(string(Password[:]))
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

//export goPasswordCallback
func goPasswordCallback(prompt *C.cchar_t, buf *C.char, size C.uint) C.int {
	C.memcpy(unsafe.Pointer(buf), unsafe.Pointer(&Password[0]), C.size_t(size))
	return C.int(0)
}

func TransferWrap(address, amount, remark string) int {
	// TODO: validate address, amount, remark

	return int(0)
}

func ValidateXdagAddress(address string) bool {
	_, err := xdagoUtils.Address2Hash(address)
	return err == nil
}

func ValidateRemark(remark string) bool {
	return utf8string.NewString(remark).IsASCII()
}

func NewWalletWindow() {
	if WalletWindow != nil {
		return
	}

	LogonWindow.Win.Hide()
	w := WalletApp.NewWindow(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version) +
		getTestTitle())
	WalletWindow = w
	w.SetMaster()
	LogonWindow.Win.Content().Resize(fyne.NewSize(0, 0))
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabAccount"),
			theme.HomeIcon(), AccountPage(XdagAddress, XdagBalance, WalletWindow)),
		container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabTransfer"),
			theme.MailSendIcon(), TransferPage(WalletWindow, TransferWrap)),
		container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabHistory"),
			theme.ContentPasteIcon(), HistoryPage(WalletWindow)),
		container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabDonate"),
			theme.NewThemedResource(donateIconRes), DonatePage(WalletWindow, TransferWrap)),
		container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabAbout"),
			theme.InfoIcon(), AboutPage(WalletWindow)),
		//container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabSettings"),
		//	theme.SettingsIcon(), SettingsPage(WalletWindow))
	)
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
			//C.xdag_get_balance_wrap()
		}
	}
}

// get xdag wallet private key
func getDefaultKey() []byte {
	p := C.xdag_get_default_key()
	if uintptr(p) > 0 {
		key := C.GoBytes(p, 32)
		//fmt.Println(hex.EncodeToString(key[:]))
		//xlog.Info("default private key:", hex.EncodeToString(key[:]))
		return key
	}
	return nil
}

func transactionBlock(from, to, remark string, value float64, key *secp256k1.PrivateKey) string {
	if key == nil {
		xlog.Error("transaction default key error")
		return ""
	}
	var inAddress string
	var err error
	if len(from) == common.XDAG_ADDRESS_SIZE { // old xdag address
		if !ValidateXdagAddress(from) {
			xlog.Error("transaction send address length error")
			return ""
		}
		hash, err := xdagoUtils.Address2Hash(from)
		if err != nil {
			xlog.Error(err)
			return ""
		}
		inAddress = hex.EncodeToString(hash[:24])
	} else { // new base58 address
		inAddress, err = checkBase58Address(from)
		if err != nil {
			xlog.Error(err)
			return ""
		}
	}

	outAddress, err := checkBase58Address(to)
	if err != nil {
		xlog.Error(err)
		return ""
	}
	var remarkBytes [common.XDAG_FIELD_SIZE]byte
	if len(remark) > 0 {
		if ValidateRemark(remark) {
			copy(remarkBytes[:], remark)
		} else {
			xlog.Error("remark error")
			return ""
		}
	}

	var valBytes [8]byte
	if value > 0.0 {
		transVal := xdagoUtils.Xdag2Amount(value)
		binary.LittleEndian.PutUint64(valBytes[:], transVal)
	} else {
		xlog.Error("transaction value is zero")
		return ""
	}

	t := xdagoUtils.GetCurrentTimestamp()
	var timeBytes [8]byte
	binary.LittleEndian.PutUint64(timeBytes[:], t)

	var sb strings.Builder
	// header: transport
	sb.WriteString("0000000000000000")

	// header: field types
	compKey := key.PubKey().SerializeCompressed()
	sb.WriteString(fieldTypes(config.GetConfig().Option.IsTestNet,
		len(remark) > 0, compKey[0] == secp256k1.PubKeyFormatCompressedEven))

	// header: timestamp
	sb.WriteString(hex.EncodeToString(timeBytes[:]))
	// header: fee
	sb.WriteString("0000000000000000")

	// input field: input address
	sb.WriteString(inAddress)
	// input field: input value
	sb.WriteString(hex.EncodeToString(valBytes[:]))
	// output field: output address
	sb.WriteString(outAddress)
	// output field: out value
	sb.WriteString(hex.EncodeToString(valBytes[:]))
	// remark field
	if len(remark) > 0 {
		sb.WriteString(hex.EncodeToString(remarkBytes[:]))
	}
	// public key field
	sb.WriteString(hex.EncodeToString(compKey[1:33]))

	r, s := transactionSign(sb.String(), key, len(remark) > 0)
	// sign field: sign_r
	sb.WriteString(r)
	// sign field: sign_s
	sb.WriteString(s)
	// zero fields
	if len(remark) > 0 {
		for i := 0; i < 18; i++ {
			sb.WriteString("00000000000000000000000000000000")
		}
	} else {
		for i := 0; i < 20; i++ {
			sb.WriteString("00000000000000000000000000000000")
		}
	}
	return sb.String()
}

func checkBase58Address(address string) (string, error) {
	addrBytes, _, err := base58.ChkDec(address)
	if err != nil {
		xlog.Error(err)
		return "", err
	}
	if len(addrBytes) != 24 {
		xlog.Error("transaction receive address length error")
		return "", errors.New("transaction receive address length error")
	}
	return hex.EncodeToString(addrBytes[:]), nil
}

func transactionSign(block string, key *secp256k1.PrivateKey, hasRemark bool) (string, string) {
	var sb strings.Builder
	sb.WriteString(block)
	if hasRemark {
		for i := 0; i < 18; i++ {
			sb.WriteString("00000000000000000000000000000000")
		}
	} else {
		for i := 0; i < 20; i++ {
			sb.WriteString("00000000000000000000000000000000")
		}
	}
	pubKey := key.PubKey().SerializeCompressed()
	sb.WriteString(hex.EncodeToString(pubKey[:]))

	b, _ := hex.DecodeString(sb.String())

	hash := sha256.Sum256(b)
	hash = sha256.Sum256(hash[:])

	r, s := cryptography.EcdsaSign(key, hash[:])

	return hex.EncodeToString(r[:]), hex.EncodeToString(s[:])
}

func fieldTypes(isTest, hasRemark, isPubKeyEven bool) string {
	// 1/8--2--3--[9]--6/7--5--5
	// header(main/test)--input--output--[remark]--pubKey(even/odd)--sign_r--sign_s
	var sb strings.Builder
	if isTest {
		sb.WriteString("28") // test net
	} else {
		sb.WriteString("21") // main net
	}

	if hasRemark { // with remark
		if isPubKeyEven {
			sb.WriteString("93560500000000") // even public key
		} else {
			sb.WriteString("93570500000000") // odd public key
		}
	} else { // without remark
		if isPubKeyEven {
			sb.WriteString("63550000000000") // even public key
		} else {
			sb.WriteString("73550000000000") // odd public key
		}
	}

	return sb.String()
}

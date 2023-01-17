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
	"encoding/hex"
	"fmt"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"goXdagWallet/wallet_state"
	"goXdagWallet/xlog"
	"os"
	"path"
	"strings"
	"time"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

var chanBalance = make(chan int, 1)
var regDone = make(chan int, 1)

func Xdag_Wallet_fount() int {
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	pathName := path.Join(pwd, "xdagj_dat", "dnet_key.dat")
	// change current working directory
	os.Chdir(pwd)

	fi, err := os.Stat(pathName)
	if err != nil {
		return -1
	}
	if fi.Size() != 2048 {
		return -1
	}
	return 0

	//res := C.xdag_dnet_crpt_found()
	//return int(res)
}
func ConnectWallet() {
	C.init_event_callback()
	C.init_password_callback()

	pa := C.CString(config.GetConfig().Option.PoolAddress)
	defer C.free(unsafe.Pointer(pa))
	var testnet int
	if config.GetConfig().Option.IsTestNet {
		testnet = 1
	}

	argv := make([]*C.char, 1)

	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	cs := C.CString(path.Join(pwd, "xdag.exe"))

	defer C.free(unsafe.Pointer(cs))

	argv[0] = cs
	result := C.xdag_init_wrap(C.int(1), (**C.char)(unsafe.Pointer(&argv[0])), pa, C.int(testnet))
	fmt.Println((int32)(result))

}

//export goEventCallback
func goEventCallback(obj unsafe.Pointer, xdagEvent *C.xdag_event) C.int {
	eventId := xdagEvent.event_id
	errCode := xdagEvent.error_no
	eventData := C.GoString(xdagEvent.event_data)
	fmt.Println(int(eventId))
	fmt.Println(eventData)

	switch eventId {
	case C.event_id_log:
		//fmt.Println("event_id_log")
		if int(errCode) > 0x3000 && int(errCode) < 0x4000 {
			setTransferError(eventData)
		}
		xlog.Trace(eventData)
		break
	case C.event_id_state_change:
		//fmt.Println("event_id_state_change")
		state, ok := wallet_state.MessageToState(eventData)
		if ok && state == wallet_state.LoadingBlocks {
			regDone <- 1
			StatusInfo.Text = wallet_state.Localize(state)
			canvas.Refresh(StatusInfo)
		} else if ok && state != wallet_state.TransferPending {
			StatusInfo.Text = wallet_state.Localize(state)
			canvas.Refresh(StatusInfo)
		} else if ok && state == wallet_state.TransferPending {
			TransStatus.Text = wallet_state.Localize(state)
			DonaTransStatus.Text = wallet_state.Localize(state)
		}
		xlog.Info(eventData)
		if (!config.GetConfig().Option.IsTestNet && strings.Contains(eventData, "Connected to the mainnet pool")) ||
			(config.GetConfig().Option.IsTestNet && strings.Contains(eventData, "Connected to the testnet pool")) {
			C.xdag_get_address_wrap()
			C.xdag_get_balance_wrap()
		}
		break
	case C.event_id_state_done:
		xlog.Info(eventData)
		//fmt.Println("event_id_state_done")
		break
	case C.event_id_address_done:
		Address = eventData
		xlog.Info(eventData)
		//fmt.Println("event_id_address_done")
		break
	case C.event_id_balance_done:
		if Balance != eventData {
			Balance = eventData
			AccountBalance.Set(Balance)
			TransStatus.Text = ""
			DonaTransStatus.Text = ""
		}
		NewWalletWindow()
		xlog.Info(eventData)
		//fmt.Println("event_id_balance_done")
		break
	case C.event_id_xfer_done:
		fmt.Println("event_id_xfer_done")
		setTransferDone()
		xlog.Info(eventData)
		break
	case C.event_id_err_exit:
		//fmt.Println("event_id_err_exit")
		xlog.Error(eventData)

		if int(errCode) == 0x1002 { // password incorrect
			StatusInfo.Text = i18n.GetString("Message_PasswordIncorrect")
			canvas.Refresh(StatusInfo)
			//WalletApp.SendNotification(&fyne.Notification{
			//	Title:   i18n.GetString("WalletWindow_Title"),
			//	Content: i18n.GetString("Message_PasswordIncorrect"),
			//})
			time.Sleep(time.Second * 2)
		}
		C.xdag_exit_wrap()
		os.Exit(1)
	default:
		break
	}

	return C.int(0)
}

//export goPasswordCallback
func goPasswordCallback(prompt *C.cchar_t, buf *C.char, size C.uint) C.int {
	C.memcpy(unsafe.Pointer(buf), unsafe.Pointer(&Password[0]), C.size_t(size))
	return C.int(0)
}

func TransferWrap(address, amount, remark string) int {
	// TODO: validate address, amount, remark

	csAddress := C.CString(address)
	defer C.free(unsafe.Pointer(csAddress))

	csAmount := C.CString(amount)
	defer C.free(unsafe.Pointer(csAmount))

	csRemark := C.CString(remark)
	defer C.free(unsafe.Pointer(csRemark))

	result := C.xdag_transfer_wrap(csAddress, csAmount, csRemark)
	fmt.Println(int(result))
	if int(result) == 0 && address != CommunityAddress {
		config.InsertAddress(address)
	}
	return int(result)
}

func ValidateAddress(address string) bool {
	csAddress := C.CString(address)
	defer C.free(unsafe.Pointer(csAddress))

	res := C.xdag_is_valid_wallet_address(csAddress)
	if res == C.int(0) {
		return true
	} else {
		return false
	}
}

func ValidateRemark(remark string) bool {
	csRemark := C.CString(remark)
	defer C.free(unsafe.Pointer(csRemark))

	res := C.xdag_is_valid_remark(csRemark)
	if res == C.int(0) {
		return true
	} else {
		return false
	}
}

func NewWalletWindow() {
	if WalletWindow != nil {
		return
	}
	k := getDefaultPubKey()
	if k == nil {
		fmt.Println("get default key failed.")
	}
	LogonWindow.Win.Hide()
	w := WalletApp.NewWindow(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version) +
		getTestTitle())
	WalletWindow = w
	w.SetMaster()
	LogonWindow.Win.Content().Resize(fyne.NewSize(0, 0))
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabAccount"),
			theme.HomeIcon(), AccountPage(Address, Balance, WalletWindow)),
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
			C.xdag_get_balance_wrap()
		}
	}
}

// get xdag wallet private key
func getDefaultPubKey() []byte {
	p := C.xdag_get_default_key()
	if uintptr(p) > 0 {
		key := C.GoBytes(p, 32)
		fmt.Println(hex.EncodeToString(key[:]))
		xlog.Info("default private key:", hex.EncodeToString(key[:]))
		return key
	}
	return nil
}

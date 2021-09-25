// Package main provides various examples of Fyne API capabilities.
package main

//#cgo linux LDFLAGS: -L${SRCDIR}/lib -lxdag_runtime -L/usr/lib -lsecp256k1 -lssl -lcrypto -lm
//#include "src/xdag_runtime.h"
//#include "callback.h"
//#include <stdlib.h>
//#include <string.h>
/*
 typedef const char cchar_t;
*/
import "C"

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"goXdagWallet/components"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"goXdagWallet/wallet_state"
	"goXdagWallet/xlog"
	"os"
	"strings"
	"unsafe"
)

var LogonWindow LogonWin
var WalletWindow fyne.Window
var WalletApp fyne.App
var globalPassword [256]byte
var globalAddress string
var globalBalance string

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func main() {
	config.InitConfig()
	if i18n.LoadI18nStrings() != nil {
		return
	}

	xlog.SetLogFile("./", "go_wallet.log")
	hasAccount := C.xdag_dnet_crpt_found() // cgo call xdag_runtime C library
	WalletApp = app.NewWithID("go.xdag.wallet")
	WalletApp.SetIcon(resourceIconPng)
	LogonWindow.NewLogonWindow(int(hasAccount))

	themes := config.GetConfig().Theme
	if themes != "Dark" && themes != "Light" {
		WalletApp.Settings().SetTheme(theme.DarkTheme())
	} else if themes == "Dark" {
		WalletApp.Settings().SetTheme(theme.DarkTheme())
	} else {
		WalletApp.Settings().SetTheme(theme.LightTheme())
	}
	LogonWindow.Win.ShowAndRun()
}

func walletWindow() {
	if WalletWindow != nil {
		return
	}

	LogonWindow.Win.Hide()
	w := WalletApp.NewWindow(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version))
	WalletWindow = w
	w.SetMaster()
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabAccount"),
			theme.HomeIcon(), components.AccountPage(globalAddress, globalBalance, WalletWindow)),
		container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabTransfer"),
			theme.MailSendIcon(), components.TransferPage(WalletWindow, transferWrap)),
		container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabHistory"),
			theme.ContentPasteIcon(), components.HistoryPage(WalletWindow)),
		container.NewTabItemWithIcon(i18n.GetString("WalletWindow_TabSettings"),
			theme.SettingsIcon(), components.SettingsPage(WalletWindow)))
	if fyne.CurrentDevice().IsMobile() {
		tabs.SetTabLocation(container.TabLocationBottom)
	} else {
		tabs.SetTabLocation(container.TabLocationLeading)
	}

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(640, 480))
	w.CenterOnScreen()
	w.SetOnClosed(func() {
		WalletApp.Quit()
		os.Exit(0)
	})
	w.Show()
}

func connectWallet() {
	C.init_event_callback()
	C.init_password_callback()

	pa := C.CString(config.GetConfig().Option.PoolAddress)
	defer C.free(unsafe.Pointer(pa))

	argv := make([]*C.char, 1)
	cs := C.CString("xdag.exe")
	defer C.free(unsafe.Pointer(cs))

	argv[0] = cs
	result := C.xdag_init_wrap(C.int(1), (**C.char)(unsafe.Pointer(&argv[0])), pa)
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
		xlog.Trace(eventData)
		break
	case C.event_id_state_change:
		//fmt.Println("event_id_state_change")
		state, ok := wallet_state.MessageToState(eventData)
		if ok && state != wallet_state.TransferPending {
			LogonWindow.StatusInfo.Text = wallet_state.Localize(state)
		}
		xlog.Info(eventData)
		if strings.Contains(eventData, "Connected to the mainnet pool") {
			C.xdag_get_address_wrap()
			C.xdag_get_balance_wrap()
		}
		break
	case C.event_id_state_done:
		xlog.Info(eventData)
		//fmt.Println("event_id_state_done")
		break
	case C.event_id_address_done:
		globalAddress = eventData
		xlog.Info(eventData)
		//fmt.Println("event_id_address_done")
		break
	case C.event_id_balance_done:
		if globalBalance != eventData {
			globalBalance = eventData
		}
		walletWindow()
		xlog.Info(eventData)
		//fmt.Println("event_id_balance_done")
		break
	case C.event_id_xfer_done:
		fmt.Println("event_id_xfer_done")
		xlog.Info(eventData)
		break
	case C.event_id_err_exit:
		//fmt.Println("event_id_err_exit")
		if int(errCode) == 0x1002 { // password incorrect
			LogonWindow.WrongPassword()
		}
		xlog.Error(eventData)
		break
	default:
		break
	}

	return C.int(0)
}

//export goPasswordCallback
func goPasswordCallback(prompt *C.cchar_t, buf *C.char, size C.uint) C.int {
	C.memcpy(unsafe.Pointer(buf), unsafe.Pointer(&globalPassword[0]), C.size_t(size))
	return C.int(0)
}

func transferWrap(address, amount, remark string) int {
	// TODO: validate address, amount, remark

	csAddress := C.CString(address)
	defer C.free(unsafe.Pointer(csAddress))

	csAmount := C.CString(amount)
	defer C.free(unsafe.Pointer(csAmount))

	csRemark := C.CString(remark)
	defer C.free(unsafe.Pointer(csRemark))

	result := C.xdag_transfer_wrap(csAddress, csAmount, csRemark)
	fmt.Println(int(result))
	return int(result)
}

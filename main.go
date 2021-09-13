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
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"goXdagWallet/components"
	"goXdagWallet/xlog"
	"os"
	"strings"
	"unsafe"
)

var WalletWindow fyne.Window
var WalletApp fyne.App
var globalPassword [256]byte

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func main() {
	initConfig()
	if LoadI18nStrings() != nil {
		return
	}

	xlog.SetLogFile("./", "go_wallet.log")
	hasAccount := C.xdag_dnet_crpt_found() // cgo call xdag_runtime C library
	WalletApp = app.NewWithID("go.xdag.wallet")
	WalletApp.SetIcon(resourceIconPng)
	w := WalletApp.NewWindow(fmt.Sprintf(GetI18nString("LogonWindow_Title"), conf.Version))
	WalletWindow = w

	//w.SetMaster()
	var btn *widget.Button
	if hasAccount == 0 { // found wallet key file
		btn = widget.NewButton(GetI18nString("LogonWindow_ConnectWallet"), connectAccount)
	} else if hasAccount == -1 { // not fount
		btn = widget.NewButton(GetI18nString("LogonWindow_RegisterWallet"), walletWindow)
	}

	btn.Importance = widget.HighImportance
	image := canvas.NewImageFromResource(resourceLogonPng)
	content := container.New(layout.NewMaxLayout(), image,
		container.New(layout.NewVBoxLayout(), layout.NewSpacer(),
			layout.NewSpacer(), layout.NewSpacer(),
			layout.NewSpacer(), layout.NewSpacer(),
			container.New(layout.NewPaddedLayout(), btn),
			layout.NewSpacer()))
	w.SetContent(content)
	w.Resize(fyne.NewSize(450, 300))
	w.CenterOnScreen()
	w.SetOnClosed(func() {
		WalletApp.Quit()
		os.Exit(0)
	})

	themes := GetConfig().Theme
	if themes != "Dark" && themes != "Light" {
		WalletApp.Settings().SetTheme(theme.DarkTheme())
	} else if themes == "Dark" {
		WalletApp.Settings().SetTheme(theme.DarkTheme())
	} else {
		WalletApp.Settings().SetTheme(theme.LightTheme())
	}
	w.ShowAndRun()
}

func connectAccount() {
	if conf.Option.PoolAddress == "" {
		dialog.ShowInformation(GetI18nString("Common_MessageTitle"),
			GetI18nString("LogonWindow_NoPoolAddress"), WalletWindow)
		return
	}

	//pwdDialog := components.NewPasswordDialog(nil)
	//pwdDialog.ShowPasswordDialog(GetI18nString("Common_MessageTitle"),
	//	GetI18nString("PasswordWindow_InputPassword"),
	//	GetI18nString("Common_Confirm"),
	//	GetI18nString("Common_Cancel"), nil, WalletWindow
	//	)
	walletWindow()
	C.init_event_callback()
	C.init_password_callback()

	pa := C.CString(conf.Option.PoolAddress)
	defer C.free(unsafe.Pointer(pa))

	argv := make([]*C.char, 1)
	cs := C.CString("xdag.exe")
	defer C.free(unsafe.Pointer(cs))

	argv[0] = cs
	result := C.xdag_init_wrap(C.int(1), (**C.char)(unsafe.Pointer(&argv[0])), pa)
	fmt.Println((int32)(result))
	//switch result {
	//
	//}
}

func walletWindow() {
	WalletWindow.Hide()
	w := WalletApp.NewWindow(fmt.Sprintf(GetI18nString("LogonWindow_Title"), conf.Version))
	WalletWindow = w
	w.SetMaster()
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(GetI18nString("WalletWindow_TabAccount"),
			theme.HomeIcon(), components.AccountPage(WalletWindow)),
		container.NewTabItemWithIcon(GetI18nString("WalletWindow_TabTransfer"),
			theme.MailSendIcon(), components.TransferPage(WalletWindow)),
		container.NewTabItemWithIcon(GetI18nString("WalletWindow_TabHistory"),
			theme.ContentPasteIcon(), components.HistoryPage(WalletWindow)),
		container.NewTabItemWithIcon(GetI18nString("WalletWindow_TabSettings"),
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

//export goEventCallback
func goEventCallback(obj unsafe.Pointer, xdagEvent *C.xdag_event) C.int {
	eventId := xdagEvent.event_id
	eventData := C.GoString(xdagEvent.event_data)
	//fmt.Println(eventData)

	switch eventId {
	case C.event_id_log:
		//fmt.Println("event_id_log")
		xlog.Trace(eventData)
		break
	case C.event_id_state_change:
		//fmt.Println("event_id_state_change")
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
		xlog.Info(eventData)
		//fmt.Println("event_id_address_done")
		break
	case C.event_id_balance_done:
		xlog.Info(eventData)
		//fmt.Println("event_id_balance_done")
		break
	case C.event_id_err_exit:
		//fmt.Println("event_id_err_exit")
		xlog.Error(eventData)
		break
	default:
		break
	}

	return C.int(0)
}

//export goPasswordCallback
func goPasswordCallback(prompt *C.cchar_t, buf *C.char, size C.uint) C.int {
	for i := range globalPassword {
		globalPassword[i] = 0
	}
	pwd := "ljr20040224"
	copy(globalPassword[:], pwd)

	C.memcpy(unsafe.Pointer(buf), unsafe.Pointer(&globalPassword[0]), C.size_t(size))
	return C.int(0)
}

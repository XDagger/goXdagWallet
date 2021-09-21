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
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"goXdagWallet/xlog"
	"os"
	"strings"
	"unsafe"
)

var WalletWindow fyne.Window
var LogonWindow fyne.Window
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
	w := WalletApp.NewWindow(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version))
	LogonWindow = w

	var btn *widget.Button
	if hasAccount == 0 { // found wallet key file
		btn = widget.NewButton(i18n.GetString("LogonWindow_ConnectWallet"), connectClick)
	} else if hasAccount == -1 { // not fount
		btn = widget.NewButton(i18n.GetString("LogonWindow_RegisterWallet"), walletWindow)
	}
	btn.Importance = widget.HighImportance

	settingBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		showLanguageDialog(i18n.GetString("SettingsWindow_ChooseLanguage"),
			i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), func(lang string) {
				config.GetConfig().CultureInfo = lang
				config.SaveConfig()
				i18n.LoadI18nStrings()
				LogonWindow.SetTitle(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version))
				if hasAccount == 0 { // found wallet key file
					btn.SetText(i18n.GetString("LogonWindow_ConnectWallet"))
				} else if hasAccount == -1 { // not fount
					btn.SetText(i18n.GetString("LogonWindow_RegisterWallet"))
				}
			}, LogonWindow)
	})
	settingBtn.Resize(fyne.NewSize(20, 20))
	settingBtn.Importance = widget.HighImportance

	image := canvas.NewImageFromResource(resourceLogonPng)
	content := container.New(layout.NewMaxLayout(), image,
		container.New(layout.NewVBoxLayout(),
			container.New(layout.NewHBoxLayout(), layout.NewSpacer(), settingBtn),
			layout.NewSpacer(), layout.NewSpacer(),
			layout.NewSpacer(), layout.NewSpacer(),
			container.New(layout.NewPaddedLayout(), btn),
			layout.NewSpacer()))
	w.SetContent(content)
	w.Resize(fyne.NewSize(410, 300))
	w.CenterOnScreen()
	w.SetOnClosed(func() {
		WalletApp.Quit()
		os.Exit(0)
	})

	themes := config.GetConfig().Theme
	if themes != "Dark" && themes != "Light" {
		WalletApp.Settings().SetTheme(theme.DarkTheme())
	} else if themes == "Dark" {
		WalletApp.Settings().SetTheme(theme.DarkTheme())
	} else {
		WalletApp.Settings().SetTheme(theme.LightTheme())
	}
	w.ShowAndRun()
}

func connectClick() {
	if config.GetConfig().Option.PoolAddress == "" {
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("LogonWindow_NoPoolAddress"), LogonWindow)
		return
	}

	showPasswordDialog(i18n.GetString("PasswordWindow_InputPassword"),
		i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), LogonWindow)

	//switch result {
	//
	//}
}

func walletWindow() {
	if WalletWindow != nil {
		return
	}

	LogonWindow.Hide()
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

func showPasswordDialog(title, ok, dismiss string, parent fyne.Window) {
	wgt := widget.NewEntry()
	wgt.Password = true

	dialog.ShowCustomConfirm(title, ok, dismiss, wgt, func(b bool) {
		for i := range globalPassword {
			globalPassword[i] = 0
		}
		str := wgt.Text
		fmt.Println("b", b, "pwd", str)
		if b && len(str) > 0 {
			copy(globalPassword[:], str)
			connectWallet()
		}
	}, parent)

}

func showLanguageDialog(title, ok, dismiss string, callback func(string), parent fyne.Window) {
	lang := "en-US"
	radio := widget.NewRadioGroup([]string{"English", "中文"}, func(value string) {
		if value == "English" {
			lang = "en-US"
		} else {
			lang = "zh-CN"
		}
	})

	dialog.ShowCustomConfirm(title, ok, dismiss, radio, func(b bool) {
		if b {
			callback(lang)
		}
	}, parent)

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
		xlog.Error(eventData)
		break
	default:
		break
	}

	return C.int(0)
}

//export goPasswordCallback
func goPasswordCallback(prompt *C.cchar_t, buf *C.char, size C.uint) C.int {
	//for i := range globalPassword {
	//	globalPassword[i] = 0
	//}
	//pwd := "ljr20040224"
	//copy(globalPassword[:], pwd)

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

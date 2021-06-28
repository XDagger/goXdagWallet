// Package main provides various examples of Fyne API capabilities.
package main

//#cgo LDFLAGS: -L${SRCDIR}/lib -lxdag_runtime -L/usr/lib -lsecp256k1 -lssl -lcrypto -lm
//#include "src/xdag_runtime.h"
import "C"

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"goXdagWallet/components"
)

var topWindow fyne.Window

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

	hasAccount:= C.xdag_dnet_crpt_found() // cgo call xdag_runtime C library
	a := app.NewWithID("go.xdag.wallet")
	a.SetIcon(resourceIconPng)
	w := a.NewWindow(fmt.Sprintf(GetI18nString("LogonWindow_Title"), conf.Version))
	topWindow = w

	w.SetMaster()
	var btn *widget.Button
	if hasAccount == 0 {			// found wallet key file
		btn = widget.NewButton("connect wallet", walletWindow)
	} else if hasAccount == -1 {	 // not fount
		btn = widget.NewButton("create wallet", walletWindow)
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
	w.Resize(fyne.NewSize(420, 300))
	w.CenterOnScreen()

	themes := GetConfig().Theme
	if themes != "Dark" && themes != "Light" {
		a.Settings().SetTheme(theme.DarkTheme())
	} else if themes == "Dark" {
		a.Settings().SetTheme(theme.DarkTheme())
	} else {
		a.Settings().SetTheme(theme.LightTheme())
	}
	w.ShowAndRun()
}

func walletWindow() {
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(GetI18nString("WalletWindow_TabAccount"),
			theme.HomeIcon(), components.AccountPage(topWindow)),
		container.NewTabItemWithIcon(GetI18nString("WalletWindow_TabTransfer"),
			theme.MailSendIcon(), components.TransferPage(topWindow)),
		container.NewTabItemWithIcon(GetI18nString("WalletWindow_TabHistory"),
			theme.ContentPasteIcon(), components.HistoryPage(topWindow)),
		container.NewTabItemWithIcon(GetI18nString("WalletWindow_TabSettings"),
			theme.SettingsIcon(), components.SettingsPage(topWindow)))
	if fyne.CurrentDevice().IsMobile() {
		tabs.SetTabLocation(container.TabLocationBottom)
	} else {
		tabs.SetTabLocation(container.TabLocationLeading)
	}

	topWindow.SetContent(tabs)
	topWindow.Resize(fyne.NewSize(640, 480))
	topWindow.CenterOnScreen()
}
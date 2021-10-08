// Package main provides various examples of Fyne API capabilities.
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"goXdagWallet/components"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"goXdagWallet/xlog"
)

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
	hasAccount := components.Xdag_Wallet_fount() // cgo call xdag_runtime C library
	components.WalletApp = app.NewWithID("go.xdag.wallet")
	components.WalletApp.SetIcon(components.GetAppIcon())
	components.LogonWindow.NewLogonWindow(hasAccount)

	themes := config.GetConfig().Theme
	if themes != "Dark" && themes != "Light" {
		components.WalletApp.Settings().SetTheme(theme.DarkTheme())
	} else if themes == "Dark" {
		components.WalletApp.Settings().SetTheme(theme.DarkTheme())
	} else {
		components.WalletApp.Settings().SetTheme(theme.LightTheme())
	}
	components.LogonWindow.Win.ShowAndRun()
}

// Package main provides various examples of Fyne API capabilities.
package main

import (
	"goXdagWallet/components"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"goXdagWallet/xlog"
	"os"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func init() {
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)

	os.Setenv("FYNE_FONT", path.Join(pwd, "data", "myFont.ttf"))

}

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func main() {
	xlog.SetLogFile("./", "go_wallet.log")
	config.InitConfig()
	if err := i18n.LoadI18nStrings(); err != nil {
		xlog.Error(err)
		return
	}

	hasAccount := components.Xdag_Wallet_fount() // cgo call xdag_runtime C library
	components.WalletApp = app.NewWithID("io.xdagj.wallet")
	components.WalletApp.SetIcon(components.GetAppIcon())
	components.LogonWindow.NewLogonWindow(hasAccount)
	components.LogonWindow.Win.ShowAndRun()
	os.Unsetenv("FYNE_FONT")
}

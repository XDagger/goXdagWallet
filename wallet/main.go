// Package main provides various examples of Fyne API capabilities.
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"goXdagWallet/components"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"goXdagWallet/xlog"
	"os"
	"path"
)

var splashDone = make(chan struct{})

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
	if components.ShowSplashWindow(splashDone) {
		go func() {
			for range splashDone {
				components.LogonWindow.NewLogonWindow(hasAccount)
				components.LogonWindow.Win.Show()
				splashDone <- struct{}{}
			}
		}()
	} else {
		components.LogonWindow.NewLogonWindow(hasAccount)
		components.LogonWindow.Win.Show()
	}
	components.WalletApp.Run()
	os.Unsetenv("FYNE_FONT")
}

// Package main provides various examples of Fyne API capabilities.
package main

import (
	"flag"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"goXdagWallet/cli"
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
	accountStatus := components.Xdag_Wallet_fount() // search wallet data files

	mode := flag.String("mode", "gui", "run mode: gui, cli, server")
	flag.Parse()

	if *mode == "gui" {
		if err := i18n.LoadI18nStrings(); err != nil {
			xlog.Error(err)
			return
		}

		components.WalletApp = app.NewWithID("io.xdagj.wallet")
		components.WalletApp.SetIcon(components.GetAppIcon())
		if components.ShowSplashWindow(splashDone) {
			go func() {
				for range splashDone {
					components.LogonWindow.NewLogonWindow(accountStatus)
					components.LogonWindow.Win.Show()
					splashDone <- struct{}{}
				}
			}()
		} else {
			components.LogonWindow.NewLogonWindow(accountStatus)
			components.LogonWindow.Win.Show()
		}
		components.WalletApp.Run()
	} else if *mode == "cli" {
		cli.NewCli(accountStatus)
	}
	os.Unsetenv("FYNE_FONT")
}

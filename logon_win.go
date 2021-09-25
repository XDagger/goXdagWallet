package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"image/color"
	"os"
)

type LogonWin struct {
	Win               fyne.Window
	BtnContainer      *fyne.Container
	ProgressContainer *fyne.Container
	StatusInfo        *canvas.Text
}

func (l *LogonWin) NewLogonWindow(hasAccount int) {
	w := WalletApp.NewWindow(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version))
	l.Win = w

	var btn *widget.Button
	if hasAccount == 0 { // found wallet key file
		btn = widget.NewButton(i18n.GetString("LogonWindow_ConnectWallet"), l.connectClick)
	} else if hasAccount == -1 { // not fount
		btn = widget.NewButton(i18n.GetString("LogonWindow_RegisterWallet"), walletWindow)
	}
	btn.Importance = widget.HighImportance

	l.StatusInfo = canvas.NewText("", color.White)
	progress := widget.NewProgressBarInfinite()
	l.BtnContainer = container.New(layout.NewPaddedLayout(), btn)
	l.ProgressContainer = container.New(layout.NewPaddedLayout(), progress)
	l.ProgressContainer.Hide()

	settingBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		showLanguageDialog(i18n.GetString("SettingsWindow_ChooseLanguage"),
			i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), func(lang string) {
				config.GetConfig().CultureInfo = lang
				config.SaveConfig()
				i18n.LoadI18nStrings()
				l.Win.SetTitle(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version))
				if hasAccount == 0 { // found wallet key file
					btn.SetText(i18n.GetString("LogonWindow_ConnectWallet"))
				} else if hasAccount == -1 { // not fount
					btn.SetText(i18n.GetString("LogonWindow_RegisterWallet"))
				}
			}, l.Win)
	})
	settingBtn.Resize(fyne.NewSize(20, 20))
	settingBtn.Importance = widget.HighImportance

	image := canvas.NewImageFromResource(resourceLogonPng)
	content := container.New(layout.NewMaxLayout(), image,
		container.New(layout.NewVBoxLayout(),
			container.New(layout.NewHBoxLayout(), layout.NewSpacer(), settingBtn),
			layout.NewSpacer(), layout.NewSpacer(),
			layout.NewSpacer(),
			container.New(layout.NewHBoxLayout(), layout.NewSpacer(), l.StatusInfo, layout.NewSpacer()),
			l.BtnContainer,
			l.ProgressContainer,
			layout.NewSpacer()))
	w.SetContent(content)
	w.Resize(fyne.NewSize(410, 300))
	w.CenterOnScreen()
	w.SetOnClosed(func() {
		WalletApp.Quit()
		os.Exit(0)
	})
}

func (l *LogonWin) StartConnect() {
	l.BtnContainer.Hide()
	l.ProgressContainer.Show()
	l.StatusInfo.Text = i18n.GetString("LogonWindow_ConnectingAccount")
}

func (l *LogonWin) WrongPassword() {
	l.BtnContainer.Show()
	l.ProgressContainer.Hide()
	l.StatusInfo.Text = ""
	dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
		i18n.GetString("Message_PasswordIncorrect"), l.Win)
}

func (l *LogonWin) connectClick() {
	if config.GetConfig().Option.PoolAddress == "" {
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("LogonWindow_NoPoolAddress"), l.Win)
		return
	}
	l.showPasswordDialog(i18n.GetString("PasswordWindow_InputPassword"),
		i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), l.Win)
}

func (l *LogonWin) showPasswordDialog(title, ok, dismiss string, parent fyne.Window) {
	wgt := widget.NewEntry()
	wgt.Password = true

	dialog.ShowCustomConfirm(title, ok, dismiss, wgt, func(b bool) {
		for i := range globalPassword {
			globalPassword[i] = 0
		}
		str := wgt.Text
		if b && len(str) > 0 {
			copy(globalPassword[:], str)
			l.StartConnect()
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

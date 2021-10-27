package components

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"image/color"
	"os"
	"time"
)

type LogonWin struct {
	Win               fyne.Window
	BtnContainer      *fyne.Container
	ProgressContainer *fyne.Container
	HasAccount        bool
	Password          string
}

var StatusInfo = canvas.NewText("", color.White)

func (l *LogonWin) NewLogonWindow(hasAccount int) {
	w := WalletApp.NewWindow(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version))
	l.Win = w

	var btn *widget.Button
	l.HasAccount = hasAccount == 0
	if hasAccount == 0 { // found wallet key file
		btn = widget.NewButton(i18n.GetString("LogonWindow_ConnectWallet"), l.connectClick)
	} else if hasAccount == -1 { // not fount
		btn = widget.NewButton(i18n.GetString("LogonWindow_RegisterWallet"), l.connectClick)
	}
	btn.Importance = widget.HighImportance

	//l.StatusInfo = canvas.NewText("", color.White)
	StatusInfo.Alignment = fyne.TextAlignCenter
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

	appearanceBtn := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		fyneSetting := WalletApp.NewWindow(i18n.GetString("LogonWindow_AppearanceSetting"))
		fyneSetting.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		fyneSetting.Resize(fyne.NewSize(480, 480))
		fyneSetting.Show()
	})
	appearanceBtn.Resize(fyne.NewSize(20, 20))
	appearanceBtn.Importance = widget.HighImportance

	image := canvas.NewImageFromResource(resourceDaggerLogoCPng)
	content := container.New(layout.NewMaxLayout(), image,
		container.New(layout.NewVBoxLayout(),
			container.New(layout.NewHBoxLayout(), layout.NewSpacer(), appearanceBtn, settingBtn),
			layout.NewSpacer(), layout.NewSpacer(),
			layout.NewSpacer(),
			//container.New(layout.NewHBoxLayout(), layout.NewSpacer(), StatusInfo, layout.NewSpacer()),
			StatusInfo,
			l.BtnContainer,
			l.ProgressContainer,
			layout.NewSpacer()))
	w.SetContent(content)
	w.Resize(fyne.NewSize(410, 300))
	w.SetOnClosed(func() {
		WalletApp.Quit()
		os.Exit(0)
	})
	w.CenterOnScreen()
	settingBtn.Refresh()
	appearanceBtn.Refresh()
}

func (l *LogonWin) StartConnect() {
	l.BtnContainer.Hide()
	l.ProgressContainer.Show()
	StatusInfo.Text = i18n.GetString("LogonWindow_ConnectingAccount")
	canvas.Refresh(StatusInfo)
}
func (l *LogonWin) StartRegister() {
	l.BtnContainer.Hide()
	l.ProgressContainer.Show()
	StatusInfo.Text = i18n.GetString("WalletState_Registering")
	canvas.Refresh(StatusInfo)
	go registerTimer()
}

func (l *LogonWin) connectClick() {
	if config.GetConfig().Option.PoolAddress == "" {
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("LogonWindow_NoPoolAddress"), l.Win)
		return
	}
	var title string
	if l.HasAccount {
		title = i18n.GetString("PasswordWindow_InputPassword")
	} else {
		title = i18n.GetString("PasswordWindow_SetPassword")
	}
	l.showPasswordDialog(title,
		i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), l.Win)
}

func (l *LogonWin) showPasswordDialog(title, ok, dismiss string, parent fyne.Window) {
	wgt := widget.NewEntry()
	wgt.Password = true

	dialog.ShowCustomConfirm(title, ok, dismiss, wgt, func(b bool) {
		for i := range Password {
			Password[i] = 0
		}
		str := wgt.Text
		if b && len(str) > 0 {
			if l.HasAccount {
				copy(Password[:], str)
				l.StartConnect()
				ConnectWallet()
			} else {
				l.Password = str
				l.ReShowPasswordDialog(i18n.GetString("PasswordWindow_RetypePassword"),
					i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), l.Win)
			}
		}

	}, parent)
}

func (l *LogonWin) ReShowPasswordDialog(title, ok, dismiss string, parent fyne.Window) {
	wgt := widget.NewEntry()
	wgt.Password = true

	dialog.ShowCustomConfirm(title, ok, dismiss, wgt, func(b bool) {
		for i := range Password {
			Password[i] = 0
		}
		str := wgt.Text
		if b && len(str) > 0 {
			if str == l.Password {
				copy(Password[:], str)
				l.StartRegister()
				ConnectWallet()
			} else {
				dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
					i18n.GetString("PasswordWindow_PasswordMismatch"), l.Win)
			}
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

func GetAppIcon() fyne.Resource {
	return resourceLogoPng
}

func registerTimer() {
	start := time.Now()
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			span := time.Now().Sub(start)
			z := time.Unix(0, 0).UTC()
			StatusInfo.Text = fmt.Sprintf(i18n.GetString("LogonWindow_InitializingElapsedTime"),
				z.Add(span).Format("04:05"))
			canvas.Refresh(StatusInfo)
			break
		case <-regDone:
			return
		default:

		}
	}
}

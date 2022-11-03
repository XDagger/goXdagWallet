package components

import (
	"encoding/json"
	"errors"
	"fmt"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"goXdagWallet/xlog"
	"image/color"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	cp "github.com/otiai10/copy"
)

type LogonWin struct {
	Win               fyne.Window
	BtnContainer      *fyne.Container
	LogonBtn          *widget.Button
	ProgressContainer *fyne.Container
	HasAccount        bool
	Password          string
}

var StatusInfo = canvas.NewText("", color.White)

func (l *LogonWin) NewLogonWindow(hasAccount int) {

	w := WalletApp.NewWindow(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version) +
		getTestTitle())
	l.Win = w

	l.HasAccount = hasAccount == 0
	if hasAccount == 0 { // found wallet key file
		l.LogonBtn = widget.NewButton(i18n.GetString("LogonWindow_ConnectWallet"), l.connectClick)
	} else if hasAccount == -1 { // not fount
		l.LogonBtn = widget.NewButton(i18n.GetString("LogonWindow_RegisterWallet"), l.connectClick)
	}
	l.LogonBtn.Importance = widget.HighImportance

	//l.StatusInfo = canvas.NewText("", color.White)
	StatusInfo.Alignment = fyne.TextAlignCenter
	progress := widget.NewProgressBarInfinite()
	progress.Hide() // for primary color changing
	l.BtnContainer = container.New(layout.NewPaddedLayout(), l.LogonBtn)
	l.ProgressContainer = container.New(layout.NewPaddedLayout(), progress)

	settingBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		showLanguageDialog(i18n.GetString("SettingsWindow_ChooseLanguage"),
			i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), func(lang string) {
				config.GetConfig().CultureInfo = lang
				config.SaveConfig()
				i18n.LoadI18nStrings()
				l.Win.SetTitle(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version) +
					getTestTitle())
				if hasAccount == 0 { // found wallet key file
					l.LogonBtn.SetText(i18n.GetString("LogonWindow_ConnectWallet"))
				} else if hasAccount == -1 { // not fount
					l.LogonBtn.SetText(i18n.GetString("LogonWindow_RegisterWallet"))
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

	image := canvas.NewImageFromResource(resourceXj3Png)
	content := container.New(layout.NewMaxLayout(), image,
		container.New(layout.NewVBoxLayout(),
			container.New(layout.NewHBoxLayout(), layout.NewSpacer(), appearanceBtn, settingBtn),
			layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(),
			layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(),
			layout.NewSpacer(),
			StatusInfo,
			l.BtnContainer,
			l.ProgressContainer,
			layout.NewSpacer()))
	w.SetContent(content)
	w.Resize(fyne.NewSize(410, 305))
	w.SetOnClosed(func() {
		xlog.CleanXdagLog(xlog.StdXdagLog)
		WalletApp.Quit()
		os.Exit(0)
	})
	w.CenterOnScreen()
	settingBtn.Refresh()
	appearanceBtn.Refresh()
}

func (l *LogonWin) StartConnect() {
	l.BtnContainer.Hide()
	l.ProgressContainer.Objects[0].Show() // for primary color changing
	StatusInfo.Text = i18n.GetString("LogonWindow_ConnectingAccount")
	canvas.Refresh(StatusInfo)
}
func (l *LogonWin) StartRegister() {
	l.BtnContainer.Hide()
	l.ProgressContainer.Objects[0].Show() // for primary color changing
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
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)

	if l.HasAccount {
		l.showPasswordDialog(i18n.GetString("PasswordWindow_InputPassword"),
			i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), l.Win)
	} else {

		confirmFrm := dialog.NewConfirm(i18n.GetString("Wallet_Choice"),
			i18n.GetString("Wallet_CreateOrImport"),
			func(b bool) {
				if b {
					dlgOpen := dialog.NewFolderOpen(
						func(uri fyne.ListableURI, err error) {
							defer func() {
								l.Win.Resize(fyne.NewSize(410, 305))
							}()

							if uri == nil || err != nil {
								return
							}
							if !checkOldWallet(uri.Path()) {
								dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
									i18n.GetString("WalletImport_WalletNotExist"), l.Win)
								return
							}
							if l.copyOldWallet(uri.Path()) != nil {
								dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
									i18n.GetString("WalletImport_FilesCopyFailed"), l.Win)
								return
							}
							l.HasAccount = true
							l.showPasswordDialog(i18n.GetString("PasswordWindow_InputPassword"),
								i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), l.Win)

						}, l.Win)
					l.Win.Resize(fyne.NewSize(800, 500))
					dlgOpen.Resize(fyne.NewSize(800, 500))
					dlgOpen.Show()

				} else {

					pathDest := path.Join(pwd, "xdagj_dat")
					if err := os.RemoveAll(pathDest); err != nil {
						dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
							i18n.GetString("WalletImport_FilesCopyFailed"), l.Win)
						return
					}
					if err := os.MkdirAll(pathDest, 0666); err != nil {
						dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
							i18n.GetString("WalletImport_FilesCopyFailed"), l.Win)
						return
					}
					l.showPasswordDialog(i18n.GetString("PasswordWindow_SetPassword"),
						i18n.GetString("Common_Confirm"), i18n.GetString("Common_Cancel"), l.Win)
				}
			}, l.Win)
		confirmFrm.SetConfirmText(i18n.GetString("Wallet_Import"))
		confirmFrm.SetDismissText(i18n.GetString("Wallet_Create"))
		confirmFrm.Show()
	}
}

func (l *LogonWin) showPasswordDialog(title, ok, dismiss string, parent fyne.Window) {
	l.LogonBtn.SetText(i18n.GetString("LogonWindow_ConnectWallet"))
	wgt := widget.NewEntry()
	wgt.Password = true

	dialog.ShowCustomConfirm(title, ok, dismiss, wgt, func(b bool) {
		for i := range Password {
			Password[i] = 0
		}
		str := wgt.Text
		if b {
			if l.HasAccount {
				if len(str) > 0 {
					copy(Password[:], str)
				}
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
		if b {
			if str == l.Password {
				if len(str) > 0 {
					copy(Password[:], str)
				}
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
	radio := widget.NewRadioGroup([]string{"English", "中文", "Français", "Русский"}, func(value string) {
		if value == "English" {
			lang = "en-US"
		} else if value == "Français" {
			lang = "fr-FR"
		} else if value == "Русский" {
			lang = "ru-RU"
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
	return resourceIconPng
}

func registerTimer() {
	start := time.Now()
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			span := time.Since(start)
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

func getTestTitle() string {
	var testNet string
	if config.GetConfig().Option.IsTestNet {
		if config.GetConfig().CultureInfo == "zh-CN" {
			testNet = "测试网"
		} else if config.GetConfig().CultureInfo == "fr-FR" {
			testNet = "Réseau Test"
		} else if config.GetConfig().CultureInfo == "ru-RU" {
			testNet = "тестовая сеть"
		} else {
			testNet = "Test Net"
		}
	}
	return testNet
}

func checkOldWallet(walletDir string) bool {

	pathName := path.Join(walletDir, "dnet_key.dat")
	fi, err := os.Stat(pathName)
	if err != nil {
		return false
	}
	if fi.Size() != 2048 {
		return false
	}

	pathName = path.Join(walletDir, "wallet.dat")
	fi, err = os.Stat(pathName)
	if err != nil {
		return false
	}
	if fi.Size() != 32 {
		return false
	}

	pathName = path.Join(walletDir, "storage")
	_, err = os.Stat(pathName)

	return err == nil
}

func (l *LogonWin) copyOldWallet(walletDir string) error {
	pwd, _ := os.Executable()
	pwd, _ = path.Split(pwd)
	pathDest := path.Join(pwd, "xdagj_dat")
	if err := os.RemoveAll(pathDest); err != nil {
		return err
	}
	if err := os.MkdirAll(pathDest, 0666); err != nil {
		return err
	}
	if err := copyFile(walletDir, pathDest, "dnet_key.dat", 2048); err != nil {
		return err
	}
	if err := copyFile(walletDir, pathDest, "wallet.dat", 32); err != nil {
		return err
	}
	if err := cp.Copy(path.Join(walletDir, "storage"), path.Join(pathDest, "storage")); err != nil {
		return err
	}

	l.importConfig(walletDir)

	return nil
}

func copyFile(walletDir, pathDest, fileName string, n int64) error {
	pathName := path.Join(walletDir, fileName)
	source, err := os.Open(pathName)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(path.Join(pathDest, fileName))
	if err != nil {
		return err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	if err != nil {
		return err
	}
	if nBytes != n {
		return errors.New(fileName + " copy  failed")
	}
	return nil
}

func (l *LogonWin) importConfig(walletDir string) {

	configName := path.Join(walletDir, "wallet-config.json")
	_, err := os.Stat(configName)
	if err != nil {
		walletDir = strings.TrimSuffix(walletDir, `\`)
		walletDir = strings.TrimSuffix(walletDir, `/`)
		configName = path.Join(path.Dir(walletDir), "wallet-config.json")
		_, err = os.Stat(configName)
		if err != nil {
			return
		}
	}
	var oldConf config.Config
	data, err := os.ReadFile(configName)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &oldConf)
	if err != nil {
		return
	}

	if config.GetConfig().CultureInfo == "en-US" && oldConf.CultureInfo != "en-US" {

		config.GetConfig().CultureInfo = oldConf.CultureInfo

		i18n.LoadI18nStrings()

		l.Win.Resize(fyne.NewSize(410, 305))

		l.Win.SetTitle(fmt.Sprintf(i18n.GetString("LogonWindow_Title"), config.GetConfig().Version) +
			getTestTitle())

		l.LogonBtn.SetText(i18n.GetString("LogonWindow_ConnectWallet"))

		l.Win.Resize(fyne.NewSize(410, 305))
	}

	config.GetConfig().Query = oldConf.Query

	for _, a := range oldConf.Addresses {
		config.InsertAddress(a)
	}
}

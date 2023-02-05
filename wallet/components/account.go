package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	qrcode "github.com/skip2/go-qrcode"
	"goXdagWallet/i18n"
)

type myEntry struct {
	widget.Entry
}

func newMyEntry() *myEntry {
	ret := &myEntry{}
	ret.ExtendBaseWidget(ret)
	return ret
}
func newMyEntryWithData(data binding.String) *myEntry {
	ret := &myEntry{}
	ret.Bind(data)
	ret.ExtendBaseWidget(ret)
	return ret
}
func (e *myEntry) MouseDown(_ *desktop.MouseEvent)    {}
func (e *myEntry) MouseUp(_ *desktop.MouseEvent)      {}
func (e *myEntry) Tapped(_ *fyne.PointEvent)          {}
func (e *myEntry) TappedSecondary(_ *fyne.PointEvent) {}
func (e *myEntry) KeyDown(_ *fyne.KeyEvent)           {}
func (e *myEntry) KeyUp(_ *fyne.KeyEvent)             {}

var AccountBalance = binding.NewString()

func AccountPage(address, balance string, w fyne.Window) *fyne.Container {
	btn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		//btn := widget.NewButtonWithIcon(i18n.GetString("WalletWindow_CopyAddress"), theme.ContentCopyIcon(), func() {
		w.Clipboard().SetContent(address)
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("WalletWindow_AddressCopied"), w)
	})

	exportBtn := widget.NewButtonWithIcon(i18n.GetString("Wallet_Export"), theme.FileIcon(),
		func() {

		})
	exportBtn.Importance = widget.HighImportance

	addr := newMyEntry()
	addr.Text = address
	addr.ActionItem = btn

	bala := newMyEntryWithData(AccountBalance)
	AccountBalance.Set(balance)
	if balance == "" {
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("Rpc_Get_Amount_fail"), WalletWindow)
	}
	exportBtnContainer := container.New(layout.NewPaddedLayout(), exportBtn)
	var png []byte
	png, _ = qrcode.Encode("xdag:"+address, qrcode.Medium, 256)

	image := canvas.NewImageFromResource(&fyne.StaticResource{
		StaticName:    "AddressQRcode",
		StaticContent: png,
	})
	image.SetMinSize(fyne.NewSize(256, 256))

	c := container.NewVBox(
		widget.NewLabel(""),
		container.New(layout.NewMaxLayout(), &widget.Form{
			Items: []*widget.FormItem{
				{Text: i18n.GetString("WalletWindow_AddressTitle"),
					Widget: addr},
				{Text: i18n.GetString("WalletWindow_BalanceTitle"),
					Widget: bala},
			},
		}),
		exportBtnContainer,
		widget.NewLabel(""),
		container.NewHBox(layout.NewSpacer(), image, layout.NewSpacer()))
	if LogonWindow.WalletType == HAS_ONLY_XDAG {
		c.Remove(exportBtnContainer)
	}
	return c
}

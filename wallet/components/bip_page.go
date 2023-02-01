package components

import (
	"goXdagWallet/i18n"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	qrcode "github.com/skip2/go-qrcode"
)

var BipAccountBalance = binding.NewString()

func BipPage(address, balance string, w fyne.Window) *fyne.Container {
	btn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		//btn := widget.NewButtonWithIcon(i18n.GetString("WalletWindow_CopyAddress"), theme.ContentCopyIcon(), func() {
		w.Clipboard().SetContent(address)
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("WalletWindow_AddressCopied"), w)
	})

	addr := newMyEntry()
	addr.Text = address
	addr.ActionItem = btn

	bala := newMyEntryWithData(BipAccountBalance)
	BipAccountBalance.Set(balance)

	var png []byte
	png, _ = qrcode.Encode("xdag:"+address, qrcode.Medium, 256)

	image := canvas.NewImageFromResource(&fyne.StaticResource{
		StaticName:    "AddressQRcode",
		StaticContent: png,
	})
	image.SetMinSize(fyne.NewSize(256, 256))

	return container.NewVBox(
		widget.NewLabel(""),
		container.New(layout.NewMaxLayout(), &widget.Form{
			Items: []*widget.FormItem{
				{Text: i18n.GetString("WalletWindow_AddressTitle"),
					Widget: addr},
				{Text: i18n.GetString("WalletWindow_BalanceTitle"),
					Widget: bala},
			},
		}),
		widget.NewLabel(""),
		container.NewHBox(layout.NewSpacer(), image, layout.NewSpacer()))
}

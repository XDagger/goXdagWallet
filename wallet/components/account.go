package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"goXdagWallet/i18n"
)

var AccountBalance = binding.NewString()

func AccountPage(address, balance string, w fyne.Window) *fyne.Container {
	addr := widget.NewEntry()
	addr.Text = address
	addr.Disable()

	bala := widget.NewEntryWithData(AccountBalance)
	AccountBalance.Set(balance)
	bala.Disable()

	return container.NewGridWithRows(3,
		layout.NewSpacer(),
		container.New(layout.NewMaxLayout(), &widget.Form{
			Items: []*widget.FormItem{ // we can specify items in the constructor
				{Text: i18n.GetString("WalletWindow_AddressTitle"), Widget: addr},
				{Text: i18n.GetString("WalletWindow_BalanceTitle"), Widget: bala},
			},
			SubmitText: i18n.GetString("WalletWindow_CopyAddress"),
			OnSubmit: func() {
				w.Clipboard().SetContent(address)
				dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
					i18n.GetString("WalletWindow_AddressCopied"), w)
			},
		}),

		layout.NewSpacer())
}

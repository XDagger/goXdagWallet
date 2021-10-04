package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"goXdagWallet/i18n"
)

var AccountBalance = binding.NewString()

func AccountPage(address, balance string, w fyne.Window) *fyne.Container {
	addr := widget.NewEntry()
	addr.Text = address

	bala := widget.NewEntryWithData(AccountBalance)
	AccountBalance.Set(balance)

	return container.NewGridWithRows(3,
		layout.NewSpacer(),
		container.New(layout.NewMaxLayout(), &widget.Form{
			Items: []*widget.FormItem{ // we can specify items in the constructor
				{Text: i18n.GetString("WalletWindow_AddressTitle"), Widget: addr},
				{Text: i18n.GetString("WalletWindow_BalanceTitle"), Widget: bala},
			},
		}),
		layout.NewSpacer())
}

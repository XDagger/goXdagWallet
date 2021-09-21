package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"goXdagWallet/i18n"
)

//var transfer func(address, amount, remark string) int

func TransferPage(w fyne.Window, transWrap func(string, string, string) int) *fyne.Container {

	//transfer = transWrap

	addr := widget.NewEntry()
	amount := widget.NewEntry()
	remark := widget.NewEntry()

	return container.NewGridWithRows(3,
		layout.NewSpacer(),
		container.New(layout.NewMaxLayout(), &widget.Form{
			Items: []*widget.FormItem{ // we can specify items in the constructor
				{Text: i18n.GetString("WalletWindow_Transfer_ToAddress"), Widget: addr},
				{Text: i18n.GetString("WalletWindow_Transfer_Amount"), Widget: amount},
				{Text: i18n.GetString("WalletWindow_Transfer_Remark"), Widget: remark},
			},
			SubmitText: i18n.GetString("TransferWindow_TransferTitle"),
			OnSubmit: func() {
				transWrap(addr.Text, amount.Text, remark.Text)
			},
		}),
		layout.NewSpacer())
}

//func doTransfer() {
//	transfer("Rr7mGuOOlpnhmYwb3eJQ4G48U4fiKD1m", "1.5", "goxdagwallet")
//}

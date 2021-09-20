package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"goXdagWallet/i18n"
)

var transfer func(address, amount, remark string) int

func TransferPage(w fyne.Window, transWrap func(string, string, string) int) fyne.Widget {
	//return widget.NewLabel("Transfer")
	transfer = transWrap
	return widget.NewButton(i18n.GetString("TransferWindow_TransferTitle"), doTransfer)
}

func doTransfer() {
	transfer("Rr7mGuOOlpnhmYwb3eJQ4G48U4fiKD1m", "1.5", "goxdagwallet")
}

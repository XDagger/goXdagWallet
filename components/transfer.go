package components

//#include "../src/xdag_runtime.h"
//#include <stdlib.h>
import "C"
import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"goXdagWallet/i18n"
	"strconv"
)

var TransStatus = widget.NewLabel("")
var TransBtnContainer *fyne.Container
var TransProgressContainer *fyne.Container

func TransferPage(w fyne.Window, transWrap func(string, string, string) int) *fyne.Container {

	addr := widget.NewEntry()
	amount := widget.NewEntry()
	remark := widget.NewEntry()

	TransProgressContainer = container.New(layout.NewPaddedLayout(), widget.NewProgressBarInfinite())
	TransProgressContainer.Hide()

	btn := widget.NewButtonWithIcon(i18n.GetString("TransferWindow_TransferTitle"), theme.ConfirmIcon(),
		func() {
			if !checkInput(addr.Text, amount.Text, remark.Text, w) {
				return
			}
			message := fmt.Sprintf(i18n.GetString("TransferWindow_ConfirmTransfer"), amount, addr)
			dialog.ShowConfirm(i18n.GetString("Common_ConfirmTitle"),
				message, func(b bool) {
					if b {
						TransProgressContainer.Show()
						TransBtnContainer.Hide()
						TransStatus.Text = i18n.GetString("TransferWindow_CommittingTransaction")
						transWrap(addr.Text, amount.Text, remark.Text)
					}
				}, w)

		})
	btn.Importance = widget.HighImportance
	TransBtnContainer = container.New(layout.NewPaddedLayout(), btn)

	return container.NewGridWithRows(4,
		layout.NewSpacer(),
		container.New(layout.NewMaxLayout(), &widget.Form{
			Items: []*widget.FormItem{ // we can specify items in the constructor
				{Text: i18n.GetString("WalletWindow_Transfer_ToAddress"), Widget: addr},
				{Text: i18n.GetString("WalletWindow_Transfer_Amount"), Widget: amount},
				{Text: i18n.GetString("WalletWindow_Transfer_Remark"), Widget: remark},
			},
		}),
		container.NewVBox(container.NewHBox(
			layout.NewSpacer(), TransStatus, layout.NewSpacer()),
			TransProgressContainer,
			TransBtnContainer),
		layout.NewSpacer())
}

//func doTransfer() {
//	transfer("Rr7mGuOOlpnhmYwb3eJQ4G48U4fiKD1m", "1.5", "goxdagwallet")
//}

func checkInput(addr, amount, remark string, window fyne.Window) bool {
	if len(addr) == 0 || !ValidateAddress(addr) {
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("TransferWindow_AccountFormatError"), window)
		return false
	}

	value, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("TransferWindow_AmountFormatError"), window)
		return false
	}

	balance, _ := strconv.ParseFloat(Balance, 64)
	if balance < value {
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("TransferWindow_InsufficientAmount"), window)
		return false
	}

	if len(remark) > 0 && !ValidateRemark(remark) {
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("TransferWindow_RemarkFormatError"), window)
		return false
	}
	return true
}

func setTransferDone() {
	TransProgressContainer.Hide()
	TransBtnContainer.Show()
	dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
		i18n.GetString("TransferWindow_CommitSuccess"), WalletWindow)
}

func setTransferError(e string) {
	TransProgressContainer.Hide()
	TransBtnContainer.Show()
	dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
		i18n.GetString("TransferWindow_CommitFailed")+e, WalletWindow)
}

package components

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"goXdagWallet/config"
	"goXdagWallet/i18n"
	"strconv"
)

var TransStatus = widget.NewLabel("")
var TransBtnContainer *fyne.Container
var TransProgressContainer *fyne.Container
var AddressList *widget.List
var AddressEntry = widget.NewEntry()

func TransferPage(w fyne.Window, transWrap func(string, string, string) int) *fyne.Container {
	amount := widget.NewEntry()
	remark := widget.NewEntry()

	TransProgressContainer = container.New(layout.NewPaddedLayout(), widget.NewProgressBarInfinite())
	TransProgressContainer.Hide()

	btn := widget.NewButtonWithIcon(i18n.GetString("TransferWindow_TransferTitle"), theme.ConfirmIcon(),
		func() {
			if !checkInput(AddressEntry.Text, amount.Text, remark.Text, w) {
				return
			}
			message := fmt.Sprintf(i18n.GetString("TransferWindow_ConfirmTransfer"), amount.Text, AddressEntry.Text)
			fmt.Println(message)
			dialog.ShowConfirm(i18n.GetString("Common_ConfirmTitle"),
				message, func(b bool) {
					if b {
						TransProgressContainer.Show()
						TransBtnContainer.Hide()
						TransStatus.Text = i18n.GetString("TransferWindow_CommittingTransaction")
						transWrap(AddressEntry.Text, amount.Text, remark.Text)
					}
				}, w)

		})
	btn.Importance = widget.HighImportance
	TransBtnContainer = container.New(layout.NewPaddedLayout(), btn)
	makeAddrList()
	top := container.NewVBox(
		container.NewHBox(
			layout.NewSpacer(), TransStatus, layout.NewSpacer()),
		container.New(layout.NewMaxLayout(), &widget.Form{
			Items: []*widget.FormItem{ // we can specify items in the constructor
				{Text: i18n.GetString("WalletWindow_Transfer_ToAddress"), Widget: AddressEntry},
				{Text: i18n.GetString("WalletWindow_Transfer_Amount"), Widget: amount},
				{Text: i18n.GetString("WalletWindow_Transfer_Remark"), Widget: remark},
			},
		}),
		container.NewVBox(
			TransProgressContainer,
			TransBtnContainer),
		widget.NewLabel(""),
		widget.NewLabel(i18n.GetString("TransferWindow_MostRecently")),
	)

	return container.New(
		layout.NewBorderLayout(top, nil, nil, nil),
		top,
		AddressList,
	)
}

func checkInput(addr, amount, remark string, window fyne.Window) bool {
	if len(addr) == 0 || !ValidateAddress(addr) {
		dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
			i18n.GetString("TransferWindow_AccountFormatError"), window)
		return false
	}

	value, err := strconv.ParseFloat(amount, 64)
	if err != nil || value <= 0.0 {
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

func makeAddrList() {
	AddressList = widget.NewList(
		func() int {
			return len(config.GetConfig().Addresses)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.ContentCopyIcon()), widget.NewLabel("address"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(config.GetConfig().Addresses[id])
		})
	AddressList.OnSelected = func(id widget.ListItemID) {
		AddressEntry.SetText(config.GetConfig().Addresses[id])
	}
}

func setTransferDone() {
	TransProgressContainer.Hide()
	TransBtnContainer.Show()
	dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
		i18n.GetString("TransferWindow_CommitSuccess"), WalletWindow)
	makeAddrList()
}

func setTransferError(e string) {
	TransProgressContainer.Hide()
	TransBtnContainer.Show()
	dialog.ShowInformation(i18n.GetString("Common_MessageTitle"),
		i18n.GetString("TransferWindow_CommitFailed")+e, WalletWindow)
}

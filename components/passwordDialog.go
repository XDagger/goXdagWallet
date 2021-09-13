package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// PasswordDialog is a variation of a dialog which prompts the user to enter some text.
type PasswordDialog struct {
	entry *widget.Entry

	onClosed func()
}

// SetOnClosed changes the callback which is run when the dialog is closed,
// which is nil by default.
//
// The callback is called unconditionally whether the user confirms or cancels.
//
// Note that the callback will be called after onConfirm, if both are non-nil.
// This way onConfirm can potential modify state that this callback needs to
// get the user input when the user confirms, while also being able to handle
// the case where the user cancelled.
func (i *PasswordDialog) SetOnClosed(callback func()) {
	i.onClosed = callback
}

// ShowPasswordDialog shows the password dialog.
func (i *PasswordDialog) ShowPasswordDialog(title, message, confirm, dismiss string, onConfirm func(string), parent fyne.Window) {

	i.entry.Password = true
	items := []*widget.FormItem{widget.NewFormItem(message, i.entry)}
	dialog.ShowForm(title, confirm, dismiss, items, func(_ bool) {
		// User has confirmed and entered an input
		if onConfirm != nil {
			onConfirm(i.entry.Text)
		}
		if i.onClosed != nil {
			i.onClosed()
		}

	}, parent)

}

// NewPasswordDialog creates a dialog over the specified window for the user to enter a password
func NewPasswordDialog(callback func()) *PasswordDialog {
	i := &PasswordDialog{entry: widget.NewEntry()}
	if callback != nil {
		i.SetOnClosed(callback)
	}
	return i
}

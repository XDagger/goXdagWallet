package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func AboutPage(w fyne.Window) *fyne.Container {

	return container.New(layout.NewMaxLayout(), layout.NewSpacer())
}

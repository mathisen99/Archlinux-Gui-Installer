package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Arch Linux Installer")
	w.Resize(fyne.NewSize(1024, 768))
	w.CenterOnScreen()

	wizard := NewWizard(w)
	w.SetContent(wizard.Layout())

	w.ShowAndRun()
}

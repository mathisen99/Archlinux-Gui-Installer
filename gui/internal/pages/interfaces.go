package pages

import (
	"archgui/gui/internal/state"

	"fyne.io/fyne/v2"
)

type WizardController interface {
	Next()
	Back()
	ShowLog(msg string)
	SetNextButtonEnabled(bool)
}

type Page interface {
	// Content returns the UI for this page
	Content(config *state.InstallConfig, ctrl WizardController) fyne.CanvasObject

	// OnNext is called when Next is clicked. Return error to block navigation.
	OnNext(config *state.InstallConfig) error

	// Title returns the title of the page
	Title() string
}

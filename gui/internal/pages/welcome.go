package pages

import (
	"archgui/gui/internal/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type WelcomePage struct{}

func (p *WelcomePage) Title() string {
	return "Welcome"
}

func (p *WelcomePage) Content(config *state.InstallConfig, ctrl WizardController) fyne.CanvasObject {
	return container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle("Welcome to Arch Linux Installer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabel("This wizard will guide you through the installation process."),
			widget.NewLabel("You can choose between Automatic and Manual partitioning."),
			widget.NewLabel(""),
			widget.NewLabel("Click 'Next' to begin."),
		),
	)
}

func (p *WelcomePage) OnNext(config *state.InstallConfig) error {
	return nil
}

func NewWelcomePage() *WelcomePage {
	return &WelcomePage{}
}

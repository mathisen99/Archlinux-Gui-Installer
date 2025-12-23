package pages

import (
	"archgui/gui/internal/state"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type AccountPage struct{}

func (p *AccountPage) Title() string {
	return "User Account"
}

func (p *AccountPage) Content(config *state.InstallConfig, ctrl WizardController) fyne.CanvasObject {
	hostEntry := widget.NewEntry()
	hostEntry.Text = config.Hostname
	hostEntry.OnChanged = func(s string) { config.Hostname = s }

	fullEntry := widget.NewEntry()
	fullEntry.SetPlaceHolder("Firstname Lastname")
	fullEntry.OnChanged = func(s string) { config.FullName = s }

	userEntry := widget.NewEntry()
	userEntry.Text = config.Username
	userEntry.OnChanged = func(s string) { config.Username = s }

	rootPass := widget.NewPasswordEntry()
	rootPass.OnChanged = func(s string) { config.RootPassword = s }

	userPass := widget.NewPasswordEntry()
	userPass.OnChanged = func(s string) { config.UserPassword = s }

	// Shell Select
	shellSelect := widget.NewSelect([]string{"bash", "zsh", "zsh-ohmyzsh"}, func(s string) {
		config.Shell = s
	})
	shellSelect.SetSelected(config.Shell)

	return container.NewVBox(
		widget.NewLabelWithStyle("System & User Account", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewForm(
			widget.NewFormItem("Hostname", hostEntry),
			widget.NewFormItem("Full Name", fullEntry),
			widget.NewFormItem("Username", userEntry),
			widget.NewFormItem("Root Password", rootPass),
			widget.NewFormItem("User Password", userPass),
			widget.NewFormItem("Preferred Shell", shellSelect),
		),
	)
}

func (p *AccountPage) OnNext(config *state.InstallConfig) error {
	if config.Username == "" || config.Hostname == "" {
		return fmt.Errorf("username and hostname are required")
	}
	if config.RootPassword == "" || config.UserPassword == "" {
		return fmt.Errorf("passwords are required")
	}
	return nil
}

func NewAccountPage() *AccountPage {
	return &AccountPage{}
}

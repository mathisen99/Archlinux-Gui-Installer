package pages

import (
	"archgui/gui/internal/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type LocalizationPage struct{}

func (p *LocalizationPage) Title() string {
	return "Localization"
}

func (p *LocalizationPage) Content(config *state.InstallConfig, ctrl WizardController) fyne.CanvasObject {
	tzEntry := widget.NewEntry()
	tzEntry.SetPlaceHolder("Europe/London")
	tzEntry.Text = config.Timezone
	tzEntry.OnChanged = func(s string) { config.Timezone = s }

	locEntry := widget.NewEntry()
	locEntry.SetPlaceHolder("en_US")
	locEntry.Text = config.Locale
	locEntry.OnChanged = func(s string) { config.Locale = s }

	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("us")
	keyEntry.Text = config.Keymap
	keyEntry.OnChanged = func(s string) { config.Keymap = s }

	return container.NewVBox(
		widget.NewLabelWithStyle("Configure Localization", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewForm(
			widget.NewFormItem("Timezone", tzEntry),
			widget.NewFormItem("Locale", locEntry),
			widget.NewFormItem("Keymap", keyEntry),
		),
	)
}

func (p *LocalizationPage) OnNext(config *state.InstallConfig) error {
	return nil
}

func NewLocalizationPage() *LocalizationPage {
	return &LocalizationPage{}
}

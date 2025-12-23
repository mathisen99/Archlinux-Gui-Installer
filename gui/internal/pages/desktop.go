package pages

import (
	"archgui/gui/internal/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type DesktopPage struct{}

func (p *DesktopPage) Title() string {
	return "Desktop Environment"
}

func (p *DesktopPage) Content(config *state.InstallConfig, ctrl WizardController) fyne.CanvasObject {
	desktopSelect := widget.NewSelect([]string{
		"none", "xfce", "gnome", "kde", "i3", "sway", "hyprland", "cinnamon", "mate", "lxqt", "budgie",
	}, func(s string) {
		config.Desktop = s
	})
	desktopSelect.SetSelected(config.Desktop)

	nvidiaCheck := widget.NewCheck("Install NVIDIA Drivers", func(b bool) {
		config.InstallNvidia = b
	})
	nvidiaCheck.Checked = config.InstallNvidia

	return container.NewVBox(
		widget.NewLabelWithStyle("Choose Desktop Environment", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewForm(
			widget.NewFormItem("Desktop", desktopSelect),
			widget.NewFormItem("Graphics", nvidiaCheck),
		),
	)
}

func (p *DesktopPage) OnNext(config *state.InstallConfig) error {
	return nil
}

func NewDesktopPage() *DesktopPage {
	return &DesktopPage{}
}

package pages

import (
	"archgui/gui/internal/state"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SummaryPage struct{}

func (p *SummaryPage) Title() string {
	return "Installation Summary"
}

func (p *SummaryPage) Content(config *state.InstallConfig, ctrl WizardController) fyne.CanvasObject {
	// Build summary text
	summary := fmt.Sprintf(`Target Disk: %s
Manual Partitioning: %v
Filesystem: %s
Encrypt: %v

Hostname: %s
User: %s (Full: %s)
Shell: %s

Timezone: %s
Locale: %s
Keymap: %s

Desktop: %s
Nvidia: %v
`,
		config.Disk, config.ManualPartitioning, config.Filesystem, config.Encrypt,
		config.Hostname, config.Username, config.FullName, config.Shell,
		config.Timezone, config.Locale, config.Keymap,
		config.Desktop, config.InstallNvidia)

	if config.ManualPartitioning {
		summary += fmt.Sprintf("\nManual Targets:\nRoot: %s (Format: %v)\nEFI: %s (Format: %v)",
			config.TargetRoot, config.FormatRoot, config.TargetEFI, config.FormatEFI)
	}

	return container.NewVBox(
		widget.NewLabelWithStyle("Ready to Install", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Please review your settings below."),
		widget.NewSeparator(),
		widget.NewLabel(summary),
		widget.NewSeparator(),
		widget.NewLabel("Click 'Install' to begin. This operation cannot be undone."),
	)
}

func (p *SummaryPage) OnNext(config *state.InstallConfig) error {
	return nil
}

func NewSummaryPage() *SummaryPage {
	return &SummaryPage{}
}

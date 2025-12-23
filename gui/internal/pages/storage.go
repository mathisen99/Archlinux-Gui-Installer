package pages

import (
	"fmt"
	"os/exec"
	"strings"

	"archgui/gui/internal/data"
	"archgui/gui/internal/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type StoragePage struct {
	// Auto Widgets
	diskSelect *widget.Select
	fsSelect   *widget.Select
	encCheck   *widget.Check
	luksPass   *widget.Entry

	// Manual Widgets
	openCfdiskBtn *widget.Button
	refreshBtn    *widget.Button
	rootSelect    *widget.Select
	formatRoot    *widget.Check
	efiSelect     *widget.Select
	formatEfi     *widget.Check

	// Logic
	modeSelect *widget.RadioGroup

	contentContainer *fyne.Container
}

func (p *StoragePage) Title() string {
	return "Storage Configuration"
}

func (p *StoragePage) Content(config *state.InstallConfig, ctrl WizardController) fyne.CanvasObject {
	// --- Auto Partitioning Widgets ---
	p.diskSelect = widget.NewSelect(data.GetDisks(), func(val string) {
		// Parse: /dev/sda (...)
		if len(val) > 0 {
			var name string
			_, _ = fmt.Sscanf(val, "/dev/%s", &name)
			config.Disk = "/dev/" + strings.TrimSpace(strings.Split(name, " ")[0]) // simple parse
		}
	})
	if len(p.diskSelect.Options) > 0 {
		p.diskSelect.SetSelected(p.diskSelect.Options[0])
	}

	p.fsSelect = widget.NewSelect([]string{"ext4", "btrfs"}, func(val string) {
		config.Filesystem = val
	})
	p.fsSelect.SetSelected(config.Filesystem)

	p.luksPass = widget.NewPasswordEntry()
	p.luksPass.Disable()
	p.luksPass.OnChanged = func(s string) { config.LuksPassword = s }

	p.encCheck = widget.NewCheck("Encrypt Disk (LUKS)", func(checked bool) {
		config.Encrypt = checked
		if checked {
			p.luksPass.Enable()
		} else {
			p.luksPass.Disable()
		}
	})
	p.encCheck.Checked = config.Encrypt

	autoContent := container.NewVBox(
		widget.NewLabelWithStyle("Automatic Partitioning", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Warning: Selected disk will be WIPED."),
		widget.NewForm(
			widget.NewFormItem("Target Disk", p.diskSelect),
			widget.NewFormItem("Filesystem", p.fsSelect),
			widget.NewFormItem("Encryption", p.encCheck),
			widget.NewFormItem("LUKS Password", p.luksPass),
		),
	)

	// --- Manual Partitioning Widgets ---
	p.openCfdiskBtn = widget.NewButton("Open Partition Manager (cfdisk)", func() {
		// Launch cfdisk in xterm
		cmd := exec.Command("xterm", "-e", "cfdisk")
		if err := cmd.Start(); err != nil {
			// Show error (ctrl needed? or log)
			ctrl.ShowLog(fmt.Sprintf("Failed to launch cfdisk: %v", err))
		}
	})

	p.rootSelect = widget.NewSelect(data.GetPartitions(), func(val string) {
		config.TargetRoot = parseDevPath(val)
	})
	p.formatRoot = widget.NewCheck("Format Root?", func(b bool) { config.FormatRoot = b })
	p.formatRoot.Checked = true

	p.efiSelect = widget.NewSelect(data.GetPartitions(), func(val string) {
		config.TargetEFI = parseDevPath(val)
	})
	p.formatEfi = widget.NewCheck("Format EFI?", func(b bool) { config.FormatEFI = b })

	p.refreshBtn = widget.NewButton("Refresh Partitions", func() {
		parts := data.GetPartitions()
		p.rootSelect.Options = parts
		p.efiSelect.Options = parts
		p.rootSelect.Refresh()
		p.efiSelect.Refresh()
	})

	manualContent := container.NewVBox(
		widget.NewLabelWithStyle("Manual Partitioning", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("1. Use cfdisk to create partitions."),
		p.openCfdiskBtn,
		widget.NewLabel("2. Select partitions to use."),
		p.refreshBtn,
		widget.NewForm(
			widget.NewFormItem("Root Partition (/)", p.rootSelect),
			widget.NewFormItem("", p.formatRoot),
			widget.NewFormItem("EFI Partition (/boot)", p.efiSelect),
			widget.NewFormItem("", p.formatEfi),
		),
	)

	// --- Mode Switching ---
	p.contentContainer = container.NewStack(autoContent, manualContent)
	// defaulting to Auto
	manualContent.Hide()

	p.modeSelect = widget.NewRadioGroup([]string{"Automatic", "Manual"}, func(val string) {
		if val == "Automatic" {
			config.ManualPartitioning = false
			manualContent.Hide()
			autoContent.Show()
		} else {
			config.ManualPartitioning = true
			autoContent.Hide()
			manualContent.Show()
		}
	})
	p.modeSelect.Horizontal = true
	p.modeSelect.Selected = "Automatic"

	return container.NewVBox(
		widget.NewLabel("Choose Partitioning Mode:"),
		p.modeSelect,
		widget.NewSeparator(),
		p.contentContainer,
	)
}

func parseDevPath(val string) string {
	// "/dev/sda1 (10G)" -> "/dev/sda1"
	if len(val) > 0 {
		var name string
		_, _ = fmt.Sscanf(val, "/dev/%s", &name)
		return "/dev/" + strings.TrimSpace(strings.Split(name, " ")[0])
	}
	return ""
}

func (p *StoragePage) OnNext(config *state.InstallConfig) error {
	// Validation
	if config.ManualPartitioning {
		if config.TargetRoot == "" {
			return fmt.Errorf("please select a Root partition")
		}
	} else {
		if config.Disk == "" {
			return fmt.Errorf("please select a Target Disk")
		}
		if config.Encrypt && config.LuksPassword == "" {
			return fmt.Errorf("LUKS Password is required")
		}
	}
	return nil
}

func NewStoragePage() *StoragePage {
	return &StoragePage{}
}

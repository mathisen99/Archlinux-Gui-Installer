package pages

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"

	"archgui/gui/internal/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type InstallPage struct {
	logOutput *widget.Entry
	started   bool
}

func (p *InstallPage) Title() string {
	return "Installing"
}

func (p *InstallPage) Content(config *state.InstallConfig, ctrl WizardController) fyne.CanvasObject {
	p.logOutput = widget.NewMultiLineEntry()
	p.logOutput.Disable()
	// Make it fill space?
	// In VContainer, we can wrap in Scroll

	// We delay start slightly to ensure UI renders
	if !p.started {
		p.started = true
		go p.RunInstall(config, ctrl)
	}

	return container.NewBorder(
		widget.NewLabelWithStyle("Installing Arch Linux...", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil, nil, nil,
		container.NewScroll(p.logOutput),
	)
}

func (p *InstallPage) RunInstall(config *state.InstallConfig, ctrl WizardController) {
	// 1. Generate Config
	// We need to disable Next/Back during install

	fyne.Do(func() {
		ctrl.SetNextButtonEnabled(false)
	})
	// Actually we can't go back either.

	time.Sleep(500 * time.Millisecond) // UI settle

	cfgStr := generateConfigEnv(config)
	configFile := "/tmp/install.env"
	if err := os.WriteFile(configFile, []byte(cfgStr), 0600); err != nil {
		p.AppendLog("Error writing config: " + err.Error())
		return
	}

	cmd := exec.Command("bash", "backend/arch-install.sh", "--config", configFile)

	// Determine Pipe
	cmdPipe, err := cmd.StdoutPipe()
	if err != nil {
		p.AppendLog("Error creating pipe: " + err.Error())
		return
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		p.AppendLog("Error starting command: " + err.Error())
		return
	}

	scanner := bufio.NewScanner(cmdPipe)
	for scanner.Scan() {
		p.AppendLog(scanner.Text())
	}

	err = cmd.Wait()
	if err != nil {
		p.AppendLog(fmt.Sprintf("\nInstallation FAILED: %v", err))
	} else {
		p.AppendLog("\nInstallation SUCCESS! You can reboot now.")
		// Maybe enable a "Finish" button?
		// ctrl.Next() to a "Done" page? Or just leave it here.
	}
	// We keep buttons disabled or maybe enable "Close"?
	// For wizard pattern, usually "Next" becomes "Finish".
	// But we are at the last page.
}

func (p *InstallPage) AppendLog(msg string) {
	str := msg + "\n"
	fyne.Do(func() {
		curr := p.logOutput.Text
		p.logOutput.SetText(curr + str)
		p.logOutput.Refresh()
	})
}

func (p *InstallPage) OnNext(config *state.InstallConfig) error {
	return nil
}

func NewInstallPage() *InstallPage {
	return &InstallPage{}
}

func boolToString(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func generateConfigEnv(c *state.InstallConfig) string {
	return fmt.Sprintf(
		"DISK=%s\nMANUAL_PARTITIONING=%s\nTARGET_ROOT=%s\nTARGET_EFI=%s\nFORMAT_ROOT=%s\nFORMAT_EFI=%s\n"+
			"HOSTNAME=%s\nFULL_NAME=%s\nUSERNAME=%s\nROOT_PASSWORD=%s\nUSER_PASSWORD=%s\n"+
			"TIMEZONE=%s\nLOCALE=%s\nKEYMAP=%s\n"+
			"FS_TYPE=%s\nUSE_LUKS=%s\nLUKS_PASSWORD=%s\n"+
			"DESKTOP_ENV=%s\nSHELL_CHOICE=%s\nHAS_NVIDIA=%s\nNONINTERACTIVE=yes\n",
		c.Disk, boolToString(c.ManualPartitioning), c.TargetRoot, c.TargetEFI, boolToString(c.FormatRoot), boolToString(c.FormatEFI),
		c.Hostname, c.FullName, c.Username, c.RootPassword, c.UserPassword,
		c.Timezone, c.Locale, c.Keymap,
		c.Filesystem,
		boolToString(c.Encrypt), c.LuksPassword,
		c.Desktop, c.Shell, boolToString(c.InstallNvidia),
	)
}

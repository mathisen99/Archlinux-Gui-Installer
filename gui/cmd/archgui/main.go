package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Disk structure for lsblk JSON parsing
type BlockDevice struct {
	Name string `json:"name"`
	Size string `json:"size"`
	Type string `json:"type"`
}

type LsblkOutput struct {
	BlockDevices []BlockDevice `json:"blockdevices"`
}

func getDisks() []string {
	cmd := exec.Command("lsblk", "-d", "-n", "-o", "NAME,SIZE,TYPE", "--json")
	output, err := cmd.Output()
	if err != nil {
		// Fallback for dev/testing if lsblk fails or not present
		return []string{"/dev/sda (Test)", "/dev/nmve0n1 (Test)"}
	}

	var data LsblkOutput
	if err := json.Unmarshal(output, &data); err != nil {
		return []string{"Error parsing disks"}
	}

	var disks []string
	for _, dev := range data.BlockDevices {
		// Filter out loop and rom devices usually
		if dev.Type == "disk" {
			disks = append(disks, fmt.Sprintf("/dev/%s (%s)", dev.Name, dev.Size))
		}
	}
	if len(disks) == 0 {
		disks = append(disks, "No disks found")
	}
	return disks
}

func main() {
	a := app.New()
	w := a.NewWindow("Arch Linux Installer")
	w.Resize(fyne.NewSize(1024, 768))
	w.CenterOnScreen()

	// Title
	title := widget.NewLabelWithStyle("Arch Linux GUI Installer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// 1. Disk Selection
	diskSelect := widget.NewSelect(getDisks(), func(value string) {
		// Handle selection if needed
	})
	if len(diskSelect.Options) > 0 {
		diskSelect.SetSelected(diskSelect.Options[0])
	}

	// 2. User Info
	hostnameEntry := widget.NewEntry()
	hostnameEntry.SetPlaceHolder("archlinux")
	hostnameEntry.Text = "archlinux"

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("user")
	usernameEntry.Text = "user"

	rootPassEntry := widget.NewPasswordEntry()
	userPassEntry := widget.NewPasswordEntry()

	// 3. System Config
	fsSelect := widget.NewSelect([]string{"ext4", "btrfs"}, nil)
	fsSelect.SetSelected("ext4")

	// Encryption
	luksPassEntry := widget.NewPasswordEntry()
	luksPassEntry.Disable()

	encryptCheck := widget.NewCheck("Encrypt Disk (LUKS)", func(checked bool) {
		if checked {
			luksPassEntry.Enable()
		} else {
			luksPassEntry.Disable()
			luksPassEntry.SetText("")
		}
	})

	// Desktop
	desktopSelect := widget.NewSelect([]string{
		"none", "xfce", "gnome", "kde", "i3", "sway", "hyprland",
	}, nil)
	desktopSelect.SetSelected("xfce")

	// Shell
	shellSelect := widget.NewSelect([]string{"bash", "zsh", "zsh-ohmyzsh"}, nil)
	shellSelect.SetSelected("bash")

	// Nvidia
	nvidiaCheck := widget.NewCheck("Install NVIDIA Drivers", nil)

	// Layout Construction
	// Using a Form widget for cleanliness
	paramForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Target Disk", Widget: diskSelect},
			{Text: "Hostname", Widget: hostnameEntry},
			{Text: "Username", Widget: usernameEntry},
			{Text: "Root Password", Widget: rootPassEntry},
			{Text: "User Password", Widget: userPassEntry},
			{Text: "Filesystem", Widget: fsSelect},
			{Text: "Encryption", Widget: encryptCheck},
			{Text: "LUKS Password", Widget: luksPassEntry},
			{Text: "Desktop Env", Widget: desktopSelect},
			{Text: "Shell", Widget: shellSelect},
			{Text: "Graphics", Widget: nvidiaCheck},
		},
	}

	// Logs
	logOutput := widget.NewMultiLineEntry()
	logOutput.Disable()
	logOutput.SetMinRowsVisible(15)

	// Install Action
	installBtn := widget.NewButton("Install Arch Linux", func() {
		logOutput.SetText("Starting installation...\n")

		// Validation
		diskVal := diskSelect.Selected
		// strip size info from disk string: "/dev/sda (20G)" -> "/dev/sda"
		var diskPath string
		if len(diskVal) > 0 {
			var name string
			_, _ = fmt.Sscanf(diskVal, "/dev/%s", &name)
			diskPath = "/dev/" + name
		}

		if diskPath == "" {
			logOutput.SetText("Error: No disk selected.")
			return
		}

		// Config Gen
		configContent := generateConfig(
			diskPath, hostnameEntry.Text, usernameEntry.Text, rootPassEntry.Text, userPassEntry.Text,
			fsSelect.Selected, luksPassEntry.Text, desktopSelect.Selected, shellSelect.Selected,
			encryptCheck.Checked, nvidiaCheck.Checked,
		)

		configFile := "/tmp/install.env"
		err := os.WriteFile(configFile, []byte(configContent), 0600)
		if err != nil {
			logOutput.SetText(fmt.Sprintf("Error writing config: %v", err))
			return
		}

		// Execute
		cmd := exec.Command("bash", "backend/arch-install.sh", "--config", configFile)

		// For MVP, capture output. In a real app, we'd stream stdout to the log widget.
		output, err := cmd.CombinedOutput()
		if err != nil {
			logOutput.SetText(fmt.Sprintf("Installation failed:\n%s\nError: %v", string(output), err))
		} else {
			logOutput.SetText(fmt.Sprintf("Installation Success:\n%s", string(output)))
		}
	})
	installBtn.Importance = widget.HighImportance

	// Container
	content := container.NewBorder(
		title, // top
		nil,   // bottom
		nil,   // left
		nil,   // right
		container.NewVBox(
			widget.NewLabel("Start by configuring your installation below:"),
			paramForm,
			installBtn,
			widget.NewLabel("Installation Log:"),
			container.NewScroll(logOutput),
		),
	)
	// Wrap content in a scroll container because form is long
	scrollContent := container.NewScroll(content)

	w.SetContent(scrollContent)
	w.ShowAndRun()
}

func boolToString(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func generateConfig(disk, host, user, rootPass, userPass, fs, luksPass, desktop, shell string, encrypt, nvidia bool) string {
	return fmt.Sprintf(
		"DISK=%s\nHOSTNAME=%s\nUSERNAME=%s\nROOT_PASSWORD=%s\nUSER_PASSWORD=%s\n"+
			"FS_TYPE=%s\nUSE_LUKS=%s\nLUKS_PASSWORD=%s\n"+
			"DESKTOP_ENV=%s\nSHELL_CHOICE=%s\nHAS_NVIDIA=%s\nNONINTERACTIVE=yes\n",
		disk, host, user, rootPass, userPass,
		fs,
		boolToString(encrypt), luksPass,
		desktop, shell, boolToString(nvidia),
	)
}

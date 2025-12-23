package main

import (
	"fmt"
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func generateConfig(disk, hostname, username, rootPass, userPass string) string {
	return fmt.Sprintf(
		"DISK=%s\nHOSTNAME=%s\nUSERNAME=%s\nROOT_PASSWORD=%s\nUSER_PASSWORD=%s\nNONINTERACTIVE=yes\n",
		disk, hostname, username, rootPass, userPass,
	)
}

func main() {
	a := app.New()
	w := a.NewWindow("Arch Linux Installer")
	w.Resize(fyne.NewSize(800, 600))

	// UI Elements
	title := widget.NewLabelWithStyle("Arch Linux GUI Installer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Form Inputs
	diskEntry := widget.NewEntry()
	diskEntry.SetPlaceHolder("/dev/sda")

	hostnameEntry := widget.NewEntry()
	hostnameEntry.SetPlaceHolder("archlinux")

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("user")

	rootPassEntry := widget.NewPasswordEntry()
	userPassEntry := widget.NewPasswordEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Target Disk", Widget: diskEntry},
			{Text: "Hostname", Widget: hostnameEntry},
			{Text: "Username", Widget: usernameEntry},
			{Text: "Root Password", Widget: rootPassEntry},
			{Text: "User Password", Widget: userPassEntry},
		},
	}

	// Log Output
	logOutput := widget.NewMultiLineEntry()
	logOutput.Disable()
	logOutput.SetMinRowsVisible(10)

	// Install Action
	installBtn := widget.NewButton("Install Arch Linux", func() {
		logOutput.SetText("Starting installation...\n")

		// 1. Generate Config
		configContent := generateConfig(diskEntry.Text, hostnameEntry.Text, usernameEntry.Text, rootPassEntry.Text, userPassEntry.Text)

		configFile := "/tmp/install.env"
		err := os.WriteFile(configFile, []byte(configContent), 0600)
		if err != nil {
			logOutput.SetText(fmt.Sprintf("Error writing config: %v", err))
			return
		}

		// 2. Run Backend
		// Note: In production we need to handle absolute paths properly or assume PWD
		cmd := exec.Command("bash", "backend/arch-install.sh", "--config", configFile)

		// Capture output
		// For MVP, we'll just capture combined output at the end or try basic streaming
		// Real streaming requires a bit more goroutine work

		output, err := cmd.CombinedOutput()
		if err != nil {
			logOutput.SetText(fmt.Sprintf("Installation failed:\n%s\nError: %v", string(output), err))
		} else {
			logOutput.SetText(fmt.Sprintf("Installation Success:\n%s", string(output)))
		}
	})

	content := container.NewVBox(
		title,
		widget.NewLabel("Start by configuring your installation below:"),
		form,
		installBtn,
		widget.NewLabel("Installation Log:"),
		container.NewScroll(logOutput), // Wrap log in scroll
	)

	w.SetContent(content)
	w.ShowAndRun()
}

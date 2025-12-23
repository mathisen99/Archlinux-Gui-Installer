package main

import (
	"archgui/gui/internal/pages"
	"archgui/gui/internal/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Wizard struct {
	window  fyne.Window
	content *fyne.Container

	pages   []pages.Page
	current int

	nextBtn *widget.Button
	backBtn *widget.Button

	config *state.InstallConfig
}

func NewWizard(w fyne.Window) *Wizard {
	wiz := &Wizard{
		window: w,
		config: state.NewInstallConfig(),
	}

	// Initialize pages
	wiz.pages = []pages.Page{
		pages.NewWelcomePage(),
		pages.NewStoragePage(),
		pages.NewLocalizationPage(),
		pages.NewAccountPage(),
		pages.NewDesktopPage(),
		pages.NewSummaryPage(),
		pages.NewInstallPage(),
	}

	// Setup UI
	wiz.nextBtn = widget.NewButton("Next", wiz.Next)
	wiz.nextBtn.Importance = widget.HighImportance

	wiz.backBtn = widget.NewButton("Back", wiz.Back)

	wiz.content = container.NewStack()

	// Initial view
	wiz.updateView()

	return wiz
}

// UI Layout
func (w *Wizard) Layout() fyne.CanvasObject {
	return container.NewBorder(
		nil, // top
		container.NewHBox(widget.NewLabel("Arch Installer Wizard"), canvasSpacer(), w.backBtn, w.nextBtn), // bottom nav
		nil,       // left
		nil,       // right
		w.content, // center
	)
}

func canvasSpacer() fyne.CanvasObject {
	return &widget.Label{Text: " "} // Hacky spacer, layout.Spacer is better but needs importing layout
}

// Controller Implementation

func (w *Wizard) Next() {
	// Check validation of current page
	currentPage := w.pages[w.current]
	if err := currentPage.OnNext(w.config); err != nil {
		w.ShowLog(err.Error()) // Show error dialog
		return
	}

	if w.current < len(w.pages)-1 {
		w.current++
		w.updateView()
	}
}

func (w *Wizard) Back() {
	if w.current > 0 {
		w.current--
		w.updateView()
	}
}

func (w *Wizard) ShowLog(msg string) {
	// Simple error dialog
	d := widget.NewModalPopUp(widget.NewLabel(msg), w.window.Canvas())
	d.Show()
	// Or prefer standard dialog
	// dialog.ShowError(errors.New(msg), w.window)
}

func (w *Wizard) SetNextButtonEnabled(enabled bool) {
	if enabled {
		w.nextBtn.Enable()
	} else {
		w.nextBtn.Disable()
	}
}

func (w *Wizard) updateView() {
	p := w.pages[w.current]

	// Update Window Title
	w.window.SetTitle("Arch Installer - " + p.Title())

	// Update Content
	// We pass 'w' (Wizard) as the controller
	pageContent := p.Content(w.config, w)

	// Wrap in scroll just in case
	scroll := container.NewScroll(pageContent)

	w.content.Objects = []fyne.CanvasObject{scroll}
	w.content.Refresh()

	// Button State
	if w.current == 0 {
		w.backBtn.Disable()
	} else {
		w.backBtn.Enable()
	}

	if w.current == len(w.pages)-1 {
		// Last page (Install logs)
		// Usually Next is hidden or disabled, or "Finish"
		w.nextBtn.Disable()
		w.backBtn.Disable() // Can't go back from running install
	} else if w.current == len(w.pages)-2 {
		w.nextBtn.SetText("Install")
	} else {
		w.nextBtn.SetText("Next")
		w.nextBtn.Enable()
	}
}

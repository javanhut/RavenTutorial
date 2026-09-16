package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// tutorial holds the state of the walkthrough: which track the user picked
// and how far through it they are.
type tutorial struct {
	app    fyne.App
	win    fyne.Window
	device string
	steps  []Step
	index  int

	// What the arrow keys should do on the page that is currently drawn.
	// Every page sets these as it is built, and either may be nil - on the
	// welcome page there is nothing to go back to.
	onNext func()
	onBack func()
}

// GuiInit builds the window and starts the tutorial on the welcome page.
func GuiInit() {
	a := app.New()
	w := a.NewWindow("Raven Linux Tutorial")
	w.Resize(fyne.NewSize(860, 640))

	t := &tutorial{app: a, win: w}

	// The keyboard drives the whole tour, which matters on a desktop where
	// most of what is being taught is keys. Registered once, against fields
	// that each page rewrites, rather than re-registered per page.
	w.Canvas().SetOnTypedKey(t.typedKey)

	t.showWelcome()
	w.ShowAndRun()
}

// setContent puts a page in the window and then refreshes it.
//
// The refresh is not decoration. A wrapping label caches the size it worked
// out the first time it was asked, and the first time is during the layout
// pass that has not yet given it a width - so a paragraph reports the height
// of one line and the section under it is drawn on top of it. Refreshing once
// the page has been laid out clears that cache and sizes every paragraph
// against the width it actually got.
func (t *tutorial) setContent(content fyne.CanvasObject) {
	t.win.SetContent(content)
	refreshDeep(content)
}

// refreshDeep refreshes a tree of widgets from the leaves up.
//
// The order is the point. A container lays itself out from the sizes its
// children report, so refreshing the parent first would only re-read the
// stale sizes; the children have to be asked to work theirs out again before
// anything above them is laid out.
func refreshDeep(o fyne.CanvasObject) {
	switch v := o.(type) {
	case *fyne.Container:
		for _, child := range v.Objects {
			refreshDeep(child)
		}
	case *container.Scroll:
		refreshDeep(v.Content)
	}
	o.Refresh()
}

// typedKey turns a keypress into a page turn.
func (t *tutorial) typedKey(e *fyne.KeyEvent) {
	switch e.Name {
	case fyne.KeyRight, fyne.KeyReturn, fyne.KeyEnter, fyne.KeySpace, fyne.KeyPageDown:
		if t.onNext != nil {
			t.onNext()
		}
	case fyne.KeyLeft, fyne.KeyBackspace, fyne.KeyPageUp:
		if t.onBack != nil {
			t.onBack()
		}
	case fyne.KeyHome:
		t.showWelcome()
	}
}

// showWelcome is the first thing the user sees: a greeting and a Next button.
func (t *tutorial) showWelcome() {
	body := widget.NewLabel(
		"This tour introduces the desktop, its keys and gestures, and the " +
			"applications Raven Linux ships with.\n\n" +
			"It takes about ten minutes, and everything in it can be tried " +
			"as you read - the tutorial stays open beside whatever you " +
			"start.\n\n" +
			"Press Next, or the right arrow key, to get started.")
	body.Alignment = fyne.TextAlignCenter
	body.Wrapping = fyne.TextWrapWord

	next := widget.NewButtonWithIcon("Next", theme.NavigateNextIcon(), t.showDeviceChoice)
	next.Importance = widget.HighImportance

	t.onNext, t.onBack = t.showDeviceChoice, nil

	t.setContent(container.NewBorder(
		nil,
		container.NewCenter(next),
		nil, nil,
		container.NewCenter(container.NewVBox(
			heading("Welcome to Raven Linux"),
			body,
		)),
	))
}

// showDeviceChoice asks which kind of machine the user is on, which decides
// the tutorial track they get.
func (t *tutorial) showDeviceChoice() {
	question := widget.NewLabel("Are you using a laptop or a desktop?")
	question.Alignment = fyne.TextAlignCenter
	question.Wrapping = fyne.TextWrapWord

	hint := widget.NewLabel(
		"Only one chapter differs: a laptop is taught the trackpad " +
			"gestures, a desktop the mouse chords that do the same things.")
	hint.Alignment = fyne.TextAlignCenter
	hint.Wrapping = fyne.TextWrapWord

	laptop := widget.NewButtonWithIcon("Laptop", theme.ComputerIcon(), func() {
		t.startTrack("Laptop", LaptopSteps)
	})
	laptop.Importance = widget.HighImportance

	desktop := widget.NewButtonWithIcon("Desktop", theme.DesktopIcon(), func() {
		t.startTrack("Desktop", DesktopSteps)
	})
	desktop.Importance = widget.HighImportance

	choices := container.New(layout.NewGridLayout(2), laptop, desktop)

	back := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), t.showWelcome)

	// Next is deliberately not bound here: which track to start is the one
	// question the tutorial cannot answer on the user's behalf.
	t.onNext, t.onBack = nil, t.showWelcome

	t.setContent(container.NewBorder(
		nil,
		container.NewBorder(nil, nil, back, nil),
		nil, nil,
		container.NewCenter(container.NewVBox(
			heading("Pick Your Machine"),
			question,
			hint,
			widget.NewSeparator(),
			choices,
		)),
	))
}

// startTrack loads the chosen set of steps and shows the first one.
func (t *tutorial) startTrack(device string, steps []Step) {
	t.device = device
	t.steps = steps
	t.index = 0
	t.showStep()
}

// showStep draws the current step of the chosen track.
func (t *tutorial) showStep() {
	if len(t.steps) == 0 {
		t.showFinish()
		return
	}

	step := t.steps[t.index]
	last := t.index == len(t.steps)-1

	// Back goes to the previous step, or out to the device question on
	// step one.
	back := func() {
		if t.index == 0 {
			t.showDeviceChoice()
			return
		}
		t.index--
		t.showStep()
	}
	next := func() {
		if last {
			t.showFinish()
			return
		}
		t.index++
		t.showStep()
	}
	t.onNext, t.onBack = next, back

	backButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), back)

	// On the last step the Next button becomes Finish.
	nextLabel, nextIcon := "Next", theme.NavigateNextIcon()
	if last {
		nextLabel, nextIcon = "Finish", theme.ConfirmIcon()
	}
	nextButton := widget.NewButtonWithIcon(nextLabel, nextIcon, next)
	nextButton.Importance = widget.HighImportance

	contents := widget.NewButtonWithIcon("Contents", theme.ListIcon(), t.showContents)

	t.setContent(container.NewBorder(
		t.stepHeader(step),
		container.NewBorder(nil, nil, backButton, nextButton, container.NewCenter(contents)),
		nil, nil,
		container.NewVScroll(t.stepBody(step)),
	))
}

// stepHeader is the title and the progress readout above every step.
func (t *tutorial) stepHeader(step Step) fyne.CanvasObject {
	progress := widget.NewProgressBar()
	progress.Min, progress.Max = 0, float64(len(t.steps))
	progress.TextFormatter = func() string {
		return fmt.Sprintf("%s tutorial - step %d of %d",
			t.device, t.index+1, len(t.steps))
	}
	progress.SetValue(float64(t.index + 1))

	return container.NewVBox(
		heading(step.Title),
		progress,
		widget.NewSeparator(),
	)
}

// stepBody assembles everything below the title: the prose, the table of
// bindings, the applications, the footnote and whatever buttons the step
// asked for.
func (t *tutorial) stepBody(step Step) fyne.CanvasObject {
	body := widget.NewLabel(step.Body)
	body.Wrapping = fyne.TextWrapWord

	page := container.NewVBox(body)

	if len(step.Bindings) > 0 {
		page.Add(widget.NewSeparator())
		page.Add(subheading(orDefault(step.BindingsTitle, "Shortcuts")))
		page.Add(bindingTable(step.Bindings))
	}

	if len(step.Apps) > 0 {
		page.Add(widget.NewSeparator())
		for _, a := range step.Apps {
			page.Add(appRow(a, t.win))
		}
	}

	if step.Launch != nil {
		// Not named app: this file imports a package called that.
		if program := (App{Exec: step.Launch.Exec, Args: step.Launch.Args}); program.Installed() {
			button := widget.NewButtonWithIcon(step.Launch.Label, theme.MediaPlayIcon(), func() {
				if err := program.Start(); err != nil {
					dialog.ShowError(err, t.win)
				}
			})
			page.Add(container.NewCenter(button))
		}
	}

	if step.Panel == "oracle" {
		page.Add(widget.NewSeparator())
		page.Add(t.oraclePanel(oracleInstalled()))
	}

	if step.Note != "" {
		page.Add(widget.NewSeparator())
		page.Add(caption(step.Note))
	}

	return page
}

// bindingTable draws the keys on the left and what they do on the right.
//
// A form layout rather than a two-column grid: the keys column takes the
// width its widest chord needs and no more, which keeps the descriptions
// starting at the same place without giving half the page to "Esc".
func bindingTable(bindings []Binding) fyne.CanvasObject {
	rows := make([]fyne.CanvasObject, 0, len(bindings)*2)
	for _, b := range bindings {
		keys := widget.NewLabelWithStyle(b.Keys, fyne.TextAlignLeading,
			fyne.TextStyle{Monospace: true})
		does := widget.NewLabel(b.Does)
		does.Wrapping = fyne.TextWrapWord
		rows = append(rows, keys, does)
	}
	return container.New(layout.NewFormLayout(), rows...)
}

// appRow is one application in the tour: its name, what it is for, and a
// button that opens it.
//
// A program that is not installed still gets its row - knowing what exists is
// half the point of the page - but says so instead of offering a button that
// could only fail.
func appRow(a App, win fyne.Window) fyne.CanvasObject {
	name := widget.NewLabelWithStyle(a.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	what := widget.NewLabel(a.What)
	what.Wrapping = fyne.TextWrapWord

	var action fyne.CanvasObject
	if a.Installed() {
		label := "Open"
		if a.Terminal {
			label = "Open in Terminal"
		}
		button := widget.NewButtonWithIcon(label, theme.MediaPlayIcon(), func() {
			if err := a.Start(); err != nil {
				dialog.ShowError(err, win)
			}
		})
		action = button
	} else {
		action = caption("not installed")
	}

	return container.NewBorder(
		container.NewBorder(nil, nil, name, action),
		widget.NewSeparator(),
		nil, nil,
		what,
	)
}

// showContents opens the list of steps in this track, so a reader can go
// straight to the page they half-remember rather than pressing Next eleven
// times to reach it.
func (t *tutorial) showContents() {
	list := widget.NewList(
		func() int { return len(t.steps) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			label.SetText(fmt.Sprintf("%2d.  %s", id+1, t.steps[id].Title))
			label.TextStyle = fyne.TextStyle{Bold: id == t.index}
		},
	)

	jump := dialog.NewCustom("Contents", "Close", list, t.win)
	jump.Resize(fyne.NewSize(460, 480))
	list.OnSelected = func(id widget.ListItemID) {
		jump.Hide()
		t.index = id
		t.showStep()
	}
	jump.Show()
}

// showFinish is the closing page once a track has been walked through.
func (t *tutorial) showFinish() {
	body := widget.NewLabel(
		"That is the end of the " + t.device + " tour.\n\n" +
			"Nothing in it needs remembering: Super + Ctrl + H lists every " +
			"shortcut the desktop answers, at any time, from anywhere.\n\n" +
			"You can walk through the other track to see how the same " +
			"things are done on the other kind of machine, or start this " +
			"one again.")
	body.Alignment = fyne.TextAlignCenter
	body.Wrapping = fyne.TextWrapWord

	restart := widget.NewButtonWithIcon("Start Over", theme.ViewRefreshIcon(), t.showWelcome)
	restart.Importance = widget.HighImportance

	other := widget.NewButtonWithIcon("Take the Other Track", theme.ComputerIcon(), func() {
		if t.device == "Laptop" {
			t.startTrack("Desktop", DesktopSteps)
			return
		}
		t.startTrack("Laptop", LaptopSteps)
	})

	contents := widget.NewButtonWithIcon("Contents", theme.ListIcon(), t.showContents)

	t.onNext, t.onBack = nil, func() {
		t.index = len(t.steps) - 1
		t.showStep()
	}

	t.setContent(container.NewCenter(container.NewVBox(
		heading("All Done"),
		body,
		container.NewCenter(container.NewHBox(restart, other, contents)),
	)))
}

// heading makes a large, bold, centred title that follows the current theme.
func heading(text string) *widget.Label {
	l := widget.NewLabelWithStyle(text, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	l.SizeName = theme.SizeNameHeadingText
	return l
}

// subheading labels a section inside a step.
func subheading(text string) *widget.Label {
	l := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	l.SizeName = theme.SizeNameSubHeadingText
	return l
}

// caption is the smaller, quieter text used for footnotes.
func caption(text string) *widget.Label {
	l := widget.NewLabel(text)
	l.SizeName = theme.SizeNameCaptionText
	l.Wrapping = fyne.TextWrapWord
	return l
}

// orDefault returns fallback when s is empty.
func orDefault(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

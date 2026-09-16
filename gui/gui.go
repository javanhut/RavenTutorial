package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
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
	a.Settings().SetTheme(currentTheme)

	w := a.NewWindow("Raven Linux Tutorial")
	w.Resize(fyne.NewSize(900, 680))
	if icon := ravenLogo(false); icon != nil {
		w.SetIcon(icon)
	}

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

// ---------------------------------------------------------------------------
// The pages before the tour proper.
// ---------------------------------------------------------------------------

// showWelcome is the first thing the user sees: the machine's own bird, and
// one button.
func (t *tutorial) showWelcome() {
	greeting := canvas.NewText("Welcome to", dim())
	greeting.TextSize = theme.Size(theme.SizeNameSubHeadingText)
	greeting.Alignment = fyne.TextAlignCenter

	name := canvas.NewText("Raven Linux", theme.Color(theme.ColorNameForeground))
	name.TextSize = 42
	name.TextStyle = fyne.TextStyle{Bold: true}
	name.Alignment = fyne.TextAlignCenter

	about := widget.NewLabel(
		"A tour of the desktop, its keys and gestures, and the " +
			"applications this machine came with.\n\n" +
			"About ten minutes. Everything in it can be tried as you read - " +
			"the tutorial stays open beside whatever you start.")
	about.Wrapping = fyne.TextWrapWord
	about.Alignment = fyne.TextAlignCenter

	begin := widget.NewButtonWithIcon("Begin", theme.NavigateNextIcon(), t.showDeviceChoice)
	begin.Importance = widget.HighImportance

	t.onNext, t.onBack = t.showDeviceChoice, nil

	t.setContent(container.NewPadded(centred(
		container.NewCenter(logo(128)),
		greeting,
		name,
		widget.NewSeparator(),
		about,
		container.NewCenter(begin),
		container.NewCenter(hint("Right arrow, Return or Space turns the page")),
	)))
}

// showDeviceChoice asks which kind of machine the user is on, which decides
// the tutorial track they get.
func (t *tutorial) showDeviceChoice() {
	question := widget.NewLabel(
		"One chapter of the tour differs. A laptop is taught the trackpad " +
			"gestures; a desktop is taught the mouse chords that do the " +
			"same things. Everything else is the same either way.")
	question.Wrapping = fyne.TextWrapWord
	question.Alignment = fyne.TextAlignCenter

	laptop := choice(theme.ComputerIcon(), "Laptop",
		"Three fingers on the trackpad drive the desktop.", func() {
			t.startTrack("Laptop", LaptopSteps)
		})
	desktop := choice(theme.DesktopIcon(), "Desktop",
		"Super and the mouse buttons do the same work.", func() {
			t.startTrack("Desktop", DesktopSteps)
		})

	back := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), t.showWelcome)
	back.Importance = widget.LowImportance

	// Next is deliberately not bound here: which track to start is the one
	// question the tutorial cannot answer on the user's behalf.
	t.onNext, t.onBack = nil, t.showWelcome

	t.setContent(container.NewPadded(container.NewBorder(
		nil,
		container.NewBorder(nil, nil, back, nil),
		nil, nil,
		centred(
			container.NewCenter(eyebrow("Before we start")),
			container.NewCenter(title("Pick Your Machine")),
			question,
			widget.NewSeparator(),
			container.New(layout.NewGridLayout(2),
				container.NewPadded(laptop), container.NewPadded(desktop)),
		),
	)))
}

// startTrack loads the chosen set of steps and shows the first one.
func (t *tutorial) startTrack(device string, steps []Step) {
	t.device = device
	t.steps = steps
	t.index = 0
	t.showStep()
}

// ---------------------------------------------------------------------------
// A page of the tour.
// ---------------------------------------------------------------------------

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

	page := container.NewVScroll(container.NewPadded(readable(t.stepBody(step))))
	page.ScrollToTop()

	t.setContent(container.NewBorder(
		container.NewPadded(readable(t.stepHeader(step))),
		container.NewPadded(readable(t.footer(back, next, last))),
		nil, nil,
		page,
	))
}

// stepHeader is where in the tour this page is, what it is called, and how
// much of the tour is behind it.
func (t *tutorial) stepHeader(step Step) fyne.CanvasObject {
	where := fmt.Sprintf("%s  ·  step %d of %d", t.device, t.index+1, len(t.steps))

	return container.NewVBox(
		eyebrow(where),
		title(step.Title),
		stepProgress(t.index, len(t.steps)),
	)
}

// footer is the row that turns the page.
func (t *tutorial) footer(back, next func(), last bool) fyne.CanvasObject {
	backButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), back)
	backButton.Importance = widget.LowImportance

	// On the last step the Next button becomes Finish.
	nextLabel, nextIcon := "Next", theme.NavigateNextIcon()
	if last {
		nextLabel, nextIcon = "Finish", theme.ConfirmIcon()
	}
	nextButton := widget.NewButtonWithIcon(nextLabel, nextIcon, next)
	nextButton.Importance = widget.HighImportance

	contents := widget.NewButtonWithIcon("Contents", theme.ListIcon(), t.showContents)
	contents.Importance = widget.LowImportance

	return container.NewVBox(
		widget.NewSeparator(),
		container.NewBorder(nil, nil, backButton, nextButton, container.NewCenter(contents)),
	)
}

// stepBody assembles everything below the title: the prose, the keys, the
// applications, the footnote and whatever buttons the step asked for.
//
// Each of those is a section on its own card. The point of the cards is not
// decoration: a page of this tutorial is prose *and* a reference table, and
// putting the table on its own surface is what stops the reader having to
// work out which of the two they are looking at.
func (t *tutorial) stepBody(step Step) fyne.CanvasObject {
	page := container.NewVBox(body(step.Body))

	if len(step.Bindings) > 0 {
		page.Add(card(container.NewVBox(
			sectionTitle(orDefault(step.BindingsTitle, "Shortcuts")),
			bindingTable(step.Bindings),
		)))
	}

	for _, a := range step.Apps {
		page.Add(t.appCard(a))
	}

	if step.Launch != nil {
		// Not named app: this file imports a package called that.
		if program := (App{Exec: step.Launch.Exec, Args: step.Launch.Args}); program.Installed() {
			button := widget.NewButtonWithIcon(step.Launch.Label, theme.MediaPlayIcon(), func() {
				t.start(program)
			})
			button.Importance = widget.HighImportance
			page.Add(container.NewCenter(button))
		}
	}

	if step.Panel == "oracle" {
		page.Add(t.oraclePanel(oracleInstalled()))
	}

	if step.Note != "" {
		page.Add(t.noteCard(step.Note))
	}

	return page
}

// bindingTable draws the keys on the left and what they do on the right.
//
// A form layout rather than a grid: the keys column takes the width its
// widest chord needs and no more, which lines the descriptions up without
// giving half the page to "Esc".
func bindingTable(bindings []Binding) fyne.CanvasObject {
	rows := make([]fyne.CanvasObject, 0, len(bindings)*2)
	for _, b := range bindings {
		does := widget.NewLabel(b.Does)
		does.Wrapping = fyne.TextWrapWord
		// Neither cell is centred: the chord reads from the left like the
		// keyboard it describes, and the description from the left like the
		// prose above it.
		rows = append(rows, keycaps(b.Keys), does)
	}
	return container.New(layout.NewFormLayout(), rows...)
}

// appCard is one application: its own icon, its name, what it is for, and a
// button that opens it.
//
// A program that is not installed still gets its card - knowing what exists
// is half the point of the page - but says so instead of offering a button
// that could only fail.
func (t *tutorial) appCard(a App) fyne.CanvasObject {
	// A label rather than canvas text, so that the name and the line under
	// it start at the same place: a widget label carries the toolkit's inner
	// padding and raw canvas text does not.
	name := widget.NewLabelWithStyle(a.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	what := widget.NewLabel(a.What)
	what.Wrapping = fyne.TextWrapWord

	var action fyne.CanvasObject
	if a.Installed() {
		label := "Open"
		if a.Terminal {
			label = "Terminal"
		}
		program := a
		button := widget.NewButtonWithIcon(label, theme.MediaPlayIcon(), func() {
			t.start(program)
		})
		action = container.NewCenter(button)
	} else {
		action = container.NewCenter(hint("not installed"))
	}

	return card(container.NewBorder(
		nil, nil,
		container.NewCenter(avatar(a, 40)), action,
		container.NewVBox(name, what),
	))
}

// noteCard is the page's closing caveat, set apart from the body so that it
// reads as an aside rather than as the next paragraph.
func (t *tutorial) noteCard(note string) fyne.CanvasObject {
	return card(container.NewVBox(
		badge("Note", theme.ColorNamePrimary),
		secondary(note),
	))
}

// start runs a program and says so if it could not.
func (t *tutorial) start(a App) {
	if err := a.Start(); err != nil {
		dialog.ShowError(err, t.win)
	}
}

// ---------------------------------------------------------------------------
// Moving around, and the end.
// ---------------------------------------------------------------------------

// showContents opens the list of steps in this track, so a reader can go
// straight to the page they half-remember rather than pressing Next eleven
// times to reach it.
func (t *tutorial) showContents() {
	list := widget.NewList(
		func() int { return len(t.steps) },
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, widget.NewLabel("00"), nil, widget.NewLabel(""))
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			row := o.(*fyne.Container)
			number := row.Objects[1].(*widget.Label)
			name := row.Objects[0].(*widget.Label)

			number.SetText(fmt.Sprintf("%d", id+1))
			number.TextStyle = fyne.TextStyle{Monospace: true}
			name.SetText(t.steps[id].Title)
			name.TextStyle = fyne.TextStyle{Bold: id == t.index}
			name.Refresh()
		},
	)
	list.Select(t.index)

	jump := dialog.NewCustom("Contents", "Close", list, t.win)
	jump.Resize(fyne.NewSize(480, 520))
	list.OnSelected = func(id widget.ListItemID) {
		if id == t.index {
			return
		}
		jump.Hide()
		t.index = id
		t.showStep()
	}
	jump.Show()
}

// showFinish is the closing page once a track has been walked through.
func (t *tutorial) showFinish() {
	done := widget.NewLabel(
		"Nothing in the tour needs remembering. Super + Ctrl + H lists " +
			"every shortcut the desktop answers, at any time, from " +
			"anywhere.")
	done.Wrapping = fyne.TextWrapWord
	done.Alignment = fyne.TextAlignCenter

	other, otherSteps := "Desktop", DesktopSteps
	if t.device == "Desktop" {
		other, otherSteps = "Laptop", LaptopSteps
	}

	switchTrack := widget.NewButtonWithIcon("Take the "+other+" Track",
		theme.ComputerIcon(), func() {
			t.startTrack(other, otherSteps)
		})
	switchTrack.Importance = widget.HighImportance

	restart := widget.NewButtonWithIcon("Start Over", theme.ViewRefreshIcon(), t.showWelcome)
	contents := widget.NewButtonWithIcon("Contents", theme.ListIcon(), t.showContents)

	t.onNext, t.onBack = nil, func() {
		t.index = len(t.steps) - 1
		t.showStep()
	}

	t.setContent(container.NewPadded(centred(
		container.NewCenter(logo(96)),
		container.NewCenter(eyebrow(t.device+" tour complete")),
		container.NewCenter(title("All Done")),
		done,
		widget.NewSeparator(),
		container.NewCenter(container.NewHBox(restart, switchTrack, contents)),
		container.NewCenter(hint("Left arrow goes back to the last page")),
	)))
}

// hint is the smallest text on a page: the aside about which key does this.
func hint(text string) fyne.CanvasObject {
	label := canvas.NewText(text, dim())
	label.TextSize = theme.Size(theme.SizeNameCaptionText)
	label.Alignment = fyne.TextAlignCenter
	return label
}

// orDefault returns fallback when s is empty.
func orDefault(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

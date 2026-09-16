package gui

import (
	"image/color"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// Every page has to say something, and every table row has to have both of
// its columns. An empty field is invisible in the window rather than an
// error, which is exactly the kind of mistake that reaches a user.
func TestStepsAreComplete(t *testing.T) {
	for _, track := range map[string][]Step{"laptop": LaptopSteps, "desktop": DesktopSteps} {
		for _, step := range track {
			if step.Title == "" {
				t.Errorf("a step has no title: %+v", step)
			}
			if strings.TrimSpace(step.Body) == "" {
				t.Errorf("%q has no body", step.Title)
			}
			for _, b := range step.Bindings {
				if b.Keys == "" || b.Does == "" {
					t.Errorf("%q has a half-empty binding: %+v", step.Title, b)
				}
			}
			for _, a := range step.Apps {
				if a.Name == "" || a.What == "" || a.Exec == "" {
					t.Errorf("%q has an incomplete application: %+v", step.Title, a)
				}
			}
		}
	}
}

// Two steps with the same title would make the Contents list ambiguous.
func TestTitlesAreUnique(t *testing.T) {
	for name, track := range map[string][]Step{"laptop": LaptopSteps, "desktop": DesktopSteps} {
		seen := map[string]bool{}
		for _, step := range track {
			if seen[step.Title] {
				t.Errorf("%s track repeats the title %q", name, step.Title)
			}
			seen[step.Title] = true
		}
	}
}

// The tracks are meant to differ in one chapter and agree everywhere else. A
// fix made to a shared page must reach both machines, and this is what says
// so if a chapter is ever accidentally added to only one of them.
func TestTracksDifferOnlyInTheDeviceChapter(t *testing.T) {
	shared := func(track, device []Step) []string {
		var out []string
		skip := map[string]bool{}
		for _, s := range device {
			skip[s.Title] = true
		}
		for _, s := range track {
			if !skip[s.Title] {
				out = append(out, s.Title)
			}
		}
		return out
	}

	laptop := shared(LaptopSteps, laptopSteps)
	desktop := shared(DesktopSteps, mouseSteps)
	if len(laptop) != len(desktop) {
		t.Fatalf("tracks share %d and %d pages: %v vs %v", len(laptop), len(desktop), laptop, desktop)
	}
	for i := range laptop {
		if laptop[i] != desktop[i] {
			t.Errorf("shared page %d differs: %q vs %q", i, laptop[i], desktop[i])
		}
	}
}

// Every track ends on the one optional thing the tutorial offers.
func TestOracleIsLast(t *testing.T) {
	for name, track := range map[string][]Step{"laptop": LaptopSteps, "desktop": DesktopSteps} {
		if last := track[len(track)-1]; last.Title != oracleStep.Title {
			t.Errorf("%s track ends on %q, not on Oracle", name, last.Title)
		}
	}
}

// Drawing every page of both tracks, which is what catches a widget combination
// that panics on a page nobody opened while developing.
func TestEveryPageDraws(t *testing.T) {
	a := test.NewApp()
	defer test.NewApp()

	win := a.NewWindow("draw")
	tut := &tutorial{app: a, win: win}

	for device, track := range map[string][]Step{"Laptop": LaptopSteps, "Desktop": DesktopSteps} {
		tut.device, tut.steps = device, track
		for i := range track {
			tut.index = i
			body := tut.stepBody(track[i])
			win.SetContent(container.NewBorder(tut.stepHeader(track[i]), nil, nil, nil, body))
			win.Resize(fyne.NewSize(860, 640))
			refreshDeep(win.Content())
		}
	}
}

// Both halves of the Oracle page draw: the offer to install it, and the offer
// to open it. Which one a machine sees depends on what is installed on it, so
// neither is covered by running the tutorial here.
func TestBothOraclePanelsDraw(t *testing.T) {
	a := test.NewApp()
	defer test.NewApp()

	win := a.NewWindow("oracle")
	tut := &tutorial{app: a, win: win}

	for _, installed := range []bool{true, false} {
		win.SetContent(container.NewVScroll(tut.oraclePanel(installed)))
		win.Resize(fyne.NewSize(860, 640))
		refreshDeep(win.Content())
	}
}

// shellWords builds the line handed to sh -c, so a word it fails to quote is
// a word the shell would re-read as syntax.
func TestShellWordsQuotes(t *testing.T) {
	got := shellWords([]string{"oracle", "ask", "why won't it start; rm -rf /", "it's"})
	want := `'oracle' 'ask' 'why won'\''t it start; rm -rf /' 'it'\''s'`
	if got != want {
		t.Errorf("shellWords:\n got %s\nwant %s", got, want)
	}
}

// The accent comes out of the desktop's own configuration file, so the three
// spellings that file allows all have to read the same.
func TestParseHexColor(t *testing.T) {
	cases := map[string]color.NRGBA{
		"#22C5DD":   {R: 0x22, G: 0xc5, B: 0xdd, A: 0xff},
		"#22c5dd":   {R: 0x22, G: 0xc5, B: 0xdd, A: 0xff},
		"#7af":      {R: 0x77, G: 0xaa, B: 0xff, A: 0xff},
		"#8022C5DD": {R: 0x22, G: 0xc5, B: 0xdd, A: 0x80},
		" #22C5DD ": {R: 0x22, G: 0xc5, B: 0xdd, A: 0xff},
	}
	for text, want := range cases {
		got, ok := parseHexColor(text)
		if !ok || got != want {
			t.Errorf("parseHexColor(%q) = %v, %v; want %v", text, got, ok, want)
		}
	}
	for _, bad := range []string{"", "#", "22C5DD!", "#12345", "not a colour"} {
		if _, ok := parseHexColor(bad); ok {
			t.Errorf("parseHexColor(%q) accepted a value it should not have", bad)
		}
	}
}

// A label on the accent has to be readable whatever accent the person chose,
// which means it cannot be one fixed colour.
func TestReadableOnPicksContrast(t *testing.T) {
	dark := readableOn(color.NRGBA{R: 0x2a, G: 0x2a, B: 0x80, A: 0xff})
	if dark.R < 0xc0 {
		t.Errorf("a deep accent should take light text, got %v", dark)
	}
	light := readableOn(color.NRGBA{R: 0x22, G: 0xc5, B: 0xdd, A: 0xff})
	if light.R > 0x40 {
		t.Errorf("a bright accent should take dark text, got %v", light)
	}
}

// A chord is drawn as the keys it is made of, and a motion as one thing.
func TestKeycapsSplitsChords(t *testing.T) {
	chord := keycaps("Super + Ctrl + H").(*fyne.Container)
	// Three caps, and a plus between each pair.
	if len(chord.Objects) != 5 {
		t.Errorf("Super + Ctrl + H drew %d pieces, want 5", len(chord.Objects))
	}
	motion := keycaps("three fingers sideways").(*fyne.Container)
	if len(motion.Objects) != 1 {
		t.Errorf("a motion drew %d pieces, want 1", len(motion.Objects))
	}
}

// The progress track fills the width it is given whatever the step count, so
// a sixteen-page tour does not turn into a row of dots on a narrow window.
func TestProgressTrackFillsItsWidth(t *testing.T) {
	segments := make([]fyne.CanvasObject, 16)
	for i := range segments {
		bar := canvas.NewRectangle(color.Transparent)
		bar.SetMinSize(fyne.NewSize(6, 4))
		segments[i] = bar
	}

	track := &trackLayout{}
	track.Layout(segments, fyne.NewSize(500, 4))

	first, last := segments[0], segments[len(segments)-1]
	if first.Position().X != 0 {
		t.Errorf("the track starts at %v, not at the left edge", first.Position().X)
	}
	if end := last.Position().X + last.Size().Width; end < 499 || end > 501 {
		t.Errorf("the track ends at %v, not at the full 500 it was given", end)
	}
	if first.Size().Width != last.Size().Width {
		t.Errorf("segments differ in width: %v and %v", first.Size().Width, last.Size().Width)
	}
}

// Prose is held to a readable column however wide the window gets.
func TestReadableCapsTheColumn(t *testing.T) {
	content := widget.NewLabel("some prose")
	column := &readableLayout{}
	objects := []fyne.CanvasObject{content}

	column.Layout(objects, fyne.NewSize(1600, 400))
	if content.Size().Width != readableWidth {
		t.Errorf("a wide window gave the column %v, want %v", content.Size().Width, readableWidth)
	}
	if left := content.Position().X; left != (1600-readableWidth)/2 {
		t.Errorf("the column sits at %v rather than centred", left)
	}

	column.Layout(objects, fyne.NewSize(400, 400))
	if content.Size().Width != 400 {
		t.Errorf("a narrow window gave the column %v, want all 400", content.Size().Width)
	}
}

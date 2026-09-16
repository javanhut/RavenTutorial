package gui

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
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

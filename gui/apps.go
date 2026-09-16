package gui

import (
	"os/exec"
	"strings"
)

// App is one program in the tour of what Raven Linux ships with.
//
// Exec is the command, which is also how the tutorial decides whether to
// offer to open it: a machine installed from the minimal profile has fewer of
// these than a desktop one, and a button for something that is not there
// would be a broken promise on the page.
type App struct {
	Name string
	What string
	Exec string
	Args []string

	// Terminal marks a command-line program, which is opened inside
	// raven-terminal rather than on its own.
	Terminal bool
}

// Installed reports whether this program can actually be started here.
//
// A terminal program needs the terminal as well as itself: raven-terminal is
// how it gets a window.
func (a App) Installed() bool {
	if _, err := exec.LookPath(a.Exec); err != nil {
		return false
	}
	if a.Terminal {
		if _, err := exec.LookPath(terminal); err != nil {
			return false
		}
	}
	return true
}

// Start runs the program, detached, and returns without waiting for it.
//
// A tutorial that blocked while the file manager was open would be a tutorial
// you had to close to carry on reading.
func (a App) Start() error {
	name, args := a.Exec, a.Args
	switch {
	case a.Terminal && len(a.Args) > 0:
		// Help text in a terminal that closed the instant it finished
		// printing would be unreadable, so the shell takes over the
		// window afterwards and you are left somewhere to try the
		// command yourself.
		line := shellWords(append([]string{a.Exec}, a.Args...)) + `; exec ${SHELL:-/bin/sh}`
		name, args = terminal, []string{"-e", "sh", "-c", line}
	case a.Terminal:
		name, args = terminal, append([]string{"-e", a.Exec}, a.Args...)
	}
	cmd := exec.Command(name, args...)
	if err := cmd.Start(); err != nil {
		return err
	}
	// Nothing here ever reads the exit status, and a child nobody waits on
	// stays a zombie for the life of the tutorial. Release lets the process
	// go entirely.
	return cmd.Process.Release()
}

// terminal is what a command-line program is opened inside. The desktop names
// this same binary in two compiled-in places, so it is the one terminal a
// Raven machine is guaranteed to have.
const terminal = "raven-terminal"

// shellWords joins a command into one line for sh -c, quoting each word so
// that nothing in it is re-read as shell syntax.
func shellWords(words []string) string {
	var b strings.Builder
	for i, w := range words {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteByte('\'')
		b.WriteString(strings.ReplaceAll(w, "'", `'\''`))
		b.WriteByte('\'')
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// The applications, grouped the way the tour introduces them.
// ---------------------------------------------------------------------------

var everydayApps = []App{
	{
		Name: "Terminal",
		What: "A GPU-accelerated terminal. Super + Ctrl + E opens one from anywhere.",
		Exec: "raven-terminal",
	},
	{
		Name: "Files",
		What: "Browse, search, preview and organise. Tabs, a sidebar and a preview panel.",
		Exec: "ravenfilemanager",
	},
	{
		Name: "EagleEye",
		What: "Images: PNG, JPEG, WebP, GIF, SVG, HEIC, AVIF and more. The bytes decide the format, so a PNG saved as .jpg still opens.",
		Exec: "eagleeye",
	},
	{
		Name: "Raven Viewer",
		What: "PDF and DOCX, with the document outline in the sidebar and annotations listed with their pages.",
		Exec: "raven-viewer",
	},
	{
		Name: "Owl Player",
		What: "Video and audio, built straight on FFmpeg with its own GPU renderer.",
		Exec: "owl-player",
	},
	{
		Name:     "Crow",
		What:     "A selection-first modal text editor, in the terminal.",
		Exec:     "crow",
		Terminal: true,
	},
}

var systemApps = []App{
	{
		Name: "Settings",
		What: "Wi-Fi, Bluetooth, sound, screens, theme and wallpaper, the dock and bar, default applications, storage, privacy and updates. Super + Ctrl + P.",
		Exec: "raven-settings",
	},
	{
		Name: "Raven Store",
		What: "Browse, install and update software without a terminal. It runs rvn underneath, so it and the command line always agree. Super + Ctrl + I.",
		Exec: "raven-store",
	},
	{
		Name: "Raven Power",
		What: "Battery telemetry, runtime estimates, power profiles and battery health.",
		Exec: "raven-power",
	},
	{
		Name: "Controls",
		What: "Keyboard backlight, fan speeds and thermal profiles.",
		Exec: "raven-controls",
	},
}

var terminalTools = []App{
	{
		Name:     "rvn",
		What:     "The package manager. install, uninstall, update, find, info, list, owns, files, sync - official repositories and the AUR in one binary, with no pacman or yay underneath.",
		Exec:     "rvn",
		Args:     []string{"--help"},
		Terminal: true,
	},
	{
		Name:     "caw",
		What:     "Wi-Fi. Speaks to the kernel directly and does WPA itself, so it replaces iw, iwctl, wpa_supplicant and dhcpcd outright.",
		Exec:     "caw",
		Args:     []string{"--help"},
		Terminal: true,
	},
	{
		Name:     "ivaldi",
		What:     "Version control - the one Raven Linux is itself developed in.",
		Exec:     "ivaldi",
		Args:     []string{"--help"},
		Terminal: true,
	},
	{
		Name:     "imlazy",
		What:     "The task runner. Projects describe their commands in lazy.toml and imlazy runs them.",
		Exec:     "imlazy",
		Args:     []string{"--help"},
		Terminal: true,
	},
	{
		Name:     "poxy",
		What:     "A universal package manager, for what rvn does not cover.",
		Exec:     "poxy",
		Args:     []string{"--help"},
		Terminal: true,
	},
	{
		Name:     "oxigen",
		What:     "The Oxigen language, shipped as a static interpreter that needs nothing else installed.",
		Exec:     "oxigen",
		Args:     []string{"--help"},
		Terminal: true,
	},
}

// ---------------------------------------------------------------------------
// The pages that show them.
// ---------------------------------------------------------------------------

var applicationSteps = []Step{
	{
		Title: "The Applications",
		Body: "Raven Linux writes its own applications rather than " +
			"assembling somebody else's. They share one look - Raven Glass, " +
			"the translucent panels you have been using all tour - and they " +
			"follow the theme, accent colour and transparency you set in " +
			"Settings.\n\n" +
			"Here is what came with your machine. Open any of them from " +
			"this page; they will appear in their own window beside the " +
			"tutorial.",
		Apps: everydayApps,
		Note: "Anything missing from this list can be installed from the " +
			"Raven Store, or with rvn install from a terminal.",
	},
	{
		Title: "Running the Machine",
		Body: "Four applications between them cover everything you would " +
			"otherwise set with a scatter of command-line tools.\n\n" +
			"Quick settings, on Super + Ctrl + S, is the short version of " +
			"the first of them, and its last row opens the full application.",
		Apps: systemApps,
	},
	{
		Title: "At the Command Line",
		Body: "Your shell is ravenshell, and the tools around it are Raven's " +
			"own too. None of them is a wrapper around something else - the " +
			"package manager is not pacman underneath, and the Wi-Fi tool is " +
			"not wpa_supplicant underneath.\n\n" +
			"Open any of them here to read its help in a terminal.",
		Apps: terminalTools,
		Note: "Two more programs you will not run by hand: roostbar draws " +
			"the status bar across the top of the screen, and ravencanvas " +
			"draws the wallpaper - including animated and computed ones. " +
			"Both are configured from Settings.",
	},
}

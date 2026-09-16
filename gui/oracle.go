package gui

import (
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Oracle is the one optional thing this tutorial offers to install.
//
// It is deliberately not part of Raven Linux: nothing in the base system
// requires it, references it, or knows it exists, it is in no install
// profile, and it adds no service, no autostart entry and no background
// process. It runs when you run it and at no other time. That is why this
// page offers it rather than assuming it, and why the tutorial never installs
// it without being asked.
const (
	oracleApp  = "raven-oracle"
	oracleCLI  = "oracle"
	oracleRepo = "https://github.com/javanhut/OracleAgent.git"
)

// oracleInstalled reports whether either half of Oracle is already here.
func oracleInstalled() bool {
	for _, bin := range []string{oracleApp, oracleCLI} {
		if _, err := exec.LookPath(bin); err == nil {
			return true
		}
	}
	return false
}

// oracleInstall is what to run to get it: clone the repository and let its
// own task runner install it, exactly as its README says. imlazy install
// builds both halves; imlazy install-cli builds the command alone, without
// GTK, for a machine with no desktop on it.
//
// The tutorial shows these rather than running them. The install step needs
// root, and a tutorial that silently asked for a password in a window you
// could not see would be worse than one that hands you the four lines.
var oracleInstall = strings.Join([]string{
	"git clone " + oracleRepo,
	"cd OracleAgent",
	"imlazy install          # or: sudo make install",
}, "\n")

// oracleModel is the optional second half: Oracle's checks need no model at
// all, but its questions do. It downloads nothing by itself and installs no
// server, so this is a separate, deliberate step.
var oracleModel = strings.Join([]string{
	"rvn install ollama",
	"ollama serve &",
	"ollama pull qwen2.5:3b-instruct",
}, "\n")

// oracleRemove is on the page for the same reason the install is: something
// you are asked to add should come with the line that takes it away again.
var oracleRemove = strings.Join([]string{
	"oracle forget           # delete the config and anything downloaded",
	"imlazy uninstall        # remove the binaries, launcher entry and icon",
}, "\n")

// oracleStep is the last page of every track.
var oracleStep = Step{
	Title: "Raven Oracle (optional)",
	Body: "Oracle is a local troubleshooting companion. It reads your " +
		"machine, tells you what looks wrong, and suggests what to try - " +
		"and it does that with no model, no downloaded weights and no " +
		"network at all. That is the floor rather than a degraded mode.\n\n" +
		"It finds the class of failure nothing else reports: a service " +
		"whose executable is missing and which fails quietly at every " +
		"boot, firmware the kernel asked for and did not get, a package " +
		"database two tools disagree about.\n\n" +
		"It is not part of Raven Linux, and that is the design. It has no " +
		"autostart entry, no tray icon and no background service. Nothing " +
		"runs when you are not running it, and it never runs a command for " +
		"you - it copies one to the clipboard and you decide.\n\n" +
		"If you give it a local model, it will also answer questions about " +
		"this machine in prose, reading the system afresh on every turn so " +
		"that \"did that fix it?\" is a question it can actually answer.",
	BindingsTitle: "What it can do",
	Bindings: []Binding{
		{"raven-oracle", "the desktop app: findings, ask, explain an error, settings"},
		{"oracle doctor", "check the machine and report, as plain text"},
		{"oracle tui", "the same findings full-screen in a terminal"},
		{"oracle ask \"...\"", "ask about this machine, and keep asking"},
		{"oracle explain", "paste or pipe in what a failed command printed"},
		{"oracle context", "show exactly what would be sent to a model, before anything is"},
		{"oracle status", "what is configured and what is missing"},
	},
	Note: "Installing it puts two binaries, a launcher entry, its metainfo " +
		"and an icon under /usr/local. Nothing in /etc is touched and " +
		"nothing is enabled.",
	Panel: "oracle",
}

// ---------------------------------------------------------------------------
// The panel drawn under the Oracle step.
// ---------------------------------------------------------------------------

// oraclePanel is the only hand-built section in the tutorial, because it is
// the only page that has to say something different depending on what is on
// the machine: an offer to install, or an offer to open.
// installed is passed in rather than looked up here so that both halves can
// be drawn in a test on a machine that has only one of them.
func (t *tutorial) oraclePanel(installed bool) fyne.CanvasObject {
	if installed {
		return t.oracleHere()
	}
	return t.oracleMissing()
}

// oracleHere is the panel when Oracle is already installed.
func (t *tutorial) oracleHere() fyne.CanvasObject {
	status := widget.NewLabelWithStyle(
		"Oracle is installed on this machine.",
		fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	buttons := container.NewHBox()
	for _, a := range []App{
		{Name: "Raven Oracle", Exec: oracleApp},
		{Name: "oracle doctor", Exec: oracleCLI, Args: []string{"doctor"}, Terminal: true},
	} {
		if !a.Installed() {
			continue
		}
		program := a
		label := "Open Raven Oracle"
		if program.Terminal {
			label = "Check This Machine"
		}
		buttons.Add(widget.NewButtonWithIcon(label, theme.MediaPlayIcon(), func() {
			if err := program.Start(); err != nil {
				dialog.ShowError(err, t.win)
			}
		}))
	}

	return container.NewVBox(
		status,
		buttons,
		t.commandBlock("Give it a local model, to ask it questions", oracleModel),
		t.commandBlock("If you want it gone again", oracleRemove),
	)
}

// oracleMissing is the panel when it is not installed, which is what a fresh
// Raven install looks like.
func (t *tutorial) oracleMissing() fyne.CanvasObject {
	status := widget.NewLabel(
		"Oracle is not installed, which is the state a fresh Raven Linux " +
			"comes in. Nothing on the system needs it.\n\n" +
			"These commands fetch and install it. The tutorial does not run " +
			"them for you - installing needs root, and a password prompt " +
			"you cannot see is worse than four lines you can read.")
	status.Wrapping = fyne.TextWrapWord

	return container.NewVBox(
		status,
		t.commandBlock("Install it", oracleInstall),
		t.commandBlock("Then, optionally, a local model for its questions", oracleModel),
	)
}

// commandBlock shows a few shell lines with a button that copies them and,
// where there is a terminal to open, one that opens it.
func (t *tutorial) commandBlock(title, commands string) fyne.CanvasObject {
	text := widget.NewLabelWithStyle(commands, fyne.TextAlignLeading,
		fyne.TextStyle{Monospace: true})

	copyButton := widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {
		t.app.Clipboard().SetContent(commands)
	})

	buttons := container.NewHBox(copyButton)
	if _, err := exec.LookPath(terminal); err == nil {
		buttons.Add(widget.NewButtonWithIcon("Open a Terminal", theme.ComputerIcon(), func() {
			if err := (App{Exec: terminal}).Start(); err != nil {
				dialog.ShowError(err, t.win)
			}
		}))
	}

	return widget.NewCard(title, "", container.NewVBox(text, buttons))
}

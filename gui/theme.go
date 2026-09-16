package gui

import (
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Raven Glass, as far as a toolkit that is not GTK can follow it.
//
// The desktop's own applications - Settings, Store, Power, Oracle - share one
// stylesheet, and its rules are worth restating because they are what this
// file implements:
//
//  1. One material. Surfaces are translucent layers of the same glass, and
//     hierarchy comes from how much light a layer lets through rather than
//     from a different colour per surface.
//  2. A hairline, never a heavy border: white at about 8% alpha.
//  3. One radius scale - 6 for controls, 8 for rows, 12 for groups, 14 for
//     cards, 20 for heroes. Nothing else.
//  4. Type earns hierarchy with weight, not with size alone; secondary text
//     is the same face at reduced opacity.
//  5. The accent is the person's, and it does exactly two jobs: the selected
//     thing, and the primary action. Everything else is grey.
//
// Rule 5 is why this file reads ~/.config/raven/desktop.toml. A tutorial that
// introduced somebody to their desktop in a colour their desktop does not use
// would be teaching the wrong thing on its very first page.
const (
	glassDarkBackground = 0x17171dff
	glassDarkForeground = 0xe8e8f0ff
	glassDarkInput      = 0x1c1c23ff
	glassDarkMenu       = 0x26262fff

	glassLightBackground = 0xf2f2f7ff
	glassLightForeground = 0x1c1c22ff
	glassLightInput      = 0xffffffff
	glassLightMenu       = 0xffffffff

	// The fallback accent is the one the stylesheet itself ships with, for a
	// machine whose desktop.toml has not been written yet.
	defaultAccent = 0x7aa2f7ff
)

// desktopSettings is the part of the desktop's own configuration this
// tutorial cares about.
type desktopSettings struct {
	accent color.Color
	light  bool
}

// readDesktopSettings reads ~/.config/raven/desktop.toml, the file Raven
// Settings writes and the desktop reads.
//
// Every failure lands on the defaults, on purpose. This is a look, not a
// policy: a machine with no file yet, or a file half-saved while somebody is
// editing it, should still get a tutorial that opens and is legible.
func readDesktopSettings() desktopSettings {
	settings := desktopSettings{accent: rgba(defaultAccent)}

	home, err := os.UserHomeDir()
	if err != nil {
		return settings
	}
	data, err := os.ReadFile(filepath.Join(home, ".config", "raven", "desktop.toml"))
	if err != nil {
		return settings
	}

	// A five-line reader rather than a TOML dependency: two keys are wanted,
	// both are plain strings in a known section, and a parser that cannot
	// fail is the right shape for something whose failure mode is "the
	// tutorial is the wrong blue".
	section := ""
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") {
			section = strings.Trim(line, "[]")
			continue
		}
		if section != "appearance" {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"`)
		switch key {
		case "accent":
			if c, ok := parseHexColor(value); ok {
				settings.accent = c
			}
		case "theme_mode":
			settings.light = value == "light"
		}
	}
	return settings
}

// parseHexColor reads #RGB, #RRGGBB or #AARRGGBB, the three forms the
// desktop's own configuration files accept.
func parseHexColor(s string) (color.NRGBA, bool) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(s) == 3 {
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	}
	if len(s) != 6 && len(s) != 8 {
		return color.NRGBA{}, false
	}
	n, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return color.NRGBA{}, false
	}
	c := color.NRGBA{R: uint8(n >> 16), G: uint8(n >> 8), B: uint8(n), A: 0xff}
	if len(s) == 8 {
		c = color.NRGBA{R: uint8(n >> 16), G: uint8(n >> 8), B: uint8(n), A: uint8(n >> 24)}
	}
	return c, true
}

// currentTheme is what this run is dressed in, read once at startup.
//
// The pieces in widgets.go ask it whether the desktop is in light or dark
// mode, because a few of them choose between two washes of the same material
// rather than between two named theme colours - and Fyne hands a widget the
// variant only while it is drawing itself, which is too late for something
// that is deciding what to build.
var currentTheme = ravenTheme{readDesktopSettings()}

// ravenTheme dresses Fyne in Raven Glass.
type ravenTheme struct {
	desktopSettings
}

var _ fyne.Theme = ravenTheme{}

// Color answers with the glass palette, in whichever mode the desktop is in.
//
// The variant Fyne offers is ignored: this desktop has one switch for light
// and dark, it is in desktop.toml, and a tutorial that disagreed with the
// desktop about which mode the machine is in would look like a bug.
func (t ravenTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	if t.light {
		return t.lightColor(name)
	}
	return t.darkColor(name)
}

func (t ravenTheme) darkColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return rgba(glassDarkBackground)
	case theme.ColorNameForeground:
		return rgba(glassDarkForeground)
	case theme.ColorNameButton, theme.ColorNameInputBackground:
		return white(0x14)
	case theme.ColorNameOverlayBackground, theme.ColorNameMenuBackground:
		return rgba(glassDarkMenu)
	case theme.ColorNameHeaderBackground:
		return rgba(glassDarkInput)
	case theme.ColorNameHover:
		return white(0x12)
	case theme.ColorNamePressed:
		return white(0x1f)
	case theme.ColorNameSeparator, theme.ColorNameInputBorder:
		return white(0x17)
	case theme.ColorNameDisabled:
		return fade(rgba(glassDarkForeground), 0x50)
	case theme.ColorNameDisabledButton:
		return white(0x0a)
	case theme.ColorNamePlaceHolder:
		return fade(rgba(glassDarkForeground), 0x8c)
	case theme.ColorNameShadow:
		return black(0x5c)
	case theme.ColorNameScrollBar:
		return white(0x2a)
	case theme.ColorNameScrollBarBackground:
		return color.Transparent
	}
	return t.sharedColor(name)
}

func (t ravenTheme) lightColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return rgba(glassLightBackground)
	case theme.ColorNameForeground:
		return rgba(glassLightForeground)
	case theme.ColorNameButton, theme.ColorNameInputBackground:
		return white(0xd9)
	case theme.ColorNameOverlayBackground, theme.ColorNameMenuBackground:
		return rgba(glassLightMenu)
	case theme.ColorNameHeaderBackground:
		return rgba(glassLightInput)
	case theme.ColorNameHover:
		return black(0x0f)
	case theme.ColorNamePressed:
		return black(0x1a)
	case theme.ColorNameSeparator, theme.ColorNameInputBorder:
		return black(0x17)
	case theme.ColorNameDisabled:
		return fade(rgba(glassLightForeground), 0x50)
	case theme.ColorNameDisabledButton:
		return black(0x0a)
	case theme.ColorNamePlaceHolder:
		return fade(rgba(glassLightForeground), 0x8c)
	case theme.ColorNameShadow:
		return black(0x1a)
	case theme.ColorNameScrollBar:
		return black(0x2a)
	case theme.ColorNameScrollBarBackground:
		return color.Transparent
	}
	return t.sharedColor(name)
}

// sharedColor is everything the two modes agree about: the accent's two jobs,
// and the three status colours the stylesheet names.
func (t ravenTheme) sharedColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNamePrimary, theme.ColorNameHyperlink:
		return t.accent
	case theme.ColorNameFocus:
		return fade(t.accent, 0x66)
	case theme.ColorNameSelection:
		return fade(t.accent, 0x40)
	case theme.ColorNameForegroundOnPrimary:
		return readableOn(t.accent)
	case theme.ColorNameSuccess:
		return rgba(0x34d399ff)
	case theme.ColorNameWarning:
		return rgba(0xfbbf24ff)
	case theme.ColorNameError, theme.ColorNameForegroundOnError:
		return rgba(0xfb7185ff)
	}
	return theme.DefaultTheme().Color(name, t.variant())
}

func (t ravenTheme) variant() fyne.ThemeVariant {
	if t.light {
		return theme.VariantLight
	}
	return theme.VariantDark
}

// Font and Icon are the toolkit's own. Raven's look is carried by colour,
// spacing and radius here; swapping the typeface as well would make the
// tutorial the one Raven application with different letterforms.
func (t ravenTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t ravenTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size is the stylesheet's one radius scale, and type that leans on weight
// rather than on size: a heading only a little larger than the body, and a
// body with enough line spacing to read a paragraph in.
func (t ravenTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 26
	case theme.SizeNameSubHeadingText:
		return 17
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameLineSpacing:
		return 5
	case theme.SizeNamePadding:
		return 5
	case theme.SizeNameInnerPadding:
		return 10
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameButtonRadius, theme.SizeNameSelectionRadius,
		theme.SizeNameInputRadius, theme.SizeNameScrollBarRadius:
		return 8
	case theme.SizeNameCardRadius:
		return 14
	case theme.SizeNameDialogRadius, theme.SizeNamePopupRadius, theme.SizeNameMenuRadius:
		return 12
	case theme.SizeNameScrollBar:
		return 10
	case theme.SizeNameScrollBarSmall:
		return 4
	}
	return theme.DefaultTheme().Size(name)
}

// ---------------------------------------------------------------------------
// Colour helpers.
// ---------------------------------------------------------------------------

// rgba turns an 0xRRGGBBAA literal into a colour.
func rgba(v uint32) color.NRGBA {
	return color.NRGBA{R: uint8(v >> 24), G: uint8(v >> 16), B: uint8(v >> 8), A: uint8(v)}
}

// white and black are the hairlines and washes the material is made of: one
// colour at a given alpha, over whatever is beneath it.
func white(alpha uint8) color.NRGBA { return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: alpha} }
func black(alpha uint8) color.NRGBA { return color.NRGBA{A: alpha} }

// fade returns a colour at the given alpha.
func fade(c color.Color, alpha uint8) color.NRGBA {
	r, g, b, _ := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: alpha}
}

// readableOn picks black or white text for a background, by its perceived
// brightness. An accent is the person's choice and can be anything - a pale
// cyan and a deep indigo both want a label on top, and they do not want the
// same one.
func readableOn(c color.Color) color.NRGBA {
	r, g, b, _ := c.RGBA()
	// Rec. 601 luma, in 8-bit terms.
	luma := (299*float64(r>>8) + 587*float64(g>>8) + 114*float64(b>>8)) / 1000
	if luma > 150 {
		return rgba(0x16161fff)
	}
	return rgba(0xffffffff)
}

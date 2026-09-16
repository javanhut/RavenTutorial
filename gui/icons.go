package gui

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
)

// Finding the icon a program is drawn with elsewhere on this desktop.
//
// The launcher and the dock draw Raven's applications with the icons their
// .desktop entries name, and a tutorial that drew the same applications as
// grey placeholders would look like a list of something else. So this does
// what they do, in the small: find the entry, read its Icon= key, and resolve
// that name against the icon directories.
//
// Every step of it is allowed to fail. A machine can have an application with
// no entry, an entry with no icon, or an icon name no theme on it provides,
// and each of those is a page with one fewer picture rather than an error.

// dataDirs are the roots that hold .desktop entries and icon themes, in
// precedence order: the user's own, then anything installed to /usr/local
// (which is where an optional application like Oracle puts itself), then the
// system's.
func dataDirs() []string {
	var dirs []string
	if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, filepath.Join(home, ".local", "share"))
	}
	if set := os.Getenv("XDG_DATA_DIRS"); set != "" {
		dirs = append(dirs, strings.Split(set, ":")...)
	} else {
		dirs = append(dirs, "/usr/local/share", "/usr/share")
	}
	return dirs
}

// appIcon finds the icon for a command, as an entry names it.
//
// Matching is on the first word of Exec rather than on the entry's file name:
// the tutorial knows programs by the command that starts them, and an entry
// is free to be called anything.
func appIcon(command string) fyne.Resource {
	name := iconName(command)
	if name == "" {
		return nil
	}
	return loadIcon(name)
}

// iconName reads the Icon= key of the entry that best matches a command.
//
// "Best" needs saying, because one command can appear in several entries.
// Crow's entry runs `raven-terminal -e crow`, so a naive match on the first
// word of Exec hands the terminal's icon to Crow, or Crow's to the terminal,
// depending on which file is read first. So every entry is scored and the
// strongest match wins:
//
//	3  the entry runs exactly this command and nothing else
//	2  the entry's command is this one, with arguments
//	1  this command appears later in the entry - it is run inside another
//
// Entries are read in precedence order and ties keep the first, so a user's
// own copy of an application still beats the system's.
func iconName(command string) string {
	best, icon := 0, ""
	for _, dir := range dataDirs() {
		entries, err := os.ReadDir(filepath.Join(dir, "applications"))
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if filepath.Ext(entry.Name()) != ".desktop" {
				continue
			}
			words, named := entryFields(filepath.Join(dir, "applications", entry.Name()))
			if named == "" || len(words) == 0 {
				continue
			}
			score := 0
			switch {
			case len(words) == 1 && words[0] == command:
				score = 3
			case words[0] == command:
				score = 2
			default:
				for _, word := range words[1:] {
					if word == command {
						score = 1
						break
					}
				}
			}
			if score > best {
				best, icon = score, named
			}
		}
	}
	return icon
}

// entryFields returns the words of an entry's Exec line and the icon it
// names.
//
// Only the first group of the file is read, which is the one that describes
// the application itself; the Desktop Action groups after it name their own
// Exec lines for menu items, and those are not what the tutorial launches.
// Field codes - %U, %F and the rest - are dropped, and each word is reduced
// to its base name so that an entry written with an absolute path still
// matches the command it runs.
func entryFields(path string) (exec []string, icon string) {
	file, err := os.Open(path)
	if err != nil {
		return nil, ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	group := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[") {
			if group != "" {
				break
			}
			group = line
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch strings.TrimSpace(key) {
		case "Exec":
			for _, word := range strings.Fields(value) {
				if strings.HasPrefix(word, "%") {
					continue
				}
				exec = append(exec, filepath.Base(word))
			}
		case "Icon":
			icon = strings.TrimSpace(value)
		}
	}
	return exec, icon
}

// loadIcon resolves an icon name to a file and reads it.
//
// An absolute path is taken as given, which the spec allows. Otherwise the
// name is looked for in the two directory layouts that exist on this machine
// - hicolor's <size>/apps and breeze's apps/<size> - largest first, since
// these are drawn at 32 to 48 points and scaling a 96 down beats scaling a 16
// up. /usr/share/pixmaps is the last resort, where entries that predate icon
// themes put their pictures.
func loadIcon(name string) fyne.Resource {
	return loadIconPreferring(name, ".svg", ".png")
}

// loadIconPreferring is loadIcon with the file formats tried in a given
// order.
//
// The order matters more than it should. Fyne rasterises SVG itself, and its
// rasteriser does not draw every SVG in existence - the machine's own logo,
// which has a gradient in it, comes out blank. So anything drawn large enough
// for that to be obvious asks for the PNG first, and the icons beside a line
// of text, where scaling matters more, ask for the SVG.
func loadIconPreferring(name string, exts ...string) fyne.Resource {
	if filepath.IsAbs(name) {
		return readIcon(name)
	}

	// Largest first: these are drawn between 32 and 128 points, and scaling a
	// 256 down beats scaling a 32 up.
	sizes := []string{"scalable", "256x256", "128x128", "96x96", "64x64", "48x48", "32x32"}
	themes := []string{"hicolor", "breeze-dark", "breeze", "Adwaita"}

	// The format is the outermost loop, so that asking for a PNG first really
	// does mean every PNG on the machine before any SVG - rather than the
	// first directory that holds either.
	for _, ext := range exts {
		for _, dir := range dataDirs() {
			for _, iconTheme := range themes {
				root := filepath.Join(dir, "icons", iconTheme)
				for _, size := range sizes {
					// hicolor keeps <size>/apps; breeze keeps apps/<n>,
					// where n is one number rather than "48x48".
					candidates := []string{filepath.Join(root, size, "apps", name+ext)}
					if n, _, ok := strings.Cut(size, "x"); ok {
						candidates = append(candidates, filepath.Join(root, "apps", n, name+ext))
					}
					for _, candidate := range candidates {
						if res := readIcon(candidate); res != nil {
							return res
						}
					}
				}
			}
			if res := readIcon(filepath.Join(dir, "pixmaps", name+ext)); res != nil {
				return res
			}
		}
	}
	return nil
}

// readIcon loads one file, or returns nil if it is not there.
func readIcon(path string) fyne.Resource {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return fyne.NewStaticResource(filepath.Base(path), data)
}

// ravenLogo is the machine's own logo, in the version meant for the mode the
// desktop is in: the light-inked bird for a dark background, and the dark one
// for a light background.
func ravenLogo(light bool) fyne.Resource {
	name := "raven-logo"
	if light {
		name = "raven-logo-dark"
	}
	return loadIconPreferring(name, ".png", ".svg")
}

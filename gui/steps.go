package gui

// Step is a single page of a tutorial track.
//
// A page is a title, some prose, and - for the pages that are really about
// the keyboard or the trackpad - a table of bindings. The table is kept out
// of Body deliberately: a chord belongs in a column where it can be scanned,
// not in the middle of a sentence where it has to be read.
type Step struct {
	Title string
	Body  string

	// Bindings is drawn as a two-column table under the body.
	Bindings []Binding

	// BindingsTitle labels that table. Empty means "Shortcuts".
	BindingsTitle string

	// Apps is drawn as a list of programs, each with a button that opens
	// it when it is installed on this machine.
	Apps []App

	// Note is a smaller closing line, for the caveat that would otherwise
	// bloat the body.
	Note string

	// Launch, when set, puts a button on the page that starts a program.
	// The button is left out entirely when the program is not installed,
	// so a page never offers something that cannot happen.
	Launch *Launch

	// Panel names a hand-built section to draw under everything else.
	// Only "oracle" exists; see oracle.go.
	Panel string
}

// Binding is one row of a step's table: the keys, and what they do.
type Binding struct {
	Keys string
	Does string
}

// Launch is a "try it" button: a program to start, and what to call it.
type Launch struct {
	Label string
	Exec  string
	Args  []string
}

// ---------------------------------------------------------------------------
// The opening chapters, shown to everyone.
// ---------------------------------------------------------------------------

var openingSteps = []Step{
	{
		Title: "Meet the Leader Key",
		Body: "Most of Raven Linux runs off one key: the Leader key.\n\n" +
			"It is the Super key - the one between Ctrl and Alt, usually " +
			"marked with a logo. Some keyboards call it the Windows key or " +
			"the Command key. Everywhere in this tutorial, Super means that " +
			"key.\n\n" +
			"Almost every desktop shortcut is Super and Ctrl together, then " +
			"a letter. That is not an accident. Super on its own belongs to " +
			"whatever application you are using, so the desktop taking a " +
			"chord can never take one out from under your editor or your " +
			"terminal.\n\n" +
			"There are three exceptions, and they are worth knowing now.",
		BindingsTitle: "The plain Super layer",
		Bindings: []Binding{
			{"Super + C", "copy, in whatever application is focused"},
			{"Super + V", "paste, in whatever application is focused"},
			{"Super + L", "lock the session"},
		},
		Note: "Super + C and Super + V are translated into Ctrl + C and " +
			"Ctrl + V for the focused application, because copying belongs " +
			"to the application rather than to the desktop. Terminals are " +
			"exempt - a terminal that received Ctrl + C when you meant copy " +
			"would interrupt whatever was running. Super + L is never handed " +
			"back to anyone: a lock chord that failed whenever a terminal " +
			"happened to be focused would fail at the exact moment somebody " +
			"walked away from the machine.",
	},
	{
		Title: "The Keybinding List",
		Body: "Raven Linux can show you every one of its shortcuts at any " +
			"time, so none of this tutorial is something you have to " +
			"memorise.\n\n" +
			"Press Super + Ctrl + H.\n\n" +
			"A panel opens listing every chord the desktop answers right " +
			"now. Press Esc, or click anywhere outside the panel, to close " +
			"it. If you forget a key on a later page, that panel is the " +
			"fastest way back to it.",
		BindingsTitle: "Getting help and getting out",
		Bindings: []Binding{
			{"Super + Ctrl + H", "show the keybinding list"},
			{"Esc", "close the list, the launcher, or any panel"},
			{"Super + Ctrl + Esc", "quit the desktop and return to the console"},
		},
		Note: "Super + Ctrl + Esc ends your session and closes everything " +
			"open in it. It is on this page so that you recognise it if you " +
			"land on it by accident, not as something to try now.",
	},
}

// ---------------------------------------------------------------------------
// The desktop: dock, launcher, panels.
// ---------------------------------------------------------------------------

var shellSteps = []Step{
	{
		Title: "The Dock",
		Body: "Move your pointer to the very bottom of the screen and leave " +
			"it there for a moment.\n\n" +
			"The dock slides up, and hides again when you move away. The " +
			"short pause before it appears is deliberate: without it, every " +
			"pointer movement that crossed the bottom of the screen would " +
			"summon it.\n\n" +
			"The dock is the taskbar. Applications you have pinned and " +
			"applications that are running share one strip, and running is a " +
			"small mark on the tile rather than a separate region of the " +
			"screen - because a program before and after you start it is the " +
			"same program.\n\n" +
			"The leftmost tile is always the launcher button. Rest on any " +
			"other tile for a moment and its open windows are pictured above " +
			"it.",
		BindingsTitle: "Clicking the dock",
		Bindings: []Binding{
			{"click", "start the application, or raise its window if one is open"},
			{"middle click", "start another copy, whether or not one is running"},
			{"Ctrl + click", "the same - start another copy"},
		},
		Note: "Raven's own applications open one window per launch, so a " +
			"second launch really is a second window. On a trackpad, use " +
			"Ctrl + click over the dock: a three-finger tap arrives as a " +
			"middle click, but the dock keeps its own meaning for it.",
	},
	{
		Title: "The Launcher",
		Body: "Press Super + Ctrl + Space, or click the leftmost tile on the " +
			"dock.\n\n" +
			"The launcher opens with the search field already focused, so " +
			"you can start typing immediately - two characters and Return is " +
			"usually enough. It searches your installed applications and " +
			"your files at once, and it will do arithmetic: type 1920/3 and " +
			"the answer is the first row.\n\n" +
			"Typed text that matches nothing still offers to run it as a " +
			"command, on the last row, where the highlight only lands when " +
			"there was nothing else to find.",
		BindingsTitle: "Inside the launcher",
		Bindings: []Binding{
			{"Super + Ctrl + Space", "open the launcher"},
			{"type", "search applications, files and arithmetic"},
			{"Up / Down", "move the selection (Ctrl + P / Ctrl + N also work)"},
			{"Return", "open the selected result"},
			{"Tab", "the selected application's other ways to start - and Pin"},
			{"Ctrl + Left / Right", "filter: everything, applications only, files only"},
			{"Ctrl + U", "empty the field without closing the launcher"},
			{"Ctrl + W", "delete the word behind the cursor"},
			{"Esc", "close without launching"},
		},
		Note: "Files need at least two characters before they are searched. " +
			"One letter matches most of an index of thousands, so f lists " +
			"Files and Firefox and no files at all.",
	},
	{
		Title: "Pinned Apps and Quick Settings",
		Body: "Two more panels sit on the same layer as the launcher.\n\n" +
			"Super + Ctrl + A shows the applications you have pinned. Things " +
			"get onto it from the launcher: press Tab on an application and " +
			"the last item of the menu is Pin.\n\n" +
			"Super + Ctrl + S opens quick settings - Wi-Fi, Bluetooth, " +
			"volume, brightness, animations, how long until the session " +
			"locks itself, suspend and power off. The last row opens the " +
			"full settings application, which is also Super + Ctrl + P.",
		BindingsTitle: "Panels",
		Bindings: []Binding{
			{"Super + Ctrl + A", "the pinned applications"},
			{"Super + Ctrl + S", "quick settings"},
			{"Super + Ctrl + P", "the settings application"},
			{"Super + Ctrl + I", "the software store"},
			{"Return / Delete", "on the pinned panel: open, or unpin"},
			{"Shift + arrows", "on the pinned panel: move a pin"},
		},
		Note: "Quick settings rows are stepped with the left and right " +
			"arrows and acted on with Return, so the whole panel works " +
			"without the pointer.",
		Launch: &Launch{Label: "Open Settings", Exec: "raven-settings"},
	},
}

// ---------------------------------------------------------------------------
// Windows and workspaces.
// ---------------------------------------------------------------------------

var windowSteps = []Step{
	{
		Title: "Windows and Tiles",
		Body: "Windows on Raven Linux are tiled: they share the screen " +
			"rather than stacking on top of each other, and a new window " +
			"takes a place in the layout instead of landing wherever the " +
			"last one did.\n\n" +
			"Open a couple of terminals and move them around. Super + Ctrl + " +
			"arrows moves the focused window between tiles; Super + Ctrl + " +
			"Return promotes it into the first tile, which is the big one.\n\n" +
			"Alt + Tab is the switcher, and it is not limited to this " +
			"workspace: it lists every window you have open anywhere, most " +
			"recently used first. Hold Alt, step with Tab, and let go.",
		BindingsTitle: "Windows",
		Bindings: []Binding{
			{"Super + Ctrl + E", "open a terminal (Super + Ctrl + T does the same)"},
			{"Super + Ctrl + Q", "close the focused window (or Super + Ctrl + X)"},
			{"Super + Ctrl + J / K", "focus the next / previous window"},
			{"Alt + Tab", "every window, most recent first (Shift for backwards)"},
			{"Super + Ctrl + arrows", "move the focused window between tiles"},
			{"Super + Ctrl + Return", "swap the focused window into the first tile"},
			{"Super + Ctrl + R", "resize the focused window with the arrows"},
		},
		Launch: &Launch{Label: "Open a Terminal", Exec: "raven-terminal"},
	},
	{
		Title: "Workspaces",
		Body: "A workspace is a whole screen of windows, and you can have " +
			"as many as you need. There is always one empty spare at the " +
			"end to move into, and a workspace you leave empty disappears " +
			"again.\n\n" +
			"Super + Ctrl + C opens the carousel: the workspace you are on " +
			"shrinks into a card in the middle, with the ones either side " +
			"of it as narrow, dimmed cards. Step along it and press Super + " +
			"Ctrl + C again to accept the one in the centre.\n\n" +
			"Every card is a real workspace, not a window from this one.",
		BindingsTitle: "Workspaces",
		Bindings: []Binding{
			{"Super + Ctrl + 1..9", "go to that workspace"},
			{"Super + Ctrl + Shift + 1..9", "send the focused window to that workspace"},
			{"Super + Ctrl + C", "open the carousel, then accept the centred workspace"},
			{"Tab / Shift + Tab", "with the carousel open, slide the row"},
			{"Super + wheel", "the workspace either side"},
			{"Super + middle click", "open or close the overview"},
			{"Super + Ctrl + Tab", "focus the next screen"},
			{"Super + Ctrl + Shift + Tab", "send the focused window to the next screen"},
		},
	},
	{
		Title: "Putting a Window Away",
		Body: "Minimising on Raven Linux is called putting a window away, " +
			"and it puts the window on the dock.\n\n" +
			"Super + Ctrl + M puts the focused window away. Super + Ctrl + " +
			"Shift + M brings back a strip in the middle of the screen " +
			"holding everything you have put away: step along it with the " +
			"arrows and press Return to bring one back into the workspace " +
			"you are on now.\n\n" +
			"The strip dismisses itself after about four seconds if you do " +
			"not choose anything. Esc closes it immediately.",
		BindingsTitle: "Put away, and bring back",
		Bindings: []Binding{
			{"Super + Ctrl + M", "put the focused window away to the dock"},
			{"Super + Ctrl + Shift + M", "show the put-away windows"},
			{"arrows / Tab", "step along the strip"},
			{"Return", "bring the highlighted window back"},
			{"Esc", "dismiss the strip"},
			{"Super + right click", "put away the window under the pointer"},
		},
		Note: "A window's own minimise button does the same thing, so " +
			"applications that draw their own titlebar agree with the rest " +
			"of the desktop.",
	},
}

// ---------------------------------------------------------------------------
// The device-specific chapter.
// ---------------------------------------------------------------------------

// laptopSteps cover trackpad gestures - laptop track only.
var laptopSteps = []Step{
	{
		Title: "Three Finger Gestures",
		Body: "Your trackpad drives the whole desktop with three fingers. " +
			"Put three fingers down and move them - the desktop follows " +
			"your fingers as they travel, rather than waiting for you to " +
			"finish, and settles wherever you lift.\n\n" +
			"Sideways is the workspaces. Up is the overview. Down puts the " +
			"window you are in away.\n\n" +
			"A flick decides by its direction, so a short quick swipe moves " +
			"one workspace; a slow one settles on whatever is nearest when " +
			"you let go. Try each of them now.",
		BindingsTitle: "Three fingers",
		Bindings: []Binding{
			{"three fingers sideways", "preview and switch between workspaces"},
			{"three fingers up", "open the overview"},
			{"three fingers down", "put the focused window away"},
			{"three fingers down, overview open", "close the overview again"},
			{"three-finger double tap", "show the put-away windows"},
			{"sideways, strip open", "highlight one of them"},
			{"three fingers up, strip open", "bring the highlighted window back"},
			{"Esc", "dismiss the strip"},
		},
		Note: "A swipe that starts vertical stays vertical for its whole " +
			"length, so putting a window away never jerks the workspaces " +
			"halfway across the screen. Two-finger scrolling is untouched - " +
			"it is not a gesture, and it goes to your applications as " +
			"normal. A three-finger tap is a middle click, except over the " +
			"dock, where Ctrl + click starts a second window instead.",
	},
	{
		Title: "The Same Thing Without the Trackpad",
		Body: "Every gesture has a key and a mouse equivalent, so nothing " +
			"is lost when you plug in a mouse or dock the laptop.\n\n" +
			"Super and the left button held down is the three fingers: the " +
			"pointer's travel is fed to the same recogniser, so the " +
			"workspaces follow the mouse sideways, the overview follows it " +
			"up, and a drag down puts a window away.\n\n" +
			"Pressed and let go without moving, it is the tap.",
		BindingsTitle: "Gestures on keys and mouse",
		Bindings: []Binding{
			{"three fingers sideways", "Super + Ctrl + C, or Super + wheel"},
			{"three fingers up", "Super + Ctrl + C, or Super + middle click"},
			{"three fingers down", "Super + Ctrl + M, or Super + right click"},
			{"three-finger double tap", "Super + Ctrl + Shift + M, or Super + click"},
			{"up, with the strip open", "Return, or wheel up"},
		},
	},
}

// mouseSteps cover the pointer and keyboard equivalents - desktop track only.
var mouseSteps = []Step{
	{
		Title: "Super and the Mouse",
		Body: "On a machine with a mouse, the gestures a trackpad would " +
			"give you live on Super and the buttons. The rule is simple: " +
			"Super and the left button held down is the three fingers.\n\n" +
			"Held and moved, it is a swipe outright - the workspaces follow " +
			"the mouse sideways and settle where you let go, the overview " +
			"follows it up, and a drag down puts away the window it started " +
			"on. The drag is fed raw pointer movement, so it keeps going " +
			"after the cursor has hit the edge of the screen.\n\n" +
			"Pressed and released without moving, the same button is a tap: " +
			"once for the strip of put-away windows, then again on a tile " +
			"to bring that window back.",
		BindingsTitle: "Pointer chords",
		Bindings: []Binding{
			{"Super + left drag", "sideways: workspaces; up: overview; down: put away"},
			{"Super + click", "the put-away strip; again on a tile brings it back"},
			{"Super + wheel", "the workspace either side"},
			{"Super + right click", "put the window under the pointer away"},
			{"Super + middle click", "open or close the overview"},
		},
		Note: "A plain click and a plain wheel always belong to the " +
			"application under the pointer, and so does every other " +
			"modifier - only exactly Super, with Ctrl, Alt and Shift up, is " +
			"the desktop's.",
	},
	{
		Title: "Lost the Pointer?",
		Body: "Shake the mouse - a few quick back-and-forths, in any " +
			"direction - and the pointer grows to three times its size, " +
			"holds for most of a second after you stop, and shrinks back.\n\n" +
			"The tip stays exactly where it was, so you can carry on using " +
			"it while it is large, and applications see nothing of it: the " +
			"pointer moves exactly as it did, only the drawing of it " +
			"changes.\n\n" +
			"Try it now, then move on.",
		Note: "An application that draws its own pointer, like a game that " +
			"has hidden it, is left alone.",
	},
}

// ---------------------------------------------------------------------------
// The closing chapters, shown to everyone.
// ---------------------------------------------------------------------------

var closingSteps = []Step{
	{
		Title: "Screenshots and Recording",
		Body: "Print is the one key besides the volume keys that works " +
			"without the Super layer, because that is where every other " +
			"desktop puts it.\n\n" +
			"It is settled before the panels, so it captures whatever is on " +
			"the screen - the launcher open, a menu down - and after the " +
			"lock, so a locked screen cannot be photographed through it.",
		BindingsTitle: "Capture",
		Bindings: []Binding{
			{"Print", "screenshot the whole screen"},
			{"Shift + Print", "screenshot a region you drag out"},
			{"Ctrl + Print", "screenshot the focused window"},
			{"Super + Print", "start or stop recording the screen"},
		},
		Note: "Locking the screen while a recording is running stops it, so " +
			"a recording shows the lock screen and never the desktop " +
			"behind it.",
	},
	{
		Title: "Notifications, Sound and Locking",
		Body: "Notifications appear in the corner and stack. Two chords " +
			"clear them, and quick settings has a Do not disturb row for " +
			"when you want none at all - under it, only critical " +
			"notifications appear and the rest wait. A Notifications row " +
			"in the same panel says how many are waiting and brings them " +
			"back as cards, so nothing that arrived while you were busy is " +
			"lost.\n\n" +
			"The volume keys work whatever is open - the launcher, quick " +
			"settings, even the lock screen - because they act on the " +
			"speakers rather than on the session. Each press shows a slider " +
			"at the bottom of the screen for a moment.\n\n" +
			"Super + L locks. Nothing else resolves while the session is " +
			"locked; every key goes to the lock screen, and the volume keys " +
			"are the only exception.",
		BindingsTitle: "Notifications and the session",
		Bindings: []Binding{
			{"Super + Ctrl + N", "dismiss the newest notification"},
			{"Super + Ctrl + Shift + N", "dismiss every notification"},
			{"volume keys", "raise, lower or mute the volume"},
			{"Super + L", "lock the session"},
		},
	},
}

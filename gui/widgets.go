package gui

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// The pieces every page is built from, in the desktop's own material: one
// glass surface, a hairline round it, the radius scale from the stylesheet,
// and the accent used for exactly two things.

// readableWidth is the widest a column of prose is allowed to get.
//
// Text does not become easier to read on a wider screen, it becomes harder:
// past roughly ninety characters the eye loses the start of the next line on
// the way back. So the window can be any size and the words stay in a column
// this wide, in the middle.
const readableWidth = 760

// readable centres an object in a column no wider than [readableWidth].
func readable(content fyne.CanvasObject) fyne.CanvasObject {
	return container.New(&readableLayout{}, content)
}

type readableLayout struct{}

func (l *readableLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	min := fyne.NewSize(0, 0)
	for _, o := range objects {
		min = min.Max(o.MinSize())
	}
	// The column is allowed to be narrower than its content is wide; what it
	// may not do is force the window to be as wide as readableWidth.
	min.Width = fyne.Min(min.Width, readableWidth)
	return min
}

func (l *readableLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	width := fyne.Min(size.Width, readableWidth)
	for _, o := range objects {
		o.Resize(fyne.NewSize(width, size.Height))
		o.Move(fyne.NewPos((size.Width-width)/2, 0))
	}
}

// centred stacks objects in the middle of whatever space it is given, in a
// column no wider than [readableWidth].
//
// The spacers are what distinguish this from simply centring: a centred
// container shrinks to the smallest size its contents will accept, which for
// a wrapping paragraph is a very narrow column of very short lines. Here the
// column keeps its full width and the free space goes above and below it.
func centred(objects ...fyne.CanvasObject) fyne.CanvasObject {
	column := container.NewVBox(append([]fyne.CanvasObject{layout.NewSpacer()},
		append(objects, layout.NewSpacer())...)...)
	return readable(column)
}

// surface is a panel of the glass: a filled, rounded rectangle with a
// hairline around it, holding padded content.
func surface(content fyne.CanvasObject, radius float32) *fyne.Container {
	back := canvas.NewRectangle(surfaceColor())
	back.CornerRadius = radius
	back.StrokeColor = hairline()
	back.StrokeWidth = 1
	return container.NewStack(back, container.NewPadded(content))
}

// card is a surface at the card radius, which is what a page's sections use.
func card(content fyne.CanvasObject) *fyne.Container {
	return surface(content, theme.Size(theme.SizeNameCardRadius))
}

// surfaceColor is the layer above the window: white over the dark glass,
// white over the light one. One material, differing only in how much light it
// lets through.
func surfaceColor() color.Color {
	if currentTheme.light {
		return white(0xb8)
	}
	return white(0x0d)
}

// hairline is the 1px edge that separates one layer from the next.
func hairline() color.Color {
	if currentTheme.light {
		return black(0x14)
	}
	return white(0x17)
}

// dim is body text that is present but secondary - the same face, carrying
// less light, which is the stylesheet's one mechanism for it.
func dim() color.Color {
	return fade(theme.Color(theme.ColorNameForeground), 0xa8)
}

// ---------------------------------------------------------------------------
// Keys.
// ---------------------------------------------------------------------------

// keycaps draws a chord as the keys it is made of.
//
// "Super + Ctrl + H" becomes three caps with a thin plus between them, which
// is how the chord looks on the keyboard the reader is being asked to find it
// on. A phrase with no plus in it - "three fingers sideways" - becomes one
// wide cap, and reads as the single motion it is.
func keycaps(chord string) fyne.CanvasObject {
	parts := strings.Split(chord, " + ")
	row := make([]fyne.CanvasObject, 0, len(parts)*2-1)
	for i, part := range parts {
		if i > 0 {
			plus := canvas.NewText("+", dim())
			plus.TextSize = theme.Size(theme.SizeNameCaptionText)
			row = append(row, container.NewCenter(plus))
		}
		row = append(row, keycap(part))
	}
	return container.NewHBox(row...)
}

// keycap is one key, drawn like one: a rounded tile with a catch-light along
// its top edge, the way the stylesheet raises a surface.
func keycap(text string) fyne.CanvasObject {
	label := canvas.NewText(text, theme.Color(theme.ColorNameForeground))
	label.TextSize = theme.Size(theme.SizeNameCaptionText) + 1
	label.TextStyle = fyne.TextStyle{Bold: true}

	size := fyne.MeasureText(label.Text, label.TextSize, label.TextStyle)

	back := canvas.NewRectangle(keycapColor())
	back.CornerRadius = 6
	back.StrokeColor = hairline()
	back.StrokeWidth = 1
	back.SetMinSize(fyne.NewSize(size.Width+16, size.Height+10))

	return container.NewStack(back, container.NewCenter(label))
}

func keycapColor() color.Color {
	if currentTheme.light {
		return white(0xf0)
	}
	return white(0x1a)
}

// ---------------------------------------------------------------------------
// Headings and rules.
// ---------------------------------------------------------------------------

// sectionTitle names a section inside a page: the name, and a short accent
// rule under it.
//
// The rule is the accent doing its "this is the thing you are looking at"
// job, at the smallest size it can do it - short and under the words rather
// than beside them, so that it reads as an underline and not as a stray mark
// against the edge of the card.
func sectionTitle(text string) fyne.CanvasObject {
	label := canvas.NewText(text, theme.Color(theme.ColorNameForeground))
	label.TextSize = theme.Size(theme.SizeNameSubHeadingText)
	label.TextStyle = fyne.TextStyle{Bold: true}

	rule := canvas.NewRectangle(theme.Color(theme.ColorNamePrimary))
	rule.CornerRadius = 1
	rule.SetMinSize(fyne.NewSize(28, 2))

	return container.NewVBox(label, container.NewHBox(rule))
}

// eyebrow is the small line above a title that says where in the tour this
// page is. Accent, small, and spaced out, so it reads as a label rather than
// as the first line of the text.
func eyebrow(text string) fyne.CanvasObject {
	label := canvas.NewText(strings.ToUpper(text), theme.Color(theme.ColorNamePrimary))
	label.TextSize = theme.Size(theme.SizeNameCaptionText) - 1
	label.TextStyle = fyne.TextStyle{Bold: true}
	return label
}

// title is a page's heading.
func title(text string) fyne.CanvasObject {
	label := canvas.NewText(text, theme.Color(theme.ColorNameForeground))
	label.TextSize = theme.Size(theme.SizeNameHeadingText)
	label.TextStyle = fyne.TextStyle{Bold: true}
	return label
}

// body is a paragraph: wrapped, at the reading size, in the foreground.
func body(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	return label
}

// secondary is a note: wrapped, smaller, carrying less light.
func secondary(text string) fyne.CanvasObject {
	label := widget.NewRichText(&widget.TextSegment{
		Text: text,
		Style: widget.RichTextStyle{
			ColorName: theme.ColorNameForeground,
			SizeName:  theme.SizeNameCaptionText,
			Inline:    false,
		},
	})
	label.Wrapping = fyne.TextWrapWord
	return label
}

// codeBlock is a few lines of shell, on their own inset surface.
//
// Commands are the one thing on a page that must be read character by
// character, so they get a monospaced face, a darker well to sit in, and no
// wrapping: a line break invented by the layout in the middle of a command is
// a command somebody will paste wrongly.
func codeBlock(commands string) fyne.CanvasObject {
	text := widget.NewLabelWithStyle(commands, fyne.TextAlignLeading,
		fyne.TextStyle{Monospace: true})

	back := canvas.NewRectangle(wellColor())
	back.CornerRadius = 8
	back.StrokeColor = hairline()
	back.StrokeWidth = 1

	return container.NewStack(back, container.NewPadded(text))
}

// wellColor is the layer *below* the surface it sits on - the one place the
// material goes the other way, because a command is set into a card rather
// than raised off it.
func wellColor() color.Color {
	if currentTheme.light {
		return black(0x0a)
	}
	return black(0x3d)
}

// badge is the small coloured word that labels a card - Note, Installed,
// Optional.
//
// Rich text rather than plain canvas text, because it carries the same inner
// padding as the paragraph beneath it: raw text does not, and a badge drawn
// without it starts a few pixels to the left of every line it introduces.
func badge(text string, name fyne.ThemeColorName) fyne.CanvasObject {
	return widget.NewRichText(&widget.TextSegment{
		Text: text,
		Style: widget.RichTextStyle{
			ColorName: name,
			SizeName:  theme.SizeNameCaptionText,
			TextStyle: fyne.TextStyle{Bold: true},
		},
	})
}

// ---------------------------------------------------------------------------
// Progress.
// ---------------------------------------------------------------------------

// stepProgress is the tour's place marker: one segment per page, filled for
// the pages behind you, accented for the one you are on.
//
// A bar with a percentage on it tells you how much is left as a number; a row
// of segments tells you at a glance, and shows that the pages are discrete
// things you can move between.
func stepProgress(index, total int) fyne.CanvasObject {
	segments := make([]fyne.CanvasObject, 0, total)
	for i := range total {
		bar := canvas.NewRectangle(progressColor(i, index))
		bar.CornerRadius = 2
		bar.SetMinSize(fyne.NewSize(6, 4))
		segments = append(segments, bar)
	}
	return container.New(&trackLayout{}, segments...)
}

func progressColor(segment, current int) color.Color {
	accent := theme.Color(theme.ColorNamePrimary)
	switch {
	case segment == current:
		return accent
	case segment < current:
		return fade(accent, 0x70)
	default:
		return fade(theme.Color(theme.ColorNameForeground), 0x24)
	}
}

// trackLayout lays segments in one row of equal width, with a hairline gap.
//
// A grid would do this, but a grid with sixteen columns on a narrow window
// gives each segment less than the gap, and the row turns into dots. This
// keeps the gap fixed and lets the segments take what is left.
type trackLayout struct{}

const segmentGap = 3

func (l *trackLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}
	min := objects[0].MinSize()
	return fyne.NewSize(min.Width*float32(len(objects))+segmentGap*float32(len(objects)-1), min.Height)
}

func (l *trackLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) == 0 {
		return
	}
	gaps := segmentGap * float32(len(objects)-1)
	width := (size.Width - gaps) / float32(len(objects))
	for i, o := range objects {
		o.Resize(fyne.NewSize(width, size.Height))
		o.Move(fyne.NewPos(float32(i)*(width+segmentGap), 0))
	}
}

// ---------------------------------------------------------------------------
// Things you can press that are not buttons.
// ---------------------------------------------------------------------------

// tapArea is an invisible sheet that makes whatever is under it pressable.
//
// Fyne's button holds a label and an icon and draws its own background, which
// is the wrong shape for a card with a heading, a paragraph and a picture in
// it. This is the other half of that: no drawing of its own, over a surface
// that does the drawing, reporting the press and the hover so the surface can
// light up under the pointer.
type tapArea struct {
	widget.BaseWidget
	onTap   func()
	onHover func(bool)
}

func newTapArea(onTap func(), onHover func(bool)) *tapArea {
	area := &tapArea{onTap: onTap, onHover: onHover}
	area.ExtendBaseWidget(area)
	return area
}

func (a *tapArea) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent))
}

func (a *tapArea) Tapped(*fyne.PointEvent) {
	if a.onTap != nil {
		a.onTap()
	}
}

func (a *tapArea) MouseIn(*desktop.MouseEvent) {
	if a.onHover != nil {
		a.onHover(true)
	}
}

func (a *tapArea) MouseMoved(*desktop.MouseEvent) {}

func (a *tapArea) MouseOut() {
	if a.onHover != nil {
		a.onHover(false)
	}
}

// Cursor turns the pointer into the one people expect over something that can
// be pressed, which is most of what tells them that it can be.
func (a *tapArea) Cursor() desktop.Cursor { return desktop.PointerCursor }

var (
	_ fyne.Tappable      = (*tapArea)(nil)
	_ desktop.Hoverable  = (*tapArea)(nil)
	_ desktop.Cursorable = (*tapArea)(nil)
)

// choice is a large pressable card: a picture, a name, and a line about what
// picking it means. The device question is asked with two of these.
func choice(icon fyne.Resource, name, detail string, onTap func()) fyne.CanvasObject {
	picture := canvas.NewImageFromResource(icon)
	picture.FillMode = canvas.ImageFillContain
	picture.SetMinSize(fyne.NewSize(56, 56))

	heading := canvas.NewText(name, theme.Color(theme.ColorNameForeground))
	heading.TextSize = theme.Size(theme.SizeNameSubHeadingText)
	heading.TextStyle = fyne.TextStyle{Bold: true}
	heading.Alignment = fyne.TextAlignCenter

	about := widget.NewLabel(detail)
	about.Wrapping = fyne.TextWrapWord
	about.Alignment = fyne.TextAlignCenter

	back := canvas.NewRectangle(surfaceColor())
	back.CornerRadius = theme.Size(theme.SizeNameCardRadius)
	back.StrokeColor = hairline()
	back.StrokeWidth = 1

	content := container.NewPadded(container.NewVBox(
		container.NewPadded(picture),
		heading,
		about,
	))

	// The card lights up and takes the accent for its edge while the pointer
	// is on it: the same two jobs the accent has everywhere else.
	hover := func(in bool) {
		if in {
			back.FillColor = theme.Color(theme.ColorNameHover)
			back.StrokeColor = theme.Color(theme.ColorNamePrimary)
		} else {
			back.FillColor = surfaceColor()
			back.StrokeColor = hairline()
		}
		back.Refresh()
	}

	return container.NewStack(back, content, newTapArea(onTap, hover))
}

// ---------------------------------------------------------------------------
// Application pictures.
// ---------------------------------------------------------------------------

// avatar is the picture beside an application's name: the icon the rest of
// the desktop draws it with, on a tile, or - for the few that name an icon
// this machine does not have - a tile with its initial.
//
// The tile is not decoration. An application icon is drawn for whatever is
// behind it on the machine it came from, and this desktop has both a light
// and a dark mode: raven-terminal's is a black bird with no background of its
// own, which is perfectly legible on the light page it was drawn for and
// invisible on dark glass. A pale tile under every icon is one rule that
// makes all of them legible in both modes, rather than a per-icon judgement
// that would be wrong for the next application somebody installs.
//
// A letter is the right fallback for the rest: unmistakably this application
// rather than a missing file, and legible without anybody having to draw
// anything.
func avatar(a App, size float32) fyne.CanvasObject {
	if icon := appIcon(a.Exec); icon != nil {
		picture := canvas.NewImageFromResource(icon)
		picture.FillMode = canvas.ImageFillContain
		// The picture fills the tile rather than sitting inside it, so an
		// icon that carries its own background covers the tile completely
		// and only a transparent one is backed by it. Otherwise every
		// square icon would be drawn inside a second, paler square.
		picture.SetMinSize(fyne.NewSize(size, size))
		return iconTile(picture, white(0xea), size)
	}

	initial := strings.ToUpper(a.Name[:1])
	label := canvas.NewText(initial, readableOn(theme.Color(theme.ColorNamePrimary)))
	label.TextSize = size * 0.45
	label.TextStyle = fyne.TextStyle{Bold: true}
	label.Alignment = fyne.TextAlignCenter

	return iconTile(label, theme.Color(theme.ColorNamePrimary), size)
}

// iconTile centres something on a rounded square of the given colour.
func iconTile(content fyne.CanvasObject, fill color.Color, size float32) fyne.CanvasObject {
	tile := canvas.NewRectangle(fill)
	tile.CornerRadius = size / 4.5
	tile.SetMinSize(fyne.NewSize(size, size))
	return container.NewStack(tile, container.NewCenter(content))
}

// logo is the machine's own bird, at the size asked for.
func logo(size float32) fyne.CanvasObject {
	resource := ravenLogo(currentTheme.light)
	if resource == nil {
		return container.NewWithoutLayout()
	}
	picture := canvas.NewImageFromResource(resource)
	picture.FillMode = canvas.ImageFillContain
	picture.SetMinSize(fyne.NewSize(size, size))
	return picture
}

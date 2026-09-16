package gui

// The two tracks, assembled from the chapters.
//
// Only one chapter differs between them. Everything a laptop is told about
// three fingers, a desktop is told about Super and the mouse buttons - the
// compositor feeds both to the same recogniser, so they are two ways of
// describing one thing rather than two features.
var (
	// LaptopSteps is the tutorial shown to people on a laptop.
	LaptopSteps = track(
		openingSteps,
		shellSteps,
		windowSteps,
		laptopSteps,
		closingSteps,
		applicationSteps,
		[]Step{oracleStep},
	)

	// DesktopSteps is the tutorial shown to people on a desktop.
	DesktopSteps = track(
		openingSteps,
		shellSteps,
		windowSteps,
		mouseSteps,
		closingSteps,
		applicationSteps,
		[]Step{oracleStep},
	)
)

// track joins step lists into one track without sharing backing storage, so
// the two tracks can never overwrite each other's steps.
func track(lists ...[]Step) []Step {
	var out []Step
	for _, list := range lists {
		out = append(out, list...)
	}
	return out
}

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.NewWithID("xyz.andy.fysionboard")
	w := a.NewWindow("FysionBoard")

	boardID := a.Preferences().StringWithFallback("currentID", "default")
	f := &fysion{app: a, win: w, id: boardID}
	w.SetContent(f.buildUI())
	w.Resize(fyne.NewSize(minColWidth*1.5, minColWidth*2.5))

	title := a.Preferences().String(boardID + ".name")
	if title != "" {
		content := a.Preferences().StringList(boardID + ".items")
		f.setBoard(title, content)
	}
	w.ShowAndRun()
}

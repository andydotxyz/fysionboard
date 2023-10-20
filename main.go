package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.NewWithID("xyz.andy.fysionboard")
	w := a.NewWindow("FysionBoard")

	f := &fysion{win: w}
	w.SetContent(f.buildUI())
	w.Resize(fyne.NewSize(minColWidth*1.5, minColWidth*2.5))
	w.ShowAndRun()
}

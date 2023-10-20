package main

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

const minColWidth = float32(240)

type layout struct {
	cols      int
	maxHeight float32
}

func (l *layout) Layout(objs []fyne.CanvasObject, s fyne.Size) {
	l.cols = int((s.Width + theme.Padding()) / (minColWidth + theme.Padding()))
	if l.cols < 1 {
		l.cols = 1
	}
	width := (s.Width - (theme.Padding() * float32(l.cols-1))) / float32(l.cols)
	l.maxHeight = 0

	tops := make([]float32, l.cols)
	col := 0
	for _, o := range objs {
		top := tops[col]

		ratio := float32(0.0)
		height := float32(0)
		switch t := o.(type) {
		case *canvas.Image:
			ratio = t.Aspect()
			height = width / ratio
		case *canvas.Rectangle:
			height = width * 0.25
		default:
			o.Resize(fyne.NewSize(width, 105)) // adjust for width then measure for wrapping
			height = o.MinSize().Height
		}
		o.Move(fyne.NewPos(float32(col)*width+float32(col)*theme.Padding(), top))
		o.Resize(fyne.NewSize(width, height))

		tops[col] = tops[col] + height + theme.Padding()
		l.maxHeight = float32(math.Max(float64(l.maxHeight), float64(top+height)))
		col++
		if col >= l.cols {
			col = 0
		}
	}
}

func (l *layout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(minColWidth, l.maxHeight) // prop open scroll container
}

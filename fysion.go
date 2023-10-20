package main

import (
	"errors"
	"image/color"
	"io"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type board struct {
	title string
	items []fyne.CanvasObject
}

type fysion struct {
	top *board

	title *widget.RichText
	body  *fyne.Container

	win fyne.Window
}

func (f *fysion) buildUI() fyne.CanvasObject {
	f.title = widget.NewRichTextFromMarkdown("# Untitled")
	f.body = container.New(&layout{})

	var add *widget.Button
	add = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		p := fyne.CurrentApp().Driver().AbsolutePositionForObject(add).AddXY(0, add.Size().Height)
		f.showAdd(p)
	})
	edit := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		input := widget.NewEntry()
		dialog.ShowForm("Set board title", "Update", "Cancel",
			[]*widget.FormItem{
				widget.NewFormItem("Vision title", input),
			},
			func(ok bool) {
				if !ok {
					return
				}
				f.title.ParseMarkdown("# " + input.Text)
			}, f.win)
	})
	top := container.NewBorder(nil, nil, nil, container.NewHBox(edit, add),
		container.NewCenter(f.title))

	return container.NewBorder(top, nil, nil, nil,
		container.NewVScroll(f.body))
}

func (f *fysion) showAdd(p fyne.Position) {
	m := fyne.NewMenu("Add",
		fyne.NewMenuItem("File", f.showAddFile),
		fyne.NewMenuItem("Text", f.showAddText),
		fyne.NewMenuItem("Color", f.showAddColor))
	widget.ShowPopUpMenuAtPosition(m, f.win.Canvas(), p)
}

func (f *fysion) showAddColor() {
	dialog.ShowColorPicker("Add a colour panel", "Pick a colour for your vision",
		func(c color.Color) {
			if c == nil {
				return
			}
			r := canvas.NewRectangle(c)
			r.SetMinSize(fyne.NewSize(150, 50))
			f.body.Add(r)
		}, f.win)
}

func (f *fysion) showAddFile() {
	o := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
		if r == nil {
			return
		}
		if err != nil {
			dialog.ShowError(err, f.win)
			return
		}

		defer r.Close()
		b, err := io.ReadAll(r)
		name := r.URI().Name()

		switch strings.ToLower(filepath.Ext(name)) {
		case ".png", ".jpeg", ".jpg":
			res := fyne.NewStaticResource(name, b)
			add := canvas.NewImageFromResource(res)
			add.FillMode = canvas.ImageFillContain
			add.ScaleMode = canvas.ImageScaleFastest
			f.body.Add(add)
		case ".txt":
			txt := widget.NewLabel(string(b))
			txt.Wrapping = fyne.TextWrapWord
			f.body.Add(txt)
		default:
			dialog.ShowError(errors.New("unsupported file type "+filepath.Ext(name)), f.win)
		}
	}, f.win)
	o.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpeg", ".jpg", ".txt"}))
	o.Show()
}

func (f *fysion) showAddText() {
	input := widget.NewMultiLineEntry()
	dialog.ShowForm("Add a text panel", "Add", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Vision text", input),
		},
		func(ok bool) {
			if !ok {
				return
			}
			r := widget.NewRichTextFromMarkdown("## " + input.Text)
			f.body.Add(r)
		}, f.win)
}

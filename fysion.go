package main

import (
	"errors"
	"image/color"
	"io"
	"path/filepath"
	"strconv"
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
	app fyne.App

	title *widget.RichText
	body  *fyne.Container

	win fyne.Window
}

func (f *fysion) addColor(c color.Color) {
	r := canvas.NewRectangle(c)
	r.SetMinSize(fyne.NewSize(150, 50))
	f.body.Add(r)
}

func (f *fysion) addFile(name string, r io.ReadCloser) {
	defer r.Close()
	b, _ := io.ReadAll(r)

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
}

func (f *fysion) addText(text string) {
	r := widget.NewRichTextFromMarkdown("## " + text)
	f.body.Add(r)
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
				f.app.Preferences().SetString(boardID+".name", input.Text)
			}, f.win)
	})
	top := container.NewBorder(nil, nil, nil, container.NewHBox(edit, add),
		container.NewCenter(f.title))

	return container.NewBorder(top, nil, nil, nil,
		container.NewVScroll(f.body))
}

func (f *fysion) setBoard(title string, items []string) {
	f.title.ParseMarkdown("# " + title)
	f.body.RemoveAll()
	for _, item := range items {
		u, err := storage.ParseURI(item)
		if err != nil {
			f.addText(item)
			continue
		}

		switch u.Scheme() {
		case "file":
			r, err := storage.Reader(u)
			if err != nil {
				fyne.LogError("Failed to open reader", err)
				continue
			}
			f.addFile(u.Name(), r)
		case "color":
			c := parseColor(item[8:]) // strip the color://
			f.addColor(c)
		default:
			f.addText(item)
		}
	}
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
			f.addColor(c)
			list := f.app.Preferences().StringList(boardID + ".items")
			list = append(list, "color://"+formatColor(c))
			f.app.Preferences().SetStringList(boardID+".items", list)
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

		f.addFile(r.URI().Name(), r)
		list := f.app.Preferences().StringList(boardID + ".items")
		list = append(list, r.URI().String())
		f.app.Preferences().SetStringList(boardID+".items", list)
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
			f.addText(input.Text)
			list := f.app.Preferences().StringList(boardID + ".items")
			list = append(list, input.Text)
			f.app.Preferences().SetStringList(boardID+".items", list)
		}, f.win)
}

func parseColor(s string) color.Color {
	parts := strings.Split(s, ",")
	if len(parts) != 4 {
		return color.Transparent
	}

	r, _ := strconv.Atoi(parts[0])
	g, _ := strconv.Atoi(parts[1])
	b, _ := strconv.Atoi(parts[2])
	a, _ := strconv.Atoi(parts[3])
	return color.NRGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)}
}

func formatColor(c color.Color) string {
	r, g, b, a := c.RGBA()

	parts := make([]string, 4)
	parts[0] = strconv.Itoa(int(r))
	parts[1] = strconv.Itoa(int(g))
	parts[2] = strconv.Itoa(int(b))
	parts[3] = strconv.Itoa(int(a))
	return strings.Join(parts, ",")
}

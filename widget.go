package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type widgetInfo struct {
	name   string
	create func() fyne.CanvasObject
	edit   func(fyne.CanvasObject) []*widget.FormItem
}

var widgets = map[string]widgetInfo{
	"*widget.Button": {
		name: "Button",
		create: func() fyne.CanvasObject {
			return widget.NewButton("Button", func() {})
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			b := obj.(*widget.Button)
			entry := widget.NewEntry()
			entry.SetText(b.Text)
			entry.OnChanged = func(text string) {
				b.SetText(text)
			}
			return []*widget.FormItem{
				widget.NewFormItem("Text", entry),
				widget.NewFormItem("Icon", widget.NewSelect(iconNames, func(selected string) {
					b.SetIcon(icons[selected])
				}))}
		},
	},
	"*widget.Hyperlink": {
		name: "Hyperlink",
		create: func() fyne.CanvasObject {
			fyneURL, _ := url.Parse("https://fyne.io")
			return widget.NewHyperlink("Link Text", fyneURL)
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			link := obj.(*widget.Hyperlink)
			title := widget.NewEntry()
			title.SetText(link.Text)
			title.OnChanged = func(text string) {
				link.SetText(text)
			}
			subtitle := widget.NewEntry()
			subtitle.SetText(link.URL.String())
			subtitle.OnChanged = func(text string) {
				fyneURL, _ := url.Parse(text)
				link.SetURL(fyneURL)
			}
			return []*widget.FormItem{
				widget.NewFormItem("Text", title),
				widget.NewFormItem("URL", subtitle)}
		},
	},
	"*widget.Card": {
		name: "Card",
		create: func() fyne.CanvasObject {
			return widget.NewCard("Title", "Subtitle", widget.NewLabel("Content here"))
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			c := obj.(*widget.Card)
			title := widget.NewEntry()
			title.SetText(c.Title)
			title.OnChanged = func(text string) {
				c.SetTitle(text)
			}
			subtitle := widget.NewEntry()
			subtitle.SetText(c.Subtitle)
			subtitle.OnChanged = func(text string) {
				c.SetSubTitle(text)
			}
			return []*widget.FormItem{
				widget.NewFormItem("Title", title),
				widget.NewFormItem("Subtitle", subtitle)}
		},
	},
	"*widget.Entry": {
		name: "Entry",
		create: func() fyne.CanvasObject {
			e := widget.NewEntry()
			e.SetPlaceHolder("Entry")
			return e
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			l := obj.(*widget.Entry)
			entry1 := widget.NewEntry()
			entry1.SetText(l.Text)
			entry1.OnChanged = func(text string) {
				l.SetText(text)
			}
			entry2 := widget.NewEntry()
			entry2.SetText(l.PlaceHolder)
			entry2.OnChanged = func(text string) {
				l.SetPlaceHolder(text)
			}
			return []*widget.FormItem{
				widget.NewFormItem("Text", entry1),
				widget.NewFormItem("PlaceHolder", entry2)}
		},
	},
	"*widget.Icon": {
		name: "Icon",
		create: func() fyne.CanvasObject {
			return widget.NewIcon(theme.HelpIcon())
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			i := obj.(*widget.Icon)
			return []*widget.FormItem{
				widget.NewFormItem("Icon", widget.NewSelect(iconNames, func(selected string) {
					i.SetResource(icons[selected])
				}))}
		},
	},
	"*widget.Label": {
		name: "Label",
		create: func() fyne.CanvasObject {
			return widget.NewLabel("Label")
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			l := obj.(*widget.Label)
			entry := widget.NewEntry()
			entry.SetText(l.Text)
			entry.OnChanged = func(text string) {
				l.SetText(text)
			}
			return []*widget.FormItem{
				widget.NewFormItem("Text", entry)}
		},
	},
	"*widget.Check": {
		name: "CheckBox",
		create: func() fyne.CanvasObject {
			return widget.NewCheck("Tick it or don't", func(b bool) {})
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			c := obj.(*widget.Check)
			title := widget.NewEntry()
			title.SetText(c.Text)
			title.OnChanged = func(text string) {
				c.Text = text
				c.Refresh()
			}
			isChecked := widget.NewCheck("", func(b bool) { c.SetChecked(b) })
			isChecked.SetChecked(c.Checked)
			return []*widget.FormItem{
				widget.NewFormItem("Title", title),
				widget.NewFormItem("isChecked", isChecked)}
		},
	},
	"*widget.RadioGroup": {
		name: "RadioGroup",
		create: func() fyne.CanvasObject {
			return widget.NewRadioGroup([]string{"Option 1", "Option 2"}, func(s string) {})
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			r := obj.(*widget.RadioGroup)
			initialOption := widget.NewRadioGroup(r.Options, func(s string) { r.SetSelected(s) })
			initialOption.SetSelected(r.Selected)
			entry := widget.NewMultiLineEntry()
			entry.SetText(strings.Join(r.Options, "\n"))
			entry.OnChanged = func(text string) {
				r.Options = strings.Split(text, "\n")
				r.Refresh()
				initialOption.Options = strings.Split(text, "\n")
				initialOption.Refresh()
			}
			return []*widget.FormItem{
				widget.NewFormItem("Options", entry),
				widget.NewFormItem("Initial Option", initialOption)}
		},
	},
	"*widget.Select": {
		name: "Select",
		create: func() fyne.CanvasObject {
			return widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {})
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			s := obj.(*widget.Select)
			initialOption := widget.NewSelect(s.Options, func(opt string) { s.SetSelected(opt) })
			initialOption.SetSelected(s.Selected)
			entry := widget.NewMultiLineEntry()
			entry.SetText(strings.Join(s.Options, "\n"))
			entry.OnChanged = func(text string) {
				s.Options = strings.Split(text, "\n")
				s.Refresh()
				initialOption.Options = strings.Split(text, "\n")
				initialOption.Refresh()
			}
			return []*widget.FormItem{
				widget.NewFormItem("Options", entry),
				widget.NewFormItem("Initial Option", initialOption)}
		},
	},
	"*widget.Accordion": {
		name: "Accordion",
		create: func() fyne.CanvasObject {
			return widget.NewAccordion(widget.NewAccordionItem("Item 1", widget.NewLabel("The content goes here")), widget.NewAccordionItem("Item 2", widget.NewLabel("Content part 2 goes here")))
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			// TODO: Need to add the properties
			// entry := widget.NewEntry()
			return []*widget.FormItem{}
		},
	},
	"*widget.List": {
		name: "List",
		create: func() fyne.CanvasObject {
			myList := []string{"Item 1", "Item 2", "Item 3", "Item 4"}
			// TODO: Need to make the list get adjusted to show the full list of items, currently it has only one item height apprx.
			return widget.NewList(func() int { return len(myList) }, func() fyne.CanvasObject {
				return container.New(layout.NewHBoxLayout(), widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
			}, func(id widget.ListItemID, item fyne.CanvasObject) {
				item.(*fyne.Container).Objects[1].(*widget.Label).SetText(myList[id])
			})
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			return []*widget.FormItem{}
		},
	},
	"*widget.Menu": {
		name: "Menu",
		create: func() fyne.CanvasObject {
			myMenu := fyne.NewMenu("Menu Name", fyne.NewMenuItem("Item 1", func() { fmt.Println("From Item 1") }), fyne.NewMenuItem("Item 2", func() { fmt.Println("From Item 2") }), fyne.NewMenuItem("Item 3", func() { fmt.Println("From Item 3") }))
			return widget.NewMenu(myMenu)
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			return []*widget.FormItem{}
		},
	},
	"*widget.Form": {
		name: "Form",
		create: func() fyne.CanvasObject {
			return widget.NewForm(widget.NewFormItem("Username", widget.NewEntry()), widget.NewFormItem("Password", widget.NewPasswordEntry()), widget.NewFormItem("", container.NewGridWithColumns(2, widget.NewButton("Submit", func() { fmt.Println("Form is submitted") }), widget.NewButton("Cancel", func() { fmt.Println("Form is Cancelled") }))))
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			return []*widget.FormItem{}
		},
	},
	"*widget.MultiLineEntry": {
		name: "Multi Line Entry",
		create: func() fyne.CanvasObject {
			mle := widget.NewMultiLineEntry()
			mle.SetPlaceHolder("Enter Some \nLong text \nHere")
			mle.Wrapping = fyne.TextWrapWord
			return mle
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			mle := obj.(*widget.Entry)
			placeholder := widget.NewMultiLineEntry()
			placeholder.Wrapping = fyne.TextWrapWord
			placeholder.SetText(mle.PlaceHolder)
			placeholder.OnChanged = func(s string) {
				mle.SetPlaceHolder(s)
			}
			value := widget.NewMultiLineEntry()
			value.Wrapping = fyne.TextWrapWord
			value.SetText(mle.Text)
			value.OnChanged = func(s string) {
				mle.SetText(s)
			}
			return []*widget.FormItem{
				widget.NewFormItem("Placeholder", placeholder),
				widget.NewFormItem("Value", value)}
		},
	},
	"*widget.PasswordEntry": {
		name: "Password Entry",
		create: func() fyne.CanvasObject {
			e := widget.NewPasswordEntry()
			e.SetPlaceHolder("Password Entry")
			return e
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			l := obj.(*widget.Entry)
			text := widget.NewPasswordEntry()
			text.SetText(l.Text)
			text.OnChanged = func(text string) {
				l.SetText(text)
			}
			placeholder := widget.NewEntry()
			placeholder.SetText(l.PlaceHolder)
			placeholder.OnChanged = func(text string) {
				l.SetPlaceHolder(text)
			}
			// hidePassword := widget.NewCheck("Hide Password", func(b bool) {})
			// hidePassword.SetChecked(l.Hidden)
			// hidePassword.OnChanged = func(b bool) {
			// 	l.Hidden = b
			// }
			return []*widget.FormItem{
				widget.NewFormItem("Text", text),
				// widget.NewFormItem("Hide password", placeholder),
				widget.NewFormItem("PlaceHolder", placeholder)}
		},
	},
	"*widget.ProgressBar": {
		name: "Progress Bar",
		create: func() fyne.CanvasObject {
			p := widget.NewProgressBar()
			p.SetValue(0.1)
			return p
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			p := obj.(*widget.ProgressBar)
			value := widget.NewEntry()
			value.SetText(fmt.Sprintf("%f", p.Value))
			value.OnChanged = func(s string) {
				if f, err := strconv.ParseFloat(s, 64); err == nil {
					p.SetValue(f)
				}
			}
			return []*widget.FormItem{
				widget.NewFormItem("Value", value)}
		},
	},
	"*widget.Separator": {
		// Separator's height(or width as you may call) and color come from the theme, so not sure if we can change the color and height here
		name: "Separator",
		create: func() fyne.CanvasObject {
			return widget.NewSeparator()
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			return []*widget.FormItem{}
		},
	},
	"*widget.Slider": {
		name: "Slider",
		create: func() fyne.CanvasObject {
			s := widget.NewSlider(0, 100)
			s.OnChanged = func(f float64) {
				fmt.Println("Slider changed to", f)
			}
			return s
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			slider := obj.(*widget.Slider)
			val := widget.NewEntry()
			val.SetText(fmt.Sprintf("%f", slider.Value))
			val.OnChanged = func(s string) {
				if f, err := strconv.ParseFloat(s, 64); err == nil {
					slider.SetValue(f)
				}
			}
			return []*widget.FormItem{
				widget.NewFormItem("Value", val)}
		},
	},
	"*widget.Table": {
		name: "Table",
		create: func() fyne.CanvasObject {
			return widget.NewTable(func() (int, int) { return 3, 3 }, func() fyne.CanvasObject {
				return widget.NewLabel("Cell 000, 000")
			}, func(id widget.TableCellID, cell fyne.CanvasObject) {
				label := cell.(*widget.Label)
				switch id.Col {
				case 0:
					label.SetText(fmt.Sprintf("%d", id.Row+1))
				case 1:
					label.SetText("A longer cell")
				default:
					label.SetText(fmt.Sprintf("Cell %d, %d", id.Row+1, id.Col+1))
				}
			})
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			return []*widget.FormItem{}
		},
	},
	"*widget.TextGrid": {
		name: "Text Grid",
		create: func() fyne.CanvasObject {
			to := widget.NewTextGrid()
			to.SetText("ABCD \nEFGH")
			return to
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			to := obj.(*widget.TextGrid)
			entry := widget.NewEntry()
			entry.SetText(to.Text())
			entry.OnChanged = func(s string) {
				to.SetText(s)
			}
			return []*widget.FormItem{
				widget.NewFormItem("Text", entry)}
		},
	},

	"*fyne.Container": {
		name: "Container",
		create: func() fyne.CanvasObject {
			return container.NewMax()
		},
		edit: func(obj fyne.CanvasObject) []*widget.FormItem {
			c := obj.(*fyne.Container)
			props := layoutProps[c]

			var items []*widget.FormItem
			var choose *widget.FormItem
			// TODO figure out how to work Border...
			choose = widget.NewFormItem("Layout", widget.NewSelect(layoutNames, func(l string) {
				lay := layouts[l]
				props["layout"] = l
				c.Layout = lay.create(props)
				c.Refresh()
				choose.Widget.Hide()

				edit := lay.edit
				items = []*widget.FormItem{choose}
				if edit != nil {
					items = append(items, edit(c, props)...)
				}

				editForm = widget.NewForm(items...)
				paletteList.Objects = []fyne.CanvasObject{editForm}
				choose.Widget.Show()
				paletteList.Refresh()
			}))
			choose.Widget.(*widget.Select).SetSelected(props["layout"])
			return items
		},
	},
}

// widgetNames is an array with the list of names of all the widgets
var widgetNames = extractWidgetNames()

// extractWidgetNames returns all the list of names of all the widgets from our data
func extractWidgetNames() []string {
	var widgetNamesFromData = make([]string, len(widgets))
	i := 0
	for k := range widgets {
		widgetNamesFromData[i] = k
		i++
	}
	return widgetNamesFromData
}

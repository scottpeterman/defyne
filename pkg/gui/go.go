package gui

import (
	"fmt"
	"go/format"
	"io"
	"reflect"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/fyne-io/defyne/internal/guidefs"
)

func ExportGo(obj fyne.CanvasObject, meta map[fyne.CanvasObject]map[string]string, w io.Writer) error {
	guidefs.InitOnce()

	packagesList := packagesRequired(obj)
	varList := varsRequired(obj, meta)
	code := exportCode(packagesList, varList, obj, meta)

	_, err := w.Write([]byte(code))
	return err
}

func ExportGoPreview(obj fyne.CanvasObject, meta map[fyne.CanvasObject]map[string]string, w io.Writer) error {
	guidefs.InitOnce()

	packagesList := packagesRequired(obj)
	packagesList = append(packagesList, "app")
	varList := varsRequired(obj, meta)
	code := exportCode(packagesList, varList, obj, meta)

	code += `
func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")
	gui := newGUI()
	myWindow.SetContent(gui.makeUI())
	myWindow.ShowAndRun()
}
`
	_, err := w.Write([]byte(code))

	return err
}

func exportCode(pkgs, vars []string, obj fyne.CanvasObject, meta map[fyne.CanvasObject]map[string]string) string {
	for i := 0; i < len(pkgs); i++ {
		if pkgs[i] != "net/url" {
			pkgs[i] = "fyne.io/fyne/v2/" + pkgs[i]
		}

		pkgs[i] = fmt.Sprintf(`	"%s"`, pkgs[i])
	}

	defs := make(map[string]string)

	_, clazz := getTypeOf(obj)
	main := guidefs.Widgets[clazz].Gostring(obj, meta, defs)
	setup := ""
	for k, v := range defs {
		setup += "g." + k + " = " + v + "\n"
	}

	code := fmt.Sprintf(`// auto-generated
// Code generated by GUI builder.

package main

import (
	"fyne.io/fyne/v2"
%s
)

type gui struct {
%s
}

func newGUI() *gui {
	return &gui{}
}

func (g *gui) makeUI() fyne.CanvasObject {
	%s

	return %s}
`,
		strings.Join(pkgs, "\n"),
		strings.Join(vars, "\n"),
		setup, main)

	formatted, err := format.Source([]byte(code))
	if err != nil {
		fyne.LogError("Failed to encode GUI code", err)
		return ""
	}
	return string(formatted)
}

func packagesRequired(obj fyne.CanvasObject) []string {
	if w, ok := obj.(fyne.Widget); ok {
		return packagesRequiredForWidget(w)
	}

	ret := []string{"container"}
	var objs []fyne.CanvasObject
	if c, ok := obj.(*fyne.Container); ok {
		objs = c.Objects
	} else if c, ok := obj.(*fyne.Container); ok {
		objs = c.Objects
	}
	for _, w := range objs {
		for _, p := range packagesRequired(w) {
			added := false
			for _, exists := range ret {
				if p == exists {
					added = true
					break
				}
			}
			if !added {
				ret = append(ret, p)
			}
		}
	}
	return ret
}

func packagesRequiredForWidget(w fyne.Widget) []string {
	name := reflect.TypeOf(w).String()
	if guidefs.Widgets[name].Packages != nil {
		return guidefs.Widgets[name].Packages(w)
	}

	return []string{"widget"}
}

func varsRequired(obj fyne.CanvasObject, props map[fyne.CanvasObject]map[string]string) []string {
	name := props[obj]["name"]
	if w, ok := obj.(fyne.Widget); ok {
		if name == "" {
			return []string{}
		}

		_, class := getTypeOf(w)
		return []string{name + " " + class}
	}

	var ret []string
	if c, ok := obj.(*fyne.Container); ok {
		if name != "" {
			ret = append(ret, name+" "+"*fyne.Container")
		}

		for _, w := range c.Objects {
			ret = append(ret, varsRequired(w, props)...)
		}
	}
	return ret
}

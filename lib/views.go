// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"
)

var (
	// The directory where all the views are being stored.
	ViewsDir = "views"
)

const (
	// The path to the layout file.
	layout = "application/layout"

	// The extension of views.
	viewsExt = "tpl"
)

// This struct holds all the data that can be passed to a view.
type ViewData struct {
	// The name of the javascript file to be used.
	JS string

	// The error message.
	Error string
}

// Returns the path to be used to open the view with the given name.
func view(name string) string {
	return path.Join(ViewsDir, name+"."+viewsExt)
}

// Render the view with the given name after evaluating the passed data. The
// rendered view will be written to the given writer.
func Render(res http.ResponseWriter, name string, data interface{}) {
	b, e := ioutil.ReadFile(view(layout))
	if e != nil {
		panic("Could not read layout file!")
	}
	t, e := template.New("l").Funcs(layoutHelpers(name, data)).Parse(string(b))
	if e != nil {
		panic("Could not parse layout file!")
	}
	t.Execute(res, data)
}

// Returns all the helpers used by the layout template. Right now only the
// "yield" helpers has been implemented.
func layoutHelpers(name string, data interface{}) template.FuncMap {
	return template.FuncMap{
		"yield": func() template.HTML {
			var buffer bytes.Buffer

			b, e := ioutil.ReadFile(view(name))
			if e != nil {
				r := fmt.Sprintf("Could not read: %v => %v", name, e)
				panic(r)
			}
			t := template.New(name).Funcs(viewHelpers())
			t, e = t.Parse(string(b))
			if e != nil {
				r := fmt.Sprintf("Could not parse: %v => %v", name, e)
				panic(r)
			}
			t.Execute(&buffer, data)
			return template.HTML(buffer.String())
		},
		"view": func() template.HTML {
			s := strings.SplitN(name, "/", 2)
			if len(s) == 2 {
				return template.HTML(s[0] + "_" + s[1])
			}
			return template.HTML("")
		},
	}
}

// Returns all the helpers available to any view. We have the following
// helpers: fmtDate and inc. The inc helper just increases the given integer
// value by one. The fmtDate helper executes the fmtDate function.
func viewHelpers() template.FuncMap {
	return template.FuncMap{
		"fmtDate": fmtDate,
		"inc": func(n int) int {
			return n + 1
		},
		"noescape": func(str string) template.HTML {
			return template.HTML(str)
		},
		"eqString": func(str1, str2 string) bool {
			return str1 == str2
		},
	}
}

// Returns a string with the given time formatted as expected by the view.
func fmtDate(t time.Time) string {
	return fmt.Sprintf("%02d/%02d/%04d", t.Day(), t.Month(), t.Year())
}

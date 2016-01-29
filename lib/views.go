// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

var (
	// ViewsDir is the directory where all the views are being stored.
	ViewsDir = "views"
)

const (
	// The path to the layout file.
	layout = "application/layout"

	// The extension of views.
	viewsExt = "tpl"
)

// ViewData holds all the data that can be passed to a view.
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
		log.Print("Could not read layout file")
		return
	}
	t, e := template.New("l").Funcs(layoutHelpers(name, data)).Parse(string(b))
	if e != nil {
		panic("Could not parse layout file")
	}
	if err := t.Execute(res, data); err != nil {
		log.Printf("Could not render template %v: %v", name, err)
	}
}

// Returns all the helpers used by the layout template. Right now only the
// "yield" helpers has been implemented.
func layoutHelpers(name string, data interface{}) template.FuncMap {
	return template.FuncMap{
		"yield": func() template.HTML {
			var buffer bytes.Buffer

			b, e := ioutil.ReadFile(view(name))
			if e != nil {
				log.Printf("Could not read: %v => %v", name, e)
				return template.HTML("")
			}
			t := template.New(name).Funcs(viewHelpers())
			t, e = t.Parse(string(b))
			if e != nil {
				log.Printf("Could not parse: %v => %v", name, e)
				return template.HTML("")
			}
			if err := t.Execute(&buffer, data); err != nil {
				log.Printf("Could not yield template %v: %v", name, err)
			}
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
	}
}

// Returns a string with the given time formatted as expected by the view.
func fmtDate(t time.Time) string {
	return fmt.Sprintf("%02d/%02d/%04d", t.Day(), t.Month(), t.Year())
}

// Copyright (C) 2014-2017 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package lib

import (
	"html/template"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestView(t *testing.T) {
	ViewsDir = "views"

	if view("path") != "views/path.tpl" {
		t.Fatal("Wrong path value expected: views/path.tpl")
	}
	if view("path/sub") != "views/path/sub.tpl" {
		t.Fatal("Wrong path value expected: views/path/sub.tpl")
	}
	if view("/path/sub") != "views/path/sub.tpl" {
		t.Fatal("Wrong path value expected: views/path/sub.tpl")
	}
}

func TestRender(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err.Error())
	}

	if filepath.Base(wd) == "lib" {
		wd = filepath.Dir(wd)
	}

	if err := os.Chdir(wd); err != nil {
		t.Fatalf("Could not change directory")
	}

	w := httptest.NewRecorder()
	Render(w, "topics/show", nil)

	body := w.Body.String()
	if !strings.HasPrefix(body, "<!DOCTYPE html>") {
		t.Fatalf("Not an HTML document?")
	}

	if !strings.Contains(body, "<form action=\"/topics\" method=\"POST\"") {
		t.Fatalf("There should be a POST /topics form")
	}
}

func TestHelpers(t *testing.T) {
	// fmtDate
	date := viewHelpers()["fmtDate"].(func(time.Time) string)
	d := time.Date(2014, time.March, 12, 12, 59, 12, 123, time.UTC)
	if dt := date(d); dt != "12/03/2014" {
		t.Fatalf("Wrong value %v, expected 12/03/2014", dt)
	}

	// inc
	inc := viewHelpers()["inc"].(func(int) int)
	if i := inc(0); i != 1 {
		t.Fatalf("Wrong value %v, expected 1", i)
	}

	// noescape
	noescape := viewHelpers()["noescape"].(func(string) template.HTML)
	if html := noescape("<a>"); html != "<a>" {
		t.Fatalf("Wrong value %v, expected <a>", html)
	}
}

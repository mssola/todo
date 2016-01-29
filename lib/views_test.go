// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
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
}

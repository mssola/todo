// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestView(t *testing.T) {
	viewsDir = "views"
	assert.Equal(t, view("path"), "views/path.tpl")
	assert.Equal(t, view("path/sub"), "views/path/sub.tpl")
	assert.Equal(t, view("/path/sub"), "views/path/sub.tpl")
}

func TestHelpers(t *testing.T) {
	// fmtDate
	date := viewHelpers()["fmtDate"].(func(time.Time) string)
	d := time.Date(2014, time.March, 12, 12, 59, 12, 123, time.UTC)
	assert.Equal(t, date(d), "12/03/2014")

	// inc
	inc := viewHelpers()["inc"].(func(int) int)
	assert.Equal(t, inc(0), 1)
}

// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {
	// Make sure that the InitSession function does something.
	assert.Nil(t, store)
	InitSession()
	assert.NotNil(t, store)

	// GetStore gets the proper store.
	req, _ := http.NewRequest("POST", "/", nil)
	s := GetStore(req)
	assert.Empty(t, s.Values)

	// GetCookie & SetCookie
	assert.Empty(t, GetCookie(req, "hello"))
	w := httptest.NewRecorder()
	SetCookie(w, req, "hello", "world")
	assert.Equal(t, GetCookie(req, "hello"), "world")

	// DeleteCookie
	SetCookie(w, req, "another", "anotherworld")
	SetCookie(w, req, "yetanother", "yetanotherworld")
	DeleteCookie(w, req, "another")
	s = GetStore(req)
	assert.Equal(t, len(s.Values), 2)
	assert.Equal(t, GetCookie(req, "hello"), "world")
	assert.Equal(t, GetCookie(req, "yetanother"), "yetanotherworld")
}

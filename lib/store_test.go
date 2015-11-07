// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSession(t *testing.T) {
	// Make sure that the InitSession function does something.
	if store != nil {
		t.Fatal("could not initalize test")
	}
	InitSession()
	if store == nil {
		t.Fatal("could not initalize test")
	}

	// GetStore gets the proper store.
	req, _ := http.NewRequest("POST", "/", nil)
	s := GetStore(req)
	if len(s.Values) != 0 {
		t.Fatalf("Expected to be empty, but: %v", s.Values)
	}

	// GetCookie & SetCookie
	if ck := GetCookie(req, "hello"); ck != nil {
		t.Fatalf("Expected to be empty: %v", ck)
	}
	w := httptest.NewRecorder()
	SetCookie(w, req, "hello", "world")
	if ck := GetCookie(req, "hello"); ck != "world" {
		t.Fatalf("Expected to be 'world': %v", ck)
	}

	// DeleteCookie
	SetCookie(w, req, "another", "anotherworld")
	SetCookie(w, req, "yetanother", "yetanotherworld")
	DeleteCookie(w, req, "another")
	s = GetStore(req)
	if len(s.Values) != 2 {
		t.Fatalf("Wrong number of cookies: %v", len(s.Values))
	}
	if ck := GetCookie(req, "hello"); ck != "world" {
		t.Fatalf("Expected to be 'world': %v", ck)
	}
	if ck := GetCookie(req, "yetanother"); ck != "yetanotherworld" {
		t.Fatalf("Expected to be 'yetanotherworld': %v", ck)
	}
}

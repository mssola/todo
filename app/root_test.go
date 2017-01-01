// Copyright (C) 2014-2017 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mssola/go-utils/security"
)

func TestCreateUserPage(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	w := httptest.NewRecorder()
	RootIndex(w, req)

	if w.Code != 200 {
		t.Fatalf("Got %v; Expected: %v", w.Code, 200)
	}
	if strings.Contains("<h1>Create user</h1>", w.Body.String()) {
		t.Fatalf("Body should've contained '<h1>Create user</h1>'")
	}
}

func TestLoginPage(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	password := security.PasswordSalt("1111")
	createUser("user", password)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	w := httptest.NewRecorder()
	RootIndex(w, req)

	if w.Code != 200 {
		t.Fatalf("Got %v, Expected: %v", w.Code, 200)
	}
	if strings.Contains("<h1>Login</h1>", w.Body.String()) {
		t.Fatalf("Body should've contained '<h1>Login</h1>'")
	}
}

func TestTopicsRedirect(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	password := security.PasswordSalt("1111")
	createUser("user", password)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	w := httptest.NewRecorder()
	login(w, req)
	RootIndex(w, req)

	if w.Code != 302 {
		t.Fatalf("Got %v, Expected: %v", w.Code, 302)
	}
	if w.HeaderMap["Location"][0] != "/topics" {
		t.Fatalf("Got %v, Expected: %v", w.HeaderMap["Location"][0], "/topics")
	}
}

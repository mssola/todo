// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/mssola/go-utils/security"
	"github.com/mssola/todo/lib"
)

func login(res http.ResponseWriter, req *http.Request) {
	var user User
	err := Db.SelectOne(&user, "select * from users")
	if err != nil {
		panic("There are no users...")
	}
	lib.SetCookie(res, req, "userId", user.ID)
}

func TestLogin(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	// This guy will be re-used throughout this test.
	param := make(url.Values)
	param["name"] = []string{"user"}
	param["password"] = []string{"1234"}

	// No users.
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.PostForm = param
	w := httptest.NewRecorder()
	Login(w, req)

	if w.Code != 302 {
		t.Fatalf("Got %v, Expected: %v", w.Code, 302)
	}
	if w.HeaderMap["Location"][0] != "/" {
		t.Fatalf("Got %v, Expected: %v", w.HeaderMap["Location"][0], "/")
	}
	if lib.GetCookie(req, "userId") != nil {
		t.Fatalf("Expected to be empty")
	}

	// Wrong password.
	password := security.PasswordSalt("1111")
	createUser("user", password)

	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.PostForm = param
	w = httptest.NewRecorder()
	Login(w, req)

	if w.Code != 302 {
		t.Fatalf("Got %v, Expected: %v", w.Code, 302)
	}
	if w.HeaderMap["Location"][0] != "/" {
		t.Fatalf("Got %v, Expected: %v", w.HeaderMap["Location"][0], "/")
	}
	if lib.GetCookie(req, "userId") != nil {
		t.Fatalf("Expected to be empty")
	}

	// Ok.
	req, err = http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	param["password"] = []string{"1111"}
	req.PostForm = param
	w = httptest.NewRecorder()
	Login(w, req)

	if w.Code != 302 {
		t.Fatalf("Got %v, Expected: %v", w.Code, 302)
	}
	if w.HeaderMap["Location"][0] != "/" {
		t.Fatalf("Got %v, Expected: %v", w.HeaderMap["Location"][0], "/")
	}
	if lib.GetCookie(req, "userId") == nil {
		t.Fatalf("Expected to be empty")
	}
	var user User
	err = Db.SelectOne(&user, "select * from users")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if lib.GetCookie(req, "userId") != user.ID {
		t.Fatalf("Wrong values")
	}
}

func TestLoginJson(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	// The body is nil.
	req, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Login(w, req)

	decoder := json.NewDecoder(w.Body)
	var r lib.Response
	err = decoder.Decode(&r)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if r.Error != "Failed!" {
		t.Fatalf("Got %v, Expected: %v", r.Error, "Failed!")
	}

	// The body is correct but there are no users.
	body := "{\"name\":\"mssola\",\"password\":\"1234\"}"
	req, err = http.NewRequest("POST", "/login", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	Login(w, req)

	decoder = json.NewDecoder(w.Body)
	err = decoder.Decode(&r)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if r.Error != "Failed!" {
		t.Fatalf("Got %v, Expected: %v", r.Error, "Failed!")
	}

	// Everything is fine.
	err = createUser("mssola", security.PasswordSalt("1234"))
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req, err = http.NewRequest("POST", "/login", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	Login(w1, req)

	decoder = json.NewDecoder(w1.Body)
	var u1, u2 User
	err = decoder.Decode(&u1)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	err = Db.SelectOne(&u2, "select * from users")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if u1.ID != u2.ID {
		t.Fatalf("Got %v, Expected: %v", u1.ID, u2)
	}

	// Malformed JSON
	body1 := "{\"password\":\"1234\""
	req, err = http.NewRequest("POST", "/login", strings.NewReader(body1))
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	Login(w2, req)

	decoder = json.NewDecoder(w2.Body)
	err = decoder.Decode(&r)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if r.Error != "Failed!" {
		t.Fatalf("Got %v, Expected: %v", r.Error, "Failed!")
	}
}

func TestLogout(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	// Create the user and loggin it in.
	password := security.PasswordSalt("1111")
	createUser("user", password)

	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	w := httptest.NewRecorder()
	login(w, req)

	// Check that the user has really been logged in.
	var user User
	err = Db.SelectOne(&user, "select * from users")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if ck := lib.GetCookie(req, "userId"); ck != user.ID {
		t.Fatalf("Got: %v; expected: %v", ck, user.ID)
	}

	// Logout
	Logout(w, req)
	if lib.GetCookie(req, "userId") != nil {
		t.Fatalf("Expected to be empty")
	}
}

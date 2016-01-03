// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mssola/go-utils/security"
)

func TestCreateUser(t *testing.T) {
	initTestDB()

	// There's nothing before.
	var u User
	err := Db.SelectOne(&u, "select * from users")
	if err == nil {
		t.Fatalf("Should be not nil")
	}
	if u.ID != "" {
		t.Fatalf("Expected to be empty")
	}

	// Now we create a user.
	err = createUser("u1", "1234")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	err = Db.SelectOne(&u, "select * from users")
	if u.ID == "" {
		t.Fatalf("Expected to not be empty")
	}
	if u.Name != "u1" {
		t.Fatalf("Got %v; Expected: %v", u.Name, "u1")
	}
	if u.PasswordHash == "" {
		t.Fatalf("Expected to not be empty")
	}

	// We cannot create another user.
	err = createUser("u2", "1234")
	if err == nil {
		t.Fatalf("Should be not nil")
	}
	if err.Error() != "too many users" {
		t.Fatalf("Got %v; Expected: %v", err.Error(), "too many users")
	}

	closeTestDB()
}

func TestMatchPassword(t *testing.T) {
	initTestDB()

	// User does not exist.
	u, err := matchPassword("u", "1234")
	if err == nil {
		t.Fatalf("Should be not nil")
	}

	// User exists but has a different password.
	password := security.PasswordSalt("1111")
	err = createUser("u", password)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	u, err = matchPassword("u", "1234")
	if err == nil {
		t.Fatalf("Should be not nil")
	}

	// User exists and has this password.
	u, err = matchPassword("u", "1111")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if u == "" {
		t.Fatalf("Expected to not be empty")
	}

	closeTestDB()
}

func TestUsersCreate(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	param := make(url.Values)
	param["name"] = []string{"user"}
	param["password"] = []string{"1234"}

	req, err := http.NewRequest("POST", "/users", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.PostForm = param
	w := httptest.NewRecorder()
	UsersCreate(w, req)

	if w.Code != 302 {
		t.Fatalf("Got %v; Expected: %v", w.Code, 302)
	}
	if w.HeaderMap["Location"][0] != "/" {
		t.Fatalf("Got %v; Expected: %v", w.HeaderMap["Location"][0], "/")
	}

	var user User
	err = Db.SelectOne(&user, "select * from users")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if user.ID == "" {
		t.Fatalf("Expected to not be empty")
	}
	if user.Name != "user" {
		t.Fatalf("Got %v; Expected: %v", user.Name, "user")
	}
	if user.PasswordHash == "" {
		t.Fatalf("Expected to not be empty")
	}
}

func TestUserCreateAlreadyExists(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	password := security.PasswordSalt("1234")
	createUser("user", password)

	param := make(url.Values)
	param["name"] = []string{"another"}
	param["password"] = []string{"1234"}

	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.PostForm = param
	w := httptest.NewRecorder()
	UsersCreate(w, req)

	if w.Code != 403 {
		t.Fatalf("Got %v; Expected: %v", w.Code, 403)
	}
	if w.HeaderMap["Location"][0] != "/" {
		t.Fatalf("Got %v; Expected: %v", w.HeaderMap["Location"][0], "/")
	}

	var user User
	err = Db.SelectOne(&user, "select * from users")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if user.ID == "" {
		t.Fatalf("Expected to not be empty")
	}
	if user.Name != "user" {
		t.Fatalf("Got %v; Expected: %v", user.Name, "user")
	}
	if user.PasswordHash == "" {
		t.Fatalf("Expected to not be empty")
	}
}

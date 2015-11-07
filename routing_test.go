// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mssola/todo/app"
	"github.com/mssola/todo/lib"
	"github.com/nu7hatch/gouuid"
)

// Initialize the database before running an unit test.
func initTestDB() {
	lib.InitSession()
	lib.ViewsDir = "../views"

	_ = os.Setenv("TODO_ENV", "test")
	app.InitDB()

	_ = app.Db.TruncateTables()
}

// Use this in the end of every unit test.
func closeTestDB() {
	_ = app.Db.TruncateTables()
	app.CloseDB()
}

func TestUserLogged(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Request was not successful: %v", err)
	}

	if val := userLogged(req, nil); val {
		t.Fatalf("Expected to be false: %v", val)
	}

	w := httptest.NewRecorder()
	lib.SetCookie(w, req, "userId", "1")

	if val := userLogged(req, nil); val {
		t.Fatalf("Expected to be false: %v", val)
	}

	uuid, _ := uuid.NewV4()
	u := &app.User{
		ID:           uuid.String(),
		Name:         "user",
		PasswordHash: "1234",
	}
	app.Db.Insert(u)
	var user app.User
	if err := app.Db.SelectOne(&user, "select * from users"); err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	w = httptest.NewRecorder()
	lib.SetCookie(w, req, "userId", user.ID)

	if val := userLogged(req, nil); !val {
		t.Fatalf("Expected to be true: %v", val)
	}
}

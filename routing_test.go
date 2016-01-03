// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/distribution/uuid"
	"github.com/mssola/todo/app"
	"github.com/mssola/todo/lib"
)

// Initialize the database before running an unit test.
func initTestDB() {
	lib.InitSession()
	lib.ViewsDir = "../views"

	app.InitDB()

	if err := app.Db.TruncateTables(); err != nil {
		log.Fatalf("Could not initialize DB: %v", err)
	}
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

	u := &app.User{ID: uuid.Generate().String(), Name: "user", PasswordHash: "1234"}
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

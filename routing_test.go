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
	"github.com/stretchr/testify/assert"
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
	assert.Nil(t, err)

	assert.False(t, userLogged(req, nil))

	w := httptest.NewRecorder()
	lib.SetCookie(w, req, "userId", "1")

	assert.False(t, userLogged(req, nil))

	uuid, _ := uuid.NewV4()
	u := &app.User{
		ID:           uuid.String(),
		Name:         "user",
		PasswordHash: "1234",
	}
	app.Db.Insert(u)
	var user app.User
	err = app.Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)

	w = httptest.NewRecorder()
	lib.SetCookie(w, req, "userId", user.ID)

	assert.True(t, userLogged(req, nil))
}

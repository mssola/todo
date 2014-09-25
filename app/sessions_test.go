// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mssola/go-utils/security"
	"github.com/mssola/todo/lib"
	"github.com/stretchr/testify/assert"
)

func login(res http.ResponseWriter, req *http.Request) {
	var user User
	err := Db.SelectOne(&user, "select * from users")
	if err != nil {
		panic("There are no users...")
	}
	lib.SetCookie(res, req, "userId", user.Id)
}

func TestLogin(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	// This guy will be re-used throughout this test.
	param := make(url.Values)
	param["name"] = []string{"user"}
	param["password"] = []string{"1234"}

	// No users.
	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	req.PostForm = param
	w := httptest.NewRecorder()
	Login(w, req)

	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/")
	assert.Empty(t, lib.GetCookie(req, "userId"))

	// Wrong password.
	password := security.PasswordSalt("1111")
	createUser("user", password)

	req, err = http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	req.PostForm = param
	w = httptest.NewRecorder()
	Login(w, req)

	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/")
	assert.Empty(t, lib.GetCookie(req, "userId"))

	// Ok.
	req, err = http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	param["password"] = []string{"1111"}
	req.PostForm = param
	w = httptest.NewRecorder()
	Login(w, req)

	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/")
	assert.NotEmpty(t, lib.GetCookie(req, "userId"))
	var user User
	err = Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)
	assert.Equal(t, lib.GetCookie(req, "userId"), user.Id)
}

func TestLogout(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	// Create the user and loggin it in.
	password := security.PasswordSalt("1111")
	createUser("user", password)

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	login(w, req)

	// Check that the user has really been logged in.
	var user User
	err = Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)
	assert.Equal(t, lib.GetCookie(req, "userId"), user.Id)

	// Logout
	Logout(w, req)
	assert.Empty(t, lib.GetCookie(req, "userId"))
}

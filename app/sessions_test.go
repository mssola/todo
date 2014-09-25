// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mssola/todo/app/lib"
	"github.com/mssola/todo/app/models"
	"github.com/stretchr/testify/assert"
)

func login(res http.ResponseWriter, req *http.Request) {
	var user models.User
	err := models.Db.SelectOne(&user, "select * from users")
	if err != nil {
		panic("There are no users...")
	}

	s := lib.GetStore(req)
	s.Values["userId"] = user.Id
	s.Save(req, res)
}

func TestLogin(t *testing.T) {
	InitTest()
	defer models.CloseDB()

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
	s := lib.GetStore(req)
	assert.Empty(t, s.Values["userId"])

	// Wrong password.
	createUser("user", "1111")
	req, err = http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	req.PostForm = param
	w = httptest.NewRecorder()
	Login(w, req)

	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/")
	s = lib.GetStore(req)
	assert.Empty(t, s.Values["userId"])

	// Ok.
	req, err = http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	param["password"] = []string{"1111"}
	req.PostForm = param
	w = httptest.NewRecorder()
	Login(w, req)

	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/")
	s = lib.GetStore(req)
	assert.NotEmpty(t, s.Values["userId"])
	var user models.User
	err = models.Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)
	assert.Equal(t, s.Values["userId"], user.Id)
}

func TestLogout(t *testing.T) {
	InitTest()
	defer models.CloseDB()

	// Create the user and loggin it in.
	createUser("user", "1111")
	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	login(w, req)

	// Check that the user has really been logged in.
	var user models.User
	err = models.Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)
	s := lib.GetStore(req)
	assert.Equal(t, s.Values["userId"], user.Id)

	// Logout
	Logout(w, req)
	s = lib.GetStore(req)
	assert.Empty(t, s.Values["userId"])
}

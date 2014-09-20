// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func login(res http.ResponseWriter, req *http.Request) {
	var user User
	err := Db.SelectOne(&user, "select * from users")
	if err != nil {
		panic("There are no users...")
	}

	s, err := store.Get(req, sessionName)
	if err != nil {
		panic("Could not get cookie store...")
	}
	s.Values["userId"] = user.Id
	s.Save(req, res)
}

func TestUserLogged(t *testing.T) {
	InitTest()
	defer CloseDB()

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	assert.False(t, UserLogged(req, nil))

	s, err := store.Get(req, sessionName)
	assert.Nil(t, err)
	s.Values["userId"] = "1"
	w := httptest.NewRecorder()
	s.Save(req, w)

	assert.False(t, UserLogged(req, nil))

	createUser("user", "1234")
	var user User
	err = Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)

	s, err = store.Get(req, sessionName)
	assert.Nil(t, err)
	s.Values["userId"] = user.Id
	w = httptest.NewRecorder()
	s.Save(req, w)

	assert.True(t, UserLogged(req, nil))
}

func TestLogin(t *testing.T) {
	InitTest()
	defer CloseDB()

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
	s, _ := store.Get(req, sessionName)
	assert.Empty(t, s.Values["userId"])

	// Wrong password.
	createUser("user", "1111")
	req, err = http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	req.PostForm = param
	w = httptest.NewRecorder()
	Login(w, req)

	assert.Equal(t, w.Code, 302)
	s, _ = store.Get(req, sessionName)
	assert.Empty(t, s.Values["userId"])

	// Ok.
	req, err = http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	param["password"] = []string{"1111"}
	req.PostForm = param
	w = httptest.NewRecorder()
	Login(w, req)

	assert.Equal(t, w.Code, 302)
	s, _ = store.Get(req, sessionName)
	assert.NotEmpty(t, s.Values["userId"])
	var user User
	err = Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)
	assert.Equal(t, s.Values["userId"], user.Id)
}

func TestLogout(t *testing.T) {
	InitTest()
	defer CloseDB()

	// Create the user and loggin it in.
	createUser("user", "1111")
	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	login(w, req)

	// Check that the user has really been logged in.
	var user User
	err = Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)
	s, _ := store.Get(req, sessionName)
	assert.Equal(t, s.Values["userId"], user.Id)

	// Logout
	Logout(w, req)
	s, _ = store.Get(req, sessionName)
	assert.Empty(t, s.Values["userId"])
}

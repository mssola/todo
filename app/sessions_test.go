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

func TestLoginJson(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	// The body is nil.
	req, err := http.NewRequest("POST", "/login", nil)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Login(w, req)

	decoder := json.NewDecoder(w.Body)
	var r lib.Response
	err = decoder.Decode(&r)
	assert.Nil(t, err)
	assert.Equal(t, r.Error, "Failed!")

	// The body is correct but there are no users.
	body := "{\"name\":\"mssola\",\"password\":\"1234\"}"
	req, err = http.NewRequest("POST", "/login", strings.NewReader(body))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	Login(w, req)

	decoder = json.NewDecoder(w.Body)
	err = decoder.Decode(&r)
	assert.Nil(t, err)
	assert.Equal(t, r.Error, "Failed!")

	// Everything is fine.
	err = createUser("mssola", security.PasswordSalt("1234"))
	assert.Nil(t, err)
	req, err = http.NewRequest("POST", "/login", strings.NewReader(body))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	Login(w1, req)

	decoder = json.NewDecoder(w1.Body)
	var u1, u2 User
	err = decoder.Decode(&u1)
	assert.Nil(t, err)
	err = Db.SelectOne(&u2, "select * from users")
	assert.Nil(t, err)
	assert.Equal(t, u1.Id, u2.Id)

	// Malformed JSON
	body1 := "{\"password\":\"1234\""
	req, err = http.NewRequest("POST", "/login", strings.NewReader(body1))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	Login(w2, req)

	decoder = json.NewDecoder(w2.Body)
	err = decoder.Decode(&r)
	assert.Nil(t, err)
	assert.Equal(t, r.Error, "Failed!")
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

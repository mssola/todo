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
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	InitTestDB()

	// There's nothing before.
	var u User
	err := Db.SelectOne(&u, "select * from users")
	assert.NotNil(t, err)
	assert.Empty(t, u.Id)

	// Now we create a user.
	err = createUser("u1", "1234")
	assert.Nil(t, err)
	err = Db.SelectOne(&u, "select * from users")
	assert.NotEmpty(t, u.Id)
	assert.Equal(t, u.Name, "u1")
	assert.NotEmpty(t, u.Password_hash)
	assert.NotEmpty(t, u.Created_at)

	// We cannot create another user.
	err = createUser("u2", "1234")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Too many users!")

	CloseTestDB()
}

func TestMatchPassword(t *testing.T) {
	InitTestDB()

	// User does not exist.
	u, err := matchPassword("u", "1234")
	assert.NotNil(t, err)

	// User exists but has a different password.
	password := security.PasswordSalt("1111")
	err = createUser("u", password)
	assert.Nil(t, err)
	u, err = matchPassword("u", "1234")
	assert.NotNil(t, err)

	// User exists and has this password.
	u, err = matchPassword("u", "1111")
	assert.Nil(t, err)
	assert.NotEmpty(t, u)

	CloseTestDB()
}

func TestUsersCreate(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	param := make(url.Values)
	param["name"] = []string{"user"}
	param["password"] = []string{"1234"}

	req, err := http.NewRequest("POST", "/users", nil)
	assert.Nil(t, err)
	req.PostForm = param
	w := httptest.NewRecorder()
	UsersCreate(w, req)

	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/")

	var user User
	err = Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)
	assert.NotEmpty(t, user.Id)
	assert.Equal(t, user.Name, "user")
	assert.NotEmpty(t, user.Password_hash)
	assert.NotEmpty(t, user.Created_at)
}

func TestUserCreateAlreadyExists(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	password := security.PasswordSalt("1234")
	createUser("user", password)

	param := make(url.Values)
	param["name"] = []string{"another"}
	param["password"] = []string{"1234"}

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	req.PostForm = param
	w := httptest.NewRecorder()
	UsersCreate(w, req)

	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/")

	var user User
	err = Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)
	assert.NotEmpty(t, user.Id)
	assert.Equal(t, user.Name, "user")
	assert.NotEmpty(t, user.Password_hash)
	assert.NotEmpty(t, user.Created_at)
}

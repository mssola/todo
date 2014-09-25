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
	"github.com/mssola/todo/app/models"
	"github.com/stretchr/testify/assert"
)

func TestUsersCreate(t *testing.T) {
	models.InitTestDB()
	defer models.CloseDB()

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

	var user models.User
	err = models.Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)
	assert.NotEmpty(t, user.Id)
	assert.Equal(t, user.Name, "user")
	assert.NotEmpty(t, user.Password_hash)
	assert.NotEmpty(t, user.Created_at)
}

func TestUserCreateAlreadyExists(t *testing.T) {
	models.InitTestDB()
	defer models.CloseDB()

	password := security.PasswordSalt("1234")
	models.CreateUser("user", password)

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

	var user models.User
	err = models.Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)
	assert.NotEmpty(t, user.Id)
	assert.Equal(t, user.Name, "user")
	assert.NotEmpty(t, user.Password_hash)
	assert.NotEmpty(t, user.Created_at)
}

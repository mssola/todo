// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mssola/go-utils/security"
	"github.com/mssola/todo/app/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserPage(t *testing.T) {
	models.InitTestDB()
	defer models.CloseTestDB()

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	RootIndex(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "<h1>Create user</h1>")
}

func TestLoginPage(t *testing.T) {
	models.InitTestDB()
	defer models.CloseTestDB()

	password := security.PasswordSalt("1111")
	models.CreateUser("user", password)

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	RootIndex(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Contains(t, w.Body.String(), "<h1>Login</h1>")
}

func TestTopicsRedirect(t *testing.T) {
	models.InitTestDB()
	defer models.CloseTestDB()

	password := security.PasswordSalt("1111")
	models.CreateUser("user", password)

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	login(w, req)
	RootIndex(w, req)

	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/topics")
}

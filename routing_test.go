// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mssola/todo/app/models"
	"github.com/mssola/todo/lib"
	"github.com/stretchr/testify/assert"
)

func TestUserLogged(t *testing.T) {
	models.InitTestDB()
	defer models.CloseTestDB()

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	assert.False(t, userLogged(req, nil))

	s := lib.GetStore(req)
	assert.Nil(t, err)
	s.Values["userId"] = "1"
	w := httptest.NewRecorder()
	s.Save(req, w)

	assert.False(t, userLogged(req, nil))

	models.CreateUser("user", "1234")
	var user models.User
	err = models.Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)

	s = lib.GetStore(req)
	assert.Nil(t, err)
	s.Values["userId"] = user.Id
	w = httptest.NewRecorder()
	s.Save(req, w)

	assert.True(t, userLogged(req, nil))
}

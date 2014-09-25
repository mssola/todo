// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mssola/todo/app"
	"github.com/mssola/todo/lib"
	"github.com/nu7hatch/gouuid"
	"github.com/stretchr/testify/assert"
)

func TestUserLogged(t *testing.T) {
	app.InitTestDB()
	defer app.CloseTestDB()

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	assert.False(t, userLogged(req, nil))

	s := lib.GetStore(req)
	assert.Nil(t, err)
	s.Values["userId"] = "1"
	w := httptest.NewRecorder()
	s.Save(req, w)

	assert.False(t, userLogged(req, nil))

	uuid, _ := uuid.NewV4()
	u := &app.User{
		Id:            uuid.String(),
		Name:          "user",
		Password_hash: "1234",
	}
	app.Db.Insert(u)
	var user app.User
	err = app.Db.SelectOne(&user, "select * from users")
	assert.Nil(t, err)

	s = lib.GetStore(req)
	assert.Nil(t, err)
	s.Values["userId"] = user.Id
	w = httptest.NewRecorder()
	s.Save(req, w)

	assert.True(t, userLogged(req, nil))
}

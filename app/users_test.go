// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/mssola/go-utils/security"
	"github.com/mssola/todo/app/lib"
	"github.com/mssola/todo/app/models"
	"github.com/nu7hatch/gouuid"
	"github.com/stretchr/testify/assert"
)

func InitTest() {
	InitSession()
	lib.ViewsDir = "../views"

	os.Setenv("TODO_ENV", "test")
	InitDB()
	tables := []string{"users", "topics"}
	for _, v := range tables {
		_, err := models.Db.Db.Exec(fmt.Sprintf("truncate table %v cascade", v))
		if err != nil {
			panic(fmt.Sprintf("Could not trucate table: %v\n", err))
		}
	}
}

func createUser(name, password string) {
	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	u := &models.User{
		Id:            uuid.String(),
		Name:          name,
		Password_hash: security.PasswordSalt(password),
		Created_at:    time.Now(),
	}
	models.Db.Insert(u)
}

func TestUsersCreate(t *testing.T) {
	InitTest()
	defer CloseDB()

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
	InitTest()
	defer CloseDB()
	createUser("user", "1234")

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

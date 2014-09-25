// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/mssola/go-utils/security"
	"github.com/nu7hatch/gouuid"
)

// There can only be one user in this application.
type User struct {
	Id            string
	Name          string
	Password_hash string
	Created_at    time.Time
}

// Create a new user with the given name and password. The given password is
// stored as-is. Therefore, it should've been encrypted before calling this
// function.
func createUser(name, password string) error {
	// Only one user is allowed in this application.
	count, err := Db.SelectInt("select count(*) from users")
	if err != nil || count > 0 {
		return errors.New("Too many users!")
	}

	// Create the user and redirect.
	uuid, _ := uuid.NewV4()
	u := &User{
		Id:            uuid.String(),
		Name:          name,
		Password_hash: password,
	}
	return Db.Insert(u)
}

// Match the user with the given name and password and its id.
func matchPassword(name, password string) (string, error) {
	var u User

	e := Db.SelectOne(&u, "select * from users where name=$1", name)
	if e != nil {
		return "", e
	}
	if !security.PasswordMatch(u.Password_hash, password) {
		return "", errors.New("Wrong password!")
	}
	return u.Id, nil
}

// Creates a user. It expects the "name" and the "password" form values to be
// present. Moreover, only one user is allowed in this application.
func UsersCreate(res http.ResponseWriter, req *http.Request) {
	password := security.PasswordSalt(req.FormValue("password"))

	createUser(req.FormValue("name"), password)
	http.Redirect(res, req, "/", http.StatusFound)
}

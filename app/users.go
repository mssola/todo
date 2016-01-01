// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/docker/distribution/uuid"
	"github.com/mssola/go-utils/security"
)

// User contains a user in this application.
type User struct {
	ID           string    `json:"token"`
	Name         string    `json:"-"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
}

// Create a new user with the given name and password. The given password is
// stored as-is. Therefore, it should've been encrypted before calling this
// function.
func createUser(name, password string) error {
	// Only one user is allowed in this application.
	if count, err := Db.SelectInt("select count(*) from users"); err != nil {
		return err
	} else if count > 0 {
		return errors.New("too many users")
	}

	// Create the user and redirect.
	id := uuid.Generate().String()
	return Db.Insert(&User{ID: id, Name: name, PasswordHash: password})
}

// Match the user with the given name and password and its id.
func matchPassword(name, password string) (string, error) {
	var u User

	e := Db.SelectOne(&u, "select * from users where name=$1", name)
	if e != nil {
		return "", e
	}
	if !security.PasswordMatch(u.PasswordHash, password) {
		return "", errors.New("Wrong password!")
	}
	return u.ID, nil
}

// UsersCreate responds to: POST /users. It expects the "name" and the
// "password" form values to be present. Moreover, only one user is
// allowed in this application.
func UsersCreate(res http.ResponseWriter, req *http.Request) {
	password := security.PasswordSalt(req.FormValue("password"))

	if err := createUser(req.FormValue("name"), password); err != nil {
		log.Printf("Could not create user: %v", err)
		http.Redirect(res, req, "/", http.StatusForbidden)
	} else {
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

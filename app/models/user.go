// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package models

import (
	"errors"
	"time"

	"github.com/mssola/go-utils/security"
	"github.com/nu7hatch/gouuid"
)

type User struct {
	Id            string
	Name          string
	Password_hash string
	Created_at    time.Time
}

// TODO: test
func CreateUser(name, password string) error {
	// Only one user is allowed in this application.
	count, err := Db.SelectInt("select count(*) from users")
	if err != nil || count > 0 {
		return errors.New("Too many users!")
	}

	// Create the user and redirect.
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	u := &User{
		Id:            uuid.String(),
		Name:          name,
		Password_hash: password,
		Created_at:    time.Now(),
	}
	return Db.Insert(u)
}

// TODO: test
func MatchPassword(name, password string) *User {
	u := &User{}

	e := Db.SelectOne(u, "select * from users where name=$1", name)
	if e != nil || !security.PasswordMatch(u.Password_hash, password) {
		return nil
	}
	return u
}

// TODO: test
func Logged(id string) bool {
	c, err := Db.SelectInt("select count(*) from users where id=$1", id)
	return err == nil && c == 1
}

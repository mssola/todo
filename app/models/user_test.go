// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package models

import (
	"testing"

	"github.com/mssola/go-utils/security"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	InitTestDB()
	defer CloseDB()

	// There's nothing before.
	var u User
	err := Db.SelectOne(&u, "select * from users")
	assert.NotNil(t, err)
	assert.Empty(t, u.Id)

	// Now we create a user.
	err = CreateUser("u1", "1234")
	assert.Nil(t, err)
	err = Db.SelectOne(&u, "select * from users")
	assert.NotEmpty(t, u.Id)
	assert.Equal(t, u.Name, "u1")
	assert.NotEmpty(t, u.Password_hash)
	assert.NotEmpty(t, u.Created_at)

	// We cannot create another user.
	err = CreateUser("u2", "1234")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Too many users!")
}

func TestMatchPassword(t *testing.T) {
	InitTestDB()
	defer CloseDB()

	// User does not exist.
	u, err := MatchPassword("u", "1234")
	assert.NotNil(t, err)

	// User exists but has a different password.
	password := security.PasswordSalt("1111")
	err = CreateUser("u", password)
	assert.Nil(t, err)
	u, err = MatchPassword("u", "1234")
	assert.NotNil(t, err)

	// User exists and has this password.
	u, err = MatchPassword("u", "1111")
	assert.Nil(t, err)
	assert.NotEmpty(t, u)
}

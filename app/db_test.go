// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	// Does not exist.
	assert.False(t, Exists("users", "1"))

	// Exists!
	createUser("u", "1234")
	var u User
	err := Db.SelectOne(&u, "select * from users")
	assert.Nil(t, err)
	assert.True(t, Exists("users", u.Id))
}

func TestCount(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	// Try to count a non-existing table.
	count := Count("doesnotexist")
	assert.Equal(t, count, 0)

	// Counting.
	count = Count("topics")
	assert.Equal(t, count, 0)
	createTopic("t1")
	createTopic("t2")
	count = Count("topics")
	assert.Equal(t, count, 2)
	count = Count("users")
	assert.Equal(t, count, 0)
}

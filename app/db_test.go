// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"os"
	"testing"

	"github.com/mssola/todo/lib"
	"github.com/stretchr/testify/assert"
)

// Initialize the database before running an unit test.
func initTestDB() {
	lib.InitSession()
	lib.ViewsDir = "../views"

	_ = os.Setenv("TODO_ENV", "test")
	InitDB()

	_ = Db.TruncateTables()
}

// Use this in the end of every unit test.
func closeTestDB() {
	_ = Db.TruncateTables()
	CloseDB()
}

func TestExists(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	// Does not exist.
	assert.False(t, Exists("users", "1"))

	// Exists!
	createUser("u", "1234")
	var u User
	err := Db.SelectOne(&u, "select * from users")
	assert.Nil(t, err)
	assert.True(t, Exists("users", u.ID))
}

func TestCount(t *testing.T) {
	initTestDB()
	defer closeTestDB()

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

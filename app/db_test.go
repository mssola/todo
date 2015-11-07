// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"os"
	"testing"

	"github.com/mssola/todo/lib"
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
	if Exists("users", "1") {
		t.Fatal("Expected to be false")
	}

	// Exists!
	createUser("u", "1234")
	var u User
	err := Db.SelectOne(&u, "select * from users")
	if err != nil {
		t.Fatalf("Should've executed correctly, but: %v", err)
	}
	if !Exists("users", u.ID) {
		t.Fatal("Expected to be true")
	}
}

func TestCount(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	// Try to count a non-existing table.
	count := Count("doesnotexist")
	if count != 0 {
		t.Fatalf("Wrong count: %v; expected: %v", count, 0)
	}

	// Counting.
	count = Count("topics")
	if count != 0 {
		t.Fatalf("Wrong count: %v; expected: %v", count, 0)
	}
	createTopic("t1")
	createTopic("t2")
	count = Count("topics")
	if count != 2 {
		t.Fatalf("Wrong count: %v; expected: %v", count, 2)
	}
	count = Count("users")
	if count != 0 {
		t.Fatalf("Wrong count: %v; expected: %v", count, 0)
	}
}

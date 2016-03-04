// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/mssola/todo/lib"
)

// Initialize the database before running an unit test.
func initTestDB() {
	log.SetOutput(ioutil.Discard)

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

func TestConfigURL(t *testing.T) {
	dbname := EnvOrElse("TODO_DB_NAME", "todo-dev")

	str := configURL()
	exp := fmt.Sprintf("user=postgres host=localhost port=5432 dbname=%s sslmode=disable", dbname)
	if str != exp {
		t.Fatalf("Got: '%s'; Expected: '%s'", str, exp)
	}

	exp = fmt.Sprintf("user=postgres host=localhost port=5432 dbname=%s sslmode=disable password=1234", dbname)
	os.Setenv("TODO_DB_PASSWORD", "1234")
	str = configURL()
	os.Setenv("TODO_DB_PASSWORD", "")
	if str != exp {
		t.Fatalf("Got: '%s'; Expected: '%s'", str, exp)
	}
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

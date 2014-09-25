// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package models

import (
	"github.com/coopernurse/gorp"
	"github.com/mssola/go-utils/db"
	"github.com/mssola/go-utils/misc"
	"github.com/mssola/go-utils/path"
)

// Global instance that holds a connection to the DB. It gets initialized after
// calling the InitDB function. You have to call CloseDB in order to close the
// connection.
var Db gorp.DbMap

// Initialize the global DB connection.
func InitDB() {
	c := db.Open(db.Options{
		Base:        path.FindRoot("todo", "."),
		Relative:    "/db/database.json",
		Environment: misc.EnvOrElse("TODO_ENV", "development"),
		DBMS:        "postgres",
		Heroku:      true,
	})
	Db = gorp.DbMap{Db: c, Dialect: gorp.PostgresDialect{}}
	Db.AddTableWithName(User{}, "users")
	Db.AddTableWithName(Topic{}, "topics")
}

// Close the global DB connection.
func CloseDB() {
	Db.Db.Close()
}

// TODO: test
func Count(name string) int64 {
	count, err := Db.SelectInt("select count(*) from " + name)
	if err != nil {
		return 0
	}
	return count
}

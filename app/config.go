// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"github.com/coopernurse/gorp"
	"github.com/mssola/go-utils/db"
	"github.com/mssola/go-utils/misc"
	"github.com/mssola/go-utils/path"
	"github.com/mssola/todo/app/models"
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
	Db.AddTableWithName(models.User{}, "users")
	Db.AddTableWithName(models.Topic{}, "topics")
}

// Close the global DB connection.
func CloseDB() {
	Db.Db.Close()
}

// This struct holds all the data that can be passed to a view.
type ViewData struct {
	// The id of the current user.
	Id string

	// Set to true if the current user is logged in.
	LoggedIn bool

	// Set to true if the views has to include Javascript.
	JS bool

	// Set to true if an error has happenned.
	Error bool
}

// This file is licensed under the MIT license.
// See the LICENSE file.

package models

import (
	"fmt"
	"os"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"github.com/mssola/go-utils/db"
	"github.com/mssola/go-utils/misc"
	"github.com/mssola/go-utils/path"
	"github.com/mssola/todo/lib"
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

// Initialize the database before running an unit test.
func InitTestDB() {
	lib.InitSession()
	lib.ViewsDir = "../views"

	os.Setenv("TODO_ENV", "test")
	InitDB()

	Db.TruncateTables()
}

// Use this in the end of every unit test.
func CloseTestDB() {
	Db.TruncateTables()
	CloseDB()
}

// Close the global DB connection.
func CloseDB() {
	Db.Db.Close()
}

// Returns true if there is a row in the given table that matches the given id.
// It returns false otherwise.
func Exists(name, id string) bool {
	q := fmt.Sprintf("select count(*) from %v where id=$1", name)
	c, err := Db.SelectInt(q, id)
	return err == nil && c == 1
}

// Count the number of rows for the given table. Returns a 0 on error. I know
// that this is not idiomatic, but it comes in handy in this case.
func Count(name string) int64 {
	count, err := Db.SelectInt("select count(*) from " + name)
	if err != nil {
		return 0
	}
	return count
}

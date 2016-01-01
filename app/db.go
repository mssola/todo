// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coopernurse/gorp"

	// Blank import because we are using postgresql
	_ "github.com/lib/pq"
)

// maxConnectionTries contains the number of connection attempts that this
// application is going to make before panic'ing.
//const maxConnectionTries = 10
const maxConnectionTries = 5

// Global instance that holds a connection to the DB. It gets initialized after
// calling the InitDB function. You have to call CloseDB in order to close the
// connection.
var Db gorp.DbMap

// EnvOrElse returns the value of the given environment variable. If this
// environment variable is not set, then it returns the provided alternative
// value.
func EnvOrElse(name, value string) string {
	if env := os.Getenv(name); env != "" {
		return env
	}
	return value
}

// configURL returns the string being used to connect with our PostgreSQL
// database.
func configURL() string {
	user := EnvOrElse("TODO_DB_USER", "postgres")
	dbname := EnvOrElse("TODO_DB_NAME", "todo-dev")
	password := EnvOrElse("TODO_DB_PASSWORD", "")
	host := EnvOrElse("TODO_DB_HOST", "localhost")
	sslmode := EnvOrElse("TODO_DB_SSLMODE", "disable")

	str := "user=%s host=%s port=5432 dbname=%s sslmode=%s"
	if password != "" {
		str += " password=%s"
		return fmt.Sprintf(str, user, host, dbname, sslmode, password)
	}
	return fmt.Sprintf(str, user, host, dbname, sslmode)
}

// establishConnection tries to establish a connection to the DB. It tries to
// do so until maxConnectionTries is reached, at which point it panics.
func establishConnection() *sql.DB {
	var err error

	str := configURL()
	log.Printf("Trying with: '%s'", str)
	d, err := sql.Open("postgres", str)

	for i := 0; i < maxConnectionTries; i++ {
		if err = d.Ping(); err == nil {
			log.Printf("postgres: connection established.")
			return d
		}
		if i < maxConnectionTries-1 {
			log.Printf("postgres: ping failed: %v", err)
			log.Printf("posgres: retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}
	log.Fatalf("postgres: could not establish connection with '%s'.", str)
	return nil
}

// InitDB initializes the global DB connection.
func InitDB() {
	d := establishConnection()

	Db = gorp.DbMap{Db: d, Dialect: gorp.PostgresDialect{}}
	Db.AddTableWithName(User{}, "users")
	Db.AddTableWithName(Topic{}, "topics")
}

// CloseDB close the global DB connection.
func CloseDB() {
	if err := Db.Db.Close(); err != nil {
		log.Printf("Could not close database: %v", err)
	}
}

// Exists returns true if there is a row in the given table that matches the
// given id. It returns false otherwise.
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

// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/mssola/go-utils/db"
	"github.com/mssola/go-utils/misc"
	"github.com/mssola/go-utils/path"
	"github.com/mssola/go-utils/security"
	"github.com/mssola/todo/app/models"
)

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

// Initialize the global DB connection.
func InitDB() {
	c := db.Open(db.Options{
		Base:        path.FindRoot("todo", "."),
		Relative:    "/db/database.json",
		Environment: misc.EnvOrElse("TODO_ENV", "development"),
		DBMS:        "postgres",
		Heroku:      true,
	})
	models.Db = gorp.DbMap{Db: c, Dialect: gorp.PostgresDialect{}}
	models.Db.AddTableWithName(models.User{}, "users")
	models.Db.AddTableWithName(models.Topic{}, "topics")
}

// Close the global DB connection.
func CloseDB() {
	models.Db.Db.Close()
}

// Global variable that holds the cookie store for this application. It gets
// initialized by calling the InitSession function.
var store *sessions.CookieStore

// The name of the session to be used for the safe cookies.
const sessionName = "todo"

// Initialize the global cookie store.
func InitSession() {
	store = sessions.NewCookieStore([]byte(security.NewAuthToken()))
	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 60 * 60 * 24 * 30 * 12, // A year.
	}
}

// A route matcher as expected by the mux package. It returns true (thus,
// accepting the route) if the current user is logged in, false otherwise.
func UserLogged(req *http.Request, rm *mux.RouteMatch) bool {
	s, _ := store.Get(req, sessionName)
	if id, ok := s.Values["userId"].(string); ok {
		return models.Logged(id)
	}
	return false
}

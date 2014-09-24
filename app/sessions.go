// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/mssola/go-utils/security"
	"github.com/mssola/todo/app/models"
)

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

// Returns true if the user with the given id exists, false otherwise.
func IsUserLogged(id interface{}) bool {
	if id == nil {
		return false
	}

	var u models.User
	e := Db.SelectOne(&u, "select * from users where id=$1", id.(string))
	return e == nil
}

// A route matcher as expected by the mux package. It returns true (thus,
// accepting the route) if the current user is logged in, false otherwise.
func UserLogged(req *http.Request, rm *mux.RouteMatch) bool {
	s, _ := store.Get(req, sessionName)
	return IsUserLogged(s.Values["userId"])
}

// Login a user. It expects the "name" and "password" form values. Regardless
// if it was successful or not, it will redirect the user to the root path.
func Login(res http.ResponseWriter, req *http.Request) {
	var u models.User

	// Check if the user exists and that the password is spot on.
	n, password := req.FormValue("name"), req.FormValue("password")
	e := Db.SelectOne(&u, "select * from users where name=$1", n)
	if e != nil || !security.PasswordMatch(u.Password_hash, password) {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	// It's ok to login this user.
	s, _ := store.Get(req, sessionName)
	s.Values["userId"] = u.Id
	s.Save(req, res)
	http.Redirect(res, req, "/", http.StatusFound)
}

// Logout the current user.
func Logout(res http.ResponseWriter, req *http.Request) {
	s, _ := store.Get(req, sessionName)
	delete(s.Values, "userId")
	s.Save(req, res)

	http.Redirect(res, req, "/", http.StatusFound)
}

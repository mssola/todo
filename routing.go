// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mssola/todo/app"
	"github.com/mssola/todo/app/models"
	"github.com/mssola/todo/lib"
)

// A route matcher as expected by the mux package. It returns true (thus,
// accepting the route) if the current user is logged in, false otherwise.
func userLogged(req *http.Request, rm *mux.RouteMatch) bool {
	id := lib.GetCookie(req, "userId")

	if value, ok := id.(string); ok {
		return models.Logged(value)
	}
	return false
}

// Handles the routing for this application. Returns a mux.Router with all our
// routes setup.
func route() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", app.RootIndex).Methods("GET")
	r.HandleFunc("/login", app.Login).Methods("POST")
	r.HandleFunc("/logout", app.Logout).Methods("POST").
		MatcherFunc(userLogged)
	r.HandleFunc("/users", app.UsersCreate).Methods("POST")
	r.HandleFunc("/topics", app.TopicsIndex).Methods("GET").
		MatcherFunc(userLogged)
	r.HandleFunc("/topics", app.TopicsCreate).Methods("POST").
		MatcherFunc(userLogged)

	return r
}

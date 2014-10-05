// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mssola/todo/app"
	"github.com/mssola/todo/lib"
)

// It returns true (thus, accepting the route) if the current user is
// logged in, false otherwise.
func userLogged(req *http.Request, rm *mux.RouteMatch) bool {
	var rid string

	if lib.JsonEncoding(req) {
		rid = req.URL.Query().Get("userId")
	} else if id, ok := lib.GetCookie(req, "userId").(string); ok {
		rid = id
	}
	return app.Exists("users", rid)
}

// Returns true if this request should not let JSON requests pass.
func private(req *http.Request, rm *mux.RouteMatch) bool {
	return !lib.JsonEncoding(req)
}

// Returns the "Not found" response.
func notFound(w http.ResponseWriter, req *http.Request) {
	if lib.JsonEncoding(req) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, lib.Response{Error: "Something wrong happened!"})
	} else {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "404 Not Found")
	}
}

// Handles the routing for this application. Returns a mux.Router with all our
// routes setup. This handles both the regular web page and the JSON API is.
// The endpoints open for the API are:
//
//  - sessions: login.
//  - topics: index, create, show, update, delete.
func route() *mux.Router {
	r := mux.NewRouter()

	// Let's get our own "Not Found" handler.
	r.NotFoundHandler = http.HandlerFunc(notFound)

	// The routing itself.
	r.HandleFunc("/", app.RootIndex).Methods("GET").
		MatcherFunc(private)
	r.HandleFunc("/login", app.Login).Methods("POST")
	r.HandleFunc("/logout", app.Logout).Methods("POST").
		MatcherFunc(userLogged).MatcherFunc(private)
	r.HandleFunc("/users", app.UsersCreate).Methods("POST").
		MatcherFunc(private)
	r.HandleFunc("/topics", app.TopicsIndex).Methods("GET").
		MatcherFunc(userLogged)
	r.HandleFunc("/topics", app.TopicsCreate).Methods("POST").
		MatcherFunc(userLogged)
	r.HandleFunc("/topics/{id}", app.TopicsShow).Methods("GET").
		MatcherFunc(userLogged)
	r.HandleFunc("/topics/{id}", app.TopicsUpdate).Methods("POST").
		MatcherFunc(userLogged)
	r.HandleFunc("/topics/{id}/delete", app.TopicsDestroy).Methods("POST").
		MatcherFunc(userLogged)

	return r
}

// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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

	if lib.JSONEncoding(req) {
		rid = req.URL.Query().Get("token")
	} else if id, ok := lib.GetCookie(req, "userId").(string); ok {
		rid = id
	}
	return app.Exists("users", rid)
}

// Returns true if this request should not let JSON requests pass.
func private(req *http.Request, rm *mux.RouteMatch) bool {
	return !lib.JSONEncoding(req)
}

// Returns the "Not found" response.
func notFound(w http.ResponseWriter, req *http.Request) {
	if lib.JSONEncoding(req) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, lib.Response{Error: "Something wrong happened!"})
	} else {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "404 Not Found")
	}
}

func license(w http.ResponseWriter, req *http.Request) {
	lib.Render(w, "application/license", &lib.ViewData{})
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
	r.HandleFunc("/license", license).Methods("GET")
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
		MatcherFunc(userLogged).MatcherFunc(private)
	r.HandleFunc("/topics/{id}", app.TopicsUpdateJSON).Methods("PATCH", "PUT").
		MatcherFunc(userLogged)
	r.HandleFunc("/topics/{id}/delete", app.TopicsDestroy).Methods("POST").
		MatcherFunc(userLogged).MatcherFunc(private)
	r.HandleFunc("/topics/{id}", app.TopicsDestroyJSON).Methods("DELETE").
		MatcherFunc(userLogged)

	return r
}

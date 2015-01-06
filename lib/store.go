// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mssola/go-utils/security"
)

// Global variable that holds the cookie store for this application. It gets
// initialized by calling the InitSession function.
var store *sessions.CookieStore

const (
	// The name of the session to be used for the safe cookies.
	sessionName = "todo"

	// Max-Age of a whole year.
	maxAge = 60 * 60 * 24 * 30 * 12
)

// Initialize the global cookie store.
func InitSession() {
	store = sessions.NewCookieStore([]byte(security.NewAuthToken()))
	store.Options = &sessions.Options{Path: "/", MaxAge: maxAge}
}

// Tries to get the cookie store for the given request. It panics if it fails.
func GetStore(req *http.Request) *sessions.Session {
	// Let's generate a new session. Moreover, Gorilla's documentation
	// says that this method *never* fails on creating a new session, so
	// it's safe to ignore the given error.
	s, _ := store.Get(req, sessionName)
	return s
}

// Returns the value for the specified cookie. If the cookie does not exist,
// then an empty interface{} gets returned.
func GetCookie(req *http.Request, name string) interface{} {
	s := GetStore(req)
	return s.Values[name]
}

// Sets the given value to the specified cookie.
func SetCookie(res http.ResponseWriter, req *http.Request, key, value string) {
	s := GetStore(req)
	s.Values[key] = value
	s.Save(req, res)
}

// Deletes the specified cookie.
func DeleteCookie(res http.ResponseWriter, req *http.Request, key string) {
	s := GetStore(req)
	delete(s.Values, key)
	s.Save(req, res)
}

// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package lib

import (
	"log"
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

	// Only the current session.
	maxAge = 0
)

// InitSession initializes the global cookie store.
func InitSession() {
	store = sessions.NewCookieStore([]byte(security.NewAuthToken()))
	store.Options = &sessions.Options{Path: "/", MaxAge: maxAge}
}

// GetStore tries to get the cookie store for the given request. It panics if
// it fails.
func GetStore(req *http.Request) *sessions.Session {
	// Let's generate a new session. Moreover, Gorilla's documentation
	// says that this method *never* fails on creating a new session, so
	// it's safe to ignore the given error.
	s, _ := store.Get(req, sessionName)
	return s
}

// GetCookie returns the value for the specified cookie. If the cookie does
// not exist, then an empty interface{} gets returned.
func GetCookie(req *http.Request, name string) interface{} {
	s := GetStore(req)
	return s.Values[name]
}

// SetCookie sets the given value to the specified cookie.
func SetCookie(res http.ResponseWriter, req *http.Request, key, value string) {
	s := GetStore(req)
	s.Values[key] = value
	if err := s.Save(req, res); err != nil {
		log.Printf("Could not save cookie: %v", err)
	}
}

// DeleteCookie deletes the specified cookie.
func DeleteCookie(res http.ResponseWriter, req *http.Request, key string) {
	s := GetStore(req)
	delete(s.Values, key)
	if err := s.Save(req, res); err != nil {
		log.Printf("Could not save cookie: %v", err)
	}
}

// Copyright (C) 2013 Miquel Sabaté Solà <mikisabate@gmail.com>
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

func GetStore(req *http.Request) *sessions.Session {
	s, err := store.Get(req, sessionName)
	if err != nil {
		panic("Could not get the cookie store!")
	}
	return s
}

func GetCookie(req *http.Request, name string) interface{} {
	s := GetStore(req)
	return s.Values[name]
}

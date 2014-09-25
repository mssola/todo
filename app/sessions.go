// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"

	"github.com/mssola/todo/lib"
)

// Login a user. It expects the "name" and "password" form values. Regardless
// if it was successful or not, it will redirect the user to the root path.
func Login(res http.ResponseWriter, req *http.Request) {
	// Check if the user exists and that the password is spot on.
	n, password := req.FormValue("name"), req.FormValue("password")
	id, err := matchPassword(n, password)
	if err != nil {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	// It's ok to login this user.
	s := lib.GetStore(req)
	s.Values["userId"] = id
	s.Save(req, res)
	http.Redirect(res, req, "/", http.StatusFound)
}

// Logout the current user.
func Logout(res http.ResponseWriter, req *http.Request) {
	s := lib.GetStore(req)
	delete(s.Values, "userId")
	s.Save(req, res)

	http.Redirect(res, req, "/", http.StatusFound)
}

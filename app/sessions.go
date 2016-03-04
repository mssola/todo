// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mssola/todo/lib"
)

// Returns the name and the password parameters as given by the request. This
// method abstracts away the origin of these values.
func getNamePassword(req *http.Request) (string, string) {
	if lib.JSONEncoding(req) {
		if req.Body == nil {
			return "", ""
		}

		decoder := json.NewDecoder(req.Body)

		var t struct{ Name, Password string }
		err := decoder.Decode(&t)
		if err != nil {
			return "", ""
		}
		return t.Name, t.Password
	}
	return req.FormValue("name"), req.FormValue("password")
}

// Login a user. It expects the "name" and "password" form values. Regardless
// if it was successful or not, it will redirect the user to the root path.
func Login(res http.ResponseWriter, req *http.Request) {
	// Check if the user exists and that the password is spot on.
	n, password := getNamePassword(req)
	id, err := matchPassword(n, password)
	if lib.CheckError(res, req, err) {
		return
	}

	// It's ok to login this user.
	if lib.JSONEncoding(req) {
		b, _ := json.Marshal(User{ID: id})
		fmt.Fprint(res, string(b))
	} else {
		lib.SetCookie(res, req, "userId", id)
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

// Logout the current user.
func Logout(res http.ResponseWriter, req *http.Request) {
	lib.DeleteCookie(res, req, "userId")
	lib.DeleteCookie(res, req, "topic")
	http.Redirect(res, req, "/", http.StatusFound)
}

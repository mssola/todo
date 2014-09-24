// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
)

// Renders the root page. It has three different options:
//
//  1. If there's no user, it renders the "Create user" page.
//  2. If the current user is not logged in, it render the "Login" page.
//  3. If the current user is logged in, then it redirects the user to the
//     /topics page.
func RootIndex(res http.ResponseWriter, req *http.Request) {
	s, _ := store.Get(req, sessionName)
	id := s.Values["userId"]

	if id == nil {
		count, err := Db.SelectInt("select count(*) from users")
		if err == nil && count == 0 {
			render(res, "users/new", &Options{})
		} else {
			render(res, "application/login", &Options{})
		}
	} else {
		http.Redirect(res, req, "/topics", http.StatusFound)
	}
}

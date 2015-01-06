// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"

	"github.com/mssola/todo/lib"
)

// Renders the root page. It has three different options:
//
//  1. If there's no user, it renders the "Create user" page.
//  2. If the current user is not logged in, it render the "Login" page.
//  3. If the current user is logged in, then it redirects the user to the
//     /topics page.
func RootIndex(res http.ResponseWriter, req *http.Request) {
	id := lib.GetCookie(req, "userId")

	if id == nil {
		count := Count("users")
		if count == 0 {
			lib.Render(res, "users/new", &lib.ViewData{})
		} else {
			lib.Render(res, "application/login", &lib.ViewData{})
		}
	} else {
		http.Redirect(res, req, "/topics", http.StatusFound)
	}
}

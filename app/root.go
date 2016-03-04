// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package app

import (
	"net/http"

	"github.com/mssola/todo/lib"
)

// RootIndex renders the root page. It has three different options:
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

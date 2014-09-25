// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mssola/todo/app/lib"
	"github.com/mssola/todo/app/models"
)

// A route matcher as expected by the mux package. It returns true (thus,
// accepting the route) if the current user is logged in, false otherwise.
func UserLogged(req *http.Request, rm *mux.RouteMatch) bool {
	id := lib.GetCookie(req, "userId")

	if value, ok := id.(string); ok {
		return models.Logged(value)
	}
	return false
}

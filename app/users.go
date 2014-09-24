// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"

	"github.com/mssola/go-utils/security"
	"github.com/mssola/todo/app/models"
)

// Creates a user. It expects the "name" and the "password" form values to be
// present. Moreover, only one user is allowed in this application.
func UsersCreate(res http.ResponseWriter, req *http.Request) {
	password := security.PasswordSalt(req.FormValue("password"))

	models.CreateUser(req.FormValue("name"), password)
	http.Redirect(res, req, "/", http.StatusFound)
}

// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
	"time"

	"github.com/mssola/go-utils/security"
	"github.com/nu7hatch/gouuid"
)

// Creates a user. It expects the "name" and the "password" form values to be
// present. Moreover, only one user is allowed in this application.
func UsersCreate(res http.ResponseWriter, req *http.Request) {
	count, err := Db.SelectInt("select count(*) from users")
	if err != nil || count > 0 {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	// Create the user and redirect.
	uuid, err := uuid.NewV4()
	if err != nil {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}
	u := &User{
		Id:            uuid.String(),
		Name:          req.FormValue("name"),
		Password_hash: security.PasswordSalt(req.FormValue("password")),
		Created_at:    time.Now(),
	}
	Db.Insert(u)
	http.Redirect(res, req, "/", http.StatusFound)
}

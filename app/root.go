// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
)

// TODO
func RootIndex(res http.ResponseWriter, req *http.Request) {
	s, _ := store.Get(req, sessionName)
	id := s.Values["userId"]

	if id == nil {
		o := &Options{}
		count, err := Db.SelectInt("select count(*) from users")
		if err == nil && count == 0 {
			render(res, "users/new", o)
		} else {
			render(res, "application/login", o)
		}
	} else {
		http.Redirect(res, req, "/topics", http.StatusFound)
	}
}

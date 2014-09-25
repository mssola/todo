// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"

	"github.com/mssola/todo/app/models"
	"github.com/mssola/todo/lib"
)

func TopicsIndex(res http.ResponseWriter, req *http.Request) {
	o := &lib.ViewData{}
	lib.Render(res, "topics/index", o)
}

func TopicsCreate(res http.ResponseWriter, req *http.Request) {
	models.CreateTopic(req.FormValue("name"))
	http.Redirect(res, req, "/topics", http.StatusFound)
}

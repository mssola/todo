// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mssola/todo/app/models"
	"github.com/mssola/todo/lib"
)

func renderShow(res http.ResponseWriter, topic *models.Topic) {
	o := &lib.ViewData{}
	lib.Render(res, "topics/show", o)
}

func TopicsIndex(res http.ResponseWriter, req *http.Request) {
	var t models.Topic

	models.Db.SelectOne(&t, "select * from topics order by name limit 1")
	renderShow(res, &t)
}

func TopicsCreate(res http.ResponseWriter, req *http.Request) {
	models.CreateTopic(req.FormValue("name"))
	http.Redirect(res, req, "/topics", http.StatusFound)
}

func TopicsShow(res http.ResponseWriter, req *http.Request) {
	var t models.Topic

	p := mux.Vars(req)
	models.Db.SelectOne(&t, "select * from topics where id=$1", p["id"])
	renderShow(res, &t)
}

func TopicsUpdate(res http.ResponseWriter, req *http.Request) {
	var t models.Topic

	// Get the original.
	p := mux.Vars(req)
	models.Db.SelectOne(&t, "select * from topics where id=$1", p["id"])

	// We can either rename, or change the contents, but not both things at the
	// same time.
	name := req.FormValue("name")
	if name != "" {
		t.name = name
	} else {
		t.contents = req.FormValue("contents")
	}
	models.Db.update(&t)
	http.Redirect(res, req, "/topics", http.StatusFound)
}

func TopicsDestroy(res http.ResponseWriter, req *http.Request) {
	p := mux.Vars(req)
	models.Db.Exec("delete from topics where id=$1", p["id"])
	http.Redirect(res, req, "/topics", http.StatusFound)
}

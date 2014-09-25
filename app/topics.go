// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mssola/todo/lib"
	"github.com/nu7hatch/gouuid"
)

// A topic is my way to divide different "contexts" inside my To Do list.
type Topic struct {
	Id         string
	Name       string
	Contents   string
	Created_at time.Time
}

// Given a name, try to create a new topic.
func createTopic(name string) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	t := &Topic{
		Id:   uuid.String(),
		Name: name,
	}
	return Db.Insert(t)
}

func renderShow(res http.ResponseWriter, topic *Topic) {
	o := &lib.ViewData{}
	lib.Render(res, "topics/show", o)
}

func TopicsIndex(res http.ResponseWriter, req *http.Request) {
	var t Topic

	Db.SelectOne(&t, "select * from topics order by name limit 1")
	renderShow(res, &t)
}

func TopicsCreate(res http.ResponseWriter, req *http.Request) {
	createTopic(req.FormValue("name"))
	http.Redirect(res, req, "/topics", http.StatusFound)
}

func TopicsShow(res http.ResponseWriter, req *http.Request) {
	var t Topic

	p := mux.Vars(req)
	Db.SelectOne(&t, "select * from topics where id=$1", p["id"])
	renderShow(res, &t)
}

func TopicsUpdate(res http.ResponseWriter, req *http.Request) {
	var t Topic

	// Get the original.
	p := mux.Vars(req)
	Db.SelectOne(&t, "select * from topics where id=$1", p["id"])

	// We can either rename, or change the contents, but not both things at the
	// same time.
	name := req.FormValue("name")
	if name != "" {
		t.Name = name
	} else {
		t.Contents = req.FormValue("contents")
	}
	Db.Update(&t)
	http.Redirect(res, req, "/topics", http.StatusFound)
}

func TopicsDestroy(res http.ResponseWriter, req *http.Request) {
	p := mux.Vars(req)
	Db.Exec("delete from topics where id=$1", p["id"])
	http.Redirect(res, req, "/topics", http.StatusFound)
}

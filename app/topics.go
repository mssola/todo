// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mssola/todo/lib"
	"github.com/nu7hatch/gouuid"
	"github.com/russross/blackfriday"
)

// A topic is my way to divide different "contexts" inside my To Do list.
type Topic struct {
	Id         string
	Name       string
	Contents   string
	Created_at time.Time
}

type TopicData struct {
	lib.ViewData

	Rendered string

	Current *Topic

	Topics []Topic
}

// Given a name, try to create a new topic.
func createTopic(name string) error {
	uuid, _ := uuid.NewV4()
	t := &Topic{
		Id:   uuid.String(),
		Name: name,
	}
	return Db.Insert(t)
}

func renderShow(res http.ResponseWriter, topic *Topic) {
	var topics []Topic
	Db.Select(&topics, "select * from topics order by name")

	// Render the contents
	unsafe := blackfriday.MarkdownCommon([]byte(topic.Contents))
	rs := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	// And render the page.
	o := &TopicData{
		Rendered: string(rs),
		Current:  topic,
		Topics:   topics,
	}
	o.JS = "topics"
	lib.Render(res, "topics/show", o)
}

func TopicsIndex(res http.ResponseWriter, req *http.Request) {
	var t Topic

	if id := lib.GetCookie(req, "topic"); id != "" {
		Db.SelectOne(&t, "select * from topics where id=$1", id)
	} else {
		Db.SelectOne(&t, "select * from topics order by name limit 1")
	}
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
	if t.Id != "" {
		lib.SetCookie(res, req, "topic", t.Id)
	}
	renderShow(res, &t)
}

func TopicsUpdate(res http.ResponseWriter, req *http.Request) {
	p := mux.Vars(req)

	// We can either rename, or change the contents, but not both things at the
	// same time.
	name := req.FormValue("name")
	if name != "" {
		Db.Exec("update topics set name=$1 where id=$2", name, p["id"])
	} else {
		cts := req.FormValue("contents")
		Db.Exec("update topics set contents=$1 where id=$2", cts, p["id"])
	}
	http.Redirect(res, req, "/topics", http.StatusFound)
}

func TopicsDestroy(res http.ResponseWriter, req *http.Request) {
	p := mux.Vars(req)
	Db.Exec("delete from topics where id=$1", p["id"])
	http.Redirect(res, req, "/topics", http.StatusFound)
}

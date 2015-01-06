// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
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
// Moreover this type also has the "Markdown" attribute. This attribute does
// not match any column from the database, but it comes handy in the API layer.
type Topic struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Contents   string    `json:"contents"`
	Created_at time.Time `json:"created_at"`
	Markdown   string    `db:"-" json:"markdown"`
}

// Generate the Markdown code for the current contents of this topic.
func (t *Topic) RenderMarkdown() {
	unsafe := blackfriday.MarkdownCommon([]byte(t.Contents))
	t.Markdown = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
}

// This is the data that will be assed to the renderShow function in order to
// render the main page.
type TopicData struct {
	lib.ViewData

	Current *Topic

	Topics []Topic
}

// Given a name, try to create a new topic.
func createTopic(name string) (*Topic, error) {
	uuid, _ := uuid.NewV4()
	t := &Topic{
		Id:         uuid.String(),
		Name:       name,
		Created_at: time.Now(),
	}
	return t, Db.Insert(t)
}

// Sends the main page with the given topic rendered in it as
// the current one.
func renderShow(res http.ResponseWriter, topic *Topic) {
	var topics []Topic
	Db.Select(&topics, "select * from topics order by name")
	topic.RenderMarkdown()

	// And render the page.
	o := &TopicData{
		Current: topic,
		Topics:  topics,
	}
	o.JS = "topics"
	lib.Render(res, "topics/show", o)
}

func TopicsIndex(res http.ResponseWriter, req *http.Request) {
	if lib.JsonEncoding(req) {
		TopicsApiIndex(res, req)
		return
	}

	var t Topic

	if id := lib.GetCookie(req, "topic"); id != "" && id != nil {
		Db.SelectOne(&t, "select * from topics where id=$1", id)
	} else {
		Db.SelectOne(&t, "select * from topics order by name limit 1")
	}
	renderShow(res, &t)
}

func TopicsCreate(res http.ResponseWriter, req *http.Request) {
	if lib.JsonEncoding(req) {
		TopicsApiCreate(res, req)
		return
	}

	t, err := createTopic(req.FormValue("name"))
	if err != nil {
		http.Redirect(res, req, "/topics", http.StatusFound)
	} else {
		http.Redirect(res, req, "/topics/"+t.Id, http.StatusFound)
	}
}

func TopicsShow(res http.ResponseWriter, req *http.Request) {
	if lib.JsonEncoding(req) {
		TopicsApiShow(res, req)
		return
	}

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

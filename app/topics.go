// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"log"
	"net/http"
	"time"

	"github.com/docker/distribution/uuid"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mssola/todo/lib"
	"github.com/russross/blackfriday"
)

// Topic is my way to divide different "contexts" inside my To Do list.
// Moreover this type also has the "Markdown" attribute. This attribute does
// not match any column from the database, but it comes handy in the API layer.
type Topic struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Contents  string    `json:"contents"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Markdown  string    `db:"-" json:"markdown"`
}

// RenderMarkdown generates the Markdown code for the current contents of this
// topic.
func (t *Topic) RenderMarkdown() {
	unsafe := blackfriday.MarkdownCommon([]byte(t.Contents))
	t.Markdown = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
}

// TopicData is the data that will be assed to the renderShow function in order
// to render the main page.
type TopicData struct {
	lib.ViewData

	Current *Topic

	Topics []Topic
}

// Given a name, try to create a new topic.
func createTopic(name string) (*Topic, error) {
	id := uuid.Generate().String()
	t := &Topic{ID: id, Name: name, CreatedAt: time.Now()}
	return t, Db.Insert(t)
}

// Sends the main page with the given topic rendered in it as
// the current one.
func renderShow(res http.ResponseWriter, topic *Topic) {
	var topics []Topic
	_, err := Db.Select(&topics, "select * from topics order by name")
	if err != nil {
		log.Printf("Select went wrong: %v", err)
	}
	topic.RenderMarkdown()

	// And render the page.
	o := &TopicData{
		Current: topic,
		Topics:  topics,
	}
	o.JS = "topics"
	lib.Render(res, "topics/show", o)
}

// TopicsIndex responds to: GET /topics
func TopicsIndex(res http.ResponseWriter, req *http.Request) {
	if lib.JSONEncoding(req) {
		TopicsIndexJSON(res, req)
		return
	}

	var err error
	var t Topic

	if id := lib.GetCookie(req, "topic"); id != "" && id != nil {
		err = Db.SelectOne(&t, "select * from topics where id=$1", id)
	} else {
		err = Db.SelectOne(&t, "select * from topics order by name limit 1")
	}
	if err != nil {
		log.Printf("Could not select topics: %v", err)
	}
	renderShow(res, &t)
}

// TopicsCreate responds to: POST /topics
func TopicsCreate(res http.ResponseWriter, req *http.Request) {
	if lib.JSONEncoding(req) {
		TopicsCreateJSON(res, req)
		return
	}

	if t, err := createTopic(req.FormValue("name")); err != nil {
		http.Redirect(res, req, "/topics", http.StatusFound)
	} else {
		http.Redirect(res, req, "/topics/"+t.ID, http.StatusFound)
	}
}

// TopicsShow responds to: GET /topics/:id
func TopicsShow(res http.ResponseWriter, req *http.Request) {
	if lib.JSONEncoding(req) {
		TopicsShowJSON(res, req)
		return
	}

	var t Topic

	p := mux.Vars(req)
	err := Db.SelectOne(&t, "select * from topics where id=$1", p["id"])
	if err != nil {
		log.Printf("Could not select topic: %v", err)
	}
	if t.ID != "" {
		lib.SetCookie(res, req, "topic", t.ID)
	}
	renderShow(res, &t)
}

// TopicsUpdate responds to: PUT/PATCH /posts/:id
func TopicsUpdate(res http.ResponseWriter, req *http.Request) {
	var err error
	p := mux.Vars(req)

	// We can either rename, or change the contents, but not both things at the
	// same time.
	name := req.FormValue("name")
	if name != "" {
		_, err = Db.Exec("update topics set name=$1 where id=$2", name, p["id"])
	} else {
		cts := req.FormValue("contents")
		_, err = Db.Exec("update topics set contents=$1 where id=$2", cts, p["id"])
	}
	if err != nil {
		log.Printf("Could not update topic: %v", err)
	}
	http.Redirect(res, req, "/topics", http.StatusFound)
}

// TopicsDestroy responds to: DELETE /posts/:id
func TopicsDestroy(res http.ResponseWriter, req *http.Request) {
	p := mux.Vars(req)
	_, err := Db.Exec("delete from topics where id=$1", p["id"])
	if err != nil {
		log.Printf("Could not perform delete of topic: %v", err)
	}
	http.Redirect(res, req, "/topics", http.StatusFound)
}

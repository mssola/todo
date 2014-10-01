// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"encoding/json"
	"fmt"
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
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Contents   string    `json:"contents"`
	Created_at time.Time `json:"created_at"`
}

type TopicData struct {
	lib.ViewData

	Rendered string

	Current *Topic

	Topics []Topic
}

// Given a name, try to create a new topic.
func createTopic(name string) (string, error) {
	uuid, _ := uuid.NewV4()
	t := &Topic{
		Id:   uuid.String(),
		Name: name,
	}
	return t.Id, Db.Insert(t)
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
	if lib.JsonEncoding(req) {
		var tr struct {
			Topics []Topic `json:"topics"`
		}
		Db.Select(&tr.Topics, "select * from topics")
		b, err := json.Marshal(tr)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprint(res, lib.Response{Error: "Failed!"})
		} else {
			fmt.Fprint(res, string(b))
		}
		return
	}

	var t Topic

	if id := lib.GetCookie(req, "topic"); id != "" {
		Db.SelectOne(&t, "select * from topics where id=$1", id)
	} else {
		Db.SelectOne(&t, "select * from topics order by name limit 1")
	}
	renderShow(res, &t)
}

func getValue(req *http.Request, name string) string {
	if lib.JsonEncoding(req) {
		decoder := json.NewDecoder(req.Body)

		var t struct{ Value string }
		err := decoder.Decode(&t)
		if err != nil {
			return ""
		}
		return t.Value
	}
	return req.FormValue(name)
}

func TopicsCreate(res http.ResponseWriter, req *http.Request) {
	id, err := createTopic(getValue(req, "name"))
	if lib.JsonEncoding(req) {
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprint(res, lib.Response{Error: "Failed!"})
		} else {
			t := struct {
				Id string `json:"id"`
			}{Id: id}
			b, err := json.Marshal(t)
			if err != nil {
				res.WriteHeader(http.StatusNotFound)
				fmt.Fprint(res, lib.Response{Error: "Failed!"})
			} else {
				fmt.Fprint(res, string(b))
			}
		}
	} else {
		http.Redirect(res, req, "/topics", http.StatusFound)
	}
}

type topicsShow struct {
	Topic

	Render string
}

func TopicsShow(res http.ResponseWriter, req *http.Request) {
	var t Topic

	p := mux.Vars(req)
	err := Db.SelectOne(&t, "select * from topics where id=$1", p["id"])
	if lib.JsonEncoding(req) {
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprint(res, lib.Response{Error: "Failed!"})
			return
		}

		var ts topicsShow
		ts.Contents = t.Contents
		ts.Id, ts.Name, ts.Created_at = t.Id, t.Name, t.Created_at
		unsafe := blackfriday.MarkdownCommon([]byte(t.Contents))
		ts.Render = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))

		b, err := json.Marshal(ts)
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprint(res, lib.Response{Error: "Failed!"})
		} else {
			fmt.Fprint(res, string(b))
		}
		return
	}

	if t.Id != "" {
		lib.SetCookie(res, req, "topic", t.Id)
	}
	renderShow(res, &t)
}

func TopicsUpdate(res http.ResponseWriter, req *http.Request) {
	p := mux.Vars(req)

	// We can either rename, or change the contents, but not both things at the
	// same time.
	name := getValue(req, "name")
	update := false
	var err error
	var cts string
	if name != "" {
		_, err = Db.Exec("update topics set name=$1 where id=$2", name, p["id"])
	} else {
		update = true
		cts = getValue(req, "contents")
		_, err = Db.Exec("update topics set contents=$1 where id=$2", cts, p["id"])
	}

	if lib.JsonEncoding(req) {
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprint(res, lib.Response{Error: "Failed!"})
		} else {
			if update {
				var ts struct {
					Render string `json:"contents"`
				}
				unsafe := blackfriday.MarkdownCommon([]byte(cts))
				ts.Render = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))

				b, err := json.Marshal(ts)
				if err != nil {
					res.WriteHeader(http.StatusNotFound)
					fmt.Fprint(res, lib.Response{Error: "Failed!"})
				} else {
					fmt.Fprint(res, string(b))
				}
			} else {
				fmt.Fprint(res, lib.Response{Message: "Ok"})
			}
		}
	} else {
		http.Redirect(res, req, "/topics", http.StatusFound)
	}
}

func TopicsDestroy(res http.ResponseWriter, req *http.Request) {
	p := mux.Vars(req)
	_, err := Db.Exec("delete from topics where id=$1", p["id"])

	if lib.JsonEncoding(req) {
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprint(res, lib.Response{Error: "Failed!"})
		} else {
			fmt.Fprint(res, lib.Response{Message: "Ok"})
		}
	} else {
		http.Redirect(res, req, "/topics", http.StatusFound)
	}
}

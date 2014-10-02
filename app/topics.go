// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"bytes"
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

// TODO: check and test for malformed JSON requests.

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
		b, _ := json.Marshal(tr)
		fmt.Fprint(res, string(b))
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

// TODO: re-think this once we split the _api thingie.
//          -> create a struct for JSON thingies with tags and shit ?

func getContents(req *http.Request, buffer bytes.Buffer) string {
	if lib.JsonEncoding(req) {
		reader := bytes.NewReader(buffer.Bytes())
		decoder := json.NewDecoder(reader)

		var t struct{ Contents string }
		err := decoder.Decode(&t)
		if err != nil {
			return ""
		}
		return t.Contents
	}
	return req.FormValue("contents")
}

func getNameFromBuffer(req *http.Request, buffer bytes.Buffer) string {
	if lib.JsonEncoding(req) {
		reader := bytes.NewReader(buffer.Bytes())
		decoder := json.NewDecoder(reader)

		var t struct{ Name string }
		err := decoder.Decode(&t)
		if err != nil {
			return ""
		}
		return t.Name
	}
	return req.FormValue("name")
}

func getName(req *http.Request) string {
	if lib.JsonEncoding(req) {
		decoder := json.NewDecoder(req.Body)

		var t struct{ Name string }
		err := decoder.Decode(&t)
		if err != nil {
			return ""
		}
		return t.Name
	}
	return req.FormValue("name")
}

func TopicsCreate(res http.ResponseWriter, req *http.Request) {
	id, err := createTopic(getName(req))
	if lib.JsonEncoding(req) {
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprint(res, lib.Response{Error: "Failed!"})
		} else {
			t := struct {
				Id string `json:"id"`
			}{Id: id}
			b, _ := json.Marshal(t)
			fmt.Fprint(res, string(b))
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

		b, _ := json.Marshal(ts)
		fmt.Fprint(res, string(b))
		return
	}

	if t.Id != "" {
		lib.SetCookie(res, req, "topic", t.Id)
	}
	renderShow(res, &t)
}

func TopicsUpdate(res http.ResponseWriter, req *http.Request) {
	p := mux.Vars(req)

	var buffer bytes.Buffer
	if lib.JsonEncoding(req) {
		// TODO: check if empty, check if couldn't read, check, check, ...
		_, _ = buffer.ReadFrom(req.Body)
	}

	// We can either rename, or change the contents, but not both things at the
	// same time.
	name := getNameFromBuffer(req, buffer)
	update := false
	var cts string
	if name != "" {
		Db.Exec("update topics set name=$1 where id=$2", name, p["id"])
	} else {
		update = true
		cts = getContents(req, buffer)
		Db.Exec("update topics set contents=$1 where id=$2", cts, p["id"])
	}

	if lib.JsonEncoding(req) {
		if update {
			var ts struct {
				Render string `json:"contents"`
			}
			unsafe := blackfriday.MarkdownCommon([]byte(cts))
			ts.Render = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))

			b, _ := json.Marshal(ts)
			fmt.Fprint(res, string(b))
		} else {
			fmt.Fprint(res, lib.Response{Message: "Ok"})
		}
	} else {
		http.Redirect(res, req, "/topics", http.StatusFound)
	}
}

func TopicsDestroy(res http.ResponseWriter, req *http.Request) {
	p := mux.Vars(req)
	Db.Exec("delete from topics where id=$1", p["id"])

	if lib.JsonEncoding(req) {
		fmt.Fprint(res, lib.Response{Message: "Ok"})
	} else {
		http.Redirect(res, req, "/topics", http.StatusFound)
	}
}

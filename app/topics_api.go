// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mssola/todo/lib"
)

// The parameters that be given through a request body.
type params struct {
	Name     string
	Contents string
}

// Get the possible parameters from the given request. Note that it will only
// check for the "name" and "contents" parameters.
func getFromBody(req *http.Request) *params {
	var p params

	if req.Body == nil {
		return nil
	}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&p); err != nil {
		return nil
	}
	return &p
}

// Safely render and send a JSON response with the given topic. This function
// should be called after performing some operation that might return an error.
// This error from the previous operation is the third parameter. The fourth
// parameter tells this function to generate the Markdown code for this topic.
func renderJson(res http.ResponseWriter, topic *Topic, err error, md bool) {
	// Try to render the given Topic.
	if err == nil {
		if md {
			topic.RenderMarkdown()
		}
		if b, err := json.Marshal(topic); err == nil {
			fmt.Fprint(res, string(b))
			return
		}
	}

	// Render a generic error.
	lib.JsonError(res)
}

func TopicsApiIndex(res http.ResponseWriter, req *http.Request) {
	var topics []Topic
	Db.Select(&topics, "select * from topics")
	b, _ := json.Marshal(topics)
	fmt.Fprint(res, string(b))
}

func TopicsApiCreate(res http.ResponseWriter, req *http.Request) {
	if p := getFromBody(req); p == nil {
		lib.JsonError(res)
	} else {
		t, err := createTopic(p.Name)
		renderJson(res, t, err, false)
	}
}

func TopicsApiShow(res http.ResponseWriter, req *http.Request) {
	var t Topic
	p := mux.Vars(req)
	err := Db.SelectOne(&t, "select * from topics where id=$1", p["id"])
	renderJson(res, &t, err, true)
}

func TopicsApiUpdate(res http.ResponseWriter, req *http.Request) {
	var str, value string
	var p *params

	if p = getFromBody(req); p == nil {
		lib.JsonError(res)
		return
	}

	// Execute the update query. Depending on the given parameters this will be
	// just a plain rename, or a full update.
	if value = p.Name; value == "" {
		str = "contents"
		if value = p.Contents; value == "" {
			lib.JsonError(res)
			return
		}
	} else {
		str = "name"
	}

	// And finally send the JSON response.
	var t Topic
	str = fmt.Sprintf("update topics set %v=$1 where id=$2 returning *", str)
	err := Db.SelectOne(&t, str, value, mux.Vars(req)["id"])
	renderJson(res, &t, err, true)
}

func TopicsApiDestroy(res http.ResponseWriter, req *http.Request) {
	p := mux.Vars(req)
	results, err := Db.Exec("delete from topics where id=$1", p["id"])

	if err != nil {
		fmt.Fprint(res, lib.Response{Error: "Could not remove topic"})
	} else if count, _ := results.RowsAffected(); count == 0 {
		fmt.Fprint(res, lib.Response{Error: "Could not remove topic"})
	} else {
		fmt.Fprint(res, lib.Response{Message: "Ok"})
	}
}

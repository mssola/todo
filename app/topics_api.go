// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mssola/todo/lib"
)

// Returns the name parameter as given by the body of the request.
func getName(req *http.Request) string {
	var m struct{ Name string }

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&m); err != nil {
		return ""
	}
	return m.Name
}

// Returns the value of the given element inside the buffer.
func getValue(name string, bf bytes.Buffer) (string, error) {
	var m map[string]string
	err := json.Unmarshal(bf.Bytes(), &m)
	if err != nil {
		return "", err
	}
	return m[name], nil
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
	if name := getName(req); name == "" {
		lib.JsonError(res)
	} else {
		t, err := createTopic(name)
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
	var err error
	var str string
	var buffer bytes.Buffer

	// Keep the body of the request in a buffer so we can read it multiple
	// times.
	if req.Body == nil {
		lib.JsonError(res)
		return
	}
	if _, err := buffer.ReadFrom(req.Body); err != nil {
		lib.JsonError(res)
		return
	}

	// Execute the update query. Depending on the given parameters this will be
	// just a plain rename, or a full update.
	value, err := getValue("name", buffer)
	if value == "" && err == nil {
		str = "contents"
		value, err = getValue("contents", buffer)
	} else {
		str = "name"
	}

	// And finally send the JSON response.
	var t Topic
	str = fmt.Sprintf("update topics set %v=$1 where id=$2 returning *", str)
	err = Db.SelectOne(&t, str, value, mux.Vars(req)["id"])
	renderJson(res, &t, err, true)
}

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
	"github.com/microcosm-cc/bluemonday"
	"github.com/mssola/todo/lib"
	"github.com/russross/blackfriday"
)

// TODO: re-think this once we split the _api thingie.
//          -> create a struct for JSON thingies with tags and shit ?

func getName(req *http.Request) string {
	decoder := json.NewDecoder(req.Body)

	var t struct{ Name string }
	err := decoder.Decode(&t)
	if err != nil {
		return ""
	}
	return t.Name
}

func getNameFromBuffer(req *http.Request, buffer bytes.Buffer) string {
	reader := bytes.NewReader(buffer.Bytes())
	decoder := json.NewDecoder(reader)

	var t struct{ Name string }
	err := decoder.Decode(&t)
	if err != nil {
		return ""
	}
	return t.Name
}

func getContents(req *http.Request, buffer bytes.Buffer) string {
	reader := bytes.NewReader(buffer.Bytes())
	decoder := json.NewDecoder(reader)

	var t struct{ Contents string }
	err := decoder.Decode(&t)
	if err != nil {
		return ""
	}
	return t.Contents
}

// TODO: check json.Marshal always.
// TODO: malformed JSON (+ tests)

func TopicsApiIndex(res http.ResponseWriter, req *http.Request) {
	var tr struct {
		Topics []Topic `json:"topics"`
	}
	Db.Select(&tr.Topics, "select * from topics")
	b, _ := json.Marshal(tr)
	fmt.Fprint(res, string(b))
}

func TopicsApiCreate(res http.ResponseWriter, req *http.Request) {
	id, err := createTopic(getName(req))

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
}

func TopicsApiShow(res http.ResponseWriter, req *http.Request) {
	var t Topic

	p := mux.Vars(req)
	err := Db.SelectOne(&t, "select * from topics where id=$1", p["id"])

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
}

func TopicsApiUpdate(res http.ResponseWriter, req *http.Request) {
	var buffer bytes.Buffer
	p := mux.Vars(req)

	// TODO: check if empty, check if couldn't read, check, check, ...
	_, _ = buffer.ReadFrom(req.Body)

	name := getNameFromBuffer(req, buffer)
	if name != "" {
		Db.Exec("update topics set name=$1 where id=$2", name, p["id"])
		fmt.Fprint(res, lib.Response{Message: "Ok"})
	} else {
		cts := getContents(req, buffer)
		Db.Exec("update topics set contents=$1 where id=$2", cts, p["id"])

		var ts struct {
			Render string `json:"contents"`
		}
		unsafe := blackfriday.MarkdownCommon([]byte(cts))
		ts.Render = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))

		b, _ := json.Marshal(ts)
		fmt.Fprint(res, string(b))
	}
}

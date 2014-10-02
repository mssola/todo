// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// The default response for short messages.
type Response struct {
	Message string `json:"msg,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Concatenate this response by marshalling it into JSON.
func (r Response) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

// Sends the standard error for this application.
func JsonError(res http.ResponseWriter) {
	res.WriteHeader(http.StatusNotFound)
	fmt.Fprint(res, Response{Error: "Failed!"})
}

// TODO: not sure about this one.
func CheckError(res http.ResponseWriter, req *http.Request, err error) bool {
	if err == nil {
		return false
	}

	if JsonEncoding(req) {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, Response{Error: "Failed!"})
	} else {
		http.Redirect(res, req, "/", http.StatusFound)
	}
	return true
}

func checkHeader(req *http.Request, name string) bool {
	ct := req.Header[name]
	if len(ct) != 1 {
		return false
	}
	return ct[0] == "application/json"
}

// TODO: test
func JsonEncoding(req *http.Request) bool {
	if checkHeader(req, "Content-Type") {
		return true
	}
	if checkHeader(req, "Accept") {
		return true
	}
	return strings.HasSuffix(req.URL.Path, ".json")
}

func GetUserId(req *http.Request) string {
	if JsonEncoding(req) {
		return req.URL.Query().Get("userId")
	}
	if id, ok := GetCookie(req, "userId").(string); ok {
		return id
	}
	return ""
}

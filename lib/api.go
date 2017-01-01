// Copyright (C) 2014-2017 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response for short messages.
type Response struct {
	Message string `json:"msg,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Concatenate this response by marshalling it into JSON.
func (r Response) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// JSONError sends the standard error for this application.
func JSONError(res http.ResponseWriter) {
	res.WriteHeader(http.StatusNotFound)
	fmt.Fprint(res, Response{Error: "Failed!"})
}

// CheckError returns true and sends the proper response if the given error is
// not nil. If the error is nil, then it just returns false and does nothing.
func CheckError(res http.ResponseWriter, req *http.Request, err error) bool {
	if err == nil {
		return false
	}

	if JSONEncoding(req) {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, Response{Error: "Failed!"})
	} else {
		http.Redirect(res, req, "/", http.StatusFound)
	}
	return true
}

// Returns true if the given header key has "application/json" as its value.
func checkHeader(req *http.Request, name string) bool {
	ct := req.Header[name]
	if len(ct) != 1 {
		return false
	}
	return ct[0] == "application/json"
}

// JSONEncoding returns true if we can assume that this is a JSON request,
// false otherwise.
func JSONEncoding(req *http.Request) bool {
	if checkHeader(req, "Content-Type") {
		return true
	}
	return checkHeader(req, "Accept")
}

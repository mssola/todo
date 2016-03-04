// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package lib

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponse(t *testing.T) {
	r := Response{
		Message: "message",
		Error:   "error",
	}
	str := "{\"msg\":\"message\",\"error\":\"error\"}"
	if s := r.String(); s != str {
		t.Fatalf("Got: %v, Expecting: %v", s, str)
	}

	r1 := Response{
		Error: "error",
	}
	str = "{\"error\":\"error\"}"
	if s := r1.String(); s != str {
		t.Fatalf("Got: %v, Expecting: %v", s, str)
	}

	r2 := Response{
		Message: "message",
	}
	str = "{\"msg\":\"message\"}"
	if s := r2.String(); s != str {
		t.Fatalf("Got: %v, Expecting: %v", s, str)
	}

	r3 := Response{}
	str = "{}"
	if s := r3.String(); s != str {
		t.Fatalf("Got: %v, Expecting: %v", s, str)
	}
}

func TestJsonError(t *testing.T) {
	w := httptest.NewRecorder()
	JSONError(w)
	if w.Code != http.StatusNotFound {
		t.Fatalf("Wrong code. Got: %v, expected 404", w.Code)
	}

	var r Response
	decoder := json.NewDecoder(w.Body)
	if err := decoder.Decode(&r); err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if err := r.Error; err != "Failed!" {
		t.Fatalf("Expected to be 'Failed!': %v", err)
	}
}

func TestCheckError(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Base case.
	if r := CheckError(w, request, nil); r {
		t.Fatalf("Expected to be false!")
	}

	// Regular HTML
	if b := CheckError(w, request, errors.New("Something")); !b {
		t.Fatalf("Expected to be true!")
	}
	if w.Code != http.StatusFound {
		t.Fatalf("Wrong code. Got: %v, expected 200", w.Code)
	}

	// JSON
	r1, _ := http.NewRequest("GET", "/", nil)
	r1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	if b := CheckError(w1, r1, errors.New("Something")); !b {
		t.Fatalf("Expected to be true!")
	}

	var r Response
	decoder := json.NewDecoder(w1.Body)
	err := decoder.Decode(&r)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if err := r.Error; err != "Failed!" {
		t.Fatalf("Expected to be 'Failed!': %v", err)
	}
}

func TestJsonEncoding(t *testing.T) {
	// Nope.
	r1, _ := http.NewRequest("GET", "/", nil)
	if res := JSONEncoding(r1); res {
		t.Fatalf("Expected to be false")
	}

	// Yes, because of the "Content-Type" header.
	r3, _ := http.NewRequest("GET", "/something", nil)
	r3.Header.Set("Content-Type", "application/json")
	if res := JSONEncoding(r3); !res {
		t.Fatalf("Expected to be true")
	}

	// Yes, because of the "Accept" header.
	r4, _ := http.NewRequest("GET", "/something", nil)
	r4.Header.Set("Accept", "application/json")
	if res := JSONEncoding(r4); !res {
		t.Fatalf("Expected to be true")
	}
}

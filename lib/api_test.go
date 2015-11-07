// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	r := Response{
		Message: "message",
		Error:   "error",
	}
	str := "{\"msg\":\"message\",\"error\":\"error\"}"
	assert.Equal(t, r.String(), str)

	r1 := Response{
		Error: "error",
	}
	str = "{\"error\":\"error\"}"
	assert.Equal(t, r1.String(), str)

	r2 := Response{
		Message: "message",
	}
	str = "{\"msg\":\"message\"}"
	assert.Equal(t, r2.String(), str)

	r3 := Response{}
	str = "{}"
	assert.Equal(t, r3.String(), str)
}

func TestJsonError(t *testing.T) {
	w := httptest.NewRecorder()
	JSONError(w)
	assert.Equal(t, w.Code, http.StatusNotFound)

	var r Response
	decoder := json.NewDecoder(w.Body)
	err := decoder.Decode(&r)
	assert.Nil(t, err)
	assert.Equal(t, r.Error, "Failed!")
}

func TestCheckError(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Base case.
	assert.False(t, CheckError(w, request, nil))

	// Regular HTML
	b := CheckError(w, request, errors.New("Something"))
	assert.True(t, b)
	assert.Equal(t, w.Code, http.StatusFound)

	// JSON
	r1, _ := http.NewRequest("GET", "/", nil)
	r1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	b = CheckError(w1, r1, errors.New("Something"))

	var r Response
	decoder := json.NewDecoder(w1.Body)
	err := decoder.Decode(&r)
	assert.Nil(t, err)
	assert.Equal(t, r.Error, "Failed!")
}

func TestJsonEncoding(t *testing.T) {
	// Nope.
	r1, _ := http.NewRequest("GET", "/", nil)
	assert.False(t, JSONEncoding(r1))

	// Yes, because of the "Content-Type" header.
	r3, _ := http.NewRequest("GET", "/something", nil)
	r3.Header.Set("Content-Type", "application/json")
	assert.True(t, JSONEncoding(r3))

	// Yes, because of the "Accept" header.
	r4, _ := http.NewRequest("GET", "/something", nil)
	r4.Header.Set("Accept", "application/json")
	assert.True(t, JSONEncoding(r4))
}

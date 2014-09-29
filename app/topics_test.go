// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestTopicsCreate(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	param := make(url.Values)
	param["name"] = []string{"user"}

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	req.PostForm = param
	w := httptest.NewRecorder()
	TopicsCreate(w, req)

	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/topics")

	var topic Topic
	err = Db.SelectOne(&topic, "select * from topics")
	assert.Nil(t, err)
	assert.NotEmpty(t, topic.Id)
	assert.Equal(t, topic.Name, "user")
	assert.NotEmpty(t, topic.Created_at)
}

func TestTopicsCreateNoName(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	param := make(url.Values)

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	req.PostForm = param
	w := httptest.NewRecorder()
	TopicsCreate(w, req)

	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/topics")

	var topic Topic
	err = Db.SelectOne(&topic, "select * from topics")
	assert.Empty(t, topic.Id)
	assert.NotNil(t, err)
	count, err := Db.SelectInt("select count(*) from users")
	assert.Equal(t, count, 0)
}

func TestTopicsRename(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	err := createTopic("topic")
	assert.Nil(t, err)

	var t1, t2 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	assert.Nil(t, err)

	param := make(url.Values)

	req, err := http.NewRequest("POST", "/topics/"+t1.Id, nil)
	assert.Nil(t, err)
	param["name"] = []string{"topic1"}
	req.PostForm = param
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsUpdate)
	m.ServeHTTP(w, req)

	// DB
	err = Db.SelectOne(&t2, "select * from topics")
	assert.Nil(t, err)
	assert.Equal(t, t2.Name, "topic1")
	assert.Equal(t, t1.Id, t2.Id)

	// HTTP
	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/topics")
}

func TestUpdateContents(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	err := createTopic("topic")
	assert.Nil(t, err)

	var t1, t2 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	assert.Nil(t, err)

	param := make(url.Values)

	req, err := http.NewRequest("POST", "/topics/"+t1.Id, nil)
	assert.Nil(t, err)
	param["contents"] = []string{"**bold**"}
	req.PostForm = param
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsUpdate)
	m.ServeHTTP(w, req)

	// DB
	err = Db.SelectOne(&t2, "select * from topics")
	assert.Nil(t, err)
	assert.Equal(t, t1.Name, t2.Name)
	assert.Equal(t, t1.Id, t2.Id)
	assert.NotEqual(t, t1.Contents, t2.Contents)
	assert.Equal(t, t2.Contents, "**bold**")

	// HTTP
	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/topics")
}

func TestDestroy(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	err := createTopic("topic")
	assert.Nil(t, err)

	var t1 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/topics/"+t1.Id+"/delete", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}/delete", TopicsDestroy)
	m.ServeHTTP(w, req)

	// DB
	c, err := Db.SelectInt("select count(*) from topics")
	assert.Nil(t, err)
	assert.Equal(t, c, 0)

	// HTTP
	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/topics")
}

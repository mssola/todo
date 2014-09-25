// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/mssola/todo/app/models"
	"github.com/stretchr/testify/assert"
)

func TestTopicsCreate(t *testing.T) {
	models.InitTestDB()
	defer models.CloseDB()

	param := make(url.Values)
	param["name"] = []string{"user"}

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	req.PostForm = param
	w := httptest.NewRecorder()
	TopicsCreate(w, req)

	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/topics")

	var topic models.Topic
	err = models.Db.SelectOne(&topic, "select * from topics")
	assert.Nil(t, err)
	assert.NotEmpty(t, topic.Id)
	assert.Equal(t, topic.Name, "user")
	assert.NotEmpty(t, topic.Created_at)
}

func TestTopicsCreateNoName(t *testing.T) {
	models.InitTestDB()
	defer models.CloseDB()

	param := make(url.Values)

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	req.PostForm = param
	w := httptest.NewRecorder()
	TopicsCreate(w, req)

	assert.Equal(t, w.Code, 302)
	assert.Equal(t, w.HeaderMap["Location"][0], "/topics")

	var topic models.Topic
	err = models.Db.SelectOne(&topic, "select * from topics")
	assert.Empty(t, topic.Id)
	assert.NotNil(t, err)
	count, err := models.Db.SelectInt("select count(*) from users")
	assert.Equal(t, count, 0)
}

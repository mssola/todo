// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mssola/todo/lib"
	"github.com/stretchr/testify/assert"
)

func TestTopicsIndex(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	_, err := createTopic("topic1")
	assert.Nil(t, err)
	_, err = createTopic("topic2")
	assert.Nil(t, err)

	var t1, t2 Topic
	err = Db.SelectOne(&t1, "select * from topics where name=$1", "topic1")
	assert.Nil(t, err)
	err = Db.SelectOne(&t2, "select * from topics where name=$1", "topic2")
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	TopicsIndex(w, req)

	str, _ := ioutil.ReadAll(w.Body)
	s := string(str)
	s1 := "<li class=\"selected\"><a href=\"/topics/%v\">topic1</a></li>"
	s2 := "<li><a href=\"/topics/%v\">topic2</a></li>"
	assert.True(t, strings.Contains(s, fmt.Sprintf(s1, t1.Id)))
	assert.True(t, strings.Contains(s, fmt.Sprintf(s2, t2.Id)))
}

func TestTopicsIndexWithCookie(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	_, err := createTopic("topic1")
	assert.Nil(t, err)
	_, err = createTopic("topic2")
	assert.Nil(t, err)

	var t1, t2 Topic
	err = Db.SelectOne(&t1, "select * from topics where name=$1", "topic1")
	assert.Nil(t, err)
	err = Db.SelectOne(&t2, "select * from topics where name=$1", "topic2")
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()
	lib.SetCookie(w, req, "topic", t2.Id)
	TopicsIndex(w, req)

	str, _ := ioutil.ReadAll(w.Body)
	s := string(str)
	s1 := "<li><a href=\"/topics/%v\">topic1</a></li>"
	s2 := "<li class=\"selected\"><a href=\"/topics/%v\">topic2</a></li>"
	assert.True(t, strings.Contains(s, fmt.Sprintf(s1, t1.Id)))
	assert.True(t, strings.Contains(s, fmt.Sprintf(s2, t2.Id)))
}

func TestTopicsIndexJson(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	_, err := createTopic("topic1")
	assert.Nil(t, err)
	_, err = createTopic("topic2")
	assert.Nil(t, err)

	var topicsDb [2]Topic
	err = Db.SelectOne(&topicsDb[0], "select * from topics where name=$1", "topic1")
	assert.Nil(t, err)
	err = Db.SelectOne(&topicsDb[1], "select * from topics where name=$1", "topic2")
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	TopicsIndex(w, req)

	var data map[string][]Topic
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&data)
	assert.Nil(t, err)

	for i := 0; i < 2; i++ {
		assert.Equal(t, data["topics"][i].Id, topicsDb[i].Id)
		assert.Equal(t, data["topics"][i].Name, topicsDb[i].Name)
	}
}

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

func TestTopicsCreateJson(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	// We try to create a topic for the first time.
	body := "{\"name\":\"mssola\"}"
	reader := strings.NewReader(body)
	req, err := http.NewRequest("POST", "/topics", reader)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics", TopicsCreate)
	m.ServeHTTP(w, req)

	var u struct{ Id string }
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&u)
	assert.Nil(t, err)

	var t1 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	assert.Nil(t, err)
	assert.Equal(t, t1.Id, u.Id)

	// Now let's try it with the same name.
	reader2 := strings.NewReader(body)
	req, err = http.NewRequest("POST", "/topics", reader2)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()

	m = mux.NewRouter()
	m.HandleFunc("/topics", TopicsCreate)
	m.ServeHTTP(w1, req)

	var resp lib.Response
	decoder = json.NewDecoder(w1.Body)
	err = decoder.Decode(&resp)
	assert.Nil(t, err)
	assert.Equal(t, resp.Error, "Failed!")
}

func TestTopicsShow(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	_, err := createTopic("topic1")
	assert.Nil(t, err)
	_, err = createTopic("topic2")
	assert.Nil(t, err)

	var topicsDb [2]Topic
	err = Db.SelectOne(&topicsDb[0], "select * from topics where name=$1", "topic1")
	assert.Nil(t, err)
	err = Db.SelectOne(&topicsDb[1], "select * from topics where name=$1", "topic2")
	assert.Nil(t, err)

	req, err := http.NewRequest("GET", "/topics/"+topicsDb[1].Id, nil)
	assert.Nil(t, err)
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsShow)
	m.ServeHTTP(w, req)

	str, _ := ioutil.ReadAll(w.Body)
	s := string(str)

	s1 := "<li><a href=\"/topics/%v\">topic1</a></li>"
	s2 := "<li class=\"selected\"><a href=\"/topics/%v\">topic2</a></li>"
	assert.True(t, strings.Contains(s, fmt.Sprintf(s1, topicsDb[0].Id)))
	assert.True(t, strings.Contains(s, fmt.Sprintf(s2, topicsDb[1].Id)))
}

func TestTopicsShowJson(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	_, err := createTopic("topic1")
	assert.Nil(t, err)
	_, err = createTopic("topic2")
	Db.Exec("update topics set contents=$1 where name=$2", "**co**", "topic2")
	assert.Nil(t, err)

	var topicsDb [2]Topic
	err = Db.SelectOne(&topicsDb[0], "select * from topics where name=$1", "topic1")
	assert.Nil(t, err)
	err = Db.SelectOne(&topicsDb[1], "select * from topics where name=$1", "topic2")
	assert.Nil(t, err)

	req, err := http.NewRequest("GET", "/topics/"+topicsDb[1].Id, nil)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsShow)
	m.ServeHTTP(w, req)

	var topic struct {
		Topic
		Render string
	}
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&topic)
	assert.Nil(t, err)

	assert.Equal(t, topicsDb[1].Id, topic.Id)
	assert.Equal(t, topicsDb[1].Name, topic.Name)
	assert.Equal(t, topicsDb[1].Contents, topic.Contents)
	rendered := strings.TrimSpace(topic.Render)
	assert.Equal(t, "<p><strong>co</strong></p>", rendered)

	// A non-existant topic.
	req, err = http.NewRequest("GET", "/topics/1", nil)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()

	m.ServeHTTP(w1, req)

	var response lib.Response
	decoder = json.NewDecoder(w1.Body)
	err = decoder.Decode(&response)
	assert.Nil(t, err)

	assert.Equal(t, response.Error, "Failed!")
	assert.Empty(t, response.Message)
}

func TestTopicsRename(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	_, err := createTopic("topic")
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

func TestTopicsRenameJson(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	_, err := createTopic("topic")
	assert.Nil(t, err)

	var t1, t2 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	assert.Nil(t, err)

	body := strings.NewReader("{\"name\":\"topic1\"}")
	req, err := http.NewRequest("POST", "/topics/"+t1.Id, body)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsUpdate)
	m.ServeHTTP(w, req)

	// DB
	err = Db.SelectOne(&t2, "select * from topics")
	assert.Nil(t, err)
	assert.Equal(t, t2.Name, "topic1")
	assert.Equal(t, t1.Id, t2.Id)

	// Response
	var resp lib.Response
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	assert.Nil(t, err)
	assert.Equal(t, resp.Message, "Ok")
}

func TestUpdateContents(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	_, err := createTopic("topic")
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

func TestUpdateJson(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	_, err := createTopic("topic")
	assert.Nil(t, err)

	var t1, t2 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	assert.Nil(t, err)

	body := strings.NewReader("{\"contents\":\"**contents**\"}")
	req, err := http.NewRequest("POST", "/topics/"+t1.Id, body)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
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
	assert.Equal(t, t2.Contents, "**contents**")

	// Response
	var st struct {
		Contents string
	}
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&st)
	assert.Nil(t, err)
	rendered := strings.TrimSpace(st.Contents)
	assert.Equal(t, rendered, "<p><strong>contents</strong></p>")
}

func TestDestroy(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	_, err := createTopic("topic")
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

func TestDestroyJson(t *testing.T) {
	InitTestDB()
	defer CloseTestDB()

	_, err := createTopic("topic")
	assert.Nil(t, err)

	var t1 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/topics/"+t1.Id+"/delete", nil)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}/delete", TopicsDestroy)
	m.ServeHTTP(w, req)

	// DB
	c, err := Db.SelectInt("select count(*) from topics")
	assert.Nil(t, err)
	assert.Equal(t, c, 0)

	// HTTP
	var resp lib.Response
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	assert.Nil(t, err)
	assert.Equal(t, resp.Message, "Ok")
}

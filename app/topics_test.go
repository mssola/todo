// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
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
)

func TestTopicsIndex(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic1")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	_, err = createTopic("topic2")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var t1, t2 Topic
	err = Db.SelectOne(&t1, "select * from topics where name=$1", "topic1")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	err = Db.SelectOne(&t2, "select * from topics where name=$1", "topic2")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	w := httptest.NewRecorder()
	TopicsIndex(w, req)

	str, _ := ioutil.ReadAll(w.Body)
	s := string(str)
	s1 := "<li class=\"selected\"><a href=\"/topics/%v\">topic1</a></li>"
	s2 := "<li><a href=\"/topics/%v\">topic2</a></li>"
	if !strings.Contains(s, fmt.Sprintf(s1, t1.ID)) {
		t.Fatalf("S: %v; Should've containerd: %v", s, fmt.Sprintf(s1, t1.ID))
	}
	if !strings.Contains(s, fmt.Sprintf(s2, t2.ID)) {
		t.Fatalf("S: %v; Should've containerd: %v", s, fmt.Sprintf(s2, t2.ID))
	}
}

func TestTopicsIndexWithCookie(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic1")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	_, err = createTopic("topic2")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var t1, t2 Topic
	err = Db.SelectOne(&t1, "select * from topics where name=$1", "topic1")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	err = Db.SelectOne(&t2, "select * from topics where name=$1", "topic2")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	w := httptest.NewRecorder()
	lib.SetCookie(w, req, "topic", t2.ID)
	TopicsIndex(w, req)

	str, _ := ioutil.ReadAll(w.Body)
	s := string(str)
	s1 := "<li><a href=\"/topics/%v\">topic1</a></li>"
	s2 := "<li class=\"selected\"><a href=\"/topics/%v\">topic2</a></li>"
	if !strings.Contains(s, fmt.Sprintf(s1, t1.ID)) {
		t.Fatalf("S: %v; Should've containerd: %v", s, fmt.Sprintf(s1, t1.ID))
	}
	if !strings.Contains(s, fmt.Sprintf(s2, t2.ID)) {
		t.Fatalf("S: %v; Should've containerd: %v", s, fmt.Sprintf(s2, t2.ID))
	}
}

func TestTopicsIndexJson(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic1")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	_, err = createTopic("topic2")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var topicsDb [2]Topic
	err = Db.SelectOne(&topicsDb[0], "select * from topics where name=$1", "topic1")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	err = Db.SelectOne(&topicsDb[1], "select * from topics where name=$1", "topic2")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	TopicsIndex(w, req)

	var topics []Topic
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&topics)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	for i := 0; i < 2; i++ {
		if topics[i].ID != topicsDb[i].ID {
			t.Fatalf("Got %v, Expected: %v", topics[i].ID, topicsDb[i].ID)
		}
		if topics[i].Name != topicsDb[i].Name {
			t.Fatalf("Got %v, Expected: %v", topics[i].Name, topicsDb[i].Name)
		}
	}
}

func TestTopicsCreate(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	param := make(url.Values)
	param["name"] = []string{"user"}

	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.PostForm = param
	w := httptest.NewRecorder()
	TopicsCreate(w, req)

	var topic Topic
	err = Db.SelectOne(&topic, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if topic.ID == "" {
		t.Fatalf("Expected to not be empty")
	}
	if topic.Name != "user" {
		t.Fatalf("Got %v, Expected: %v", topic.Name, "user")
	}

	if w.Code != 302 {
		t.Fatalf("Got %v, Expected: %v", w.Code, 302)
	}
	if w.HeaderMap["Location"][0] != "/topics/"+topic.ID {
		t.Fatalf("Got %v, Expected: %v", w.HeaderMap["Location"][0], "/topics/"+topic.ID)
	}
}

func TestTopicsCreateNoName(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	param := make(url.Values)

	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.PostForm = param
	w := httptest.NewRecorder()
	TopicsCreate(w, req)

	if w.Code != 302 {
		t.Fatalf("Got %v; Expected: %v", w.Code, 302)
	}
	if w.HeaderMap["Location"][0] != "/topics" {
		t.Fatalf("Got %v; Expected: %v", w.HeaderMap["Location"][0], "/topics")
	}

	var topic Topic
	err = Db.SelectOne(&topic, "select * from topics")
	if topic.ID != "" {
		t.Fatalf("Should be empty")
	}
	if err == nil {
		t.Fatalf("Expected to not be nil")
	}
	count, err := Db.SelectInt("select count(*) from users")
	if count != 0 {
		t.Fatalf("Got %v; Expected: %v", count, 0)
	}
}

func TestTopicsCreateJson(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	// We try to create a topic for the first time.
	body := "{\"name\":\"mssola\"}"
	reader := strings.NewReader(body)
	req, err := http.NewRequest("POST", "/topics", reader)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics", TopicsCreate)
	m.ServeHTTP(w, req)

	var topic Topic
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&topic)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var t1 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if t1.ID != topic.ID {
		t.Fatalf("Got %v; Expected: %v", t1.ID, topic.ID)
	}
	if t1.Name != topic.Name {
		t.Fatalf("Got %v; Expected: %v", t1.Name, topic.Name)
	}

	// Now let's try it with the same name.
	reader2 := strings.NewReader(body)
	req, err = http.NewRequest("POST", "/topics", reader2)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()

	m = mux.NewRouter()
	m.HandleFunc("/topics", TopicsCreate)
	m.ServeHTTP(w1, req)

	var resp lib.Response
	decoder = json.NewDecoder(w1.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if resp.Error != "Failed!" {
		t.Fatalf("Got %v; Expected: %v", resp.Error, "Failed!")
	}
}

func TestTopicsCreateJsonMalformed(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	// We try to create a topic for the first time.
	body := "{\"name\":\"mssola\""
	reader := strings.NewReader(body)
	req, err := http.NewRequest("POST", "/topics", reader)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics", TopicsCreate)
	m.ServeHTTP(w, req)

	var topic lib.Response
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&topic)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if topic.Error != "Failed!" {
		t.Fatalf("Got %v; Expected: %v", topic.Error, "Failed!")
	}
	if w.Code != 404 {
		t.Fatalf("Got %v; Expected: %v", w.Code, 404)
	}
}

func TestTopicsShow(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic1")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	_, err = createTopic("topic2")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var topicsDb [2]Topic
	err = Db.SelectOne(&topicsDb[0], "select * from topics where name=$1", "topic1")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	err = Db.SelectOne(&topicsDb[1], "select * from topics where name=$1", "topic2")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	req, err := http.NewRequest("GET", "/topics/"+topicsDb[1].ID, nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsShow)
	m.ServeHTTP(w, req)

	str, _ := ioutil.ReadAll(w.Body)
	s := string(str)

	s1 := "<li><a href=\"/topics/%v\">topic1</a></li>"
	s2 := "<li class=\"selected\"><a href=\"/topics/%v\">topic2</a></li>"
	if !strings.Contains(s, fmt.Sprintf(s1, topicsDb[0].ID)) {
		t.Fatalf("S: %v; Should've containerd: %v", s, fmt.Sprintf(s1, topicsDb[0].ID))
	}
	if !strings.Contains(s, fmt.Sprintf(s2, topicsDb[1].ID)) {
		t.Fatalf("S: %v; Should've containerd: %v", s, fmt.Sprintf(s2, topicsDb[1].ID))
	}
}

func TestTopicsShowJson(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic1")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	_, err = createTopic("topic2")
	Db.Exec("update topics set contents=$1 where name=$2", "**co**", "topic2")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var topicsDb [2]Topic
	err = Db.SelectOne(&topicsDb[0], "select * from topics where name=$1", "topic1")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	err = Db.SelectOne(&topicsDb[1], "select * from topics where name=$1", "topic2")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	req, err := http.NewRequest("GET", "/topics/"+topicsDb[1].ID, nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsShow)
	m.ServeHTTP(w, req)

	var topic Topic
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&topic)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	if topicsDb[1].ID != topic.ID {
		t.Fatalf("Got %v; Expected: %v", topicsDb[1].ID, topic.ID)
	}
	if topicsDb[1].Name != topic.Name {
		t.Fatalf("Got %v; Expected: %v", topicsDb[1].Name, topic.Name)
	}
	if topicsDb[1].Contents != topic.Contents {
		t.Fatalf("Got %v; Expected: %v", topicsDb[1].Contents, topic.Contents)
	}
	rendered := strings.TrimSpace(topic.Markdown)
	if "<p><strong>co</strong></p>" != rendered {
		t.Fatalf("Got %v; Expected: %v", "<p><strong>co</strong></p>", rendered)
	}

	// A non-existant topic.
	req, err = http.NewRequest("GET", "/topics/1", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()

	m.ServeHTTP(w1, req)

	var response lib.Response
	decoder = json.NewDecoder(w1.Body)
	err = decoder.Decode(&response)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	if response.Error != "Failed!" {
		t.Fatalf("Got %v; Expected: %v", response.Error, "Failed!")
	}
	if response.Message != "" {
		t.Fatal("Should not be empty")
	}
}

func TestTopicsShowJsonFail(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	req, err := http.NewRequest("GET", "/topics/1", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsShow)
	m.ServeHTTP(w, req)

	var resp lib.Response
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if resp.Error != "Failed!" {
		t.Fatalf("Got %v; Expected: %v", resp.Error, "Failed!")
	}
	if w.Code != 404 {
		t.Fatalf("Got %v; Expected: %v", w.Code, 404)
	}
}

func TestTopicsRename(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var t1, t2 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	param := make(url.Values)

	req, err := http.NewRequest("POST", "/topics/"+t1.ID, nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	param["name"] = []string{"topic1"}
	req.PostForm = param
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsUpdate)
	m.ServeHTTP(w, req)

	// DB
	err = Db.SelectOne(&t2, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if t2.Name != "topic1" {
		t.Fatalf("Got %v; Expected: %v", t2.Name, "topic1")
	}
	if t1.ID != t2.ID {
		t.Fatalf("Got %v; Expected: %v", t1.ID, t2.ID)
	}

	// HTTP
	if w.Code != 302 {
		t.Fatalf("Got %v; Expected: %v", w.Code, 302)
	}
	if w.HeaderMap["Location"][0] != "/topics" {
		t.Fatalf("Got %v; Expected: %v", w.HeaderMap["Location"][0], "/topics")
	}
}

func TestTopicsRenameJson(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var t1, t2 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	body := strings.NewReader("{\"name\":\"topic1\"}")
	req, err := http.NewRequest("PUT", "/topics/"+t1.ID, body)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsUpdateJSON)
	m.ServeHTTP(w, req)

	// DB
	var resp Topic
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	err = Db.SelectOne(&t2, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if t2.Name != "topic1" {
		t.Fatalf("Got %v; Expected: %v", t2.Name, "topic1")
	}
	if t2.Name != resp.Name {
		t.Fatalf("Got %v; Expected: %v", t2.Name, resp.Name)
	}
	if t1.ID != t2.ID {
		t.Fatalf("Got %v; Expected: %v", t1.ID, t2.ID)
	}
	if t2.ID != resp.ID {
		t.Fatalf("Got %v; Expected: %v", t2.ID, resp.ID)
	}
}

func TestTopicsRenameJsonMalformed(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var t1 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	body := strings.NewReader("{\"name\":\"topic1\"")
	req, err := http.NewRequest("PUT", "/topics/"+t1.ID, body)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsUpdateJSON)
	m.ServeHTTP(w, req)

	// DB
	var resp lib.Response
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if resp.Error != "Failed!" {
		t.Fatalf("Got %v; Expected: %v", resp.Error, "Failed!")
	}
	if w.Code != 404 {
		t.Fatalf("Got %v; Expected: %v", w.Code, 404)
	}
}

func TestTopicsRenameJsonFail(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	_, err = createTopic("topic1")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var t1 Topic
	err = Db.SelectOne(&t1, "select * from topics where name=$1", "topic")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	body := strings.NewReader("{\"name\":\"topic1\"}")
	req, err := http.NewRequest("PUT", "/topics/"+t1.ID, body)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsUpdateJSON)
	m.ServeHTTP(w, req)

	// DB
	var resp lib.Response
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if resp.Error != "Failed!" {
		t.Fatalf("Got %v; Expected: %v", resp.Error, "Failed!")
	}
	if w.Code != 404 {
		t.Fatalf("Got %v; Expected: %v", w.Code, 404)
	}
}

func TestUpdateContents(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var t1, t2 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	param := make(url.Values)

	req, err := http.NewRequest("POST", "/topics/"+t1.ID, nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	param["contents"] = []string{"**bold**"}
	req.PostForm = param
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsUpdate)
	m.ServeHTTP(w, req)

	// DB
	err = Db.SelectOne(&t2, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if t1.Name != t2.Name {
		t.Fatalf("Got %v; Expected: %v", t1.Name, t2.Name)
	}
	if t1.ID != t2.ID {
		t.Fatalf("Got %v; Expected: %v", t1.ID, t2.ID)
	}
	if t1.Contents == t2.Contents {
		t.Fatalf("%v -- %v;; should be different", t1.Contents, t2.Contents)
	}
	if t2.Contents != "**bold**" {
		t.Fatalf("Got %v; Expected: %v", t2.Contents, "**bold**")
	}

	// HTTP
	if w.Code != 302 {
		t.Fatalf("Got %v; Expected: %v", w.Code, 302)
	}
	if w.HeaderMap["Location"][0] != "/topics" {
		t.Fatalf("Got %v; Expected: %v", w.HeaderMap["Location"][0], "/topics")
	}
}

func TestUpdateJson(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var t1, t2 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	body := strings.NewReader("{\"contents\":\"**contents**\"}")
	req, err := http.NewRequest("PUT", "/topics/"+t1.ID, body)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsUpdateJSON)
	m.ServeHTTP(w, req)

	// DB
	err = Db.SelectOne(&t2, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if t1.Name != t2.Name {
		t.Fatalf("Got %v; Expected: %v", t1.Name, t2.Name)
	}
	if t1.ID != t2.ID {
		t.Fatalf("Got %v; Expected: %v", t1.ID, t2.ID)
	}
	if t1.Contents == t2.Contents {
		t.Fatalf("%v -- %v;; should be different", t1.Contents, t2.Contents)
	}
	if t2.Contents != "**contents**" {
		t.Fatalf("Got %v; Expected: %v", t2.Contents, "**contents**")
	}

	// Response
	var st Topic
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&st)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	rendered := strings.TrimSpace(st.Markdown)
	if rendered != "<p><strong>contents</strong></p>" {
		t.Fatalf("Got %v; Expected: %v", rendered, "<p><strong>contents</strong></p>")
	}
}

func TestTopicsUpdateJsonMalformed(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var t1 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	body := strings.NewReader("{\"contents\":\"topic1\"")
	req, err := http.NewRequest("PUT", "/topics/"+t1.ID, body)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsUpdateJSON)
	m.ServeHTTP(w, req)

	// DB
	var resp lib.Response
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if resp.Error != "Failed!" {
		t.Fatalf("Got %v; Expected: %v", resp.Error, "Failed!")
	}
	if w.Code != 404 {
		t.Fatalf("Got %v; Expected: %v", w.Code, 404)
	}
}

func TestTopicsUpdateNoBody(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var t1 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	req, err := http.NewRequest("PUT", "/topics/"+t1.ID, nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsUpdateJSON)
	m.ServeHTTP(w, req)

	// DB
	var resp lib.Response
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if resp.Error != "Failed!" {
		t.Fatalf("Got %v; Expected: %v", resp.Error, "Failed!")
	}
	if w.Code != 404 {
		t.Fatalf("Got %v; Expected: %v", w.Code, 404)
	}
}

func TestTopicsUpdateNoValidParametersGiven(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var t1 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	body := strings.NewReader("{\"something\":\"**contents**\"}")
	req, err := http.NewRequest("PUT", "/topics/"+t1.ID, body)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsUpdateJSON)
	m.ServeHTTP(w, req)

	// DB
	var resp lib.Response
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if resp.Error != "Failed!" {
		t.Fatalf("Got %v; Expected: %v", resp.Error, "Failed!")
	}
	if w.Code != 404 {
		t.Fatalf("Got %v; Expected: %v", w.Code, 404)
	}
}

func TestDestroy(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var t1 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	req, err := http.NewRequest("POST", "/topics/"+t1.ID+"/delete", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}/delete", TopicsDestroy)
	m.ServeHTTP(w, req)

	// DB
	c, err := Db.SelectInt("select count(*) from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if c != 0 {
		t.Fatalf("Got %v; Expected: %v", c, 0)
	}

	// HTTP
	if w.Code != 302 {
		t.Fatalf("Got %v; Expected: %v", w.Code, 302)
	}
	if w.HeaderMap["Location"][0] != "/topics" {
		t.Fatalf("Got %v; Expected: %v", w.HeaderMap["Location"][0], "/topics")
	}
}

func TestDestroyJson(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	_, err := createTopic("topic")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	var t1 Topic
	err = Db.SelectOne(&t1, "select * from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}

	req, err := http.NewRequest("DELETE", "/topics/"+t1.ID, nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsDestroyJSON)
	m.ServeHTTP(w, req)

	// DB
	c, err := Db.SelectInt("select count(*) from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if c != 0 {
		t.Fatalf("Got %v; Expected: %v", c, 0)
	}

	// HTTP
	var resp lib.Response
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if resp.Message != "Ok" {
		t.Fatalf("Got %v; Expected: %v", resp.Message, "Ok")
	}
}

func TestDestroyJsonError(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	req, err := http.NewRequest("DELETE",
		"/topics/7a0a771a-cc11-4079-59ba-81df690a0588", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsDestroyJSON)
	m.ServeHTTP(w, req)

	// DB
	c, err := Db.SelectInt("select count(*) from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if c != 0 {
		t.Fatalf("Got %v; Expected: %v", c, 0)
	}

	// HTTP
	var resp lib.Response
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if resp.Error != "Could not remove topic" {
		t.Fatalf("Got %v; Expected: %v", resp.Error, "Could not remove topic")
	}
}

func TestWrongUuidFormatApi(t *testing.T) {
	initTestDB()
	defer closeTestDB()

	req, err := http.NewRequest("DELETE", "/topics/1", nil)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	m := mux.NewRouter()
	m.HandleFunc("/topics/{id}", TopicsDestroyJSON)
	m.ServeHTTP(w, req)

	// DB
	c, err := Db.SelectInt("select count(*) from topics")
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if c != 0 {
		t.Fatalf("Got %v; Expected: %v", c, 0)
	}

	// HTTP
	var resp lib.Response
	decoder := json.NewDecoder(w.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		t.Fatalf("Expected to be nil: %v", err)
	}
	if resp.Error != "Could not remove topic" {
		t.Fatalf("Got %v; Expected: %v", resp.Error, "Could not remove topic")
	}
}

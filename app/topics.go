// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

import (
	"net/http"
	"time"

	"github.com/mssola/todo/app/lib"
	"github.com/nu7hatch/gouuid"
)

func TopicsIndex(res http.ResponseWriter, req *http.Request) {
	o := &ViewData{}
	lib.Render(res, "topics/index", o)
}

func TopicsCreate(res http.ResponseWriter, req *http.Request) {
	uuid, err := uuid.NewV4()
	if err != nil {
		http.Redirect(res, req, "/topics", http.StatusFound)
		return
	}

	t := &Topic{
		Id:         uuid.String(),
		Name:       req.FormValue("name"),
		Created_at: time.Now(),
	}
	Db.Insert(t)
	http.Redirect(res, req, "/topics", http.StatusFound)
}

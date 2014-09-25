// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"fmt"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/mssola/go-utils/misc"
	"github.com/mssola/todo/app"
	"github.com/mssola/todo/app/lib"
	"github.com/mssola/todo/app/models"
)

func main() {
	// Because Martini was too mainstream :P
	n := negroni.Classic()

	// Sessions.
	lib.InitSession()

	// Database.
	models.InitDB()
	defer models.CloseDB()

	// Routing.
	r := mux.NewRouter()
	r.HandleFunc("/", app.RootIndex).Methods("GET")
	r.HandleFunc("/login", app.Login).Methods("POST")
	r.HandleFunc("/logout", app.Logout).Methods("POST").
		MatcherFunc(UserLogged)
	r.HandleFunc("/users", app.UsersCreate).Methods("POST")
	r.HandleFunc("/topics", app.TopicsIndex).Methods("GET").
		MatcherFunc(UserLogged)
	r.HandleFunc("/topics", app.TopicsCreate).Methods("POST").
		MatcherFunc(UserLogged)
	n.UseHandler(r)

	// Run, Forrest, run!
	port := fmt.Sprintf(":%v", misc.EnvOrElse("PORT", "3000"))
	n.Run(port)
}

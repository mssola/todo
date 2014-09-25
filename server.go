// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"fmt"

	"github.com/codegangsta/negroni"
	_ "github.com/lib/pq"
	"github.com/mssola/go-utils/misc"
	"github.com/mssola/todo/app/models"
	"github.com/mssola/todo/lib"
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
	r := route()
	n.UseHandler(r)

	// Run, Forrest, run!
	port := fmt.Sprintf(":%v", misc.EnvOrElse("PORT", "3000"))
	n.Run(port)
}

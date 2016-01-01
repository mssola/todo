// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"fmt"

	"github.com/codegangsta/negroni"
	"github.com/mssola/todo/app"
	"github.com/mssola/todo/lib"
)

func main() {
	// Because Martini was too mainstream :P
	n := negroni.Classic()

	// Sessions.
	lib.InitSession()

	// Database.
	app.InitDB()
	defer app.CloseDB()

	// Routing.
	r := route()
	n.UseHandler(r)

	// Run, Forrest, run!
	port := fmt.Sprintf(":%v", app.EnvOrElse("PORT", "3000"))
	n.Run(port)
}

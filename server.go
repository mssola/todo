// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/mssola/todo/app"
	"github.com/mssola/todo/lib"
)

func main() {
	// Initialize app.
	lib.InitSession()
	app.InitDB()
	defer app.CloseDB()

	// Routing.
	n := negroni.Classic()
	r := route()
	n.UseHandler(r)

	port := fmt.Sprintf(":%v", app.EnvOrElse("TODO_PORT", "3000"))

	// Try to run on HTTPS.
	cert := os.Getenv("TODO_CERT_PATH")
	key := os.Getenv("TODO_KEY_PATH")
	if cert != "" && key != "" {
		log.Printf("Running on port %s", port)
		if err := http.ListenAndServeTLS(port, cert, key, n); err != nil {
			log.Fatalf("Could not start server: %v", err)
		}
		os.Exit(0)
	}

	// Falling back to normal HTTP.
	log.Printf("Warning: this server does not use a safe connection!")
	n.Run(port)
}


# DB

This package implements a way to open DB connections. Basically, this package
opens DB connections by assuming that a JSON config file exists somewhere. This
JSON file has the following structure:

~~~ json
{
  "development": {
    "user": "mssola",
      "dbname": "project-dev",
      "password": "1234",
      "sslmode": "disable"
  },
  "production": {
    "user": "mssola",
    "dbname": "project",
    "password": "1234",
    "sslmode": "require"
  },
  "test": {
    "user": "mssola",
    "dbname": "project-test",
    "password": "1234",
    "sslmode": "disable"
  }
}
~~~

Moreover, I've also added the possibility to tell the Open function if we're
using Heroku's PostgreSQL add-on. An example:

~~~ go
package main

import (
	"github.com/mssola/go-utils/db"
	"github.com/mssola/go-utils/path"
)

func main() {
	// The Options struct has a default value for each of its fields. Take a
	// look at the documentation of the Options struct for more information.
	// The Base field is mandatory (it panics if it's not set).

	// It also panics if the JSON file could not be found.

	d := db.Open(db.Options{
		Base:        path.FindRoot("project", "."),
		Relative:    "db/database.json",
		Environment: "production",
		DBMS:        "postgres",
	})

	// Doing something useful with the d connection...

	d.Close()

	// Heroku.
	d = db.Open(db.Options{
		Base:        path.FindRoot("project", "."),
		Relative:    "db/database.json",
		Environment: "production",
		Heroku:      true,
	})

	// Doing something useful with the d connection...

	d.Close()
}
~~~

Copyright &copy; 2014 Miquel Sabaté Solà, released under the MIT License.

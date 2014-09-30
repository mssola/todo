// Copyright (C) 2014 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

// The options that can be passed to the Open function.
type Options struct {
	// An absolute path containing the base directory of the project.
	// NOTE: it panics if it's not set.
	Base string

	// Path relative to the Base path pointing to the JSON file.
	// Default: "database.json"
	Relative string

	// The name of the current environment. Default: "development".
	Environment string

	// The name of the DBMS driver. If Heroku is set to true, this option will
	// be ignored since we're always picking "postgres" in Heroku. If Heroku
	// is set to false, then it defaults to "postgres".
	DBMS string

	// Set to true if the Open function has to deal with Heroku's config.
	Heroku bool
}

// Fill the options that haven't been set with the defaults.
func mergeDefaults(options Options) Options {
	if options.Base == "" {
		panic("You have to set the Base value")
	}
	if options.Relative == "" {
		options.Relative = "database.json"
	}
	if options.Environment == "" {
		options.Environment = "development"
	}
	if options.DBMS == "" {
		options.DBMS = "postgres"
	}
	return options
}

// Returns a connection string after evaluating the given options and reading
// the JSON config file.
func configUrl(options Options) string {
	// Read the contents and unmarshal the thing.
	url := path.Join(options.Base, options.Relative)
	contents, err := ioutil.ReadFile(url)
	if err != nil {
		panic("Could not find config file!")
	}

	// Unmarshal the contents of the config file.
	m := map[string]map[string]string{}
	json.Unmarshal(contents, &m)

	// Put it in a fancy string.
	current, cfg, i := m[options.Environment], "", 0
	size := len(current)
	for k, v := range current {
		if v != "" {
			cfg += k + "=" + v
			if i != size-1 {
				cfg += " "
			}
		}
		i++
	}
	return strings.TrimSpace(cfg)
}

// This function parses the given URL as it's provided by Heroku and returns a
// string that works with Go's sql.Open function. Note that it assumes that
// PostgreSQL is the DBMS being used. The SSLMode is guessed from the PGSSL
// environment variable (defaults to "disable").
func herokuUrl(url string) string {
	// Black magic to get the PostgreSQL config.
	rg := "(?i)^postgres://(?:([^:@]+):([^@]*)@)?([^@/:]+):(\\d+)/(.*)$"
	regex := regexp.MustCompile(rg)
	matches := regex.FindStringSubmatch(url)
	if matches == nil {
		log.Fatalf("Wrong URL format!")
	}

	// In Heroku, SSLMode is stored in the PGSSL environment variable.
	sslmode := os.Getenv("PGSSL")
	if sslmode == "" {
		sslmode = "disable"
	}

	// And now we can build a proper url for PostgreSQL.
	s := "user=%s password=%s host=%s port=%s dbname=%s sslmode=%s"
	spec := fmt.Sprintf(s, matches[1], matches[2], matches[3], matches[4],
		matches[5], sslmode)
	return spec
}

func Open(options Options) *sql.DB {
	// Get the options straight.
	url := ""
	opt := mergeDefaults(options)
	if opt.Heroku {
		if h := os.Getenv("DATABASE_URL"); h != "" {
			url = herokuUrl(h)
		} else {
			url = configUrl(options)
		}
	} else {
		url = configUrl(options)
	}

	// Finally connect to the DB and return the connection.
	d, err := sql.Open(options.DBMS, url)
	if err != nil {
		panic(err)
	}
	return d
}

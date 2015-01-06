// Copyright (C) 2014-2015 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package db

import (
	"os"
	"strings"
	"testing"

	"github.com/mssola/go-utils/path"
)

func TestMergeDefaults(t *testing.T) {
	opts := Options{
		Base: "/home/mssola",
	}
	res := mergeDefaults(opts)
	if res.Relative != "database.json" {
		t.Errorf("Wrong value for Relative")
	}
	if res.Environment != "development" {
		t.Errorf("Wrong value for Environment")
	}
	if res.DBMS != "postgres" {
		t.Errorf("Wrong value for DBMS")
	}

	opts = Options{
		Base:        "/home/mssola",
		Relative:    "db/database.json",
		Environment: "production",
		DBMS:        "sqlite",
		Heroku:      true,
	}
	res = mergeDefaults(opts)
	if res.Relative != "db/database.json" {
		t.Errorf("Wrong value for Relative")
	}
	if res.Environment != "production" {
		t.Errorf("Wrong value for Environment")
	}
	if res.DBMS != "sqlite" {
		t.Errorf("Wrong value for DBMS")
	}
	if !res.Heroku {
		t.Errorf("Wrong value for Heroku")
	}
}

func mapify(str string) map[string]string {
	res := make(map[string]string)

	s := strings.Split(str, " ")
	for _, v := range s {
		pair := strings.Split(v, "=")
		res[pair[0]] =
			pair[1]
	}
	return res
}

func TestConfigUrl(t *testing.T) {
	opts := Options{
		Base:        path.FindRoot("db", "."),
		Relative:    "/test/test.json",
		Environment: "production",
	}
	d := mapify(configUrl(opts))
	if d["user"] != "postgres" || d["dbname"] != "ontop" ||
		d["sslmode"] != "require" {

		t.Errorf("Wrong value")
	}

	opts = Options{
		Base:        path.FindRoot("db", "."),
		Relative:    "test/subdir/test.json",
		Environment: "test",
	}
	d = mapify(configUrl(opts))
	if d["user"] != "postgres" || d["dbname"] != "subdir-test" ||
		d["sslmode"] != "disable" {

		t.Errorf("Wrong value")
	}
}

func TestHerokuUrl(t *testing.T) {
	url := "postgres://mssola:1234abcd@127.0.0.1:5432/kitty"
	str := herokuUrl(url)
	if str != "user=mssola password=1234abcd host=127.0.0.1 port=5432"+
		" dbname=kitty sslmode=disable" {
		t.Errorf("Wrong value for Heroku URL")
	}

	os.Setenv("PGSSL", "require")
	str = herokuUrl(url)
	if str != "user=mssola password=1234abcd host=127.0.0.1 port=5432"+
		" dbname=kitty sslmode=require" {
		t.Errorf("Wrong value for Heroku URL")
	}
}

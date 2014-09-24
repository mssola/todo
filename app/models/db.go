// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package models

import "github.com/coopernurse/gorp"

// Global instance that holds a connection to the DB. It gets initialized after
// calling the InitDB function. You have to call CloseDB in order to close the
// connection.
var Db gorp.DbMap

// TODO: test
func Count(name string) int64 {
	count, err := Db.SelectInt("select count(*) from " + name)
	if err != nil {
		return 0
	}
	return count
}

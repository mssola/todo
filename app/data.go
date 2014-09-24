// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package app

// This struct holds all the data that can be passed to a view.
type ViewData struct {
	// The id of the current user.
	Id string

	// Set to true if the current user is logged in.
	LoggedIn bool

	// Set to true if the views has to include Javascript.
	JS bool

	// Set to true if an error has happenned.
	Error bool
}

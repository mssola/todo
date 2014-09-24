// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package models

import "time"

type Topic struct {
	Id         string
	Name       string
	Created_at time.Time
}

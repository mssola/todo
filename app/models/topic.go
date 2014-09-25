// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package models

import (
	"time"

	"github.com/nu7hatch/gouuid"
)

// A topic is my way to divide different "contexts" inside my To Do list.
type Topic struct {
	Id         string
	Name       string
	Contents   string
	Created_at time.Time
}

// Given a name, try to create a new topic.
func CreateTopic(name string) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	t := &Topic{
		Id:   uuid.String(),
		Name: name,
	}
	return Db.Insert(t)
}

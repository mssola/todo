// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package models

import (
	"time"

	"github.com/nu7hatch/gouuid"
)

type Topic struct {
	Id         string
	Name       string
	Created_at time.Time
}

func CreateTopic(name string) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	t := &Topic{
		Id:         uuid.String(),
		Name:       name,
		Created_at: time.Now(),
	}
	return Db.Insert(t)
}

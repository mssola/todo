// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTopic(t *testing.T) {
	InitTestDB()
	defer CloseDB()

	// There's nothing before.
	var topic Topic
	err := Db.SelectOne(&topic, "select * from topics")
	assert.NotNil(t, err)
	assert.Empty(t, topic.Id)

	// Now we create two topics.
	err = CreateTopic("t1")
	assert.Nil(t, err)
	err = CreateTopic("t2")
	assert.Nil(t, err)

	var topics []Topic
	_, err = Db.Select(&topics, "select * from topics order by name")
	assert.NotEmpty(t, topics[0].Id)
	assert.Equal(t, topics[0].Name, "t1")
	assert.NotEmpty(t, topics[0].Created_at)
	assert.NotEmpty(t, topics[1].Id)
	assert.Equal(t, topics[1].Name, "t2")
	assert.NotEmpty(t, topics[1].Created_at)

	// We can't create a topic with an existing name.
	err = CreateTopic("t1")
	assert.NotNil(t, err)
}

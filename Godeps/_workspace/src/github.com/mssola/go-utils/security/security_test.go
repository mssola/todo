// Copyright (C) 2014-2015 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package security

import (
	"testing"
	"time"
)

func TestPassword(t *testing.T) {
	salt := PasswordSalt("1234")
	if salt == "" {
		t.Errorf("Empty string!")
	}
	if PasswordMatch(salt, "123") {
		t.Errorf("Invalid password has been matched.")
	}
	if !PasswordMatch(salt, "1234") {
		t.Errorf("It didn't match.")
	}
}

func TestNewAuthToken(t *testing.T) {
	auth := NewAuthToken()
	if auth == "" {
		t.Errorf("Empty string!")
	}
	time.Sleep(500 * time.Millisecond)
	another := NewAuthToken()
	if another == "" {
		t.Errorf("Empty string!")
	}
	if auth == another {
		t.Errorf("Didn't generate a new auth token!")
	}
}

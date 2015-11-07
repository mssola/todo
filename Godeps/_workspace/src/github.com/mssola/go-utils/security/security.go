// Copyright (C) 2014-2015 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package security

import (
	"crypto/md5"
	"encoding/base64"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Public: generate a salt for the given password.
// Returns a string containing the salted version of the given password.
func PasswordSalt(password string) string {
	pass := []byte(password)
	salt, _ := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	return string(salt)
}

// Public: check if the given password matches with the given hash.
// Returns true if the hashed version and the password match, false otherwise.
func PasswordMatch(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

// Public: Returns a new pseudo-random authentication token as a string.
func NewAuthToken() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	raw := strconv.Itoa(r.Int())
	md := md5.New()
	data := base64.StdEncoding.EncodeToString(md.Sum([]byte(raw)))
	return data
}

// Copyright (C) 2014-2015 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package misc

import (
	"os"
	"testing"
)

func TestEnvOrElse(t *testing.T) {
	os.Clearenv()
	if env := EnvOrElse("ENV", "development"); env != "development" {
		t.Errorf("Expected 'development'")
	}
	os.Setenv("ENV", "production")
	if env := EnvOrElse("ENV", "development"); env != "production" {
		t.Errorf("Expected 'production'")
	}
	os.Setenv("ENV", "")
	if env := EnvOrElse("ENV", "development"); env != "development" {
		t.Errorf("Expected 'development'")
	}
}

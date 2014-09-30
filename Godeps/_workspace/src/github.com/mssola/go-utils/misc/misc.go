// Copyright (C) 2014 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package misc

import (
	"os"
)

// Get the value of the given environment variable. If this environment
// variable is not set, then the given default value will be returned.
func EnvOrElse(name, value string) string {
	if env := os.Getenv(name); env != "" {
		return env
	}
	return value
}

// Copyright (C) 2014-2015 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package path

import "testing"

func TestFindRoot(t *testing.T) {
	abs := FindRoot("mssola", "/home/mssola/lala")
	if abs != "/home/mssola" {
		t.Errorf("Expected '/home/mssola'")
	}
	abs = FindRoot("home", "/home/mssola/lala")
	if abs != "/home" {
		t.Errorf("Expected '/home'")
	}
	abs = FindRoot("/", "/home/mssola/lala")
	if abs != "/" {
		t.Errorf("Expected '/'")
	}
}

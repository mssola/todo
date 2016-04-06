
## security [![GoDoc](https://godoc.org/github.com/mssola/go-utils/security?status.png)](http://godoc.org/github.com/mssola/go-utils/security)

This package implements a few functions regarding security.

~~~ go
package main

import (
	"fmt"

	"github.com/mssola/go-utils/security"
)

func main() {
	// PasswordSalt & PasswordMatch.
	salted := security.PasswordSalt("1234")
	fmt.Printf("%v\n", security.PasswordMatch(salted, "123"))  // => false
	fmt.Printf("%v\n", security.PasswordMatch(salted, "1234")) // => true

	// NewAuthToken. It generates a pseudo-random authentication string. This
	// string can be used to create a cookie store for safe sessions. It can
	// also be used as an authentication token.
	fmt.Printf("%v\n", security.NewAuthToken())
}
~~~

Copyright &copy; 2014-2016 Miquel Sabaté Solà, released under the MIT License.

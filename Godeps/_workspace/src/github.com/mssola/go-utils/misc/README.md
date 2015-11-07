
## misc [![GoDoc](https://godoc.org/github.com/mssola/go-utils/misc?status.png)](http://godoc.org/github.com/mssola/go-utils/misc)

This package contains random functions that I didn't know where to put. By now
there's only the EnvOrElse function. This function is similar to Python's get
but this one deals with environment variables. An example:

~~~ go
package main

import (
	"fmt"

	"github.com/mssola/go-utils/misc"
)

func main() {
	home := misc.EnvOrElse("HOME", "/home/user")
	fmt.Printf("%v\n", home) // => /home/mssola
	home = misc.EnvOrElse("HOM", "/home/user")
	fmt.Printf("%v\n", home) // => /home/user
}
~~~

Copyright &copy; 2014-2015 Miquel Sabaté Solà, released under the MIT License.

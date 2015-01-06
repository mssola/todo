
# Path

This package contains functions that deal with paths. Right now there's only
the FindRoot function. An example:

~~~ go
package main

import (
	"fmt"

	"github.com/mssola/go-utils/path"
)

func main() {
	base := path.FindRoot("mssola", "/home/mssola/another/mssola/dir")
	fmt.Printf("%v\n", base) // => /home/mssola/another/mssola

	// We can also specify the current path with a ".".
	// In this example image that the current pwd is "/home/mssola/test"
	base = path.FindRoot("mssola", ".")
	fmt.Printf("%v\n", base) // => /home/mssola
}
~~~

Copyright &copy; 2014-2015 Miquel Sabaté Solà, released under the MIT License.

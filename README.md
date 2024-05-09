# Goodreads Parser

This is a library with basic functionality to parse goodreads website for books. Since API is going to shut down and they do not give API keys anymore.

Example of usage:

```
package main

import goodreads_parser "github.com/sonac/goodreads_parser/api"
import "fmt"

func main() {
	prs := goodreads_parser.NewParser()
	bks, err := prs.FindBooks("harry potter", 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", bks)
}
```

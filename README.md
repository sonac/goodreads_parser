# Goodreads Parser

This is a library with basic functionality to parse goodreads website for books. Since API is going to shut down and they do not give API keys anymore.

## Examples

### Basic usage

See example of basic usage in `main.go`. You can run run this locally to check out via `go run main.go`

To use as a library:

```aiignore
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
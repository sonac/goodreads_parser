# Goodreads Parser

This is a library with basic functionality to parse goodreads website for books. Since API is going to shut down and they do not give API keys anymore.

## Examples

### Basic usage

```go
package main

import goodreads_parser "github.com/sonac/goodreads_parser/api"
import "fmt"

func main() {
  parser := goodreads_parser.NewParser()
  options := goodreads_parser.DefaultBookOptions() // Use sensible defaults

  books, err := parser.FindBooks("harry potter", 5, options)
  if err != nil {
    fmt.Println("Error:", err)
    return
  }

  fmt.Printf("Found %d books\n", len(books))
  for i, book := range books {
    fmt.Printf("%d. %s by %s (%.1f stars, %d ratings)\n", 
              i+1, book.Title, book.Author, book.Rating.Avg, book.Rating.Count)
  }
}
```

### With high quality filters

```go
package main

import goodreads_parser "github.com/sonac/goodreads_parser/api"
import "fmt"

func main() {
  parser := goodreads_parser.NewParser()

  // Use built-in high quality options (1000+ ratings, 4.0+ avg rating)
  options := goodreads_parser.NewHighQualityBookOptions()

  books, err := parser.FindBooks("science fiction", 10, options)
  if err != nil {
    fmt.Println("Error:", err)
    return
  }

  fmt.Printf("Found %d high-quality books\n", len(books))
  for i, book := range books {
    fmt.Printf("%d. %s by %s (%.1f stars, %d ratings)\n", 
              i+1, book.Title, book.Author, book.Rating.Avg, book.Rating.Count)
  }
}
```

### With no filters

```go
package main

import goodreads_parser "github.com/sonac/goodreads_parser/api"
import "fmt"

func main() {
  parser := goodreads_parser.NewParser()
  options := goodreads_parser.NewBookOptionsNoFilters() // No quality filtering

  books, err := parser.FindBooks("harry potter", 5, options)
  if err != nil {
    fmt.Println("Error:", err)
    return
  }

  fmt.Printf("%+v", books)
}
```

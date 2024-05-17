# Distinct

This micro-library estimates the number of distinct elements in a stream:

```go
import (
    "fmt"
    "github.com/betamos/distinct"
)

// Uses a map of size 200
c := distinct.NewCounter[int](200, nil)
for i := 0; i < 100000; i++ {
	c.Add(i % 300)
}
fmt.Println(c.Estimate()) // Should be somewhere around 300
```

It's based on [Distinct Elements in Streams: An Algorithm for the (Text) Book](https://arxiv.org/pdf/2301.10191).

# Simple Cache

With this library you can easily add thread safe LRU Cache to your service

## Example

```go
package main

import "github.com/offluck/slru"

var (
	LoginAttempts = slru.NewCache[string, int](0)
)

func main() {
	LoginAttempts.Push(slru.NewItem("someone@mail.com", 0))
	attempts, ok := LoginAttempts.Get("someone@mail.com")
	if !ok {
		// Error
	}

	err := LoginAttempts.Set("someone@mail.com", attempts+1)
	if err != nil {
		// Error
	}
}

```

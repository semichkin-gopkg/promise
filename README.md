# Library for async job waiting

```go
package main

import (
	"context"
	"errors"
	"github.com/semichkin-gopkg/promise"
	"log"
)

func asyncDivide(a, b int) *promise.Promise[int] {
	p := promise.NewPromise[int](
		promise.WithRejectHandler[int](func(err error) {
			log.Println("logging: ", err)
		}),
	)

	go func() {
		if b == 0 {
			p.Reject(errors.New("division by zero"))
			return
		}

		p.Resolve(a / b)
	}()

	return p
}

func main() {
	log.Println(asyncDivide(1, 1).Wait(context.Background()).Payload) // 1
	log.Println(asyncDivide(1, 0).Wait(context.Background()).Error)   // division by zero
}
```